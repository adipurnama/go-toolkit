package main

import (
	"context"
	"time"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// BuildInfo app build version, should be set from build phase
// go build -ldflags "-X main.BuildInfo=${BUILD_TIME}.${GIT_COMMIT}" -o ./out/echo-restapi example/echo-restapi/main.go
var BuildInfo = "NOT_SET"

const (
	dPort       = 8088
	dWaitDur    = 3 * time.Second
	dTimeoutDur = 4 * time.Second
)

func main() {
	appName := "example-echo-restapi"

	// setup logging
	// development mode - logfmt
	_ = log.NewDevLogger(log.LevelDebug, appName, nil, nil, "default_key1", "default_value1").Set()
	// production mode - json
	// _ = log.NewLogger(log.LevelDebug, appName, nil, nil, "default_key1", "default_value1").Set()

	cfg := echokit.RuntimeConfig{
		Port:                    dPort,
		ShutdownWaitDuration:    dWaitDur,
		ShutdownTimeoutDuration: dTimeoutDur,
		HealthCheckFunc:         healthCheck(),
		Name:                    appName,
		BuildInfo:               BuildInfo,
	}

	e := echo.New()
	e.Use(echokit.RequestIDLoggerMiddleware(), echokit.BodyDumpHandler(middleware.DefaultSkipper))

	echokit.RunServer(e, &cfg)
}

func healthCheck() echokit.HealthCheckFunc {
	return func(ctx context.Context) error {
		return nil
	}
}
