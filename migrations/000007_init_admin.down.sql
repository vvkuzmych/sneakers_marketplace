-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_user_stats ON orders;
DROP FUNCTION IF EXISTS update_user_stats();

-- Drop audit_logs table
DROP TABLE IF EXISTS audit_logs;

-- Remove admin columns from products
ALTER TABLE products DROP COLUMN IF EXISTS is_featured;

-- Remove admin columns from users
ALTER TABLE users DROP COLUMN IF EXISTS last_login;
ALTER TABLE users DROP COLUMN IF EXISTS total_spent;
ALTER TABLE users DROP COLUMN IF EXISTS total_orders;
ALTER TABLE users DROP COLUMN IF EXISTS banned_by;
ALTER TABLE users DROP COLUMN IF EXISTS banned_at;
ALTER TABLE users DROP COLUMN IF EXISTS ban_reason;
ALTER TABLE users DROP COLUMN IF EXISTS is_banned;
ALTER TABLE users DROP COLUMN IF EXISTS role;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_is_banned;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_audit_logs_entity;
DROP INDEX IF EXISTS idx_audit_logs_action_type;
DROP INDEX IF EXISTS idx_audit_logs_admin_id;
