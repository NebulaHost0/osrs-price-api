# Fix Cloudflare SSL Error

## Problem

```
curl: (35) SSL routines:ST_CONNECT:tlsv1 alert protocol version
```

This means:
- ✅ Server is reachable
- ❌ SSL/TLS handshake failing
- **Cause**: Your Go API serves HTTP, but Cloudflare expects HTTPS on port 8443

## ✅ Solution 1: Use Cloudflare Flexible SSL (Easiest)

Let Cloudflare handle SSL, keep your API on HTTP.

### Setup:

1. **Change your Portainer port back to 80 or 8080**:
```yaml
ports:
  - "80:8080"  # OR "8080:8080"
```

2. **Cloudflare** → **SSL/TLS** → Overview
   - Mode: **Flexible** ⚠️
   - This means: Browser → Cloudflare (HTTPS), Cloudflare → Your Server (HTTP)

3. **Cloudflare** → **DNS**
   - Type: A
   - Name: api
   - Target: Your server IP
   - **Proxy**: 🟠 **Proxied**

4. **Cloudflare** → **Rules** → **Page Rules** (or Configuration Rules)
   - Create rule for `api.grandexchange.gg/*`
   - Setting: **Forwarding URL**
   - Status: 301
   - Destination: Remove port from URL OR
   
   **Better: Create an origin rule**
   - Go to **Rules** → **Origin Rules**
   - Set destination port to `80` or `8080`

5. **Update Frontend**:
```env
NEXT_PUBLIC_GO_API_URL=https://api.grandexchange.gg
```

**Result**: Access API at `https://api.grandexchange.gg/health` (no port!)

### Flexible SSL Mode:
```
Browser ----HTTPS----> Cloudflare ----HTTP----> Your Server (port 80/8080)
```

---

## ✅ Solution 2: Standard HTTPS Port 443 (Recommended)

Use the standard HTTPS port and let Cloudflare proxy it.

### Setup:

1. **Update Portainer**:
```yaml
ports:
  - "443:8080"
```

2. **Cloudflare DNS**:
   - Type: A
   - Name: api
   - Proxy: 🟠 Proxied

3. **Cloudflare SSL/TLS**:
   - Mode: **Flexible**

4. **Frontend**:
```env
NEXT_PUBLIC_GO_API_URL=https://api.grandexchange.gg
```

5. **Test**:
```bash
curl https://api.grandexchange.gg/health
```

---

## ✅ Solution 3: Add Real SSL to Go API (Most Secure)

Add actual SSL/TLS to your Go API for end-to-end encryption.

### Step 1: Get Cloudflare Origin Certificate

1. **Cloudflare** → **SSL/TLS** → **Origin Server**
2. **Create Certificate**
3. **Save**:
   - Certificate → `origin-cert.pem`
   - Private Key → `origin-key.pem`

### Step 2: Update Go API to Use TLS

Create `internal/api/tls.go`:
```go
package api

import (
    "github.com/gin-gonic/gin"
)

func RunWithTLS(router *gin.Engine, port, certFile, keyFile string) error {
    return router.RunTLS(":"+port, certFile, keyFile)
}
```

Update `main.go`:
```go
// Check if SSL certificates exist
certFile := os.Getenv("SSL_CERT_FILE")
keyFile := os.Getenv("SSL_KEY_FILE")

if certFile != "" && keyFile != "" {
    log.Printf("Starting HTTPS server on port %s", port)
    if err := router.RunTLS(":"+port, certFile, keyFile); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
} else {
    log.Printf("Starting HTTP server on port %s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

### Step 3: Mount Certificates in Docker

```yaml
services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    ports:
      - "8443:8443"
    volumes:
      - ./origin-cert.pem:/app/certs/cert.pem:ro
      - ./origin-key.pem:/app/certs/key.pem:ro
    environment:
      DATABASE_URL: postgresql://...
      PORT: 8443
      SSL_CERT_FILE: /app/certs/cert.pem
      SSL_KEY_FILE: /app/certs/key.pem
```

### Step 4: Cloudflare SSL Mode

- Mode: **Full (strict)** ✅
- Now you have end-to-end encryption

---

## What I Recommend

**For Your Case**: Use **Solution 2** (Standard Port 443)

### Why:
✅ Simplest setup  
✅ No custom port in URLs  
✅ Cloudflare handles all SSL  
✅ No certificate management  
✅ Professional appearance  
✅ Works immediately  

### To Implement:

**1. Update Portainer**:
```yaml
ports:
  - "443:8080"
```

**2. Update Frontend**:
```env
NEXT_PUBLIC_GO_API_URL=https://api.grandexchange.gg
```

**3. Cloudflare**:
- SSL Mode: Flexible
- Proxy: On (orange cloud)

**4. Test**:
```bash
curl https://api.grandexchange.gg/health
```

Done! Clean, simple, and it works. 🎉