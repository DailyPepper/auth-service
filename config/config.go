package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	GRPCAddr    string
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		GRPCAddr:    getEnv("GRPC_ADDR", ":50051"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/auth_db?sslmode=disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..")
}
