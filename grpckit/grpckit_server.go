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

const (
	defaultPort                = 8288
	defaultReqTimeout          = 7 * time.Second
	defaultShutdownWaitTimeout = 7 * time.Second
)

// RuntimeConfig defines runtime configuration for grpc service with health check.
type RuntimeConfig struct {
	ShutdownWaitDuration time.Duration `json:"shutdown_wait_duration,omitempty"`
	RequestTimeout       time.Duration `json:"request_timeout,omitempty"`
	Port                 int           `json:"port,omitempty"`
	Name                 string        `json:"name,omitempty"`
	EnableReflection     bool          `json:"enable_reflection,omitempty"`
	HealthCheckFunc      `json:"-"`
}

func (cfg *RuntimeConfig) validate() {
	// port
	if cfg.Port == 0 {
		cfg.Port = defaultPort
	}

	// check for timeout setting
	if cfg.RequestTimeout == 0 {
		cfg.RequestTimeout = defaultReqTimeout
	}

	if cfg.ShutdownWaitDuration == 0 {
		cfg.ShutdownWaitDuration = defaultShutdownWaitTimeout
	}
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

	cfg.validate()

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

	log.FromCtx(appCtx).Info("serving gRPC service", "config", cfg)

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
