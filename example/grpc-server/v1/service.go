package v1

import (
	context "context"
	"fmt"
)

var _ ExampleServiceServer = (*Service)(nil)

// Service - ExampleServiceServer example.
type Service struct {
}

// Greet - ExampleServiceServer impl.
func (s *Service) Greet(_ context.Context, req *HelloRequest) (*HelloResponse, error) {
	return &HelloResponse{
		Greeting: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}
