package main

import (
	"context"
	"net/http"
	"time"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// BuildInfo app build version, should be set from build phase
// go build -ldflags "-X main.BuildInfo=${BUILD_TIME}.${GIT_COMMIT}" -o ./out/echo-restapi example/echo-restapi/main.go.
var BuildInfo = "NOT_SET"

const (
	dPort           = 8088
	dWaitDur        = 3 * time.Second
	dTimeoutDur     = 4 * time.Second
	reqTimeout      = 5 * time.Second
	longProcessTime = 10 * time.Second
	isProdMode      = false
)

func main() {
	// appname should be snake_cased
	appName := "example_echo_restapi"

	// setup logging
	if isProdMode {
		// production mode - json
		_ = log.NewLogger(log.LevelDebug, appName, nil, nil, "default_key1", "default_value1").Set()
	} else {
		// development mode - logfmt
		_ = log.NewDevLogger(log.LevelDebug, appName, nil, nil, "default_key1", "default_value1").Set()
	}

	cfg := echokit.RuntimeConfig{
		Port:                    dPort,
		ShutdownWaitDuration:    dWaitDur,
		ShutdownTimeoutDuration: dTimeoutDur,
		HealthCheckFunc:         healthCheck(),
		Name:                    appName,
		BuildInfo:               BuildInfo,
		RequestTimeoutConfig: &echokit.TimeoutConfig{
			Timeout: reqTimeout,
		},
	}

	e := echo.New()
	e.Use(
		middleware.Recover(),
		echokit.RequestIDLoggerMiddleware(&cfg),
		echokit.BodyDumpHandler(func(ctx echo.Context) bool {
			path := ctx.Request().URL.Path
			skippedPaths := []string{"/health", "/metrics"}

			for _, v := range skippedPaths {
				if path == v {
					return true
				}
			}

			return false
		}),
	)

	// routes
	e.POST("/users", createUser)
	e.GET("/longsleep", longOperation)

	echokit.RunServer(e, &cfg)
}

func healthCheck() echokit.HealthCheckFunc {
	return func(ctx context.Context) error {
		return nil
	}
}

func createUser(ctx echo.Context) error {
	type User struct {
		Name string `validate:"required" json:"name"`
	}

	var u User

	if err := ctx.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusAccepted, u)
}

func longOperation(ctx echo.Context) error {
	rCtx := ctx.Request().Context()

	now := time.Now()

	defer func() {
		log.FromCtx(rCtx).Debug("request ends", "elapsed_time_ms", time.Since(now).Milliseconds())
	}()

	// usually, we call external system with context as parameter
	// e.g. db.QueryRow(ctx, ...), redis.WithContext(ctx).Get(key).Result(), etc..
	// this example using time.Sleep to simulate long running process
	// since it doesn't support ctx param, we use select-case
	select {
	case <-rCtx.Done():
		return echo.NewHTTPError(http.StatusRequestTimeout, "request timeout")
	case <-time.After(longProcessTime):
		return ctx.String(http.StatusOK, "OK. Done.")
	}
}
