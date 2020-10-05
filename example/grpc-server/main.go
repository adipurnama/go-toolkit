package main

import (
	"context"
	"time"

	"github.com/adipurnama/go-toolkit/grpckit"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

const (
	port  = 30331
	waitS = 3 * time.Second
)

func main() {
	cfg := grpckit.RuntimeConfig{
		Port:                 port,
		ShutdownWaitDuration: waitS,
		Name:                 "example-grpc",
		EnableReflection:     true,
		HealthCheckFunc:      healthCheck(),
	}

	uIntOpt := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor()),
	)

	s := grpc.NewServer(uIntOpt)

	grpckit.Run(s, &cfg)
}

// if necesssary do complex checking here
// check db / cache.. call api, etc..
func healthCheck() grpckit.HealthCheckFunc {
	return func(ctx context.Context) error {
		return nil
	}
}
