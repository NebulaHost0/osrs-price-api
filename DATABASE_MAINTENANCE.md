# Database Maintenance Strategy

## Current Growth Rate

- **Items tracked**: ~3,700
- **Update frequency**: Every 5 minutes (12 times/hour)
- **Records per day**: 3,700 × 12 × 24 = **1,065,600 records/day**
- **Records per month**: ~32 million
- **Records per year**: ~388 million

## Storage Estimates

Each record is approximately:
- `id` (BIGINT): 8 bytes
- `item_id` (INT): 4 bytes
- `high`, `low`, `high_time`, `low_time`, `high_volume`, `low_volume` (6 × BIGINT): 48 bytes
- `timestamp`, `created_at` (2 × TIMESTAMP): 16 bytes
- Indexes overhead: ~20 bytes
- **Total per record**: ~96 bytes

### Projected Growth:
- **1 month**: 32M records × 96 bytes = ~3 GB
- **6 months**: ~18 GB
- **1 year**: ~37 GB
- **2 years**: ~74 GB

## ✅ IMPLEMENTED: Tiered Retention Strategy

### Tier 1: High-Resolution Data (Recent)
- **Period**: Last 7 days
- **Resolution**: 5-minute intervals (raw data)
- **Purpose**: Detailed short-term analysis, 1D-1W charts
- **Records**: ~7.5 million
- **Table**: `price_history`

### Tier 2: Medium-Resolution Data (1-3 months)
- **Period**: 8-90 days ago
- **Resolution**: 1-hour aggregates
- **Purpose**: Medium-term trends, 1M-3M charts
- **Records**: ~620K (12× reduction from raw)
- **Table**: `price_history_hourly`

### Tier 3: Low-Resolution Data (Historical)
- **Period**: 91+ days to 5 years
- **Resolution**: 24-hour aggregates
- **Purpose**: Long-term trends, 6M-5Y charts
- **Records**: ~6.75M for 5 years (24× reduction from hourly)
- **Table**: `price_history_daily`

### Result:
Instead of **388M records/year**, we keep **~15M records/year** (96% reduction)

### Storage Breakdown:
- **Raw (7 days)**: 7.5M × 96 bytes = 720 MB
- **Hourly (90 days)**: 620K × 120 bytes = 74 MB
- **Daily (5 years)**: 6.75M × 140 bytes = 945 MB
- **Total for 5 years**: ~1.7 GB (instead of 185 GB!)

### Cost Impact:
- **Without aggregation**: $150+/month
- **With aggregation**: $5-10/month
- **Savings**: 95%+ cost reduction

## Implementation

### 1. Data Aggregation Tables

```sql
-- Hourly aggregates
CREATE TABLE price_history_hourly (
    id BIGSERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL,
    avg_high BIGINT,
    avg_low BIGINT,
    max_high BIGINT,
    min_low BIGINT,
    total_high_volume BIGINT,
    total_low_volume BIGINT,
    data_points INTEGER,
    hour_timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_hourly_item_time (item_id, hour_timestamp)
);

-- Daily aggregates
CREATE TABLE price_history_daily (
    id BIGSERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL,
    avg_high BIGINT,
    avg_low BIGINT,
    max_high BIGINT,
    min_low BIGINT,
    opening_high BIGINT,
    opening_low BIGINT,
    closing_high BIGINT,
    closing_low BIGINT,
    total_high_volume BIGINT,
    total_low_volume BIGINT,
    volatility FLOAT,
    data_points INTEGER,
    day_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_daily_item_date (item_id, day_date)
);
```

### 2. Automated Cleanup Worker

The system will automatically:
1. **Aggregate old data** (nightly)
2. **Delete raw data** after aggregation
3. **Compress historical data** (monthly)

### 3. Retention Policies

- **Raw 5-min data**: Keep 7 days
- **Hourly aggregates**: Keep 90 days
- **Daily aggregates**: Keep forever (minimal space)

## Cost Analysis

### Before Optimization:
- **1 year**: 388M records = ~37 GB
- **PostgreSQL cost** (Railway): ~$20-40/month for 40GB

### After Optimization:
- **1 year**: 18M records = ~1.7 GB
- **PostgreSQL cost** (Railway): ~$5-10/month for 5GB

**Savings**: 75-80% cost reduction

## Monitoring

Track these metrics:
- Total database size
- Records per table
- Oldest record date
- Daily growth rate
- Cleanup job success/failure

## Backup Strategy

1. **Daily backups** of last 7 days (high-res)
2. **Weekly backups** of aggregated data
3. **Monthly exports** to cold storage (S3/Backblaze)

## Future Optimizations

1. **Table Partitioning** (by date range)
2. **Compression** (pg_compress)
3. **Time-series DB** (TimescaleDB) for better performance
4. **Archive to S3** for old data retrieval