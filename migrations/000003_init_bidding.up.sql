-- Create bids table
CREATE TABLE IF NOT EXISTS bids (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    size_id BIGINT NOT NULL REFERENCES sizes(id) ON DELETE CASCADE,
    price DECIMAL(10, 2) NOT NULL CHECK (price > 0),
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'matched', 'cancelled', 'expired')),
    expires_at TIMESTAMP,
    matched_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create asks table
CREATE TABLE IF NOT EXISTS asks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    size_id BIGINT NOT NULL REFERENCES sizes(id) ON DELETE CASCADE,
    price DECIMAL(10, 2) NOT NULL CHECK (price > 0),
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'matched', 'cancelled', 'expired')),
    expires_at TIMESTAMP,
    matched_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create matches table (when bid meets ask)
CREATE TABLE IF NOT EXISTS matches (
    id BIGSERIAL PRIMARY KEY,
    bid_id BIGINT NOT NULL REFERENCES bids(id) ON DELETE CASCADE,
    ask_id BIGINT NOT NULL REFERENCES asks(id) ON DELETE CASCADE,
    buyer_id BIGINT NOT NULL REFERENCES users(id),
    seller_id BIGINT NOT NULL REFERENCES users(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    size_id BIGINT NOT NULL REFERENCES sizes(id),
    price DECIMAL(10, 2) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed')),
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for bids
CREATE INDEX idx_bids_user_id ON bids(user_id);
CREATE INDEX idx_bids_product_size ON bids(product_id, size_id);
CREATE INDEX idx_bids_status ON bids(status) WHERE status = 'active';
CREATE INDEX idx_bids_price ON bids(product_id, size_id, price DESC) WHERE status = 'active';
CREATE INDEX idx_bids_created_at ON bids(created_at DESC);

-- Create indexes for asks
CREATE INDEX idx_asks_user_id ON asks(user_id);
CREATE INDEX idx_asks_product_size ON asks(product_id, size_id);
CREATE INDEX idx_asks_status ON asks(status) WHERE status = 'active';
CREATE INDEX idx_asks_price ON asks(product_id, size_id, price ASC) WHERE status = 'active';
CREATE INDEX idx_asks_created_at ON asks(created_at DESC);

-- Create indexes for matches
CREATE INDEX idx_matches_bid_id ON matches(bid_id);
CREATE INDEX idx_matches_ask_id ON matches(ask_id);
CREATE INDEX idx_matches_buyer_id ON matches(buyer_id);
CREATE INDEX idx_matches_seller_id ON matches(seller_id);
CREATE INDEX idx_matches_product_size ON matches(product_id, size_id);
CREATE INDEX idx_matches_status ON matches(status);
CREATE INDEX idx_matches_created_at ON matches(created_at DESC);

-- Create triggers for updated_at
CREATE TRIGGER update_bids_updated_at BEFORE UPDATE ON bids
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_asks_updated_at BEFORE UPDATE ON asks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE bids IS 'Buyer bids - offers to buy at specific price';
COMMENT ON TABLE asks IS 'Seller asks - offers to sell at specific price';
COMMENT ON TABLE matches IS 'Matched bids and asks that resulted in transactions';
COMMENT ON COLUMN bids.status IS 'active: waiting for match, matched: completed, cancelled: user cancelled, expired: time expired';
COMMENT ON COLUMN asks.status IS 'active: waiting for match, matched: completed, cancelled: user cancelled, expired: time expired';
