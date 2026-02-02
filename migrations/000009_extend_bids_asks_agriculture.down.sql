-- Rollback: Remove agricultural extensions from bids/asks

BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS idx_asks_contract_type;
DROP INDEX IF EXISTS idx_bids_contract_type;
DROP INDEX IF EXISTS idx_asks_delivery_window;
DROP INDEX IF EXISTS idx_bids_delivery_window;

-- Drop columns from bids
ALTER TABLE bids 
    DROP COLUMN IF EXISTS payment_terms,
    DROP COLUMN IF EXISTS quality_requirements,
    DROP COLUMN IF EXISTS contract_type,
    DROP COLUMN IF EXISTS delivery_location,
    DROP COLUMN IF EXISTS delivery_window_end,
    DROP COLUMN IF EXISTS delivery_window_start;

-- Drop columns from asks
ALTER TABLE asks 
    DROP COLUMN IF EXISTS payment_terms,
    DROP COLUMN IF EXISTS quality_specs,
    DROP COLUMN IF EXISTS contract_type,
    DROP COLUMN IF EXISTS delivery_location,
    DROP COLUMN IF EXISTS delivery_window_end,
    DROP COLUMN IF EXISTS delivery_window_start;

COMMIT;
