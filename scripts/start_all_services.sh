#!/bin/bash

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}🚀 Starting all microservices...${NC}"
echo ""

# Get project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Create logs directory
mkdir -p logs

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo -e "${YELLOW}⚠️  Warning: .env file not found${NC}"
fi

# Check if binaries exist
if [ ! -f bin/user-service ]; then
    echo -e "${YELLOW}Building user-service...${NC}"
    go build -o bin/user-service ./cmd/user-service
fi

if [ ! -f bin/product-service ]; then
    echo -e "${YELLOW}Building product-service...${NC}"
    go build -o bin/product-service ./cmd/product-service
fi

if [ ! -f bin/bidding-service ]; then
    echo -e "${YELLOW}Building bidding-service...${NC}"
    go build -o bin/bidding-service ./cmd/bidding-service
fi

# Start User Service
echo -e "${CYAN}Starting User Service (port 50051)...${NC}"
SERVER_PORT=50051 nohup ./bin/user-service > logs/user-service.log 2>&1 &
USER_PID=$!
echo $USER_PID > logs/user-service.pid
echo -e "${GREEN}✅ User Service started (PID: $USER_PID)${NC}"
sleep 1

# Start Product Service
echo -e "${CYAN}Starting Product Service (port 50052)...${NC}"
SERVER_PORT=50052 nohup ./bin/product-service > logs/product-service.log 2>&1 &
PRODUCT_PID=$!
echo $PRODUCT_PID > logs/product-service.pid
echo -e "${GREEN}✅ Product Service started (PID: $PRODUCT_PID)${NC}"
sleep 1

# Start Bidding Service
echo -e "${CYAN}Starting Bidding Service (port 50053)...${NC}"
SERVER_PORT=50053 nohup ./bin/bidding-service > logs/bidding-service.log 2>&1 &
BIDDING_PID=$!
echo $BIDDING_PID > logs/bidding-service.pid
echo -e "${GREEN}✅ Bidding Service started (PID: $BIDDING_PID)${NC}"
sleep 2

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✨ All services started successfully! ✨${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${CYAN}Services:${NC}"
echo "  • User Service:    http://localhost:50051 (PID: $USER_PID)"
echo "  • Product Service: http://localhost:50052 (PID: $PRODUCT_PID)"
echo "  • Bidding Service: http://localhost:50053 (PID: $BIDDING_PID)"
echo ""
echo -e "${CYAN}Logs:${NC}"
echo "  • tail -f logs/user-service.log"
echo "  • tail -f logs/product-service.log"
echo "  • tail -f logs/bidding-service.log"
echo ""
echo -e "${CYAN}To stop services:${NC}"
echo "  • ./scripts/stop_all_services.sh"
echo ""
