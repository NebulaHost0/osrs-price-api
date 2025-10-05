# Portainer Deployment Guide - Simple Setup

Deploy OSRS Price API to Portainer with environment variables configured in the UI.

## Quick Deploy

### Step 1: Ensure Code is Pushed

```bash
git add .
git commit -m "Add Docker support"
git push
```

Wait for GitHub Actions to complete (~3-5 minutes).

### Step 2: Deploy in Portainer

#### Method 1: Using Stacks (Recommended)

1. **Open Portainer** → **Stacks** → **Add stack**

2. **Name**: `osrs-price-api`

3. **Web editor** - Paste:

```yaml
version: '3.8'

services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    container_name: osrs-price-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgresql://postgres:VGXGodgBrcUFmlKlUQDGFmLjBzUogxhJ@switchyard.proxy.rlwy.net:22411/railway
      PORT: 8080
```

4. **Replace**:
   - `YOUR-USERNAME` → Your GitHub username
   - `DATABASE_URL` → Your actual Railway database URL

5. **Click** "Deploy the stack"

#### Method 2: Using Container UI

1. **Containers** → **Add container**

2. **Configuration**:
   - **Name**: `osrs-price-api`
   - **Image**: `ghcr.io/YOUR-USERNAME/osrs-price-api:latest`

3. **Network ports**:
   - Click "publish a new network port"
   - **host**: `8080` | **container**: `8080`

4. **Env tab** - Click "Add environment variable":
   - **name**: `DATABASE_URL`
   - **value**: `postgresql://postgres:YOUR_PASSWORD@YOUR_HOST:PORT/railway`
   
   - **name**: `PORT`
   - **value**: `8080`

5. **Restart policy**: Unless stopped

6. **Deploy container**

### Step 3: Verify

1. **Check container status** - Should show "healthy" (green)

2. **View logs**:
   - Containers → `osrs-price-api` → Logs
   - Look for: "Starting OSRS Price API server on port 8080"

3. **Test API**:
```bash
curl http://YOUR-SERVER-IP:8080/health
```

## Environment Variables

Set these in Portainer:

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://user:pass@host:5432/db` |
| `PORT` | API port | `8080` |

### DATABASE_URL Examples:

**Railway**:
```
postgresql://postgres:password@containers-us-west-123.railway.app:5432/railway
```

**Supabase**:
```
postgresql://postgres.projectid:password@aws-0-us-west-1.pooler.supabase.com:5432/postgres
```

**DigitalOcean**:
```
postgresql://doadmin:password@db-postgresql-nyc3-12345.ondigitalocean.com:25060/defaultdb?sslmode=require
```

**Local PostgreSQL on Docker**:
```
postgresql://postgres:password@host.docker.internal:5432/osrs_prices
```

## Updating the Container

When new code is pushed:

1. **Portainer** → **Containers** → `osrs-price-api`
2. Click **Recreate**
3. ✅ Check "Pull latest image version"
4. Click "Recreate"

Or for stacks:
1. **Stacks** → `osrs-price-api`
2. Click **Editor** 
3. Click **Update the stack** (will pull latest)

## Adding to Existing Network

If you have other containers that need to access the API:

```yaml
services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgresql://...
      PORT: 8080
    networks:
      - your-existing-network

networks:
  your-existing-network:
    external: true
```

## Using with Nginx Proxy Manager

If you're using Nginx Proxy Manager in Portainer:

1. **Don't expose port 8080** to host:
```yaml
services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    # No ports section - only accessible via internal network
    networks:
      - proxy-network
```

2. **In Nginx Proxy Manager**:
   - Add Proxy Host
   - Domain: `api.yourdomain.com`
   - Forward to: `osrs-price-api` port `8080`
   - SSL: Let's Encrypt

3. **Update frontend**:
```env
NEXT_PUBLIC_GO_API_URL=https://api.yourdomain.com
```

## Security Best Practices

✅ **No .env in image** - Credentials not baked into image  
✅ **Environment variables** - Set directly in Portainer  
✅ **SSL connections** - Use `?sslmode=require` for database  
✅ **Non-root user** - Container runs as unprivileged user  
✅ **Health checks** - Monitors container health  
✅ **Limited logs** - Prevents disk space issues  

## Monitoring

### Container Stats in Portainer:
- CPU usage
- Memory usage
- Network I/O
- Logs

### Health Status:
- 🟢 **Healthy**: API is running and responding
- 🟡 **Starting**: Container is starting up (40s grace period)
- 🔴 **Unhealthy**: API not responding (check logs)

## Troubleshooting

### "failed to connect to database"
```bash
# Check DATABASE_URL is correct
docker exec osrs-price-api env | grep DATABASE_URL

# Test database connection
docker exec osrs-price-api wget -O- your-db-host:5432
```

### "port already in use"
```bash
# Check what's using port 8080
netstat -tulpn | grep 8080

# Use different port in Portainer
ports:
  - "8081:8080"  # Map host 8081 to container 8080
```

### Container keeps restarting
```bash
# Check logs for errors
docker logs osrs-price-api --tail 100

# Common issues:
# - Invalid DATABASE_URL
# - Database not reachable
# - Migration failed
```

## Complete Setup Example

**Portainer Stack YAML**:
```yaml
version: '3.8'

services:
  osrs-api:
    image: ghcr.io/yourusername/osrs-price-api:latest
    container_name: osrs-price-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgresql://postgres:VGXGodgBrcUFmlKlUQDGFmLjBzUogxhJ@switchyard.proxy.rlwy.net:22411/railway
      PORT: 8080
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

**That's all you need!** No PostgreSQL container, just your API connecting to Railway.

Now you can:
- ✅ Push code → GitHub Actions builds Docker image
- ✅ Update in Portainer → Pull latest & recreate
- ✅ No secrets in Git → All credentials in Portainer

Perfect setup for your Portainer server! 🐳