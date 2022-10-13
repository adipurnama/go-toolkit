package main

import (
	"context"
	"net/http"
	"os"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/labstack/echo/v4"
	"github.com/pinpoint-apm/pinpoint-go-agent"
	"google.golang.org/grpc"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/echokit/echoapmkit"
	v1 "github.com/adipurnama/go-toolkit/examples/grpc-server/v1"
	"github.com/adipurnama/go-toolkit/grpckit"
	"github.com/adipurnama/go-toolkit/grpckit/grpcapmkit"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/pinpointkit"
	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/adipurnama/go-toolkit/springcloud"
	"github.com/adipurnama/go-toolkit/tracer"
	"github.com/adipurnama/go-toolkit/web/httpclient"
)

const (
	restPort = 8082
	port     = 30331
	wait     = 3 * time.Second
	timeout  = 5 * time.Second
)

func main() {
	appContext, cancel := runtimekit.NewRuntimeContext()
	defer cancel()

	appName := "example_echo_grpc"

	// gRPC server
	cfg := grpckit.RuntimeConfig{
		Port:                 port,
		ShutdownWaitDuration: wait,
		Name:                 appName,
		EnableReflection:     true,
		HealthCheckFunc:      healthCheck(),
	}

	httpClient := httpclient.NewStdHTTPClient()
	springConfig := springcloud.NewRemoteConfig(httpClient)

	err := springConfig.Load(appContext)
	if err != nil {
		return
	}

	var pAgent pinpoint.Agent

	pHost, ok := os.LookupEnv("PINPONT_HOST")
	if ok {
		pOpts := pinpointkit.WithOptions(pinpointkit.Options{
			AppName: appName,
			Env:     "test",
			Host:    pHost,
		})

		pagent, err := pinpointkit.NewAgent(pOpts)
		if err != nil {
			log.FromCtx(appContext).Error(err, "failed call pinpoint agent")
		}
		pAgent = pagent
	}

	tracer.Setup(tracer.WithPinpoint())

	grpcServerOpts := []grpcapmkit.ServerOption{}
	if pAgent != nil {
		grpcServerOpts = append(grpcServerOpts, grpcapmkit.WithPinpointAgent(pAgent))
	}

	sOpts := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpcapmkit.NewUnaryServerInterceptor(grpcServerOpts...),
			grpckit.RequestTimeoutInterceptor(timeout),
			grpckit.RequestIDInterceptor(grpckit.DefaultRequestIDProvider()),
			grpckit.ErrorResponseWriterInterceptor(grpckit.DefaultGRPCErrorHandler),
			grpckit.LoggerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
		))

	s := grpc.NewServer(
		sOpts,
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpcapmkit.NewStreamServerInterceptor(grpcServerOpts...),
			),
		),
	)

	exampleSvc := &v1.Server{}

	v1.RegisterExampleServiceServer(s, exampleSvc)

	// run as separate go routine
	go grpckit.RunWithContext(appContext, s, &cfg)

	// rest server
	restCfg := echokit.RuntimeConfig{
		Port:                    restPort,
		ShutdownWaitDuration:    wait,
		ShutdownTimeoutDuration: timeout,
		HealthCheckFunc: func(c context.Context) error {
			return nil
		},
	}

	apmOpts := []echoapmkit.APMOption{}
	if pAgent != nil {
		apmOpts = append(apmOpts, echoapmkit.WithPinpointAgent(pAgent))
	}

	e := echo.New()
	e.Use(echoapmkit.RecoverMiddleware(apmOpts...))

	e.GET("/hello", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "HELLO WORLD")
	})

	e.GET("/crash", crash())

	// IMPORTANT: run echo server with same context as gRPC server
	echokit.RunServerWithContext(appContext, e, &restCfg)
}

// if necesssary do complex checking here
// check db / cache.. call api, etc..
func healthCheck() grpckit.HealthCheckFunc {
	return func(ctx context.Context) error {
		return nil
	}
}

func crash() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		_, span := tracer.NewSpan(ctx.Request().Context(), tracer.SpanLvlHandler)
		defer span.End()

		zero := 0
		ten := 10

		divided := ten / zero

		return ctx.JSON(http.StatusInternalServerError, divided)
	}
}
