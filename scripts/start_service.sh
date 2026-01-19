#!/bin/bash

# Usage: ./scripts/start_service.sh user-service

SERVICE_NAME=$1

if [ -z "$SERVICE_NAME" ]; then
    echo "Usage: $0 <service-name>"
    echo "Example: $0 user-service"
    echo "Available services:"
    echo "  - user-service (port 50051)"
    echo "  - product-service (port 50052)"
    echo "  - bidding-service (port 50053)"
    echo "  - order-service (port 50054)"
    echo "  - payment-service (port 50055)"
    echo "  - notification-service (port 50056)"
    echo "  - api-gateway (port 8080)"
    exit 1
fi

cd /Users/vkuzm/GolandProjects/sneakers_marketplace

# Export common environment variables
export DATABASE_URL="postgres://postgres:postgres@localhost:5435/sneakers_marketplace?sslmode=disable"
export REDIS_URL="redis://localhost:6380/0"
export JWT_SECRET="your-super-secret-key-change-in-production"
export STRIPE_MODE="demo"

# Load additional vars from .env if exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | grep -v '^$' | xargs)
fi

# Set service-specific port environment variable
case "$SERVICE_NAME" in
    "user-service")
        export USER_SERVICE_PORT=50051
        echo "🚀 Starting User Service on port 50051..."
        ;;
    "product-service")
        export PRODUCT_SERVICE_PORT=50052
        echo "🚀 Starting Product Service on port 50052..."
        ;;
    "bidding-service")
        export BIDDING_SERVICE_PORT=50053
        echo "🚀 Starting Bidding Service on port 50053..."
        ;;
    "order-service")
        export ORDER_SERVICE_PORT=50054
        echo "🚀 Starting Order Service on port 50054..."
        ;;
    "payment-service")
        export PAYMENT_SERVICE_PORT=50055
        echo "🚀 Starting Payment Service on port 50055..."
        ;;
    "notification-service")
        export NOTIFICATION_SERVICE_PORT=50056
        echo "🚀 Starting Notification Service on port 50056..."
        ;;
    "api-gateway")
        export HTTP_PORT=8080
        echo "🚀 Starting API Gateway on port 8080..."
        ;;
    *)
        echo "❌ Unknown service: $SERVICE_NAME"
        exit 1
        ;;
esac

./bin/$SERVICE_NAME
