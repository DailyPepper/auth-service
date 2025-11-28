# Makefile

# Project variables
PROJECT_NAME=auth-service
PROTO_DIR=./docs/proto
GENERATED_DIR=./pkg/generated
PROTO_FILE=auth/auth.proto
FULL_PROTO_PATH=$(PROTO_DIR)/$(PROTO_FILE)
GOBIN=$(shell go env GOPATH)/bin
DB_URL=postgres://postgres:password@localhost:5432/auth_db?sslmode=disable

.PHONY: all proto build clean check-deps help
.PHONY: docker-up docker-down docker-logs
.PHONY: migrate-create migrate-up migrate-down migrate-status
.PHONY: run test reset dev

# Default target
all: check-deps proto build

# =============================================================================
# DEVELOPMENT
# =============================================================================

# Запуск в режиме разработки (база + миграции + приложение)
dev: docker-up migrate-up run

# Запуск приложения
run:
	@echo "🚀 Starting $(PROJECT_NAME)..."
	go run cmd/auth/main.go

# Запуск тестов
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# =============================================================================
# PROTOBUF
# =============================================================================

# Проверка зависимостей для protobuf (без goose)
check-proto-deps:
	@echo "🔍 Checking protobuf dependencies..."
	@which protoc > /dev/null || (echo "❌ Error: protoc not installed. Run: make install-deps" && false)
	@test -f "$(GOBIN)/protoc-gen-go" || (echo "❌ Error: protoc-gen-go not installed. Run: make install-deps" && false)
	@test -f "$(GOBIN)/protoc-gen-go-grpc" || (echo "❌ Error: protoc-gen-go-grpc not installed. Run: make install-deps" && false)
	@echo "✅ All protobuf dependencies found!"

# Проверка всех зависимостей (включая goose)
check-deps: check-proto-deps
	@which goose > /dev/null || (echo "⚠️  Warning: goose not installed. Run: go install github.com/pressly/goose/v3/cmd/goose@latest" && true)
	@echo "✅ All dependencies checked!"

# Генерация protobuf кода (требует только protobuf зависимости)
proto: check-proto-deps
	@echo "📝 Generating protobuf code..."
	@echo "📄 Proto file: $(FULL_PROTO_PATH)"
	@echo "📁 Output dir: $(GENERATED_DIR)"
	
	# Создаем директорию если не существует
	mkdir -p $(GENERATED_DIR)/auth
	
	# Добавляем GOBIN в PATH для этой команды
	PATH="$(GOBIN):$$PATH" protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(GENERATED_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GENERATED_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)
	
	@echo "✅ Protobuf code generated successfully!"
	@echo "📁 Generated files in: $(GENERATED_DIR)/auth"

# Установка всех зависимостей
install-deps:
	@echo "📦 Installing dependencies..."
	
	# Protobuf compiler
	@if [ "$$(uname)" = "Darwin" ]; then \
		echo "🍎 Installing protobuf on macOS..."; \
		brew install protobuf; \
	elif [ "$$(uname)" = "Linux" ]; then \
		echo "🐧 Installing protobuf on Linux..."; \
		sudo apt-get update && sudo apt-get install -y protobuf-compiler; \
	else \
		echo "❌ Unsupported OS"; \
		exit 1; \
	fi
	
	# Go protobuf plugins
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	
	# Goose for migrations
	go install github.com/pressly/goose/v3/cmd/goose@latest
	
	@echo "✅ Dependencies installed successfully!"
	@echo "📝 Please run: source ~/.zshrc or source ~/.bashrc"

# =============================================================================
# DATABASE & MIGRATIONS
# =============================================================================

# Запуск базы данных
docker-up:
	@echo "🐘 Starting PostgreSQL database..."
	docker-compose up -d postgres
	@echo "⏳ Waiting for database to be ready..."
	@sleep 5
	@echo "✅ Database is ready!"

# Остановка базы данных
docker-down:
	@echo "🛑 Stopping database..."
	docker-compose down

# Просмотр логов базы данных
docker-logs:
	docker-compose logs -f postgres

# Создание новой миграции (требует goose)
migrate-create: check-goose
	@read -p "📝 Enter migration name: " name; \
	$(GOBIN)/goose -dir migrations create $${name} sql
	@echo "✅ Migration created in migrations/ directory"

# Применить миграции (требует goose)
migrate-up: check-goose
	@echo "🔄 Applying database migrations..."
	$(GOBIN)/goose -dir migrations postgres "$(DB_URL)" up
	@echo "✅ Migrations applied successfully"

# Откатить последнюю миграцию (требует goose)
migrate-down: check-goose
	@echo "↩️  Rolling back last migration..."
	$(GOBIN)/goose -dir migrations postgres "$(DB_URL)" down
	@echo "✅ Migration rolled back"

# Показать статус миграций (требует goose)
migrate-status: check-goose
	@echo "📊 Migration status:"
	$(GOBIN)/goose -dir migrations postgres "$(DB_URL)" status

# Проверка наличия goose
check-goose:
	@which goose > /dev/null || (echo "❌ Error: goose not installed. Run: make install-deps" && false)

# =============================================================================
# BUILD & DEPLOY
# =============================================================================

# Сборка проекта
build: proto
	@echo "🔨 Building $(PROJECT_NAME)..."
	mkdir -p bin
	go build -o bin/$(PROJECT_NAME) cmd/auth/main.go
	@echo "✅ Build completed: bin/$(PROJECT_NAME)"

# Полный перезапуск (база + миграции)
reset: docker-down docker-up migrate-up
	@echo "🔄 System reset completed"

# Очистка
clean:
	@echo "🧹 Cleaning generated files..."
	rm -rf $(GENERATED_DIR)/*
	rm -rf bin/
	@echo "✅ Clean completed"

# Полная очистка (включая Docker volumes)
clean-all: clean
	@echo "🧹 Cleaning Docker volumes..."
	docker-compose down -v
	@echo "✅ Full clean completed"

# =============================================================================
# UTILS
# =============================================================================

# Показать информацию о путях
paths:
	@echo "=== Path Information ==="
	@echo "📁 GOPATH: $(shell go env GOPATH)"
	@echo "📁 GOBIN: $(GOBIN)"
	@echo "🔧 which protoc: $(shell which protoc)"
	@echo "🔧 which protoc-gen-go: $(shell which protoc-gen-go 2>/dev/null || echo '❌ NOT FOUND')"
	@echo "🔧 which protoc-gen-go-grpc: $(shell which protoc-gen-go-grpc 2>/dev/null || echo '❌ NOT FOUND')"
	@echo "🔧 which goose: $(shell which goose 2>/dev/null || echo '⚠️  NOT FOUND')"

# Подключение к базе данных
db-connect:
	@echo "🔗 Connecting to database..."
	psql "$(DB_URL)"

# Проверка здоровья базы данных
db-health:
	@echo "❤️  Checking database health..."
	@pg_isready -d "$(DB_URL)" || echo "❌ Database is not ready"

# Help
help:
	@echo "🏗️  $(PROJECT_NAME) - Available commands:"
	@echo ""
	@echo "📦 DEPENDENCIES:"
	@echo "  make install-deps    - Install all dependencies"
	@echo "  make check-deps      - Check if dependencies are installed"
	@echo "  make paths           - Show path information"
	@echo ""
	@echo "🔧 DEVELOPMENT:"
	@echo "  make dev             - Full dev setup (db + migrations + app)"
	@echo "  make run             - Run application"
	@echo "  make test            - Run tests"
	@echo ""
	@echo "📝 PROTOBUF:"
	@echo "  make proto           - Generate protobuf code"
	@echo ""
	@echo "🗄️  DATABASE:"
	@echo "  make docker-up       - Start database"
	@echo "  make docker-down     - Stop database"
	@echo "  make docker-logs     - View database logs"
	@echo "  make migrate-create  - Create new migration"
	@echo "  make migrate-up      - Apply migrations"
	@echo "  make migrate-down    - Rollback last migration"
	@echo "  make migrate-status  - Show migration status"
	@echo "  make db-connect      - Connect to database"
	@echo "  make db-health       - Check database health"
	@echo ""
	@echo "🏗️  BUILD:"
	@echo "  make build           - Build application"
	@echo "  make reset           - Full reset (db + migrations)"
	@echo "  make clean           - Clean generated files"
	@echo "  make clean-all       - Clean everything (including Docker volumes)"
	@echo ""
	@echo "❓ HELP:"
	@echo "  make help            - Show this help message"

# =============================================================================
# SHORTCUTS
# =============================================================================

# Alias for common commands
up: docker-up
down: docker-down
logs: docker-logs
migrate: migrate-up
status: migrate-status
db: docker-up