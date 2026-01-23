-- Migration: Seed default subscription plans
-- Phase 2, Day 1 (Part 4/5)
-- Date: 2026-01-23

-- =====================================================
-- INSERT DEFAULT SUBSCRIPTION PLANS
-- Free, Pro, Elite with different fee structures
-- =====================================================

INSERT INTO subscription_plans (
    name, 
    display_name, 
    description, 
    price_monthly, 
    price_yearly,
    buyer_fee_percent,
    seller_fee_percent,
    features,
    max_active_listings,
    max_monthly_transactions,
    sort_order
) VALUES 
    -- Free Plan (Default for all users)
    (
        'free',
        'Free',
        'Perfect for casual buyers and sellers',
        0.00,
        0.00,
        1.00,  -- 1% buyer fee
        2.00,  -- 2% seller fee
        '["Basic support", "Standard processing", "Email notifications"]'::jsonb,
        10,    -- Max 10 active listings
        NULL,  -- Unlimited transactions
        1
    ),
    
    -- Pro Plan ($9.99/month or $99.99/year)
    (
        'pro',
        'Pro',
        'Best for active traders',
        9.99,
        99.99,
        1.00,  -- 1% buyer fee (same as Free)
        1.00,  -- 1% seller fee (50% discount!)
        '["Priority support", "Lower seller fees", "Advanced analytics", "Featured listings", "Early access"]'::jsonb,
        50,    -- Max 50 active listings
        NULL,  -- Unlimited transactions
        2
    ),
    
    -- Elite Plan ($49.99/month or $499.99/year)
    (
        'elite',
        'Elite',
        'For power sellers and collectors',
        49.99,
        499.99,
        0.50,  -- 0.5% buyer fee (50% discount!)
        0.50,  -- 0.5% seller fee (75% discount!)
        '["Premium support", "Lowest fees", "Dedicated account manager", "API access", "Custom integrations", "White-glove service", "Priority verification"]'::jsonb,
        NULL,  -- Unlimited listings
        NULL,  -- Unlimited transactions
        3
    )
ON CONFLICT (name) DO NOTHING;

-- Verify plans were created
SELECT 
    name, 
    display_name, 
    price_monthly, 
    buyer_fee_percent, 
    seller_fee_percent 
FROM subscription_plans 
ORDER BY sort_order;
