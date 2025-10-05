# Docker Deployment with .env File

Deploy the OSRS Price API using Docker with your existing external database (no PostgreSQL container needed).

## Quick Start

### Option 1: Build and Run Locally

```bash
# 1. Make sure your .env file is configured
cat .env
# Should contain:
# DATABASE_URL=postgresql://user:password@your-db-host:5432/osrs_prices
# PORT=8080

# 2. Build the Docker image
docker build -t osrs-price-api .

# 3. Run the container
docker run -d \
  --name osrs-price-api \
  --restart unless-stopped \
  -p 8080:8080 \
  osrs-price-api

# 4. Check it's working
curl http://localhost:8080/health
```

### Option 2: Use Docker Compose (External DB)

```bash
# Use the external DB compose file
docker-compose -f docker-compose-external-db.yml up -d

# Check logs
docker-compose logs -f api
```

## For Portainer Deployment

### Method 1: Build Image Locally, Push to Your Registry

If you don't want to use GitHub Container Registry:

```bash
# 1. Build locally with your .env file included
docker build -t your-registry.com/osrs-price-api:latest .

# 2. Push to your private registry
docker push your-registry.com/osrs-price-api:latest

# 3. In Portainer:
#    - Containers → Add Container
#    - Image: your-registry.com/osrs-price-api:latest
#    - Port: 8080:8080
#    - Deploy!
```

### Method 2: Use Portainer Stack with Build

1. **Portainer** → Stacks → Add Stack
2. **Name**: `osrs-price-api`
3. **Build method**: Repository (or upload)
4. **Paste this**:

```yaml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: osrs-price-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    # .env file is included in the image during build
```

5. **Upload your project folder** (including .env)
6. **Deploy**

### Method 3: Pass .env via Docker Volume (Best for Portainer)

```yaml
version: '3.8'

services:
  api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    container_name: osrs-price-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      # Mount your .env file from host
      - /path/to/your/.env:/app/.env:ro
    networks:
      - osrs-network

networks:
  osrs-network:
    driver: bridge
```

#### In Portainer:

1. **Stacks** → Add Stack
2. Use the YAML above
3. **Replace** `/path/to/your/.env` with actual path on your server
4. **Deploy**

Example paths:
- Linux: `/home/user/osrs-api/.env`
- Portainer volume: `/portainer/files/osrs-api/.env`

## Verify .env is Loaded

```bash
# Check container logs for database connection
docker logs osrs-price-api

# Should see:
# "Starting OSRS Price API server on port 8080"
# "Starting price fetcher worker"

# Exec into container and check
docker exec osrs-price-api cat .env
```

## Using External Database

Your `.env` file should point to your external database:

```env
# Railway
DATABASE_URL=postgresql://postgres:password@containers-us-west-123.railway.app:5432/railway

# Supabase
DATABASE_URL=postgresql://postgres:password@db.abcdefghij.supabase.co:5432/postgres

# DigitalOcean
DATABASE_URL=postgresql://user:password@db-postgresql-sfo2-12345.ondigitalocean.com:25060/defaultdb?sslmode=require

# Your own server
DATABASE_URL=postgresql://postgres:password@192.168.1.100:5432/osrs_prices

# Port (default 8080)
PORT=8080
```

## Security Considerations

### ⚠️ Warning: .env File in Image

Including `.env` in the Docker image means:
- ✅ **Pros**: Easy deployment, no external configuration needed
- ⚠️ **Cons**: Database credentials are in the image

### Better for Production: Use Environment Variables

**Option A**: Pass env vars directly to Portainer:

```yaml
services:
  api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    environment:
      DATABASE_URL: postgresql://user:password@external-db:5432/osrs_prices
      PORT: 8080
```

**Option B**: Use Docker Secrets (Portainer supports this):

1. **Portainer** → Secrets → Add Secret
   - Name: `db_connection`
   - Value: Your DATABASE_URL
   
2. **Stack configuration**:
```yaml
services:
  api:
    image: ghcr.io/YOUR-USERNAME/osrs-price-api:latest
    secrets:
      - db_connection
    environment:
      PORT: 8080

secrets:
  db_connection:
    external: true
```

3. **Update your Go app** to read from `/run/secrets/db_connection`

## Testing the Setup

### 1. Build locally with your .env:
```bash
# Ensure .env exists
ls -la .env

# Build
docker build -t osrs-test .

# Run
docker run --rm -p 8080:8080 osrs-test

# Test
curl http://localhost:8080/health
```

### 2. Check database connection:
```bash
docker logs osrs-test

# Look for:
# ✅ "Connected to database"
# ✅ "Starting price fetcher worker"
# ❌ "Failed to connect to database" (check DATABASE_URL)
```

## Portainer Stack Template

Save this as your Portainer stack:

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
      # Replace with your actual database URL
      DATABASE_URL: postgresql://postgres:password@your-db-host.railway.app:5432/railway
      PORT: 8080
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## Update Workflow

After making these changes, commit and push:

```bash
git add Dockerfile .dockerignore docker-compose-external-db.yml
git commit -m "Add Docker support with external database"
git push
```

GitHub Actions will:
1. ✅ Build Docker image with .env included
2. ✅ Push to ghcr.io
3. ✅ Build standalone binaries
4. ✅ Create release

Then in Portainer:
1. Pull the latest image
2. Deploy with your stack configuration
3. Done! 🎉

## Troubleshooting

### Image won't pull in Portainer
```bash
# Make the image public on GitHub:
# Go to: https://github.com/users/YOUR-USERNAME/packages/container/osrs-price-api/settings
# Change visibility to: Public
```

### Container exits immediately
```bash
# Check logs
docker logs osrs-price-api

# Common issues:
# - Invalid DATABASE_URL
# - Database not accessible from container
# - Port conflict
```

### Database connection timeout
```bash
# Test database connectivity from container
docker run --rm osrs-price-api ping your-db-host

# Check firewall allows connections from Docker network
```

## Ready to Deploy!

Your Docker container will:
- ✅ Include migrations
- ✅ Load .env file automatically
- ✅ Connect to your external database
- ✅ Run price fetcher and cleanup workers
- ✅ Expose API on port 8080
- ✅ Auto-restart on failure