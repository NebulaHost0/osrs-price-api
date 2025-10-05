# Full Stack Setup Guide

Complete setup guide for the OSRS Grand Exchange Tracker with Go backend and Next.js frontend.

## System Architecture

```
┌─────────────────┐
│   Next.js UI    │ ← User interaction
│   (Port 3000)   │
└────────┬────────┘
         │
         ├──────────────┐
         │              │
         ▼              ▼
┌────────────────┐  ┌──────────────┐
│   Supabase     │  │   Go API     │
│   (Auth/User   │  │   (Prices)   │
│    Features)   │  │ (Port 8080)  │
└────────────────┘  └───────┬──────┘
                            │
                            ▼
                    ┌──────────────┐
                    │  PostgreSQL  │
                    │ (Price Data) │
                    └──────────────┘
```

**Data Flow:**
- **Supabase**: Handles authentication, user profiles, watchlists, alerts, marketplace
- **Go API**: Handles OSRS price data, fetching, storage, and historical queries
- **PostgreSQL**: Stores price history (managed by Go API)
- **Next.js**: Frontend UI that connects to both Supabase and Go API

## Prerequisites

- **Go** 1.21 or higher
- **Node.js** 18 or higher
- **PostgreSQL** 12 or higher
- **Supabase** account (free tier works fine)

## Step 1: Database Setup

### PostgreSQL for Price Data

1. **Install PostgreSQL** (if not already installed):

```bash
# macOS
brew install postgresql@16
brew services start postgresql@16

# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql

# Windows - Download from https://www.postgresql.org/download/windows/
```

2. **Create the database**:

```bash
# Connect to PostgreSQL
psql postgres

# Create database
CREATE DATABASE osrs_prices;

# Create user (optional)
CREATE USER osrs_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE osrs_prices TO osrs_user;

# Exit
\q
```

### Supabase for User Data

1. Go to [supabase.com](https://supabase.com) and create a new project
2. Wait for your project to be provisioned (2-3 minutes)
3. Go to **Settings** → **API** and copy:
   - Project URL
   - Anon public key

## Step 2: Go Backend Setup

1. **Navigate to project root**:

```bash
cd /path/to/Project1
```

2. **Install Go dependencies**:

```bash
go mod download
```

3. **Configure environment variables**:

Create a `.env` file in the project root:

```env
# PostgreSQL connection
DATABASE_URL=postgresql://postgres:your_password@localhost:5432/osrs_prices?sslmode=disable

# Server port
PORT=8080
```

For cloud databases (Railway, Render, etc.), see [ENV_SETUP.md](ENV_SETUP.md).

4. **Run database migrations**:

```bash
# Migrations run automatically on startup, but you can run manually:
go run cmd/migrate/main.go
```

5. **Start the Go API**:

```bash
go run main.go
```

You should see:
```
Starting OSRS Price API server on port 8080
```

The API will automatically:
- ✅ Connect to PostgreSQL
- ✅ Run migrations
- ✅ Start price fetching worker (every 5 minutes)
- ✅ Expose REST API on port 8080

6. **Verify the Go API is working**:

```bash
# Health check
curl http://localhost:8080/health

# Get current prices
curl http://localhost:8080/api/v1/prices | jq
```

## Step 3: Frontend Setup

1. **Navigate to frontend directory**:

```bash
cd osrs-ge-tracker
```

2. **Install Node.js dependencies**:

```bash
npm install
# or
pnpm install
# or
yarn install
```

3. **Configure environment variables**:

Create `.env.local` file in the `osrs-ge-tracker` directory:

```env
# Supabase configuration (for authentication and user features)
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key-here

# Go API configuration (for price data)
NEXT_PUBLIC_GO_API_URL=http://localhost:8080
```

**Important:** Replace `your-project` and `your-anon-key-here` with your actual Supabase values from Step 1.

4. **Start the Next.js development server**:

```bash
npm run dev
# or
pnpm dev
```

The frontend will start on `http://localhost:3000`.

5. **Open your browser**:

Navigate to `http://localhost:3000` and you should see the OSRS GE Tracker!

## Step 4: Testing the Integration

### Test Price Data (Go API)

1. **View all prices**: Go to homepage and search for items
2. **View item detail**: Click on any item to see its details
3. **View price chart**: Charts should load with historical data
4. **View price changes**: 24h price changes should show for authenticated users

### Test User Features (Supabase)

1. **Create account**: Click "Register" and create a new account
2. **Login**: Sign in with your credentials
3. **Watchlist**: Add items to your watchlist (star icon)
4. **Alerts**: Set price alerts on items (bell icon)
5. **Marketplace**: Browse or create marketplace listings

## Common Issues & Solutions

### Go API Issues

**Problem**: `failed to connect to database`
```bash
# Solution: Check your PostgreSQL is running
brew services list | grep postgresql  # macOS
sudo systemctl status postgresql      # Linux

# Verify connection string
psql "postgresql://postgres:password@localhost:5432/osrs_prices"
```

**Problem**: `No price data available`
```bash
# Solution: Wait 5 minutes for first price fetch, or check logs
# The worker fetches prices every 5 minutes

# Check if data is being stored
psql osrs_prices -c "SELECT COUNT(*) FROM price_history;"
```

### Frontend Issues

**Problem**: `Failed to fetch prices from Go API`
```bash
# Solution 1: Ensure Go API is running
curl http://localhost:8080/health

# Solution 2: Check NEXT_PUBLIC_GO_API_URL in .env.local
echo $NEXT_PUBLIC_GO_API_URL  # Should be http://localhost:8080
```

**Problem**: Supabase authentication errors
```bash
# Solution: Verify your Supabase credentials
# Check .env.local has correct NEXT_PUBLIC_SUPABASE_URL and NEXT_PUBLIC_SUPABASE_ANON_KEY

# Test connection
curl "https://your-project.supabase.co/rest/v1/?apikey=your-anon-key"
```

**Problem**: CORS errors
```bash
# Solution: The Go API allows all origins by default
# If you see CORS errors, check that you're using http://localhost:3000 for the frontend
```

## Production Deployment

### Deploy Go API

**Option 1: Railway**
1. Create account at [railway.app](https://railway.app)
2. Create new project
3. Add PostgreSQL service
4. Add Go app and connect to GitHub
5. Set `DATABASE_URL` environment variable (Railway auto-provides this)
6. Deploy!

**Option 2: Fly.io**
```bash
fly launch
fly secrets set DATABASE_URL="your-connection-string"
fly deploy
```

**Option 3: DigitalOcean/AWS/GCP**
- Build binary: `go build -o osrs-api main.go`
- Upload to server
- Set up systemd service
- Configure reverse proxy (nginx/caddy)

### Deploy Frontend (Vercel)

1. Push code to GitHub
2. Go to [vercel.com](https://vercel.com)
3. Import your repository
4. Set environment variables:
   ```
   NEXT_PUBLIC_SUPABASE_URL=your-supabase-url
   NEXT_PUBLIC_SUPABASE_ANON_KEY=your-supabase-key
   NEXT_PUBLIC_GO_API_URL=https://your-go-api.com
   ```
5. Deploy!

## Environment Variables Reference

### Go Backend (`.env` in project root)

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://user:pass@localhost:5432/osrs_prices` |
| `PORT` | API server port | `8080` |

### Next.js Frontend (`osrs-ge-tracker/.env.local`)

| Variable | Description | Required For |
|----------|-------------|--------------|
| `NEXT_PUBLIC_SUPABASE_URL` | Supabase project URL | Authentication, watchlists, alerts |
| `NEXT_PUBLIC_SUPABASE_ANON_KEY` | Supabase anonymous key | Authentication, watchlists, alerts |
| `NEXT_PUBLIC_GO_API_URL` | Go API endpoint | Price data, charts |

## Development Workflow

### Starting Development

```bash
# Terminal 1: Start Go API
go run main.go

# Terminal 2: Start Next.js
cd osrs-ge-tracker && npm run dev
```

### Making Changes

**Backend changes:**
1. Edit Go files in `internal/`
2. Restart Go server (Ctrl+C, then `go run main.go`)
3. Test API with curl or frontend

**Frontend changes:**
1. Edit files in `osrs-ge-tracker/src/`
2. Next.js hot reloads automatically
3. Check browser console for errors

### Database Migrations

```bash
# Check current migration status
go run cmd/migrate/main.go -command=status

# Create new migration
go run cmd/migrate/main.go -command=create -name=add_new_feature

# Run migrations
go run cmd/migrate/main.go -command=up
```

See [MIGRATIONS.md](MIGRATIONS.md) for details.

## Architecture Details

For more information:
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [osrs-ge-tracker/MIGRATION_TO_GO_API.md](osrs-ge-tracker/MIGRATION_TO_GO_API.md) - Migration details
- [MIGRATIONS.md](MIGRATIONS.md) - Database migrations
- [ENV_SETUP.md](ENV_SETUP.md) - Cloud database setup

## Support

If you encounter issues:
1. Check the logs (Go API prints to stdout)
2. Check browser console (F12)
3. Verify environment variables are set correctly
4. Ensure both servers are running

## Next Steps

- ✅ Set up production deployment
- ✅ Configure custom domain
- ✅ Set up monitoring (Sentry, LogRocket, etc.)
- ✅ Enable SSL/TLS certificates
- ✅ Set up CI/CD pipeline
- ✅ Configure backup strategy for databases