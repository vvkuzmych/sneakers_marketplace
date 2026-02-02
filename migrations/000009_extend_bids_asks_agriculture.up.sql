-- Migration: Extend bids/asks for agricultural commodities
-- Adds delivery window, location, contract type fields
-- Date: 2026-02-02

BEGIN;

-- Extend bids table
ALTER TABLE bids 
    ADD COLUMN IF NOT EXISTS delivery_window_start TIMESTAMP,
    ADD COLUMN IF NOT EXISTS delivery_window_end TIMESTAMP,
    ADD COLUMN IF NOT EXISTS delivery_location TEXT,
    ADD COLUMN IF NOT EXISTS contract_type VARCHAR(50),      -- 'spot', 'forward', 'futures'
    ADD COLUMN IF NOT EXISTS quality_requirements JSONB,
    ADD COLUMN IF NOT EXISTS payment_terms VARCHAR(100);     -- 'Net 30', 'Cash on Delivery'

-- Extend asks table
ALTER TABLE asks 
    ADD COLUMN IF NOT EXISTS delivery_window_start TIMESTAMP,
    ADD COLUMN IF NOT EXISTS delivery_window_end TIMESTAMP,
    ADD COLUMN IF NOT EXISTS delivery_location TEXT,
    ADD COLUMN IF NOT EXISTS contract_type VARCHAR(50),
    ADD COLUMN IF NOT EXISTS quality_specs JSONB,
    ADD COLUMN IF NOT EXISTS payment_terms VARCHAR(100);

-- Indexes for delivery dates (important for agricultural trading)
CREATE INDEX IF NOT EXISTS idx_bids_delivery_window ON bids(delivery_window_start, delivery_window_end);
CREATE INDEX IF NOT EXISTS idx_asks_delivery_window ON asks(delivery_window_start, delivery_window_end);
CREATE INDEX IF NOT EXISTS idx_bids_contract_type ON bids(contract_type);
CREATE INDEX IF NOT EXISTS idx_asks_contract_type ON asks(contract_type);

-- Comments
COMMENT ON COLUMN bids.delivery_window_start IS 'Desired delivery start date (for agricultural)';
COMMENT ON COLUMN bids.contract_type IS 'spot, forward, or futures contract';
COMMENT ON COLUMN bids.quality_requirements IS 'Required quality specs (protein, moisture, etc.)';

COMMIT;
