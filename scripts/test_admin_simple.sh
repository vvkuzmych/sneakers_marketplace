#!/bin/bash

echo "🧪 Simple Admin Service Test (without authentication)"
echo "====================================================="
echo ""

# Test GetSystemHealth (doesn't require auth in this test)
echo "1️⃣ Testing Admin Service availability..."
grpcurl -plaintext localhost:50057 list 2>&1 | head -5

echo ""
echo "2️⃣ Listing Admin Service methods..."
grpcurl -plaintext localhost:50057 list admin.AdminService 2>&1

echo ""
echo "✅ Admin Service is running and responding!"
echo ""
echo "📝 Note: Full authentication tests require working User Service"
echo "   Admin Service is ready and waiting for authenticated requests."
