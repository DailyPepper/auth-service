package main

import (
	"auth-service/config"
	"auth-service/internal/repository"
	"auth-service/internal/server"
	"auth-service/internal/service"
	"auth-service/pkg/logger"
	"auth-service/pkg/migrations"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Инициализация логгера
	log := logger.New("info")
	log.Info("🔧 Initializing auth service...")

	// 1. Загрузка конфигурации
	log.Info("1. Loading configuration...")
	cfg := config.Load()
	log.Info("✅ Configuration loaded successfully")

	// 2. Подключение к базе данных
	log.Info("2. Connecting to database...")
	userRepo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to database: %v", err)
	}
	defer userRepo.Close()
	log.Info("✅ Database connection established")

	// 3. Запуск миграций
	log.Info("3. Running database migrations...")
	if err := migrations.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatal("❌ Failed to run migrations: %v", err)
	}
	log.Info("✅ Database migrations completed")

	// 4. Создание сервисов
	log.Info("4. Creating services...")
	registrService := service.NewRegistrService(userRepo)
	if registrService == nil {
		log.Fatal("❌ Failed to create registr service - returned nil")
	}
	log.Info("✅ Services created successfully")

	// 5. Создание и запуск gRPC сервера
	log.Info("5. Creating gRPC server...")
	grpcServer := server.NewGRPCServer(registrService)
	if grpcServer == nil {
		log.Fatal("❌ Failed to create gRPC server - returned nil")
	}
	log.Info("✅ gRPC server created successfully")

	log.Info("6. Starting gRPC server on %s...", cfg.GRPCAddr)
	go func() {
		if err := grpcServer.Start(cfg.GRPCAddr); err != nil {
			log.Fatal("❌ Failed to start gRPC server: %v", err)
		}
	}()

	log.Info("✅ Auth Service started successfully")
	log.Info("📍 gRPC Port: %s", cfg.GRPCAddr)
	log.Info("🗄️  Database: %s", cfg.DatabaseURL)
	log.Info("📡 Ready to accept gRPC requests")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Info("⏳ Server is running. Press Ctrl+C to stop...")
	<-quit

	log.Info("🛑 Shutting down server...")
	grpcServer.Stop()
	log.Info("👋 Server stopped")
}
