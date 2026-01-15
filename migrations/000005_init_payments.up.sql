-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    payment_id VARCHAR(100) UNIQUE NOT NULL,
    
    -- Relations
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    
    -- Stripe details
    stripe_payment_intent_id VARCHAR(255) UNIQUE,
    stripe_charge_id VARCHAR(255),
    stripe_customer_id VARCHAR(255),
    
    -- Amount
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (
        status IN (
            'pending',
            'processing',
            'succeeded',
            'failed',
            'cancelled',
            'refunded',
            'partially_refunded'
        )
    ),
    
    -- Payment method
    payment_method VARCHAR(50),
    card_last4 VARCHAR(4),
    card_brand VARCHAR(20),
    
    -- Refund
    refunded_amount DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (refunded_amount >= 0),
    refund_reason TEXT,
    
    -- Timestamps
    processed_at TIMESTAMP,
    refunded_at TIMESTAMP,
    
    -- Metadata (for additional Stripe data)
    metadata JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for payments
CREATE INDEX idx_payments_payment_id ON payments(payment_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_stripe_payment_intent_id ON payments(stripe_payment_intent_id);
CREATE INDEX idx_payments_stripe_charge_id ON payments(stripe_charge_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at DESC);

-- Create payouts table (for seller payouts)
CREATE TABLE IF NOT EXISTS payouts (
    id BIGSERIAL PRIMARY KEY,
    payout_id VARCHAR(100) UNIQUE NOT NULL,
    
    -- Relations
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    seller_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    payment_id BIGINT REFERENCES payments(id) ON DELETE RESTRICT,
    
    -- Stripe Connect details
    stripe_transfer_id VARCHAR(255),
    stripe_account_id VARCHAR(255),
    
    -- Amount
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (
        status IN (
            'pending',
            'processing',
            'paid',
            'failed',
            'reversed'
        )
    ),
    
    -- Failure info
    failure_reason TEXT,
    
    -- Timestamps
    processed_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for payouts
CREATE INDEX idx_payouts_payout_id ON payouts(payout_id);
CREATE INDEX idx_payouts_order_id ON payouts(order_id);
CREATE INDEX idx_payouts_seller_id ON payouts(seller_id);
CREATE INDEX idx_payouts_payment_id ON payouts(payment_id);
CREATE INDEX idx_payouts_stripe_transfer_id ON payouts(stripe_transfer_id);
CREATE INDEX idx_payouts_status ON payouts(status);
CREATE INDEX idx_payouts_created_at ON payouts(created_at DESC);

-- Triggers for updated_at
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payouts_updated_at BEFORE UPDATE ON payouts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to generate payment ID
CREATE OR REPLACE FUNCTION generate_payment_id()
RETURNS TEXT AS $$
DECLARE
    new_id TEXT;
    year_part TEXT;
    seq_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    
    SELECT LPAD(
        (COUNT(*) + 1)::TEXT, 
        8, 
        '0'
    ) INTO seq_part
    FROM payments
    WHERE EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM NOW());
    
    -- Format: PAY-2026-00000001
    new_id := 'PAY-' || year_part || '-' || seq_part;
    
    RETURN new_id;
END;
$$ LANGUAGE plpgsql;

-- Function to generate payout ID
CREATE OR REPLACE FUNCTION generate_payout_id()
RETURNS TEXT AS $$
DECLARE
    new_id TEXT;
    year_part TEXT;
    seq_part TEXT;
BEGIN
    year_part := TO_CHAR(NOW(), 'YYYY');
    
    SELECT LPAD(
        (COUNT(*) + 1)::TEXT, 
        8, 
        '0'
    ) INTO seq_part
    FROM payouts
    WHERE EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM NOW());
    
    -- Format: PAYOUT-2026-00000001
    new_id := 'PAYOUT-' || year_part || '-' || seq_part;
    
    RETURN new_id;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-generate payment_id
CREATE OR REPLACE FUNCTION set_payment_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.payment_id IS NULL OR NEW.payment_id = '' THEN
        NEW.payment_id := generate_payment_id();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_payment_id_trigger
    BEFORE INSERT ON payments
    FOR EACH ROW
    EXECUTE FUNCTION set_payment_id();

-- Trigger to auto-generate payout_id
CREATE OR REPLACE FUNCTION set_payout_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.payout_id IS NULL OR NEW.payout_id = '' THEN
        NEW.payout_id := generate_payout_id();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_payout_id_trigger
    BEFORE INSERT ON payouts
    FOR EACH ROW
    EXECUTE FUNCTION set_payout_id();

-- Comments for documentation
COMMENT ON TABLE payments IS 'Payment transactions processed via Stripe';
COMMENT ON TABLE payouts IS 'Seller payouts via Stripe Connect';
COMMENT ON COLUMN payments.payment_id IS 'Unique payment identifier (e.g., PAY-2026-00000001)';
COMMENT ON COLUMN payments.stripe_payment_intent_id IS 'Stripe PaymentIntent ID';
COMMENT ON COLUMN payments.stripe_charge_id IS 'Stripe Charge ID (after payment succeeds)';
COMMENT ON COLUMN payouts.payout_id IS 'Unique payout identifier (e.g., PAYOUT-2026-00000001)';
COMMENT ON COLUMN payouts.stripe_transfer_id IS 'Stripe Transfer ID (Stripe Connect)';
COMMENT ON COLUMN payouts.stripe_account_id IS 'Seller Stripe Connect account ID';
