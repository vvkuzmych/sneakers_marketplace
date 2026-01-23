# 👟 Sneakers Marketplace

> **Production-ready microservices platform for sneakers trading with real-time Bid/Ask matching engine**

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go)](https://golang.org)
[![gRPC](https://img.shields.io/badge/gRPC-Protocol%20Buffers-4285F4?style=flat&logo=google)](https://grpc.io)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)](https://postgresql.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**A modern, scalable sneakers marketplace built with Go microservices architecture. Features a sophisticated Bid/Ask matching engine inspired by StockX and GOAT.**

---

## 🚀 Features

### 🔐 User Service (Port 50051)
- **Secure Authentication** - JWT-based auth with access & refresh tokens
- **User Management** - Profile, addresses, and session tracking
- **Password Security** - bcrypt hashing with cost 12
- **Multi-address Support** - Shipping and billing addresses

### 📦 Product Service (Port 50052)
- **Product Catalog** - Complete CRUD with SKU management
- **Smart Inventory** - Size-based inventory with reservation system
- **Image Management** - Multiple images per product with primary flag
- **Full-text Search** - Fast product search using PostgreSQL GIN index
- **Audit Trail** - Complete inventory transaction history

### 🎯 Bidding Service (Port 50053) ⭐
- **Bid/Ask System** - Place buy and sell orders
- **Automatic Matching Engine** - Instant order matching when prices cross
- **Real-time Market Data** - Highest bid, lowest ask, spread calculation
- **Order Book** - View active bids and asks
- **Match History** - Complete transaction history
- **FIFO Algorithm** - Fair, time-priority matching
- **Auto-Notifications** - Instant alerts on match creation

### 📦 Order Service (Port 50054)
- **Order Management** - Complete order lifecycle (11 states)
- **Fee Calculation** - Buyer (5%) + Seller (4%) fees
- **Order Numbers** - Auto-generated (ORD-YYYYMMDD-XXXX)
- **Shipping Tracking** - Tracking number integration
- **Status History** - Complete audit trail

### 💳 Payment Service (Port 50055)
- **Stripe Integration** - Real + Demo mode
- **Payment Intents** - Secure payment processing
- **Refunds** - Full and partial refunds
- **Seller Payouts** - Stripe Connect integration
- **Payment History** - Complete transaction tracking

### 🔔 Notification Service (Port 50056) 🆕
- **Email Notifications** - SMTP integration (Mailhog)
- **Real-time Alerts** - WebSocket push notifications
- **8 Notification Types** - Match, Order, Payment, Refund, etc.
- **User Preferences** - Control email/push settings
- **Read Tracking** - Mark as read/unread
- **History** - Paginated notification list

### 🌐 API Gateway (Port 8080) 🆕
- **HTTP REST API** - User-friendly HTTP endpoints
- **WebSocket Support** - Real-time bidirectional communication
- **JWT Authentication** - Secure token-based auth
- **gRPC Proxy** - Routes to all microservices
- **CORS Support** - Cross-origin requests
- **Health Checks** - Service status monitoring

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Client Applications                   │
│              (REST API Gateway - Future)                 │
└────────────────────┬────────────────────────────────────┘
                     │ gRPC
          ┌──────────┼──────────┐
          │          │          │
┌─────────▼────┐ ┌──▼──────┐ ┌─▼────────────┐
│User Service  │ │ Product │ │   Bidding    │
│   :50051     │ │ Service │ │   Service    │
│              │ │ :50052  │ │   :50053     │
│• Auth & JWT  │ │• Catalog│ │• Bid/Ask     │
│• Profile     │ │• Invntry│ │• Matching 🔥 │
│• Addresses   │ │• Search │ │• Order Book  │
└──────┬───────┘ └────┬────┘ └──────┬───────┘
       │              │              │
       └──────────────┼──────────────┘
                      │
            ┌─────────▼─────────┐
            │   PostgreSQL 16   │
            │  • 11 tables      │
            │  • Migrations     │
            │  • Full-text idx  │
            └───────────────────┘
```

**Microservices Pattern:**
- Independent deployment
- Technology freedom
- Horizontal scaling
- Fault isolation

---

## 📋 Prerequisites

- **Go** 1.25+
- **PostgreSQL** 16+
- **Docker** & Docker Compose
- **golang-migrate** (for migrations)
- **grpcurl** (for testing)

### Installation

```bash
# macOS
brew install go postgresql docker golang-migrate grpcurl

# Verify versions
go version          # 1.25+
psql --version      # 16+
migrate -version    # 4.x+
```

---

## 🚀 Quick Start

> **New!** We now have a comprehensive Makefile with 60+ commands! See [MAKEFILE_GUIDE.md](MAKEFILE_GUIDE.md) for full documentation.

### 🎯 Super Quick Start (3 commands!)

```bash
# 1. Install all dependencies (Go, npm, tools)
make install

# 2. Setup database (create, migrate, seed with test data)
make db-setup

# 3. Start everything (backend + frontend + mailhog)
make dev
```

**That's it!** Your application is now running:
- 🌐 Frontend: http://localhost:5173
- 🔌 API Gateway: http://localhost:8080
- 📧 Mailhog: http://localhost:8025

**Test users (password: `password123`):**
- john@example.com
- jane@example.com
- test@example.com

### 📋 Common Makefile Commands

```bash
make help              # Show all available commands
make dev               # Start all services (recommended!)
make stop              # Stop all services
make restart           # Rebuild & restart backend
make status            # Check running services
make logs              # Show recent logs
make logs-follow       # Follow logs in real-time

make build             # Build all backend services
make test              # Run all tests with coverage
make lint              # Run all linters
make lint-fix          # Auto-fix linting issues

make db-reset          # Drop + recreate database
make db-backup         # Backup database
make proto             # Generate gRPC code from .proto files

make version           # Show Go, Node, npm versions
make health            # API health check
```

**Full documentation:** [MAKEFILE_GUIDE.md](MAKEFILE_GUIDE.md)

---

### 🔧 Manual Setup (Alternative)

<details>
<summary>Click to expand manual setup instructions</summary>

### 1. Clone & Setup

```bash
git clone https://github.com/vvkuzmych/sneakers_marketplace
cd sneakers_marketplace

# Create .env file
cp env.example .env

# Generate JWT secret
echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env
```

### 2. Start Infrastructure

```bash
# Start PostgreSQL, Redis, Kafka, etc.
docker-compose up -d

# Verify containers
docker-compose ps

# Check PostgreSQL
psql postgres://postgres:postgres@localhost:5432/sneakers_marketplace -c "\dt"
```

### 3. Run Migrations

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/sneakers_marketplace?sslmode=disable"
migrate -path migrations -database "${DATABASE_URL}" up

# Verify tables
psql ${DATABASE_URL} -c "\dt"
```

### 4. Build Services

```bash
go build -o bin/user-service ./cmd/user-service
go build -o bin/product-service ./cmd/product-service
go build -o bin/bidding-service ./cmd/bidding-service
go build -o bin/api-gateway ./cmd/api-gateway
go build -o bin/notification-service ./cmd/notification-service
```

### 5. Run Services

```bash
# Start all services in background
export $(cat .env | grep -v '^#' | xargs)
nohup ./bin/user-service > /tmp/user-service.log 2>&1 &
nohup ./bin/product-service > /tmp/product-service.log 2>&1 &
nohup ./bin/bidding-service > /tmp/bidding-service.log 2>&1 &
nohup ./bin/notification-service > /tmp/notification-service.log 2>&1 &
nohup ./bin/api-gateway > /tmp/api-gateway.log 2>&1 &

# Frontend
cd frontend && npm run dev
```

### 6. Test Services

```bash
# Test API Gateway
curl http://localhost:8080/health

# Test login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

</details>

---

## 🧪 Testing

### User Service Test

```bash
# Register new user
grpcurl -plaintext -d '{
  "email": "alice@example.com",
  "password": "SecurePass123!",
  "first_name": "Alice",
  "last_name": "Smith"
}' localhost:50051 user.UserService/Register

# Login
grpcurl -plaintext -d '{
  "email": "alice@example.com",
  "password": "SecurePass123!"
}' localhost:50051 user.UserService/Login
```

### Product Service Test

```bash
# Create product
grpcurl -plaintext -d '{
  "sku": "AJ1-001",
  "name": "Air Jordan 1 Chicago",
  "brand": "Nike",
  "price": 170.00
}' localhost:50052 product.ProductService/CreateProduct

# List products
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50052 product.ProductService/ListProducts
```

### Bidding Service Test

```bash
# Place bid
grpcurl -plaintext -d '{
  "user_id": 1,
  "product_id": 1,
  "size_id": 1,
  "price": 200.00,
  "quantity": 1
}' localhost:50053 bidding.BiddingService/PlaceBid

# Get market price
grpcurl -plaintext -d '{
  "product_id": 1,
  "size_id": 1
}' localhost:50053 bidding.BiddingService/GetMarketPrice
```

---

## 📚 API Documentation

### User Service (50051)

| Method | Description |
|--------|-------------|
| `Register` | Create new user account |
| `Login` | Authenticate user, get JWT tokens |
| `RefreshToken` | Get new access token |
| `Logout` | Invalidate session |
| `GetProfile` | Get user profile |
| `UpdateProfile` | Update user information |
| `AddAddress` | Add shipping/billing address |
| `GetAddresses` | List user addresses |
| `UpdateAddress` | Update address |
| `DeleteAddress` | Remove address |

### Product Service (50052)

| Method | Description |
|--------|-------------|
| `CreateProduct` | Add new product |
| `GetProduct` | Get product with images & sizes |
| `ListProducts` | List with pagination & filters |
| `UpdateProduct` | Update product details |
| `DeleteProduct` | Remove product |
| `SearchProducts` | Full-text search |
| `AddProductImage` | Add product image |
| `DeleteProductImage` | Remove image |
| `AddSize` | Add size with inventory |
| `GetAvailableSizes` | List sizes & stock |
| `UpdateInventory` | Adjust stock quantity |
| `ReserveInventory` | Reserve for order |
| `ReleaseInventory` | Release reservation |

### Bidding Service (50053)

| Method | Description |
|--------|-------------|
| `PlaceBid` | Place buy order (auto-match) |
| `PlaceAsk` | Place sell order (auto-match) |
| `GetBid` | Get bid details |
| `GetAsk` | Get ask details |
| `GetUserBids` | List user's bids |
| `GetUserAsks` | List user's asks |
| `GetProductBids` | View order book (bids) |
| `GetProductAsks` | View order book (asks) |
| `CancelBid` | Cancel buy order |
| `CancelAsk` | Cancel sell order |
| `GetHighestBid` | Get top bid |
| `GetLowestAsk` | Get lowest ask |
| `GetMarketPrice` | Get market data & spread |
| `GetMatch` | Get match details |
| `GetUserMatches` | List user's transactions |

---

## 🔧 Configuration

### Environment Variables

```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5435/sneakers_marketplace?sslmode=disable
DATABASE_PORT=5435

# Redis
REDIS_URL=redis://localhost:6380/0
REDIS_PORT=6380

# Kafka
KAFKA_BROKERS=localhost:9094

# JWT
JWT_SECRET=<generated-secret>
JWT_EXPIRATION=24h
REFRESH_TOKEN_EXPIRATION=168h

# Server Ports
SERVER_PORT=50051  # Override per service
USER_SERVICE_ADDR=localhost:50051
PRODUCT_SERVICE_ADDR=localhost:50052
BIDDING_SERVICE_ADDR=localhost:50053
```

### Database Connection Pooling

```go
MaxConns:           25
MaxConnLifetime:    1 hour
MaxConnIdleTime:    30 minutes
```

---

## 🎯 Matching Engine Algorithm

### How It Works

The bidding service implements a **real-time matching engine** that automatically matches bids and asks when prices cross:

```
Bid (Buyer):  "I want to BUY at $X or HIGHER"
Ask (Seller): "I want to SELL at $Y or LOWER"

Match when: Bid Price >= Ask Price
Final Price: Ask Price (seller's price)
```

### Example Flow

```
1. Buyer places BID: $200  → active (waiting for seller)
2. Seller places ASK: $220 → active (no match, price too high)
3. Market shows spread: $200 / $220

4. New buyer places BID: $225 → INSTANT MATCH! ⚡
   - Bid $225 >= Ask $220 ✓
   - Match created at $220 (seller's price)
   - Both orders marked as 'matched'
   - Match record created
   - Orders removed from order book
```

### Matching Rules

- **Price Priority:** Best price gets matched first
- **Time Priority:** Earlier orders matched first (FIFO)
- **Exact Quantity:** Only exact quantity matches (for simplicity)
- **Atomic Transactions:** All updates in single DB transaction
- **Instant Execution:** Matching happens immediately on order placement

---

## 📊 Database Schema

### Users & Auth
- `users` - user accounts
- `addresses` - shipping/billing addresses
- `sessions` - JWT session management

### Products & Inventory
- `products` - sneaker catalog
- `product_images` - multiple images per product
- `sizes` - inventory by size
- `inventory_transactions` - audit trail

### Bidding & Trading
- `bids` - buyer orders
- `asks` - seller orders
- `matches` - completed transactions

---

## 🛠️ Development

### Project Structure

```
sneakers_marketplace/
├── cmd/                    # Service entry points
│   ├── user-service/
│   ├── product-service/
│   └── bidding-service/
├── internal/               # Private application code
│   ├── user/              # User service
│   ├── product/           # Product service
│   └── bidding/           # Bidding service
├── pkg/                    # Shared packages
│   ├── auth/              # JWT & Password
│   ├── config/            # Configuration
│   ├── database/          # DB connections
│   ├── logger/            # Logging
│   └── proto/             # gRPC definitions
├── migrations/             # SQL migrations
├── scripts/                # Helper scripts
└── docs/                   # Documentation
```

### Adding a New Service

1. Create proto definition in `pkg/proto/service/`
2. Generate Go code: `protoc --go_out=. --go-grpc_out=. pkg/proto/service/service.proto`
3. Implement model, repository, service, handler in `internal/service/`
4. Create main.go in `cmd/service/`
5. Add migration if needed
6. Write tests

---

## 🚧 Roadmap

### Phase 1 ✅ (Completed)
- [x] User Service (Auth, Profile)
- [x] Product Service (Catalog, Inventory)
- [x] Bidding Service (Matching Engine)
- [x] Database migrations
- [x] Test scripts

### Phase 2 🔄 (Next)
- [ ] Order Service (process matches)
- [ ] Payment Service (Stripe integration)
- [ ] Notification Service (email, websockets)
- [ ] API Gateway (HTTP REST + Swagger)
- [ ] WebSocket for real-time updates

### Phase 3 📅 (Future)
- [ ] Frontend (React/Next.js)
- [ ] Kubernetes deployment
- [ ] CI/CD pipeline
- [ ] Monitoring & Alerts
- [ ] Load testing

---

## 📈 Performance

### Benchmarks (Local)

- **User Registration:** ~50ms
- **Product Search:** ~15ms (with full-text index)
- **Bid/Ask Matching:** ~10ms (including transaction)
- **Market Price Calculation:** ~5ms

### Scalability

- **Horizontal Scaling:** Each service can run multiple instances
- **Database Pooling:** 25 connections per service instance
- **Stateless Services:** No shared state, perfect for k8s
- **Matching Engine:** O(1) lookup with indexed price columns

---

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👥 Authors

- **Volodymyr Kuzmych** - [@vvkuzmych](https://github.com/vvkuzmych)

---

## 🙏 Acknowledgments

- Inspired by **StockX** and **GOAT** marketplace platforms
- Built with **Go** and **gRPC**
- Uses **PostgreSQL** for reliability
- Matching engine concept from traditional stock exchanges

---

## 📞 Support

- 📧 Email: [support@example.com](mailto:support@example.com)
- 🐛 Issues: [GitHub Issues](https://github.com/vvkuzmych/sneakers_marketplace/issues)
- 📚 Documentation: [Full Docs](https://github.com/vvkuzmych/sneakers_marketplace/tree/main/docs)

---

**⭐ Star this repo if you find it useful!**

**Made with ❤️ and Go**
