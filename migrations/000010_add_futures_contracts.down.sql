-- Rollback: Remove futures_contracts table

BEGIN;

DROP INDEX IF EXISTS idx_futures_contracts_expiration;
DROP INDEX IF EXISTS idx_futures_contracts_status;
DROP INDEX IF EXISTS idx_futures_contracts_delivery;
DROP INDEX IF EXISTS idx_futures_contracts_seller;
DROP INDEX IF EXISTS idx_futures_contracts_buyer;
DROP INDEX IF EXISTS idx_futures_contracts_product;

DROP TABLE IF EXISTS futures_contracts;

COMMIT;
