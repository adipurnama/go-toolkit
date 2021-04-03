package grpckit

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestIDProvider generates new request id for each requests.
type RequestIDProvider interface {
	NewRequestID() string
}

// DefaultRequestIDProvider generates uuid.V4 string for request id.
func DefaultRequestIDProvider() RequestIDProvider {
	return &uuidV4RequestIDProvider{}
}

type uuidV4RequestIDProvider struct{}

// NewRequestID implements RequestIDProvider interface.
func (p *uuidV4RequestIDProvider) NewRequestID() string {
	return uuid.NewV4().String()
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
