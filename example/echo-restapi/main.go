package main

import (
	"context"
	"time"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/labstack/echo/v4"
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
	cfg := echokit.RuntimeConfig{
		Port:                    dPort,
		ShutdownWaitDuration:    dWaitDur,
		ShutdownTimeoutDuration: dTimeoutDur,
		HealthCheckFunc:         healthCheck(),
		Name:                    "example-echo-restapi",
		BuildInfo:               BuildInfo,
	}

	e := echo.New()

	echokit.RunServer(e, &cfg)
}

func healthCheck() echokit.HealthCheckFunc {
	return func(ctx context.Context) error {
		return nil
	}
}
