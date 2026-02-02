-- Rollback: Remove agricultural_products table

BEGIN;

DROP INDEX IF EXISTS idx_agricultural_products_harvest;
DROP INDEX IF EXISTS idx_agricultural_products_certifications;
DROP INDEX IF EXISTS idx_agricultural_products_grade;
DROP INDEX IF EXISTS idx_agricultural_products_commodity;

DROP TABLE IF EXISTS agricultural_products;

COMMIT;
