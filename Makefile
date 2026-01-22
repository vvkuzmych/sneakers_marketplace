.PHONY: help lint lint-fix test build clean

# Default target
help:
	@echo "Available commands:"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make lint-fix       - Run golangci-lint with auto-fix"
	@echo "  make frontend-lint  - Run ESLint on frontend"
	@echo "  make frontend-fix   - Run ESLint with auto-fix on frontend"
	@echo "  make test           - Run Go tests"
	@echo "  make build          - Build all services"
	@echo "  make clean          - Clean build artifacts"

# Backend linting
lint:
	@echo "🔍 Running golangci-lint..."
	golangci-lint run --timeout=10m

lint-fix:
	@echo "🔧 Running golangci-lint with auto-fix..."
	golangci-lint run --fix --timeout=10m

# Frontend linting
frontend-lint:
	@echo "🔍 Running ESLint on frontend..."
	cd frontend && npm run lint

frontend-fix:
	@echo "🔧 Running ESLint with auto-fix on frontend..."
	cd frontend && npm run lint:fix

# Testing
test:
	@echo "🧪 Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Building
build:
	@echo "🔨 Building all services..."
	./scripts/build.sh

# Cleaning
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	cd frontend && rm -rf dist/
