package v1

import (
	context "context"
	"fmt"

	"github.com/adipurnama/go-toolkit/tracer"
)

var _ ExampleServiceServer = (*Service)(nil)

// Service - ExampleServiceServer example.
type Service struct {
}

// Greet - ExampleServiceServer impl.
func (s *Service) Greet(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	span := tracer.ServiceFuncSpan(ctx)
	defer span.End()

	return &HelloResponse{
		Greeting: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}
