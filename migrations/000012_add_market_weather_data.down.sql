-- Rollback: Remove market_data and weather_events tables

BEGIN;

-- Drop weather_events indexes and table
DROP INDEX IF EXISTS idx_weather_events_commodities;
DROP INDEX IF EXISTS idx_weather_events_start;
DROP INDEX IF EXISTS idx_weather_events_type;
DROP INDEX IF EXISTS idx_weather_events_region;

DROP TABLE IF EXISTS weather_events;

-- Drop market_data indexes and table
DROP INDEX IF EXISTS idx_market_data_source;
DROP INDEX IF EXISTS idx_market_data_recorded;
DROP INDEX IF EXISTS idx_market_data_commodity;

DROP TABLE IF EXISTS market_data;

COMMIT;
