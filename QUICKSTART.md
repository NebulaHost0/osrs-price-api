# Quick Start Guide

Get the OSRS Price API running in under 5 minutes!

## Option 1: Cloud Database (Railway, Render, etc.) ⭐

If you already have a PostgreSQL database on Railway, Render, Heroku, or Supabase:

```bash
# 1. Create/Edit your .env file
cp .env.example .env
# Add your DATABASE_URL to .env

# 2. Run the application
go run main.go
```

The API will be available at `http://localhost:8080` and will automatically:
- Load your `.env` file
- Connect to your database
- Run migrations
- Start fetching prices

## Option 2: Docker (Local Development)

Using Docker for PostgreSQL:

```bash
# 1. Start PostgreSQL
docker-compose up -d

# 2. Wait a few seconds for the database to be ready
sleep 5

# 3. Make sure .env has local connection
# DATABASE_URL=postgresql://postgres:postgres@localhost:5432/osrs_prices?sslmode=disable

# 4. Run the application
go run main.go
```

## Option 3: Local PostgreSQL

If you have PostgreSQL installed locally:

```bash
# 1. Create the database
createdb osrs_prices

# 2. Create .env file with local connection
echo 'DATABASE_URL=postgresql://postgres:postgres@localhost:5432/osrs_prices?sslmode=disable' > .env

# 3. Run the application
go run main.go
```

## Option 3: Using Makefile

```bash
# Start everything (database + app)
make dev

# Or step by step:
make docker-up    # Start database
make run          # Start app
```

## Testing the API

Once the server is running:

```bash
# Health check
curl http://localhost:8080/health

# Get current price for Abyssal Whip (ID: 4151)
curl http://localhost:8080/api/v1/prices/4151

# Wait 5-10 minutes for some data to accumulate, then try:

# Get price history
curl http://localhost:8080/api/v1/history/4151?hours=1

# Get top gainers
curl http://localhost:8080/api/v1/gainers?limit=5

# Get top losers
curl http://localhost:8080/api/v1/losers?limit=5
```

## What Happens on Startup?

1. **Database Connection**: Connects to PostgreSQL
2. **Auto-Migration**: Creates the `price_history` table if it doesn't exist
3. **Initial Fetch**: Immediately fetches current prices from OSRS Wiki
4. **Background Worker**: Starts fetching prices every 5 minutes
5. **API Server**: Starts listening on port 8080

## Managing Migrations

### Check Migration Status
```bash
make migrate-status
```

### Run Migrations Manually
```bash
make migrate-up
```

For more details on migrations, see [MIGRATIONS.md](MIGRATIONS.md).

## Stopping the Application

Press `Ctrl+C` to gracefully shut down the server.

To stop the database:
```bash
docker-compose down
```

## Troubleshooting

### "Failed to connect to database"
- Make sure PostgreSQL is running: `docker-compose ps`
- Check your connection settings in environment variables

### "No price history data"
- The app needs to run for at least 5-10 minutes to accumulate historical data
- Current prices work immediately

### Port 8080 already in use
```bash
PORT=3000 go run main.go
```

## Next Steps

- Read the full [README.md](README.md) for detailed API documentation
- Check out example queries in the README
- Set up your own item tracking scripts!