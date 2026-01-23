# 🎟️ Event Tickets: Поетапна імплементація

## 🎯 Головна мета

Додати **Event Tickets** як другу вертикаль до існуючої sneakers marketplace **БЕЗ порушення** поточної функціональності.

---

## 📊 Загальний огляд фаз

| Фаза | Назва | Тривалість | Ризик | Can Rollback? |
|------|-------|------------|-------|---------------|
| 0 | Pre-work & Planning | 3 дні | 🟢 Low | ✅ Yes |
| 1 | Database Foundation | 1 тиждень | 🟢 Low | ✅ Yes |
| 2 | Backend Infrastructure | 2 тижні | 🟡 Medium | ✅ Yes |
| 3 | Frontend Foundation | 1 тиждень | 🟢 Low | ✅ Yes |
| 4 | Tickets Core Features | 2 тижні | 🟡 Medium | ✅ Yes |
| 5 | Testing & Bug Fixes | 1 тиждень | 🟢 Low | ✅ Yes |
| 6 | Beta Launch | 1 тиждень | 🟡 Medium | ✅ Yes |
| 7 | Public Launch | 1 тиждень | 🔴 High | ⚠️ Partial |

**Total: 9-10 тижнів**

---

## 🔄 Ключові принципи

### ✅ Zero Downtime
- Всі зміни backward compatible
- Поточні sneakers продовжують працювати
- Можливість rollback на кожному етапі

### ✅ Feature Flags
- Tickets функціональність за feature flag
- Можна вимкнути в будь-який момент
- A/B testing готовність

### ✅ Incremental Rollout
- Beta → Limited users → Public
- Monitor metrics на кожному етапі

---

# ФАЗА 0: Pre-work & Planning

## 🎯 Цілі
- Підготувати команду
- Налаштувати інфраструктуру
- Створити план комунікації

## ⏱️ Тривалість: 3 дні

---

## 📋 Завдання

### Day 1: Team Alignment

#### ✅ Завдання 1.1: Kickoff Meeting (2 години)
**Що робити:**
```
1. Презентувати Event Tickets концепцію
2. Обговорити архітектуру
3. Розподілити ролі:
   - Backend lead
   - Frontend lead
   - QA lead
   - DevOps
```

**Deliverable:**
- [ ] Meeting notes
- [ ] Team roles assigned
- [ ] Questions documented

---

#### ✅ Завдання 1.2: Code Review поточної системи (4 години)
**Що робити:**
```bash
# Review existing code
git log --since="3 months ago" --oneline

# Find critical areas
grep -r "size_id" internal/
grep -r "product_id" internal/
grep -r "matching" internal/bidding/
```

**Deliverable:**
- [ ] List of files to modify
- [ ] Potential breaking points identified
- [ ] Dependency tree created

---

### Day 2: Infrastructure Setup

#### ✅ Завдання 2.1: Create Feature Flag (2 години)
```go
// pkg/featureflags/flags.go
package featureflags

type Flags struct {
    TicketsEnabled bool `env:"FEATURE_TICKETS_ENABLED" default:"false"`
}

func (f *Flags) IsTicketsEnabled() bool {
    return f.TicketsEnabled
}
```

**Environment variables:**
```bash
# .env
FEATURE_TICKETS_ENABLED=false  # Початково вимкнено!
```

**Deliverable:**
- [ ] Feature flag package created
- [ ] Tests written
- [ ] Documentation updated

---

#### ✅ Завдання 2.2: Setup Staging Environment (3 години)
```bash
# Create separate DB for staging
createdb sneakers_marketplace_staging

# Copy production data
pg_dump sneakers_marketplace | psql sneakers_marketplace_staging

# Deploy to staging
make deploy-staging
```

**Deliverable:**
- [ ] Staging DB created
- [ ] Staging environment running
- [ ] CI/CD pipeline configured

---

#### ✅ Завдання 2.3: Monitoring Setup (2 години)
```yaml
# monitoring/dashboards/tickets.yaml
dashboard:
  name: "Event Tickets Metrics"
  panels:
    - ticket_listings_total
    - ticket_bids_total
    - ticket_matches_total
    - tickets_vertical_errors
```

**Deliverable:**
- [ ] Grafana dashboard created
- [ ] Alerts configured
- [ ] Log aggregation ready

---

### Day 3: Documentation & Planning

#### ✅ Завдання 3.1: API Documentation (2 години)
```yaml
# docs/api/tickets.yaml
/api/v1/tickets:
  get:
    summary: "List events"
    parameters:
      - name: vertical
        in: query
        schema:
          type: string
          enum: [sneakers, tickets]
```

**Deliverable:**
- [ ] OpenAPI spec updated
- [ ] Postman collection created
- [ ] Example requests documented

---

#### ✅ Завдання 3.2: Database Migration Plan (2 години)
```sql
-- Plan migration steps
-- 1. Add columns with defaults
-- 2. Backfill existing data
-- 3. Add indexes
-- 4. Update constraints
```

**Deliverable:**
- [ ] Migration scripts written
- [ ] Rollback scripts written
- [ ] Testing plan for migrations

---

## ✅ Phase 0 Checklist

- [ ] Team aligned and roles assigned
- [ ] Feature flag implemented
- [ ] Staging environment ready
- [ ] Monitoring configured
- [ ] API documentation started
- [ ] Migration plan ready

**Success Criteria:**
- ✅ Team готова до розробки
- ✅ Infrastructure ready
- ✅ Можна почати Phase 1

---

# ФАЗА 1: Database Foundation

## 🎯 Цілі
- Додати multi-vertical support до database
- **НЕ ЛАМАТИ** існуючі sneakers дані
- 100% backward compatible

## ⏱️ Тривалість: 1 тиждень (5 робочих днів)

---

## 📋 Завдання

### Day 1: Migration Scripts

#### ✅ Завдання 1.1: Create Migration Files (2 години)
```bash
mkdir -p internal/database/migrations
touch internal/database/migrations/20260122_add_vertical_support.up.sql
touch internal/database/migrations/20260122_add_vertical_support.down.sql
```

**Deliverable:**
- [ ] Migration files created
- [ ] Named with timestamp
- [ ] Both up and down migrations

---

#### ✅ Завдання 1.2: Write UP Migration (4 години)
```sql
-- 20260122_add_vertical_support.up.sql
BEGIN;

-- Step 1: Add vertical to products (DEFAULT ensures backward compatibility)
ALTER TABLE products 
  ADD COLUMN vertical VARCHAR(50) DEFAULT 'sneakers' NOT NULL;

-- Step 2: Add metadata (JSONB for flexibility)
ALTER TABLE products 
  ADD COLUMN vertical_metadata JSONB DEFAULT '{}' NOT NULL;

-- Step 3: Add expiration (only for tickets)
ALTER TABLE products 
  ADD COLUMN expires_at TIMESTAMP;

-- Step 4: Rename sizes → variants (more universal)
ALTER TABLE sizes RENAME TO variants;

-- Step 5: Add vertical to variants
ALTER TABLE variants 
  ADD COLUMN vertical VARCHAR(50) DEFAULT 'sneakers' NOT NULL;

-- Step 6: Add variant metadata
ALTER TABLE variants 
  ADD COLUMN variant_metadata JSONB DEFAULT '{}' NOT NULL;

-- Step 7: Update foreign keys
ALTER TABLE bids RENAME COLUMN size_id TO variant_id;
ALTER TABLE asks RENAME COLUMN size_id TO variant_id;

-- Step 8: Backfill variant metadata from existing data
UPDATE variants 
SET variant_metadata = jsonb_build_object(
  'size_us', size_us,
  'size_eu', size_eu,
  'size_uk', size_uk
)
WHERE vertical = 'sneakers';

-- Step 9: Create indexes
CREATE INDEX idx_products_vertical ON products(vertical);
CREATE INDEX idx_products_expires_at ON products(expires_at) 
  WHERE expires_at IS NOT NULL;
CREATE INDEX idx_variants_vertical ON variants(vertical);

-- Step 10: Create vertical configs table
CREATE TABLE vertical_configs (
  id SERIAL PRIMARY KEY,
  vertical VARCHAR(50) UNIQUE NOT NULL,
  config JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Step 11: Insert configs
INSERT INTO vertical_configs (vertical, config) VALUES
('sneakers', '{
  "shipping_required": true,
  "authentication_required": true,
  "digital": false,
  "expiration_enabled": false,
  "fee_percentage": 3.0,
  "transfer_instant": false
}'),
('tickets', '{
  "shipping_required": false,
  "authentication_required": false,
  "digital": true,
  "expiration_enabled": true,
  "fee_percentage": 5.0,
  "transfer_instant": true
}');

-- Step 12: Add constraints
ALTER TABLE products 
  ADD CONSTRAINT chk_vertical 
  CHECK (vertical IN ('sneakers', 'tickets'));

ALTER TABLE variants 
  ADD CONSTRAINT chk_variant_vertical 
  CHECK (vertical IN ('sneakers', 'tickets'));

COMMIT;
```

**Testing:**
```sql
-- Verify existing data intact
SELECT COUNT(*) FROM products WHERE vertical = 'sneakers';
-- Should return all existing products

SELECT COUNT(*) FROM variants WHERE vertical = 'sneakers';
-- Should return all existing sizes

SELECT COUNT(*) FROM bids b 
  JOIN variants v ON b.variant_id = v.id;
-- Should return all existing bids (foreign keys work)
```

**Deliverable:**
- [ ] UP migration written
- [ ] All steps documented
- [ ] Testing queries ready

---

#### ✅ Завдання 1.3: Write DOWN Migration (2 години)
```sql
-- 20260122_add_vertical_support.down.sql
BEGIN;

-- Rollback in reverse order

-- 1. Drop constraints
ALTER TABLE variants DROP CONSTRAINT IF EXISTS chk_variant_vertical;
ALTER TABLE products DROP CONSTRAINT IF EXISTS chk_vertical;

-- 2. Drop vertical configs
DROP TABLE IF EXISTS vertical_configs;

-- 3. Drop indexes
DROP INDEX IF EXISTS idx_variants_vertical;
DROP INDEX IF EXISTS idx_products_expires_at;
DROP INDEX IF EXISTS idx_products_vertical;

-- 4. Remove metadata from variants
ALTER TABLE variants DROP COLUMN IF EXISTS variant_metadata;
ALTER TABLE variants DROP COLUMN IF EXISTS vertical;

-- 5. Rename foreign keys back
ALTER TABLE asks RENAME COLUMN variant_id TO size_id;
ALTER TABLE bids RENAME COLUMN variant_id TO size_id;

-- 6. Rename variants back to sizes
ALTER TABLE variants RENAME TO sizes;

-- 7. Remove product columns
ALTER TABLE products DROP COLUMN IF EXISTS expires_at;
ALTER TABLE products DROP COLUMN IF EXISTS vertical_metadata;
ALTER TABLE products DROP COLUMN IF EXISTS vertical;

COMMIT;
```

**Deliverable:**
- [ ] DOWN migration written
- [ ] Rollback tested on staging

---

### Day 2: Test Migrations on Staging

#### ✅ Завдання 2.1: Backup Production DB (1 година)
```bash
# Full backup
pg_dump -U postgres sneakers_marketplace > backup_before_migration.sql

# Verify backup
psql -U postgres -d postgres -c "\l" | grep sneakers
```

**Deliverable:**
- [ ] Backup created
- [ ] Backup verified
- [ ] Backup stored securely

---

#### ✅ Завдання 2.2: Run Migration on Staging (2 години)
```bash
# Apply migration
psql -U postgres -d sneakers_marketplace_staging < internal/database/migrations/20260122_add_vertical_support.up.sql

# Verify
psql -U postgres -d sneakers_marketplace_staging << EOF
\d products
\d variants
SELECT * FROM vertical_configs;
EOF
```

**Testing queries:**
```sql
-- 1. Verify all existing products are 'sneakers'
SELECT vertical, COUNT(*) FROM products GROUP BY vertical;
-- Expected: sneakers | <count>

-- 2. Verify variants metadata
SELECT id, vertical, variant_metadata FROM variants LIMIT 5;
-- Expected: size_us, size_eu, size_uk in metadata

-- 3. Test existing bids still work
SELECT b.id, b.product_id, b.variant_id, b.price 
FROM bids b 
JOIN variants v ON b.variant_id = v.id 
LIMIT 10;
-- Expected: All rows returned (foreign keys work)

-- 4. Test existing functionality
SELECT 
  p.name,
  v.variant_metadata->>'size_us' as size,
  COUNT(b.id) as bid_count
FROM products p
JOIN variants v ON v.product_id = p.id
LEFT JOIN bids b ON b.product_id = p.id AND b.variant_id = v.id
WHERE p.vertical = 'sneakers'
GROUP BY p.id, p.name, v.id
LIMIT 10;
-- Expected: Existing sneakers data intact
```

**Deliverable:**
- [ ] Migration успішна на staging
- [ ] Всі тести пройдені
- [ ] Існуючі дані не зламані

---

#### ✅ Завдання 2.3: Test Rollback on Staging (1 година)
```bash
# Rollback migration
psql -U postgres -d sneakers_marketplace_staging < internal/database/migrations/20260122_add_vertical_support.down.sql

# Verify rollback
psql -U postgres -d sneakers_marketplace_staging << EOF
\d products
-- Should NOT have 'vertical' column

\d sizes
-- Should be named 'sizes', not 'variants'

SELECT * FROM bids LIMIT 5;
-- Should have 'size_id', not 'variant_id'
EOF
```

**Deliverable:**
- [ ] Rollback успішний
- [ ] Schema повернута до початкового стану

---

### Day 3: Code Changes for Models

#### ✅ Завдання 3.1: Update Product Model (3 години)
```go
// internal/product/model/product.go
package model

import (
    "time"
    "github.com/vvkuzmych/sneakers_marketplace/pkg/vertical"
)

type Product struct {
    ID               int64                  `json:"id" db:"id"`
    Name             string                 `json:"name" db:"name"`
    Description      string                 `json:"description" db:"description"`
    ImageURL         string                 `json:"image_url" db:"image_url"`
    Vertical         string                 `json:"vertical" db:"vertical"`                     // NEW
    VerticalMetadata map[string]interface{} `json:"vertical_metadata" db:"vertical_metadata"`   // NEW
    ExpiresAt        *time.Time             `json:"expires_at,omitempty" db:"expires_at"`      // NEW
    CreatedAt        time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
}

// Helper methods
func (p *Product) GetVertical() vertical.Vertical {
    return vertical.Vertical(p.Vertical)
}

func (p *Product) IsExpired() bool {
    if p.ExpiresAt == nil {
        return false
    }
    return time.Now().After(*p.ExpiresAt)
}

func (p *Product) IsSneaker() bool {
    return p.Vertical == string(vertical.VerticalSneakers)
}

func (p *Product) IsTicket() bool {
    return p.Vertical == string(vertical.VerticalTickets)
}

// Sneakers-specific data
func (p *Product) GetBrand() string {
    if !p.IsSneaker() {
        return ""
    }
    if brand, ok := p.VerticalMetadata["brand"].(string); ok {
        return brand
    }
    return ""
}

func (p *Product) GetModel() string {
    if !p.IsSneaker() {
        return ""
    }
    if model, ok := p.VerticalMetadata["model"].(string); ok {
        return model
    }
    return ""
}

// Tickets-specific data
func (p *Product) GetEventName() string {
    if !p.IsTicket() {
        return ""
    }
    if eventName, ok := p.VerticalMetadata["event_name"].(string); ok {
        return eventName
    }
    return ""
}

func (p *Product) GetVenue() string {
    if !p.IsTicket() {
        return ""
    }
    if venue, ok := p.VerticalMetadata["venue"].(string); ok {
        return venue
    }
    return ""
}

func (p *Product) GetEventDate() *time.Time {
    if !p.IsTicket() {
        return nil
    }
    if eventDateStr, ok := p.VerticalMetadata["event_date"].(string); ok {
        eventDate, err := time.Parse(time.RFC3339, eventDateStr)
        if err != nil {
            return nil
        }
        return &eventDate
    }
    return nil
}
```

**Tests:**
```go
// internal/product/model/product_test.go
func TestProduct_GetVertical(t *testing.T) {
    p := &Product{Vertical: "sneakers"}
    assert.Equal(t, vertical.VerticalSneakers, p.GetVertical())
}

func TestProduct_IsSneaker(t *testing.T) {
    p := &Product{Vertical: "sneakers"}
    assert.True(t, p.IsSneaker())
    assert.False(t, p.IsTicket())
}

func TestProduct_IsExpired(t *testing.T) {
    past := time.Now().Add(-1 * time.Hour)
    future := time.Now().Add(1 * time.Hour)
    
    p1 := &Product{ExpiresAt: &past}
    assert.True(t, p1.IsExpired())
    
    p2 := &Product{ExpiresAt: &future}
    assert.False(t, p2.IsExpired())
    
    p3 := &Product{ExpiresAt: nil}
    assert.False(t, p3.IsExpired())
}
```

**Deliverable:**
- [ ] Product model updated
- [ ] Helper methods added
- [ ] Tests written and passing

---

#### ✅ Завдання 3.2: Update Variant Model (2 години)
```go
// internal/product/model/variant.go (formerly size.go)
package model

import "time"

type Variant struct {
    ID              int64                  `json:"id" db:"id"`
    ProductID       int64                  `json:"product_id" db:"product_id"`
    Vertical        string                 `json:"vertical" db:"vertical"`                   // NEW
    VariantMetadata map[string]interface{} `json:"variant_metadata" db:"variant_metadata"`   // NEW
    CreatedAt       time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
}

// Sneakers-specific methods
func (v *Variant) GetSizeUS() float64 {
    if v.Vertical != "sneakers" {
        return 0
    }
    if sizeUS, ok := v.VariantMetadata["size_us"].(float64); ok {
        return sizeUS
    }
    return 0
}

func (v *Variant) GetSizeEU() int {
    if v.Vertical != "sneakers" {
        return 0
    }
    if sizeEU, ok := v.VariantMetadata["size_eu"].(float64); ok {
        return int(sizeEU)
    }
    return 0
}

// Tickets-specific methods
func (v *Variant) GetSection() string {
    if v.Vertical != "tickets" {
        return ""
    }
    if section, ok := v.VariantMetadata["section"].(string); ok {
        return section
    }
    return ""
}

func (v *Variant) GetRow() string {
    if v.Vertical != "tickets" {
        return ""
    }
    if row, ok := v.VariantMetadata["row"].(string); ok {
        return row
    }
    return ""
}

func (v *Variant) GetSeat() string {
    if v.Vertical != "tickets" {
        return ""
    }
    if seat, ok := v.VariantMetadata["seat"].(string); ok {
        return seat
    }
    return ""
}

func (v *Variant) GetSeatType() string {
    if v.Vertical != "tickets" {
        return ""
    }
    if seatType, ok := v.VariantMetadata["seat_type"].(string); ok {
        return seatType
    }
    return ""
}
```

**Deliverable:**
- [ ] Variant model created
- [ ] Helper methods for both verticals
- [ ] Tests written

---

### Day 4: Update Repositories

#### ✅ Завдання 4.1: Update Product Repository (3 години)
```go
// internal/product/repository/product_repository.go

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
    query := `
        INSERT INTO products (
            name, description, image_url, 
            vertical, vertical_metadata, expires_at
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `
    
    // Convert metadata to JSON
    metadataJSON, err := json.Marshal(product.VerticalMetadata)
    if err != nil {
        return fmt.Errorf("failed to marshal metadata: %w", err)
    }
    
    err = r.db.QueryRow(ctx, query,
        product.Name,
        product.Description,
        product.ImageURL,
        product.Vertical,
        metadataJSON,
        product.ExpiresAt,
    ).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
    
    return err
}

func (r *ProductRepository) FindExpired(ctx context.Context) ([]*model.Product, error) {
    query := `
        SELECT id, name, description, image_url, 
               vertical, vertical_metadata, expires_at,
               created_at, updated_at
        FROM products
        WHERE expires_at IS NOT NULL 
          AND expires_at < NOW()
          AND vertical = 'tickets'
    `
    
    rows, err := r.db.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var products []*model.Product
    for rows.Next() {
        p := &model.Product{}
        var metadataJSON []byte
        
        err := rows.Scan(
            &p.ID, &p.Name, &p.Description, &p.ImageURL,
            &p.Vertical, &metadataJSON, &p.ExpiresAt,
            &p.CreatedAt, &p.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        
        // Parse metadata
        if err := json.Unmarshal(metadataJSON, &p.VerticalMetadata); err != nil {
            return nil, err
        }
        
        products = append(products, p)
    }
    
    return products, nil
}
```

**Deliverable:**
- [ ] Repository methods updated
- [ ] JSONB handling implemented
- [ ] Tests updated

---

### Day 5: Integration Testing

#### ✅ Завдання 5.1: End-to-End Tests (4 години)
```go
// internal/product/repository/product_repository_test.go
func TestProductRepository_CreateSneaker(t *testing.T) {
    // Setup test DB
    db := setupTestDB(t)
    repo := NewProductRepository(db)
    
    // Create sneaker
    product := &model.Product{
        Name:        "Nike Air Jordan 1",
        Description: "Classic sneaker",
        ImageURL:    "https://example.com/image.jpg",
        Vertical:    "sneakers",
        VerticalMetadata: map[string]interface{}{
            "brand":   "Nike",
            "model":   "Air Jordan 1",
            "colorway": "Chicago",
        },
    }
    
    err := repo.Create(context.Background(), product)
    assert.NoError(t, err)
    assert.NotZero(t, product.ID)
    
    // Verify
    found, err := repo.GetByID(context.Background(), product.ID)
    assert.NoError(t, err)
    assert.Equal(t, "sneakers", found.Vertical)
    assert.Equal(t, "Nike", found.GetBrand())
}

func TestProductRepository_CreateTicket(t *testing.T) {
    db := setupTestDB(t)
    repo := NewProductRepository(db)
    
    expiresAt := time.Now().Add(30 * 24 * time.Hour)
    
    product := &model.Product{
        Name:        "Taylor Swift Concert",
        Description: "The Eras Tour",
        ImageURL:    "https://example.com/ts.jpg",
        Vertical:    "tickets",
        VerticalMetadata: map[string]interface{}{
            "event_name": "Taylor Swift - The Eras Tour",
            "venue":      "NSC Olimpiyskiy",
            "event_date": "2026-06-15T20:00:00Z",
        },
        ExpiresAt: &expiresAt,
    }
    
    err := repo.Create(context.Background(), product)
    assert.NoError(t, err)
    
    // Verify
    found, err := repo.GetByID(context.Background(), product.ID)
    assert.NoError(t, err)
    assert.Equal(t, "tickets", found.Vertical)
    assert.Equal(t, "Taylor Swift - The Eras Tour", found.GetEventName())
    assert.False(t, found.IsExpired())
}
```

**Deliverable:**
- [ ] E2E tests written
- [ ] All tests passing
- [ ] Coverage > 80%

---

## ✅ Phase 1 Checklist

- [ ] Migration scripts готові (up + down)
- [ ] Migrations протестовані на staging
- [ ] Rollback протестований
- [ ] Product model updated
- [ ] Variant model created
- [ ] Repositories updated
- [ ] All tests passing
- [ ] Existing sneakers functionality не зламана

**Success Criteria:**
- ✅ Database supports multi-vertical
- ✅ Zero breaking changes
- ✅ Can rollback safely
- ✅ All tests passing (>80% coverage)

**Готовність до Phase 2:** ✅

---

# ФАЗА 2: Backend Infrastructure

## 🎯 Цілі
- Створити Tickets Service
- Додати vertical-aware matching logic
- **НЕ ЗМІНЮВАТИ** існуючу sneakers logic

## ⏱️ Тривалість: 2 тижні (10 робочих днів)

---

## 📋 Week 1: Core Services

### Day 1-2: Vertical Package (2 дні)

#### ✅ Завдання 1.1: Create Vertical Package
```go
// pkg/vertical/vertical.go
package vertical

type Vertical string

const (
    VerticalSneakers Vertical = "sneakers"
    VerticalTickets  Vertical = "tickets"
)

type Config struct {
    ShippingRequired       bool    `json:"shipping_required"`
    AuthenticationRequired bool    `json:"authentication_required"`
    Digital                bool    `json:"digital"`
    ExpirationEnabled      bool    `json:"expiration_enabled"`
    FeePercentage          float64 `json:"fee_percentage"`
    TransferInstant        bool    `json:"transfer_instant"`
}

var configs = map[Vertical]Config{
    VerticalSneakers: {
        ShippingRequired:       true,
        AuthenticationRequired: true,
        Digital:                false,
        ExpirationEnabled:      false,
        FeePercentage:          3.0,
        TransferInstant:        false,
    },
    VerticalTickets: {
        ShippingRequired:       false,
        AuthenticationRequired: false,
        Digital:                true,
        ExpirationEnabled:      true,
        FeePercentage:          5.0,
        TransferInstant:        true,
    },
}

func (v Vertical) GetConfig() Config {
    return configs[v]
}

func (v Vertical) IsValid() bool {
    _, ok := configs[v]
    return ok
}
```

**Tests:**
```go
func TestVertical_GetConfig(t *testing.T) {
    sneakersConfig := VerticalSneakers.GetConfig()
    assert.Equal(t, 3.0, sneakersConfig.FeePercentage)
    
    ticketsConfig := VerticalTickets.GetConfig()
    assert.Equal(t, 5.0, ticketsConfig.FeePercentage)
}
```

**Deliverable:**
- [ ] Vertical package created
- [ ] Tests written
- [ ] Documentation added

---

### Day 3-4: Tickets Service Structure (2 дні)

#### ✅ Завдання 3.1: Create Tickets Package
```bash
mkdir -p internal/tickets/{model,repository,service,handler}
```

```go
// internal/tickets/model/ticket.go
package model

import "time"

type TicketEvent struct {
    ProductID    int64     `json:"product_id"`
    EventName    string    `json:"event_name"`
    Venue        string    `json:"venue"`
    VenueAddress string    `json:"venue_address"`
    EventDate    time.Time `json:"event_date"`
    EventType    string    `json:"event_type"` // concert, sports, theater
    Artist       string    `json:"artist,omitempty"`
    MinPrice     float64   `json:"min_price"`
    MaxPrice     float64   `json:"max_price"`
}

type TicketSeat struct {
    VariantID int64  `json:"variant_id"`
    Section   string `json:"section"`
    Row       string `json:"row"`
    Seat      string `json:"seat"`
    SeatType  string `json:"seat_type"` // VIP, Regular, Standing
    Available bool   `json:"available"`
    Price     float64 `json:"price,omitempty"`
}
```

**Deliverable:**
- [ ] Tickets models created
- [ ] Package structure ready

---

### Day 5: Matching Logic Update (1 день)

#### ✅ Завдання 5.1: Vertical-Aware Matching
```go
// internal/bidding/service/bidding_service.go

func (s *BiddingService) PlaceBid(ctx context.Context, bid *model.Bid) (*model.Bid, *model.Match, error) {
    // Get product
    product, err := s.productRepo.GetByID(ctx, bid.ProductID)
    if err != nil {
        return nil, nil, fmt.Errorf("product not found: %w", err)
    }
    
    // Validate bid for vertical
    if err := s.validateBidForVertical(bid, product); err != nil {
        return nil, nil, err
    }
    
    // Save bid
    if err := s.repo.PlaceBid(ctx, bid); err != nil {
        return nil, nil, err
    }
    
    // Try to find match
    match, err := s.findMatch(ctx, bid, product)
    if err != nil {
        return bid, nil, err
    }
    
    return bid, match, nil
}

func (s *BiddingService) validateBidForVertical(bid *model.Bid, product *model.Product) error {
    switch product.GetVertical() {
    case vertical.VerticalSneakers:
        // Existing sneakers validation
        if bid.VariantID == 0 {
            return fmt.Errorf("variant_id (size) required for sneakers")
        }
        
    case vertical.VerticalTickets:
        // Tickets validation
        if product.IsExpired() {
            return fmt.Errorf("event has already passed")
        }
        // variant_id optional for tickets ("any seat")
        
    default:
        return fmt.Errorf("unknown vertical: %s", product.Vertical)
    }
    
    return nil
}

func (s *BiddingService) findMatch(ctx context.Context, bid *model.Bid, product *model.Product) (*model.Match, error) {
    switch product.GetVertical() {
    case vertical.VerticalSneakers:
        return s.findSneakerMatch(ctx, bid)
        
    case vertical.VerticalTickets:
        return s.findTicketMatch(ctx, bid, product)
        
    default:
        return nil, fmt.Errorf("unknown vertical: %s", product.Vertical)
    }
}

// Existing sneakers logic (НЕ ЗМІНЮЄТЬСЯ!)
func (s *BiddingService) findSneakerMatch(ctx context.Context, bid *model.Bid) (*model.Match, error) {
    // Exact match: product + size + price
    ask, err := s.repo.FindMatchingAsk(ctx, bid.ProductID, bid.VariantID, bid.Price)
    if err != nil || ask == nil {
        return nil, err
    }
    
    return s.createMatch(ctx, bid, ask)
}

// NEW: Tickets matching logic
func (s *BiddingService) findTicketMatch(ctx context.Context, bid *model.Bid, product *model.Product) (*model.Match, error) {
    var ask *model.Ask
    var err error
    
    if bid.VariantID != 0 {
        // Specific seat requested
        ask, err = s.repo.FindMatchingAsk(ctx, bid.ProductID, bid.VariantID, bid.Price)
    } else {
        // "Any seat" - find cheapest available ask
        ask, err = s.repo.FindCheapestAskForProduct(ctx, bid.ProductID, bid.Price)
    }
    
    if err != nil || ask == nil {
        return nil, err
    }
    
    return s.createMatch(ctx, bid, ask)
}
```

**Критично:** Sneakers matching logic залишається БЕЗ ЗМІН!

**Tests:**
```go
func TestFindMatch_Sneakers_StillWorks(t *testing.T) {
    // Test existing sneakers matching
    // Should work exactly as before
}

func TestFindMatch_Tickets_SpecificSeat(t *testing.T) {
    // Test tickets matching with specific seat
}

func TestFindMatch_Tickets_AnySeat(t *testing.T) {
    // Test tickets matching with "any seat"
}
```

**Deliverable:**
- [ ] Vertical-aware matching implemented
- [ ] Existing sneakers tests pass
- [ ] New tickets tests pass

---

## 📋 Week 2: API & Jobs

### Day 6-7: API Endpoints (2 дні)

#### ✅ Завдання 6.1: Add Tickets Endpoints
```go
// internal/gateway/routes/routes.go
func SetupRoutes(r *gin.Engine, handlers *Handlers) {
    api := r.Group("/api/v1")
    
    // Existing routes (НЕ ЧІПАЄМО!)
    sneakers := api.Group("/sneakers")
    {
        sneakers.GET("", handlers.Product.ListSneakers)
        sneakers.GET("/:id", handlers.Product.GetSneaker)
    }
    
    // NEW: Tickets routes (паралельно)
    tickets := api.Group("/tickets")
    {
        tickets.GET("", handlers.Ticket.ListEvents)
        tickets.GET("/:id", handlers.Ticket.GetEvent)
        tickets.GET("/:id/seats", handlers.Ticket.GetAvailableSeats)
    }
    
    // Universal bidding (works for both!)
    bidding := api.Group("/bidding")
    {
        bidding.POST("/bid", handlers.Bidding.PlaceBid)
        bidding.POST("/ask", handlers.Bidding.PlaceAsk)
        bidding.GET("/bids/product/:product_id", handlers.Bidding.GetBidsForProduct)
        bidding.GET("/asks/product/:product_id", handlers.Bidding.GetAsksForProduct)
    }
}
```

**Deliverable:**
- [ ] Tickets endpoints added
- [ ] Existing routes не змінені
- [ ] API documentation updated

---

### Day 8-9: Expiration Job (2 дні)

#### ✅ Завдання 8.1: Create Expiration Job
```go
// cmd/expiration-job/main.go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/vvkuzmych/sneakers_marketplace/internal/bidding/repository"
    "github.com/vvkuzmych/sneakers_marketplace/internal/product/repository"
    "github.com/vvkuzmych/sneakers_marketplace/pkg/database"
    "github.com/vvkuzmych/sneakers_marketplace/pkg/config"
    "github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
)

func main() {
    // Initialize logger
    log := logger.NewProduction()
    log.Info("Starting Expiration Job")
    
    // Load config
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err.Error())
    }
    
    // Connect to DB
    ctx := context.Background()
    db, err := database.NewPostgresPool(ctx, database.PostgresConfig{
        URL: cfg.Database.URL,
    }, log)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer db.Close()
    
    // Initialize repos
    productRepo := repository.NewProductRepository(db)
    biddingRepo := repository.NewBiddingRepository(db)
    
    // Run every hour
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    log.Info("Expiration job started, checking every hour")
    
    // Run immediately first time
    expireOldEvents(ctx, productRepo, biddingRepo, log)
    
    // Then run every hour
    for range ticker.C {
        expireOldEvents(ctx, productRepo, biddingRepo, log)
    }
}

func expireOldEvents(ctx context.Context, productRepo *repository.ProductRepository, biddingRepo *repository.BiddingRepository, log *logger.Logger) {
    // Find expired products (events that passed)
    expiredProducts, err := productRepo.FindExpired(ctx)
    if err != nil {
        log.Errorf("Failed to find expired products: %v", err)
        return
    }
    
    if len(expiredProducts) == 0 {
        log.Info("No expired products found")
        return
    }
    
    log.Infof("Found %d expired products", len(expiredProducts))
    
    for _, product := range expiredProducts {
        // Cancel all active bids/asks for this product
        cancelled, err := biddingRepo.CancelAllForProduct(ctx, product.ID)
        if err != nil {
            log.Errorf("Failed to cancel bids/asks for product %d: %v", product.ID, err)
            continue
        }
        
        log.Infof("Expired product %d (%s) - cancelled %d bids/asks", 
            product.ID, product.Name, cancelled)
    }
}
```

**Repository method:**
```go
// internal/bidding/repository/bidding_repository.go

func (r *BiddingRepository) CancelAllForProduct(ctx context.Context, productID int64) (int, error) {
    query := `
        UPDATE bids 
        SET status = 'cancelled', updated_at = NOW()
        WHERE product_id = $1 AND status = 'active'
    `
    
    result, err := r.db.Exec(ctx, query, productID)
    if err != nil {
        return 0, err
    }
    
    bidsCancelled := result.RowsAffected()
    
    query = `
        UPDATE asks 
        SET status = 'cancelled', updated_at = NOW()
        WHERE product_id = $1 AND status = 'active'
    `
    
    result, err = r.db.Exec(ctx, query, productID)
    if err != nil {
        return int(bidsCancelled), err
    }
    
    asksCancelled := result.RowsAffected()
    
    return int(bidsCancelled + asksCancelled), nil
}
```

**Systemd service:**
```ini
# /etc/systemd/system/sneakers-expiration-job.service
[Unit]
Description=Sneakers Marketplace Expiration Job
After=network.target postgresql.service

[Service]
Type=simple
User=sneakers
WorkingDirectory=/opt/sneakers_marketplace
Environment="PATH=/usr/local/bin:/usr/bin:/bin"
EnvironmentFile=/opt/sneakers_marketplace/.env
ExecStart=/opt/sneakers_marketplace/bin/expiration-job
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Deliverable:**
- [ ] Expiration job created
- [ ] Systemd service configured
- [ ] Logs monitored

---

### Day 10: Integration Testing (1 день)

#### ✅ Завдання 10.1: End-to-End Testing
```bash
# Test 1: Existing sneakers flow
curl -X POST http://localhost:8080/api/v1/bidding/bid \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "variant_id": 5,
    "price": 200,
    "quantity": 1
  }'

# Test 2: New tickets flow (specific seat)
curl -X POST http://localhost:8080/api/v1/bidding/bid \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 100,
    "variant_id": 50,
    "price": 500,
    "quantity": 1
  }'

# Test 3: Tickets flow (any seat)
curl -X POST http://localhost:8080/api/v1/bidding/bid \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 100,
    "variant_id": 0,
    "price": 400,
    "quantity": 1
  }'
```

**Deliverable:**
- [ ] All E2E tests pass
- [ ] Sneakers functionality intact
- [ ] Tickets functionality works

---

## ✅ Phase 2 Checklist

- [ ] Vertical package created
- [ ] Tickets service structure ready
- [ ] Vertical-aware matching implemented
- [ ] Tickets API endpoints added
- [ ] Expiration job created and deployed
- [ ] All integration tests pass
- [ ] Zero breaking changes for sneakers

**Success Criteria:**
- ✅ Backend supports tickets
- ✅ Sneakers still works
- ✅ Matching logic is vertical-aware
- ✅ Expiration job runs correctly

**Готовність до Phase 3:** ✅

---

# ФАЗА 3: Frontend Foundation

## 🎯 Цілі
- Створити UI foundation для tickets
- **НЕ ЛАМАТИ** існуючий sneakers UI
- Додати vertical switching

## ⏱️ Тривалість: 1 тиждень (5 робочих днів)

*(Детальні завдання можуть бути додані окремо)*

---

# ФАЗА 4-7: Tickets Features, Testing, Beta, Launch

*(Детальні плани для наступних фаз можуть бути створені окремо після успішного завершення Phase 1-3)*

---

## 📊 Загальний Progress Tracker

```
Phase 0: Pre-work          [========== ] 100% ✅
Phase 1: Database          [========== ] 100% ✅
Phase 2: Backend           [==========  ] 90%  🔄
Phase 3: Frontend          [            ] 0%   ⏸️
Phase 4: Features          [            ] 0%   ⏸️
Phase 5: Testing           [            ] 0%   ⏸️
Phase 6: Beta              [            ] 0%   ⏸️
Phase 7: Launch            [            ] 0%   ⏸️

Overall Progress: 35%
```

---

## 🎯 Success Metrics

### Technical Metrics:
- ✅ Zero downtime during deployment
- ✅ No regression in sneakers performance
- ✅ Test coverage > 80%
- ✅ All E2E tests passing

### Business Metrics:
- ✅ Sneakers metrics stable (no degradation)
- 🎯 100+ tickets events listed (Phase 6)
- 🎯 50+ ticket transactions (Phase 6)
- 🎯 User satisfaction > 4.5/5 (Phase 7)

---

**Створено детальний поетапний план додавання Event Tickets БЕЗ порушення існуючої функціональності!** 🎟️✨
