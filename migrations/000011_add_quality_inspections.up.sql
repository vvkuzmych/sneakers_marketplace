-- Migration: Add quality_inspections table
-- For lab testing and quality certification of agricultural products
-- Date: 2026-02-02

BEGIN;

CREATE TABLE quality_inspections (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES agricultural_products(id) ON DELETE CASCADE,
    order_id BIGINT REFERENCES orders(id),
    
    -- Inspector
    inspector_name VARCHAR(255) NOT NULL,
    inspector_license VARCHAR(100),
    inspection_company VARCHAR(255),
    
    -- Inspection details
    inspection_date TIMESTAMP NOT NULL,
    inspection_location TEXT,
    
    -- Test Results (JSONB for flexibility)
    test_results JSONB NOT NULL,
    -- Example: {
    --   "moisture": 13.2,
    --   "protein": 12.8,
    --   "test_weight": 60.5,
    --   "foreign_matter": 0.5
    -- }
    
    -- Grade & Quality
    assigned_grade VARCHAR(50),               -- 'Grade 1', 'Grade 2'
    quality_score DECIMAL(5,2),               -- 0-100
    
    -- Pass/Fail
    passed BOOLEAN NOT NULL,
    failed_reasons TEXT[],
    
    -- Certification
    certificate_number VARCHAR(100) UNIQUE,
    certificate_url TEXT,
    certificate_issued_at TIMESTAMP,
    
    -- Compliance
    usda_approved BOOLEAN DEFAULT false,
    organic_verified BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP DEFAULT NOW(),
    
    CHECK (quality_score IS NULL OR (quality_score >= 0 AND quality_score <= 100))
);

-- Indexes
CREATE INDEX idx_quality_inspections_product ON quality_inspections(product_id);
CREATE INDEX idx_quality_inspections_order ON quality_inspections(order_id);
CREATE INDEX idx_quality_inspections_date ON quality_inspections(inspection_date);
CREATE INDEX idx_quality_inspections_passed ON quality_inspections(passed);

-- Comments
COMMENT ON TABLE quality_inspections IS 'Lab testing and quality certification records';
COMMENT ON COLUMN quality_inspections.test_results IS 'JSONB field with test results (moisture, protein, etc.)';
COMMENT ON COLUMN quality_inspections.passed IS 'Whether inspection passed quality requirements';

COMMIT;
