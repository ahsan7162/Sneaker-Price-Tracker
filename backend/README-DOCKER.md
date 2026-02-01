# Docker Setup Guide

This guide explains how to run the Sneaker Price Tracker backend using Docker.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+

## Quick Start

### Production Mode

1. **Build and start all services:**
   ```bash
   docker-compose up -d
   ```

2. **View logs:**
   ```bash
   docker-compose logs -f backend
   ```

3. **Stop services:**
   ```bash
   docker-compose down
   ```

### Development Mode

For development with hot-reload and volume mounting:

```bash
docker-compose -f docker-compose.dev.yml up -d
```

## Available Commands

Using the Makefile (recommended):

```bash
# Build images
make build

# Start services
make up

# Stop services
make down

# View logs
make logs

# Run migrations manually
make migrate-up
make migrate-down

# Clean up everything
make clean

# Development mode
make dev-up
make dev-down
make dev-logs
make dev-migrate-up
make dev-migrate-down
```

Or using docker-compose directly:

```bash
# Build
docker-compose build

# Start
docker-compose up -d

# Stop
docker-compose down

# View logs
docker-compose logs -f

# Run migrations
docker-compose exec backend ./migrate up
docker-compose exec backend ./migrate down
```

## Services

### PostgreSQL Database
- **Container**: `sneaker-tracker-db`
- **Port**: `5432`
- **Database**: `sneaker_tracker`
- **User**: `postgres`
- **Password**: `postgres`
- **Data**: Persisted in Docker volume `postgres_data`

### Backend
- **Container**: `sneaker-tracker-backend`
- **Port**: `8080` (exposed for future API server)
- **Environment Variables**: Configured via docker-compose.yml

## Environment Variables

You can customize the database connection by setting environment variables:

```bash
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=sneaker_tracker
DB_SSLMODE=disable
```

Or modify `docker-compose.yml`:

```yaml
environment:
  DB_HOST: postgres
  DB_PORT: 5432
  DB_USER: your_user
  DB_PASSWORD: your_password
  DB_NAME: sneaker_tracker
  DB_SSLMODE: disable
```

## Database Access

### From Host Machine

```bash
# Using psql
psql -h localhost -p 5432 -U postgres -d sneaker_tracker

# Password: postgres
```

### From Another Container

```bash
docker-compose exec postgres psql -U postgres -d sneaker_tracker
```

## Running Migrations

Migrations run automatically when the backend container starts. To run manually:

```bash
# Run migrations up
docker-compose exec backend ./migrate up

# Rollback migrations
docker-compose exec backend ./migrate down
```

## Troubleshooting

### Container won't start

1. Check logs:
   ```bash
   docker-compose logs backend
   ```

2. Verify PostgreSQL is healthy:
   ```bash
   docker-compose ps
   ```

3. Check database connection:
   ```bash
   docker-compose exec postgres pg_isready -U postgres
   ```

### Migrations fail

1. Check if database exists:
   ```bash
   docker-compose exec postgres psql -U postgres -l
   ```

2. Manually create database if needed:
   ```bash
   docker-compose exec postgres psql -U postgres -c "CREATE DATABASE sneaker_tracker;"
   ```

### Reset Everything

```bash
# Stop and remove containers, networks, and volumes
docker-compose down -v

# Remove images
docker-compose down --rmi all

# Start fresh
docker-compose up -d
```

## Data Persistence

Database data is stored in a Docker volume named `postgres_data`. To backup:

```bash
# Backup
docker-compose exec postgres pg_dump -U postgres sneaker_tracker > backup.sql

# Restore
docker-compose exec -T postgres psql -U postgres sneaker_tracker < backup.sql
```

## Network

Both services run on the `sneaker-network` bridge network, allowing them to communicate using service names (e.g., `postgres` instead of `localhost`).

## Production Considerations

For production deployment:

1. **Change default passwords** in `docker-compose.yml`
2. **Use secrets management** for sensitive data
3. **Enable SSL** for database connections (`DB_SSLMODE=require`)
4. **Set resource limits** for containers
5. **Use a managed database** service instead of containerized PostgreSQL
6. **Set up proper logging** and monitoring
7. **Use environment-specific compose files**
