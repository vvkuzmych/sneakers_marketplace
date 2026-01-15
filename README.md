# 👟 Sneakers Marketplace

**Production-ready microservices platform for sneaker trading with real-time auction system**

> Inspired by StockX and GOAT - A marketplace where sneaker enthusiasts can buy and sell limited edition sneakers through a sophisticated bid/ask system with authentication verification.

---

## 🎯 Project Overview

Sneakers Marketplace is a full-stack e-commerce platform built with Go microservices that demonstrates:

- **Real-time auction system** (Bid/Ask matching like a stock exchange)
- **Microservices architecture** (9 independent services)
- **Event-driven design** (Kafka for async communication)
- **Authentication workflow** (Multi-step order verification)
- **Production-ready patterns** (CQRS, Saga, Circuit Breaker)

### 🔥 Key Features

#### For Buyers
- 🔍 Browse sneakers with real-time market data
- 💰 Place bids or buy instantly
- 📊 View price history charts
- 🔔 Get notified when price drops
- 📦 Track orders in real-time
- 💼 Track portfolio value

#### For Sellers
- 📝 List sneakers with ask price
- 💵 Sell instantly to highest bidder
- ✅ Authentication verification
- 💳 Secure payouts
- 📈 Sales analytics

#### Platform
- 🤖 Automatic bid/ask matching
- ✓ Product authentication flow
- 📱 Real-time updates (WebSockets)
- 🔐 Secure payments (Stripe)
- 📊 Market analytics

---

## 🏗️ Architecture

### Microservices (9 services)

```
┌──────────────────┐
│   API Gateway    │ ← External traffic (REST)
└────────┬─────────┘
         │
    ┌────┴────────────────────┐
    │                         │
┌───▼────────┐        ┌──────▼──────┐
│   User     │◄─gRPC─►│   Product   │
│  Service   │        │   Service   │
└────────────┘        └──────┬──────┘
                             │
                      ┌──────▼──────┐
                      │   Bidding   │ ← Core business logic
                      │   Service   │ ← Matching Engine 🔥
                      └──────┬──────┘
                             │
            ┌────────────────┼────────────────┐
            │                │                │
    ┌───────▼────┐   ┌──────▼──────┐  ┌─────▼──────┐
    │   Order    │   │   Payment   │  │  Matching  │
    │  Service   │   │   Service   │  │   Engine   │
    └────────────┘   └─────────────┘  └────────────┘
            │                │                │
    ┌───────▼────────────────▼────────────────▼──────┐
    │         Event Bus (Kafka / RabbitMQ)            │
    └──┬──────────────┬──────────────┬────────────┬──┘
       │              │              │            │
┌──────▼────┐  ┌─────▼─────┐  ┌────▼────┐  ┌───▼─────┐
│Notification│  │ Analytics │  │ Search  │  │  Auth   │
│  Service   │  │  Service  │  │Service  │  │ Service │
└────────────┘  └───────────┘  └─────────┘  └─────────┘
```

### Service Responsibilities

| Service | Technology | Purpose |
|---------|-----------|---------|
| **User Service** | Go + PostgreSQL | Authentication, profiles, addresses, wishlist |
| **Product Service** | Go + PostgreSQL | Catalog, inventory, variants (sizes) |
| **Bidding Service** | Go + PostgreSQL + Redis | Bid/Ask management, matching engine |
| **Order Service** | Go + PostgreSQL | Order orchestration (Saga pattern) |
| **Payment Service** | Go + PostgreSQL + Stripe | Payment processing, refunds |
| **Notification Service** | Go + Kafka | Emails, SMS, push notifications |
| **Search Service** | Go + Elasticsearch | Full-text search, filters |
| **Analytics Service** | Go + InfluxDB | Market data, price charts, reports |
| **Authentication Service** | Go + PostgreSQL | Product verification workflow |

---

## 💰 How It Works

### The Bid/Ask System

```
Example: Nike Air Jordan 1 "Chicago" - Size US 9

Current Market State:
├─ Last Sale: $420
├─ Highest Bid: $410 ← Buyer wants to buy
├─ Lowest Ask: $450 ← Seller wants to sell
└─ Gap: $40

Scenario 1: Instant Buy
- Buyer clicks "Buy Now" at $450
- Matches with Lowest Ask
- Order created immediately

Scenario 2: Place Bid
- Buyer places Bid at $430
- Added to order book
- When seller lists Ask ≤ $430 → AUTO MATCH! 🎉

Scenario 3: Market Movement
- New Bid: $440
- New Bid: $445
- New Ask: $445
- MATCH! Order created at $445
```

### Order Flow (Authentication)

```
1. Match Created (Bid meets Ask)
   ↓
2. Payment Processed
   ↓
3. Seller ships to Authentication Center
   Status: "En route to verification"
   ↓
4. Authentication Team inspects:
   - Box condition
   - Authenticity (stitching, materials, tags)
   - Size verification
   ↓
5a. PASS ✅
    - Ship to Buyer
    - Release payment to Seller (minus fees)
    - Order Complete
    
5b. FAIL ❌
    - Return to Seller
    - Refund Buyer
    - Seller gets warning
```

---

## 🛠️ Tech Stack

### Backend
- **Language:** Go 1.25+
- **Communication:** gRPC (inter-service), REST (client API)
- **Databases:** 
  - PostgreSQL (primary data)
  - Redis (cache, order book)
  - Elasticsearch (search)
  - InfluxDB (time-series metrics)
- **Message Queue:** Kafka / RabbitMQ
- **Payment:** Stripe API
- **Email/SMS:** SendGrid, Twilio

### Infrastructure
- **Containers:** Docker
- **Orchestration:** Kubernetes
- **Service Discovery:** Consul
- **API Gateway:** Kong / Custom
- **Monitoring:** Prometheus + Grafana
- **Tracing:** Jaeger
- **Logging:** ELK Stack
- **CI/CD:** GitHub Actions

### Frontend (Optional - Future)
- React / Next.js
- WebSocket client (real-time updates)
- Chart.js (price charts)

---

## 📊 Database Design

### Key Tables

**Users & Auth**
- `users` - user accounts
- `addresses` - shipping/billing addresses
- `sessions` - JWT sessions

**Products**
- `products` - sneaker catalog (brand, model, colorway)
- `product_variants` - size-specific inventory
- `product_images` - product photos

**Trading** 🔥
- `bids` - buyer bids (with expiration)
- `asks` - seller listings
- `matches` - completed bid/ask matches
- `market_data` - price history (for charts)

**Orders**
- `orders` - order records
- `order_items` - line items
- `order_events` - event sourcing log

**Payments**
- `payments` - payment transactions
- `payouts` - seller payouts

See [DATABASE_SCHEMA.md](./docs/DATABASE_SCHEMA.md) for details.

---

## 🚀 Development Plan

### Phase 1: Foundation (Week 1-2)
- [x] Project setup & structure
- [ ] User Service (auth, profiles)
- [ ] Product Service (catalog)
- [ ] Basic CRUD operations
- [ ] Docker Compose for local dev

### Phase 2: Core Trading Logic (Week 3-4) 🔥
- [ ] Bidding Service (Bid/Ask management)
- [ ] **Matching Engine** (goroutines + channels)
- [ ] Order Service (Saga pattern)
- [ ] Payment Service (Stripe)
- [ ] Kafka setup

### Phase 3: Order Flow (Week 5)
- [ ] Authentication Service (verification workflow)
- [ ] Multi-step order states
- [ ] Notification Service
- [ ] Email templates

### Phase 4: Real-time & Analytics (Week 6)
- [ ] WebSocket server (real-time bidding)
- [ ] Search Service (Elasticsearch)
- [ ] Analytics Service (price charts)
- [ ] Market data aggregation

### Phase 5: Production Ready (Week 7-8)
- [ ] API Gateway (Kong)
- [ ] Service Discovery (Consul)
- [ ] Kubernetes deployment
- [ ] Monitoring (Prometheus + Grafana)
- [ ] Distributed tracing (Jaeger)
- [ ] Load testing
- [ ] Documentation

See [DEVELOPMENT_PLAN.md](./docs/DEVELOPMENT_PLAN.md) for detailed timeline.

---

## 📁 Project Structure

```
sneakers_marketplace/
├── cmd/
│   ├── user-service/
│   ├── product-service/
│   ├── bidding-service/      ← Matching Engine
│   ├── order-service/
│   ├── payment-service/
│   ├── notification-service/
│   ├── search-service/
│   ├── analytics-service/
│   └── auth-service/          ← Product authentication
├── internal/
│   ├── user/
│   ├── product/
│   ├── bidding/
│   ├── order/
│   ├── payment/
│   └── ...
├── pkg/
│   ├── proto/                 ← gRPC definitions
│   ├── kafka/                 ← Kafka client
│   ├── middleware/            ← Shared middleware
│   └── utils/
├── migrations/                ← SQL migrations
├── deployments/
│   ├── docker-compose.yml
│   ├── kubernetes/
│   └── terraform/
├── docs/                      ← Documentation
├── scripts/                   ← Helper scripts
└── tests/                     ← Integration tests
```

---

## 🎓 Go Concepts Demonstrated

### Week 1-2 (Basics)
- ✅ HTTP servers (Gin/Chi)
- ✅ PostgreSQL operations (pgx)
- ✅ Error handling & wrapping
- ✅ Testing (unit + integration)
- ✅ Structs, interfaces, methods

### Week 3-4 (Intermediate)
- ✅ Context (timeouts, cancellation)
- ✅ Custom error types
- ✅ Middleware (auth, logging)
- ✅ Environment config
- ✅ Database transactions

### Week 5 (Goroutines & Channels) 🔥
- ✅ **Worker pools** (notification service)
- ✅ **Pipeline pattern** (analytics)
- ✅ **Fan-out/fan-in** (parallel matching)
- ✅ **Channels** (bid/ask streaming)
- ✅ **Select** (event multiplexing)
- ✅ **Graceful shutdown** (all services)

### Advanced (Production)
- ✅ **Matching Engine** (custom algorithm)
- ✅ **Saga pattern** (distributed transactions)
- ✅ **Event Sourcing** (order events)
- ✅ **CQRS** (command/query separation)
- ✅ **gRPC** (inter-service communication)
- ✅ **Kafka** (event streaming)
- ✅ **WebSockets** (real-time updates)
- ✅ **Circuit Breaker** (resilience)
- ✅ **Distributed Tracing** (Jaeger)

---

## 🔧 Local Development

### Prerequisites
```bash
# Required
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

# Optional (for full stack)
- Kafka
- Elasticsearch
- Node.js (for frontend)
```

### Quick Start

```bash
# 1. Clone repository
git clone https://github.com/vvkuzmych/sneakers_marketplace.git
cd sneakers_marketplace

# 2. Start infrastructure
docker-compose up -d

# 3. Run migrations
make migrate-up

# 4. Start services (in separate terminals)
make run-user-service
make run-product-service
make run-bidding-service

# 5. (Optional) Seed database
make seed
```

### Environment Setup

```bash
# Copy example env file
cp .env.example .env

# Edit with your values
vim .env
```

Required env vars:
- `DATABASE_URL` - PostgreSQL connection
- `REDIS_URL` - Redis connection
- `KAFKA_BROKERS` - Kafka brokers
- `STRIPE_SECRET_KEY` - Stripe API key
- `JWT_SECRET` - JWT signing secret

---

## 📚 Documentation

- [Architecture](./docs/ARCHITECTURE.md) - Detailed system design
- [Database Schema](./docs/DATABASE_SCHEMA.md) - All tables explained
- [API Documentation](./docs/API.md) - REST & gRPC endpoints
- [Matching Engine](./docs/MATCHING_ENGINE.md) - How bid/ask matching works
- [Development Plan](./docs/DEVELOPMENT_PLAN.md) - Week-by-week roadmap
- [Deployment Guide](./docs/DEPLOYMENT.md) - Kubernetes setup

---

## 🧪 Testing

```bash
# Unit tests
make test

# Integration tests
make test-integration

# E2E tests
make test-e2e

# Load tests
make test-load

# Test coverage
make coverage
```

---

## 📈 Monitoring & Observability

### Metrics (Prometheus)
```
http://localhost:9090
```

Key metrics:
- `bids_total` - Total bids placed
- `matches_total` - Successful matches
- `order_duration_seconds` - Order processing time
- `payment_errors_total` - Payment failures

### Dashboards (Grafana)
```
http://localhost:3000
```

### Distributed Tracing (Jaeger)
```
http://localhost:16686
```

### Logs (Elasticsearch)
```
http://localhost:9200
```

---

## 🎯 Key Highlights for Portfolio

### 1. Matching Engine 🔥
**Problem:** How to efficiently match bids and asks in real-time?

**Solution:**
- In-memory order book (Redis)
- Goroutines for parallel matching
- Channels for bid/ask streaming
- Pessimistic locking for race conditions

### 2. Distributed Transactions (Saga Pattern)
**Problem:** Order involves multiple services (Product, Payment, Notification)

**Solution:**
- Saga coordinator in Order Service
- Compensation logic for rollbacks
- Event sourcing for audit trail

### 3. Real-time Updates (WebSockets)
**Problem:** Users need live price updates

**Solution:**
- WebSocket server with goroutine per connection
- Redis pub/sub for broadcasting
- Channels for message routing

### 4. Scalability
**Problem:** Handle 1000+ concurrent bids

**Solution:**
- Horizontal scaling (Kubernetes)
- Service discovery (Consul)
- Load balancing (Envoy)
- Circuit breakers (resilience)

---

## 🤝 Contributing

This is a learning project, but contributions are welcome!

1. Fork the repo
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

---

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details

---

## 🙏 Acknowledgments

- **StockX** & **GOAT** - Inspiration for the platform
- **Go community** - Amazing ecosystem and libraries
- **Open source projects** - Kafka, Kubernetes, Prometheus, and more

---

## 📞 Contact

**Project Maintainer:** Your Name
- GitHub: [@yourusername](https://github.com/yourusername)
- Email: your.email@example.com
- LinkedIn: [Your Profile](https://linkedin.com/in/yourprofile)

---

## 🎯 Project Status

**Current Phase:** Phase 1 - Foundation ✅
**Next Milestone:** User & Product Services (Week 1-2)
**Target Completion:** 8 weeks

---

**Built with ❤️ and Go**

**Star ⭐ this repo if you find it useful!**
