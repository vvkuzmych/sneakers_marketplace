# 🎟️ Multi-Vertical Implementation Plan: Event Tickets

## 🎯 Мета

Додати **Event Tickets** як другу вертикаль до існуючої sneakers marketplace **БЕЗ порушення** поточної функціональності.

---

## 📋 Table of Contents

1. [Бізнес-логіка Event Tickets](#бізнес-логіка-event-tickets)
2. [Відмінності від Sneakers](#відмінності-від-sneakers)
3. [Архітектура Multi-Vertical](#архітектура-multi-vertical)
4. [Database Schema Changes](#database-schema-changes)
5. [Backend Changes](#backend-changes)
6. [Frontend Changes](#frontend-changes)
7. [Поетапний план (15 кроків)](#поетапний-план)
8. [Testing Strategy](#testing-strategy)
9. [Rollback Plan](#rollback-plan)

---

## 🎟️ Бізнес-логіка Event Tickets

### Що таке Event Tickets Marketplace?

**Концепція:** Купівля/продаж квитків на концерти, спортивні події, театр, фестивалі.

### Відмінності від звичайного продажу квитків:

```
Традиційна модель (Ticketmaster):
- Fixed price
- First-come first-served
- Scalpers проблема

Bid-Ask модель (наша):
- Dynamic pricing
- Market determines price
- Anti-scalping механізми
- Transparent pricing
```

### Приклад використання:

```
Концерт: Taylor Swift - Kyiv Arena
Date: 2026-06-15 20:00
Seats: Section A, Row 5, Seats 10-11

BIDS (хочу купити):
User1: $500 for Section A
User2: $450 for Section A
User3: $400 for Section B

ASKS (хочу продати):
User4: $550 for Section A, Row 5, Seat 10
User5: $480 for Section B, Row 10, Seat 15

→ Коли BID ≥ ASK → INSTANT MATCH!
```

### Ключові особливості квитків:

1. **Expiration** - квитки "згорають" після події
2. **Seats** - конкретні місця (Section, Row, Seat)
3. **Transfer** - можливість передачі квитка
4. **Verification** - QR код, barcode
5. **Last-minute** - ціни падають перед подією
6. **No refunds** - як правило, повернення немає

---

## 🔍 Відмінності від Sneakers

### Порівняльна таблиця:

| Характеристика | Sneakers | Event Tickets |
|----------------|----------|---------------|
| **Expiration** | ❌ Немає | ✅ Дата події |
| **Uniques** | Розмір (7-13) | Місце (Section/Row/Seat) |
| **Transferability** | Доставка | Цифрова передача |
| **Physical/Digital** | Фізичний товар | Цифровий |
| **Verification** | Автентифікація | QR код |
| **Price dynamics** | Стабільна | Волатильна (час до події) |
| **Refunds** | Можливі | Немає |
| **Delivery** | 3-7 днів | Миттєво |
| **Inventory** | У продавця | Інтеграція з venues |
| **Match logic** | Розмір + ціна | Seat + ціна |

### Нова бізнес-логіка для квитків:

```go
// Особливості Event Tickets:

1. Expiration logic:
   - Квиток недійсний після події
   - Auto-cancel bids/asks після події
   - Refund logic якщо подія скасована

2. Seat specificity:
   - Section A != Section B
   - Row 1 != Row 10
   - BID може бути на "any seat in Section A"
   - ASK завжди конкретне місце

3. Price dynamics:
   - Ціна падає ближче до події
   - Surge pricing для популярних подій
   - Last-minute deals

4. Transfer mechanics:
   - Instant transfer через QR/barcode
   - Verify ownership
   - Prevent duplicate sales
```

---

## 🏗️ Архітектура Multi-Vertical

### Поточна архітектура (Single Vertical):

```
products table
  → shoe specific fields (size, brand, colorway)
  
bids/asks
  → size_id (shoe sizes)
  
matching logic
  → product_id + size_id + price
```

### Нова архітектура (Multi-Vertical):

```
products table
  → vertical (enum: 'sneakers', 'tickets')
  → vertical_metadata (JSONB) ← гнучке поле
  
bids/asks
  → variant_id (universal: size OR seat)
  → variant_metadata (JSONB)
  
matching logic
  → vertical-aware matching
```

### Database Design:

```sql
-- 1. Add vertical support to products
ALTER TABLE products ADD COLUMN vertical VARCHAR(50) DEFAULT 'sneakers';
ALTER TABLE products ADD COLUMN vertical_metadata JSONB;

-- Sneakers metadata example:
{
  "brand": "Nike",
  "model": "Air Jordan 1",
  "colorway": "Chicago",
  "release_date": "2024-01-15"
}

-- Tickets metadata example:
{
  "event_name": "Taylor Swift Concert",
  "venue": "Kyiv Arena",
  "venue_address": "Kyiv, Ukraine",
  "event_date": "2026-06-15T20:00:00Z",
  "event_type": "concert",
  "artist": "Taylor Swift",
  "min_price": 50,
  "max_price": 5000
}

-- 2. Rename sizes → variants (universal)
ALTER TABLE sizes RENAME TO variants;
ALTER TABLE variants ADD COLUMN vertical VARCHAR(50) DEFAULT 'sneakers';
ALTER TABLE variants ADD COLUMN variant_metadata JSONB;

-- Sneakers variant example:
{
  "size_us": 10,
  "size_eu": 44,
  "size_uk": 9
}

-- Tickets variant example:
{
  "section": "A",
  "row": "5",
  "seat": "10",
  "seat_type": "VIP"  // or "Regular", "Standing"
}

-- 3. Update foreign keys
ALTER TABLE bids RENAME COLUMN size_id TO variant_id;
ALTER TABLE asks RENAME COLUMN size_id TO variant_id;

-- 4. Add expiration for tickets
ALTER TABLE products ADD COLUMN expires_at TIMESTAMP;
CREATE INDEX idx_products_expires_at ON products(expires_at) WHERE expires_at IS NOT NULL;

-- 5. Vertical-specific settings
CREATE TABLE vertical_configs (
  id SERIAL PRIMARY KEY,
  vertical VARCHAR(50) UNIQUE NOT NULL,
  config JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO vertical_configs (vertical, config) VALUES
('sneakers', '{
  "shipping_required": true,
  "authentication_required": true,
  "digital": false,
  "expiration_enabled": false,
  "fee_percentage": 3.0
}'),
('tickets', '{
  "shipping_required": false,
  "authentication_required": false,
  "digital": true,
  "expiration_enabled": true,
  "fee_percentage": 5.0,
  "transfer_instant": true
}');
```

---

## 🔧 Backend Changes

### 1. Create Vertical Package

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

func (v Vertical) GetConfig() Config {
    // Load from database or config
}

func (v Vertical) IsValid() bool {
    return v == VerticalSneakers || v == VerticalTickets
}
```

### 2. Update Product Model

```go
// internal/product/model/product.go
type Product struct {
    ID               int64                  `json:"id"`
    Name             string                 `json:"name"`
    Description      string                 `json:"description"`
    ImageURL         string                 `json:"image_url"`
    Vertical         string                 `json:"vertical"` // NEW
    VerticalMetadata map[string]interface{} `json:"vertical_metadata"` // NEW
    ExpiresAt        *time.Time             `json:"expires_at,omitempty"` // NEW
    CreatedAt        time.Time              `json:"created_at"`
    UpdatedAt        time.Time              `json:"updated_at"`
}

// Sneakers-specific methods
func (p *Product) AsSneaker() *SneakerProduct {
    if p.Vertical != "sneakers" {
        return nil
    }
    return &SneakerProduct{
        Product: p,
        Brand:   p.VerticalMetadata["brand"].(string),
        Model:   p.VerticalMetadata["model"].(string),
        // ...
    }
}

// Tickets-specific methods
func (p *Product) AsTicket() *TicketProduct {
    if p.Vertical != "tickets" {
        return nil
    }
    return &TicketProduct{
        Product:    p,
        EventName:  p.VerticalMetadata["event_name"].(string),
        Venue:      p.VerticalMetadata["venue"].(string),
        EventDate:  p.VerticalMetadata["event_date"].(string),
        // ...
    }
}
```

### 3. Update Variant Model (бувший Size)

```go
// internal/product/model/variant.go
type Variant struct {
    ID              int64                  `json:"id"`
    ProductID       int64                  `json:"product_id"`
    Vertical        string                 `json:"vertical"` // NEW
    VariantMetadata map[string]interface{} `json:"variant_metadata"` // NEW
    CreatedAt       time.Time              `json:"created_at"`
}

// For sneakers
type SneakerVariant struct {
    *Variant
    SizeUS float64 `json:"size_us"`
    SizeEU int     `json:"size_eu"`
    SizeUK float64 `json:"size_uk"`
}

// For tickets
type TicketVariant struct {
    *Variant
    Section  string `json:"section"`
    Row      string `json:"row"`
    Seat     string `json:"seat"`
    SeatType string `json:"seat_type"` // VIP, Regular, Standing
}
```

### 4. Update Matching Logic

```go
// internal/bidding/service/matching.go

func (s *BiddingService) FindMatch(bid *model.Bid) (*model.Ask, error) {
    // Get product to determine vertical
    product, err := s.productRepo.GetByID(bid.ProductID)
    if err != nil {
        return nil, err
    }
    
    // Vertical-specific matching
    switch product.Vertical {
    case "sneakers":
        return s.findSneakerMatch(bid)
    case "tickets":
        return s.findTicketMatch(bid, product)
    default:
        return nil, fmt.Errorf("unknown vertical: %s", product.Vertical)
    }
}

func (s *BiddingService) findSneakerMatch(bid *model.Bid) (*model.Ask, error) {
    // Existing logic
    // Match: product_id + variant_id (size) + price
    return s.repo.FindMatchingAsk(bid.ProductID, bid.VariantID, bid.Price)
}

func (s *BiddingService) findTicketMatch(bid *model.Bid, product *model.Product) (*model.Ask, error) {
    // Ticket-specific logic
    
    // 1. Check if event expired
    if product.ExpiresAt != nil && time.Now().After(*product.ExpiresAt) {
        return nil, fmt.Errorf("event has passed")
    }
    
    // 2. Match logic
    if bid.VariantID != 0 {
        // Specific seat requested
        return s.repo.FindMatchingAsk(bid.ProductID, bid.VariantID, bid.Price)
    } else {
        // "Any seat" - find cheapest available
        return s.repo.FindCheapestAskForProduct(bid.ProductID, bid.Price)
    }
}
```

### 5. Add Expiration Job

```go
// internal/bidding/jobs/expiration_job.go

type ExpirationJob struct {
    biddingRepo *repository.BiddingRepository
    productRepo *repository.ProductRepository
}

func (j *ExpirationJob) Run() {
    // Run every hour
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        j.expireOldEvents()
    }
}

func (j *ExpirationJob) expireOldEvents() {
    // Find expired products (events that passed)
    expiredProducts, err := j.productRepo.FindExpired()
    if err != nil {
        log.Errorf("Failed to find expired products: %v", err)
        return
    }
    
    for _, product := range expiredProducts {
        // Cancel all active bids/asks for this product
        err := j.biddingRepo.CancelAllForProduct(product.ID)
        if err != nil {
            log.Errorf("Failed to cancel bids/asks for product %d: %v", product.ID, err)
            continue
        }
        
        log.Infof("Expired product %d (%s) - cancelled all bids/asks", product.ID, product.Name)
    }
}
```

---

## 💻 Frontend Changes

### 1. Vertical Context

```typescript
// src/contexts/VerticalContext.tsx
import React, { createContext, useContext } from 'react';

type Vertical = 'sneakers' | 'tickets';

interface VerticalConfig {
  shippingRequired: boolean;
  authenticationRequired: boolean;
  digital: boolean;
  expirationEnabled: boolean;
  feePercentage: number;
}

const verticalConfigs: Record<Vertical, VerticalConfig> = {
  sneakers: {
    shippingRequired: true,
    authenticationRequired: true,
    digital: false,
    expirationEnabled: false,
    feePercentage: 3.0,
  },
  tickets: {
    shippingRequired: false,
    authenticationRequired: false,
    digital: true,
    expirationEnabled: true,
    feePercentage: 5.0,
  },
};

interface VerticalContextType {
  vertical: Vertical;
  config: VerticalConfig;
  setVertical: (v: Vertical) => void;
}

const VerticalContext = createContext<VerticalContextType | undefined>(undefined);

export const VerticalProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [vertical, setVertical] = React.useState<Vertical>('sneakers');
  
  const config = verticalConfigs[vertical];
  
  return (
    <VerticalContext.Provider value={{ vertical, config, setVertical }}>
      {children}
    </VerticalContext.Provider>
  );
};

export const useVertical = () => {
  const context = useContext(VerticalContext);
  if (!context) throw new Error('useVertical must be used within VerticalProvider');
  return context;
};
```

### 2. Vertical Selector Component

```typescript
// src/components/VerticalSelector.tsx
import React from 'react';
import { useVertical } from '../contexts/VerticalContext';

export const VerticalSelector: React.FC = () => {
  const { vertical, setVertical } = useVertical();
  
  return (
    <div className="flex gap-4 p-4">
      <button
        className={`px-6 py-3 rounded-lg ${
          vertical === 'sneakers' 
            ? 'bg-blue-600 text-white' 
            : 'bg-gray-200 text-gray-700'
        }`}
        onClick={() => setVertical('sneakers')}
      >
        👟 Sneakers
      </button>
      
      <button
        className={`px-6 py-3 rounded-lg ${
          vertical === 'tickets' 
            ? 'bg-blue-600 text-white' 
            : 'bg-gray-200 text-gray-700'
        }`}
        onClick={() => setVertical('tickets')}
      >
        🎟️ Event Tickets
      </button>
    </div>
  );
};
```

### 3. Vertical-Specific Product Card

```typescript
// src/components/ProductCard.tsx
import React from 'react';
import { Product } from '../types';
import { SneakerCard } from './SneakerCard';
import { TicketCard } from './TicketCard';

interface ProductCardProps {
  product: Product;
}

export const ProductCard: React.FC<ProductCardProps> = ({ product }) => {
  switch (product.vertical) {
    case 'sneakers':
      return <SneakerCard product={product} />;
    case 'tickets':
      return <TicketCard product={product} />;
    default:
      return <div>Unknown product type</div>;
  }
};
```

### 4. Ticket-Specific Components

```typescript
// src/components/TicketCard.tsx
import React from 'react';
import { Product } from '../types';
import { formatDate, getTimeUntilEvent } from '../utils';

interface TicketCardProps {
  product: Product;
}

export const TicketCard: React.FC<TicketCardProps> = ({ product }) => {
  const metadata = product.vertical_metadata;
  const timeUntil = getTimeUntilEvent(metadata.event_date);
  
  return (
    <div className="border rounded-lg p-4 hover:shadow-lg transition">
      <img src={product.image_url} alt={product.name} className="w-full h-48 object-cover rounded" />
      
      <h3 className="text-xl font-bold mt-4">{metadata.event_name}</h3>
      
      <div className="mt-2 text-gray-600">
        <div>📍 {metadata.venue}</div>
        <div>📅 {formatDate(metadata.event_date)}</div>
        <div className="text-red-600 font-semibold">⏰ {timeUntil}</div>
      </div>
      
      <div className="mt-4 flex justify-between items-center">
        <div>
          <div className="text-sm text-gray-500">Lowest Ask</div>
          <div className="text-2xl font-bold">${metadata.lowest_ask || '-'}</div>
        </div>
        
        <button className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700">
          View Tickets
        </button>
      </div>
    </div>
  );
};
```

### 5. Seat Selection Component

```typescript
// src/components/SeatSelector.tsx
import React, { useState } from 'react';

interface Seat {
  id: number;
  section: string;
  row: string;
  seat: string;
  available: boolean;
  price: number;
}

interface SeatSelectorProps {
  productId: number;
  onSeatSelect: (seatId: number) => void;
}

export const SeatSelector: React.FC<SeatSelectorProps> = ({ productId, onSeatSelect }) => {
  const [selectedSection, setSelectedSection] = useState<string>('A');
  
  // Fetch available seats for this section
  // const { data: seats } = useGetSeatsQuery({ productId, section: selectedSection });
  
  return (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium mb-2">Select Section</label>
        <select 
          className="w-full border rounded-lg p-2"
          value={selectedSection}
          onChange={(e) => setSelectedSection(e.target.value)}
        >
          <option value="A">Section A (VIP)</option>
          <option value="B">Section B (Premium)</option>
          <option value="C">Section C (Regular)</option>
        </select>
      </div>
      
      {/* Venue map visualization */}
      <div className="border rounded-lg p-4">
        <div className="text-center mb-4 text-sm text-gray-500">🎭 STAGE</div>
        
        {/* Simplified seat grid */}
        <div className="grid grid-cols-10 gap-2">
          {/* Render seats */}
        </div>
      </div>
    </div>
  );
};
```

---

## 📝 Поетапний план (15 кроків)

### ✅ PHASE 1: Foundation (не ламаємо існуюче)

#### **Крок 1: Database Migration - Add Vertical Support** (2 години)
```bash
# Create migration
cd internal/database/migrations
touch 001_add_vertical_support.sql
```

```sql
-- 001_add_vertical_support.sql
-- Add vertical columns (default 'sneakers' - existing data safe)
ALTER TABLE products ADD COLUMN vertical VARCHAR(50) DEFAULT 'sneakers';
ALTER TABLE products ADD COLUMN vertical_metadata JSONB DEFAULT '{}';
ALTER TABLE products ADD COLUMN expires_at TIMESTAMP;

-- Rename sizes to variants (more universal name)
ALTER TABLE sizes RENAME TO variants;
ALTER TABLE variants ADD COLUMN vertical VARCHAR(50) DEFAULT 'sneakers';
ALTER TABLE variants ADD COLUMN variant_metadata JSONB DEFAULT '{}';

-- Update existing data
UPDATE variants SET variant_metadata = jsonb_build_object(
  'size_us', size_us,
  'size_eu', size_eu,
  'size_uk', size_uk
);

-- Update foreign keys
ALTER TABLE bids RENAME COLUMN size_id TO variant_id;
ALTER TABLE asks RENAME COLUMN size_id TO variant_id;

-- Create vertical configs table
CREATE TABLE vertical_configs (
  id SERIAL PRIMARY KEY,
  vertical VARCHAR(50) UNIQUE NOT NULL,
  config JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Insert configs
INSERT INTO vertical_configs (vertical, config) VALUES
('sneakers', '{"shipping_required": true, "fee_percentage": 3.0}'),
('tickets', '{"shipping_required": false, "fee_percentage": 5.0}');

-- Create indexes
CREATE INDEX idx_products_vertical ON products(vertical);
CREATE INDEX idx_products_expires_at ON products(expires_at) WHERE expires_at IS NOT NULL;
```

**Тестування:**
```bash
# Run migration
make migrate-up

# Verify existing data intact
psql -U postgres -d sneakers_marketplace -c "SELECT COUNT(*) FROM products WHERE vertical = 'sneakers';"
# Should return all existing products

# Verify foreign keys work
psql -U postgres -d sneakers_marketplace -c "SELECT COUNT(*) FROM bids b JOIN variants v ON b.variant_id = v.id;"
# Should return all existing bids
```

**Rollback:**
```sql
-- 001_add_vertical_support_down.sql
ALTER TABLE bids RENAME COLUMN variant_id TO size_id;
ALTER TABLE asks RENAME COLUMN variant_id TO size_id;
ALTER TABLE variants RENAME TO sizes;
ALTER TABLE variants DROP COLUMN vertical;
ALTER TABLE variants DROP COLUMN variant_metadata;
ALTER TABLE products DROP COLUMN vertical;
ALTER TABLE products DROP COLUMN vertical_metadata;
ALTER TABLE products DROP COLUMN expires_at;
DROP TABLE vertical_configs;
```

---

#### **Крок 2: Create Vertical Package** (1 година)

```bash
mkdir -p pkg/vertical
touch pkg/vertical/vertical.go
```

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
}

var configs = map[Vertical]Config{
    VerticalSneakers: {
        ShippingRequired:       true,
        AuthenticationRequired: true,
        Digital:                false,
        ExpirationEnabled:      false,
        FeePercentage:          3.0,
    },
    VerticalTickets: {
        ShippingRequired:       false,
        AuthenticationRequired: false,
        Digital:                true,
        ExpirationEnabled:      true,
        FeePercentage:          5.0,
    },
}

func (v Vertical) GetConfig() Config {
    return configs[v]
}

func (v Vertical) IsValid() bool {
    _, ok := configs[v]
    return ok
}

func (v Vertical) String() string {
    return string(v)
}
```

**Тестування:**
```go
// pkg/vertical/vertical_test.go
func TestVerticalConfig(t *testing.T) {
    sneakersConfig := VerticalSneakers.GetConfig()
    assert.Equal(t, 3.0, sneakersConfig.FeePercentage)
    assert.True(t, sneakersConfig.ShippingRequired)
    
    ticketsConfig := VerticalTickets.GetConfig()
    assert.Equal(t, 5.0, ticketsConfig.FeePercentage)
    assert.False(t, ticketsConfig.ShippingRequired)
}
```

---

#### **Крок 3: Update Product Model (Backward Compatible)** (2 години)

```go
// internal/product/model/product.go
type Product struct {
    ID               int64                  `json:"id"`
    Name             string                 `json:"name"`
    Description      string                 `json:"description"`
    ImageURL         string                 `json:"image_url"`
    Vertical         string                 `json:"vertical"` // NEW
    VerticalMetadata map[string]interface{} `json:"vertical_metadata"` // NEW
    ExpiresAt        *time.Time             `json:"expires_at,omitempty"` // NEW
    CreatedAt        time.Time              `json:"created_at"`
    UpdatedAt        time.Time              `json:"updated_at"`
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

// Sneakers-specific
func (p *Product) IsSneaker() bool {
    return p.Vertical == string(vertical.VerticalSneakers)
}

// Tickets-specific
func (p *Product) IsTicket() bool {
    return p.Vertical == string(vertical.VerticalTickets)
}
```

**Тестування:** Existing sneakers tests should pass без змін!

---

### ✅ PHASE 2: Tickets Infrastructure (паралельно з sneakers)

#### **Крок 4: Create Tickets Service** (3 години)

```bash
mkdir -p internal/tickets
mkdir -p internal/tickets/model
mkdir -p internal/tickets/repository
mkdir -p internal/tickets/service
mkdir -p internal/tickets/handler
```

```go
// internal/tickets/model/ticket.go
package model

import "time"

type TicketProduct struct {
    ProductID    int64     `json:"product_id"`
    EventName    string    `json:"event_name"`
    Venue        string    `json:"venue"`
    VenueAddress string    `json:"venue_address"`
    EventDate    time.Time `json:"event_date"`
    EventType    string    `json:"event_type"` // concert, sports, theater
    MinPrice     float64   `json:"min_price"`
    MaxPrice     float64   `json:"max_price"`
}

type TicketVariant struct {
    VariantID int64  `json:"variant_id"`
    Section   string `json:"section"`
    Row       string `json:"row"`
    Seat      string `json:"seat"`
    SeatType  string `json:"seat_type"` // VIP, Regular, Standing
}
```

---

#### **Крок 5: Add Tickets Endpoints (без втручання в sneakers)** (2 години)

```go
// internal/gateway/routes/routes.go
func SetupRoutes(r *gin.Engine, handlers *Handlers) {
    api := r.Group("/api/v1")
    
    // Existing sneakers routes (не чіпаємо)
    products := api.Group("/products")
    {
        products.GET("", handlers.Product.ListProducts)
        products.GET("/:id", handlers.Product.GetProduct)
    }
    
    // NEW: Tickets routes (паралельно)
    tickets := api.Group("/tickets")
    {
        tickets.GET("", handlers.Ticket.ListEvents)
        tickets.GET("/:id", handlers.Ticket.GetEvent)
        tickets.GET("/:id/seats", handlers.Ticket.GetAvailableSeats)
    }
    
    // Bidding routes (universal - працює для обох)
    bidding := api.Group("/bidding")
    {
        bidding.POST("/bid", handlers.Bidding.PlaceBid)  // Works for both!
        bidding.POST("/ask", handlers.Bidding.PlaceAsk)  // Works for both!
    }
}
```

---

#### **Крок 6: Update Matching Logic (vertical-aware)** (3 години)

```go
// internal/bidding/service/bidding_service.go

func (s *BiddingService) PlaceBid(ctx context.Context, bid *model.Bid) (*model.Bid, *model.Match, error) {
    // Get product to determine vertical
    product, err := s.productRepo.GetByID(ctx, bid.ProductID)
    if err != nil {
        return nil, nil, err
    }
    
    // Vertical-specific validation
    if err := s.validateBidForVertical(bid, product); err != nil {
        return nil, nil, err
    }
    
    // Save bid
    if err := s.repo.PlaceBid(ctx, bid); err != nil {
        return nil, nil, err
    }
    
    // Try to find match (vertical-aware)
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
            return fmt.Errorf("size_id required for sneakers")
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

func (s *BiddingService) findSneakerMatch(ctx context.Context, bid *model.Bid) (*model.Match, error) {
    // Existing logic: exact match on product + size + price
    ask, err := s.repo.FindMatchingAsk(ctx, bid.ProductID, bid.VariantID, bid.Price)
    if err != nil || ask == nil {
        return nil, err
    }
    
    return s.createMatch(ctx, bid, ask)
}

func (s *BiddingService) findTicketMatch(ctx context.Context, bid *model.Bid, product *model.Product) (*model.Match, error) {
    var ask *model.Ask
    var err error
    
    if bid.VariantID != 0 {
        // Specific seat requested
        ask, err = s.repo.FindMatchingAsk(ctx, bid.ProductID, bid.VariantID, bid.Price)
    } else {
        // "Any seat" - find cheapest available
        ask, err = s.repo.FindCheapestAsk(ctx, bid.ProductID, bid.Price)
    }
    
    if err != nil || ask == nil {
        return nil, err
    }
    
    return s.createMatch(ctx, bid, ask)
}
```

**Критично:** Existing sneakers matching НЕ ЗМІНЮЄТЬСЯ!

---

### ✅ PHASE 3: Frontend Multi-Vertical UI

#### **Крок 7: Add Vertical Context** (1 година)

Створити `src/contexts/VerticalContext.tsx` (код вище)

---

#### **Крок 8: Add Vertical Selector** (1 година)

Створити `src/components/VerticalSelector.tsx` (код вище)

---

#### **Крок 9: Create Ticket Components** (4 години)

- `TicketCard.tsx`
- `EventDetailPage.tsx`
- `SeatSelector.tsx`
- `TicketBiddingPage.tsx`

---

### ✅ PHASE 4: Testing & Rollout

#### **Крок 10: Unit Tests** (2 години)

```go
// internal/bidding/service/bidding_service_test.go

func TestPlaceBid_Sneakers_StillWorks(t *testing.T) {
    // Ensure existing sneakers logic unchanged
}

func TestPlaceBid_Tickets_NewLogic(t *testing.T) {
    // Test tickets-specific logic
}

func TestMatchingLogic_Sneakers(t *testing.T) {
    // Exact match on size
}

func TestMatchingLogic_Tickets_SpecificSeat(t *testing.T) {
    // Exact match on seat
}

func TestMatchingLogic_Tickets_AnySeat(t *testing.T) {
    // Find cheapest available
}
```

---

#### **Крок 11: Integration Tests** (2 години)

```bash
# Test sneakers flow (existing)
curl -X POST http://localhost:8080/api/v1/bids \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"product_id": 1, "variant_id": 5, "price": 200}'

# Test tickets flow (new)
curl -X POST http://localhost:8080/api/v1/bids \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"product_id": 100, "variant_id": 0, "price": 500}'
```

---

#### **Крок 12: Seed Tickets Data** (1 година)

```sql
-- Seed first event
INSERT INTO products (name, description, image_url, vertical, vertical_metadata, expires_at)
VALUES (
  'Taylor Swift Concert - Kyiv',
  'The Eras Tour comes to Ukraine!',
  'https://example.com/taylor-swift.jpg',
  'tickets',
  '{
    "event_name": "Taylor Swift - The Eras Tour",
    "venue": "NSC Olimpiyskiy",
    "venue_address": "Kyiv, Ukraine",
    "event_date": "2026-06-15T20:00:00Z",
    "event_type": "concert",
    "artist": "Taylor Swift"
  }',
  '2026-06-15 20:00:00'
);

-- Seed seats for this event
INSERT INTO variants (product_id, vertical, variant_metadata)
VALUES
  (100, 'tickets', '{"section": "A", "row": "1", "seat": "1", "seat_type": "VIP"}'),
  (100, 'tickets', '{"section": "A", "row": "1", "seat": "2", "seat_type": "VIP"}'),
  (100, 'tickets', '{"section": "B", "row": "5", "seat": "10", "seat_type": "Regular"}');
```

---

#### **Крок 13: Expiration Job** (2 години)

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
)

func main() {
    // Connect to DB
    db, err := database.NewPostgresPool(context.Background(), ...)
    if err != nil {
        log.Fatal(err)
    }
    
    biddingRepo := repository.NewBiddingRepository(db)
    productRepo := repository.NewProductRepository(db)
    
    // Run every hour
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    log.Println("Expiration job started")
    
    for range ticker.C {
        expireOldEvents(productRepo, biddingRepo)
    }
}

func expireOldEvents(productRepo *repository.ProductRepository, biddingRepo *repository.BiddingRepository) {
    ctx := context.Background()
    
    // Find expired products
    expiredProducts, err := productRepo.FindExpired(ctx)
    if err != nil {
        log.Printf("Error finding expired products: %v", err)
        return
    }
    
    log.Printf("Found %d expired products", len(expiredProducts))
    
    for _, product := range expiredProducts {
        // Cancel all bids/asks
        err := biddingRepo.CancelAllForProduct(ctx, product.ID)
        if err != nil {
            log.Printf("Error cancelling bids/asks for product %d: %v", product.ID, err)
            continue
        }
        
        log.Printf("Expired product %d: %s", product.ID, product.Name)
    }
}
```

---

#### **Крок 14: Deploy to Staging** (2 години)

```bash
# Build all services
make build

# Run database migrations
make migrate-up

# Start services
./bin/api-gateway &
./bin/user-service &
./bin/product-service &
./bin/bidding-service &
./bin/notification-service &
./bin/expiration-job &  # NEW

# Start frontend
cd frontend && npm run dev
```

---

#### **Крок 15: Gradual Rollout to Production** (1 тиждень)

**Week 1: Sneakers Only**
- Verify existing sneakers functionality
- Monitor metrics

**Week 2: Tickets Beta (Invite-Only)**
- 100 selected users
- Test tickets flow
- Gather feedback

**Week 3: Tickets Public Launch**
- Open to all users
- Marketing campaign
- Monitor performance

---

## 🧪 Testing Strategy

### 1. Unit Tests

```bash
# Backend
go test ./internal/bidding/service/... -v
go test ./internal/tickets/... -v

# Frontend
npm test
```

### 2. Integration Tests

```bash
# Sneakers flow (должен работать как раньше)
./scripts/test-sneakers-flow.sh

# Tickets flow (новый)
./scripts/test-tickets-flow.sh
```

### 3. Load Testing

```bash
# Use k6 or artillery
k6 run load-tests/bidding-sneakers.js
k6 run load-tests/bidding-tickets.js
```

### 4. Manual Testing Checklist

- [ ] Existing sneakers продолжают работать
- [ ] Можно создать BID на sneakers
- [ ] Можно создать ASK на sneakers
- [ ] Sneakers matching работает
- [ ] Можно создать BID на tickets
- [ ] Можно создать ASK на tickets (specific seat)
- [ ] Tickets matching работает (specific seat)
- [ ] Tickets matching работает (any seat)
- [ ] Expiration job отменяет старые tickets
- [ ] Email notifications работают для обоих
- [ ] WebSocket updates работают для обоих

---

## 🔄 Rollback Plan

### If something goes wrong:

```bash
# 1. Stop new services
pkill -f expiration-job

# 2. Rollback database migration
make migrate-down

# 3. Revert code
git revert <commit-hash>

# 4. Redeploy old version
make build && make deploy

# 5. Verify sneakers still works
./scripts/test-sneakers-flow.sh
```

### Database Rollback SQL:

```sql
-- Revert to single vertical (sneakers only)
DELETE FROM products WHERE vertical = 'tickets';
DELETE FROM variants WHERE vertical = 'tickets';

ALTER TABLE bids RENAME COLUMN variant_id TO size_id;
ALTER TABLE asks RENAME COLUMN variant_id TO size_id;
ALTER TABLE variants RENAME TO sizes;

ALTER TABLE products DROP COLUMN vertical;
ALTER TABLE products DROP COLUMN vertical_metadata;
ALTER TABLE products DROP COLUMN expires_at;

DROP TABLE vertical_configs;
```

---

## 📊 Success Metrics

### Week 1 (Sneakers baseline):
- ✅ No regression in sneakers metrics
- ✅ Response time < 200ms
- ✅ Match rate > 15%

### Week 2 (Tickets beta):
- ✅ 100 tickets listed
- ✅ 50+ tickets matched
- ✅ No critical bugs

### Week 3 (Full launch):
- ✅ 1000+ tickets listed
- ✅ 500+ tickets matched
- ✅ User satisfaction > 4.5/5

---

## 🎯 Conclusion

### ✅ Що ми досягли:

1. **Multi-vertical architecture** - підтримка кількох типів товарів
2. **Zero breaking changes** - існуючі sneakers продовжують працювати
3. **Tickets support** - нова вертикаль з унікальною логікою
4. **Scalable** - легко додати 3-ю, 4-ту вертикаль

### 🚀 Наступні кроки:

- Phase 1: Sneakers (✅ done)
- Phase 2: Tickets (📝 цей план)
- Phase 3: Electronics (майбутнє)
- Phase 4: Luxury Goods (майбутнє)

---

**Готовий детальний план для додавання Event Tickets без порушення існуючої системи!** 🎟️✨

**Estimated Time: 3-4 тижні розробки + 1-2 тижні testing**
