# Fix Cloudflare Connection NOW

## The Problem

- Cloudflare is proxying HTTPS traffic (port 443)
- Your Go API serves HTTP (not HTTPS)
- SSL handshake fails

## ✅ Quick Fix: Change Cloudflare SSL Mode

### Step 1: Change SSL Mode to Flexible

1. **Go to Cloudflare Dashboard**
2. **SSL/TLS** → **Overview**
3. **Change mode to: Flexible** ⚠️

This tells Cloudflare:
- Use HTTPS between browser and Cloudflare ✅
- Use HTTP between Cloudflare and your server ✅

### Step 2: Update Portainer Port

Change your container to use port **80** (HTTP):

```yaml
services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    restart: unless-stopped
    ports:
      - "80:8080"  # HTTP port
    environment:
      DATABASE_URL: postgresql://postgres:VGXGodgBrcUFmlKlUQDGFmLjBzUogxhJ@switchyard.proxy.rlwy.net:22411/railway
      PORT: 8080
```

Update the stack and redeploy.

### Step 3: Test

```bash
curl https://api.grandexchange.gg/health
```

Should work now!

---

## Alternative: Use Non-Standard Port with DNS Only

If you want to keep port 8888:

### 1. Cloudflare DNS - Turn OFF Proxy

1. **Cloudflare** → **DNS** → `api` record
2. Click the **orange cloud** → Changes to **gray cloud** (DNS only)
3. This bypasses Cloudflare proxy

### 2. Portainer - Keep Your Port

```yaml
ports:
  - "8888:8080"
```

### 3. Frontend

```env
NEXT_PUBLIC_GO_API_URL=http://api.grandexchange.gg:8888
```

**Note**: This won't have SSL (HTTP only), but will work.

---

## My Recommendation: Use Port 80 + Flexible SSL

**Best balance of simplicity and security**:

### Portainer:
```yaml
services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    container_name: osrs-price-api
    restart: unless-stopped
    ports:
      - "80:8080"
    environment:
      DATABASE_URL: postgresql://postgres:VGXGodgBrcUFmlKlUQDGFmLjBzUogxhJ@switchyard.proxy.rlwy.net:22411/railway
      PORT: 8080
```

### Cloudflare:
- **SSL Mode**: Flexible
- **DNS Proxy**: 🟠 On

### Frontend:
```env
NEXT_PUBLIC_GO_API_URL=https://api.grandexchange.gg
```

### Result:
- ✅ Clean URL (no port)
- ✅ HTTPS for users
- ✅ Free SSL from Cloudflare
- ✅ Simple setup

---

## Quick Test Commands

```bash
# Test if server responds on port 80
curl http://YOUR-SERVER-IP:80/health

# Test if server responds on port 8888
curl http://YOUR-SERVER-IP:8888/health

# Test through Cloudflare (after fixing)
curl https://api.grandexchange.gg/health
```

## Action Items

**Do this now**:

1. ✅ **Cloudflare**: SSL Mode → **Flexible**
2. ✅ **Portainer**: Change port to `80:8080`
3. ✅ **Test**: `curl https://api.grandexchange.gg/health`
4. ✅ **Frontend**: Already updated to use `https://api.grandexchange.gg`

That's it! Should work in 2 minutes. 🚀