package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"osrs-price-api/internal/api"
	"osrs-price-api/internal/cache"
	"osrs-price-api/internal/database"
	"osrs-price-api/internal/osrs"
	"osrs-price-api/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (optional, won't fail if not found)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	// Load database configuration
	dbConfig := database.LoadConfig()

	// Connect to database
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repository
	repo := database.NewRepository(db)

	// Initialize cache
	priceCache := cache.NewPriceCache()

	// Initialize OSRS client
	osrsClient := osrs.NewClient()

	// Start background worker for periodic price fetching
	// Fetch prices every 5 minutes (aligned with OSRS Wiki update frequency)
	priceFetcher := worker.NewPriceFetcher(osrsClient, repo, 5*time.Minute)
	priceFetcher.Start()

	// Start cleanup worker to manage database size
	// Runs daily at 3 AM to delete old data and keep costs down
	cleanupWorker := worker.NewCleanupWorker(repo, 24*time.Hour)
	cleanupWorker.Start()

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware to allow frontend access
	router.Use(api.CORSMiddleware())

	// Setup API routes
	apiHandler := api.NewHandler(osrsClient, priceCache, repo)
	api.SetupRoutes(router, apiHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Check for SSL certificate files
	certFile := os.Getenv("SSL_CERT_FILE")
	keyFile := os.Getenv("SSL_KEY_FILE")

	// Setup graceful shutdown
	go func() {
		if certFile != "" && keyFile != "" {
			log.Printf("Starting OSRS Price API server with TLS on port %s", port)
			log.Printf("Using certificate: %s", certFile)
			if err := router.RunTLS(":"+port, certFile, keyFile); err != nil {
				log.Fatalf("Failed to start TLS server: %v", err)
			}
		} else {
			log.Printf("Starting OSRS Price API server (HTTP) on port %s", port)
			log.Println("Note: Set SSL_CERT_FILE and SSL_KEY_FILE for HTTPS")
			if err := router.Run(":" + port); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	priceFetcher.Stop()
	cleanupWorker.Stop()
	log.Println("Server stopped")
}
