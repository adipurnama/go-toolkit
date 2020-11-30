package main

import (
	"context"
	"time"

	v1 "github.com/adipurnama/go-toolkit/example/grpc-server/v1"
	"github.com/adipurnama/go-toolkit/grpckit"
	"github.com/adipurnama/go-toolkit/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

const (
	port  = 30331
	waitS = 3 * time.Second
)

func main() {
	appName := "example-grpc-app"

	// setup logging
	// development mode - logfmt
	_ = log.NewDevLogger(log.LevelDebug, appName, nil, nil, "default_key1", "default_value1").Set()
	// production mode - json
	// _ = log.NewLogger(log.LevelDebug, appName, nil, nil, "default_key1", "default_value1").Set()

	cfg := grpckit.RuntimeConfig{
		Port:                 port,
		ShutdownWaitDuration: waitS,
		Name:                 appName,
		EnableReflection:     true,
		HealthCheckFunc:      healthCheck(),
	}

	uIntOpt := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpckit.LoggerInterceptor(),
		))

	s := grpc.NewServer(uIntOpt)

	exampleSvc := &v1.Service{}

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
