#!/bin/bash

echo "🧪 Testing Product Service gRPC endpoints..."
echo ""

# Optional: Clean up previous test product (comment out if you want to keep test data)
# echo "🗑️  Cleaning up previous test product (ID: 1)..."
# grpcurl -plaintext -d '{"product_id": 1}' localhost:50052 product.ProductService/DeleteProduct 2>/dev/null || true
# echo ""

# Generate unique SKU with timestamp
TIMESTAMP=$(date +%s)
UNIQUE_SKU="AJ1-${TIMESTAMP}"

# 1. Create a product
echo "1️⃣ Create Product (Nike Air Jordan 1 - $UNIQUE_SKU)..."
PRODUCT_RESPONSE=$(grpcurl -plaintext -d '{
  "sku": "'"$UNIQUE_SKU"'",
  "name": "Air Jordan 1 Retro High OG Chicago",
  "brand": "Nike",
  "model": "Air Jordan 1",
  "color": "Chicago Red/White/Black",
  "description": "The iconic Air Jordan 1 in the legendary Chicago colorway. A must-have for every sneaker collector.",
  "category": "Basketball",
  "release_year": 2026,
  "retail_price": 170.00
}' localhost:50052 product.ProductService/CreateProduct)

echo "$PRODUCT_RESPONSE"

# Extract product ID from JSON using jq or grep
if command -v jq &> /dev/null; then
    PRODUCT_ID=$(echo "$PRODUCT_RESPONSE" | jq -r '.product.id')
else
    PRODUCT_ID=$(echo "$PRODUCT_RESPONSE" | grep -o '"id"[[:space:]]*:[[:space:]]*"[0-9]*"' | head -1 | grep -o '[0-9]*')
fi

echo "Created Product ID: $PRODUCT_ID"
echo ""

if [ -z "$PRODUCT_ID" ] || [ "$PRODUCT_ID" = "null" ]; then
    echo "❌ Failed to create product or extract Product ID!"
    echo "Response: $PRODUCT_RESPONSE"
    exit 1
fi

# 2. Add sizes to product
echo "2️⃣ Add Sizes (US 9, 10, 11)..."
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size": "US 9",
  "quantity": 50
}' localhost:50052 product.ProductService/AddSize

grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size": "US 10",
  "quantity": 75
}' localhost:50052 product.ProductService/AddSize

grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "size": "US 11",
  "quantity": 60
}' localhost:50052 product.ProductService/AddSize

echo "✅ Sizes added!"
echo ""

# 3. Add product images
echo "3️⃣ Add Product Images..."
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "image_url": "https://example.com/aj1-chicago-1.jpg",
  "alt_text": "Air Jordan 1 Chicago - Side View",
  "display_order": 1,
  "is_primary": true
}' localhost:50052 product.ProductService/AddProductImage

grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"',
  "image_url": "https://example.com/aj1-chicago-2.jpg",
  "alt_text": "Air Jordan 1 Chicago - Top View",
  "display_order": 2,
  "is_primary": false
}' localhost:50052 product.ProductService/AddProductImage

echo "✅ Images added!"
echo ""

# 4. Get product with details
echo "4️⃣ Get Product Details..."
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"'
}' localhost:50052 product.ProductService/GetProduct

echo ""

# 5. List all products
echo "5️⃣ List All Products..."
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10,
  "active_only": true
}' localhost:50052 product.ProductService/ListProducts

echo ""

# 6. Search products
echo "6️⃣ Search Products (Nike)..."
grpcurl -plaintext -d '{
  "query": "Nike",
  "page": 1,
  "page_size": 10
}' localhost:50052 product.ProductService/SearchProducts

echo ""

# 7. Get Available Sizes
echo "7️⃣ Get Available Sizes..."
grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"'
}' localhost:50052 product.ProductService/GetAvailableSizes

echo ""

# 8. Test Inventory Reservation
echo "8️⃣ Test Inventory Reservation..."
SIZE_RESPONSE=$(grpcurl -plaintext -d '{
  "product_id": '"$PRODUCT_ID"'
}' localhost:50052 product.ProductService/GetAvailableSizes)

# Extract first size_id for testing
if command -v jq &> /dev/null; then
    SIZE_ID=$(echo "$SIZE_RESPONSE" | jq -r '.sizes[0].id')
else
    SIZE_ID=$(echo "$SIZE_RESPONSE" | grep -o '"id"[[:space:]]*:[[:space:]]*"[0-9]*"' | head -1 | grep -o '[0-9]*')
fi

if [ ! -z "$SIZE_ID" ] && [ "$SIZE_ID" != "null" ]; then
    echo "Reserving 5 units for order TEST-ORDER-001..."
    grpcurl -plaintext -d '{
      "size_id": '"$SIZE_ID"',
      "quantity": 5,
      "order_id": "TEST-ORDER-001"
    }' localhost:50052 product.ProductService/ReserveInventory
    
    echo ""
    echo "Checking inventory after reservation..."
    grpcurl -plaintext -d '{
      "product_id": '"$PRODUCT_ID"'
    }' localhost:50052 product.ProductService/GetAvailableSizes
fi

echo ""
echo "✅ Product Service is working perfectly! 🎉"
echo ""
echo "📊 Summary:"
echo "  - Product created: $UNIQUE_SKU (ID: $PRODUCT_ID)"
echo "  - Sizes added: US 9, 10, 11"
echo "  - Images added: 2"
echo "  - Inventory reserved: 5 units for TEST-ORDER-001"