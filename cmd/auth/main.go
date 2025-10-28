package main

import (
	"auth-service/internal/server"
	"auth-service/internal/service"
	"auth-service/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := logger.New("info")

	log.Info("🔧 Initializing auth service...")

	log.Info("1. Creating registr service...")
	registrService := service.NewRegistrService()
	if registrService == nil {
		log.Fatal("❌ Failed to create registr service - returned nil")
	}
	log.Info("✅ Registr service created successfully")

	log.Info("2. Creating gRPC server...")
	grpcServer := server.NewGRPCServer(registrService)
	if grpcServer == nil {
		log.Fatal("❌ Failed to create gRPC server - returned nil")
	}
	log.Info("✅ gRPC server created successfully")

	log.Info("3. Starting gRPC server on port 50051...")
	go func() {
		if err := grpcServer.Start("50051"); err != nil {
			log.Fatal("❌ Failed to start gRPC server: %v", err)
		}
	}()

	log.Info("✅ Auth Service gRPC server started successfully")
	log.Info("📍 Port: 50051")
	log.Info("📡 Ready to accept gRPC requests")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Info("⏳ Server is running. Press Ctrl+C to stop...")
	<-quit

	log.Info("🛑 Shutting down server...")
	grpcServer.Stop()
	log.Info("👋 Server stopped")
}
