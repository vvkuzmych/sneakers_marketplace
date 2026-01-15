# 🌐 API Gateway Documentation

**Base URL:** `http://localhost:8080`  
**Version:** v1  
**Protocol:** HTTP REST (proxies to gRPC services)

---

## 🔑 Authentication

Most endpoints require JWT authentication via `Authorization` header:

```bash
Authorization: Bearer <your_jwt_token>
```

Get your token via `/api/v1/auth/register` or `/api/v1/auth/login`.

---

## 📋 API Endpoints

### Health Check

```bash
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "api-gateway"
}
```

---

## 👤 Authentication & Users

### Register User

```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890"
}
```

**Response:**
```json
{
  "user": {
    "id": "1",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "is_active": true,
    "created_at": "2026-01-15T12:00:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1...",
  "refresh_token": "eyJhbGciOiJIUzI1..."
}
```

---

### Login User

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** Same as Register

---

### Get User Profile

🔒 **Protected** - Requires JWT

```bash
GET /api/v1/users/{user_id}
Authorization: Bearer <token>
```

**Response:**
```json
{
  "user": {
    "id": "1",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "is_active": true,
    "created_at": "2026-01-15T12:00:00Z",
    "updated_at": "2026-01-15T12:00:00Z"
  }
}
```

---

## 👟 Products

### List Products

✅ **Public**

```bash
GET /api/v1/products?page=1&page_size=10
```

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 10, max: 100)

**Response:**
```json
{
  "products": [
    {
      "id": "1",
      "sku": "AJ1-001-2026",
      "name": "Air Jordan 1 Retro High OG Chicago",
      "brand": "Nike",
      "model": "Air Jordan 1",
      "color": "Chicago Red/White/Black",
      "description": "The iconic Air Jordan 1...",
      "category": "Basketball",
      "release_year": "2026",
      "retail_price": 170,
      "is_active": true,
      "created_at": "2026-01-15T12:00:00Z"
    }
  ],
  "total": "25",
  "page": 1,
  "page_size": 10
}
```

---

### Get Product by ID

✅ **Public**

```bash
GET /api/v1/products/{id}
```

**Response:**
```json
{
  "product": {
    "id": "1",
    "sku": "AJ1-001-2026",
    "name": "Air Jordan 1 Retro High OG Chicago",
    "brand": "Nike",
    "model": "Air Jordan 1",
    "color": "Chicago Red/White/Black",
    "description": "The iconic Air Jordan 1...",
    "category": "Basketball",
    "release_year": "2026",
    "retail_price": 170,
    "is_active": true,
    "images": [
      {
        "id": "1",
        "image_url": "https://example.com/aj1-1.jpg",
        "is_primary": true,
        "display_order": 1
      }
    ],
    "sizes": [
      {
        "id": "1",
        "size": "US 9",
        "quantity": 10,
        "reserved": 0
      }
    ]
  }
}
```

---

### Search Products

✅ **Public**

```bash
GET /api/v1/products/search?q=Nike
```

**Query Parameters:**
- `q` (required) - Search query

**Response:**
```json
{
  "products": [...],
  "total": "15"
}
```

---

## 💰 Bidding & Market

### Place Bid (Buy Order)

🔒 **Protected** - Requires JWT

```bash
POST /api/v1/bids
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 1,
  "product_id": 1,
  "size_id": 1,
  "price": 200,
  "quantity": 1,
  "expires_at": "2026-01-17T12:00:00Z"
}
```

**Response:**
```json
{
  "bid": {
    "id": "1",
    "user_id": "1",
    "product_id": "1",
    "size_id": "1",
    "price": 200,
    "quantity": 1,
    "status": "active",
    "expires_at": "2026-01-17T12:00:00Z",
    "created_at": "2026-01-15T12:00:00Z"
  },
  "match": null
}
```

**If matched immediately:**
```json
{
  "bid": {...},
  "match": {
    "id": "1",
    "bid_id": "1",
    "ask_id": "2",
    "buyer_id": "1",
    "seller_id": "2",
    "product_id": "1",
    "size_id": "1",
    "price": 195,
    "quantity": 1,
    "status": "pending"
  }
}
```

---

### Place Ask (Sell Order)

🔒 **Protected** - Requires JWT

```bash
POST /api/v1/asks
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 2,
  "product_id": 1,
  "size_id": 1,
  "price": 220,
  "quantity": 1,
  "expires_at": "2026-01-17T12:00:00Z"
}
```

**Response:** Similar to Place Bid

---

### Get Market Price

✅ **Public**

```bash
GET /api/v1/market/{product_id}/{size_id}
```

**Example:**
```bash
GET /api/v1/market/1/1
```

**Response:**
```json
{
  "highest_bid": 200,
  "lowest_ask": 220,
  "total_bids": "5",
  "total_asks": "3"
}
```

---

## 📦 Orders

### Get Order

🔒 **Protected** - Requires JWT

```bash
GET /api/v1/orders/{id}
Authorization: Bearer <token>
```

**Response:**
```json
{
  "order": {
    "id": "1",
    "order_number": "ORD-20260115-0001",
    "match_id": "1",
    "buyer_id": "1",
    "seller_id": "2",
    "product_id": "1",
    "size_id": "1",
    "price": 200,
    "buyer_fee": 10,
    "seller_fee": 8,
    "total_amount": 210,
    "status": "pending",
    "created_at": "2026-01-15T12:00:00Z"
  }
}
```

---

### List Buyer Orders

🔒 **Protected** - Requires JWT

```bash
GET /api/v1/orders/buyer/{buyer_id}
Authorization: Bearer <token>
```

**Response:**
```json
{
  "orders": [
    {
      "id": "1",
      "order_number": "ORD-20260115-0001",
      "status": "pending",
      "total_amount": 210,
      "created_at": "2026-01-15T12:00:00Z"
    }
  ],
  "total": "3"
}
```

---

## 💳 Payments

### Create Payment Intent

🔒 **Protected** - Requires JWT

```bash
POST /api/v1/payments/intent
Authorization: Bearer <token>
Content-Type: application/json

{
  "order_id": 1,
  "user_id": 1,
  "amount": 210,
  "currency": "USD"
}
```

**Response:**
```json
{
  "payment": {
    "id": "1",
    "order_id": "1",
    "user_id": "1",
    "amount": 210,
    "currency": "USD",
    "status": "pending",
    "created_at": "2026-01-15T12:00:00Z"
  },
  "client_secret": "pi_3ABC123..._secret_xyz"
}
```

> **Note:** `client_secret` is used for Stripe frontend integration

---

### Get Payment

🔒 **Protected** - Requires JWT

```bash
GET /api/v1/payments/{id}
Authorization: Bearer <token>
```

**Response:**
```json
{
  "payment": {
    "id": "1",
    "order_id": "1",
    "user_id": "1",
    "amount": 210,
    "currency": "USD",
    "status": "succeeded",
    "payment_method": "card",
    "card_last4": "4242",
    "card_brand": "visa",
    "created_at": "2026-01-15T12:00:00Z"
  }
}
```

---

## 🔒 Authentication Examples

### Register and Get Token

```bash
# 1. Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  }'

# Save the access_token from response
```

### Use Token for Protected Endpoints

```bash
# 2. Use token
TOKEN="eyJhbGciOiJIUzI1..."

curl http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🚨 Error Responses

### 400 Bad Request
```json
{
  "error": "invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "authorization header required"
}
```

### 404 Not Found
```json
{
  "error": "resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal server error"
}
```

---

## 🧪 Testing

Run the test script:

```bash
./scripts/test_api_gateway.sh
```

This will test:
- ✅ Health check
- ✅ User registration & login
- ✅ JWT authentication
- ✅ Public endpoints (products)
- ✅ Protected endpoints (bids, orders, payments)

---

## 🏗️ Architecture

```
HTTP REST         gRPC Services
┌─────────┐      ┌──────────────────┐
│         │      │  User Service    │
│         │──────│  :50051          │
│         │      └──────────────────┘
│   API   │      ┌──────────────────┐
│ Gateway │──────│  Product Service │
│ :8080   │      │  :50052          │
│         │      └──────────────────┘
│         │      ┌──────────────────┐
│         │──────│  Bidding Service │
│         │      │  :50053          │
│         │      └──────────────────┘
│         │      ┌──────────────────┐
│         │──────│  Order Service   │
│         │      │  :50054          │
└─────────┘      └──────────────────┘
                 ┌──────────────────┐
                 │  Payment Service │
                 │  :50055          │
                 └──────────────────┘
```

**Key Features:**
- 🔐 JWT-based authentication
- 🔓 Public & protected endpoints
- 🌐 CORS support
- ⚡ HTTP/REST to gRPC translation
- 🛡️ Centralized auth middleware
- 📊 Request logging

---

**Made with ❤️ for Sneakers Marketplace**
