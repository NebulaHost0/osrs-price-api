-- Configuration table for cleanup settings
CREATE TABLE IF NOT EXISTS cleanup_config (
    id SERIAL PRIMARY KEY,
    retention_days INTEGER NOT NULL DEFAULT 7,
    last_cleanup TIMESTAMP,
    records_deleted_last_run BIGINT DEFAULT 0,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default configuration
INSERT INTO cleanup_config (retention_days, enabled) 
VALUES (7, true);

-- Cleanup logs table to track cleanup operations
CREATE TABLE IF NOT EXISTS cleanup_logs (
    id BIGSERIAL PRIMARY KEY,
    records_deleted BIGINT NOT NULL,
    cutoff_date TIMESTAMP NOT NULL,
    duration_ms INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_cleanup_logs_created (created_at DESC)
);

-- Add comment
COMMENT ON TABLE cleanup_config IS 'Configuration for automatic database cleanup';
COMMENT ON TABLE cleanup_logs IS 'Audit log of cleanup operations';