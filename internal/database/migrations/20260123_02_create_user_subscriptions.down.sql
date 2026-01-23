-- Rollback: Drop user_subscriptions table
-- Phase 2, Day 1 (Part 2/5)
-- Date: 2026-01-23

DROP TABLE IF EXISTS user_subscriptions CASCADE;
