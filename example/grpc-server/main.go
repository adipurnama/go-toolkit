package main

import (
	"context"
	"time"

	v1 "github.com/adipurnama/go-toolkit/example/grpc-server/v1"
	"github.com/adipurnama/go-toolkit/grpckit"
	"github.com/adipurnama/go-toolkit/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
)

const (
	port       = 30331
	wait       = 3 * time.Second
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

	sOpts := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
			grpckit.RequestIDInterceptor(grpckit.DefaultRequestIDProvider()),
			grpckit.ErrorResponseWriterInterceptor(grpckit.DefaultGRPCErrorHandler),
			grpckit.LoggerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
		))

	s := grpc.NewServer(sOpts)

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
