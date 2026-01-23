-- Migration: Add helper functions and views
-- Phase 2, Day 1 (Part 6/6)
-- Date: 2026-01-23

-- =====================================================
-- AUTO-UPDATE updated_at TRIGGER
-- =====================================================

CREATE OR REPLACE FUNCTION update_subscription_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to subscription_plans
CREATE TRIGGER trigger_subscription_plans_updated_at
    BEFORE UPDATE ON subscription_plans
    FOR EACH ROW
    EXECUTE FUNCTION update_subscription_updated_at();

-- Apply trigger to user_subscriptions
CREATE TRIGGER trigger_user_subscriptions_updated_at
    BEFORE UPDATE ON user_subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_subscription_updated_at();

-- =====================================================
-- HELPER VIEW: Current User Subscriptions
-- Joins users, subscriptions, and plans for easy querying
-- =====================================================

CREATE OR REPLACE VIEW v_user_current_subscriptions AS
SELECT 
    us.id AS subscription_id,
    us.user_id,
    u.email AS user_email,
    u.first_name,
    u.last_name,
    sp.id AS plan_id,
    sp.name AS plan_name,
    sp.display_name AS plan_display_name,
    sp.price_monthly,
    sp.price_yearly,
    sp.buyer_fee_percent,
    sp.seller_fee_percent,
    sp.features,
    us.status,
    us.billing_cycle,
    us.current_period_start,
    us.current_period_end,
    us.cancel_at_period_end,
    us.stripe_subscription_id,
    us.stripe_customer_id,
    CASE 
        WHEN us.trial_end > CURRENT_TIMESTAMP THEN true
        ELSE false
    END AS is_trial_active,
    us.trial_start,
    us.trial_end,
    us.created_at AS subscription_created_at,
    CASE 
        WHEN us.current_period_end < CURRENT_TIMESTAMP THEN true
        ELSE false
    END AS is_expired
FROM user_subscriptions us
JOIN users u ON u.id = us.user_id
JOIN subscription_plans sp ON sp.id = us.plan_id
WHERE us.status = 'active';

COMMENT ON VIEW v_user_current_subscriptions IS 'View of all active user subscriptions with plan details and user info';

-- Show summary
SELECT 
    'Subscription system installed!' AS message,
    (SELECT COUNT(*) FROM subscription_plans) AS total_plans,
    (SELECT COUNT(*) FROM user_subscriptions) AS total_subscriptions,
    (SELECT COUNT(*) FROM v_user_current_subscriptions) AS active_subscriptions;
