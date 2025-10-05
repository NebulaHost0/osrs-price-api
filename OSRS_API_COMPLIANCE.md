# OSRS Wiki API Compliance

This document explains how our API complies with the [OSRS Wiki Real-time Prices API guidelines](https://oldschool.runescape.wiki/w/RuneScape:Real-time_Prices).

## Compliance Summary

✅ **User-Agent**: Set to "grandexchange.gg - OSRS GE Price Tracker"  
✅ **Rate Limiting**: Max 1 request per 5 minutes to the Wiki API  
✅ **Bulk Endpoints**: Uses `/latest` bulk endpoint, never individual item requests  
✅ **Efficient Design**: Caches prices and serves from database  
✅ **Respectful Usage**: Minimal impact on Wiki API servers

## Implementation Details

### 1. User-Agent Header

**Required by OSRS Wiki**: All requests must include a descriptive User-Agent.

**Our Implementation**:
```go
// internal/osrs/client.go
const userAgent = "grandexchange.gg - OSRS GE Price Tracker"

req.Header.Set("User-Agent", userAgent)
```

**Frontend** (for item mapping and time series):
```typescript
// osrs-ge-tracker/src/lib/osrs-api.ts
const USER_AGENT = 'grandexchange.gg - OSRS GE Price Tracker';
```

### 2. Rate Limiting

**OSRS Wiki Guideline**: "Don't spam the API"

**Our Implementation**:
- Background worker fetches prices **once every 5 minutes** (matches Wiki update frequency)
- Prices are stored in PostgreSQL database
- API serves from database, not Wiki API
- Frontend never makes bulk price requests to Wiki API

```go
// main.go
priceFetcher := worker.NewPriceFetcher(osrsClient, repo, 5*time.Minute)
```

**Request Frequency**:
- Wiki API: ~12 requests per hour (one every 5 minutes)
- Our API: Unlimited (serves from database)

### 3. Bulk Endpoint Usage

**OSRS Wiki Guideline**: "Use the bulk endpoint /latest, don't make 3700 individual requests"

**Our Implementation**:
```go
// internal/osrs/client.go

// ✅ CORRECT: Uses bulk endpoint
func (c *Client) GetLatestPrices() (map[string]models.ItemPrice, error) {
    req, err := http.NewRequest("GET", "https://prices.runescape.wiki/api/v1/osrs/latest", nil)
    // Fetches ALL items in one request
}

// ✅ CORRECT: Uses bulk endpoint internally
func (c *Client) GetItemPrice(itemID string) (*models.ItemPrice, error) {
    prices, err := c.GetLatestPrices() // Uses bulk endpoint
    return &prices[itemID], nil        // Extracts single item
}
```

We **NEVER** make individual requests like `/latest?id=123` in a loop.

### 4. Caching Strategy

**Our Approach**:
1. **Background Worker**: Fetches all prices every 5 minutes
2. **Database Storage**: Stores historical data in PostgreSQL
3. **In-Memory Cache**: Caches current prices for instant access
4. **API Layer**: Serves from cache/database, not Wiki API

```
OSRS Wiki API → (5 min intervals) → Our Go Worker → PostgreSQL
                                                      ↓
Frontend ← Our API ← Cache/Database ← Historical queries
```

**Result**: Our users get instant responses without impacting Wiki API.

### 5. What We Fetch from Wiki API

#### Backend (Go API)
- **Current Prices**: `/latest` - Once per 5 minutes
- **That's it!** All other data served from our database

#### Frontend (Next.js)
- **Item Mapping**: `/mapping` - Once at page load (cached by browser)
- **Time Series** (charts): `/timeseries` - On-demand when user views charts
  - Uses appropriate timesteps (5m, 1h, 6h, 24h)
  - Fetched directly from Wiki (historical data not in our DB yet)

### 6. Best Practices We Follow

✅ **Descriptive User-Agent**: Identifies our service clearly  
✅ **Minimal Requests**: Only 12 requests/hour for price updates  
✅ **Bulk Operations**: Always use bulk endpoints  
✅ **Proper Timesteps**: Use appropriate intervals for time series  
✅ **Caching**: Store and serve from our own database  
✅ **Error Handling**: Gracefully handle API failures  
✅ **Documentation**: Clear documentation of API usage  

### 7. Request Breakdown

**Per Hour from Our Backend**:
- Current prices (`/latest`): 12 requests (every 5 minutes)

**Per User Session from Frontend**:
- Item mapping (`/mapping`): 1 request (cached)
- Time series (`/timeseries`): ~1-3 requests (only when viewing charts)

**Total Impact**: Very low, well within acceptable limits.

### 8. Future Improvements

Potential optimizations to reduce Wiki API usage even further:

1. **Store Time Series Data**: Cache historical data in our database
2. **Item Mapping Cache**: Store item mapping in database, update weekly
3. **CDN Layer**: Add CDN for static data like item mapping

## Verification

To verify our implementation:

```bash
# Check worker is running with 5-minute interval
grep "5.*Minute" main.go

# Verify User-Agent is set
grep "grandexchange.gg" internal/osrs/client.go

# Confirm bulk endpoint usage
grep "wikiAPIURL" internal/osrs/client.go
```

## Contact

If the OSRS Wiki team has any concerns about our API usage, please contact us at:
- Website: grandexchange.gg
- Issue tracker: GitHub repository

## Acknowledgments

Thank you to the OSRS Wiki team for providing this excellent API and clear guidelines. We strive to be a respectful and efficient user of your services.

## References

- [OSRS Wiki Real-time Prices API](https://oldschool.runescape.wiki/w/RuneScape:Real-time_Prices)
- [API Documentation](https://oldschool.runescape.wiki/w/RuneScape:Real-time_Prices/API)