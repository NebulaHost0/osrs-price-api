# Database Migrations

This directory contains SQL migration files for the OSRS Price API database.

## Migration Files

### 001_create_price_history.sql
Creates the main `price_history` table with:
- `id`: Primary key
- `item_id`: OSRS item identifier
- `high`: High (buy) price
- `high_time`: Timestamp of high price
- `low`: Low (sell) price  
- `low_time`: Timestamp of low price
- `timestamp`: When the price was recorded
- `created_at`: Row creation timestamp

**Indexes:**
- `idx_item_timestamp`: Composite index on (item_id, timestamp) for fast item-specific queries
- `idx_timestamp`: Index on timestamp for time-range queries

## Running Migrations

### Automatic Migration (Recommended for Development)
Migrations run automatically when you start the application:
```bash
go run main.go
```

### Manual Migration
Use the migration tool for more control:

```bash
# Run migrations
go run cmd/migrate/main.go -command=up

# Check migration status
go run cmd/migrate/main.go -command=status

# Reset database (WARNING: destroys all data!)
go run cmd/migrate/main.go -command=reset
```

### Using Makefile
```bash
# Run migrations
make migrate-up

# Check status
make migrate-status

# Reset database
make migrate-reset
```

### Direct SQL Execution
You can also run migrations manually:
```bash
psql $DATABASE_URL -f migrations/001_create_price_history.sql
```

## Creating New Migrations

When creating new migrations, follow this naming convention:
```
XXX_description.sql       # Up migration
XXX_description.down.sql  # Down migration (rollback)
```

Where `XXX` is a zero-padded sequential number (001, 002, etc.).

Example:
```
002_add_item_names.sql
002_add_item_names.down.sql
```

## Migration Best Practices

1. **Always create both up and down migrations**
2. **Test migrations in development first**
3. **Backup production data before migrating**
4. **Make migrations idempotent** (safe to run multiple times)
5. **Use transactions for data migrations**
6. **Document breaking changes**

## Current Schema Version

**Version**: 001
**Last Updated**: 2025-10-05
**Tables**: price_history