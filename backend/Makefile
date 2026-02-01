.PHONY: build up down restart logs clean test migrate-up migrate-down

# Build the Docker images
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# Restart all services
restart: down up

# View logs
logs:
	docker-compose logs -f

# Clean up containers and volumes
clean:
	docker-compose down -v
	docker system prune -f

# Run migrations up
migrate-up:
	docker-compose exec backend ./migrate up

# Run migrations down
migrate-down:
	docker-compose exec backend ./migrate down

# Run tests (requires test database)
test:
	docker-compose exec backend go test ./...

# Development mode
dev-up:
	docker-compose -f docker-compose.dev.yml up -d

dev-down:
	docker-compose -f docker-compose.dev.yml down

dev-logs:
	docker-compose -f docker-compose.dev.yml logs -f

# Run migrations in dev mode
dev-migrate-up:
	docker-compose -f docker-compose.dev.yml exec backend go run cmd/migrate/main.go up

dev-migrate-down:
	docker-compose -f docker-compose.dev.yml exec backend go run cmd/migrate/main.go down
