# Cost Optimization Guide

## Problem: Database Growth

Without maintenance, the database grows at:
- **1 million+ records per day**
- **32 million records per month**
- **388 million records per year**
- **~37 GB per year**

This costs **$20-40/month** for PostgreSQL hosting.

## Solution: Automated Cleanup

### ✅ Implemented

1. **Automatic Data Retention**
   - Keeps last 7 days of detailed data (5-min intervals)
   - Automatically deletes older data
   - Runs daily at 3 AM

2. **Cleanup Worker**
   - Monitors database size
   - Logs all cleanup operations
   - Provides statistics on dashboard

3. **Configurable Retention**
   - Default: 7 days
   - Adjustable via database config
   - Can be tuned based on needs

### Cost Savings

**Before**: 388M records/year = ~37 GB = $20-40/month
**After**: 7.5M records (7 days) = ~720 MB = $5-10/month

**Savings: 75-80% cost reduction**

## How It Works

```
Day 1-7: Keep all 5-minute data (detailed charts, recent analysis)
Day 8+:  Automatically deleted (saves space and costs)
```

### Data Flow

```
OSRS Wiki API → Go API → PostgreSQL
                           ↓
                    [7 days retention]
                           ↓
                    [Auto cleanup at Day 8]
```

## Monitoring

### Database Stats Endpoint

```bash
# Check database health
curl http://localhost:8080/api/v1/admin/stats
```

Returns:
```json
{
  "total_records": 7500000,
  "oldest_record": "2025-09-28",
  "newest_record": "2025-10-05",
  "estimated_size": "720 MB",
  "retention_days": 7
}
```

### Cleanup Logs

All cleanup operations are logged:
- Records deleted
- Execution time
- Any errors
- Database size before/after

## Configuration

### Change Retention Period

```sql
-- Keep 14 days instead of 7
UPDATE cleanup_config SET retention_days = 14;

-- Keep only 3 days (minimal)
UPDATE cleanup_config SET retention_days = 3;

-- Disable cleanup (not recommended)
UPDATE cleanup_config SET enabled = false;
```

### Manual Cleanup

```bash
# Force cleanup right now
curl -X POST http://localhost:8080/api/v1/admin/cleanup
```

## Recommendations

### For Different Use Cases

**Day Trading / Real-time Analysis**:
- Retention: 3-7 days
- Cost: $5-10/month
- Use case: Short-term price movements

**Long-term Analysis**:
- Retention: 30 days
- Cost: $15-20/month
- Use case: Monthly trend analysis

**Historical Research**:
- Retention: 90 days
- Cost: $30-40/month
- Consider implementing data aggregation

## Future Enhancements

### Phase 2: Data Aggregation (Not Yet Implemented)

Instead of deleting old data, aggregate it:

```
Days 1-7:   5-minute intervals (raw)
Days 8-30:  1-hour aggregates
Days 31+:   24-hour aggregates
```

Benefits:
- Keep historical trends
- Minimal storage cost
- Enable long-term analysis

### Phase 3: Cold Storage

Archive old data to S3/Backblaze:
- $0.005/GB/month
- On-demand retrieval
- Unlimited history

## Best Practices

1. **Monitor Regularly**
   - Check database size weekly
   - Review cleanup logs
   - Adjust retention as needed

2. **Set Alerts**
   - Alert if database > 2 GB
   - Alert if cleanup fails
   - Alert if growth rate spikes

3. **Plan for Growth**
   - More items tracked = more records
   - Consider aggregation for long-term
   - Budget for ~$10-20/month for production

## Cost Comparison

### Railway PostgreSQL Pricing

| Storage | Cost/Month | Records Supported |
|---------|------------|-------------------|
| 512 MB  | Free       | ~5.3M (3-4 days)  |
| 1 GB    | $5         | ~10M (7 days)     |
| 5 GB    | $10        | ~52M (35 days)    |
| 10 GB   | $15        | ~104M (70 days)   |

### Other Providers

**Supabase**:
- Free: 500 MB
- Pro: $25/month for 8 GB

**Neon**:
- Free: 3 GB
- Pro: $19/month for 10 GB

**DigitalOcean**:
- $15/month for 10 GB

## Summary

✅ **Automatic cleanup enabled**
✅ **7-day retention by default**
✅ **75-80% cost savings**
✅ **Configurable and monitored**
✅ **Production-ready**

Your database will now automatically maintain a healthy size and keep costs low! 🎉