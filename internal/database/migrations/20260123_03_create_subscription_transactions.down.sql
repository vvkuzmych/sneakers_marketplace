-- Rollback: Drop subscription_transactions table
-- Phase 2, Day 1 (Part 3/5)
-- Date: 2026-01-23

DROP TABLE IF EXISTS subscription_transactions CASCADE;
