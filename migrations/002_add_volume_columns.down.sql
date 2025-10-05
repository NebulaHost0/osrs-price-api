-- Rollback volume tracking columns
DROP INDEX IF EXISTS idx_price_history_volume;

ALTER TABLE price_history 
DROP COLUMN IF EXISTS high_volume,
DROP COLUMN IF EXISTS low_volume;