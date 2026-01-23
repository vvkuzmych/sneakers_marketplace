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
