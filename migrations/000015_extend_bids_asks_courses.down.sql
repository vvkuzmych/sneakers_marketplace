-- Rollback: Remove course extensions from bids/asks

BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS idx_asks_min_enrollments;
DROP INDEX IF EXISTS idx_bids_corporate_purchase;
DROP INDEX IF EXISTS idx_bids_desired_start_date;

-- Drop columns from bids
ALTER TABLE bids 
    DROP COLUMN IF EXISTS budget_max,
    DROP COLUMN IF EXISTS corporate_purchase,
    DROP COLUMN IF EXISTS learning_goals,
    DROP COLUMN IF EXISTS group_size,
    DROP COLUMN IF EXISTS flexible_dates,
    DROP COLUMN IF EXISTS desired_start_date;

-- Drop columns from asks
ALTER TABLE asks 
    DROP COLUMN IF EXISTS bulk_discount_percent,
    DROP COLUMN IF EXISTS early_bird_price,
    DROP COLUMN IF EXISTS min_enrollments;

COMMIT;
