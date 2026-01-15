# 🏗️ Architecture - Sneakers Marketplace

## Overview

Sneakers Marketplace використовує **мікросервісну архітектуру** з event-driven communication та gRPC для synchronous calls.

---

## 🎯 Architectural Principles

1. **Service Independence** - Кожен сервіс має власну БД
2. **Domain-Driven Design** - Bounded contexts
3. **API First** - gRPC для internal, REST для external
4. **Event-Driven** - Async communication via Kafka
5. **CQRS** - Separation of reads and writes where needed
6. **Circuit Breaker** - Resilience patterns
7. **Observability** - Metrics, logs, traces everywhere

---

## 🔄 Communication Patterns

### Synchronous (gRPC)
Used when immediate response required:
```
Order Service → Product Service (check stock)
Order Service → Payment Service (process payment)
Cart → Product Service (get prices)
```

### Asynchronous (Kafka)
Used for fire-and-forget:
```
Order Service → "order.created" → Notification Service
                                → Analytics Service
                                → Search Service
```

---

## 📦 Services Deep Dive

### 1. User Service

**Purpose:** User management, authentication, profiles

**Technology:**
- Go + Gin (HTTP)
- PostgreSQL (users, addresses, sessions)
- JWT for auth
- bcrypt for passwords

**API:**
```protobuf
service UserService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc GetProfile(GetProfileRequest) returns (GetProfileResponse);
  rpc UpdateProfile(UpdateProfileRequest) returns (UpdateProfileResponse);
  rpc AddAddress(AddAddressRequest) returns (AddAddressResponse);
  rpc GetAddresses(GetAddressesRequest) returns (GetAddressesResponse);
}
```

**Database:**
- `users` (id, email, password_hash, first_name, last_name)
- `addresses` (id, user_id, street, city, postal_code)
- `sessions` (id, user_id, token_hash, expires_at)

**Dependencies:**
- Notification Service (via Kafka for welcome email)

---

### 2. Product Service

**Purpose:** Product catalog, inventory management

**Technology:**
- Go + Chi (HTTP)
- PostgreSQL (products, inventory)
- Redis (cache frequently accessed products)

**API:**
```protobuf
service ProductService {
  rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc CheckStock(CheckStockRequest) returns (CheckStockResponse);
  rpc ReserveStock(ReserveStockRequest) returns (ReserveStockResponse);
  rpc ReleaseStock(ReleaseStockRequest) returns (ReleaseStockResponse);
}
```

**Database:**
- `products` (id, brand, model, colorway, retail_price, release_date)
- `product_images` (id, product_id, url, position)
- `sizes` (id, product_id, size, stock_quantity)

**Key Features:**
- Stock reservation (pessimistic locking)
- Image management
- Size-based inventory

---

### 3. Bidding Service ⚡ (Core!)

**Purpose:** Bid/Ask management + Matching Engine

**Technology:**
- Go + gRPC
- PostgreSQL (bids, asks, matches)
- Redis (order book for fast matching)
- Goroutines + Channels (matching algorithm)

**API:**
```protobuf
service BiddingService {
  // Buyer actions
  rpc PlaceBid(PlaceBidRequest) returns (PlaceBidResponse);
  rpc CancelBid(CancelBidRequest) returns (CancelBidResponse);
  rpc GetMyBids(GetMyBidsRequest) returns (GetMyBidsResponse);
  
  // Seller actions
  rpc PlaceAsk(PlaceAskRequest) returns (PlaceAskResponse);
  rpc CancelAsk(CancelAskRequest) returns (CancelAskResponse);
  rpc GetMyAsks(GetMyAsksRequest) returns (GetMyAsksResponse);
  
  // Market data
  rpc GetMarketData(GetMarketDataRequest) returns (GetMarketDataResponse);
  rpc GetOrderBook(GetOrderBookRequest) returns (GetOrderBookResponse);
}
```

**Database:**
```sql
-- Bids (buyer intentions)
CREATE TABLE bids (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    product_id BIGINT REFERENCES products(id),
    size VARCHAR(10) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_active_bids (product_id, size, status, amount DESC)
);

-- Asks (seller listings)
CREATE TABLE asks (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT REFERENCES users(id),
    product_id BIGINT REFERENCES products(id),
    size VARCHAR(10) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    condition VARCHAR(20) DEFAULT 'new',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_active_asks (product_id, size, status, amount ASC)
);

-- Matches (successful bid/ask matches)
CREATE TABLE matches (
    id BIGSERIAL PRIMARY KEY,
    bid_id BIGINT REFERENCES bids(id),
    ask_id BIGINT REFERENCES asks(id),
    product_id BIGINT,
    size VARCHAR(10),
    price DECIMAL(10, 2),
    buyer_id BIGINT,
    seller_id BIGINT,
    status VARCHAR(30) DEFAULT 'pending',
    matched_at TIMESTAMP DEFAULT NOW()
);
```

**Matching Engine Algorithm:**

```go
// Main matching loop (goroutine)
func (e *MatchingEngine) Run(ctx context.Context) {
    ticker := time.NewTicker(100 * time.Millisecond)
    
    for {
        select {
        case <-ctx.Done():
            return
        case newBid := <-e.bidChan:
            go e.tryMatchBid(newBid)
        case newAsk := <-e.askChan:
            go e.tryMatchAsk(newAsk)
        case <-ticker.C:
            // Periodic scan for expired bids
            e.expireBids()
        }
    }
}

func (e *MatchingEngine) tryMatchBid(bid Bid) {
    // Get lowest ask for same product + size
    ask := e.redis.ZRange(
        fmt.Sprintf("asks:%d:%s", bid.ProductID, bid.Size), 
        0, 0, // Get lowest price
    )
    
    if ask == nil {
        return // No asks available
    }
    
    if bid.Amount >= ask.Amount {
        // MATCH! Atomic transaction
        e.createMatch(bid, ask)
        
        // Publish event
        e.kafka.Publish("matches.created", MatchEvent{
            BidID: bid.ID,
            AskID: ask.ID,
            Price: ask.Amount,
        })
    }
}
```

**Redis Order Book Structure:**
```redis
# Sorted set by price (ascending for asks, descending for bids)
ZADD asks:123:9 450.00 "ask:456"    # Product 123, Size 9, $450
ZADD asks:123:9 460.00 "ask:789"
ZADD asks:123:9 470.00 "ask:101"

ZADD bids:123:9 440.00 "bid:321"    # Highest bid $440
ZADD bids:123:9 430.00 "bid:654"
ZADD bids:123:9 420.00 "bid:987"

# Get best prices
ZRANGE asks:123:9 0 0  # Lowest ask: $450
ZREVRANGE bids:123:9 0 0  # Highest bid: $440
```

---

### 4. Order Service

**Purpose:** Order orchestration (Saga pattern)

**Technology:**
- Go + gRPC
- PostgreSQL (orders, order_items, order_events)
- Saga coordinator

**API:**
```protobuf
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
  rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
  rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (UpdateOrderStatusResponse);
  rpc CancelOrder(CancelOrderRequest) returns (CancelOrderResponse);
}
```

**Saga Flow:**
```
CreateOrder():
1. Reserve stock (Product Service) ← gRPC
   ↓ success
2. Process payment (Payment Service) ← gRPC
   ↓ success
3. Create order record
   ↓
4. Publish "order.created" event ← Kafka
   ↓
5. Return order ID

If any step fails → Compensate:
- Release stock
- Refund payment
- Cancel order
```

**Database:**
```sql
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGINT REFERENCES matches(id),
    buyer_id BIGINT,
    seller_id BIGINT,
    product_id BIGINT,
    size VARCHAR(10),
    price DECIMAL(10, 2),
    status VARCHAR(30), -- pending, confirmed, seller_shipping, 
                       -- authenticating, auth_passed, shipping_to_buyer, 
                       -- delivered, cancelled, auth_failed
    created_at TIMESTAMP DEFAULT NOW()
);

-- Event Sourcing for audit trail
CREATE TABLE order_events (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES orders(id),
    event_type VARCHAR(50), -- created, paid, shipped, authenticated, etc
    event_data JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

### 5. Payment Service

**Purpose:** Payment processing via Stripe

**Technology:**
- Go + Stripe SDK
- PostgreSQL (payments, payouts)

**API:**
```protobuf
service PaymentService {
  rpc ProcessPayment(ProcessPaymentRequest) returns (ProcessPaymentResponse);
  rpc CreatePayout(CreatePayoutRequest) returns (CreatePayoutResponse);
  rpc RefundPayment(RefundRequest) returns (RefundResponse);
  rpc HandleWebhook(WebhookRequest) returns (WebhookResponse);
}
```

**Flow:**
```
1. Create Payment Intent (Stripe)
2. Confirm payment
3. Hold funds (escrow) until order complete
4. On delivery confirmed:
   - Release payment to seller (minus platform fee)
5. If order cancelled:
   - Refund buyer
```

---

### 6. Notification Service

**Purpose:** Email, SMS, Push notifications

**Technology:**
- Go + Kafka Consumer
- SendGrid (email)
- Twilio (SMS)
- Worker pool pattern

**Events Consumed:**
- `user.registered` → Welcome email
- `match.created` → Notify buyer & seller
- `order.shipped` → Shipping notification
- `order.delivered` → Delivery confirmation
- `price.dropped` → Price alert

**Worker Pool:**
```go
type NotificationService struct {
    workers       int
    notifications chan Notification
    workerPool    chan struct{} // semaphore
}

func (s *NotificationService) Start(ctx context.Context) {
    // Start workers
    for i := 0; i < s.workers; i++ {
        go s.worker(ctx, i)
    }
    
    // Consume Kafka
    go s.consumeEvents(ctx)
}

func (s *NotificationService) worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            return
        case notif := <-s.notifications:
            s.workerPool <- struct{}{} // acquire
            s.send(notif)
            <-s.workerPool // release
        }
    }
}
```

---

### 7. Search Service

**Purpose:** Full-text search, filters

**Technology:**
- Go + Elasticsearch
- Kafka Consumer (for indexing)

**API:**
```protobuf
service SearchService {
  rpc Search(SearchRequest) returns (SearchResponse);
  rpc Autocomplete(AutocompleteRequest) returns (AutocompleteResponse);
  rpc GetFacets(GetFacetsRequest) returns (GetFacetsResponse);
}
```

**Elasticsearch Index:**
```json
{
  "mappings": {
    "properties": {
      "brand": {"type": "keyword"},
      "model": {"type": "text"},
      "colorway": {"type": "text"},
      "retail_price": {"type": "float"},
      "current_price": {"type": "float"},
      "sizes": {"type": "keyword"},
      "release_date": {"type": "date"},
      "suggest": {
        "type": "completion",
        "contexts": [{"name": "brand", "type": "category"}]
      }
    }
  }
}
```

---

### 8. Analytics Service

**Purpose:** Market data, price charts

**Technology:**
- Go + InfluxDB (time-series)
- Kafka Consumer
- Pipeline pattern

**Data Flow:**
```
Kafka (matches.created) 
  → Process event 
  → Aggregate metrics 
  → Store in InfluxDB

Metrics:
- sales_total (count)
- revenue (sum)
- average_price (mean)
- price_by_size (per size)
```

**InfluxDB Schema:**
```
measurement: sales
tags: product_id, size
fields: price, timestamp
```

---

### 9. Authentication Service

**Purpose:** Product verification workflow

**Technology:**
- Go + PostgreSQL
- Admin panel for authenticators

**API:**
```protobuf
service AuthenticationService {
  rpc CreateAuthTask(CreateAuthTaskRequest) returns (CreateAuthTaskResponse);
  rpc GetAuthTasks(GetAuthTasksRequest) returns (GetAuthTasksResponse);
  rpc CompleteAuth(CompleteAuthRequest) returns (CompleteAuthResponse);
}
```

**Workflow:**
```
1. Order created → Auth task created
2. Seller ships to warehouse
3. Admin opens task in panel
4. Checks:
   - Box condition
   - Authenticity markers
   - Size verification
5. Decision:
   a) PASS ✅ → Update order → Ship to buyer
   b) FAIL ❌ → Return to seller → Refund buyer
```

---

## 🔐 Security

### Authentication
- JWT tokens (HS256)
- Refresh tokens
- Token blacklist (Redis)

### Authorization
- Role-based (buyer, seller, admin, authenticator)
- Middleware per service

### Data Protection
- Passwords: bcrypt (cost 12)
- Sensitive data: encrypted at rest
- TLS for all communication

### Rate Limiting
- Per-user limits
- Per-IP limits
- Token bucket algorithm

---

## 📊 Data Flow Examples

### Scenario 1: Place Bid

```
1. User clicks "Place Bid $450"
   ↓
2. API Gateway → Bidding Service (gRPC)
   ↓
3. Bidding Service:
   - Validate bid
   - Save to PostgreSQL
   - Add to Redis order book
   - Try match with existing asks
   ↓
4. If match found:
   - Publish "match.created" event
   - Order Service handles order creation
   ↓
5. WebSocket broadcasts price update
```

### Scenario 2: Complete Order

```
1. Match created
   ↓
2. Order Service (Saga):
   a) Reserve stock (Product Service)
   b) Process payment (Payment Service)
   c) Create order
   d) Publish "order.created"
   ↓
3. Notification Service (async):
   - Email to buyer
   - Email to seller (shipping label)
   ↓
4. Seller ships to warehouse
   ↓
5. Authentication Service:
   - Verify product
   - Update order status
   ↓
6. Ship to buyer
   ↓
7. Delivered → Release payment to seller
```

---

## 🚀 Scalability Considerations

### Horizontal Scaling
- All services stateless (except Redis/DBs)
- Load balancing via Kubernetes

### Database Scaling
- Read replicas for heavy reads
- Sharding by user_id if needed
- Connection pooling

### Cache Strategy
- Redis for:
  - Order book (hot data)
  - Session tokens
  - Frequently accessed products
- Cache invalidation on writes

### Message Queue
- Kafka partitions for parallelism
- Consumer groups for load distribution

---

## 🔍 Monitoring

### Metrics (Prometheus)
Per service:
- Request count
- Request duration
- Error rate
- Active goroutines
- DB connection pool

Custom metrics:
- `bids_total`
- `matches_total`
- `order_processing_duration`

### Logging (ELK)
Structured logs (JSON):
- Timestamp
- Service
- Level (debug, info, warn, error)
- Trace ID
- Message

### Tracing (Jaeger)
- End-to-end request tracing
- Service dependencies
- Latency breakdown

---

## 🎯 Trade-offs & Decisions

### Why Microservices?
**Pros:**
- Independent scaling
- Tech diversity
- Team autonomy
- Fault isolation

**Cons:**
- Complexity
- Network latency
- Distributed transactions

**Decision:** Worth it for learning + portfolio

### Why gRPC?
**Pros:**
- Fast (binary protocol)
- Type-safe
- Bi-directional streaming

**Cons:**
- Not browser-friendly
- Steep learning curve

**Decision:** gRPC internal, REST for public API

### Why Kafka?
**Pros:**
- High throughput
- Durable
- Replay events

**Cons:**
- Operational complexity
- Eventual consistency

**Decision:** Perfect for event-driven architecture

---

**Next:** [DATABASE_SCHEMA.md](./DATABASE_SCHEMA.md)
