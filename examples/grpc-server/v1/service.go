package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/tracer"
)

var _ ExampleServiceServer = (*Server)(nil)

var errInternalServer = errors.New("internal server error")

// Server - ExampleServiceServer example.
type Server struct{}

// Greet - ExampleServiceServer impl.
func (s *Server) Greet(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	_, span := tracer.NewSpan(ctx, tracer.SpanLvlHandler)
	defer span.End()

	if !getRandomBool() {
		return nil, errors.WithStack(errors.Wrap(errInternalServer, "failed because of random bool"))
	}

	return &HelloResponse{
		Greeting: fmt.Sprintf("Hello %s. Your are %d years old.", req.Name, req.Age),
	}, nil
}

// Crash - ExampleServiceServer impl.
func (s *Server) Crash(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	zero := 0
	ten := 10
	divByZero := ten / zero

	log.FromCtx(ctx).Info("handling crash", "divByZero", divByZero)

	return nil, nil
}

func getRandomBool() bool {
	now := time.Now()
	return now.Second()%2 == 0
}
