package server

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func (s *GRPCServer) unaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("📨 gRPC method called: %s", info.FullMethod)

		// Можно добавить дополнительную логику:
		// - Аутентификацию
		// - Валидацию
		// - Метрики
		// - Трассировку

		resp, err := handler(ctx, req)

		if err != nil {
			log.Printf("❌ gRPC method %s failed: %v", info.FullMethod, err)
		} else {
			log.Printf("✅ gRPC method %s completed successfully", info.FullMethod)
		}

		return resp, err
	}
}
