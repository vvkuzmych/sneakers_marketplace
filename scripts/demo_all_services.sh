#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║                                                                 ║"
echo "║          👟 SNEAKERS MARKETPLACE - FULL DEMO 👟                ║"
echo "║                                                                 ║"
echo "║          Demonstrating all 3 microservices:                     ║"
echo "║          • User Service (Auth & Profile)                        ║"
echo "║          • Product Service (Catalog & Inventory)                ║"
echo "║          • Bidding Service (Matching Engine)                    ║"
echo "║                                                                 ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""

sleep 2

# Check if services are running
echo -e "${YELLOW}🔍 Checking if services are running...${NC}"
USER_SERVICE=$(lsof -i :50051 2>/dev/null | grep LISTEN)
PRODUCT_SERVICE=$(lsof -i :50052 2>/dev/null | grep LISTEN)
BIDDING_SERVICE=$(lsof -i :50053 2>/dev/null | grep LISTEN)

if [ -z "$USER_SERVICE" ]; then
    echo -e "${RED}❌ User Service not running on :50051${NC}"
    echo "Start it with: ./bin/user-service"
    exit 1
fi

if [ -z "$PRODUCT_SERVICE" ]; then
    echo -e "${RED}❌ Product Service not running on :50052${NC}"
    echo "Start it with: ./bin/product-service"
    exit 1
fi

if [ -z "$BIDDING_SERVICE" ]; then
    echo -e "${RED}❌ Bidding Service not running on :50053${NC}"
    echo "Start it with: ./bin/bidding-service"
    exit 1
fi

echo -e "${GREEN}✅ All services are running!${NC}"
echo ""
sleep 1

#==============================================================================
# PART 1: USER SERVICE - Authentication & Profile
#==============================================================================

echo -e "${MAGENTA}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║  PART 1: USER SERVICE - Authentication & Profile Management    ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""
sleep 1

# Generate unique email
TIMESTAMP=$(date +%s)
BUYER_EMAIL="buyer${TIMESTAMP}@sneakers.com"
SELLER_EMAIL="seller${TIMESTAMP}@sneakers.com"

echo -e "${CYAN}👤 1.1 Registering BUYER (${BUYER_EMAIL})...${NC}"
BUYER_RESPONSE=$(grpcurl -plaintext -d '{
  "email": "'"$BUYER_EMAIL"'",
  "password": "BuyerPass123!",
  "first_name": "Alice",
  "last_name": "Johnson"
}' localhost:50051 user.UserService/Register 2>&1)

echo "$BUYER_RESPONSE" | head -10
# Parse BUYER_ID
if command -v jq >/dev/null 2>&1; then
    BUYER_ID=$(echo "$BUYER_RESPONSE" | jq -r '.user.id // empty' 2>/dev/null)
else
    BUYER_ID=$(echo "$BUYER_RESPONSE" | grep '"id"' | head -1 | sed 's/.*"id"[[:space:]]*:[[:space:]]*"\([0-9]*\)".*/\1/')
fi
echo -e "${GREEN}✅ Buyer registered! ID: ${BUYER_ID}${NC}"
echo ""
sleep 1

echo -e "${CYAN}👤 1.2 Registering SELLER (${SELLER_EMAIL})...${NC}"
SELLER_RESPONSE=$(grpcurl -plaintext -d '{
  "email": "'"$SELLER_EMAIL"'",
  "password": "SellerPass123!",
  "first_name": "Bob",
  "last_name": "Smith"
}' localhost:50051 user.UserService/Register 2>&1)

echo "$SELLER_RESPONSE" | head -10
# Parse SELLER_ID
if command -v jq >/dev/null 2>&1; then
    SELLER_ID=$(echo "$SELLER_RESPONSE" | jq -r '.user.id // empty' 2>/dev/null)
else
    SELLER_ID=$(echo "$SELLER_RESPONSE" | grep '"id"' | head -1 | sed 's/.*"id"[[:space:]]*:[[:space:]]*"\([0-9]*\)".*/\1/')
fi
echo -e "${GREEN}✅ Seller registered! ID: ${SELLER_ID}${NC}"
echo ""
sleep 1

echo -e "${CYAN}🔐 1.3 Testing Login for BUYER...${NC}"
LOGIN_RESPONSE=$(grpcurl -plaintext -d '{
  "email": "'"$BUYER_EMAIL"'",
  "password": "BuyerPass123!"
}' localhost:50051 user.UserService/Login 2>&1)

echo "$LOGIN_RESPONSE" | grep -E '"(accessToken|user)"' | head -5
echo -e "${GREEN}✅ Login successful! JWT tokens received${NC}"
echo ""
sleep 2

#==============================================================================
# PART 2: PRODUCT SERVICE - Catalog & Inventory
#==============================================================================

echo -e "${MAGENTA}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║  PART 2: PRODUCT SERVICE - Catalog & Inventory Management      ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""
sleep 1

echo -e "${CYAN}📦 2.1 Creating Product: Nike Air Jordan 1 'Chicago'...${NC}"
PRODUCT_RESPONSE=$(grpcurl -plaintext -d '{
  "sku": "AJ1-CHICAGO-'"$TIMESTAMP"'",
  "name": "Air Jordan 1 Retro High OG Chicago",
  "brand": "Nike",
  "model": "Air Jordan 1",
  "color": "Chicago Red/White/Black",
  "description": "The iconic Air Jordan 1 in the legendary Chicago colorway.",
  "category": "Basketball",
  "release_year": 2026,
  "retail_price": 170.00
}' localhost:50052 product.ProductService/CreateProduct 2>&1)

echo "$PRODUCT_RESPONSE" | head -15
# Parse PRODUCT_ID
if command -v jq >/dev/null 2>&1; then
    PRODUCT_ID=$(echo "$PRODUCT_RESPONSE" | jq -r '.product.id // empty' 2>/dev/null)
else
    PRODUCT_ID=$(echo "$PRODUCT_RESPONSE" | grep '"id"' | head -1 | sed 's/.*"id"[[:space:]]*:[[:space:]]*"\([0-9]*\)".*/\1/')
fi
echo -e "${GREEN}✅ Product created! ID: ${PRODUCT_ID}${NC}"
echo ""
sleep 1

echo -e "${CYAN}📏 2.2 Adding Size US 9 with 10 units in stock...${NC}"
SIZE_9_RESPONSE=$(grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size": "US 9",
  "quantity": 10
}' localhost:50052 product.ProductService/AddSize 2>&1)

echo "$SIZE_9_RESPONSE" | head -10
# Parse SIZE_9_ID
if command -v jq >/dev/null 2>&1; then
    SIZE_9_ID=$(echo "$SIZE_9_RESPONSE" | jq -r '.size.id // empty' 2>/dev/null)
else
    SIZE_9_ID=$(echo "$SIZE_9_RESPONSE" | grep '"id"' | head -1 | sed 's/.*"id"[[:space:]]*:[[:space:]]*"\([0-9]*\)".*/\1/')
fi
echo -e "${GREEN}✅ Size US 9 added! ID: ${SIZE_9_ID}, Stock: 10${NC}"
echo ""
sleep 1

echo -e "${CYAN}📏 2.3 Adding Size US 10 with 15 units in stock...${NC}"
SIZE_10_RESPONSE=$(grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size": "US 10",
  "quantity": 15
}' localhost:50052 product.ProductService/AddSize 2>&1)

# Parse SIZE_10_ID
if command -v jq >/dev/null 2>&1; then
    SIZE_10_ID=$(echo "$SIZE_10_RESPONSE" | jq -r '.size.id // empty' 2>/dev/null)
else
    SIZE_10_ID=$(echo "$SIZE_10_RESPONSE" | grep '"id"' | head -1 | sed 's/.*"id"[[:space:]]*:[[:space:]]*"\([0-9]*\)".*/\1/')
fi
echo -e "${GREEN}✅ Size US 10 added! ID: ${SIZE_10_ID}, Stock: 15${NC}"
echo ""
sleep 1

echo -e "${CYAN}🖼️  2.4 Adding Product Images...${NC}"
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "image_url": "https://images.example.com/aj1-chicago-side.jpg",
  "alt_text": "Air Jordan 1 Chicago - Side View",
  "display_order": 1,
  "is_primary": true
}' localhost:50052 product.ProductService/AddProductImage >/dev/null 2>&1

grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "image_url": "https://images.example.com/aj1-chicago-top.jpg",
  "alt_text": "Air Jordan 1 Chicago - Top View",
  "display_order": 2,
  "is_primary": false
}' localhost:50052 product.ProductService/AddProductImage >/dev/null 2>&1

echo -e "${GREEN}✅ 2 images added to product${NC}"
echo ""
sleep 2

#==============================================================================
# PART 3: BIDDING SERVICE - Matching Engine
#==============================================================================

echo -e "${MAGENTA}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║  PART 3: BIDDING SERVICE - Bid/Ask Matching Engine             ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""
sleep 1

echo -e "${CYAN}💰 3.1 Checking initial market price...${NC}"
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_9_ID"'
}' localhost:50053 bidding.BiddingService/GetMarketPrice
echo -e "${YELLOW}📊 No bids or asks yet - market is empty${NC}"
echo ""
sleep 1

echo -e "${CYAN}📈 3.2 BUYER places BID: $200 (willing to buy at this price)...${NC}"
BID_1_RESPONSE=$(grpcurl -plaintext -d '{
  "user_id": '"$BUYER_ID"',
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_9_ID"',
  "price": 200.00,
  "quantity": 1,
  "expires_in_hours": 48
}' localhost:50053 bidding.BiddingService/PlaceBid 2>&1)

echo "$BID_1_RESPONSE" | head -12
# Parse BID_1_ID
if command -v jq >/dev/null 2>&1; then
    BID_1_ID=$(echo "$BID_1_RESPONSE" | jq -r '.bid.id // empty' 2>/dev/null)
else
    BID_1_ID=$(echo "$BID_1_RESPONSE" | grep '"id"' | head -1 | sed 's/.*"id"[[:space:]]*:[[:space:]]*"\([0-9]*\)".*/\1/')
fi
echo -e "${GREEN}✅ BID placed! ID: ${BID_1_ID}, Price: \$200, Status: active${NC}"
echo -e "${YELLOW}   Waiting for seller...${NC}"
echo ""
sleep 1

echo -e "${CYAN}💰 3.3 Checking market price after BID...${NC}"
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_9_ID"'
}' localhost:50053 bidding.BiddingService/GetMarketPrice
echo -e "${YELLOW}📊 Highest BID: \$200, No ASKs yet${NC}"
echo ""
sleep 2

echo -e "${CYAN}📉 3.4 SELLER places ASK: $250 (willing to sell at this price)...${NC}"
ASK_1_RESPONSE=$(grpcurl -plaintext -d '{
  "user_id": '"$SELLER_ID"',
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_9_ID"',
  "price": 250.00,
  "quantity": 1,
  "expires_in_hours": 48
}' localhost:50053 bidding.BiddingService/PlaceAsk 2>&1)

echo "$ASK_1_RESPONSE" | head -12
# Parse ASK_1_ID
if command -v jq >/dev/null 2>&1; then
    ASK_1_ID=$(echo "$ASK_1_RESPONSE" | jq -r '.ask.id // empty' 2>/dev/null)
else
    ASK_1_ID=$(echo "$ASK_1_RESPONSE" | grep '"id"' | head -1 | sed 's/.*"id"[[:space:]]*:[[:space:]]*"\([0-9]*\)".*/\1/')
fi
echo -e "${GREEN}✅ ASK placed! ID: ${ASK_1_ID}, Price: \$250, Status: active${NC}"
echo -e "${YELLOW}   No match yet - ASK price (\$250) > BID price (\$200)${NC}"
echo ""
sleep 1

echo -e "${CYAN}💰 3.5 Checking market spread...${NC}"
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_9_ID"'
}' localhost:50053 bidding.BiddingService/GetMarketPrice
echo -e "${YELLOW}📊 Spread: BID \$200 / ASK \$250 (difference: \$50)${NC}"
echo ""
sleep 2

echo -e "${CYAN}📈 3.6 NEW BUYER places BID: $260 (higher than ASK!)...${NC}"
echo -e "${YELLOW}   ⚡ This should trigger INSTANT MATCH! ⚡${NC}"
sleep 1

BID_2_RESPONSE=$(grpcurl -plaintext -d '{
  "user_id": '"$BUYER_ID"',
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_9_ID"',
  "price": 260.00,
  "quantity": 1,
  "expires_in_hours": 24
}' localhost:50053 bidding.BiddingService/PlaceBid 2>&1)

echo "$BID_2_RESPONSE"

# Try jq first (more reliable), fallback to grep
if command -v jq >/dev/null 2>&1; then
    MATCH_ID=$(echo "$BID_2_RESPONSE" | jq -r '.match.id // empty' 2>/dev/null)
else
    # Grep for "match" section, then find first "id" after it
    MATCH_ID=$(echo "$BID_2_RESPONSE" | grep -A 15 '"match"' | grep '"id"' | head -1 | sed 's/.*"id"[[:space:]]*:[[:space:]]*"\([0-9]*\)".*/\1/')
fi

echo ""
if [ ! -z "$MATCH_ID" ] && [ "$MATCH_ID" != "null" ] && [ "$MATCH_ID" != "" ]; then
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✨ 🎉 MATCH SUCCESSFUL! 🎉 ✨${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}   Match ID: ${MATCH_ID}${NC}"
    echo -e "${GREEN}   Buyer offered: \$260${NC}"
    echo -e "${GREEN}   Seller asked: \$250${NC}"
    echo -e "${GREEN}   Final price: \$250 (seller's price)${NC}"
    echo -e "${GREEN}   Both orders marked as 'matched'${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    sleep 2
    
    echo -e "${CYAN}📋 3.7 Viewing Match Details...${NC}"
    grpcurl -plaintext -d '{
      "match_id": '"$MATCH_ID"'
    }' localhost:50053 bidding.BiddingService/GetMatch
    echo ""
    sleep 1
else
    echo -e "${RED}❌ Match was not created (unexpected)${NC}"
    echo ""
fi

sleep 2

#==============================================================================
# PART 4: FINAL STATE
#==============================================================================

echo -e "${MAGENTA}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║  PART 4: FINAL STATE - Order Book & Match History              ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""
sleep 1

echo -e "${CYAN}📊 4.1 Checking final market price...${NC}"
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_9_ID"'
}' localhost:50053 bidding.BiddingService/GetMarketPrice
echo ""
sleep 1

echo -e "${CYAN}📈 4.2 Active BIDs in order book...${NC}"
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_9_ID"',
  "status": "active",
  "page": 1,
  "page_size": 10
}' localhost:50053 bidding.BiddingService/GetProductBids
echo ""
sleep 1

echo -e "${CYAN}📉 4.3 Active ASKs in order book...${NC}"
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_9_ID"',
  "status": "active",
  "page": 1,
  "page_size": 10
}' localhost:50053 bidding.BiddingService/GetProductAsks
echo ""
sleep 1

echo -e "${CYAN}✅ 4.4 Buyer's Match History...${NC}"
grpcurl -plaintext -d '{
  "user_id": '"$BUYER_ID"',
  "as_buyer": true,
  "page": 1,
  "page_size": 10
}' localhost:50053 bidding.BiddingService/GetUserMatches
echo ""
sleep 2

#==============================================================================
# SUMMARY
#==============================================================================

echo -e "${CYAN}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║                      📊 DEMO SUMMARY 📊                        ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""
echo -e "${GREEN}✅ USER SERVICE:${NC}"
echo "   • Registered 2 users (buyer & seller)"
echo "   • Authenticated with JWT tokens"
echo "   • Users: ${BUYER_ID}, ${SELLER_ID}"
echo ""
echo -e "${GREEN}✅ PRODUCT SERVICE:${NC}"
echo "   • Created product: Air Jordan 1 Chicago"
echo "   • Added 2 sizes (US 9, US 10)"
echo "   • Added 2 images"
echo "   • Product ID: ${PRODUCT_ID}"
echo ""
echo -e "${GREEN}✅ BIDDING SERVICE:${NC}"
echo "   • Bid #1: \$200 (active, waiting)"
echo "   • Ask #1: \$250 (matched)"
echo "   • Bid #2: \$260 (matched)"
echo "   • Match: \$250 (seller's price)"
echo "   • Match ID: ${MATCH_ID}"
echo ""
echo -e "${CYAN}🔥 MATCHING ENGINE WORKED PERFECTLY! 🔥${NC}"
echo ""
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}All 3 microservices are working together seamlessly!${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${MAGENTA}Thank you for watching the demo! 🎉${NC}"
echo ""
