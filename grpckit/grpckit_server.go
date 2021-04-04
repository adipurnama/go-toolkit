package grpckit

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/iancoleman/strcase"

	"github.com/adipurnama/go-toolkit/grpckit/grpc_health_v1"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/runtimekit"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RuntimeConfig defines runtime configuration for grpc service with health check.
type RuntimeConfig struct {
	ShutdownWaitDuration time.Duration
	Port                 int
	Name                 string
	EnableReflection     bool
	HealthCheckFunc
}

// Run grpc server with health check, creating new app context.
func Run(s *grpc.Server, cfg *RuntimeConfig) {
	appCtx, done := runtimekit.NewRuntimeContext()
	defer done()

	RunWithContext(appCtx, s, cfg)
}

// RunWithContext runs grpc server with health check using existing background context.
func RunWithContext(appCtx context.Context, s *grpc.Server, cfg *RuntimeConfig) {
	cfg.Name = strcase.ToSnake(cfg.Name)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.FromCtx(appCtx).Error(err, "net.Listen", "grpc_app_name", cfg.Name)
		return
	}

	hs := NewHealthcheckServer(cfg.HealthCheckFunc)
	grpc_health_v1.RegisterHealthServer(s, hs)

	if cfg.EnableReflection {
		reflection.Register(s)
	}

	log.FromCtx(appCtx).Info("serving gRPC service", "port", cfg.Port, "grpc_app_name", cfg.Name)

	go func() {
		<-appCtx.Done()

		hs.Serving = false

		log.FromCtx(appCtx).Info(fmt.Sprintf("shutting down gRPC server in %d ms...", cfg.ShutdownWaitDuration.Milliseconds()))
		<-time.After(cfg.ShutdownWaitDuration)

		s.GracefulStop()
	}()

	if err := s.Serve(lis); err != nil {
		log.FromCtx(appCtx).Error(err, "s.Serve", "grpc_app_name", cfg.Name)
	}
}
