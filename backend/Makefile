# Load environment variables from .env (if present)
-include .env
export $(shell sed 's/=.*//' .env 2>/dev/null)

MIGRATIONS_DIR ?= ./migrations

# Postgres connection string for golang-migrate
DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Create new migration file (with user input)
migratenew:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

# Run all up migrations
migrateup:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) up

# Roll back the last migration
migratedown:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) down 1

# Check current migration version
migrateversion:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) version

# Force to specific version
migrateforce:
	@read -p "Enter version to force: " ver; \
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) force $$ver

# Development
run: ## Run the application
	go run ./main.go