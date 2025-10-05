-- Hourly aggregated price data (for 1 week - 3 months charts)
CREATE TABLE IF NOT EXISTS price_history_hourly (
    id BIGSERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL,
    avg_high BIGINT NOT NULL,
    avg_low BIGINT NOT NULL,
    max_high BIGINT NOT NULL,
    min_low BIGINT NOT NULL,
    opening_high BIGINT,
    opening_low BIGINT,
    closing_high BIGINT,
    closing_low BIGINT,
    total_high_volume BIGINT DEFAULT 0,
    total_low_volume BIGINT DEFAULT 0,
    data_points INTEGER NOT NULL,
    hour_timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_hourly_item_time ON price_history_hourly (item_id, hour_timestamp);

-- Daily aggregated price data (for 6 months - 5 years charts)
CREATE TABLE IF NOT EXISTS price_history_daily (
    id BIGSERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL,
    avg_high BIGINT NOT NULL,
    avg_low BIGINT NOT NULL,
    max_high BIGINT NOT NULL,
    min_low BIGINT NOT NULL,
    opening_high BIGINT,
    opening_low BIGINT,
    closing_high BIGINT,
    closing_low BIGINT,
    total_high_volume BIGINT DEFAULT 0,
    total_low_volume BIGINT DEFAULT 0,
    volatility FLOAT,
    data_points INTEGER NOT NULL,
    day_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(item_id, day_date)
);

CREATE INDEX idx_daily_item_date ON price_history_daily (item_id, day_date);

-- Comments
COMMENT ON TABLE price_history_hourly IS 'Hourly aggregates of price data for medium-term charts';
COMMENT ON TABLE price_history_daily IS 'Daily aggregates of price data for long-term charts';

COMMENT ON COLUMN price_history_hourly.data_points IS 'Number of 5-minute samples aggregated into this hour';
COMMENT ON COLUMN price_history_daily.data_points IS 'Number of hourly samples aggregated into this day';
COMMENT ON COLUMN price_history_daily.volatility IS 'Price volatility measure: (max_high - min_low) / avg_high';