-- Migration: Add bulk_bids table
-- For group/corporate course purchases
-- Date: 2026-02-02

BEGIN;

CREATE TABLE bulk_bids (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    
    -- Organizer
    organizer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organizer_type VARCHAR(50) NOT NULL,      -- 'individual', 'company', 'university'
    company_name VARCHAR(255),
    
    -- Group details
    target_participants INTEGER NOT NULL,     -- Want X seats
    current_participants INTEGER DEFAULT 1,   -- Current sign-ups
    min_participants INTEGER,                 -- Minimum to proceed
    
    -- Pricing
    bid_price_per_seat DECIMAL(10,2) NOT NULL,
    total_bid_value DECIMAL(10,2) NOT NULL,
    
    -- Deadline
    deadline TIMESTAMP NOT NULL,
    auto_proceed BOOLEAN DEFAULT false,       -- Auto-purchase if target reached
    
    -- Status
    status VARCHAR(50) NOT NULL,              -- 'open', 'matched', 'cancelled', 'expired'
    matched_at TIMESTAMP,
    
    -- Participants (JSONB array)
    participants JSONB,
    -- Example: [{"user_id": 123, "email": "...", "joined_at": "..."}]
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CHECK (target_participants > 0),
    CHECK (current_participants <= target_participants),
    CHECK (total_bid_value = bid_price_per_seat * target_participants)
);

-- Indexes
CREATE INDEX idx_bulk_bids_course ON bulk_bids(course_id);
CREATE INDEX idx_bulk_bids_organizer ON bulk_bids(organizer_id);
CREATE INDEX idx_bulk_bids_status ON bulk_bids(status);
CREATE INDEX idx_bulk_bids_deadline ON bulk_bids(deadline);

-- Comments
COMMENT ON TABLE bulk_bids IS 'Group/corporate course purchases';
COMMENT ON COLUMN bulk_bids.participants IS 'JSONB array of participant info';

COMMIT;
