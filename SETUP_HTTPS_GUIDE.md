# Setup HTTPS on Go API with Cloudflare

Complete guide to enable HTTPS on your Go API server with Cloudflare Origin Certificates.

## Overview

```
Browser → Cloudflare (HTTPS) → Your Server (HTTPS with Origin Cert) → Go API
```

This provides **end-to-end encryption** with Cloudflare's free certificates.

---

## Step 1: Generate Cloudflare Origin Certificate

1. **Log into Cloudflare Dashboard**

2. **Navigate to**: SSL/TLS → Origin Server

3. **Click**: "Create Certificate"

4. **Configuration**:
   - **Private key type**: RSA (2048)
   - **Hostnames**: 
     - `api.grandexchange.gg`
     - `*.grandexchange.gg` (optional, for wildcards)
   - **Certificate Validity**: 15 years

5. **Click**: "Create"

6. **Save the files**:
   
   **Origin Certificate** (copy all text):
   ```
   -----BEGIN CERTIFICATE-----
   MIIEpDCCA...
   -----END CERTIFICATE-----
   ```
   Save as: `origin-cert.pem`
   
   **Private Key** (copy all text):
   ```
   -----BEGIN PRIVATE KEY-----
   MIIEvgIBA...
   -----END PRIVATE KEY-----
   ```
   Save as: `origin-key.pem`

⚠️ **Important**: Save these files securely. You can't retrieve the private key later!

---

## Step 2: Upload Certificates to Your Server

### Option A: Via Portainer Volumes

1. **Portainer** → **Volumes** → **Add volume**
   - Name: `osrs-api-certs`

2. **Upload certificate files** to the volume:
   - You'll need to access the server and copy files to volume location
   - Or use Portainer's file browser if available

### Option B: Direct Server Upload

SSH into your server and create the directory:

```bash
# Create directory for certificates
mkdir -p /opt/osrs-api/certs
cd /opt/osrs-api/certs

# Upload your certificate files here
# origin-cert.pem
# origin-key.pem

# Set proper permissions
chmod 600 origin-key.pem
chmod 644 origin-cert.pem
```

---

## Step 3: Update Portainer Stack

Use this updated stack configuration:

```yaml
version: '3.8'

services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    container_name: osrs-price-api
    restart: unless-stopped
    ports:
      - "443:8443"
    volumes:
      # Mount certificate files from your server
      - /opt/osrs-api/certs/origin-cert.pem:/app/certs/cert.pem:ro
      - /opt/osrs-api/certs/origin-key.pem:/app/certs/key.pem:ro
    environment:
      DATABASE_URL: postgresql://postgres:VGXGodgBrcUFmlKlUQDGFmLjBzUogxhJ@switchyard.proxy.rlwy.net:22411/railway
      PORT: 8443
      SSL_CERT_FILE: /app/certs/cert.pem
      SSL_KEY_FILE: /app/certs/key.pem
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8443/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

**Important**: Update `/opt/osrs-api/certs/` to match where you uploaded the certs on your server.

---

## Step 4: Configure Cloudflare SSL Mode

1. **Cloudflare Dashboard** → **SSL/TLS** → **Overview**

2. **Change encryption mode to**: **Full (strict)** ✅

   This means:
   - Browser → Cloudflare: HTTPS (Cloudflare cert)
   - Cloudflare → Your Server: HTTPS (Origin cert)
   - End-to-end encryption!

3. **SSL/TLS** → **Edge Certificates**:
   - ✅ Always Use HTTPS: **On**
   - ✅ Minimum TLS Version: **TLS 1.2**

---

## Step 5: Deploy and Test

1. **Deploy the stack** in Portainer

2. **Check container logs**:
   ```
   Starting OSRS Price API server with TLS on port 8443
   Using certificate: /app/certs/cert.pem
   ```

3. **Test locally on server**:
   ```bash
   curl https://localhost:8443/health
   ```

4. **Test through Cloudflare**:
   ```bash
   curl https://api.grandexchange.gg/health
   ```

5. **Test frontend**:
   - Start: `npm run dev`
   - Should connect successfully!

---

## Alternative: Use Let's Encrypt with Certbot

If you prefer Let's Encrypt certificates instead of Cloudflare Origin:

### Setup Certbot on Server:

```bash
# Install certbot
sudo apt install certbot  # Ubuntu/Debian
brew install certbot      # macOS

# Get certificate (stop API first)
sudo certbot certonly --standalone -d api.grandexchange.gg

# Certificates will be at:
# /etc/letsencrypt/live/api.grandexchange.gg/fullchain.pem
# /etc/letsencrypt/live/api.grandexchange.gg/privkey.pem
```

### Update Portainer Stack:

```yaml
volumes:
  - /etc/letsencrypt/live/api.grandexchange.gg/fullchain.pem:/app/certs/cert.pem:ro
  - /etc/letsencrypt/live/api.grandexchange.gg/privkey.pem:/app/certs/key.pem:ro
```

---

## Troubleshooting

### Container won't start

```bash
# Check logs
docker logs osrs-price-api

# Common issues:
# "open /app/certs/cert.pem: no such file or directory"
# → Certificate files not mounted correctly

# "permission denied"
# → Certificate files need proper permissions
```

### SSL handshake still fails

```bash
# Test directly on server (bypass Cloudflare)
curl -k https://YOUR-SERVER-IP:443/health

# Check certificate is valid
openssl s_client -connect localhost:8443 -servername api.grandexchange.gg
```

### Cloudflare shows 525 error

- SSL handshake failed between Cloudflare and origin
- Check certificate is valid for `api.grandexchange.gg`
- Ensure port 443 is open on your server firewall

---

## Complete Portainer Stack (Production Ready)

```yaml
version: '3.8'

services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    container_name: osrs-price-api
    restart: unless-stopped
    ports:
      - "443:8443"
    volumes:
      # IMPORTANT: Update this path to where you uploaded the certificates
      - /opt/osrs-api/certs/origin-cert.pem:/app/certs/cert.pem:ro
      - /opt/osrs-api/certs/origin-key.pem:/app/certs/key.pem:ro
    environment:
      DATABASE_URL: postgresql://postgres:VGXGodgBrcUFmlKlUQDGFmLjBzUogxhJ@switchyard.proxy.rlwy.net:22411/railway
      PORT: 8443
      SSL_CERT_FILE: /app/certs/cert.pem
      SSL_KEY_FILE: /app/certs/key.pem
    logging:
      driver: json-file
      options:
        max-size: 10m
        max-file: "3"
```

---

## Quick Start Checklist

- [ ] Generate Cloudflare Origin Certificate
- [ ] Save `origin-cert.pem` and `origin-key.pem`
- [ ] Upload certificates to server at `/opt/osrs-api/certs/`
- [ ] Update Portainer stack with certificate volume mounts
- [ ] Set environment variables: `SSL_CERT_FILE`, `SSL_KEY_FILE`, `PORT=8443`
- [ ] Deploy stack
- [ ] Set Cloudflare SSL mode to **Full (strict)**
- [ ] Test: `curl https://api.grandexchange.gg/health`
- [ ] Frontend should connect successfully

---

## Benefits of This Setup

✅ **End-to-end encryption** - HTTPS all the way  
✅ **Free certificates** - Cloudflare Origin certs are free  
✅ **15-year validity** - No renewal needed  
✅ **Cloudflare Full (strict)** - Most secure mode  
✅ **No external CA** - Origin certs work only with Cloudflare  

Your Go API will now serve HTTPS directly! 🔐