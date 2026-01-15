-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    match_id BIGINT NOT NULL REFERENCES matches(id) ON DELETE RESTRICT,
    
    -- Parties
    buyer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    seller_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    
    -- Product details
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    size_id BIGINT NOT NULL REFERENCES sizes(id) ON DELETE RESTRICT,
    
    -- Pricing
    price DECIMAL(10, 2) NOT NULL CHECK (price > 0),
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    
    -- Fees (marketplace takes a cut)
    buyer_fee DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (buyer_fee >= 0),
    seller_fee DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (seller_fee >= 0),
    platform_fee DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (platform_fee >= 0),
    
    total_amount DECIMAL(10, 2) NOT NULL CHECK (total_amount >= price),
    seller_payout DECIMAL(10, 2) NOT NULL CHECK (seller_payout <= price),
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending_payment' CHECK (
        status IN (
            'pending_payment',
            'paid',
            'processing',
            'shipped',
            'delivered',
            'completed',
            'cancelled',
            'refunded'
        )
    ),
    
    -- Shipping
    shipping_address_id BIGINT REFERENCES addresses(id) ON DELETE SET NULL,
    tracking_number VARCHAR(100),
    carrier VARCHAR(50),
    
    -- Timestamps for status changes
    payment_at TIMESTAMP,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    
    -- Notes
    buyer_notes TEXT,
    seller_notes TEXT,
    admin_notes TEXT,
    cancellation_reason TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for orders
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_match_id ON orders(match_id);
CREATE INDEX idx_orders_buyer_id ON orders(buyer_id);
CREATE INDEX idx_orders_seller_id ON orders(seller_id);
CREATE INDEX idx_orders_product_id ON orders(product_id);
CREATE INDEX idx_orders_size_id ON orders(size_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX idx_orders_buyer_status ON orders(buyer_id, status);
CREATE INDEX idx_orders_seller_status ON orders(seller_id, status);

-- Create order_status_history table for audit trail
CREATE TABLE IF NOT EXISTS order_status_history (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    from_status VARCHAR(50),
    to_status VARCHAR(50) NOT NULL,
    note TEXT,
    created_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id, created_at DESC);

-- Trigger for updated_at on orders
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to generate order number
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS TEXT AS $$
DECLARE
    new_number TEXT;
    year_part TEXT;
    seq_part TEXT;
BEGIN
    -- Get current year
    year_part := TO_CHAR(NOW(), 'YYYY');
    
    -- Get next sequence number (padded to 6 digits)
    SELECT LPAD(
        (COUNT(*) + 1)::TEXT, 
        6, 
        '0'
    ) INTO seq_part
    FROM orders
    WHERE EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM NOW());
    
    -- Format: ORD-2026-000001
    new_number := 'ORD-' || year_part || '-' || seq_part;
    
    RETURN new_number;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-generate order_number
CREATE OR REPLACE FUNCTION set_order_number()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.order_number IS NULL OR NEW.order_number = '' THEN
        NEW.order_number := generate_order_number();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_order_number_trigger
    BEFORE INSERT ON orders
    FOR EACH ROW
    EXECUTE FUNCTION set_order_number();

-- Trigger to log status changes
CREATE OR REPLACE FUNCTION log_order_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status != NEW.status THEN
        INSERT INTO order_status_history (
            order_id,
            from_status,
            to_status,
            created_by
        ) VALUES (
            NEW.id,
            OLD.status,
            NEW.status,
            'system'
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER log_order_status_change_trigger
    AFTER UPDATE OF status ON orders
    FOR EACH ROW
    EXECUTE FUNCTION log_order_status_change();

-- Comments for documentation
COMMENT ON TABLE orders IS 'Orders created from matched bids/asks';
COMMENT ON TABLE order_status_history IS 'Audit trail for order status changes';
COMMENT ON COLUMN orders.order_number IS 'Unique order number (e.g., ORD-2026-000001)';
COMMENT ON COLUMN orders.total_amount IS 'price + buyer_fee (what buyer pays)';
COMMENT ON COLUMN orders.seller_payout IS 'price - seller_fee (what seller receives)';
COMMENT ON COLUMN orders.platform_fee IS 'marketplace commission (for reporting)';
