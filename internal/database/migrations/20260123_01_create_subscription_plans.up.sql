-- Migration: Create subscription_plans table
-- Phase 2, Day 1 (Part 1/5)
-- Date: 2026-01-23

-- =====================================================
-- SUBSCRIPTION PLANS TABLE
-- Defines available subscription tiers (Free, Pro, Elite)
-- =====================================================
CREATE TABLE IF NOT EXISTS subscription_plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,           -- 'free', 'pro', 'elite'
    display_name VARCHAR(100) NOT NULL,         -- 'Free', 'Pro', 'Elite'
    description TEXT,
    price_monthly DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    price_yearly DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    
    -- Fee structure for this plan
    buyer_fee_percent DECIMAL(5,2) NOT NULL DEFAULT 1.00,
    seller_fee_percent DECIMAL(5,2) NOT NULL DEFAULT 2.00,
    
    -- Features (JSON array)
    features JSONB DEFAULT '[]'::jsonb,
    
    -- Limits (NULL = unlimited)
    max_active_listings INTEGER,
    max_monthly_transactions INTEGER,
    
    -- Metadata
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,
    stripe_price_id_monthly VARCHAR(255),
    stripe_price_id_yearly VARCHAR(255),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_subscription_plans_name ON subscription_plans(name);
CREATE INDEX idx_subscription_plans_active ON subscription_plans(is_active);

-- Comments
COMMENT ON TABLE subscription_plans IS 'Subscription plan definitions (Free, Pro, Elite)';
COMMENT ON COLUMN subscription_plans.features IS 'JSON array of plan features';
COMMENT ON COLUMN subscription_plans.buyer_fee_percent IS 'Transaction fee for buyers on this plan';
COMMENT ON COLUMN subscription_plans.seller_fee_percent IS 'Transaction fee for sellers on this plan';
