-- Migration: Assign all existing users to Free plan
-- Phase 2, Day 1 (Part 5/5)
-- Date: 2026-01-23

-- =====================================================
-- ASSIGN ALL EXISTING USERS TO FREE PLAN
-- Ensures backward compatibility - all users get Free tier
-- =====================================================

INSERT INTO user_subscriptions (
    user_id,
    plan_id,
    status,
    billing_cycle,
    current_period_start,
    current_period_end
)
SELECT 
    u.id,
    (SELECT id FROM subscription_plans WHERE name = 'free'),
    'active',
    'lifetime',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + INTERVAL '100 years'  -- Effectively no expiration for free plan
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM user_subscriptions us 
    WHERE us.user_id = u.id AND us.status = 'active'
);

-- Show how many users were assigned
SELECT COUNT(*) AS users_assigned_to_free
FROM user_subscriptions 
WHERE plan_id = (SELECT id FROM subscription_plans WHERE name = 'free');
