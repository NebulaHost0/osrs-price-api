# Cloudflare SSL Setup for API

Guide to secure `api.grandexchange.gg` with SSL through Cloudflare.

## Option 1: Cloudflare Proxy (Easiest - Recommended)

Cloudflare handles SSL automatically when you proxy through them.

### Setup:

1. **Cloudflare Dashboard** → DNS → Records
   - **Type**: A
   - **Name**: `api`
   - **IPv4**: Your server IP
   - **Proxy status**: 🟠 **Proxied** (orange cloud)
   - **TTL**: Auto

2. **SSL/TLS Settings**:
   - Go to: SSL/TLS → Overview
   - Select: **Full** or **Full (strict)**
   - Full = Cloudflare ↔ Your Server (any SSL)
   - Full (strict) = Cloudflare ↔ Your Server (valid SSL required)

3. **That's it!** Your API is now accessible at:
   ```
   https://api.grandexchange.gg:8888
   ```

### Benefits:
✅ **Free SSL** - Cloudflare provides certificate  
✅ **DDoS Protection** - Built-in protection  
✅ **Caching** - Optional API response caching  
✅ **Analytics** - Traffic insights  
✅ **Zero configuration** - Works immediately  

### Cloudflare Settings:

**SSL/TLS** → **Edge Certificates**:
- ✅ Always Use HTTPS: **On**
- ✅ Minimum TLS Version: **TLS 1.2**
- ✅ Automatic HTTPS Rewrites: **On**

**Firewall**:
- Create a rule to allow only Cloudflare IPs (optional)

---

## Option 2: Cloudflare Origin Certificate (End-to-End SSL)

For full encryption between Cloudflare and your server.

### Step 1: Create Origin Certificate

1. **Cloudflare** → **SSL/TLS** → **Origin Server**
2. Click **Create Certificate**
3. **Private key type**: RSA (2048)
4. **Hostnames**: 
   - `api.grandexchange.gg`
   - `*.grandexchange.gg`
5. **Certificate validity**: 15 years
6. Click **Create**

7. **Save both**:
   - Origin Certificate → Save as `cloudflare-cert.pem`
   - Private Key → Save as `cloudflare-key.pem`

### Step 2: Add Caddy as Reverse Proxy

Create a new stack in Portainer:

```yaml
version: '3.8'

services:
  caddy:
    image: caddy:2-alpine
    container_name: caddy-proxy
    restart: unless-stopped
    ports:
      - "443:443"
      - "8888:8888"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - ./cloudflare-cert.pem:/etc/caddy/cloudflare-cert.pem
      - ./cloudflare-key.pem:/etc/caddy/cloudflare-key.pem
      - caddy_data:/data
      - caddy_config:/config
    networks:
      - osrs-network

  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    container_name: osrs-price-api
    restart: unless-stopped
    environment:
      DATABASE_URL: postgresql://...
      PORT: 8080
    networks:
      - osrs-network

volumes:
  caddy_data:
  caddy_config:

networks:
  osrs-network:
    driver: bridge
```

**Caddyfile**:
```
api.grandexchange.gg:8888 {
    tls /etc/caddy/cloudflare-cert.pem /etc/caddy/cloudflare-key.pem
    reverse_proxy osrs-api:8080
}
```

---

## Option 3: Let's Encrypt with Nginx Proxy Manager

If you're using Nginx Proxy Manager in Portainer:

### Setup:

1. **Deploy API** (without exposing port publicly):
```yaml
services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    restart: unless-stopped
    environment:
      DATABASE_URL: postgresql://...
      PORT: 8080
    networks:
      - proxy-network
```

2. **In Nginx Proxy Manager**:
   - Add Proxy Host
   - **Domain Names**: `api.grandexchange.gg`
   - **Scheme**: `http`
   - **Forward Hostname/IP**: `osrs-api`
   - **Forward Port**: `8080`
   - **SSL** tab:
     - ✅ Force SSL
     - ✅ HTTP/2 Support
     - ✅ Request new SSL Certificate (Let's Encrypt)
   - **Advanced** tab (optional):
     ```nginx
     # Custom port
     listen 8888 ssl http2;
     ```

---

## Option 4: Expose Port 8888 with Custom Port

If you want to keep the custom port 8888:

### Update Portainer Stack:

```yaml
services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    restart: unless-stopped
    ports:
      - "8888:8080"  # Map external 8888 to internal 8080
    environment:
      DATABASE_URL: postgresql://...
      PORT: 8080
```

### Cloudflare Configuration:

**Problem**: Cloudflare proxy only works on specific ports:
- 80, 443 (HTTP/HTTPS)
- 8080, 8443 (HTTP/HTTPS alternative)
- Not 8888 ❌

**Solutions**:

**A) Use port 8080 or 8443**:
```yaml
ports:
  - "8443:8080"  # Use 8443 instead
```
Update frontend: `https://api.grandexchange.gg:8443`

**B) Use standard HTTPS port 443**:
```yaml
ports:
  - "443:8080"
```
Update frontend: `https://api.grandexchange.gg` (no port needed)

**C) DNS only (no proxy)**:
- Cloudflare DNS → Gray cloud (DNS only)
- Then use port 8888
- ⚠️ Need your own SSL certificate

---

## Recommended Setup

**For Production** (Easiest):

### 1. Update Backend Port Mapping:
```yaml
services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    restart: unless-stopped
    ports:
      - "443:8080"  # Map HTTPS port to API
    environment:
      DATABASE_URL: postgresql://postgres:VGXGodgBrcUFmlKlUQDGFmLjBzUogxhJ@switchyard.proxy.rlwy.net:22411/railway
      PORT: 8080
```

### 2. Cloudflare DNS:
- **Type**: A
- **Name**: api
- **Target**: Your server IP
- **Proxy**: 🟠 **Proxied**

### 3. Cloudflare SSL:
- **Mode**: Full (or Full strict if you add origin cert)

### 4. Update Frontend:
```typescript
const GO_API_URL = process.env.NEXT_PUBLIC_GO_API_URL || 'https://api.grandexchange.gg';
```

**Result**: Your API works at `https://api.grandexchange.gg` with SSL!

---

## Quick Fix for Current Setup (Port 8888)

If you want to keep port 8888:

### Option A: Turn off Cloudflare Proxy
1. Cloudflare DNS → `api` record
2. Click orange cloud → Gray cloud (DNS only)
3. Install SSL certificate on your server

### Option B: Change to Supported Port
1. Stop container in Portainer
2. Change port mapping to `8443:8080` or `443:8080`
3. Update Cloudflare proxy (keeps orange cloud)
4. Update frontend URL

---

## Testing

```bash
# Test with current port
curl https://api.grandexchange.gg:8888/health

# If using standard HTTPS port
curl https://api.grandexchange.gg/health
```

## What I Recommend

**Use port 443** (standard HTTPS) with Cloudflare proxy:

1. **No custom port needed** in URL
2. **Free SSL** from Cloudflare
3. **Cleaner URLs**: `https://api.grandexchange.gg`
4. **Better compatibility**
5. **DDoS protection**

Just update your Portainer stack to use port `443:8080` instead of `8888:8080`!