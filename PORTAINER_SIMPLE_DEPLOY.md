# Simple Portainer Deployment (External Database)

Deploy the OSRS Price API to Portainer using your existing external database.

## Prerequisites

✅ Your `.env` file is configured with external database  
✅ Backend code is pushed to GitHub  
✅ GitHub Actions has built the Docker image

## Your .env File

Make sure your `.env` file looks like this before building:

```env
# External database (Railway, Supabase, DigitalOcean, etc.)
DATABASE_URL=postgresql://postgres:VGXGodgBrcUFmlKlUQDGFmLjBzUogxhJ@switchyard.proxy.rlwy.net:22411/railway

# Server port
PORT=8080
```

## Step 1: Build Docker Image with .env

The `.env` file will be included in your Docker image automatically.

**Option A - Let GitHub Actions Build**:
```bash
# Just push to GitHub
git add .
git commit -m "Add Docker support with external database"
git push

# GitHub Actions will build and push to ghcr.io
```

**Option B - Build Locally**:
```bash
# Build with .env included
docker build -t osrs-price-api .

# Test locally first
docker run -p 8080:8080 osrs-price-api

# Tag for your registry (optional)
docker tag osrs-price-api your-registry/osrs-price-api:latest
docker push your-registry/osrs-price-api:latest
```

## Step 2: Deploy to Portainer

### Using Portainer Stacks (Recommended)

1. **Open Portainer** → Stacks → Add Stack

2. **Name**: `osrs-price-api`

3. **Web editor** - Paste this:

```yaml
version: '3.8'

services:
  osrs-api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    container_name: osrs-price-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

4. **Replace** `YOUR-USERNAME` with your GitHub username

5. **Deploy the stack**

### Using Portainer Containers (Alternative)

1. **Containers** → Add Container

2. **Configuration**:
   - **Name**: `osrs-price-api`
   - **Image**: `ghcr.io/YOUR-USERNAME/osrs-price-api:latest`

3. **Network ports**:
   - **Host**: `8080`
   - **Container**: `8080`

4. **Restart policy**: Unless stopped

5. **Deploy container**

## Step 3: Verify Deployment

### Check Container Status:
1. **Portainer** → Containers → `osrs-price-api`
2. Status should be **healthy** (green)

### View Logs:
1. **Portainer** → Containers → `osrs-price-api` → Logs
2. Look for:
   ```
   No .env file found, using environment variables
   Starting OSRS Price API server on port 8080
   Starting price fetcher worker (interval: 5m0s)
   Starting cleanup worker (interval: 24h0m0s)
   Fetching latest prices from OSRS Wiki API...
   Successfully saved prices to database
   ```

### Test API:
```bash
# Replace YOUR-SERVER-IP with your server's IP
curl http://YOUR-SERVER-IP:8080/health

# Should return:
{"status":"ok","service":"osrs-price-api"}

# Test prices endpoint
curl http://YOUR-SERVER-IP:8080/api/v1/prices/4151
```

## Important Notes

### Database Connectivity

Your Docker container needs to reach your external database:

✅ **Railway/Render/External Cloud**:
- Use the public connection string
- Example: `postgresql://user:pass@containers-us-west-123.railway.app:5432/railway`

✅ **Same Server/Network**:
- Use internal IP or hostname
- Example: `postgresql://postgres:pass@192.168.1.10:5432/osrs_prices`

✅ **Firewall Rules**:
- Ensure your database allows connections from Docker container IP
- Railway/Render/Supabase allow connections from anywhere by default

## Updating the Container

When you push new code to GitHub:

1. **GitHub Actions** builds new Docker image
2. **In Portainer**:
   - Go to Containers → `osrs-price-api`
   - Click **Recreate**
   - Check ✅ **Pull latest image**
   - Click **Recreate container**

Or use the stack update:
1. Stacks → `osrs-price-api`
2. Click **Pull and redeploy**

## .env File Contents Reference

```env
# Required
DATABASE_URL=postgresql://user:password@host:port/database?sslmode=require

# Optional
PORT=8080

# Examples of DATABASE_URL formats:

# Railway
DATABASE_URL=postgresql://postgres:password@containers-us-west-123.railway.app:5432/railway

# Supabase (connection pooling)
DATABASE_URL=postgresql://postgres.abcdefghij:password@aws-0-us-west-1.pooler.supabase.com:5432/postgres

# DigitalOcean
DATABASE_URL=postgresql://doadmin:password@db-postgresql-nyc3-12345.ondigitalocean.com:25060/defaultdb?sslmode=require

# AWS RDS
DATABASE_URL=postgresql://admin:password@mydb.abcdefg.us-east-1.rds.amazonaws.com:5432/osrs_prices

# Local/Self-hosted
DATABASE_URL=postgresql://postgres:password@192.168.1.100:5432/osrs_prices?sslmode=disable
```

## Troubleshooting

### Container starts but exits immediately

```bash
# Check logs
docker logs osrs-price-api

# Common issues:
# "failed to connect to database" → Check DATABASE_URL
# "bind: address already in use" → Port 8080 is taken
```

### Database connection fails

```bash
# Test connectivity from inside container
docker exec -it osrs-price-api sh
# Then try to connect to your database
```

### .env file not found

```bash
# Verify .env is in the image
docker run --rm ghcr.io/YOUR-USERNAME/osrs-price-api:latest ls -la

# Should see .env in the output
```

### Can't pull image from ghcr.io

**Make image public**:
1. Go to: `https://github.com/users/YOUR-USERNAME/packages`
2. Click on `osrs-price-api`
3. **Package settings** → Change visibility → Public

Or **authenticate Portainer**:
1. Portainer → Registries → Add registry
2. Type: Custom
3. URL: `ghcr.io`
4. Username: Your GitHub username
5. Password: GitHub Personal Access Token (with `read:packages` scope)

## Complete Deployment Checklist

- [ ] `.env` file configured with external database
- [ ] Code pushed to GitHub
- [ ] GitHub Actions completed successfully
- [ ] Docker image available at ghcr.io
- [ ] Portainer stack created
- [ ] Container deployed and running
- [ ] Container shows "healthy" status
- [ ] API health check passes: `curl http://SERVER-IP:8080/health`
- [ ] Prices are being fetched (check logs)
- [ ] Frontend connected to API

## Next Steps

1. ✅ Set up reverse proxy (Nginx Proxy Manager)
2. ✅ Add SSL certificate
3. ✅ Configure domain (api.yourdomain.com)
4. ✅ Set up monitoring alerts
5. ✅ Configure automatic backups