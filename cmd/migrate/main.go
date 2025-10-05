package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"osrs-price-api/internal/database"
	
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	var command string
	flag.StringVar(&command, "command", "up", "Migration command: up, down, status, create")
	flag.Parse()

	// Load database configuration
	dbConfig := database.LoadConfig()

	// Connect to database
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	switch command {
	case "up":
		log.Println("Running migrations...")
		if err := database.AutoMigrate(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("✓ Migrations completed successfully")

	case "status":
		log.Println("Checking migration status...")
		if err := database.MigrationStatus(db); err != nil {
			log.Fatalf("Failed to check status: %v", err)
		}

	case "reset":
		log.Println("WARNING: This will drop all tables and recreate them!")
		fmt.Print("Are you sure? Type 'yes' to continue: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			log.Println("Reset cancelled")
			os.Exit(0)
		}
		
		if err := database.ResetDatabase(db); err != nil {
			log.Fatalf("Reset failed: %v", err)
		}
		log.Println("✓ Database reset completed successfully")

	default:
		log.Fatalf("Unknown command: %s (available: up, status, reset)", command)
	}
}