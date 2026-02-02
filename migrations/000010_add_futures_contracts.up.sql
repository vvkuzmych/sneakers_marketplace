-- Migration: Add futures_contracts table
-- For forward and futures contracts in agricultural trading
-- Date: 2026-02-02

BEGIN;

CREATE TABLE futures_contracts (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES agricultural_products(id) ON DELETE CASCADE,
    
    -- Contract details
    contract_number VARCHAR(100) UNIQUE NOT NULL,
    contract_type VARCHAR(50) NOT NULL,       -- 'forward', 'futures'
    
    -- Pricing
    strike_price DECIMAL(10,2) NOT NULL,      -- Agreed price per unit
    quantity DECIMAL(10,2) NOT NULL,          -- Amount (tons, bushels)
    unit_of_measure VARCHAR(20) NOT NULL,
    
    -- Dates
    contract_date TIMESTAMP DEFAULT NOW(),
    delivery_date TIMESTAMP NOT NULL,
    expiration_date TIMESTAMP NOT NULL,
    
    -- Parties
    buyer_id BIGINT NOT NULL REFERENCES users(id),
    seller_id BIGINT NOT NULL REFERENCES users(id),
    
    -- Status
    status VARCHAR(50) NOT NULL,              -- 'active', 'settled', 'expired', 'cancelled'
    settled_price DECIMAL(10,2),
    settled_at TIMESTAMP,
    
    -- Margin & Collateral
    margin_requirement DECIMAL(10,2),
    margin_posted DECIMAL(10,2),
    
    -- Delivery
    delivery_location TEXT,
    delivery_status VARCHAR(50),              -- 'pending', 'in_transit', 'delivered'
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CHECK (strike_price > 0),
    CHECK (quantity > 0),
    CHECK (expiration_date > contract_date)
);

-- Indexes
CREATE INDEX idx_futures_contracts_product ON futures_contracts(product_id);
CREATE INDEX idx_futures_contracts_buyer ON futures_contracts(buyer_id);
CREATE INDEX idx_futures_contracts_seller ON futures_contracts(seller_id);
CREATE INDEX idx_futures_contracts_delivery ON futures_contracts(delivery_date);
CREATE INDEX idx_futures_contracts_status ON futures_contracts(status);
CREATE INDEX idx_futures_contracts_expiration ON futures_contracts(expiration_date);

-- Comments
COMMENT ON TABLE futures_contracts IS 'Forward and futures contracts for agricultural commodities';
COMMENT ON COLUMN futures_contracts.strike_price IS 'Agreed price per unit';
COMMENT ON COLUMN futures_contracts.margin_requirement IS 'Required margin for contract';

COMMIT;
