# Migration Cheatsheet

Quick reference for database migration commands.

## 🚀 Quick Start

```bash
# Easiest: Complete setup
make setup

# Or step by step:
make docker-up     # Start database
make migrate-up    # Run migrations
make run          # Start API
```

## 📋 Common Commands

### Check Status
```bash
make migrate-status
# OR
go run cmd/migrate/main.go -command=status
```

### Run Migrations
```bash
make migrate-up
# OR
go run cmd/migrate/main.go -command=up
```

### Reset Database (⚠️ Destructive!)
```bash
make migrate-reset
# OR
go run cmd/migrate/main.go -command=reset
```

### Manual SQL Migration
```bash
make migrate-sql FILE=migrations/001_create_price_history.sql
# OR
psql $DATABASE_URL -f migrations/001_create_price_history.sql
```

## 🎯 Development Workflows

### First Time Setup
```bash
# 1. Start database
docker-compose up -d

# 2. Wait for it to be ready
sleep 3

# 3. Run migrations
make migrate-up

# 4. Start app
go run main.go
```

### Daily Development
```bash
# Migrations run automatically!
go run main.go
```

### Fresh Start
```bash
# Nuclear option - clean slate
make db-drop
make db-create
make migrate-up
```

## 🏭 Production

### Deploy New Version
```bash
# 1. Backup first!
pg_dump $DATABASE_URL > backup.sql

# 2. Run migrations
export DATABASE_URL="postgresql://user:pass@prod:5432/osrs_prices"
go run cmd/migrate/main.go -command=up

# 3. Verify
go run cmd/migrate/main.go -command=status

# 4. Start app
./osrs-price-api
```

## 🔍 Troubleshooting

### Check Connection
```bash
psql $DATABASE_URL -c "SELECT 1;"
```

### View Current Schema
```bash
psql $DATABASE_URL -c "\dt"
psql $DATABASE_URL -c "\d price_history"
```

### Count Records
```bash
psql $DATABASE_URL -c "SELECT COUNT(*) FROM price_history;"
```

## 📚 Files

- `main.go` - Auto-migrates on app start
- `cmd/migrate/main.go` - Migration tool
- `migrations/*.sql` - SQL migration files
- `internal/database/connection.go` - Migration logic

## 🎨 Migration Modes

| Mode | Command | When to Use |
|------|---------|-------------|
| **Automatic** | `go run main.go` | Development (default) |
| **Manual Tool** | `make migrate-up` | Production, testing |
| **SQL Files** | `psql $DATABASE_URL -f ...` | Advanced users |

## ⚙️ Environment Setup

```bash
# Connection string (recommended)
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/osrs_prices?sslmode=disable"

# OR individual parameters
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=osrs_prices
```

## 📖 More Info

- Full guide: [MIGRATIONS.md](MIGRATIONS.md)
- Quick start: [QUICKSTART.md](QUICKSTART.md)
- API docs: [README.md](README.md)