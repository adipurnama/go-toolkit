package echokit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/adipurnama/go-toolkit/web"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/iancoleman/strcase"
	echo_prometheus "github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
)

// InjectedPaths is auto-registered paths for observability.
var InjectedPaths = []string{"/actuator/info", "/actuator/health", "/metrics"}

// RuntimeConfig restapi runtime config with healthcheck
type RuntimeConfig struct {
	ShutdownWaitDuration    time.Duration
	ShutdownTimeoutDuration time.Duration
	Port                    int
	Name                    string
	BuildInfo               string
	HealthCheckFunc
}

type healthStatus struct {
	serving bool
	Status  string `json:"status"`
}

// HealthCheckFunc is healthcheck interface func
type HealthCheckFunc func(ctx context.Context) error

// RunServer run graceful restapi server
func RunServer(e *echo.Echo, cfg *RuntimeConfig) {
	appCtx, done := web.NewRuntimeContext()
	defer done()

	RunServerWithContext(appCtx, e, cfg)
}

// RunServerWithContext run graceful restapi server with existing background context
func RunServerWithContext(appCtx context.Context, e *echo.Echo, cfg *RuntimeConfig) {
	cfg.Name = strcase.ToSnake(cfg.Name)

	logger := log.FromCtx(appCtx)
	logger.AddField("restapi_name", cfg.Name)

	hs := &healthStatus{
		serving: true,
	}

	e.HideBanner = true
	e.GET("/actuator/health", func(c echo.Context) error {
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
	e.GET("/actuator/info", func(c echo.Context) error {
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

	logger.Info("serving REST HTTP server", "port", cfg.Port)

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err, "starting http server")
	}
}
