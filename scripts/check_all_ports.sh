#!/bin/bash

echo "🔍 Checking all service ports..."
echo ""

check_port() {
    local port=$1
    local service=$2
    
    if lsof -i :$port > /dev/null 2>&1; then
        echo "✅ Port $port - $service (RUNNING)"
    else
        echo "❌ Port $port - $service (NOT RUNNING)"
    fi
}

check_port 50051 "User Service"
check_port 50052 "Product Service"
check_port 50053 "Bidding Service"
check_port 50054 "Order Service"
check_port 50055 "Payment Service"
check_port 50056 "Notification Service"
check_port 8080  "API Gateway"

echo ""
echo "📊 Summary:"
RUNNING=$(lsof -i :50051 -i :50052 -i :50053 -i :50054 -i :50055 -i :50056 -i :8080 2>/dev/null | grep LISTEN | wc -l | xargs)
echo "   Running services: $RUNNING / 7"
