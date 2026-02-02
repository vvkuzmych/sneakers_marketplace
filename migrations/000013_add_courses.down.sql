-- Rollback: Remove courses table

BEGIN;

DROP INDEX IF EXISTS idx_courses_format;
DROP INDEX IF EXISTS idx_courses_tags;
DROP INDEX IF EXISTS idx_courses_rating;
DROP INDEX IF EXISTS idx_courses_start_date;
DROP INDEX IF EXISTS idx_courses_status;
DROP INDEX IF EXISTS idx_courses_level;
DROP INDEX IF EXISTS idx_courses_category;
DROP INDEX IF EXISTS idx_courses_instructor;

DROP TABLE IF EXISTS courses;

COMMIT;
