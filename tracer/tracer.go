// Package tracer to trace function span to apm
// supported apm is elastic apm, pinpoint apm & opentelemetry
package tracer

import (
	"go.opentelemetry.io/otel"
	otelTrace "go.opentelemetry.io/otel/trace"
)

var o option

type option struct {
	tPinpoint   bool
	tElastic    bool
	tOtel       bool
	tOtelTracer otelTrace.Tracer
}

// Option ...
type Option func(opt *option)

// WithPinpoint setup tracer package with pinpoint apm tracing.
func WithPinpoint() Option {
	return func(opt *option) {
		opt.tPinpoint = true
	}
}

// WithElastic setup tracer package with elastic apm tracing.
func WithElastic() Option {
	return func(opt *option) {
		opt.tElastic = true
	}
}

// WithOpenTelemetry setup tracer package with opentelemetry apm tracing.
func WithOpenTelemetry(tracerName string) Option {
	return func(opt *option) {
		opt.tOtel = true
		opt.tOtelTracer = otel.Tracer(tracerName)
	}
}

// Setup ...
func Setup(opts ...Option) {
	for _, opt := range opts {
		opt(&o)
	}
}
