-- Rollback: Remove course_enrollments table

BEGIN;

DROP INDEX IF EXISTS idx_course_enrollments_access_expires;
DROP INDEX IF EXISTS idx_course_enrollments_completed;
DROP INDEX IF EXISTS idx_course_enrollments_user;
DROP INDEX IF EXISTS idx_course_enrollments_course;

DROP TABLE IF EXISTS course_enrollments;

COMMIT;
