-- Rollback: Remove default subscription plans
-- Phase 2, Day 1 (Part 4/5)
-- Date: 2026-01-23

DELETE FROM subscription_plans WHERE name IN ('free', 'pro', 'elite');
