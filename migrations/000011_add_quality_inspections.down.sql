-- Rollback: Remove quality_inspections table

BEGIN;

DROP INDEX IF EXISTS idx_quality_inspections_passed;
DROP INDEX IF EXISTS idx_quality_inspections_date;
DROP INDEX IF EXISTS idx_quality_inspections_order;
DROP INDEX IF EXISTS idx_quality_inspections_product;

DROP TABLE IF EXISTS quality_inspections;

COMMIT;
