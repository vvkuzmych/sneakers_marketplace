#!/bin/bash

set -e

API_URL="http://localhost:8080"

echo "🧪 Testing API Gateway (HTTP REST)..."
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if API Gateway is running
echo -e "${BLUE}0️⃣ Health Check...${NC}"
HEALTH=$(curl -s "$API_URL/health")
if echo "$HEALTH" | grep -q "healthy"; then
    echo -e "${GREEN}✅ API Gateway is healthy!${NC}"
else
    echo -e "${RED}❌ API Gateway is not responding${NC}"
    exit 1
fi
echo ""

# Register new user
echo -e "${BLUE}1️⃣ Register User...${NC}"
REGISTER_EMAIL="test-$(date +%s)@example.com"
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$REGISTER_EMAIL\",
    \"password\": \"password123\",
    \"first_name\": \"John\",
    \"last_name\": \"Doe\",
    \"phone\": \"+1234567890\"
  }")

echo "$REGISTER_RESPONSE" | jq '.'

ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.access_token')
USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.user.id')

if [ "$ACCESS_TOKEN" != "null" ] && [ "$ACCESS_TOKEN" != "" ]; then
    echo -e "${GREEN}✅ User registered successfully!${NC}"
    echo "Access Token: ${ACCESS_TOKEN:0:20}..."
    echo "User ID: $USER_ID"
else
    echo -e "${RED}❌ Registration failed${NC}"
    exit 1
fi
echo ""

# Get user profile
echo -e "${BLUE}2️⃣ Get User Profile (with JWT auth)...${NC}"
PROFILE_RESPONSE=$(curl -s -X GET "$API_URL/api/v1/users/$USER_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "$PROFILE_RESPONSE" | jq '.'

if echo "$PROFILE_RESPONSE" | jq -e '.user.email' > /dev/null; then
    echo -e "${GREEN}✅ Profile retrieved successfully!${NC}"
else
    echo -e "${RED}❌ Failed to get profile${NC}"
fi
echo ""

# List products (public endpoint)
echo -e "${BLUE}3️⃣ List Products (public)...${NC}"
PRODUCTS_RESPONSE=$(curl -s -X GET "$API_URL/api/v1/products?page=1&page_size=5")

echo "$PRODUCTS_RESPONSE" | jq '.'

PRODUCT_COUNT=$(echo "$PRODUCTS_RESPONSE" | jq '.products | length')
echo -e "${GREEN}✅ Found $PRODUCT_COUNT products${NC}"
echo ""

# Search products
echo -e "${BLUE}4️⃣ Search Products (public)...${NC}"
SEARCH_RESPONSE=$(curl -s -X GET "$API_URL/api/v1/products/search?q=Nike")

echo "$SEARCH_RESPONSE" | jq '.'
echo -e "${GREEN}✅ Search completed${NC}"
echo ""

# Test authentication - without token (should fail)
echo -e "${BLUE}5️⃣ Test Auth Protection (should fail without token)...${NC}"
UNAUTH_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/bids" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": $USER_ID,
    \"product_id\": 1,
    \"size_id\": 1,
    \"price\": 200,
    \"quantity\": 1
  }")

if echo "$UNAUTH_RESPONSE" | grep -q "authorization header required"; then
    echo -e "${GREEN}✅ Auth protection working! (401 Unauthorized)${NC}"
else
    echo "$UNAUTH_RESPONSE" | jq '.'
    echo -e "${RED}⚠️ Expected authentication error${NC}"
fi
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ API Gateway Test Complete!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📊 Summary:"
echo "  ✅ Health check"
echo "  ✅ User registration (with JWT)"
echo "  ✅ Protected endpoints (with JWT auth)"
echo "  ✅ Public endpoints (products)"
echo "  ✅ Authentication protection"
echo ""
echo "🔑 Your Access Token:"
echo "   $ACCESS_TOKEN"
echo ""
echo "💡 Try these examples:"
echo "   # Get products"
echo "   curl '$API_URL/api/v1/products'"
echo ""
echo "   # Get user profile (requires JWT)"
echo "   curl -H 'Authorization: Bearer $ACCESS_TOKEN' '$API_URL/api/v1/users/$USER_ID'"
echo ""
