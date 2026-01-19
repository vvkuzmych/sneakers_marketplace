#!/bin/bash

# ═══════════════════════════════════════════════════════════════════
# GORM vs Raw SQL Comparison Test Runner
# ═══════════════════════════════════════════════════════════════════

set -e  # Exit on error

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo ""
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║                                                                  ║"
echo "║       🔬 GORM vs Raw SQL - Comparison Test Runner               ║"
echo "║                                                                  ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""

# ═══════════════════════════════════════════════════════════════════
# Menu
# ═══════════════════════════════════════════════════════════════════

echo -e "${BOLD}Виберіть режим:${NC}"
echo ""
echo "  1) 🎬 Demo - Повна демонстрація (рекомендовано)"
echo "  2) ⚡ Benchmarks - Тести продуктивності"
echo "  3) 🔥 ВСЕ - Demo + Benchmarks"
echo "  4) 🛠️  Build - Тільки збірка"
echo ""
echo -e "${YELLOW}Або запусти напряму:${NC}"
echo "  ./run_comparison.sh demo       - Demo"
echo "  ./run_comparison.sh bench      - Benchmarks"
echo "  ./run_comparison.sh all        - Все"
echo ""

# Parse command line argument or ask
if [ $# -eq 0 ]; then
    read -p "Введи номер (1-4): " choice
else
    case "$1" in
        demo|d|1)
            choice=1
            ;;
        bench|benchmark|b|2)
            choice=2
            ;;
        all|a|3)
            choice=3
            ;;
        build|4)
            choice=4
            ;;
        *)
            echo -e "${RED}❌ Невідома опція: $1${NC}"
            exit 1
            ;;
    esac
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# ═══════════════════════════════════════════════════════════════════
# Check Database Connection
# ═══════════════════════════════════════════════════════════════════

check_database() {
    echo -e "${BLUE}🔍 Перевіряю підключення до бази даних...${NC}"
    
    DB_URL="${DATABASE_URL:-postgresql://postgres:postgres@localhost:5432/sneakers_marketplace?sslmode=disable}"
    
    # Extract host and port
    HOST=$(echo "$DB_URL" | sed -n 's/.*@\(.*\):.*/\1/p')
    PORT=$(echo "$DB_URL" | sed -n 's/.*:\([0-9]*\)\/.*/\1/p')
    
    if command -v psql &> /dev/null; then
        if psql "$DB_URL" -c "SELECT 1" &> /dev/null; then
            echo -e "${GREEN}✅ База даних доступна${NC}"
        else
            echo -e "${YELLOW}⚠️  Не можу підключитись до PostgreSQL${NC}"
            echo -e "${YELLOW}   Переконайся, що база запущена: docker-compose up -d${NC}"
            read -p "Продовжити? (y/n): " continue
            if [[ ! "$continue" =~ ^[Yy]$ ]]; then
                exit 1
            fi
        fi
    else
        echo -e "${YELLOW}⚠️  psql не встановлено, пропускаю перевірку БД${NC}"
    fi
    echo ""
}

# ═══════════════════════════════════════════════════════════════════
# Run Demo
# ═══════════════════════════════════════════════════════════════════

run_demo() {
    echo "╔══════════════════════════════════════════════════════════════════╗"
    echo "║                     🎬 Running Demo                              ║"
    echo "╚══════════════════════════════════════════════════════════════════╝"
    echo ""
    
    echo -e "${BLUE}📦 Building demo...${NC}"
    if go build -o demo main.go; then
        echo -e "${GREEN}✅ Build successful${NC}"
        echo ""
        
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${BOLD}                    🎯 DEMO OUTPUT                              ${NC}"
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        
        ./demo
        
        echo ""
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${GREEN}✅ Demo completed successfully!${NC}"
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    else
        echo -e "${RED}❌ Build failed${NC}"
        exit 1
    fi
}

# ═══════════════════════════════════════════════════════════════════
# Run Benchmarks
# ═══════════════════════════════════════════════════════════════════

run_benchmarks() {
    echo ""
    echo "╔══════════════════════════════════════════════════════════════════╗"
    echo "║                  ⚡ Running Benchmarks                            ║"
    echo "╚══════════════════════════════════════════════════════════════════╝"
    echo ""
    
    echo -e "${BLUE}🏃 Starting benchmark tests...${NC}"
    echo -e "${YELLOW}(це може зайняти 30-60 секунд)${NC}"
    echo ""
    
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BOLD}                  📊 BENCHMARK RESULTS                          ${NC}"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    # Run benchmarks and save output
    BENCH_OUTPUT=$(mktemp)
    if go test -bench=. -benchmem -benchtime=2s 2>&1 | tee "$BENCH_OUTPUT"; then
        echo ""
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${BOLD}                  📈 PERFORMANCE SUMMARY                        ${NC}"
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        
        # Parse and compare results
        echo -e "${BLUE}🔍 Аналіз результатів:${NC}"
        echo ""
        
        # Extract CREATE benchmarks
        CREATE_RAW=$(grep "BenchmarkCreate_RawSQL" "$BENCH_OUTPUT" | awk '{print $3}')
        CREATE_GORM=$(grep "BenchmarkCreate_GORM" "$BENCH_OUTPUT" | awk '{print $3}')
        
        if [ -n "$CREATE_RAW" ] && [ -n "$CREATE_GORM" ]; then
            echo -e "  ${BOLD}CREATE:${NC}"
            echo -e "    Raw SQL:  ${GREEN}${CREATE_RAW}${NC}"
            echo -e "    GORM:     ${YELLOW}${CREATE_GORM}${NC}"
            
            # Calculate overhead (simple integer division)
            RAW_NS=$(echo "$CREATE_RAW" | sed 's/[^0-9]//g')
            GORM_NS=$(echo "$CREATE_GORM" | sed 's/[^0-9]//g')
            if [ -n "$RAW_NS" ] && [ -n "$GORM_NS" ] && [ "$RAW_NS" -gt 0 ]; then
                OVERHEAD=$((GORM_NS * 100 / RAW_NS - 100))
                echo -e "    Overhead: ${RED}+${OVERHEAD}%${NC}"
            fi
            echo ""
        fi
        
        # Extract GET benchmarks
        GET_RAW=$(grep "BenchmarkGetByEmail_RawSQL" "$BENCH_OUTPUT" | awk '{print $3}')
        GET_GORM=$(grep "BenchmarkGetByEmail_GORM" "$BENCH_OUTPUT" | awk '{print $3}')
        
        if [ -n "$GET_RAW" ] && [ -n "$GET_GORM" ]; then
            echo -e "  ${BOLD}GET BY EMAIL:${NC}"
            echo -e "    Raw SQL:  ${GREEN}${GET_RAW}${NC}"
            echo -e "    GORM:     ${YELLOW}${GET_GORM}${NC}"
            
            RAW_NS=$(echo "$GET_RAW" | sed 's/[^0-9]//g')
            GORM_NS=$(echo "$GET_GORM" | sed 's/[^0-9]//g')
            if [ -n "$RAW_NS" ] && [ -n "$GORM_NS" ] && [ "$RAW_NS" -gt 0 ]; then
                OVERHEAD=$((GORM_NS * 100 / RAW_NS - 100))
                echo -e "    Overhead: ${RED}+${OVERHEAD}%${NC}"
            fi
            echo ""
        fi
        
        # Extract UPDATE benchmarks
        UPDATE_RAW=$(grep "BenchmarkUpdate_RawSQL" "$BENCH_OUTPUT" | awk '{print $3}')
        UPDATE_GORM=$(grep "BenchmarkUpdate_GORM" "$BENCH_OUTPUT" | awk '{print $3}')
        
        if [ -n "$UPDATE_RAW" ] && [ -n "$UPDATE_GORM" ]; then
            echo -e "  ${BOLD}UPDATE:${NC}"
            echo -e "    Raw SQL:  ${GREEN}${UPDATE_RAW}${NC}"
            echo -e "    GORM:     ${YELLOW}${UPDATE_GORM}${NC}"
            
            RAW_NS=$(echo "$UPDATE_RAW" | sed 's/[^0-9]//g')
            GORM_NS=$(echo "$UPDATE_GORM" | sed 's/[^0-9]//g')
            if [ -n "$RAW_NS" ] && [ -n "$GORM_NS" ] && [ "$RAW_NS" -gt 0 ]; then
                OVERHEAD=$((GORM_NS * 100 / RAW_NS - 100))
                echo -e "    Overhead: ${RED}+${OVERHEAD}%${NC}"
            fi
            echo ""
        fi
        
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${GREEN}✅ Benchmarks completed successfully!${NC}"
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        
        # Cleanup
        rm -f "$BENCH_OUTPUT"
    else
        echo -e "${RED}❌ Benchmarks failed${NC}"
        rm -f "$BENCH_OUTPUT"
        exit 1
    fi
}

# ═══════════════════════════════════════════════════════════════════
# Main
# ═══════════════════════════════════════════════════════════════════

case $choice in
    1)
        check_database
        run_demo
        ;;
    2)
        check_database
        run_benchmarks
        ;;
    3)
        check_database
        run_demo
        run_benchmarks
        ;;
    4)
        echo -e "${BLUE}🛠️  Building...${NC}"
        if go build -o demo main.go; then
            echo -e "${GREEN}✅ Build successful: ./demo${NC}"
        else
            echo -e "${RED}❌ Build failed${NC}"
            exit 1
        fi
        ;;
    *)
        echo -e "${RED}❌ Невірний вибір${NC}"
        exit 1
        ;;
esac

echo ""
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║                     ✅ All Done!                                 ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""
echo -e "${BLUE}📚 Більше інформації:${NC}"
echo "  • README.md - детальні інструкції"
echo "  • ../../docs/GORM_QUICKSTART.md - швидкий старт"
echo "  • ../../docs/GORM_INVESTIGATION.md - повний аналіз"
echo ""
echo -e "${YELLOW}💡 Tip: Запусти './run_comparison.sh all' для повного тесту${NC}"
echo ""
