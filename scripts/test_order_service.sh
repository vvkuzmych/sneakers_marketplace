#!/bin/bash

set -e

echo "🧪 Testing Order Service..."
echo ""

# First, create a match using bidding service
echo "1️⃣ Creating test match via Bidding Service..."

# Place a BID
BID_RESPONSE=$(grpcurl -plaintext -d '{
  "user_id": 1,
  "product_id": 1,
  "size_id": 1,
  "price": 200,
  "quantity": 1,
  "expires_at": "'$(date -u -v+2d +%Y-%m-%dT%H:%M:%SZ)'"
}' localhost:50053 bidding.BiddingService/PlaceBid)

echo "Bid placed"

# Place an ASK that matches
ASK_RESPONSE=$(grpcurl -plaintext -d '{
  "user_id": 2,
  "product_id": 1,
  "size_id": 1,
  "price": 200,
  "quantity": 1,
  "expires_at": "'$(date -u -v+2d +%Y-%m-%dT%H:%M:%SZ)'"
}' localhost:50053 bidding.BiddingService/PlaceAsk)

echo "$ASK_RESPONSE"

# Extract match ID
if command -v jq &> /dev/null; then
    MATCH_ID=$(echo "$ASK_RESPONSE" | jq -r '.match.id // empty')
else
    MATCH_ID=$(echo "$ASK_RESPONSE" | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | sed 's/.*"\([0-9]*\)".*/\1/')
fi

if [ -z "$MATCH_ID" ] || [ "$MATCH_ID" = "null" ]; then
    echo "❌ No match created. Try running demo_all_services.sh first to see a match example."
    echo "Continuing with example Order ID 1..."
    MATCH_ID=1
else
    echo "✅ Match created with ID: $MATCH_ID"
fi

echo ""

# Create Order from Match
echo "2️⃣ Create Order from Match..."
CREATE_ORDER_RESPONSE=$(grpcurl -plaintext -d '{
  "match_id": '$MATCH_ID',
  "buyer_id": 1,
  "seller_id": 2,
  "shipping_address_id": 1
}' localhost:50054 order.OrderService/CreateOrder)

echo "$CREATE_ORDER_RESPONSE"

# Extract order ID
if command -v jq &> /dev/null; then
    ORDER_ID=$(echo "$CREATE_ORDER_RESPONSE" | jq -r '.order.id // empty')
else
    ORDER_ID=$(echo "$CREATE_ORDER_RESPONSE" | grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | sed 's/.*"\([0-9]*\)".*/\1/')
fi

if [ -z "$ORDER_ID" ] || [ "$ORDER_ID" = "null" ]; then
    echo "❌ Order not created. Using example Order ID 1..."
    ORDER_ID=1
else
    echo "✅ Order created with ID: $ORDER_ID"
fi

echo ""

# Get Order
echo "3️⃣ Get Order Details..."
grpcurl -plaintext -d '{
  "order_id": '$ORDER_ID'
}' localhost:50054 order.OrderService/GetOrder

echo ""

# Get Buyer Orders
echo "4️⃣ Get All Orders for Buyer..."
grpcurl -plaintext -d '{
  "buyer_id": 1
}' localhost:50054 order.OrderService/GetBuyerOrders

echo ""

# Get Seller Orders
echo "5️⃣ Get All Orders for Seller..."
grpcurl -plaintext -d '{
  "seller_id": 2
}' localhost:50054 order.OrderService/GetSellerOrders

echo ""

# Mark as Paid (simulated)
echo "6️⃣ Mark Order as Paid..."
grpcurl -plaintext -d '{
  "order_id": '$ORDER_ID',
  "payment_id": 1
}' localhost:50054 order.OrderService/MarkAsPaid

echo ""

# Add Tracking
echo "7️⃣ Add Shipping Tracking..."
grpcurl -plaintext -d '{
  "order_id": '$ORDER_ID',
  "tracking_number": "USPS123456789",
  "carrier": "USPS"
}' localhost:50054 order.OrderService/AddTracking

echo ""

# Get Order again (should show updated status)
echo "8️⃣ Get Order (after updates)..."
grpcurl -plaintext -d '{
  "order_id": '$ORDER_ID'
}' localhost:50054 order.OrderService/GetOrder

echo ""
echo "✅ Order Service Test Complete! 🎉"
echo ""
echo "📊 Summary:"
echo "  - Match ID: $MATCH_ID"
echo "  - Order ID: $ORDER_ID"
echo "  - Buyer ID: 1"
echo "  - Seller ID: 2"
echo ""
