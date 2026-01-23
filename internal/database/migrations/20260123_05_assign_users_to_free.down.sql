-- Rollback: Remove free subscriptions assigned to users
-- Phase 2, Day 1 (Part 5/5)
-- Date: 2026-01-23

DELETE FROM user_subscriptions 
WHERE plan_id = (SELECT id FROM subscription_plans WHERE name = 'free')
AND billing_cycle = 'lifetime';
