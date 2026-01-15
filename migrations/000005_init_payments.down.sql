-- Drop triggers
DROP TRIGGER IF EXISTS set_payout_id_trigger ON payouts;
DROP TRIGGER IF EXISTS set_payment_id_trigger ON payments;
DROP TRIGGER IF EXISTS update_payouts_updated_at ON payouts;
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;

-- Drop functions
DROP FUNCTION IF EXISTS set_payout_id();
DROP FUNCTION IF EXISTS set_payment_id();
DROP FUNCTION IF EXISTS generate_payout_id();
DROP FUNCTION IF EXISTS generate_payment_id();

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS payouts;
DROP TABLE IF EXISTS payments;
