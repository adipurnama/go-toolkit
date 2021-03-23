package echokit

import (
	"context"
	"errors"
	"fmt"
	stdLog "log"
	"net/http"
	"time"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/web"
	validator "github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	echo_prometheus "github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	defaultInfoPath   = "/actuator/info"
	defaultHealthPath = "/actuator/health"
	defaultReqTimeout = 7 * time.Second
)

// RuntimeConfig defines echo REST API runtime config with healthcheck.
type RuntimeConfig struct {
	Port                    int
	Name                    string
	BuildInfo               string
	ShutdownWaitDuration    time.Duration
	ShutdownTimeoutDuration time.Duration
	RequestTimeoutConfig    *TimeoutConfig
	HealthCheckPath         string
	InfoCheckPath           string
	HealthCheckFunc
}

type healthStatus struct {
	serving bool
	Status  string `json:"status"`
}

// HealthCheckFunc is healthcheck interface func.
type HealthCheckFunc func(ctx context.Context) error

// RunServer run graceful restapi server.
func RunServer(e *echo.Echo, cfg *RuntimeConfig) {
	appCtx, done := web.NewRuntimeContext()
	defer done()

	RunServerWithContext(appCtx, e, cfg)
}

// RunServerWithContext run graceful restapi server with existing background context
// provides default '/actuator/health' as healthcheck endpoint
// provides '/metrics' as prometheus metrics endpoint.
// set echo.Validator using `web.Validator` from `web` package.
func RunServerWithContext(appCtx context.Context, e *echo.Echo, cfg *RuntimeConfig) {
	cfg.Name = strcase.ToSnake(cfg.Name)

	logger := log.FromCtx(appCtx)

	hs := &healthStatus{
		serving: true,
	}

	e.HideBanner = true
	validator := validator.New()

	if cfg.RequestTimeoutConfig == nil {
		cfg.RequestTimeoutConfig = &TimeoutConfig{
			Timeout: defaultReqTimeout,
			Skipper: middleware.DefaultSkipper,
		}
	}

	// request validator setup
	e.Use(ValidatorTranslatorMiddleware(validator), TimeoutMiddleware(cfg.RequestTimeoutConfig))
	e.Validator = web.NewValidator(validator)

	if cfg.HealthCheckPath == "" {
		cfg.HealthCheckPath = defaultHealthPath
	}

	e.GET(cfg.HealthCheckPath, func(c echo.Context) error {
		if !hs.serving {
			hs.Status = "OUT_OF_SERVICE"

			return c.JSON(http.StatusOK, hs)
		}

		if cfg.HealthCheckFunc == nil {
			hs.Status = "UP"

			return c.JSON(http.StatusOK, hs)
		}

		err := cfg.HealthCheckFunc(c.Request().Context())
		if err != nil {
			hs.Status = "OUT_OF_SERVICE"

			return c.JSON(http.StatusOK, hs)
		}

		hs.Status = "UP"

		return c.JSON(http.StatusOK, hs)
	})

	if cfg.InfoCheckPath == "" {
		cfg.InfoCheckPath = defaultInfoPath
	}

	e.GET(cfg.InfoCheckPath, func(c echo.Context) error {
		var v struct {
			Version string `json:"version"`
		}
		v.Version = cfg.BuildInfo

		return c.JSON(http.StatusOK, v)
	})

	// prometheus
	p := echo_prometheus.NewPrometheus(cfg.Name, nil)
	p.Use(e)

	go func() {
		<-appCtx.Done()

		hs.serving = false

		logger.Info(fmt.Sprintf("shutting down REST HTTP server in %d ms", cfg.ShutdownWaitDuration.Milliseconds()))
		<-time.After(cfg.ShutdownWaitDuration)

		// stop the server
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeoutDuration)
		defer cancel()

		if err := e.Shutdown(shutdownCtx); err != nil {
			logger.Error(err, "shutdown http server")
		}
	}()

	// error fallback handler
	e.HTTPErrorHandler = loggerHTTPErrorHandler(e.HTTPErrorHandler)

	PrintRoutes(e)

	// start server
	logger.Info("serving REST HTTP server", "port", cfg.Port)

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err, "starting http server")
	}
}

// PrintRoutes prints *echo.Echo routes.
func PrintRoutes(e *echo.Echo) {
	stdLog.Println("=== initializing http routes")

	for _, r := range e.Routes() {
		stdLog.Printf("===> %s %s %s", r.Method, r.Path, r.Name)
	}
}
