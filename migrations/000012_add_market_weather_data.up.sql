-- Migration: Add market_data and weather_events tables
-- For tracking commodity prices and weather impacts
-- Date: 2026-02-02

BEGIN;

-- Market Data (pricing, volume)
CREATE TABLE market_data (
    id BIGSERIAL PRIMARY KEY,
    commodity_type VARCHAR(100) NOT NULL,
    
    -- Pricing
    spot_price DECIMAL(10,4) NOT NULL,        -- Current market price
    futures_price DECIMAL(10,4),              -- 3-month futures
    
    -- Volume
    daily_volume DECIMAL(15,2),               -- Trading volume
    open_interest INTEGER,                    -- Open futures contracts
    
    -- Price Range
    daily_high DECIMAL(10,4),
    daily_low DECIMAL(10,4),
    
    -- Change
    price_change DECIMAL(10,4),
    price_change_percent DECIMAL(5,2),
    
    -- Data source
    data_source VARCHAR(100),                 -- 'CME', 'ICE', 'internal'
    
    recorded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    
    CHECK (spot_price > 0),
    CHECK (daily_volume IS NULL OR daily_volume >= 0)
);

-- Weather Events (impacts agricultural prices)
CREATE TABLE weather_events (
    id BIGSERIAL PRIMARY KEY,
    region VARCHAR(255) NOT NULL,
    event_type VARCHAR(100) NOT NULL,         -- 'drought', 'flood', 'frost', 'heatwave'
    severity VARCHAR(50),                     -- 'low', 'medium', 'high', 'severe'
    
    -- Impact
    affected_commodities TEXT[],              -- ['wheat', 'corn']
    estimated_impact_percent DECIMAL(5,2),    -- -15% (yield reduction)
    
    -- Dates
    event_start TIMESTAMP NOT NULL,
    event_end TIMESTAMP,
    
    -- Source
    data_source VARCHAR(100),                 -- 'NOAA', 'USDA', 'local'
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_market_data_commodity ON market_data(commodity_type);
CREATE INDEX idx_market_data_recorded ON market_data(recorded_at);
CREATE INDEX idx_market_data_source ON market_data(data_source);

CREATE INDEX idx_weather_events_region ON weather_events(region);
CREATE INDEX idx_weather_events_type ON weather_events(event_type);
CREATE INDEX idx_weather_events_start ON weather_events(event_start);
CREATE INDEX idx_weather_events_commodities ON weather_events USING gin(affected_commodities);

-- Comments
COMMENT ON TABLE market_data IS 'Historical and real-time commodity prices';
COMMENT ON TABLE weather_events IS 'Weather events that impact agricultural production';
COMMENT ON COLUMN weather_events.estimated_impact_percent IS 'Expected yield reduction (negative) or increase (positive)';

COMMIT;
