.PHONY: help build run stop test lint clean db-migrate db-rollback db-seed proto install dev

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

#############
# HELP
#############

help: ## Show this help message
	@echo "$(BLUE)═══════════════════════════════════════════════════════════$(NC)"
	@echo "$(GREEN)  Sneakers Marketplace - Makefile Commands$(NC)"
	@echo "$(BLUE)═══════════════════════════════════════════════════════════$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-25s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(BLUE)Quick Start:$(NC)"
	@echo "  1. $(GREEN)make install$(NC)     - Install dependencies"
	@echo "  2. $(GREEN)make db-setup$(NC)    - Setup database"
	@echo "  3. $(GREEN)make dev$(NC)         - Start all services"
	@echo ""

#############
# INSTALL
#############

install: install-backend install-frontend install-tools ## Install all dependencies
	@echo "$(GREEN)✅ All dependencies installed!$(NC)"

install-backend: ## Install Go dependencies
	@echo "$(BLUE)📦 Installing Go dependencies...$(NC)"
	@go mod download
	@go mod tidy

install-frontend: ## Install Frontend dependencies
	@echo "$(BLUE)📦 Installing Frontend dependencies...$(NC)"
	@cd frontend && npm install

install-tools: ## Install development tools
	@echo "$(BLUE)🔧 Installing development tools...$(NC)"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@which golangci-lint > /dev/null || brew install golangci-lint
	@which protoc > /dev/null || brew install protobuf
	@which mailhog > /dev/null || brew install mailhog
	@echo "$(GREEN)✅ Tools installed!$(NC)"

#############
# BUILD
#############

build: build-all ## Build all services

build-all: ## Build all backend services
	@echo "$(BLUE)🔨 Building all services...$(NC)"
	@mkdir -p bin
	@$(MAKE) build-gateway
	@$(MAKE) build-user
	@$(MAKE) build-product
	@$(MAKE) build-bidding
	@$(MAKE) build-order
	@$(MAKE) build-payment
	@$(MAKE) build-notification
	@echo "$(GREEN)✅ All services built successfully!$(NC)"

build-gateway: ## Build API Gateway
	@echo "$(YELLOW)Building API Gateway...$(NC)"
	@go build -o bin/api-gateway ./cmd/api-gateway

build-user: ## Build User Service
	@echo "$(YELLOW)Building User Service...$(NC)"
	@go build -o bin/user-service ./cmd/user-service

build-product: ## Build Product Service
	@echo "$(YELLOW)Building Product Service...$(NC)"
	@go build -o bin/product-service ./cmd/product-service

build-bidding: ## Build Bidding Service
	@echo "$(YELLOW)Building Bidding Service...$(NC)"
	@go build -o bin/bidding-service ./cmd/bidding-service

build-order: ## Build Order Service
	@echo "$(YELLOW)Building Order Service...$(NC)"
	@go build -o bin/order-service ./cmd/order-service

build-payment: ## Build Payment Service
	@echo "$(YELLOW)Building Payment Service...$(NC)"
	@go build -o bin/payment-service ./cmd/payment-service

build-notification: ## Build Notification Service
	@echo "$(YELLOW)Building Notification Service...$(NC)"
	@go build -o bin/notification-service ./cmd/notification-service

build-frontend: ## Build Frontend for production
	@echo "$(BLUE)🔨 Building Frontend...$(NC)"
	@cd frontend && npm run build
	@echo "$(GREEN)✅ Frontend built!$(NC)"

#############
# RUN
#############

dev: ## Start all services in development mode
	@echo "$(BLUE)🚀 Starting all services...$(NC)"
	@$(MAKE) db-check
	@$(MAKE) build-all
	@$(MAKE) run-all
	@$(MAKE) run-frontend
	@echo "$(GREEN)✅ All services running!$(NC)"
	@echo ""
	@echo "$(BLUE)═══════════════════════════════════════════════════════════$(NC)"
	@echo "$(GREEN)Frontend:        $(NC)http://localhost:5173"
	@echo "$(GREEN)API Gateway:     $(NC)http://localhost:8080"
	@echo "$(GREEN)Mailhog UI:      $(NC)http://localhost:8025"
	@echo "$(BLUE)═══════════════════════════════════════════════════════════$(NC)"

run-all: ## Start all backend services
	@echo "$(BLUE)🚀 Starting backend services...$(NC)"
	@source .env && nohup ./bin/user-service > /tmp/user-service.log 2>&1 & echo $$! > /tmp/user-service.pid
	@source .env && nohup ./bin/product-service > /tmp/product-service.log 2>&1 & echo $$! > /tmp/product-service.pid
	@source .env && nohup ./bin/bidding-service > /tmp/bidding-service.log 2>&1 & echo $$! > /tmp/bidding-service.pid
	@source .env && nohup ./bin/notification-service > /tmp/notification-service.log 2>&1 & echo $$! > /tmp/notification-service.pid
	@source .env && nohup ./bin/api-gateway > /tmp/api-gateway.log 2>&1 & echo $$! > /tmp/api-gateway.pid
	@sleep 2
	@echo "$(GREEN)✅ Backend services started!$(NC)"

run-frontend: ## Start Frontend dev server
	@echo "$(BLUE)🚀 Starting Frontend...$(NC)"
	@cd frontend && nohup npm run dev > /tmp/vite.log 2>&1 & echo $$! > /tmp/vite.pid
	@sleep 2
	@echo "$(GREEN)✅ Frontend started!$(NC)"

run-mailhog: ## Start Mailhog (email testing)
	@echo "$(BLUE)📧 Starting Mailhog...$(NC)"
	@nohup mailhog > /tmp/mailhog.log 2>&1 & echo $$! > /tmp/mailhog.pid
	@echo "$(GREEN)✅ Mailhog started at http://localhost:8025$(NC)"

stop: ## Stop all services
	@echo "$(BLUE)🛑 Stopping all services...$(NC)"
	@-pkill -f api-gateway || true
	@-pkill -f user-service || true
	@-pkill -f product-service || true
	@-pkill -f bidding-service || true
	@-pkill -f order-service || true
	@-pkill -f payment-service || true
	@-pkill -f notification-service || true
	@-pkill -f "vite" || true
	@-pkill -f mailhog || true
	@-rm -f /tmp/*.pid
	@echo "$(GREEN)✅ All services stopped!$(NC)"

restart: stop build-all run-all ## Restart all services
	@echo "$(GREEN)✅ All services restarted!$(NC)"

logs: ## Show logs from all services
	@echo "$(BLUE)📋 Showing logs...$(NC)"
	@echo "$(YELLOW)API Gateway:$(NC)"
	@tail -20 /tmp/api-gateway.log
	@echo ""
	@echo "$(YELLOW)Bidding Service:$(NC)"
	@tail -20 /tmp/bidding-service.log

logs-follow: ## Follow logs from all services
	@echo "$(BLUE)📋 Following logs... (Ctrl+C to stop)$(NC)"
	@tail -f /tmp/*.log

status: ## Check status of all services
	@echo "$(BLUE)📊 Service Status:$(NC)"
	@ps aux | grep -E 'api-gateway|user-service|product-service|bidding-service|notification-service|vite|mailhog' | grep -v grep || echo "$(RED)No services running$(NC)"

#############
# DATABASE
#############

db-check: ## Check database connection
	@echo "$(BLUE)🔍 Checking database connection...$(NC)"
	@psql -U postgres -d sneakers_marketplace -c "SELECT 1;" > /dev/null 2>&1 && echo "$(GREEN)✅ Database connected!$(NC)" || (echo "$(RED)❌ Database not accessible$(NC)" && exit 1)

db-setup: db-create db-migrate db-seed ## Setup database (create, migrate, seed)
	@echo "$(GREEN)✅ Database setup complete!$(NC)"

db-create: ## Create database
	@echo "$(BLUE)🗄️  Creating database...$(NC)"
	@psql -U postgres -c "CREATE DATABASE sneakers_marketplace;" 2>/dev/null || echo "$(YELLOW)Database already exists$(NC)"
	@echo "$(GREEN)✅ Database ready!$(NC)"

db-drop: ## Drop database (DANGEROUS!)
	@echo "$(RED)⚠️  Dropping database...$(NC)"
	@psql -U postgres -c "DROP DATABASE IF EXISTS sneakers_marketplace;"
	@echo "$(GREEN)✅ Database dropped!$(NC)"

db-migrate: ## Run database migrations
	@echo "$(BLUE)🔄 Running migrations...$(NC)"
	@for file in internal/database/migrations/*.up.sql; do \
		echo "$(YELLOW)Running migration: $$file$(NC)"; \
		psql -U postgres -d sneakers_marketplace -f $$file; \
	done
	@echo "$(GREEN)✅ Migrations complete!$(NC)"

db-rollback: ## Rollback last migration
	@echo "$(BLUE)⏪ Rolling back migrations...$(NC)"
	@for file in internal/database/migrations/*.down.sql; do \
		echo "$(YELLOW)Rolling back: $$file$(NC)"; \
		psql -U postgres -d sneakers_marketplace -f $$file; \
	done
	@echo "$(GREEN)✅ Rollback complete!$(NC)"

db-seed: ## Seed database with test data
	@echo "$(BLUE)🌱 Seeding database...$(NC)"
	@psql -U postgres -d sneakers_marketplace -f scripts/seed.sql 2>/dev/null || echo "$(YELLOW)Seed file not found$(NC)"
	@echo "$(GREEN)✅ Database seeded!$(NC)"

db-backup: ## Backup database
	@echo "$(BLUE)💾 Backing up database...$(NC)"
	@mkdir -p backups
	@pg_dump -U postgres sneakers_marketplace > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)✅ Backup created!$(NC)"

db-restore: ## Restore database from latest backup
	@echo "$(BLUE)♻️  Restoring database...$(NC)"
	@psql -U postgres -d sneakers_marketplace -f $(shell ls -t backups/*.sql | head -1)
	@echo "$(GREEN)✅ Database restored!$(NC)"

db-reset: db-drop db-create db-migrate db-seed ## Reset database (drop, create, migrate, seed)
	@echo "$(GREEN)✅ Database reset complete!$(NC)"

#############
# TESTING
#############

test: ## Run all tests
	@echo "$(BLUE)🧪 Running tests...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Tests complete! Coverage report: coverage.html$(NC)"

test-unit: ## Run unit tests only
	@echo "$(BLUE)🧪 Running unit tests...$(NC)"
	@go test -v -short ./...

test-integration: ## Run integration tests
	@echo "$(BLUE)🧪 Running integration tests...$(NC)"
	@go test -v -run Integration ./...

test-coverage: ## Generate test coverage report
	@echo "$(BLUE)📊 Generating coverage report...$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@open coverage.html

test-frontend: ## Run frontend tests
	@echo "$(BLUE)🧪 Running frontend tests...$(NC)"
	@cd frontend && npm run test

#############
# LINTING
#############

lint: lint-backend lint-frontend ## Run all linters

lint-backend: ## Run Go linter
	@echo "$(BLUE)🔍 Running golangci-lint...$(NC)"
	@golangci-lint run --timeout=10m

lint-frontend: ## Run Frontend linter
	@echo "$(BLUE)🔍 Running ESLint...$(NC)"
	@cd frontend && npm run lint

lint-fix: ## Fix linting issues automatically
	@echo "$(BLUE)🔧 Fixing linting issues...$(NC)"
	@golangci-lint run --fix --timeout=10m
	@cd frontend && npm run lint:fix
	@echo "$(GREEN)✅ Linting issues fixed!$(NC)"

#############
# PROTO
#############

proto: ## Generate code from .proto files
	@echo "$(BLUE)🔧 Generating protobuf code...$(NC)"
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/user/*.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/product/*.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/bidding/*.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/notification/*.proto
	@echo "$(GREEN)✅ Protobuf code generated!$(NC)"

#############
# DOCKER
#############

docker-build: ## Build Docker images
	@echo "$(BLUE)🐳 Building Docker images...$(NC)"
	@docker-compose build
	@echo "$(GREEN)✅ Docker images built!$(NC)"

docker-up: ## Start services with Docker Compose
	@echo "$(BLUE)🐳 Starting Docker containers...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✅ Docker containers started!$(NC)"

docker-down: ## Stop Docker containers
	@echo "$(BLUE)🐳 Stopping Docker containers...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✅ Docker containers stopped!$(NC)"

docker-logs: ## Show Docker logs
	@docker-compose logs -f

#############
# CLEAN
#############

clean: ## Clean build artifacts and temp files
	@echo "$(BLUE)🧹 Cleaning...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f /tmp/*-service.log /tmp/*-service.pid
	@rm -f /tmp/vite.log /tmp/vite.pid
	@rm -f /tmp/mailhog.log /tmp/mailhog.pid
	@cd frontend && rm -rf dist/ node_modules/.vite
	@echo "$(GREEN)✅ Cleaned!$(NC)"

clean-all: clean ## Clean everything including dependencies
	@echo "$(BLUE)🧹 Deep cleaning...$(NC)"
	@cd frontend && rm -rf node_modules
	@go clean -modcache
	@echo "$(GREEN)✅ Deep cleaned!$(NC)"

#############
# MONITORING
#############

health: ## Check health of all services
	@echo "$(BLUE)🏥 Checking service health...$(NC)"
	@curl -s http://localhost:8080/health | jq . || echo "$(RED)API Gateway not responding$(NC)"
	@echo "$(GREEN)✅ Health check complete!$(NC)"

#############
# UTILITIES
#############

env-check: ## Check environment variables
	@echo "$(BLUE)🔍 Checking environment variables...$(NC)"
	@test -f .env && echo "$(GREEN)✅ .env file exists$(NC)" || echo "$(RED)❌ .env file missing$(NC)"
	@grep -q "JWT_SECRET" .env && echo "$(GREEN)✅ JWT_SECRET set$(NC)" || echo "$(RED)❌ JWT_SECRET missing$(NC)"
	@grep -q "DATABASE_URL" .env && echo "$(GREEN)✅ DATABASE_URL set$(NC)" || echo "$(RED)❌ DATABASE_URL missing$(NC)"

version: ## Show version information
	@echo "$(BLUE)ℹ️  Version Information:$(NC)"
	@echo "Go: $(shell go version)"
	@echo "Node: $(shell node --version)"
	@echo "npm: $(shell npm --version)"
	@echo "protoc: $(shell protoc --version)"

#############
# DEFAULT
#############

.DEFAULT_GOAL := help
