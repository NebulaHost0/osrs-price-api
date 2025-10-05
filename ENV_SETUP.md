# Environment Setup Guide

## .env File Support

The OSRS Price API now automatically loads environment variables from a `.env` file using [godotenv](https://github.com/joho/godotenv).

## Setup

### 1. Create .env File

Copy the example file:
```bash
cp .env.example .env
```

### 2. Edit .env File

Open `.env` and add your database connection string:

```env
# For Railway, Render, Heroku, or other cloud providers
DATABASE_URL=postgresql://user:password@host.railway.app:5432/railway?sslmode=require

# For local development
# DATABASE_URL=postgresql://postgres:postgres@localhost:5432/osrs_prices?sslmode=disable

# Server Configuration
PORT=8080
```

### 3. Run the App

The app will automatically load the `.env` file:
```bash
go run main.go
```

You should see:
```
2025/10/05 00:42:48 No .env file found, using environment variables
# OR if .env exists, no message (it loads silently)
```

## Cloud Database Examples

### Railway
```env
DATABASE_URL=postgresql://postgres:password@containers-us-west-123.railway.app:5432/railway?sslmode=require
```

### Render
```env
DATABASE_URL=postgresql://user:password@dpg-xxx-a.oregon-postgres.render.com/dbname?sslmode=require
```

### Heroku
```env
DATABASE_URL=postgresql://user:password@ec2-xxx.compute-1.amazonaws.com:5432/dbname?sslmode=require
```

### Supabase
```env
DATABASE_URL=postgresql://postgres:password@db.xxx.supabase.co:5432/postgres?sslmode=require
```

### Local PostgreSQL
```env
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/osrs_prices?sslmode=disable
```

## Fallback to Environment Variables

If no `.env` file exists, the app will use environment variables:

```bash
# Export directly in your shell
export DATABASE_URL="postgresql://user:password@host:5432/database"
go run main.go
```

## Fallback to Individual Parameters

If `DATABASE_URL` is not set, the app will construct the connection string from individual parameters:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=osrs_prices
DB_SSLMODE=disable
```

## Testing Connection

### Check if .env is loaded
```bash
# This should work now
go run main.go
```

### Test connection directly
```bash
# Source the .env file in your shell
export $(cat .env | grep -v '^#' | xargs)

# Test with psql
psql $DATABASE_URL -c "SELECT 1;"
```

### Debug connection issues
```bash
# Check what DATABASE_URL is being used
echo $DATABASE_URL

# View connection without exposing password
echo $DATABASE_URL | sed 's/:[^@]*@/:***@/'
```

## Security

### Never commit .env to git!

The `.gitignore` file already includes:
```
.env
.env.local
```

### For production, use:
- Environment variables in your deployment platform
- Secrets management (AWS Secrets Manager, Vault, etc.)
- Don't use `.env` files in production

## Troubleshooting

### "connection refused"
- Database is not running
- Wrong host/port in connection string
- Firewall blocking connection

### "authentication failed"
- Wrong username/password
- Check connection string syntax

### "database does not exist"
- Create the database first
- Or use a different database name

### "No .env file found"
- This is OK! It will use environment variables
- Create a `.env` file if you want to use one

## Migration Tool

The migration tool also loads `.env`:

```bash
go run cmd/migrate/main.go -command=up
go run cmd/migrate/main.go -command=status
```

## Make Commands

All make commands respect the `.env` file:

```bash
make migrate-up     # Uses .env
make migrate-status # Uses .env  
make run           # Uses .env
```