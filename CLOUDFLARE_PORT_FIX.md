# Cloudflare Port 8888 Fix

## Problem

Port 8888 is **NOT supported** by Cloudflare's proxy (orange cloud).

Cloudflare HTTPS proxy only works on these ports:
- 443, 2053, 2083, 2087, 2096, **8443**

## ✅ Solution: Change to Port 8443

Port 8443 is the closest alternative to 8888 and **IS supported** by Cloudflare.

### Step 1: Update Your Portainer Container

1. **Portainer** → Stacks → `osrs-price-api` → Editor

2. **Change port mapping**:
```yaml
ports:
  - "8443:8080"  # Changed from 8888:8080
```

3. **Update the stack**

**Or via Container UI**:
1. Containers → `osrs-price-api` → Duplicate/Edit
2. Port mapping → Change `8888` to `8443`
3. Deploy

### Step 2: Update Cloudflare DNS

1. **Cloudflare Dashboard** → DNS
2. **Record**: `api.grandexchange.gg`
3. **Proxy status**: 🟠 **Proxied** (orange cloud - enabled)

### Step 3: Test

```bash
# Test the new port
curl https://api.grandexchange.gg:8443/health

# Should return:
{"status":"ok","service":"osrs-price-api"}
```

### Step 4: Update Frontend (Already Done!)

The frontend is already updated to use port 8443.

---

## Alternative: Use Standard HTTPS Port (443)

Even better - use the standard HTTPS port:

### Portainer:
```yaml
ports:
  - "443:8080"
```

### Cloudflare:
- Same setup, orange cloud enabled

### Frontend:
```typescript
NEXT_PUBLIC_GO_API_URL=https://api.grandexchange.gg
// No port needed!
```

### Benefits:
- ✅ Clean URLs (no port number)
- ✅ Standard HTTPS
- ✅ Better compatibility
- ✅ Professional appearance

---

## Current vs Recommended Setup

### Current (Port 8888 - NOT WORKING):
```
Browser → Cloudflare (❌ Port 8888 blocked) → Your Server
```

### Fix Option 1 (Port 8443 - WORKS):
```
Browser → https://api.grandexchange.gg:8443
        ↓
   Cloudflare Proxy (✅ Port 8443 allowed)
        ↓
   Your Server:8443 → Container:8080
```

### Fix Option 2 (Port 443 - BEST):
```
Browser → https://api.grandexchange.gg
        ↓
   Cloudflare Proxy (✅ Standard HTTPS)
        ↓
   Your Server:443 → Container:8080
```

---

## Cloudflare SSL Settings

After changing to supported port:

1. **SSL/TLS** → Overview
   - Mode: **Full** (Cloudflare validates your server's cert)

2. **SSL/TLS** → Edge Certificates
   - ✅ Always Use HTTPS: **On**
   - ✅ Automatic HTTPS Rewrites: **On**
   - ✅ Minimum TLS Version: TLS 1.2

3. **Firewall** (Optional)
   - Create rule to only allow Cloudflare IPs

---

## Quick Action Items

**To fix immediately**:

1. **In Portainer**: Change port `8888` → `8443`
2. **Test**: `curl https://api.grandexchange.gg:8443/health`
3. **Done!** Frontend already configured

**Better long-term**:

1. **In Portainer**: Change port to `443:8080`
2. **Test**: `curl https://api.grandexchange.gg/health`
3. **Update frontend**: Remove `:8443` from URL

---

## Verifying SSL Works

```bash
# Check SSL certificate
openssl s_client -connect api.grandexchange.gg:8443 -servername api.grandexchange.gg

# Should show Cloudflare certificate
```

## Summary

✅ **Change port 8888 → 8443** (Cloudflare supported)  
✅ **Or use 443** (even better)  
✅ **Keep Cloudflare proxy on** (orange cloud)  
✅ **Free SSL** from Cloudflare  
✅ **Frontend already updated** to use new port