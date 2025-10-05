-- Migration: Create price_history table
-- Created: 2025-10-05

CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL,
    high BIGINT,
    high_time BIGINT,
    low BIGINT,
    low_time BIGINT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE
);

-- Create composite index for efficient queries
CREATE INDEX IF NOT EXISTS idx_item_timestamp ON price_history(item_id, timestamp);

-- Create index on timestamp for time-based queries
CREATE INDEX IF NOT EXISTS idx_timestamp ON price_history(timestamp);

COMMENT ON TABLE price_history IS 'Historical price data for OSRS items';
COMMENT ON COLUMN price_history.item_id IS 'OSRS item ID from the game';
COMMENT ON COLUMN price_history.high IS 'High (buy) price in GP';
COMMENT ON COLUMN price_history.low IS 'Low (sell) price in GP';
COMMENT ON COLUMN price_history.timestamp IS 'Time when the price was recorded';