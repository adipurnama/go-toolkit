package grpckit

import (
	"context"

	"github.com/adipurnama/go-toolkit/log"
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

		if code == codes.Unknown {
			var spbStatus *spb.Status

			if errHandler != nil {
				spbStatus = errHandler(err)
			} else {
				spbStatus = DefaultGRPCErrorHandler(err)
			}

			err = status.FromProto(spbStatus).Err()

			log.Println("got unknown error ", err)

			return resp, err
		}

		return resp, err
	}
}

// DefaultGRPCErrorHandler provides default error to standard protobuff status.
func DefaultGRPCErrorHandler(err error) *spb.Status {
	code := status.Code(err)

	return &spb.Status{
		Code:    int32(code),
		Message: err.Error(),
	}
}
