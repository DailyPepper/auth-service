package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DailyPepper/auth-service/config"
	"github.com/DailyPepper/auth-service/internal/repository"
	"github.com/DailyPepper/auth-service/internal/server"
	"github.com/DailyPepper/auth-service/internal/service"
	"github.com/DailyPepper/auth-service/pkg/logger"
	"github.com/DailyPepper/auth-service/pkg/migrations"
)

func main() {
	log := logger.New("info")
	log.Info("🔧 Initializing auth service...")

	log.Info("1. Loading configuration...")
	cfg := config.Load()
	log.Info("✅ Configuration loaded successfully")

	log.Info("2. Connecting to database...")
	userRepo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to database: %v", err)
	}
	defer func() {
		if err := userRepo.Close(); err != nil {
			log.Error("❌ Failed to close database connection: %v", err)
		}
	}()
	log.Info("✅ Database connection established")

	log.Info("3. Running database migrations...")
	if err := migrations.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatal("❌ Failed to run migrations: %v", err)
	}
	log.Info("✅ Database migrations completed")

	log.Info("4. Creating services...")
	registrService := service.NewRegistrService(userRepo)
	if registrService == nil {
		log.Fatal("❌ Failed to create registr service - returned nil")
	}
	log.Info("✅ Services created successfully")

	log.Info("5. Creating gRPC server...")
	grpcServer := server.NewGRPCServer(registrService)
	if grpcServer == nil {
		log.Fatal("❌ Failed to create gRPC server - returned nil")
	}
	log.Info("✅ gRPC server created successfully")

	log.Info("6. Starting gRPC server on %s...", cfg.GRPCAddr)

	serverErr := make(chan error, 1)

	go func() {
		if err := grpcServer.Start(cfg.GRPCAddr); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		log.Fatal("❌ Failed to start gRPC server: %v", err)
	default:
		// Сервер запустился успешно
	}

	log.Info("✅ Auth Service started successfully")
	log.Info("📍 gRPC Port: %s", cfg.GRPCAddr)
	log.Info("🗄️  Database: %s", cfg.DatabaseURL)
	log.Info("📡 Ready to accept gRPC requests")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	log.Info("⏳ Server is running. Press Ctrl+C to stop...")

	select {
	case <-quit:
		log.Info("🚨 Received shutdown signal...")
	case err := <-serverErr:
		log.Error("🚨 Server error: %v", err)
	}

	log.Info("🛑 Shutting down server...")
	grpcServer.Stop()
	log.Info("👋 Server stopped gracefully")
}
