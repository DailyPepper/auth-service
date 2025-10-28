package server

import (
	"auth-service/internal/service"
	"auth-service/pkg/generated/auth"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	auth.UnimplementedAuthServiceServer
	registrService service.Registr
	server         *grpc.Server
}

func NewGRPCServer(registrService service.Registr) *GRPCServer {
	return &GRPCServer{
		registrService: registrService,
	}
}

func (s *GRPCServer) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(s.unaryInterceptor()),
	)

	auth.RegisterAuthServiceServer(s.server, s)

	log.Printf("gRPC server starting on port %s", port)
	return s.server.Serve(lis)
}

func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
		log.Println("gRPC server stopped")
	}
}
