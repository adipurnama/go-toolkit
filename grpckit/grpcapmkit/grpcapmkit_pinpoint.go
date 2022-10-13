package grpcapmkit

import (
	"context"

	"github.com/pinpoint-apm/pinpoint-go-agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const serviceTypeGrpcServer = 1130

type pinpointServerStream struct {
	grpc.ServerStream
	context context.Context
}

func (s *pinpointServerStream) Context() context.Context {
	return s.context
}

type distributedTracingContextReaderMD struct {
	md metadata.MD
}

func (m distributedTracingContextReaderMD) Get(key string) string {
	v := m.md.Get(key)
	if len(v) == 0 {
		return ""
	}

	return v[0]
}

func startPinpointSpan(ctx context.Context, agent pinpoint.Agent, rpcName string) pinpoint.Tracer {
	md, _ := metadata.FromIncomingContext(ctx) // nil is ok
	reader := &distributedTracingContextReaderMD{md}
	tracer := agent.NewSpanTracerWithReader("gRPC Server", rpcName, reader)
	tracer.Span().SetServiceType(serviceTypeGrpcServer)

	return tracer
}
