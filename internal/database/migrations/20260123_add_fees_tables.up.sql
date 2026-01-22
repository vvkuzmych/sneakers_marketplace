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
