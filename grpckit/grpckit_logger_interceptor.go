// Package grpckit is helper package for gRPC
package grpckit

import (
	"context"
	"path"
	"time"

	"github.com/adipurnama/go-toolkit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const keyRequestID = "x-request-id"

// LoggerInterceptor adds logger to request context.Context & logs the upstream call output.
func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		service := path.Dir(info.FullMethod)[1:]
		if service == "grpc.health.v1.Health" {
			return handler(ctx, req)
		}

		start := time.Now()
		newCtx := newCtxWithLogger(ctx, info.FullMethod, start)

		resp, err = handler(newCtx, req)

		code := status.Code(err)
		fields := []interface{}{
			"grpc.code", code.String(),
			"grpc.time_ms", time.Since(start).Milliseconds(),
			"grpc.response", resp,
			"grpc.request", req,
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			fields = append(fields, "metadata", md)
		}

		if err != nil {
			if clientRequestErrorCode(code) {
				fields = append(fields, "error", err)
				log.FromCtx(newCtx).Info("gRPC request completed with error", fields...)
			} else {
				log.FromCtx(newCtx).Error(err, "gRPC request completed with error", fields...)
			}

			return nil, err
		}

		switch codeToLogLevel(code) {
		case log.LevelDebug:
			log.FromCtx(newCtx).Debug("gRPC request completed", fields...)
		case log.LevelWarn:
			log.FromCtx(newCtx).Warn("gRPC request completed", fields...)
		default:
			log.FromCtx(newCtx).Info("gRPC request completed", fields...)
		}

		return resp, nil
	}
}

func clientRequestErrorCode(c codes.Code) bool {
	codes := []codes.Code{
		codes.Unauthenticated,
		codes.PermissionDenied,
		codes.NotFound,
		codes.InvalidArgument,
	}

	for _, v := range codes {
		if c == v {
			return true
		}
	}

	return false
}

func newCtxWithLogger(ctx context.Context, fullMethodString string, start time.Time) context.Context {
	method := path.Base(fullMethodString)
	service := path.Dir(fullMethodString)[1:]

	fields := []interface{}{
		"grpc.service", service,
		"grpc.method", method,
		"grpc.start_time", start.Format(time.RFC3339),
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, "grpc.request.deadline", d.Format(time.RFC3339))
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if rIDs, ok := md[keyRequestID]; ok && len(rIDs) > 0 {
			fields = append(fields, "request_id", rIDs[0])
		}
	}

	return log.NewLoggingContext(ctx, fields...)
}

// logic copied from https://github.com/rs/zerolog/issues/211
// complete non-zerolog impl at https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/logging
func codeToLogLevel(code codes.Code) log.Level {
	switch code {
	case codes.OK:
		return log.LevelDebug
	case codes.Canceled:
		return log.LevelDebug
	case codes.Unknown:
		return log.LevelInfo
	case codes.InvalidArgument:
		return log.LevelDebug
	case codes.DeadlineExceeded:
		return log.LevelInfo
	case codes.NotFound:
		return log.LevelDebug
	case codes.AlreadyExists:
		return log.LevelDebug
	case codes.PermissionDenied:
		return log.LevelInfo
	case codes.Unauthenticated:
		return log.LevelInfo
	case codes.ResourceExhausted:
		return log.LevelDebug
	case codes.FailedPrecondition:
		return log.LevelDebug
	case codes.Aborted:
		return log.LevelDebug
	case codes.OutOfRange:
		return log.LevelDebug
	case codes.Unimplemented:
		return log.LevelWarn
	case codes.Internal:
		return log.LevelWarn
	case codes.Unavailable:
		return log.LevelWarn
	case codes.DataLoss:
		return log.LevelWarn
	default:
		return log.LevelInfo
	}
}
