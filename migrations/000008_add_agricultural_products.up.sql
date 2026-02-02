-- Migration: Add agricultural_products table
-- Vertical: Agricultural Commodities
-- Date: 2026-02-02

BEGIN;

CREATE TABLE agricultural_products (
    id BIGSERIAL PRIMARY KEY,
    vertical VARCHAR(50) DEFAULT 'agriculture',
    
    -- Basic Info
    commodity_type VARCHAR(100) NOT NULL,     -- 'wheat', 'corn', 'soybeans'
    variety VARCHAR(100),                     -- 'Hard Red Winter', 'Yellow Dent'
    grade VARCHAR(50),                        -- 'Grade 1', 'Grade 2'
    
    -- Quantity & Units
    unit_of_measure VARCHAR(20) NOT NULL,     -- 'tons', 'bushels', 'kg'
    min_order_quantity DECIMAL(10,2),
    
    -- Origin
    country_of_origin VARCHAR(100),
    state_province VARCHAR(100),
    farm_name VARCHAR(255),
    
    -- Certifications & Quality (JSONB for flexibility)
    certifications JSONB,                     -- ["organic", "non-gmo"]
    quality_specs JSONB,                      -- {protein: 12.5%, moisture: 13%}
    
    -- Lab Testing
    lab_tested BOOLEAN DEFAULT false,
    lab_certificate_url TEXT,
    test_date TIMESTAMP,
    
    -- Harvest
    harvest_year INTEGER,
    harvest_season VARCHAR(50),
    
    -- Storage
    storage_location VARCHAR(255),
    storage_type VARCHAR(100),                -- 'silo', 'warehouse'
    
    -- Compliance
    usda_certified BOOLEAN DEFAULT false,
    organic_certified BOOLEAN DEFAULT false,
    non_gmo_certified BOOLEAN DEFAULT false,
    
    -- Images
    images JSONB,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CHECK (min_order_quantity IS NULL OR min_order_quantity > 0)
);

-- Indexes
CREATE INDEX idx_agricultural_products_commodity ON agricultural_products(commodity_type);
CREATE INDEX idx_agricultural_products_grade ON agricultural_products(grade);
CREATE INDEX idx_agricultural_products_certifications ON agricultural_products USING gin(certifications);
CREATE INDEX idx_agricultural_products_harvest ON agricultural_products(harvest_year, harvest_season);

-- Comments
COMMENT ON TABLE agricultural_products IS 'Agricultural commodities (wheat, corn, etc.)';
COMMENT ON COLUMN agricultural_products.commodity_type IS 'Type: wheat, corn, soybeans';
COMMENT ON COLUMN agricultural_products.quality_specs IS 'Quality attributes (protein, moisture)';

COMMIT;
