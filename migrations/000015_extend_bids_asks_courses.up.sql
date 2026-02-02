-- Migration: Extend bids/asks for courses
-- Adds course-specific fields (group size, learning goals, dates)
-- Date: 2026-02-02

BEGIN;

-- Extend bids table (students)
ALTER TABLE bids 
    ADD COLUMN IF NOT EXISTS desired_start_date TIMESTAMP,
    ADD COLUMN IF NOT EXISTS flexible_dates BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS group_size INTEGER,            -- For group/bulk bids
    ADD COLUMN IF NOT EXISTS learning_goals TEXT[],
    ADD COLUMN IF NOT EXISTS corporate_purchase BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS budget_max DECIMAL(10,2);

-- Extend asks table (instructors)
ALTER TABLE asks 
    ADD COLUMN IF NOT EXISTS min_enrollments INTEGER,       -- Minimum to run course
    ADD COLUMN IF NOT EXISTS early_bird_price DECIMAL(10,2),
    ADD COLUMN IF NOT EXISTS bulk_discount_percent DECIMAL(5,2);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_bids_desired_start_date ON bids(desired_start_date);
CREATE INDEX IF NOT EXISTS idx_bids_corporate_purchase ON bids(corporate_purchase);
CREATE INDEX IF NOT EXISTS idx_asks_min_enrollments ON asks(min_enrollments);

-- Comments
COMMENT ON COLUMN bids.desired_start_date IS 'Preferred course start date (for courses)';
COMMENT ON COLUMN bids.group_size IS 'Number of participants for bulk bid (corporate/group purchases)';
COMMENT ON COLUMN asks.min_enrollments IS 'Minimum students required to run course';

COMMIT;
