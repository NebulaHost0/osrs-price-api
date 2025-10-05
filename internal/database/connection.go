package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"osrs-price-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	ConnectionString string
}

// LoadConfig loads database configuration from environment variables
func LoadConfig() *Config {
	// Check for connection string first
	connStr := os.Getenv("DATABASE_URL")
	
	// If no connection string, build from individual parameters (fallback)
	if connStr == "" {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		dbName := getEnv("DB_NAME", "osrs_prices")
		sslMode := getEnv("DB_SSLMODE", "disable")
		
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbName, sslMode,
		)
	}
	
	return &Config{
		ConnectionString: connStr,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Connect establishes a database connection
func Connect(config *Config) (*gorm.DB, error) {
	dsn := config.ConnectionString

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")
	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	
	if err := db.AutoMigrate(&models.PriceHistory{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// MigrationStatus checks the current migration status
func MigrationStatus(db *gorm.DB) error {
	log.Println("Checking database tables...")
	
	// Check if price_history table exists
	if db.Migrator().HasTable(&models.PriceHistory{}) {
		log.Println("✓ Table 'price_history' exists")
		
		// Get table info
		var count int64
		db.Model(&models.PriceHistory{}).Count(&count)
		log.Printf("  - Record count: %d", count)
		
		// Check indexes
		if db.Migrator().HasIndex(&models.PriceHistory{}, "idx_item_timestamp") {
			log.Println("  - Index 'idx_item_timestamp' exists")
		} else {
			log.Println("  - Index 'idx_item_timestamp' missing")
		}
	} else {
		log.Println("✗ Table 'price_history' does not exist")
		log.Println("  Run migrations with: go run cmd/migrate/main.go -command=up")
	}
	
	return nil
}

// ResetDatabase drops all tables and recreates them
func ResetDatabase(db *gorm.DB) error {
	log.Println("Dropping all tables...")
	
	if err := db.Migrator().DropTable(&models.PriceHistory{}); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}
	
	log.Println("Recreating tables...")
	if err := AutoMigrate(db); err != nil {
		return err
	}
	
	return nil
}