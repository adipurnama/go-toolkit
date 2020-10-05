package grpckit

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/adipurnama/go-toolkit/grpckit/grpc_health_v1"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/sethvargo/go-signalcontext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RuntimeConfig is runtime configuration for grpc service with healthcheck
type RuntimeConfig struct {
	ShutdownWaitDuration time.Duration
	Port                 int
	Name                 string
	EnableReflection     bool
	HealthCheckFunc
}

// Run grpc server with healthcheck, creating new app context
func Run(s *grpc.Server, cfg RuntimeConfig) {
	appCtx, done := signalcontext.Wrap(
		log.NewContextLogger(context.Background()),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer done()

	RunWithContext(appCtx, s, cfg)
}

// RunWithContext run grpc server with healthcheck using existing background context
func RunWithContext(appCtx context.Context, s *grpc.Server, cfg RuntimeConfig) {
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
