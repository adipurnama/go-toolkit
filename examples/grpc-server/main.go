package main

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	pgrpc "github.com/pinpoint-apm/pinpoint-go-agent/plugin/grpc"
	"google.golang.org/grpc"

	v1 "github.com/adipurnama/go-toolkit/examples/grpc-server/v1"
	"github.com/adipurnama/go-toolkit/grpckit"
	"github.com/adipurnama/go-toolkit/grpckit/grpcapmkit"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/pinpointkit"
	"github.com/adipurnama/go-toolkit/springcloud"
	"github.com/adipurnama/go-toolkit/web/httpclient"
)

const (
	port       = 30331
	wait       = 3 * time.Second
	timeout    = 5 * time.Second
	isProdMode = false
)

func main() {
	appName := "example_grpc_app"

	// setup logging
	if isProdMode {
		// production mode - json
		_ = log.NewLogger(log.LevelDebug, appName, nil, nil, "default_key1", "default_value1").Set()
	} else {
		// development mode - logfmt
		_ = log.NewDevLogger(nil, nil, "default_key1", "default_value1").Set()
	}

	cfg := grpckit.RuntimeConfig{
		Port:                 port,
		ShutdownWaitDuration: wait,
		Name:                 appName,
		EnableReflection:     true,
		HealthCheckFunc:      healthCheck(),
	}

	httpClient := httpclient.NewStdHTTPClient()
	springConfig := springcloud.NewRemoteConfig(httpClient)

	popt := pinpointkit.WithOptionsFromConfig(springConfig, "app.pinpoint")

	pagent, err := pinpointkit.NewAgent(popt)
	if err != nil {
		log.FromCtx(context.Background()).Error(err, "failed call pinpoint agent")
	}

	sOpts := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpcapmkit.NewUnaryServerInterceptor(grpcapmkit.WithPinpointAgent(pagent)),
			grpckit.RequestTimeoutInterceptor(timeout),
			grpckit.RequestIDInterceptor(grpckit.DefaultRequestIDProvider()),
			grpckit.ErrorResponseWriterInterceptor(grpckit.DefaultGRPCErrorHandler),
			grpckit.LoggerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
		))

	s := grpc.NewServer(
		sOpts,
		grpc.StreamInterceptor(pgrpc.StreamServerInterceptor(pagent)),
	)

	exampleSvc := &v1.Server{}

	v1.RegisterExampleServiceServer(s, exampleSvc)

	grpckit.Run(s, &cfg)
}

// if necesssary do complex checking here
// check db / cache.. call api, etc..
func healthCheck() grpckit.HealthCheckFunc {
	return func(ctx context.Context) error {
		return nil
	}
}
