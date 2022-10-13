package grpcapmkit

import (
	"fmt"
	"runtime/debug"

	"github.com/pinpoint-apm/pinpoint-go-agent"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/adipurnama/go-toolkit/log"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgrpc"
)

var errPanicInternal = errors.New("found panic while serving request")

// NewUnaryServerInterceptor returns a grpc.UnaryServerInterceptor that
// traces gRPC requests with the given options.
//
// The interceptor will trace transactions with the "request" type for
// each incoming request. The transaction will be added to the context,
// so server methods can use apm.StartSpan with the provided context.
//
// By default, the interceptor will trace with apm.DefaultTracer,
// and will not recover any panics. Use WithTracer to specify an
// alternative tracer, and WithRecovery to enable panic recovery.
func NewUnaryServerInterceptor(o ...ServerOption) grpc.UnaryServerInterceptor {
	opts := serverOptions{
		elasticTracer:  apm.DefaultTracer,
		requestIgnorer: apmgrpc.DefaultServerRequestIgnorer(),
		streamIgnorer:  apmgrpc.DefaultServerStreamIgnorer(),
	}
	for _, o := range o {
		o(&opts)
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		var tx *apm.Transaction

		if opts.elasticTracer.Recording() && !opts.requestIgnorer(info) {
			tx, ctx = startTransaction(ctx, opts.elasticTracer, info.FullMethod)
			defer tx.End()
		}

		var pt pinpoint.Tracer

		if opts.pinpointAgent != nil && opts.pinpointAgent.Enable() {
			pt = startPinpointSpan(ctx, opts.pinpointAgent, info.FullMethod)
			defer pt.EndSpan()
			defer pt.NewSpanEvent(info.FullMethod).EndSpanEvent()

			ctx = pinpoint.NewContext(ctx, pt)
		}

		defer func() {
			r := recover()
			if r != nil {
				errStack := debug.Stack()

				if tx != nil {
					e := opts.elasticTracer.Recovered(r)
					e.SetTransaction(tx)
					e.Context.SetFramework("grpc", grpc.Version)
					e.Handled = true

					e.Send()
				}

				logErr, ok := r.(error)
				if !ok {
					err = errors.Wrap(errPanicInternal, fmt.Sprint(r))
				} else if pt != nil {
					pt.Span().SetError(logErr)
				}

				log.FromCtx(ctx).Error(
					logErr,
					"recovered from panic",
					"panic_stack", errStack,
				)

				err = status.Errorf(codes.Internal, "%s", r)
			}

			setTransactionResult(tx, err)
		}()

		resp, err = handler(ctx, req)

		return resp, err
	}
}

// NewStreamServerInterceptor returns a grpc.StreamServerInterceptor that
// traces gRPC stream requests with the given options.
//
// The interceptor will trace transactions with the "request" type for each
// incoming stream request. The transaction will be added to the context, so
// server methods can use apm.StartSpan with the provided context.
//
// By default, the interceptor will trace with apm.DefaultTracer, and will
// not recover any panics. Use WithTracer to specify an alternative tracer,
// and WithRecovery to enable panic recovery.
func NewStreamServerInterceptor(o ...ServerOption) grpc.StreamServerInterceptor {
	opts := serverOptions{
		elasticTracer: apm.DefaultTracer,
		streamIgnorer: apmgrpc.DefaultServerStreamIgnorer(),
	}
	for _, o := range o {
		o(&opts)
	}

	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		ctx := stream.Context()

		var tx *apm.Transaction

		if opts.elasticTracer.Recording() && !opts.streamIgnorer(info) {
			tx, ctx = startTransaction(ctx, opts.elasticTracer, info.FullMethod)
			defer tx.End()
		}

		var (
			pinpointWrappedStream *pinpointServerStream
			pt                    pinpoint.Tracer
		)

		if opts.pinpointAgent != nil && opts.pinpointAgent.Enable() {
			tracer := startPinpointSpan(stream.Context(), opts.pinpointAgent, info.FullMethod)
			defer tracer.EndSpan()
			defer tracer.NewSpanEvent(info.FullMethod).EndSpanEvent()

			ctx = pinpoint.NewContext(stream.Context(), tracer)
			pinpointWrappedStream = &pinpointServerStream{stream, ctx}
		}

		defer func() {
			r := recover()
			if r != nil {
				errStack := debug.Stack()

				if tx != nil {
					e := opts.elasticTracer.Recovered(r)
					e.SetTransaction(tx)
					e.Context.SetFramework("grpc", grpc.Version)
					e.Handled = true
					e.Send()
				}

				logErr, ok := r.(error)
				if !ok {
					err = errors.Wrap(errPanicInternal, fmt.Sprint(r))
				} else if pt != nil {
					pt.Span().SetError(logErr)
				}

				log.FromCtx(ctx).Error(
					logErr,
					"recovered from panic",
					"panic_stack", errStack,
				)

				err = status.Errorf(codes.Internal, "%s", r)
			}

			setTransactionResult(tx, err)
		}()

		if pinpointWrappedStream != nil {
			return handler(srv, pinpointWrappedStream)
		}

		return handler(srv, stream)
	}
}

func statusCodeFromError(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	statusCode := codes.Unknown

	if s, ok := status.FromError(err); ok {
		statusCode = s.Code()
	}

	return statusCode
}
