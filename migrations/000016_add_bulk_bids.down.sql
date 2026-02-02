-- Rollback: Remove bulk_bids table

BEGIN;

DROP INDEX IF EXISTS idx_bulk_bids_deadline;
DROP INDEX IF EXISTS idx_bulk_bids_status;
DROP INDEX IF EXISTS idx_bulk_bids_organizer;
DROP INDEX IF EXISTS idx_bulk_bids_course;

DROP TABLE IF EXISTS bulk_bids;

COMMIT;
