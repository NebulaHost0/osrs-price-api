# Portainer Deployment Guide

Deploy the OSRS Price API backend to your Portainer server using Docker.

## Overview

The backend is automatically built as a Docker image and published to GitHub Container Registry (ghcr.io) on every push to `main`.

## Prerequisites

- ✅ Portainer installed and running
- ✅ PostgreSQL database (can be in Portainer or external)
- ✅ Backend pushed to GitHub with Actions enabled

## Step 1: Get Your Docker Image URL

After pushing to GitHub, the workflow builds and publishes:
```
ghcr.io/YOUR-USERNAME/osrs-price-api:latest
```

Example:
```
ghcr.io/yourusername/osrs-price-api:latest
```

## Step 2: Deploy via Portainer

### Option A: Using Portainer Stacks (Recommended)

1. **Log into Portainer**
2. **Go to**: Stacks → Add Stack
3. **Name**: `osrs-price-api`
4. **Web editor**: Paste this docker-compose:

```yaml
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:16-alpine
    container_name: osrs-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: YOUR_SECURE_PASSWORD_HERE
      POSTGRES_DB: osrs_prices
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - osrs-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  # OSRS Price API
  api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    container_name: osrs-price-api
    restart: unless-stopped
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      DATABASE_URL: postgresql://postgres:YOUR_SECURE_PASSWORD_HERE@postgres:5432/osrs_prices?sslmode=disable
      PORT: 8080
    ports:
      - "8080:8080"
    networks:
      - osrs-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

volumes:
  postgres_data:

networks:
  osrs-network:
    driver: bridge
```

5. **Replace**:
   - `YOUR-USERNAME` with your GitHub username
   - `YOUR_SECURE_PASSWORD_HERE` with a strong password

6. **Click**: Deploy the stack

### Option B: Using Portainer Containers

#### Deploy PostgreSQL:

1. **Containers** → Add Container
2. **Name**: `osrs-postgres`
3. **Image**: `postgres:16-alpine`
4. **Port mapping**: `5432:5432`
5. **Environment variables**:
   - `POSTGRES_USER=postgres`
   - `POSTGRES_PASSWORD=YOUR_PASSWORD`
   - `POSTGRES_DB=osrs_prices`
6. **Restart policy**: Unless stopped
7. **Deploy**

#### Deploy API:

1. **Containers** → Add Container
2. **Name**: `osrs-price-api`
3. **Image**: `ghcr.io/YOUR-USERNAME/osrs-price-api:latest`
4. **Port mapping**: `8080:8080`
5. **Environment variables**:
   - `DATABASE_URL=postgresql://postgres:PASSWORD@osrs-postgres:5432/osrs_prices?sslmode=disable`
   - `PORT=8080`
6. **Network**: Link to `osrs-postgres` container
7. **Restart policy**: Unless stopped
8. **Deploy**

## Step 3: Authenticate with GitHub Container Registry

If the image is private, you need to authenticate Portainer:

1. **Portainer** → Registries → Add Registry
2. **Type**: Custom Registry
3. **Name**: GitHub Container Registry
4. **Registry URL**: `ghcr.io`
5. **Username**: Your GitHub username
6. **Password**: Your GitHub Personal Access Token
   - Create token: GitHub → Settings → Developer settings → Personal access tokens
   - Scopes needed: `read:packages`

## Step 4: Verify Deployment

### Check Container Logs:
1. **Portainer** → Containers → `osrs-price-api` → Logs
2. Look for:
   ```
   Starting OSRS Price API server on port 8080
   Starting price fetcher worker (interval: 5m)
   Starting cleanup worker (interval: 24h)
   ```

### Test the API:
```bash
# From your local machine or server
curl http://YOUR-SERVER-IP:8080/health

# Should return:
{"status":"ok","service":"osrs-price-api"}
```

## Step 5: Update Frontend Configuration

Update your frontend `.env.local`:
```env
NEXT_PUBLIC_GO_API_URL=http://YOUR-SERVER-IP:8080
```

Or with a domain:
```env
NEXT_PUBLIC_GO_API_URL=https://api.yourdomain.com
```

## Using External PostgreSQL

If you have an existing PostgreSQL instance:

```yaml
services:
  api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    restart: unless-stopped
    environment:
      DATABASE_URL: postgresql://user:password@external-db-host:5432/dbname?sslmode=require
      PORT: 8080
    ports:
      - "8080:8080"
```

Examples:
- **Railway**: `postgresql://user:pass@xxx.railway.app:5432/dbname`
- **Supabase**: `postgresql://postgres:pass@db.xxx.supabase.co:5432/postgres`
- **AWS RDS**: `postgresql://user:pass@xxx.rds.amazonaws.com:5432/dbname`

## Updating the Container

When new code is pushed and CI builds a new image:

### Auto-update (using Watchtower):
```yaml
services:
  watchtower:
    image: containrrr/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --interval 300 osrs-price-api
```

### Manual update via Portainer:
1. **Containers** → `osrs-price-api`
2. **Duplicate/Edit** → Check "Re-pull image"
3. **Deploy**

### Manual update via CLI:
```bash
docker pull ghcr.io/YOUR-USERNAME/osrs-price-api:latest
docker stop osrs-price-api
docker rm osrs-price-api
docker run -d \
  --name osrs-price-api \
  --restart unless-stopped \
  -p 8080:8080 \
  -e DATABASE_URL="postgresql://..." \
  ghcr.io/YOUR-USERNAME/osrs-price-api:latest
```

## Environment Variables Reference

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://user:pass@host:5432/db` |
| `PORT` | API port | `8080` |

## Reverse Proxy Setup (Optional)

### Using Nginx Proxy Manager in Portainer:

1. Deploy Nginx Proxy Manager
2. Add Proxy Host:
   - **Domain**: `api.yourdomain.com`
   - **Forward to**: `osrs-price-api:8080`
   - **SSL**: Let's Encrypt
3. Update frontend to use: `https://api.yourdomain.com`

### Using Traefik:

Add labels to your API container:
```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.api.rule=Host(`api.yourdomain.com`)"
  - "traefik.http.services.api.loadbalancer.server.port=8080"
```

## Monitoring

### Container Health:

Portainer shows:
- ✅ Green = Healthy
- 🟡 Yellow = Starting
- ❌ Red = Unhealthy

### Check Logs:
```bash
# Via Portainer UI
Containers → osrs-price-api → Logs

# Via Docker CLI
docker logs -f osrs-price-api
```

### Monitor Performance:
```bash
# Via Portainer UI
Containers → osrs-price-api → Stats

# Shows:
- CPU usage
- Memory usage
- Network I/O
```

## Backup Strategy

### Database Backups:

```bash
# Automated backup script
docker exec osrs-postgres pg_dump -U postgres osrs_prices > backup_$(date +%Y%m%d).sql

# Restore
docker exec -i osrs-postgres psql -U postgres osrs_prices < backup_20251005.sql
```

Add to cron:
```bash
# Backup daily at 2 AM
0 2 * * * docker exec osrs-postgres pg_dump -U postgres osrs_prices > /backups/osrs_$(date +\%Y\%m\%d).sql
```

## Troubleshooting

### Container won't start:
```bash
# Check logs
docker logs osrs-price-api

# Common issues:
# - Database connection failed → Check DATABASE_URL
# - Port already in use → Change port mapping
# - Permission denied → Check user permissions
```

### Can't pull image:
```bash
# Make image public or add registry auth in Portainer
# Settings → Registries → Add GitHub Container Registry
```

### Database connection failed:
```bash
# Test database connectivity
docker exec osrs-price-api wget -O- postgres:5432

# Check if postgres is running
docker ps | grep postgres
```

### High memory usage:
```bash
# Set memory limits in docker-compose.yml
deploy:
  resources:
    limits:
      memory: 512M
    reservations:
      memory: 256M
```

## Performance Tuning

### For Production:

```yaml
services:
  api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
    environment:
      DATABASE_URL: postgresql://...
      PORT: 8080
      # Optional: Tune Go runtime
      GOGC: 100
      GOMAXPROCS: 2
```

## Security Best Practices

1. ✅ **Use secrets** for passwords (Portainer → Secrets)
2. ✅ **Run as non-root** (already configured in Dockerfile)
3. ✅ **Enable health checks** (already configured)
4. ✅ **Use SSL** with reverse proxy
5. ✅ **Restrict ports** - Only expose what's needed
6. ✅ **Regular updates** - Pull latest images weekly

## Quick Commands Reference

```bash
# Pull latest image
docker pull ghcr.io/YOUR-USERNAME/osrs-price-api:latest

# Start container
docker run -d --name osrs-api -p 8080:8080 -e DATABASE_URL="..." ghcr.io/YOUR-USERNAME/osrs-price-api:latest

# View logs
docker logs -f osrs-price-api

# Restart container
docker restart osrs-price-api

# Check health
curl http://localhost:8080/health

# Execute command in container
docker exec -it osrs-price-api /bin/sh

# Remove container
docker stop osrs-price-api && docker rm osrs-price-api
```

## Support

Having issues? Check:
1. Container logs in Portainer
2. Database connectivity
3. Environment variables are set correctly
4. Port 8080 is not already in use
5. GitHub Actions workflow succeeded