-- Notifications table
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,  -- 'match_created', 'order_created', 'order_shipped', etc.
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB,  -- Additional context (order_id, match_id, product_id, etc.)
    
    -- Delivery channels
    email_sent BOOLEAN DEFAULT FALSE,
    email_sent_at TIMESTAMP,
    push_sent BOOLEAN DEFAULT FALSE,
    push_sent_at TIMESTAMP,
    
    -- Read status
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;

-- Notification preferences table
CREATE TABLE notification_preferences (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Email preferences
    email_enabled BOOLEAN DEFAULT TRUE,
    email_match_created BOOLEAN DEFAULT TRUE,
    email_order_created BOOLEAN DEFAULT TRUE,
    email_order_shipped BOOLEAN DEFAULT TRUE,
    email_payment_received BOOLEAN DEFAULT TRUE,
    email_payout_completed BOOLEAN DEFAULT TRUE,
    
    -- Push/WebSocket preferences
    push_enabled BOOLEAN DEFAULT TRUE,
    push_match_created BOOLEAN DEFAULT TRUE,
    push_order_updates BOOLEAN DEFAULT TRUE,
    push_payment_updates BOOLEAN DEFAULT TRUE,
    
    -- In-app preferences
    inapp_enabled BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to update updated_at
CREATE TRIGGER update_notification_preferences_updated_at
    BEFORE UPDATE ON notification_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create default preferences for existing users
INSERT INTO notification_preferences (user_id)
SELECT id FROM users
ON CONFLICT (user_id) DO NOTHING;

COMMENT ON TABLE notifications IS 'User notifications across all channels (email, push, in-app)';
COMMENT ON TABLE notification_preferences IS 'User-specific notification channel preferences';
COMMENT ON COLUMN notifications.data IS 'JSON data with context like {"order_id": 123, "match_id": 456, "product_name": "Nike Air Jordan 1"}';
