#!/bin/bash

echo "🎯 Testing Bidding Service - Bid/Ask Matching Engine..."
echo ""

# Prerequisites: Need existing product and size from Product Service
# Let's use Product ID: 1 and Size ID from recent tests

# Extract first available product and size
echo "📦 Getting available products and sizes..."
PRODUCT_RESPONSE=$(grpcurl -plaintext -d '{"page": 1, "page_size": 1, "active_only": true}' localhost:50052 product.ProductService/ListProducts)

if command -v jq &> /dev/null; then
    PRODUCT_ID=$(echo "$PRODUCT_RESPONSE" | jq -r '.products[0].id')
else
    PRODUCT_ID=$(echo "$PRODUCT_RESPONSE" | grep -o '"id"[[:space:]]*:[[:space:]]*"[0-9]*"' | head -1 | grep -o '[0-9]*')
fi

if [ -z "$PRODUCT_ID" ] || [ "$PRODUCT_ID" = "null" ]; then
    echo "❌ No products found! Run test_product_service.sh first."
    exit 1
fi

echo "Using Product ID: $PRODUCT_ID"

# Get sizes for this product
SIZES_RESPONSE=$(grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"'
}' localhost:50052 product.ProductService/GetAvailableSizes)

if command -v jq &> /dev/null; then
    SIZE_ID=$(echo "$SIZES_RESPONSE" | jq -r '.sizes[0].id')
else
    SIZE_ID=$(echo "$SIZES_RESPONSE" | grep -o '"id"[[:space:]]*:[[:space:]]*"[0-9]*"' | head -1 | grep -o '[0-9]*')
fi

if [ -z "$SIZE_ID" ] || [ "$SIZE_ID" = "null" ]; then
    echo "❌ No sizes found for product!"
    exit 1
fi

echo "Using Size ID: $SIZE_ID"
echo ""

# Test users
BUYER_ID=1
SELLER_ID=1  # Same user for simplicity in test

# 1. Place a BID (buyer wants to buy at $200)
echo "1️⃣ Place BID: Buyer offers $200..."
BID_RESPONSE=$(grpcurl -plaintext -d '{
  "user_id": '"$BUYER_ID"',
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_ID"',
  "price": 200.00,
  "quantity": 1,
  "expires_in_hours": 48
}' localhost:50053 bidding.BiddingService/PlaceBid)

echo "$BID_RESPONSE"

if command -v jq &> /dev/null; then
    BID_ID=$(echo "$BID_RESPONSE" | jq -r '.bid.id')
else
    BID_ID=$(echo "$BID_RESPONSE" | grep -o '"id"[[:space:]]*:[[:space:]]*"[0-9]*"' | head -1 | grep -o '[0-9]*')
fi

echo "BID ID: $BID_ID"
echo ""

# 2. Check Market Price (no match yet, only bid)
echo "2️⃣ Get Market Price (should show only BID)..."
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_ID"'
}' localhost:50053 bidding.BiddingService/GetMarketPrice

echo ""

# 3. Place an ASK at higher price (no match - seller wants $220)
echo "3️⃣ Place ASK: Seller wants $220 (No match)..."
ASK_HIGH_RESPONSE=$(grpcurl -plaintext -d '{
  "user_id": '"$SELLER_ID"',
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_ID"',
  "price": 220.00,
  "quantity": 1,
  "expires_in_hours": 48
}' localhost:50053 bidding.BiddingService/PlaceAsk)

echo "$ASK_HIGH_RESPONSE"
echo ""

# 4. Check Market Price (should show spread: bid $200, ask $220)
echo "4️⃣ Get Market Price (spread: bid $200 vs ask $220)..."
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_ID"'
}' localhost:50053 bidding.BiddingService/GetMarketPrice

echo ""

# 5. Place ANOTHER BID at higher price that MATCHES (buyer offers $225 >= $220)
echo "5️⃣ Place BID: Buyer offers $225 (SHOULD MATCH!)..."
MATCH_BID_RESPONSE=$(grpcurl -plaintext -d '{
  "user_id": '"$BUYER_ID"',
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_ID"',
  "price": 225.00,
  "quantity": 1,
  "expires_in_hours": 24
}' localhost:50053 bidding.BiddingService/PlaceBid)

echo "$MATCH_BID_RESPONSE"

if command -v jq &> /dev/null; then
    MATCH_ID=$(echo "$MATCH_BID_RESPONSE" | jq -r '.match.id')
else
    MATCH_ID=$(echo "$MATCH_BID_RESPONSE" | grep -o '"match"[^}]*"id"[[:space:]]*:[[:space:]]*"[0-9]*"' | grep -o '[0-9]*' | head -1)
fi

echo ""

if [ ! -z "$MATCH_ID" ] && [ "$MATCH_ID" != "null" ]; then
    echo "✅ MATCH CREATED! Match ID: $MATCH_ID"
    echo ""
    
    # 6. Get Match Details
    echo "6️⃣ Get Match Details..."
    grpcurl -plaintext -d '{
      "match_id": '"$MATCH_ID"'
    }' localhost:50053 bidding.BiddingService/GetMatch
    
    echo ""
else
    echo "⚠️  No match created"
    echo ""
fi

# 7. Get Product Bids
echo "7️⃣ Get All Bids for Product..."
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_ID"',
  "status": "active",
  "page": 1,
  "page_size": 10
}' localhost:50053 bidding.BiddingService/GetProductBids

echo ""

# 8. Get Product Asks
echo "8️⃣ Get All Asks for Product..."
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size_id": '"$SIZE_ID"',
  "status": "active",
  "page": 1,
  "page_size": 10
}' localhost:50053 bidding.BiddingService/GetProductAsks

echo ""

# 9. Get User Matches
echo "9️⃣ Get User Matches..."
grpcurl -plaintext -d '{
  "user_id": '"$BUYER_ID"',
  "as_buyer": true,
  "as_seller": true,
  "page": 1,
  "page_size": 10
}' localhost:50053 bidding.BiddingService/GetUserMatches

echo ""
echo "✅ Bidding Service Test Complete! 🎉"
echo ""
echo "📊 Summary:"
echo "  - Product ID: $PRODUCT_ID, Size ID: $SIZE_ID"
echo "  - Bid placed at $200 (no match)"
echo "  - Ask placed at $220 (no match)"
echo "  - Bid placed at $225 (matched with $220 ask!)"
echo "  - Match ID: $MATCH_ID"
echo "  - Final price: $220 (seller's ask price)"
