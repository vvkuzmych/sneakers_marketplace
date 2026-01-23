-- Migration: Create subscription_transactions table
-- Phase 2, Day 1 (Part 3/5)
-- Date: 2026-01-23

-- =====================================================
-- SUBSCRIPTION TRANSACTIONS TABLE
-- Payment history for subscriptions
-- =====================================================
CREATE TABLE IF NOT EXISTS subscription_transactions (
    id SERIAL PRIMARY KEY,
    user_subscription_id INTEGER NOT NULL REFERENCES user_subscriptions(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id INTEGER NOT NULL REFERENCES subscription_plans(id),
    
    -- Transaction details
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Type and status
    transaction_type VARCHAR(20) NOT NULL,
    -- Options: 'payment', 'refund', 'upgrade', 'downgrade'
    status VARCHAR(20) NOT NULL,
    -- Options: 'pending', 'succeeded', 'failed', 'refunded'
    
    -- Stripe references
    stripe_payment_intent_id VARCHAR(255),
    stripe_invoice_id VARCHAR(255),
    stripe_charge_id VARCHAR(255),
    
    -- Additional info
    description TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP
);

-- Indexes for queries
CREATE INDEX idx_subscription_transactions_user_subscription_id ON subscription_transactions(user_subscription_id);
CREATE INDEX idx_subscription_transactions_user_id ON subscription_transactions(user_id);
CREATE INDEX idx_subscription_transactions_status ON subscription_transactions(status);
CREATE INDEX idx_subscription_transactions_stripe_payment_intent_id ON subscription_transactions(stripe_payment_intent_id);
CREATE INDEX idx_subscription_transactions_created_at ON subscription_transactions(created_at DESC);

-- Comments
COMMENT ON TABLE subscription_transactions IS 'Subscription payment and billing history';
COMMENT ON COLUMN subscription_transactions.transaction_type IS 'payment, refund, upgrade, downgrade';
COMMENT ON COLUMN subscription_transactions.status IS 'pending, succeeded, failed, refunded';
