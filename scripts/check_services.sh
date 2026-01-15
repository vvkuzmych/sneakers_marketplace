#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}🔍 Checking microservices status...${NC}"
echo ""

# Function to check if port is listening
check_port() {
    local SERVICE_NAME=$1
    local PORT=$2
    
    PID=$(lsof -ti :$PORT 2>/dev/null)
    if [ ! -z "$PID" ]; then
        PROCESS=$(ps -p $PID -o comm= 2>/dev/null)
        echo -e "${GREEN}✅ ${SERVICE_NAME}${NC}"
        echo "   Port: $PORT"
        echo "   PID:  $PID"
        echo "   Process: $PROCESS"
        return 0
    else
        echo -e "${RED}❌ ${SERVICE_NAME}${NC}"
        echo "   Port: $PORT"
        echo "   Status: Not running"
        return 1
    fi
}

# Check all services
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
check_port "User Service" 50051
USER_STATUS=$?
echo ""

check_port "Product Service" 50052
PRODUCT_STATUS=$?
echo ""

check_port "Bidding Service" 50053
BIDDING_STATUS=$?
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Summary
TOTAL_RUNNING=$((3 - USER_STATUS - PRODUCT_STATUS - BIDDING_STATUS))

if [ $TOTAL_RUNNING -eq 3 ]; then
    echo -e "${GREEN}✨ All services running! ($TOTAL_RUNNING/3) ✨${NC}"
elif [ $TOTAL_RUNNING -eq 0 ]; then
    echo -e "${RED}❌ No services running (0/3)${NC}"
    echo ""
    echo -e "${YELLOW}To start services:${NC}"
    echo "  ./scripts/start_all_services.sh"
else
    echo -e "${YELLOW}⚠️  Partial: $TOTAL_RUNNING/3 services running${NC}"
    echo ""
    echo -e "${YELLOW}To start all services:${NC}"
    echo "  ./scripts/start_all_services.sh"
fi

echo ""

# Check Docker infrastructure
echo -e "${CYAN}🐳 Infrastructure Status:${NC}"
echo ""

check_docker_container() {
    local CONTAINER_NAME=$1
    local STATUS=$(docker inspect -f '{{.State.Status}}' $CONTAINER_NAME 2>/dev/null)
    
    if [ -z "$STATUS" ]; then
        echo -e "  ${RED}❌${NC} $CONTAINER_NAME: not found"
    elif [ "$STATUS" = "running" ]; then
        echo -e "  ${GREEN}✅${NC} $CONTAINER_NAME: running"
    else
        echo -e "  ${YELLOW}⚠️${NC}  $CONTAINER_NAME: $STATUS"
    fi
}

# Check main infrastructure containers
check_docker_container "sneakers_postgres"
check_docker_container "sneakers_redis"
check_docker_container "sneakers_kafka"
check_docker_container "sneakers_zookeeper"

echo ""

# Check logs
echo -e "${CYAN}📝 Recent Logs:${NC}"
echo ""

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

if [ -d logs ]; then
    for SERVICE in user-service product-service bidding-service; do
        if [ -f "logs/${SERVICE}.log" ]; then
            echo -e "${CYAN}Last 3 lines from ${SERVICE}:${NC}"
            tail -3 "logs/${SERVICE}.log" 2>/dev/null | sed 's/^/  /'
            echo ""
        fi
    done
else
    echo -e "${YELLOW}  No logs directory found${NC}"
    echo ""
fi

echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
