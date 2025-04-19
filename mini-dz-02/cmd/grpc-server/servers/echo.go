package servers

import (
	"context"
	"mini-dz-02/internal/proto/zoo"
)

// EchoServer implements the EchoService
type EchoServer struct {
	zoo.UnimplementedEchoServiceServer
}

// NewEchoServer creates a new EchoServer
func NewEchoServer() *EchoServer {
	return &EchoServer{}
}

// Echo implements the Echo method of the EchoService
func (s *EchoServer) Echo(ctx context.Context, req *zoo.EchoRequest) (*zoo.EchoResponse, error) {
	return &zoo.EchoResponse{
		Message: req.Message,
	}, nil
}
