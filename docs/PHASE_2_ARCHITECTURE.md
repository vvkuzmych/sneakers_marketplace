# 🏗️ Phase 2 - Architecture Design

**Version:** 1.0  
**Date:** 2026-01-15  
**Status:** Planning

---

## 📋 Overview

Phase 2 extends the Sneakers Marketplace with order processing, payment integration, and HTTP API Gateway.

**New Services:**
1. **Order Service** - Order lifecycle management
2. **Payment Service** - Stripe payment processing
3. **API Gateway** - HTTP REST interface

---

## 🎯 Goals

- ✅ Process matched bids/asks into orders
- ✅ Handle payments with Stripe
- ✅ Provide HTTP REST API for clients
- ✅ Track order status and shipping
- ✅ Handle refunds and cancellations
- ✅ Maintain audit trail

---

## 🔄 Data Flow

```
┌─────────────┐
│   Client    │
│  (Web/App)  │
└──────┬──────┘
       │ HTTP REST
       ▼
┌─────────────────┐
│  API Gateway    │  ← New!
│   (HTTP/gRPC)   │
└────────┬────────┘
         │ gRPC
    ┬────┼────┬────┬
    │    │    │    │
    ▼    ▼    ▼    ▼
┌──────┐┌────┐┌────┐┌────────┐
│User  ││Prod││Bid ││ Order  │ ← New!
│Svc   ││Svc ││Svc ││  Svc   │
└──┬───┘└──┬─┘└─┬──┘└───┬────┘
   │       │    │ match │
   │       │    │ event │
   │       │    └───────┤
   │       │            ▼
   │       │      ┌──────────┐
   │       │      │ Payment  │ ← New!
   │       │      │   Svc    │
   │       │      └────┬─────┘
   │       │           │
   │       │           │ Stripe
   └───┬───┴───┬───────┴─┐
       │       │         │
       ▼       ▼         ▼
   ┌────────────────────────┐
   │     PostgreSQL         │
   └────────────────────────┘
```

---

## 📦 1. Order Service

### Purpose
Convert matched bids/asks into trackable orders with full lifecycle management.

### Features
- Create order from match
- Order status tracking
- Shipping address management
- Order history and audit trail
- Order cancellation and refunds
- Seller and buyer views

### Database Schema

#### `orders` table
```sql
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,  -- e.g., "ORD-2026-001234"
    match_id BIGINT NOT NULL REFERENCES matches(id),
    
    -- Parties
    buyer_id BIGINT NOT NULL REFERENCES users(id),
    seller_id BIGINT NOT NULL REFERENCES users(id),
    
    -- Product details
    product_id BIGINT NOT NULL REFERENCES products(id),
    size_id BIGINT NOT NULL REFERENCES sizes(id),
    
    -- Pricing
    price DECIMAL(10, 2) NOT NULL,  -- Final agreed price
    quantity INT NOT NULL DEFAULT 1,
    
    -- Fees (marketplace takes a cut)
    buyer_fee DECIMAL(10, 2) DEFAULT 0,      -- Processing fee for buyer
    seller_fee DECIMAL(10, 2) DEFAULT 0,     -- Commission from seller
    platform_fee DECIMAL(10, 2) DEFAULT 0,   -- Platform fee
    
    total_amount DECIMAL(10, 2) NOT NULL,    -- price + buyer_fee
    seller_payout DECIMAL(10, 2) NOT NULL,   -- price - seller_fee
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending_payment',
    -- Statuses: pending_payment, paid, processing, shipped, 
    --           delivered, completed, cancelled, refunded
    
    -- Shipping
    shipping_address_id BIGINT REFERENCES addresses(id),
    tracking_number VARCHAR(100),
    carrier VARCHAR(50),  -- UPS, FedEx, USPS, DHL
    
    -- Timestamps
    payment_at TIMESTAMP,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    
    -- Notes
    buyer_notes TEXT,
    seller_notes TEXT,
    admin_notes TEXT,
    cancellation_reason TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_match_id ON orders(match_id);
CREATE INDEX idx_orders_buyer_id ON orders(buyer_id);
CREATE INDEX idx_orders_seller_id ON orders(seller_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
```

#### `order_status_history` table
```sql
CREATE TABLE order_status_history (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    note TEXT,
    created_by VARCHAR(50),  -- system, buyer, seller, admin
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id, created_at DESC);
```

### Order Lifecycle

```
1. Match Created (Bidding Service)
   ↓
2. Order Created (status: pending_payment)
   ↓
3. Payment Processed (Payment Service)
   → status: paid
   ↓
4. Seller Notified
   → status: processing
   ↓
5. Seller Ships Product
   → status: shipped (+ tracking number)
   ↓
6. Buyer Receives
   → status: delivered
   ↓
7. Buyer Confirms (auto after 7 days)
   → status: completed
   → Seller receives payout
```

### gRPC Methods

```protobuf
service OrderService {
    // Order management
    rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
    rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
    rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
    rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (UpdateOrderStatusResponse);
    rpc CancelOrder(CancelOrderRequest) returns (CancelOrderResponse);
    
    // Shipping
    rpc AddTrackingNumber(AddTrackingNumberRequest) returns (AddTrackingNumberResponse);
    rpc GetShippingStatus(GetShippingStatusRequest) returns (GetShippingStatusResponse);
    
    // Buyer/Seller views
    rpc GetBuyerOrders(GetBuyerOrdersRequest) returns (GetBuyerOrdersResponse);
    rpc GetSellerOrders(GetSellerOrdersRequest) returns (GetSellerOrdersResponse);
    
    // Status history
    rpc GetOrderStatusHistory(GetOrderStatusHistoryRequest) returns (GetOrderStatusHistoryResponse);
}
```

---

## 💳 2. Payment Service

### Purpose
Handle all payment processing with Stripe integration.

### Features
- Create payment intent (Stripe)
- Process payment
- Handle webhooks
- Refund processing
- Payment history
- Stripe Connect for seller payouts

### Database Schema

#### `payments` table
```sql
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    payment_id VARCHAR(100) UNIQUE NOT NULL,  -- Internal ID
    
    order_id BIGINT NOT NULL REFERENCES orders(id),
    user_id BIGINT NOT NULL REFERENCES users(id),  -- Buyer
    
    -- Stripe details
    stripe_payment_intent_id VARCHAR(255) UNIQUE,
    stripe_charge_id VARCHAR(255),
    stripe_customer_id VARCHAR(255),
    
    -- Amount
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- Statuses: pending, processing, succeeded, failed, 
    --           cancelled, refunded, partially_refunded
    
    -- Payment method
    payment_method VARCHAR(50),  -- card, apple_pay, google_pay
    card_last4 VARCHAR(4),
    card_brand VARCHAR(20),  -- visa, mastercard, amex
    
    -- Refund
    refunded_amount DECIMAL(10, 2) DEFAULT 0,
    refund_reason TEXT,
    
    -- Timestamps
    processed_at TIMESTAMP,
    refunded_at TIMESTAMP,
    
    -- Metadata
    metadata JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_stripe_payment_intent_id ON payments(stripe_payment_intent_id);
CREATE INDEX idx_payments_status ON payments(status);
```

#### `payouts` table (for sellers)
```sql
CREATE TABLE payouts (
    id BIGSERIAL PRIMARY KEY,
    payout_id VARCHAR(100) UNIQUE NOT NULL,
    
    order_id BIGINT NOT NULL REFERENCES orders(id),
    seller_id BIGINT NOT NULL REFERENCES users(id),
    
    -- Stripe Connect
    stripe_transfer_id VARCHAR(255),
    stripe_account_id VARCHAR(255),  -- Seller's Stripe Connect account
    
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- Statuses: pending, processing, paid, failed, reversed
    
    processed_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payouts_order_id ON payouts(order_id);
CREATE INDEX idx_payouts_seller_id ON payouts(seller_id);
CREATE INDEX idx_payouts_status ON payouts(status);
```

### Stripe Flow

```
1. User places bid/ask → match created
2. Order created → Payment Service called
3. Create Stripe Payment Intent
   → Client secret returned to frontend
4. Frontend confirms payment with Stripe.js
5. Stripe webhook → Payment Service
6. Update payment status → Update order status
7. After delivery confirmed:
   → Create Stripe Transfer to seller
```

### gRPC Methods

```protobuf
service PaymentService {
    // Payment processing
    rpc CreatePaymentIntent(CreatePaymentIntentRequest) returns (CreatePaymentIntentResponse);
    rpc ConfirmPayment(ConfirmPaymentRequest) returns (ConfirmPaymentResponse);
    rpc GetPayment(GetPaymentRequest) returns (GetPaymentResponse);
    rpc ListPayments(ListPaymentsRequest) returns (ListPaymentsResponse);
    
    // Refunds
    rpc CreateRefund(CreateRefundRequest) returns (CreateRefundResponse);
    rpc GetRefund(GetRefundRequest) returns (GetRefundResponse);
    
    // Webhooks
    rpc HandleStripeWebhook(HandleStripeWebhookRequest) returns (HandleStripeWebhookResponse);
    
    // Payouts (Stripe Connect)
    rpc CreatePayout(CreatePayoutRequest) returns (CreatePayoutResponse);
    rpc GetPayout(GetPayoutRequest) returns (GetPayoutResponse);
    rpc ListPayouts(ListPayoutsRequest) returns (ListPayoutsResponse);
}
```

---

## 🌐 3. API Gateway

### Purpose
Provide HTTP REST API for web/mobile clients, translating to gRPC calls.

### Features
- HTTP REST endpoints
- JWT authentication middleware
- Rate limiting
- Request validation
- Swagger/OpenAPI documentation
- CORS handling
- Response caching (Redis)

### Tech Stack
- **Gin** or **Fiber** (HTTP framework)
- **go-swagger** (OpenAPI generation)
- **rate** (rate limiting)
- **gRPC client** (call microservices)

### REST Endpoints

```
Authentication:
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout

Users:
GET    /api/v1/users/profile
PUT    /api/v1/users/profile
GET    /api/v1/users/addresses
POST   /api/v1/users/addresses
PUT    /api/v1/users/addresses/:id
DELETE /api/v1/users/addresses/:id

Products:
GET    /api/v1/products
GET    /api/v1/products/:id
GET    /api/v1/products/search?q=nike
GET    /api/v1/products/:id/sizes

Bidding:
POST   /api/v1/bids
GET    /api/v1/bids/:id
GET    /api/v1/bids
DELETE /api/v1/bids/:id
POST   /api/v1/asks
GET    /api/v1/asks/:id
GET    /api/v1/asks
DELETE /api/v1/asks/:id
GET    /api/v1/market/:product_id/:size_id

Orders:
GET    /api/v1/orders
GET    /api/v1/orders/:id
PUT    /api/v1/orders/:id/cancel
PUT    /api/v1/orders/:id/tracking
GET    /api/v1/orders/buyer
GET    /api/v1/orders/seller

Payments:
POST   /api/v1/payments/intent
POST   /api/v1/payments/confirm
POST   /api/v1/payments/:id/refund
POST   /api/v1/webhooks/stripe  (public)
```

### Middleware Stack

```
Request
  ↓
[CORS]
  ↓
[Logger]
  ↓
[Rate Limiter]
  ↓
[JWT Auth] (except public routes)
  ↓
[Request Validation]
  ↓
[Handler] → gRPC Call
  ↓
[Response Formatter]
  ↓
Response
```

---

## 📊 Service Communication

### Event-Driven (Future)

For Phase 2, we'll use direct gRPC calls. In Phase 3, we can introduce Kafka events:

```
Bidding Service → "match.created" event
                    ↓
                Order Service listens
                    ↓
                Creates order
                    ↓
                Emits "order.created" event
                    ↓
        ┌───────────┴────────────┐
        ▼                        ▼
Payment Service          Notification Service
  (process payment)        (email buyer/seller)
```

---

## 🔒 Security

### API Gateway
- JWT verification on all protected routes
- Rate limiting (100 req/min per IP)
- Request size limits
- SQL injection prevention (prepared statements)
- XSS protection (sanitize inputs)

### Payment Service
- Stripe webhook signature verification
- PCI compliance (no card data stored)
- TLS/SSL required
- Stripe Connect for seller payouts

---

## 🧪 Testing Strategy

### Unit Tests
- Models, repositories, services (each service)

### Integration Tests
- API Gateway → gRPC services
- Payment Service → Stripe sandbox

### E2E Tests
- Full flow: Register → Bid → Match → Order → Payment → Ship

---

## 📈 Performance

### Caching (Redis)
- Product catalog (5 min TTL)
- Market prices (10 sec TTL)
- User profiles (1 min TTL)

### Database
- Connection pooling (25 per service)
- Read replicas (future)
- Indexes on all foreign keys

---

## 🚀 Deployment (Phase 3)

- Docker images for each service
- Kubernetes manifests
- Horizontal pod autoscaling
- Health checks and liveness probes

---

## 📋 Phase 2 Milestones

### Milestone 1: Order Service (Week 2)
- [x] Database migration
- [ ] gRPC proto
- [ ] Models & repository
- [ ] Business logic
- [ ] gRPC handler
- [ ] Test scripts

### Milestone 2: Payment Service (Week 3)
- [ ] Database migration
- [ ] Stripe SDK integration
- [ ] Payment intent creation
- [ ] Webhook handling
- [ ] Refund logic
- [ ] Test scripts (Stripe sandbox)

### Milestone 3: API Gateway (Week 4)
- [ ] HTTP server setup
- [ ] gRPC client connections
- [ ] REST endpoints
- [ ] JWT middleware
- [ ] Swagger docs
- [ ] Integration tests

---

## ✅ Success Criteria

Phase 2 is complete when:
- ✅ Orders are created automatically from matches
- ✅ Payments can be processed with Stripe
- ✅ Refunds work correctly
- ✅ HTTP REST API functional for all services
- ✅ Swagger docs generated
- ✅ All test scripts pass
- ✅ E2E flow works: Bid → Match → Order → Payment

---

**Next:** Start with Order Service implementation! 🚀
