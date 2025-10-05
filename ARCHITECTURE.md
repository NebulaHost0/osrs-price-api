# Architecture Overview

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         OSRS Price API                          │
└─────────────────────────────────────────────────────────────────┘

┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   Client     │────────▶│  Gin Router  │────────▶│   Handler    │
│  (HTTP)      │         │              │         │              │
└──────────────┘         └──────────────┘         └───────┬──────┘
                                                           │
                         ┌─────────────────────────────────┼──────┐
                         │                                 │      │
                         ▼                                 ▼      ▼
                  ┌──────────┐                      ┌─────────┐  │
                  │  Cache   │                      │  OSRS   │  │
                  │(Memory)  │                      │ Client  │  │
                  └──────────┘                      └────┬────┘  │
                                                         │       │
                  ┌──────────────┐                      │       │
                  │  PostgreSQL  │◀─────────────────────┘       │
                  │   Database   │                              │
                  └──────┬───────┘                              │
                         ▲                                      │
                         │                                      │
                  ┌──────┴──────┐                               │
                  │ Repository  │◀──────────────────────────────┘
                  └─────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    Background Worker                            │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐   │
│  │ Price        │────▶│  OSRS Wiki   │────▶│  Database    │   │
│  │ Fetcher      │     │     API      │     │  (Batch)     │   │
│  │ (5 min)      │     │              │     │              │   │
│  └──────────────┘     └──────────────┘     └──────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## Component Description

### 1. HTTP Layer
- **Gin Router**: High-performance HTTP router
- **Handler**: Request processing and response formatting
- **Routes**: API endpoint definitions

### 2. Business Logic
- **OSRS Client**: Communicates with OSRS Wiki API
- **Repository**: Database operations and queries
- **Cache**: In-memory caching for current prices (5-minute TTL)

### 3. Data Storage
- **PostgreSQL**: Persistent storage for historical price data
- **In-Memory Cache**: Fast access to current prices

### 4. Background Worker
- **Price Fetcher**: Periodically fetches prices every 5 minutes
- Runs asynchronously from the main API
- Stores data in batches for performance

## Data Flow

### Real-time Price Request
```
1. Client → GET /api/v1/prices/4151
2. Handler checks cache
3. If cached: Return immediately
4. If not: Fetch from OSRS Wiki API
5. Store in cache (5 min TTL)
6. Return to client
```

### Historical Data Request
```
1. Client → GET /api/v1/history/4151?hours=24
2. Handler validates parameters
3. Repository queries PostgreSQL
4. Apply time range filter
5. Return aggregated data
```

### Background Price Collection
```
1. Timer triggers every 5 minutes
2. Fetch all prices from OSRS Wiki
3. Batch insert into PostgreSQL (100 records at a time)
4. Log success/failure
5. Wait for next interval
```

## Database Schema

### price_history Table
```sql
CREATE TABLE price_history (
    id          SERIAL PRIMARY KEY,
    item_id     INTEGER NOT NULL,
    high        BIGINT,
    high_time   BIGINT,
    low         BIGINT,
    low_time    BIGINT,
    timestamp   TIMESTAMP NOT NULL,
    created_at  TIMESTAMP,
    
    INDEX idx_item_timestamp (item_id, timestamp)
);
```

**Indexes:**
- `idx_item_timestamp`: Composite index for fast queries by item and time range

## API Endpoints

### Current Prices
- `GET /api/v1/prices` - All current prices
- `GET /api/v1/prices/:id` - Specific item price

### Historical Data
- `GET /api/v1/history/:id?hours=24` - Price history
- `GET /api/v1/change/:id?hours=24` - Price change analysis
- `GET /api/v1/stats/:id?hours=168` - Statistical analysis

### Market Analysis
- `GET /api/v1/gainers?limit=10&hours=24` - Top price increases
- `GET /api/v1/losers?limit=10&hours=24` - Top price decreases

### Admin
- `POST /api/v1/cache/clear` - Clear cache
- `GET /health` - Health check

## Performance Optimizations

1. **Caching Strategy**
   - Current prices cached for 5 minutes
   - Reduces API calls to OSRS Wiki
   - Faster response times

2. **Database Indexing**
   - Composite index on (item_id, timestamp)
   - Optimizes historical queries
   - Faster aggregations

3. **Batch Processing**
   - Insert prices in batches of 100
   - Reduces database round trips
   - Better throughput

4. **Connection Pooling**
   - Max 100 database connections
   - 10 idle connections
   - 1-hour connection lifetime

## Error Handling

- Graceful degradation when OSRS Wiki API is unavailable
- Database connection retry logic
- Comprehensive error messages
- Logging for debugging

## Scalability Considerations

1. **Horizontal Scaling**: Stateless API design allows multiple instances
2. **Database Partitioning**: Table can be partitioned by timestamp
3. **Read Replicas**: Historical queries can use read replicas
4. **Caching**: Can be replaced with Redis for distributed caching

## Security

- No authentication required (public data)
- Rate limiting recommended for production
- SQL injection prevention via parameterized queries
- CORS configuration needed for web clients