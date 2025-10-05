.PHONY: help run build test clean docker-up docker-down db-create db-drop migrate-up migrate-status migrate-reset

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	go run main.go

build: ## Build the application
	go build -o bin/osrs-price-api main.go

build-migrate: ## Build the migration tool
	go build -o bin/migrate cmd/migrate/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

docker-up: ## Start PostgreSQL with Docker Compose
	docker-compose up -d

docker-down: ## Stop PostgreSQL
	docker-compose down

docker-logs: ## View PostgreSQL logs
	docker-compose logs -f postgres

db-create: ## Create the database
	@echo "Creating database..."
	@psql -U postgres -h localhost -c "CREATE DATABASE osrs_prices;" || echo "Database may already exist"

db-drop: ## Drop the database (WARNING: This will delete all data!)
	@echo "Dropping database..."
	@psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS osrs_prices;"

migrate-up: ## Run database migrations
	go run cmd/migrate/main.go -command=up

migrate-status: ## Check migration status
	go run cmd/migrate/main.go -command=status

migrate-reset: ## Reset database (WARNING: This will delete all data!)
	go run cmd/migrate/main.go -command=reset

migrate-sql: ## Run SQL migration manually
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make migrate-sql FILE=migrations/001_create_price_history.sql"; \
		exit 1; \
	fi
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL environment variable not set"; \
		exit 1; \
	fi
	@echo "Running migration: $(FILE)"
	@psql $(DATABASE_URL) -f $(FILE)

tidy: ## Tidy up Go modules
	go mod tidy

deps: ## Download dependencies
	go mod download

dev: docker-up ## Start development environment (database + app)
	@echo "Waiting for database to be ready..."
	@sleep 3
	@make migrate-up
	@make run

setup: docker-up ## Complete setup (database + migrations)
	@echo "Setting up development environment..."
	@sleep 3
	@make migrate-up
	@echo "✓ Setup complete! Run 'make run' to start the server"