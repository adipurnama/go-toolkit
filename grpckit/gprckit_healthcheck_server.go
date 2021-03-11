package grpckit

import (
	"context"

	"github.com/adipurnama/go-toolkit/grpckit/grpc_health_v1"
	"github.com/adipurnama/go-toolkit/log"
)

// HealthCheckServer is default grpc health check provider.
type HealthCheckServer struct {
	Serving         bool
	healthCheckFunc HealthCheckFunc
}

// NewHealthcheckServer - factory.
func NewHealthcheckServer(hcFunc HealthCheckFunc) *HealthCheckServer {
	return &HealthCheckServer{
		Serving:         true,
		healthCheckFunc: hcFunc,
	}
}

// HealthCheckFunc - health check template func.
type HealthCheckFunc func(context.Context) error

// Check - grpc_health_v1.Server impl.
func (s *HealthCheckServer) Check(ctx context.Context, _ *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	resp := grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}

	if !s.Serving {
		resp.Status = grpc_health_v1.HealthCheckResponse_NOT_SERVING

		return &resp, nil
	}

	if s.healthCheckFunc == nil {
		resp.Status = grpc_health_v1.HealthCheckResponse_SERVING

		return &resp, nil
	}

	err := s.healthCheckFunc(ctx)
	if err != nil {
		resp.Status = grpc_health_v1.HealthCheckResponse_NOT_SERVING

		log.FromCtx(ctx).Error(err, "healthCheck.Check")
	}

	return &resp, nil
}

// Watch - grpc_health_v1.Server impl.
func (s *HealthCheckServer) Watch(_ *grpc_health_v1.HealthCheckRequest, _ grpc_health_v1.Health_WatchServer) error {
	log.Println("healthCheck.Watch", "status=", "not implemented")
	return nil
}
