package v1

import (
	context "context"
	"fmt"
	"time"

	"github.com/adipurnama/go-toolkit/grpckit/grpcapmkit"
	"github.com/pkg/errors"
)

var _ ExampleServiceServer = (*Server)(nil)

var errInternalServer = errors.New("internal server error")

// Server - ExampleServiceServer example.
type Server struct {
}

// Greet - ExampleServiceServer impl.
func (s *Server) Greet(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	span := grpcapmkit.ServerSpan(ctx)
	defer span.End()

	if !getRandomBool() {
		return nil, errors.WithStack(errors.Wrap(errInternalServer, "failed because of random bool"))
	}

	return &HelloResponse{
		Greeting: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}

func getRandomBool() bool {
	now := time.Now()
	return now.Second()%2 == 0
}
