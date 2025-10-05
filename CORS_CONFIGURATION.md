# CORS Configuration Guide

## Allowed Origins

The API allows CORS requests from:

### Development:
- `http://localhost:3000` ✅
- `http://localhost:3001` ✅
- `http://127.0.0.1:3000` ✅
- `https://localhost:3000` ✅

### Production:
- `https://grandexchange.gg` ✅
- `https://www.grandexchange.gg` ✅
- `http://grandexchange.gg` ✅
- `http://www.grandexchange.gg` ✅

### Additional Origins (via Environment Variable):
Set `ALLOWED_ORIGINS` to add more domains:

```env
ALLOWED_ORIGINS=https://staging.grandexchange.gg,https://dev.grandexchange.gg
```

## CORS Headers Set

The API returns these headers:

```
Access-Control-Allow-Origin: <origin>
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: POST, OPTIONS, GET, PUT, DELETE, PATCH
Access-Control-Allow-Headers: Content-Type, Authorization, ...
Access-Control-Max-Age: 3600
```

## Testing CORS

### Test preflight request:
```bash
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: GET" \
     -X OPTIONS \
     https://api.grandexchange.gg/api/v1/prices \
     -v
```

Should return:
```
HTTP/1.1 204 No Content
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: POST, OPTIONS, GET, PUT, DELETE, PATCH
```

### Test actual request:
```bash
curl -H "Origin: http://localhost:3000" \
     https://api.grandexchange.gg/health
```

Should return JSON with CORS headers.

## Browser Testing

Open browser console on `http://localhost:3000`:

```javascript
// Should work
fetch('https://api.grandexchange.gg/health')
  .then(r => r.json())
  .then(console.log)

// Should work
fetch('https://api.grandexchange.gg/api/v1/prices')
  .then(r => r.json())
  .then(console.log)
```

## Adding New Origins

To allow additional domains:

### Option 1: Update Code

Edit `internal/api/cors.go`:
```go
allowedOrigins := []string{
    "http://localhost:3000",
    // ... existing origins ...
    "https://your-new-domain.com", // Add here
}
```

### Option 2: Environment Variable

In Portainer or `.env`:
```env
ALLOWED_ORIGINS=https://staging.example.com,https://preview.example.com
```

## Troubleshooting

### CORS error in browser:
```
Access to fetch at 'https://api.grandexchange.gg/...' from origin 
'http://localhost:3000' has been blocked by CORS policy
```

**Check**:
1. Origin is in allowed list
2. API is serving HTTPS (if frontend uses HTTPS)
3. Preflight OPTIONS request succeeds
4. Response has correct CORS headers

### Testing from command line:
```bash
# With origin header
curl -H "Origin: http://localhost:3000" \
     -v https://api.grandexchange.gg/health 2>&1 | grep -i "access-control"

# Should show:
# < access-control-allow-origin: http://localhost:3000
```

## Security Notes

### Why not use `*` wildcard?

Using specific origins is more secure:
- ✅ Prevents unauthorized domains from accessing your API
- ✅ Allows credentials (cookies, auth headers)
- ✅ Better control over who can use your API

### Credentials Support

With `Access-Control-Allow-Credentials: true`:
- Frontend can send cookies
- Frontend can send Authorization headers
- More secure authentication

## Production Deployment

When deploying, ensure:
1. ✅ Production domain is in allowed origins
2. ✅ HTTPS is used (Cloudflare handles this)
3. ✅ CORS headers are properly set
4. ✅ Preflight caching works (max-age: 3600s)

## Current Configuration

**File**: `internal/api/cors.go`

**Default Allowed**:
- All localhost variants (dev)
- grandexchange.gg (all variants)
- Additional origins via `ALLOWED_ORIGINS` env var

**Methods Allowed**:
- GET, POST, PUT, DELETE, PATCH, OPTIONS

**Headers Allowed**:
- Content-Type, Authorization, and all standard headers

Perfect for your frontend! ✅