package main

import (
	"context"
	"os"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/pinpoint-apm/pinpoint-go-agent"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/echokit/echoapmkit"
	"github.com/adipurnama/go-toolkit/examples/echo-restapi/internal/handler"
	"github.com/adipurnama/go-toolkit/examples/echo-restapi/internal/repository"
	"github.com/adipurnama/go-toolkit/examples/echo-restapi/internal/service"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/pinpointkit"
	"github.com/adipurnama/go-toolkit/tracer"
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

	var pAgent pinpoint.Agent

	pHost, set := os.LookupEnv("PINPOINT_HOST")
	if set {
		pOpts := pinpointkit.WithOptions(pinpointkit.Options{
			AppName: appName,
			Env:     "",
			Host:    pHost,
		})

		pagent, err := pinpointkit.NewAgent(pOpts)
		if err != nil {
			log.FromCtx(context.Background()).Error(err, "failed call pinpoint agent")
		}

		pAgent = pagent
	}

	cfg := echokit.RuntimeConfig{
		Port:                    dPort,
		ShutdownWaitDuration:    dWaitDur,
		ShutdownTimeoutDuration: dTimeoutDur,
		HealthCheckPath:         "/healthz",
		HealthCheckFunc:         handler.HealthCheck(&db),
		Name:                    appName,
		BuildInfo:               BuildInfo,
		RequestTimeoutConfig: &echokit.TimeoutConfig{
			Timeout: reqTimeout,
		},
	}

	e := echo.New()
	e.HTTPErrorHandler = handler.ErrorHandler

	apmOpts := []echoapmkit.APMOption{}
	if pAgent != nil {
		apmOpts = append(apmOpts, echoapmkit.WithPinpointAgent(pAgent))
		tracer.Setup(tracer.WithPinpoint())
	}

	e.Use(
		echoapmkit.RecoverMiddleware(apmOpts...),
		echokit.RequestIDLoggerMiddleware(&cfg),
	)

	if !isProdMode {
		e.Use(
			echokit.BodyDumpHandler(func(ctx echo.Context) bool {
				path := ctx.Request().URL.Path
				skippedPaths := []string{"/healthz", "/metrics"}

				for _, v := range skippedPaths {
					if path == v {
						return true
					}
				}

				return false
			}),
		)
	}

	// routes
	e.POST("/users", handler.CreateUser(svc))
	e.GET("/users/:id", handler.GetUser(svc))
	e.GET("/longsleep", handler.LongOperation)
	e.GET("/panic", handler.PanicGuaranteed)
	e.GET("/error/:code", handler.EmitError)

	// run echo-http server
	echokit.RunServer(e, &cfg)
}
