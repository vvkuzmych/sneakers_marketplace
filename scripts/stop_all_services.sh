#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}🛑 Stopping all microservices...${NC}"
echo ""

# Get project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Function to stop service
stop_service() {
    local SERVICE_NAME=$1
    local PID_FILE="logs/${SERVICE_NAME}.pid"
    
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p $PID > /dev/null 2>&1; then
            echo -e "${YELLOW}Stopping ${SERVICE_NAME} (PID: $PID)...${NC}"
            kill $PID
            sleep 1
            
            # Check if still running
            if ps -p $PID > /dev/null 2>&1; then
                echo -e "${RED}Process still running, force killing...${NC}"
                kill -9 $PID
            fi
            
            echo -e "${GREEN}✅ ${SERVICE_NAME} stopped${NC}"
        else
            echo -e "${YELLOW}⚠️  ${SERVICE_NAME} not running (PID: $PID)${NC}"
        fi
        rm "$PID_FILE"
    else
        echo -e "${YELLOW}⚠️  No PID file found for ${SERVICE_NAME}${NC}"
    fi
}

# Stop all services
stop_service "user-service"
stop_service "product-service"
stop_service "bidding-service"

echo ""

# Also kill any processes on those ports (cleanup)
echo -e "${CYAN}Checking for remaining processes on ports...${NC}"

for PORT in 50051 50052 50053; do
    PID=$(lsof -ti :$PORT 2>/dev/null)
    if [ ! -z "$PID" ]; then
        echo -e "${YELLOW}Found process on port $PORT (PID: $PID), killing...${NC}"
        kill -9 $PID 2>/dev/null
    fi
done

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✨ All services stopped successfully! ✨${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${CYAN}To start services again:${NC}"
echo "  • ./scripts/start_all_services.sh"
echo ""
