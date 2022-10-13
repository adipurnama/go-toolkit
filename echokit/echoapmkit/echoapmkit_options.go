package echoapmkit

import (
	"github.com/pinpoint-apm/pinpoint-go-agent"
	"go.elastic.co/apm"
	apmhttp "go.elastic.co/apm/module/apmhttp"
)

type options struct {
	elasticTracer         *apm.Tracer
	pinpointAgent         pinpoint.Agent
	elasticRequestIgnorer apmhttp.RequestIgnorerFunc
}

// APMOption sets options for tracing.
type APMOption func(*options)

// WithElasticTracer returns an APMOption which sets t as the tracer
// to use for tracing server requests.
func WithElasticTracer(t *apm.Tracer) APMOption {
	if t == nil {
		panic("tracer: elastic *apm.Tracer cannot be nil")
	}

	return func(o *options) {
		o.elasticTracer = t
	}
}

// WithElasticRequestIgnorer returns a APMOption which sets r as the
// function to use to determine whether or not a request should
// be ignored. If r is nil, all requests will be reported.
func WithElasticRequestIgnorer(r apmhttp.RequestIgnorerFunc) APMOption {
	if r == nil {
		r = apmhttp.IgnoreNone
	}

	return func(o *options) {
		o.elasticRequestIgnorer = r
	}
}

// WithPinpointAgent ...
func WithPinpointAgent(p pinpoint.Agent) APMOption {
	return func(o *options) {
		o.pinpointAgent = p
	}
}
