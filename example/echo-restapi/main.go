package main

import (
	"time"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/example/echo-restapi/internal/handler"
	"github.com/adipurnama/go-toolkit/example/echo-restapi/internal/repository"
	"github.com/adipurnama/go-toolkit/example/echo-restapi/internal/service"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.elastic.co/apm/module/apmechov4"
)

// BuildInfo app build version, should be set from build phase
// go build -ldflags "-X main.BuildInfo=${BUILD_TIME}.${GIT_COMMIT}" -o ./out/echo-restapi example/echo-restapi/main.go.
var BuildInfo = "NOT_SET"

const (
	dPort       = 8088
	dWaitDur    = 1 * time.Second
	dTimeoutDur = 4 * time.Second
	reqTimeout  = 15 * time.Second
	isProdMode  = false
)

func main() {
	// appname should be snake_cased
	appName := "example_echo_restapi"

	// setup logging
	if isProdMode {
		// production mode - json
		_ = log.NewLogger(log.LevelDebug, appName, nil, nil, "additionalKey1", "additional_value1").Set()
	} else {
		// development mode - logfmt
		_ = log.NewDevLogger(nil, nil, "default_key1", "default_value1").Set()
	}

	// dependencies
	db := repository.DummyDB{}
	repo := repository.NewUserRepository(&db)
	svc := service.NewService(repo)

	cfg := echokit.RuntimeConfig{
		Port:                    dPort,
		ShutdownWaitDuration:    dWaitDur,
		ShutdownTimeoutDuration: dTimeoutDur,
		HealthCheckFunc:         handler.HealthCheck(&db),
		Name:                    appName,
		BuildInfo:               BuildInfo,
		RequestTimeoutConfig: &echokit.TimeoutConfig{
			Timeout: reqTimeout,
		},
	}

	e := echo.New()
	e.HTTPErrorHandler = handler.ErrorHandler
	e.Use(
		middleware.Recover(),
		echokit.RequestIDLoggerMiddleware(&cfg),
		apmechov4.Middleware(),
		// echokit.BodyDumpHandler(func(ctx echo.Context) bool {
		// 	path := ctx.Request().URL.Path
		// 	skippedPaths := []string{"/health", "/metrics"}

		// 	for _, v := range skippedPaths {
		// 		if path == v {
		// 			return true
		// 		}
		// 	}

		// 	return false
		// }),
	)

	// routes
	e.POST("/users", handler.CreateUser(svc))
	e.GET("/users/:id", handler.GetUser(svc))
	e.GET("/longsleep", handler.LongOperation)

	// run echo-http server
	echokit.RunServer(e, &cfg)
}
