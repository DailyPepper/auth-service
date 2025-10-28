# Makefile

PROTO_DIR=./docs/proto
GENERATED_DIR=./pkg/generated
PROTO_FILE=auth/auth.proto
FULL_PROTO_PATH=$(PROTO_DIR)/$(PROTO_FILE)
GOBIN=$(shell go env GOPATH)/bin

.PHONY: all proto build clean check-deps

# Проверка зависимостей
check-deps:
	@echo "Checking dependencies..."
	@which protoc > /dev/null || (echo "Error: protoc not installed. Run: brew install protobuf" && false)
	@test -f "$(GOBIN)/protoc-gen-go" || (echo "Error: protoc-gen-go not installed. Run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" && false)
	@test -f "$(GOBIN)/protoc-gen-go-grpc" || (echo "Error: protoc-gen-go-grpc not installed. Run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest" && false)
	@echo "All dependencies found!"

# Генерация protobuf кода
proto: check-deps
	@echo "Generating protobuf code..."
	@echo "Proto file: $(FULL_PROTO_PATH)"
	@echo "Output dir: $(GENERATED_DIR)"
	@echo "Go bin path: $(GOBIN)"
	
	# Создаем директорию если не существует
	mkdir -p $(GENERATED_DIR)/auth
	
	# Добавляем GOBIN в PATH для этой команды
	PATH="$(GOBIN):$$PATH" protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(GENERATED_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GENERATED_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)
	
	@echo "Protobuf code generated successfully!"
	@echo "Generated files in: $(GENERATED_DIR)/auth"

# Установка всех зависимостей
install-deps:
	@echo "Installing dependencies..."
	brew install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Please run: source ~/.zshrc"

# Показать информацию о путях
paths:
	@echo "=== Path Information ==="
	@echo "GOPATH: $(shell go env GOPATH)"
	@echo "GOBIN: $(GOBIN)"
	@echo "which protoc: $(shell which protoc)"
	@echo "which protoc-gen-go: $(shell which protoc-gen-go 2>/dev/null || echo 'NOT FOUND')"
	@echo "which protoc-gen-go-grpc: $(shell which protoc-gen-go-grpc 2>/dev/null || echo 'NOT FOUND')"

# Сборка проекта
build: proto
	@echo "Building application..."
	go build -o bin/auth-server cmd/auth-server/main.go

# Очистка
clean:
	rm -rf $(GENERATED_DIR)/*
	rm -rf bin/

help:
	@echo "Available commands:"
	@echo "  make install-deps - Install all dependencies"
	@echo "  make paths        - Show path information"
	@echo "  make proto        - Generate protobuf code"
	@echo "  make build        - Build application"
	@echo "  make clean        - Clean generated files"