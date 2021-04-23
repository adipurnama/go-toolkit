package grpckit

import (
	"context"

	shortuuid "github.com/lithammer/shortuuid/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestIDProvider generates new request id for each requests.
type RequestIDProvider interface {
	NewRequestID() string
}

// DefaultRequestIDProvider generates shortuuid string for request id.
func DefaultRequestIDProvider() RequestIDProvider {
	return &shortuuidRequestIDProvider{}
}

type shortuuidRequestIDProvider struct{}

// NewRequestID implements RequestIDProvider interface.
func (p *shortuuidRequestIDProvider) NewRequestID() string {
	return shortuuid.New()
}

// RequestIDInterceptor add request id to incoming request if it doesn't exists yet.
func RequestIDInterceptor(rIDProvider RequestIDProvider) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if rIDProvider == nil {
			rIDProvider = DefaultRequestIDProvider()
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if _, rIDFound := md[keyRequestID]; !rIDFound {
				md.Set(keyRequestID, rIDProvider.NewRequestID())
				ctx = metadata.NewIncomingContext(ctx, md)
			}
		}

		resp, err = handler(ctx, req)

		return resp, err
	}
}
