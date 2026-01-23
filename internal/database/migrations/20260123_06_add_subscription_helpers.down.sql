-- Rollback: Remove helper functions and views
-- Phase 2, Day 1 (Part 6/6)
-- Date: 2026-01-23

-- Drop view
DROP VIEW IF EXISTS v_user_current_subscriptions;

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_user_subscriptions_updated_at ON user_subscriptions;
DROP TRIGGER IF EXISTS trigger_subscription_plans_updated_at ON subscription_plans;

-- Drop function
DROP FUNCTION IF EXISTS update_subscription_updated_at();
