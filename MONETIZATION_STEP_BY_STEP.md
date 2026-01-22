# 💰 Монетизація: Покрокова інструкція

## 🎯 Загальний план

```
✅ PHASE 1: Transaction Fees (Week 1-2)
   ├─ Day 1: Database migrations ← ПОЧИНАЄМО ТУТ!
   ├─ Day 2: Backend models & repository
   ├─ Day 3: Fee service & business logic
   └─ Day 4: API endpoints & frontend

⏸️ PHASE 2: Subscriptions (Week 3-4)
⏸️ PHASE 3-10: Інші features
```

---

# PHASE 1, DAY 1: Database Migrations

## ⏱️ Тривалість: 2-3 години

---

## ✅ Step 1: Створити migration files

### Завдання:
Створити SQL файли для додавання fee tables.

### Команди:

```bash
cd /Users/vkuzm/GolandProjects/sneakers_marketplace

# Create migrations directory if not exists
mkdir -p internal/database/migrations

# Create migration files
touch internal/database/migrations/20260123_add_fees_tables.up.sql
touch internal/database/migrations/20260123_add_fees_tables.down.sql
```

### Перевірка:
```bash
ls -la internal/database/migrations/
# Повинно показати:
# 20260123_add_fees_tables.up.sql
# 20260123_add_fees_tables.down.sql
```

---

## ✅ Step 2: Написати UP migration

### Файл: `internal/database/migrations/20260123_add_fees_tables.up.sql`

### Код:

```sql
-- Migration: Add fees tables
-- Author: Sneakers Marketplace Team
-- Date: 2026-01-23
-- Description: Add fee_configs and transaction_fees tables for monetization

BEGIN;

-- ============================================================================
-- TABLE: fee_configs
-- Purpose: Store fee configuration per vertical (sneakers, tickets, etc.)
-- ============================================================================

CREATE TABLE fee_configs (
    id SERIAL PRIMARY KEY,
    vertical VARCHAR(50) NOT NULL,
    
    -- Fee percentages and fixed amounts
    transaction_fee_percent DECIMAL(5,2) NOT NULL DEFAULT 3.00,
    processing_fee_fixed DECIMAL(10,2) NOT NULL DEFAULT 5.00,
    authentication_fee DECIMAL(10,2) NOT NULL DEFAULT 10.00,
    
    -- Shipping fees
    shipping_buyer_charge DECIMAL(10,2) NOT NULL DEFAULT 15.00,
    shipping_seller_cost DECIMAL(10,2) NOT NULL DEFAULT 10.00,
    
    -- Min/max limits
    min_transaction_fee DECIMAL(10,2) NOT NULL DEFAULT 1.00,
    max_transaction_fee DECIMAL(10,2),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Constraints
    UNIQUE(vertical),
    CHECK (transaction_fee_percent >= 0 AND transaction_fee_percent <= 100),
    CHECK (processing_fee_fixed >= 0),
    CHECK (min_transaction_fee >= 0)
);

-- Insert default configurations
INSERT INTO fee_configs (
    vertical, 
    transaction_fee_percent, 
    processing_fee_fixed, 
    authentication_fee,
    shipping_buyer_charge,
    shipping_seller_cost
) VALUES
('sneakers', 3.00, 5.00, 10.00, 15.00, 10.00),
('tickets', 5.00, 3.00, 0.00, 0.00, 0.00);

-- ============================================================================
-- TABLE: transaction_fees
-- Purpose: Record all fees for each transaction (match)
-- ============================================================================

CREATE TABLE transaction_fees (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGINT NOT NULL,
    order_id BIGINT,
    vertical VARCHAR(50) NOT NULL,
    
    -- Original amounts
    sale_price DECIMAL(10,2) NOT NULL,
    
    -- Buyer fees
    buyer_processing_fee DECIMAL(10,2) DEFAULT 0.00,
    buyer_shipping_fee DECIMAL(10,2) DEFAULT 0.00,
    buyer_total DECIMAL(10,2) NOT NULL,
    
    -- Seller fees
    seller_transaction_fee DECIMAL(10,2) NOT NULL,
    seller_authentication_fee DECIMAL(10,2) DEFAULT 0.00,
    seller_shipping_cost DECIMAL(10,2) DEFAULT 0.00,
    seller_payout DECIMAL(10,2) NOT NULL,
    
    -- Platform revenue (total earnings)
    platform_revenue DECIMAL(10,2) NOT NULL,
    
    -- Metadata (store config snapshot for audit trail)
    fee_config_snapshot JSONB,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Constraints
    CHECK (sale_price > 0),
    CHECK (buyer_total >= sale_price),
    CHECK (seller_payout <= sale_price),
    CHECK (platform_revenue >= 0)
);

-- ============================================================================
-- INDEXES
-- Purpose: Optimize queries
-- ============================================================================

CREATE INDEX idx_transaction_fees_match_id ON transaction_fees(match_id);
CREATE INDEX idx_transaction_fees_order_id ON transaction_fees(order_id) WHERE order_id IS NOT NULL;
CREATE INDEX idx_transaction_fees_vertical ON transaction_fees(vertical);
CREATE INDEX idx_transaction_fees_created_at ON transaction_fees(created_at);

-- For revenue reporting
CREATE INDEX idx_transaction_fees_revenue_date ON transaction_fees(created_at, platform_revenue);

-- ============================================================================
-- COMMENTS
-- Purpose: Documentation
-- ============================================================================

COMMENT ON TABLE fee_configs IS 'Fee configuration per vertical (sneakers, tickets)';
COMMENT ON TABLE transaction_fees IS 'Record of all fees charged per transaction';
COMMENT ON COLUMN transaction_fees.platform_revenue IS 'Total revenue earned by platform from this transaction';
COMMENT ON COLUMN transaction_fees.fee_config_snapshot IS 'Snapshot of fee config at time of transaction (for audit)';

COMMIT;

-- ============================================================================
-- VERIFICATION QUERIES (run after migration)
-- ============================================================================

-- SELECT * FROM fee_configs;
-- Should show 2 rows (sneakers, tickets)

-- SELECT COUNT(*) FROM transaction_fees;
-- Should show 0 rows (new table)
```

---

## ✅ Step 3: Написати DOWN migration (rollback)

### Файл: `internal/database/migrations/20260123_add_fees_tables.down.sql`

### Код:

```sql
-- Rollback: Remove fees tables
-- Author: Sneakers Marketplace Team
-- Date: 2026-01-23
-- Description: Rollback fees tables migration

BEGIN;

-- Drop indexes first
DROP INDEX IF EXISTS idx_transaction_fees_revenue_date;
DROP INDEX IF EXISTS idx_transaction_fees_created_at;
DROP INDEX IF EXISTS idx_transaction_fees_vertical;
DROP INDEX IF EXISTS idx_transaction_fees_order_id;
DROP INDEX IF EXISTS idx_transaction_fees_match_id;

-- Drop tables (order matters due to potential future foreign keys)
DROP TABLE IF EXISTS transaction_fees;
DROP TABLE IF EXISTS fee_configs;

COMMIT;
```

---

## ✅ Step 4: Backup existing database

### Завдання:
Створити backup перед запуском міграції.

### Команди:

```bash
# Create backup directory
mkdir -p backups

# Full database backup
pg_dump -U postgres sneakers_marketplace > backups/backup_before_fees_$(date +%Y%m%d_%H%M%S).sql

# Verify backup exists and has content
ls -lh backups/
# Should show backup file with size > 0
```

### Перевірка:
```bash
wc -l backups/backup_before_fees_*.sql
# Should show many lines (not empty)
```

---

## ✅ Step 5: Test migration on LOCAL database

### Завдання:
Запустити міграцію на локальній БД.

### Команди:

```bash
# Connect to database
psql -U postgres -d sneakers_marketplace

# Run UP migration
\i internal/database/migrations/20260123_add_fees_tables.up.sql

# Verify tables created
\dt fee_configs
\dt transaction_fees

# Check data
SELECT * FROM fee_configs;
-- Should show 2 rows (sneakers, tickets)

# Check indexes
\di transaction_fees*
-- Should show multiple indexes

# Exit psql
\q
```

### Expected output:

```
BEGIN
CREATE TABLE
CREATE TABLE
INSERT 0 2
CREATE INDEX
CREATE INDEX
CREATE INDEX
CREATE INDEX
CREATE INDEX
COMMENT
COMMENT
COMMENT
COMMENT
COMMIT
```

---

## ✅ Step 6: Verify migration

### SQL Verification Queries:

```sql
-- 1. Check fee_configs table structure
SELECT 
    column_name, 
    data_type, 
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_name = 'fee_configs'
ORDER BY ordinal_position;

-- Expected: 11 columns (id, vertical, transaction_fee_percent, etc.)

-- 2. Check fee_configs data
SELECT 
    vertical,
    transaction_fee_percent,
    processing_fee_fixed,
    authentication_fee
FROM fee_configs;

-- Expected:
-- sneakers | 3.00 | 5.00 | 10.00
-- tickets  | 5.00 | 3.00 | 0.00

-- 3. Check transaction_fees table structure
SELECT 
    column_name, 
    data_type 
FROM information_schema.columns
WHERE table_name = 'transaction_fees'
ORDER BY ordinal_position;

-- Expected: 15+ columns

-- 4. Check indexes
SELECT 
    indexname, 
    indexdef 
FROM pg_indexes 
WHERE tablename = 'transaction_fees';

-- Expected: 5 indexes

-- 5. Verify constraints
SELECT 
    conname, 
    contype,
    pg_get_constraintdef(oid) as definition
FROM pg_constraint
WHERE conrelid = 'transaction_fees'::regclass;

-- Expected: CHECK constraints on prices
```

---

## ✅ Step 7: Test rollback

### Завдання:
Переконатися, що rollback працює.

### Команди:

```bash
psql -U postgres -d sneakers_marketplace

-- Run DOWN migration
\i internal/database/migrations/20260123_add_fees_tables.down.sql

-- Verify tables removed
\dt fee_configs
-- Should show: Did not find any relation named "fee_configs"

\dt transaction_fees
-- Should show: Did not find any relation named "transaction_fees"

\q
```

---

## ✅ Step 8: Re-apply migration (final)

### Завдання:
Застосувати міграцію остаточно.

### Команди:

```bash
psql -U postgres -d sneakers_marketplace

-- Run UP migration again
\i internal/database/migrations/20260123_add_fees_tables.up.sql

-- Verify
SELECT vertical, transaction_fee_percent FROM fee_configs;

\q
```

---

## ✅ Step 9: Update .env file

### Завдання:
Додати конфігурацію для fees (якщо потрібно).

### Файл: `.env`

### Додати:

```bash
# Monetization settings
FEES_ENABLED=true
DEFAULT_TRANSACTION_FEE_PERCENT=3.0
DEFAULT_PROCESSING_FEE=5.0
```

---

## 🎉 Day 1 Complete Checklist

- [ ] Migration files created (`up.sql` + `down.sql`)
- [ ] Database backed up
- [ ] UP migration tested
- [ ] Tables verified (structure + data)
- [ ] Indexes verified
- [ ] DOWN migration tested (rollback works)
- [ ] UP migration re-applied
- [ ] `.env` updated

---

## 📊 Результат Day 1:

```
✅ Database готова для fees!

Tables created:
  ✅ fee_configs (2 rows: sneakers, tickets)
  ✅ transaction_fees (0 rows, ready for data)

Indexes created:
  ✅ 5 indexes for performance

Ready for Day 2:
  → Backend models & repository
```

---

## 🔜 Next Steps (Day 2)

```
Day 2: Backend Models & Repository
  Step 1: Create fee_config.go model
  Step 2: Create transaction_fee.go model
  Step 3: Create fee_repository.go
  Step 4: Write repository tests
  Step 5: Integration testing
```

---

## 🆘 Troubleshooting

### Problem: Migration fails with "relation already exists"

**Solution:**
```sql
-- Drop existing tables first
DROP TABLE IF EXISTS transaction_fees CASCADE;
DROP TABLE IF EXISTS fee_configs CASCADE;

-- Then re-run UP migration
```

### Problem: Can't connect to database

**Solution:**
```bash
# Check PostgreSQL is running
pg_isready -U postgres

# If not running, start it:
brew services start postgresql@14
# or
sudo systemctl start postgresql
```

### Problem: Permission denied

**Solution:**
```bash
# Grant permissions
psql -U postgres -d sneakers_marketplace -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_user;"
```

---

## 📝 Notes

- Migration file naming: `YYYYMMDD_description.up/down.sql`
- Always test rollback BEFORE production
- Keep backups for at least 30 days
- Document all schema changes

---

**Day 1 готовий для виконання! Чекаю на вашу команду, щоб почати.** 🚀
