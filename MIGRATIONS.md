# Database Migration Guide

Complete guide for managing database migrations in the OSRS Price API.

## Quick Start

### For Development (Easiest)
```bash
# Complete setup with one command
make setup

# Or step by step:
make docker-up      # Start PostgreSQL
make migrate-up     # Run migrations
make run            # Start the app
```

### For Production
```bash
# Set your database connection string
export DATABASE_URL="postgresql://user:pass@host:5432/osrs_prices?sslmode=require"

# Run migrations
go run cmd/migrate/main.go -command=up
```

## How Migrations Work

The OSRS Price API uses **automatic migrations** by default, powered by GORM's AutoMigrate feature. When you start the application, it:

1. Connects to PostgreSQL
2. Checks if tables exist
3. Creates missing tables
4. Adds missing columns
5. Creates indexes

This happens every time the app starts, so you never have to worry about schema drift.

## Migration Methods

### Method 1: Automatic (Default) ⭐ Recommended for Development

Migrations run automatically when you start the app:

```bash
go run main.go
```

**Output:**
```
Database connection established successfully
Running database migrations...
Database migrations completed successfully
Starting OSRS Price API server on port 8080
```

**Pros:**
- ✅ Zero configuration
- ✅ Always up to date
- ✅ Perfect for development

**Cons:**
- ❌ Less control in production
- ❌ Can't roll back automatically

---

### Method 2: Migration Tool ⭐ Recommended for Production

Use the dedicated migration tool for more control:

```bash
# Run migrations
go run cmd/migrate/main.go -command=up

# Check what's been applied
go run cmd/migrate/main.go -command=status

# Reset everything (destructive!)
go run cmd/migrate/main.go -command=reset
```

**Example Output:**
```bash
$ go run cmd/migrate/main.go -command=status

Checking database tables...
✓ Table 'price_history' exists
  - Record count: 15234
  - Index 'idx_item_timestamp' exists
```

**Pros:**
- ✅ Full control over when migrations run
- ✅ Can check status before running
- ✅ Safer for production

**Cons:**
- ❌ Must run manually before starting the app

---

### Method 3: Manual SQL ⭐ For Advanced Users

Run SQL migrations directly:

```bash
# Using psql
psql $DATABASE_URL -f migrations/001_create_price_history.sql

# Or using make
make migrate-sql FILE=migrations/001_create_price_history.sql

# Rollback (down migration)
psql $DATABASE_URL -f migrations/001_create_price_history.down.sql
```

**Pros:**
- ✅ Maximum control
- ✅ Can modify on the fly
- ✅ Great for debugging

**Cons:**
- ❌ Manual tracking required
- ❌ Easy to make mistakes

---

### Method 4: Makefile Commands ⭐ Easiest for Daily Work

Use convenient Make commands:

```bash
make migrate-up       # Run migrations
make migrate-status   # Check status
make migrate-reset    # Reset database (with confirmation)
```

---

## Migration Commands Reference

### Check Migration Status
```bash
# Shows which tables exist and record counts
go run cmd/migrate/main.go -command=status
make migrate-status
```

### Apply Migrations
```bash
# Create tables and indexes
go run cmd/migrate/main.go -command=up
make migrate-up
```

### Reset Database
```bash
# Drops all tables and recreates them
# ⚠️  WARNING: Destroys all data!
go run cmd/migrate/main.go -command=reset
make migrate-reset
```

The reset command will ask for confirmation:
```
WARNING: This will drop all tables and recreate them!
Are you sure? Type 'yes' to continue:
```

---

## Current Database Schema

### price_history Table

```sql
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL,
    high BIGINT,
    high_time BIGINT,
    low BIGINT,
    low_time BIGINT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_item_timestamp ON price_history(item_id, timestamp);
CREATE INDEX idx_timestamp ON price_history(timestamp);
```

**Columns:**
- `id` - Auto-incrementing primary key
- `item_id` - OSRS item identifier (e.g., 4151 for Abyssal whip)
- `high` - High (buy) price in GP
- `high_time` - Unix timestamp when high price was recorded
- `low` - Low (sell) price in GP
- `low_time` - Unix timestamp when low price was recorded
- `timestamp` - When this price snapshot was taken
- `created_at` - When the row was inserted into the database

**Indexes:**
- `idx_item_timestamp` - Composite index for fast item-specific queries
- `idx_timestamp` - Time-range queries

---

## Production Deployment

### First Time Setup

```bash
# 1. Set environment variable
export DATABASE_URL="postgresql://user:pass@prod.example.com:5432/osrs_prices?sslmode=require"

# 2. Create database (if needed)
createdb osrs_prices

# 3. Run migrations
go run cmd/migrate/main.go -command=up

# 4. Verify
go run cmd/migrate/main.go -command=status

# 5. Start the application
./osrs-price-api
```

### Updates

When deploying updates with schema changes:

```bash
# 1. Backup database
pg_dump $DATABASE_URL > backup-$(date +%Y%m%d).sql

# 2. Run new migrations
go run cmd/migrate/main.go -command=up

# 3. Verify
go run cmd/migrate/main.go -command=status

# 4. Deploy application
```

---

## Troubleshooting

### "Table already exists"
This is normal! GORM's AutoMigrate is idempotent - it safely skips existing tables.

### "Failed to connect to database"
Check your connection string:
```bash
# Test connection
psql $DATABASE_URL -c "SELECT 1;"

# Check environment variable
echo $DATABASE_URL
```

### "Permission denied"
Ensure your database user has sufficient privileges:
```sql
GRANT ALL PRIVILEGES ON DATABASE osrs_prices TO your_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_user;
```

### Reset Everything
If you need a fresh start:
```bash
# Drop and recreate database
make db-drop
make db-create
make migrate-up
```

---

## Adding New Migrations

When you add new features that require schema changes:

1. **Update the Go model** in `internal/models/`
2. **Create SQL migration** in `migrations/`
3. **Test locally** with `make migrate-reset && make migrate-up`
4. **Update AutoMigrate** in `internal/database/connection.go`

Example:
```go
// Add new model
type ItemMetadata struct {
    ID     uint   `gorm:"primaryKey"`
    ItemID int    `gorm:"index"`
    Name   string
}

// Add to AutoMigrate
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.PriceHistory{},
        &models.ItemMetadata{}, // new
    )
}
```

---

## FAQ

**Q: Do migrations run automatically?**  
A: Yes, when you start the app with `go run main.go`. Use the migration tool if you want manual control.

**Q: Can I roll back migrations?**  
A: Currently, only by running the `.down.sql` files manually or using `migrate-reset` (which destroys all data).

**Q: What happens if a migration fails?**  
A: The app won't start. Check the error message and fix the schema issue.

**Q: Do I need to run migrations after every update?**  
A: If you use automatic migrations (default), no. Otherwise, run `migrate-up` after pulling new code.

**Q: Can I skip migrations?**  
A: Not recommended. Migrations ensure your database schema matches the application code.

---

## Best Practices

1. ✅ **Always backup production** before migrations
2. ✅ **Test migrations locally** before production
3. ✅ **Use `migrate-status`** to verify after deploying
4. ✅ **Keep SQL migration files** even if using AutoMigrate
5. ✅ **Document breaking changes** in migration files
6. ✅ **Use transactions** for complex data migrations

---

## Related Documentation

- [README.md](README.md) - General API documentation
- [QUICKSTART.md](QUICKSTART.md) - Getting started guide
- [migrations/README.md](migrations/README.md) - Migration file reference