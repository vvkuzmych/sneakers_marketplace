.PHONY: help build test clean docker-up docker-down proto migrate-up migrate-down seed

# Variables
APP_NAME=sneakers_marketplace
VERSION=$(shell git describe --tags --always --dirty)
GO=go
DOCKER_COMPOSE=docker-compose

# Colors for output
GREEN=\033[0;32m
NC=\033[0m # No Color

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Development

run-user-service: ## Run User Service
	@echo "$(GREEN)Starting User Service...$(NC)"
	cd cmd/user-service && $(GO) run main.go

run-product-service: ## Run Product Service
	@echo "$(GREEN)Starting Product Service...$(NC)"
	cd cmd/product-service && $(GO) run main.go

run-bidding-service: ## Run Bidding Service (Matching Engine)
	@echo "$(GREEN)Starting Bidding Service...$(NC)"
	cd cmd/bidding-service && $(GO) run main.go

run-order-service: ## Run Order Service
	@echo "$(GREEN)Starting Order Service...$(NC)"
	cd cmd/order-service && $(GO) run main.go

run-payment-service: ## Run Payment Service
	@echo "$(GREEN)Starting Payment Service...$(NC)"
	cd cmd/payment-service && $(GO) run main.go

run-notification-service: ## Run Notification Service
	@echo "$(GREEN)Starting Notification Service...$(NC)"
	cd cmd/notification-service && $(GO) run main.go

run-admin-service: ## Run Admin Service
	@echo "$(GREEN)Starting Admin Service...$(NC)"
	cd cmd/admin-service && $(GO) run main.go

run-search-service: ## Run Search Service
	@echo "$(GREEN)Starting Search Service...$(NC)"
	cd cmd/search-service && $(GO) run main.go

run-analytics-service: ## Run Analytics Service
	@echo "$(GREEN)Starting Analytics Service...$(NC)"
	cd cmd/analytics-service && $(GO) run main.go

run-auth-service: ## Run Authentication Service
	@echo "$(GREEN)Starting Authentication Service...$(NC)"
	cd cmd/auth-service && $(GO) run main.go

## Build

build: ## Build all services
	@echo "$(GREEN)Building all services...$(NC)"
	$(GO) build -o bin/user-service cmd/user-service/main.go
	$(GO) build -o bin/product-service cmd/product-service/main.go
	$(GO) build -o bin/bidding-service cmd/bidding-service/main.go
	$(GO) build -o bin/order-service cmd/order-service/main.go
	$(GO) build -o bin/payment-service cmd/payment-service/main.go
	$(GO) build -o bin/notification-service cmd/notification-service/main.go
	$(GO) build -o bin/admin-service cmd/admin-service/main.go
	$(GO) build -o bin/api-gateway cmd/api-gateway/main.go
	@echo "$(GREEN)✅ All services built successfully!$(NC)"

build-user: ## Build User Service
	$(GO) build -o bin/user-service cmd/user-service/main.go

build-product: ## Build Product Service
	$(GO) build -o bin/product-service cmd/product-service/main.go

build-bidding: ## Build Bidding Service
	$(GO) build -o bin/bidding-service cmd/bidding-service/main.go

build-admin: ## Build Admin Service
	$(GO) build -o bin/admin-service cmd/admin-service/main.go

## Testing

test: ## Run all tests
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO) test -v -race -cover ./...

test-coverage: ## Run tests with coverage
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(NC)"
	$(GO) test -v -race -tags=integration ./tests/integration/...

test-e2e: ## Run end-to-end tests
	@echo "$(GREEN)Running E2E tests...$(NC)"
	$(GO) test -v -tags=e2e ./tests/e2e/...

## Database

migrate-up: ## Run database migrations up
	@echo "$(GREEN)Running migrations up...$(NC)"
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down: ## Run database migrations down
	@echo "$(GREEN)Running migrations down...$(NC)"
	migrate -path migrations -database "${DATABASE_URL}" down

migrate-create: ## Create new migration (use name=your_migration_name)
	@echo "$(GREEN)Creating migration: $(name)$(NC)"
	migrate create -ext sql -dir migrations -seq $(name)

seed: ## Seed database with test data
	@echo "$(GREEN)Seeding database...$(NC)"
	$(GO) run scripts/seed/main.go

## Docker

docker-up: ## Start all services with Docker Compose
	@echo "$(GREEN)Starting Docker services...$(NC)"
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop all Docker services
	@echo "$(GREEN)Stopping Docker services...$(NC)"
	$(DOCKER_COMPOSE) down

docker-logs: ## View Docker logs
	$(DOCKER_COMPOSE) logs -f

docker-build: ## Build Docker images
	@echo "$(GREEN)Building Docker images...$(NC)"
	$(DOCKER_COMPOSE) build

docker-clean: ## Remove all containers, volumes, and images
	@echo "$(GREEN)Cleaning Docker...$(NC)"
	$(DOCKER_COMPOSE) down -v --rmi all

## Protobuf

proto: ## Generate code from .proto files
	@echo "$(GREEN)Generating protobuf code...$(NC)"
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/**/*.proto

proto-install: ## Install protoc plugins
	@echo "$(GREEN)Installing protoc plugins...$(NC)"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

## Dependencies

deps: ## Download Go dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GO) mod download

deps-tidy: ## Tidy Go dependencies
	@echo "$(GREEN)Tidying dependencies...$(NC)"
	$(GO) mod tidy

deps-vendor: ## Vendor dependencies
	@echo "$(GREEN)Vendoring dependencies...$(NC)"
	$(GO) mod vendor

## Linting

lint: ## Run linter
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run ./...

fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GO) fmt ./...

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GO) vet ./...

## Monitoring

prometheus: ## Open Prometheus dashboard
	@echo "$(GREEN)Opening Prometheus...$(NC)"
	open http://localhost:9090

grafana: ## Open Grafana dashboard
	@echo "$(GREEN)Opening Grafana...$(NC)"
	open http://localhost:3000

jaeger: ## Open Jaeger dashboard
	@echo "$(GREEN)Opening Jaeger...$(NC)"
	open http://localhost:16686

## Cleanup

clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning...$(NC)"
	rm -rf bin/
	rm -rf vendor/
	rm -f coverage.out coverage.html
	$(GO) clean

clean-all: clean docker-clean ## Clean everything including Docker

## Git

git-init: ## Initialize git repository
	@echo "$(GREEN)Initializing git repository...$(NC)"
	git init
	git add .
	git commit -m "Initial commit: Project structure"
	git branch -M main

git-push: ## Push to GitHub (first time)
	@echo "$(GREEN)Pushing to GitHub...$(NC)"
	git remote add origin https://github.com/vvkuzmych/sneakers_marketplace.git
	git push -u origin main

## Project Setup

setup: ## Initial project setup
	@echo "$(GREEN)Setting up project...$(NC)"
	$(MAKE) deps
	$(MAKE) proto-install
	@echo "$(GREEN)Creating directories...$(NC)"
	mkdir -p bin logs tmp
	@echo "$(GREEN)Setup complete!$(NC)"
	@echo "$(GREEN)Next steps:$(NC)"
	@echo "  1. Copy .env.example to .env and configure"
	@echo "  2. Run 'make docker-up' to start infrastructure"
	@echo "  3. Run 'make migrate-up' to create database tables"
	@echo "  4. Run 'make seed' to populate test data"

## Documentation

docs: ## Generate documentation
	@echo "$(GREEN)Generating documentation...$(NC)"
	godoc -http=:6060

.DEFAULT_GOAL := help
