package tracer

import (
	"context"

	"github.com/pinpoint-apm/pinpoint-go-agent"
	"go.elastic.co/apm"
	"go.opentelemetry.io/otel/attribute"
	otelTrace "go.opentelemetry.io/otel/trace"

	"github.com/adipurnama/go-toolkit/runtimekit"
)

// Span describes an operation within a transaction.
type Span struct {
	pSpan pinpoint.Tracer
	eSpan *apm.Span
	oSpan otelTrace.Span
}

const skipFuncCount = 2

// SpanLevel describes span at source code logic level.
type SpanLevel string

var (
	// SpanLvlServiceLogic for service function / methods.
	SpanLvlServiceLogic SpanLevel = "service"
	// SpanLvlDBQuery for db query call / repository level.
	SpanLvlDBQuery SpanLevel = "database"
	// SpanLvlHTTPCall for HTTP api client api call level.
	SpanLvlHTTPCall SpanLevel = "http.client"
	// SpanLvlGrpcCall for gRPC api client api call level.
	SpanLvlGrpcCall SpanLevel = "grpc.client"
	// SpanLvlHandler for http controller / handler func.
	SpanLvlHandler SpanLevel = "handler"
)

// NewSpan returns span based on layer type.
func NewSpan(ctx context.Context, lvl SpanLevel) (context.Context, *Span) {
	span := Span{}
	childCtx := ctx

	if o.tElastic {
		tx := apm.TransactionFromContext(ctx)
		span.eSpan = tx.StartSpan(runtimekit.SkippedFunctionName(skipFuncCount), string(lvl), nil)
	}

	if o.tOtel {
		spanCtx, s := o.tOtelTracer.Start(ctx, runtimekit.SkippedFunctionName(skipFuncCount))
		s.SetAttributes(attribute.String("level", string(lvl)))
		span.oSpan = s
		childCtx = spanCtx
	}

	if o.tPinpoint {
		span.pSpan = pinpoint.FromContext(ctx).NewSpanEvent(runtimekit.SkippedFunctionName(skipFuncCount))
	}

	return childCtx, &span
}

// ServiceFuncSpan returns span for service layer type.
func ServiceFuncSpan(ctx context.Context) (context.Context, *Span) {
	lvl := SpanLvlServiceLogic
	span := Span{}
	childCtx := ctx

	if o.tElastic {
		tx := apm.TransactionFromContext(ctx)
		span.eSpan = tx.StartSpan(runtimekit.SkippedFunctionName(skipFuncCount), string(lvl), nil)
	}

	if o.tOtel {
		spanCtx, s := o.tOtelTracer.Start(ctx, runtimekit.SkippedFunctionName(skipFuncCount))
		s.SetAttributes(attribute.String("level", string(lvl)))
		span.oSpan = s
		childCtx = spanCtx
	}

	if o.tPinpoint {
		span.pSpan = pinpoint.FromContext(ctx).NewSpanEvent(runtimekit.SkippedFunctionName(skipFuncCount))
	}

	return childCtx, &span
}

// RepositoryFuncSpan returns span for repository layer type.
func RepositoryFuncSpan(ctx context.Context) (context.Context, *Span) {
	lvl := SpanLvlDBQuery
	span := Span{}
	childCtx := ctx

	if o.tElastic {
		tx := apm.TransactionFromContext(ctx)
		span.eSpan = tx.StartSpan(runtimekit.SkippedFunctionName(skipFuncCount), string(lvl), nil)
	}

	if o.tOtel {
		spanCtx, s := o.tOtelTracer.Start(ctx, runtimekit.SkippedFunctionName(skipFuncCount))
		s.SetAttributes(attribute.String("level", string(lvl)))
		span.oSpan = s
		childCtx = spanCtx
	}

	if o.tPinpoint {
		span.pSpan = pinpoint.FromContext(ctx).NewSpanEvent(runtimekit.SkippedFunctionName(skipFuncCount))
	}

	return childCtx, &span
}

// AddEvent proxy opentelemetry func.
func (s *Span) AddEvent(eventName string, opts ...otelTrace.EventOption) {
	if o.tOtel {
		s.oSpan.AddEvent(eventName, opts...)
	}
}

// SetAttributes proxy opentelemetry func.
func (s *Span) SetAttributes(kv ...attribute.KeyValue) {
	if o.tOtel {
		s.oSpan.SetAttributes(kv...)
	}
}

// End ends span event.
func (s *Span) End() {
	if o.tPinpoint {
		s.pSpan.EndSpanEvent()
	}

	if o.tElastic {
		s.eSpan.End()
	}

	if o.tOtel {
		s.oSpan.End()
	}
}
