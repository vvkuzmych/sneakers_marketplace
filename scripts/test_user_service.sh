#!/bin/bash

# Install grpcurl if not installed
# brew install grpcurl

echo "🧪 Testing User Service gRPC endpoints..."
echo ""

# 1. Register a new user
echo "1️⃣ Register new user..."
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890"
}' localhost:50051 user.UserService/Register

echo ""
echo "✅ If you see access_token and refresh_token - User Service works!"
echo ""

# 2. Login
echo "2️⃣ Login..."
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "password": "SecurePassword123!"
}' localhost:50051 user.UserService/Login

echo ""
echo "Done! 🎉"
