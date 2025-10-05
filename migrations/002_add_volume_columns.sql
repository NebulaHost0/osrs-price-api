-- Add volume tracking columns to price_history table
ALTER TABLE price_history 
ADD COLUMN IF NOT EXISTS high_volume BIGINT DEFAULT 0,
ADD COLUMN IF NOT EXISTS low_volume BIGINT DEFAULT 0;

-- Create index for volume queries
CREATE INDEX IF NOT EXISTS idx_price_history_volume ON price_history (item_id, timestamp, high_volume, low_volume);

-- Add comment to explain the columns
COMMENT ON COLUMN price_history.high_volume IS 'Trading volume at high price from OSRS Wiki API';
COMMENT ON COLUMN price_history.low_volume IS 'Trading volume at low price from OSRS Wiki API';