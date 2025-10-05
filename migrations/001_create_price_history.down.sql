-- Rollback: Drop price_history table

DROP INDEX IF EXISTS idx_timestamp;
DROP INDEX IF EXISTS idx_item_timestamp;
DROP TABLE IF EXISTS price_history;