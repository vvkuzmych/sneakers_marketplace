#!/bin/bash

echo "🛑 Stopping all services..."

pkill -f user-service
pkill -f product-service
pkill -f bidding-service
pkill -f order-service
pkill -f payment-service
pkill -f notification-service
pkill -f api-gateway

echo "✅ All services stopped"

# Wait a moment
sleep 1

# Check if any still running
RUNNING=$(ps aux | grep -E "user-service|product-service|bidding-service|order-service|payment-service|notification-service|api-gateway" | grep -v grep | wc -l)

if [ "$RUNNING" -eq 0 ]; then
    echo "✅ No services running"
else
    echo "⚠️  Still $RUNNING service(s) running:"
    ps aux | grep -E "user-service|product-service|bidding-service|order-service|payment-service|notification-service|api-gateway" | grep -v grep
fi
