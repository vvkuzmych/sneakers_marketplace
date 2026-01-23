-- Migration: Create user_subscriptions table
-- Phase 2, Day 1 (Part 2/5)
-- Date: 2026-01-23

-- =====================================================
-- USER SUBSCRIPTIONS TABLE
-- Tracks individual user subscription instances
-- =====================================================
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id INTEGER NOT NULL REFERENCES subscription_plans(id),
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    -- Options: 'active', 'cancelled', 'expired', 'past_due', 'trialing'
    
    -- Billing
    billing_cycle VARCHAR(20) NOT NULL,  -- 'monthly', 'yearly', 'lifetime'
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    cancel_at_period_end BOOLEAN NOT NULL DEFAULT false,
    
    -- Stripe integration
    stripe_subscription_id VARCHAR(255),
    stripe_customer_id VARCHAR(255),
    
    -- Trial period
    trial_start TIMESTAMP,
    trial_end TIMESTAMP,
    
    -- Additional metadata
    metadata JSONB DEFAULT '{}'::jsonb,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cancelled_at TIMESTAMP,
    
    -- Ensure only one active subscription per user
    CONSTRAINT unique_active_user_subscription UNIQUE (user_id, status)
);

-- Indexes for performance
CREATE INDEX idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX idx_user_subscriptions_status ON user_subscriptions(status);
CREATE INDEX idx_user_subscriptions_plan_id ON user_subscriptions(plan_id);
CREATE INDEX idx_user_subscriptions_stripe_subscription_id ON user_subscriptions(stripe_subscription_id);
CREATE INDEX idx_user_subscriptions_period_end ON user_subscriptions(current_period_end);

-- Comments
COMMENT ON TABLE user_subscriptions IS 'User subscription records';
COMMENT ON COLUMN user_subscriptions.status IS 'active, cancelled, expired, past_due, trialing';
COMMENT ON COLUMN user_subscriptions.cancel_at_period_end IS 'If true, subscription will cancel at period end';
COMMENT ON COLUMN user_subscriptions.billing_cycle IS 'monthly, yearly, or lifetime (for free plan)';
