-- Add admin role and ban functionality to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user' NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_banned BOOLEAN DEFAULT FALSE NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS ban_reason TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS banned_at TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS banned_by INTEGER REFERENCES users(id);

-- Add total orders and total spent for user statistics
ALTER TABLE users ADD COLUMN IF NOT EXISTS total_orders INTEGER DEFAULT 0 NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS total_spent DECIMAL(10,2) DEFAULT 0 NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login TIMESTAMP;

-- Add featured flag to products for admin highlighting
ALTER TABLE products ADD COLUMN IF NOT EXISTS is_featured BOOLEAN DEFAULT FALSE NOT NULL;

-- Create index on user role for fast admin queries
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_banned ON users(is_banned);

-- Create audit_logs table for tracking admin actions
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER NOT NULL REFERENCES users(id),
    action_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT NOT NULL,
    details JSONB DEFAULT '{}'::jsonb,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for audit logs
CREATE INDEX IF NOT EXISTS idx_audit_logs_admin_id ON audit_logs(admin_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action_type ON audit_logs(action_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- Insert first admin user (update this with your email)
-- Password: admin123 (hashed with bcrypt cost 12)
INSERT INTO users (email, password_hash, first_name, last_name, phone, role, is_active, created_at, updated_at)
VALUES (
    'admin@sneakersmarketplace.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIRXlw3Pf.', -- admin123
    'Admin',
    'User',
    '+1234567890',
    'admin',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;

-- Create function to update user statistics
CREATE OR REPLACE FUNCTION update_user_stats()
RETURNS TRIGGER AS $$
BEGIN
    -- Update buyer stats on order creation
    IF NEW.status IN ('paid', 'processing', 'shipped', 'delivered') THEN
        UPDATE users 
        SET 
            total_orders = total_orders + 1,
            total_spent = total_spent + NEW.total
        WHERE id = NEW.buyer_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-update user stats
DROP TRIGGER IF EXISTS trigger_update_user_stats ON orders;
CREATE TRIGGER trigger_update_user_stats
    AFTER INSERT OR UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_user_stats();

-- Comments for documentation
COMMENT ON COLUMN users.role IS 'User role: user, admin';
COMMENT ON COLUMN users.is_banned IS 'Whether user is banned from the platform';
COMMENT ON COLUMN users.ban_reason IS 'Reason for ban (if banned)';
COMMENT ON COLUMN users.banned_at IS 'Timestamp when user was banned';
COMMENT ON COLUMN users.banned_by IS 'Admin user who performed the ban';
COMMENT ON COLUMN users.total_orders IS 'Total number of orders completed by user';
COMMENT ON COLUMN users.total_spent IS 'Total amount spent on orders';
COMMENT ON COLUMN users.last_login IS 'Last login timestamp';

COMMENT ON TABLE audit_logs IS 'Admin action audit trail';
COMMENT ON COLUMN audit_logs.action_type IS 'Type of action: user_banned, user_deleted, order_cancelled, etc.';
COMMENT ON COLUMN audit_logs.entity_type IS 'Entity affected: user, order, product';
COMMENT ON COLUMN audit_logs.entity_id IS 'ID of the affected entity';
COMMENT ON COLUMN audit_logs.details IS 'JSON details about the action';
