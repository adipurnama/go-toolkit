package grpckit

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// RequestTimeoutInterceptor adds timeout to incoming request if it doesn't exists yet.
func RequestTimeoutInterceptor(t time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel()

		resp, err = handler(ctx, req)

		return resp, err
	}
}
