# Database Migrations

This directory contains SQL migrations for the Sneakers Marketplace database.

## Migrations

### 000001 - Users & Authentication
- Creates `users` table
- Creates `addresses` table (shipping/billing)
- Creates `sessions` table (JWT token management)
- Adds indexes and triggers

### 000002 - Products & Inventory
- Creates `products` table (sneaker catalog)
- Creates `product_images` table
- Creates `sizes` table (size-based inventory)
- Creates `inventory_transactions` table (audit trail)
- Adds automatic inventory logging

## Running Migrations

### Install golang-migrate

```bash
# macOS
brew install golang-migrate

# Or download from: https://github.com/golang-migrate/migrate
```

### Apply Migrations

```bash
# Set database URL
export DATABASE_URL="postgres://postgres:postgres@localhost:5435/sneakers_marketplace?sslmode=disable"

# Run all migrations
migrate -path migrations -database "${DATABASE_URL}" up

# Run specific number of migrations
migrate -path migrations -database "${DATABASE_URL}" up 1
```

### Rollback Migrations

```bash
# Rollback all
migrate -path migrations -database "${DATABASE_URL}" down

# Rollback one
migrate -path migrations -database "${DATABASE_URL}" down 1
```

### Check Migration Status

```bash
migrate -path migrations -database "${DATABASE_URL}" version
```

### Create New Migration

```bash
migrate create -ext sql -dir migrations -seq add_new_feature
```

## Using Makefile

```bash
# Run migrations
make migrate-up

# Rollback
make migrate-down

# Create new migration
make migrate-create name=add_new_feature
```

## Database Schema

### Users
- Authentication and user profiles
- Multiple addresses per user
- Session management for JWT

### Products
- Sneaker catalog
- Multiple images per product
- Size-based inventory (each size is tracked separately)
- Automatic inventory transaction logging

## Next Migrations

Planned migrations:
- 000003 - Bids & Asks (auction system)
- 000004 - Orders & Matches
- 000005 - Payments & Transactions
- 000006 - Notifications
- 000007 - Analytics & Reports
