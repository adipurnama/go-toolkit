package log

import (
	"context"

	"github.com/adipurnama/go-toolkit/web"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"go.opencensus.io/trace"
)

// Panic - zero log.Panic
func Panic() *zerolog.Event {
	return zlog.Panic()
}

// Info - zero log.Info
func Info() *zerolog.Event {
	return zlog.Info()
}

// Debug - zero log.Debug
func Debug() *zerolog.Event {
	return zlog.Debug()
}

// Warn - zero log.Warn
func Warn() *zerolog.Event {
	return zlog.Warn()
}

// Fatal - zero log.Fatal
func Fatal() *zerolog.Event {
	return zlog.Fatal()
}

// Error - zero log.Error
func Error() *zerolog.Event {
	return zlog.Error()
}

// InfoCtx equals to :
// 	zerolog.Ctx(ctx).Info().Str("request_id", requestID(ctx)).Str("trace_id", traceID(ctx))
func InfoCtx(ctx context.Context) *zerolog.Event {
	return eventWithMetadata(ctx, zlog.Info())
}

// ErrorCtx equals to :
// 	zerolog.Ctx(ctx).Error().Str("request_id", requestID(ctx)).Str("trace_id", traceID(ctx))
func ErrorCtx(ctx context.Context) *zerolog.Event {
	return eventWithMetadata(ctx, zlog.Error())
}

// WarnCtx equals to :
// 	zerolog.Ctx(ctx).Warn().Str("request_id", requestID(ctx)).Str("trace_id", traceID(ctx))
func WarnCtx(ctx context.Context) *zerolog.Event {
	return eventWithMetadata(ctx, zlog.Warn())
}

// DebugCtx equals to :
// 	zerolog.Ctx(ctx).Info().Str("request_id", requestID(ctx)).Str("trace_id", traceID(ctx))
func DebugCtx(ctx context.Context) *zerolog.Event {
	return eventWithMetadata(ctx, zlog.Debug())
}

// FatalCtx equals to :
// 	zerolog.Ctx(ctx).Fatal().Str("request_id", requestID(ctx)).Str("trace_id", traceID(ctx))
func FatalCtx(ctx context.Context) *zerolog.Event {
	return eventWithMetadata(ctx, zlog.Fatal())
}

// PanicCtx equals to :
// 	zerolog.Ctx(ctx).Panic().Str("request_id", requestID(ctx)).Str("trace_id", traceID(ctx))
func PanicCtx(ctx context.Context) *zerolog.Event {
	return eventWithMetadata(ctx, zlog.Panic())
}

func eventWithMetadata(ctx context.Context, e *zerolog.Event) *zerolog.Event {
	traceAttrs := []trace.Attribute{}
	if requestID := web.ValueFromContext(ctx, web.ContextKeyRequestID); requestID != "" {
		e = e.Str("request_id", requestID)
		traceAttrs = append(traceAttrs, trace.StringAttribute("request_id", requestID))
	}
	if tID := web.ValueFromContext(ctx, web.ContextKeyTraceID); tID != "" {
		e = e.Str("trace_id", tID)
		traceAttrs = append(traceAttrs, trace.StringAttribute("trace_id", tID))
	}
	// if span := trace.FromContext(ctx); span != nil {
	// 	span.AddAttributes(traceAttrs...)
	// }
	return e
}
