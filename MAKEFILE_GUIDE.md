# 📘 Makefile Guide - Quick Reference

## 🚀 Quick Start

```bash
# 1. Install all dependencies
make install

# 2. Setup database
make db-setup

# 3. Start all services (Backend + Frontend)
make dev
```

**That's it!** Your application is now running:
- 🌐 Frontend: http://localhost:5173
- 🔌 API Gateway: http://localhost:8080
- 📧 Mailhog: http://localhost:8025

---

## 📚 Table of Contents

- [Installation](#installation)
- [Development](#development)
- [Database](#database)
- [Testing](#testing)
- [Building](#building)
- [Linting](#linting)
- [Protobuf](#protobuf)
- [Docker](#docker)
- [Monitoring](#monitoring)
- [Utilities](#utilities)

---

## 🔧 Installation

### Install Everything

```bash
make install
```

This will:
- ✅ Install Go dependencies (`go mod download`)
- ✅ Install Frontend dependencies (`npm install`)
- ✅ Install dev tools (protoc, golangci-lint, mailhog)

### Install Separately

```bash
make install-backend     # Go dependencies only
make install-frontend    # npm dependencies only
make install-tools       # Dev tools only
```

---

## 💻 Development

### Start Everything

```bash
make dev
```

Starts:
- ✅ All backend services (User, Product, Bidding, Notification, API Gateway)
- ✅ Frontend (Vite dev server)
- ✅ Mailhog (email testing)

### Start Services Separately

```bash
make run-all          # All backend services
make run-frontend     # Frontend only
make run-mailhog      # Mailhog only
```

### Stop Everything

```bash
make stop
```

### Restart Services

```bash
make restart          # Rebuild and restart all backend services
```

### Check Service Status

```bash
make status           # Show running services
make health           # API health check
```

### View Logs

```bash
make logs             # Show last 20 lines from each service
make logs-follow      # Follow logs in real-time (Ctrl+C to stop)
```

**Log files location:** `/tmp/*-service.log`

---

## 🗄️ Database

### Setup Database (Fresh Install)

```bash
make db-setup
```

This will:
1. Create database (`sneakers_marketplace`)
2. Run migrations (create tables)
3. Seed with test data

### Individual Commands

```bash
make db-create        # Create database
make db-drop          # Drop database (DANGEROUS!)
make db-migrate       # Run migrations
make db-rollback      # Rollback migrations
make db-seed          # Seed test data
make db-reset         # Drop + Create + Migrate + Seed
make db-check         # Test database connection
```

### Backup & Restore

```bash
make db-backup        # Backup to backups/backup_YYYYMMDD_HHMMSS.sql
make db-restore       # Restore from latest backup
```

### Test Data

After `make db-seed`, you'll have:
- **5 users** (password: `password123`)
  - john@example.com
  - jane@example.com
  - bob@example.com
  - alice@example.com
  - test@example.com
- **8 products** (sneakers)
- **Sizes 7-13** for each product
- **Sample BIDs and ASKs**

---

## 🧪 Testing

### Run All Tests

```bash
make test
```

Generates:
- ✅ Test results
- ✅ Coverage report (`coverage.html`)

### Test Types

```bash
make test-unit          # Unit tests only
make test-integration   # Integration tests only
make test-coverage      # Generate coverage report
make test-frontend      # Frontend tests (Jest/Vitest)
```

---

## 🔨 Building

### Build Everything

```bash
make build              # Build all backend services
```

Creates binaries in `bin/`:
- `bin/api-gateway`
- `bin/user-service`
- `bin/product-service`
- `bin/bidding-service`
- `bin/order-service`
- `bin/payment-service`
- `bin/notification-service`

### Build Individual Services

```bash
make build-gateway
make build-user
make build-product
make build-bidding
make build-order
make build-payment
make build-notification
```

### Build Frontend

```bash
make build-frontend     # Production build → frontend/dist/
```

---

## 🔍 Linting

### Lint Everything

```bash
make lint               # Run all linters
```

### Lint Separately

```bash
make lint-backend       # golangci-lint (Go)
make lint-frontend      # ESLint (TypeScript/React)
```

### Auto-Fix Issues

```bash
make lint-fix           # Fix linting issues automatically
```

---

## 📦 Protobuf

### Generate gRPC Code

```bash
make proto
```

Regenerates Go code from `.proto` files in:
- `pkg/proto/user/*.proto`
- `pkg/proto/product/*.proto`
- `pkg/proto/bidding/*.proto`
- `pkg/proto/notification/*.proto`

**When to run:**
- ✅ After modifying `.proto` files
- ✅ After adding new services
- ✅ After changing message definitions

---

## 🐳 Docker

### Docker Commands

```bash
make docker-build       # Build Docker images
make docker-up          # Start containers
make docker-down        # Stop containers
make docker-logs        # Follow logs
```

---

## 📊 Monitoring

### Check Health

```bash
make health             # API Gateway health check
make status             # Show running processes
```

### Logs

```bash
make logs               # Show recent logs
make logs-follow        # Follow logs in real-time
```

---

## 🛠️ Utilities

### Environment Check

```bash
make env-check          # Verify .env file and required variables
```

### Version Info

```bash
make version            # Show Go, Node, npm, protoc versions
```

### Clean Up

```bash
make clean              # Remove build artifacts, logs, temp files
make clean-all          # Deep clean (remove node_modules, Go cache)
```

---

## 📋 Common Workflows

### 1️⃣ First Time Setup

```bash
make install            # Install dependencies
make db-setup           # Setup database
make dev                # Start everything
```

### 2️⃣ Daily Development

```bash
make dev                # Start services
# ... do your work ...
make stop               # Stop when done
```

### 3️⃣ After Code Changes

```bash
make stop               # Stop services
make build              # Rebuild
make run-all            # Restart backend
# Frontend auto-reloads via Vite HMR
```

### 4️⃣ Before Committing

```bash
make lint               # Check linting
make test               # Run tests
make lint-fix           # Fix auto-fixable issues
```

### 5️⃣ Database Changes

```bash
# After creating new migration
make db-migrate         # Apply migration

# If migration fails
make db-rollback        # Rollback
# ... fix migration ...
make db-migrate         # Try again
```

### 6️⃣ Debugging

```bash
make logs-follow        # Watch logs in real-time
make status             # Check if services are running
make health             # Test API connectivity
```

---

## 🎯 Service Ports

| Service | Port | URL |
|---------|------|-----|
| Frontend | 5173 | http://localhost:5173 |
| API Gateway | 8080 | http://localhost:8080 |
| User Service | 50051 | gRPC only |
| Product Service | 50052 | gRPC only |
| Bidding Service | 50053 | gRPC only |
| Order Service | 50054 | gRPC only |
| Payment Service | 50055 | gRPC only |
| Notification Service | 50056 | gRPC only |
| Mailhog UI | 8025 | http://localhost:8025 |
| Mailhog SMTP | 1025 | Internal |
| PostgreSQL | 5432 | Internal |

---

## 🆘 Troubleshooting

### Services Won't Start

```bash
make stop               # Stop all
make clean              # Clean artifacts
ps aux | grep -E 'api-gateway|user-service' | grep -v grep
# Kill any lingering processes
make build              # Rebuild
make run-all            # Start again
```

### Database Issues

```bash
make db-check           # Test connection
make db-reset           # Nuclear option (drops and recreates)
```

### Port Already in Use

```bash
lsof -i :8080           # Find process using port 8080
kill -9 <PID>           # Kill the process
make stop               # Clean up
make run-all            # Restart
```

### Build Errors

```bash
make clean              # Clean old builds
go mod tidy             # Fix Go dependencies
make install-backend    # Reinstall dependencies
make build              # Rebuild
```

### Frontend Issues

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install             # Fresh install
npm run dev             # Test
```

---

## 🎨 Makefile Output

The Makefile uses colored output:
- 🔵 **Blue**: Info messages
- 🟡 **Yellow**: Warnings, progress
- 🟢 **Green**: Success messages
- 🔴 **Red**: Errors

---

## 📝 Help

Show all available commands:

```bash
make help
```

Or just:

```bash
make
```

---

## 🔗 Related Docs

- [gRPC Explained](GRPC_EXPLAINED.md)
- [Monetization Implementation](MONETIZATION_IMPLEMENTATION.md)
- [Phase 2: Subscriptions](PHASE_2_SUBSCRIPTIONS.md)
- [Production Email Setup](PRODUCTION_EMAIL_SETUP.md)

---

## 💡 Tips

1. **Always run `make stop` before `make restart`** to avoid port conflicts
2. **Use `make logs-follow`** to debug real-time issues
3. **Run `make lint` before committing** to catch issues early
4. **Use `make db-backup` before risky migrations**
5. **Check `make status`** if something seems off
6. **Frontend auto-reloads** with Vite HMR, no need to restart
7. **Backend requires rebuild** (`make restart`) after code changes

---

**Happy coding! 🚀**
