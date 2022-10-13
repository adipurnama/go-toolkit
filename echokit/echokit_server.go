package echokit

import (
	"context"
	"errors"
	"fmt"
	stdLog "log"
	"net/http"
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	echo_prometheus "github.com/labstack/echo-contrib/prometheus"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/adipurnama/go-toolkit/web"
)

var errInvalidHealthCheckFunc = errors.New("echokit_server: valid HealthCheckFunc is required")

const (
	defaultInfoPath   = "/actuator/info"
	defaultHealthPath = "/actuator/health"
	defaultReqTimeout = 7 * time.Second
	defaultPort       = 8088
)

// RuntimeConfig defines echo REST API runtime config with healthcheck.
type RuntimeConfig struct {
	Port                    int            `json:"port,omitempty"`
	Name                    string         `json:"name,omitempty"`
	BuildInfo               string         `json:"build_info,omitempty"`
	ShutdownWaitDuration    time.Duration  `json:"shutdown_wait_duration,omitempty"`
	ShutdownTimeoutDuration time.Duration  `json:"shutdown_timeout_duration,omitempty"`
	RequestTimeoutConfig    *TimeoutConfig `json:"request_timeout_config,omitempty"`
	HealthCheckPath         string         `json:"health_check_path,omitempty"`
	InfoCheckPath           string         `json:"info_check_path,omitempty"`
	HealthCheckFunc         `json:"-"`
}

func (cfg *RuntimeConfig) validate() {
	// port
	if cfg.Port == 0 {
		cfg.Port = defaultPort
	}

	// healthcheck
	if cfg.HealthCheckPath == "" {
		cfg.HealthCheckPath = defaultHealthPath
	}

	// check for timeout setting
	if cfg.RequestTimeoutConfig == nil {
		cfg.RequestTimeoutConfig = &TimeoutConfig{}
	}

	if cfg.RequestTimeoutConfig.Timeout == 0 {
		cfg.RequestTimeoutConfig.Timeout = defaultReqTimeout
	}

	if cfg.RequestTimeoutConfig.Skipper == nil {
		cfg.RequestTimeoutConfig.Skipper = middleware.DefaultSkipper
	}
}

type healthStatus struct {
	serving bool
	Status  string `json:"status"`
}

// HealthCheckFunc is healthcheck interface func.
type HealthCheckFunc func(ctx context.Context) error

// RunServer run graceful restapi server.
func RunServer(e *echo.Echo, cfg *RuntimeConfig) {
	appCtx, done := runtimekit.NewRuntimeContext()
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

	e.HideBanner = true
	validator := validator.New()

	cfg.validate()

	// request validator setup
	e.Use(ValidatorTranslatorMiddleware(validator), TimeoutMiddleware(cfg.RequestTimeoutConfig))
	e.Validator = web.NewValidator(validator)

	if cfg.HealthCheckFunc == nil {
		log.FromCtx(appCtx).Error(errInvalidHealthCheckFunc, "please provide healthcheck function to runtime config")
		return
	}

	// healthcheck
	hs := &healthStatus{
		serving: true,
	}

	e.GET(cfg.HealthCheckPath, func(c echo.Context) error {
		if !hs.serving {
			hs.Status = "OUT_OF_SERVICE"

			return c.JSON(http.StatusServiceUnavailable, hs)
		}

		err := cfg.HealthCheckFunc(c.Request().Context())
		if err != nil {
			hs.Status = "OUT_OF_SERVICE"

			return c.JSON(http.StatusServiceUnavailable, hs)
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
	originalErrHandler := e.HTTPErrorHandler
	e.HTTPErrorHandler = loggerHTTPErrorHandler(func(err error, respCtx echo.Context) {
		var errEcho *echo.HTTPError

		// for `echo.HTTPError`, let echo server directly handles it
		// i.e. returns message & status code from it
		if errors.As(err, &errEcho) {
			return
		}

		originalErrHandler(err, respCtx)
	})

	PrintRoutes(e)

	// start server
	logger.Info("serving REST HTTP server", "config", cfg)

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err, "starting http server")
	}
}

// PrintRoutes prints *echo.Echo routes.
func PrintRoutes(e *echo.Echo) {
	stdLog.Println("== initializing http routes")

	for _, r := range e.Routes() {
		handlerNames := strings.Split(r.Name, "/")
		stdLog.Printf("=====> %s %s %s", r.Method, r.Path, handlerNames[len(handlerNames)-1:][0])
	}
}
