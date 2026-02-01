#!/bin/sh

set -e

echo "Waiting for PostgreSQL to be ready..."

# Wait for PostgreSQL to be ready
until pg_isready -h "${DB_HOST:-postgres}" -p "${DB_PORT:-5432}" -U "${DB_USER:-postgres}"; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is up - executing migrations"

# Run migrations
./migrate "${1:-up}"

echo "Migrations completed"

# Keep container running if needed
if [ "$2" = "keep-alive" ]; then
  echo "Keeping container alive..."
  tail -f /dev/null
fi
