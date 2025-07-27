include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build run test clean docker-build docker-up docker-down docker-logs help

# Help command
help:
	@echo "Available commands:"
	@echo ""
	@echo "Go commands:"
	@echo "  build    - Build the Go application"
	@echo "  run      - Run the application locally"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build - Build Docker containers"
	@echo "  docker-up    - Start services"
	@echo "  docker-down  - Stop services"
	@echo "  docker-logs  - View logs"
	@echo "  docker-clean - Clean up containers and volumes"
	@echo ""
	@echo "Database commands:"
	@echo "  db-migrate       - Check database status"
	@echo "  db-migrate-manual - Run manual migration"
	@echo "  db-connect       - Connect to database"
	@echo "  db-reset         - Reset database (deletes all data)"
	@echo "  db-status        - Show database tables"
	@echo "  db-seed          - Seed database with sample data"
	@echo ""
	@echo "Development:"
	@echo "  dev  - Start development environment"
	@echo "  prod - Start production environment"

# Go commands
build:
	go build -o bin/main ./cmd/main.go

run:
	go run ./cmd/main.go

test:
	go test ./...

clean:
	rm -rf bin/

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v --remove-orphans

# Database commands
db-migrate:
	@echo "Running database migrations..."
	docker-compose exec postgres psql -U wschatuser -d wschatdb -c "SELECT 'Database is ready' as status;"

db-migrate-manual:
	@echo "Running manual migration..."
	docker-compose exec postgres psql -U wschatuser -d wschatdb -f /docker-entrypoint-initdb.d/01-init.sql

db-connect:
	docker-compose exec postgres psql -U wschatuser -d wschatdb

db-reset:
	@echo "Resetting database (this will delete all data)..."
	docker-compose down -v
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 10
	@echo "Database reset complete"

db-status:
	@echo "Checking database status..."
	docker-compose exec postgres psql -U wschatuser -d wschatdb -c "\dt"

db-seed:
	@echo "Seeding database with sample data..."
	docker-compose exec postgres psql -U wschatuser -d wschatdb -c "INSERT INTO users (username, password_hash) VALUES ('admin', '\$$2a\$$10\$$example_hash_here') ON CONFLICT (username) DO NOTHING;"

# Development
dev: docker-up
	@echo "Starting development environment..."
	@echo "PostgreSQL: localhost:5432"
	@echo "Application: http://localhost:9080"
	@echo "Use 'make docker-logs' to view logs"

# Production
prod: docker-build docker-up
	@echo "Production environment started"