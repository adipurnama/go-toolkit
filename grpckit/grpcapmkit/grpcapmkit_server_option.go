package grpcapmkit

import (
	"github.com/pinpoint-apm/pinpoint-go-agent"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgrpc"
)

type serverOptions struct {
	elasticTracer  *apm.Tracer
	pinpointAgent  pinpoint.Agent
	requestIgnorer apmgrpc.RequestIgnorerFunc
	streamIgnorer  apmgrpc.StreamIgnorerFunc
}

// ServerOption sets options for server-side tracing.
type ServerOption func(*serverOptions)

// WithElasticTracer returns a ServerOption which sets t as the tracer
// to use for tracing server requests.
func WithElasticTracer(t *apm.Tracer) ServerOption {
	if t == nil {
		panic("t == nil")
	}

	return func(o *serverOptions) {
		o.elasticTracer = t
	}
}

// WithPinpointAgent returns a ServerOption which sets a as the agent
// to use for tracing server requests.
func WithPinpointAgent(a pinpoint.Agent) ServerOption {
	return func(o *serverOptions) {
		o.pinpointAgent = a
	}
}
