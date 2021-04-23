package grpckit

import (
	"context"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCErrorHandler transforms error to valid grpc standard protobuff status.
type GRPCErrorHandler func(err error) *spb.Status

// ErrorResponseWriterInterceptor writes grpc status based on error found from upstream call.
func ErrorResponseWriterInterceptor(errHandler GRPCErrorHandler) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		code := status.Code(err)
		if code != codes.Unknown {
			return resp, err
		}

		var spbStatus *spb.Status

		if errHandler != nil {
			spbStatus = errHandler(err)
		} else {
			spbStatus = DefaultGRPCErrorHandler(err)
		}

		if spbStatus == nil {
			return resp, err
		}

		return resp, status.FromProto(spbStatus).Err()
	}
}

// DefaultGRPCErrorHandler provides default error to standard protobuff status.
func DefaultGRPCErrorHandler(err error) *spb.Status {
	return &spb.Status{
		Code:    int32(codes.Internal),
		Message: err.Error(),
	}
}
