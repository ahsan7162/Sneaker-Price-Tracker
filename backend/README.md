# Sneaker Price Tracker - Backend

Go backend for the Sneaker Price Tracker application with PostgreSQL database support.

## Project Structure

```
backend/
├── cmd/
│   └── migrate/
│       └── main.go          # Migration runner
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── db/
│   │   └── db.go            # Database connection
│   ├── migrations/
│   │   ├── migrate.go       # Migration utilities
│   │   ├── 000001_init_schema.up.sql
│   │   └── 000001_init_schema.down.sql
│   ├── models/
│   │   ├── product.go
│   │   ├── country.go
│   │   ├── product_variant.go
│   │   └── price_history.go
│   └── repositories/
│       ├── product_repository.go
│       ├── product_repository_test.go
│       ├── country_repository.go
│       ├── country_repository_test.go
│       ├── product_variant_repository.go
│       ├── product_variant_repository_test.go
│       ├── price_history_repository.go
│       └── price_history_repository_test.go
├── go.mod
└── README.md
```

## Database Schema

### Tables

1. **products** - Stores general shoe model information
   - id (Primary Key)
   - brand_name (String)
   - shoe_name (String)
   - base_url (Text)
   - created_at (Timestamp)

2. **countries** - Lookup table for markets
   - id (Primary Key)
   - country_code (String, Unique)
   - currency (String)

3. **product_variants** - Maps specific Color and Size combinations
   - id (Primary Key)
   - product_id (Foreign Key → products.id)
   - color (String)
   - shoe_size (String)
   - unique_identifier (Composite unique key: product_id, color, shoe_size)

4. **price_history** - Time-series table for price checks
   - id (Primary Key)
   - variant_id (Foreign Key → product_variants.id)
   - country_id (Foreign Key → countries.id)
   - price (Numeric(10, 2))
   - is_in_stock (Boolean)
   - captured_at (Timestamp)

## Setup

### Option 1: Docker (Recommended)

The easiest way to run the backend is using Docker Compose:

```bash
# Build and start all services (PostgreSQL + Backend)
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down
```

Migrations run automatically when the backend container starts.

For detailed Docker instructions, see [README-DOCKER.md](./README-DOCKER.md).

### Option 2: Local Development

#### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher

#### Installation

1. Install dependencies:
```bash
go mod download
```

2. Set up environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=sneaker_tracker
export DB_SSLMODE=disable
```

Or create a `.env` file (not tracked in git):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=sneaker_tracker
DB_SSLMODE=disable
```

3. Create the database:
```sql
CREATE DATABASE sneaker_tracker;
```

4. Run migrations:
```bash
go run cmd/migrate/main.go up
```

To rollback migrations:
```bash
go run cmd/migrate/main.go down
```

## Running Tests

### Local Development

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run tests for a specific package:
```bash
go test ./internal/repositories/...
```

**Note**: Tests require a test database. Set `DB_NAME=sneaker_tracker_test` for running tests.

### Docker

Run tests in Docker:
```bash
docker-compose exec backend go test ./...
```

Or using Makefile:
```bash
make test
```

## CRUD Operations

The backend provides complete CRUD operations for all tables through repository pattern:

### Products
- `Create(req CreateProductRequest) (*Product, error)`
- `GetByID(id int64) (*Product, error)`
- `GetAll() ([]Product, error)`
- `Update(id int64, req UpdateProductRequest) (*Product, error)`
- `Delete(id int64) error`

### Countries
- `Create(req CreateCountryRequest) (*Country, error)`
- `GetByID(id int64) (*Country, error)`
- `GetByCountryCode(countryCode string) (*Country, error)`
- `GetAll() ([]Country, error)`
- `Update(id int64, req UpdateCountryRequest) (*Country, error)`
- `Delete(id int64) error`

### Product Variants
- `Create(req CreateProductVariantRequest) (*ProductVariant, error)`
- `GetByID(id int64) (*ProductVariant, error)`
- `GetByProductID(productID int64) ([]ProductVariant, error)`
- `GetAll() ([]ProductVariant, error)`
- `Update(id int64, req UpdateProductVariantRequest) (*ProductVariant, error)`
- `Delete(id int64) error`

### Price History
- `Create(req CreatePriceHistoryRequest) (*PriceHistory, error)`
- `GetByID(id int64) (*PriceHistory, error)`
- `GetByVariantID(variantID int64) ([]PriceHistory, error)`
- `GetByVariantAndCountry(variantID, countryID int64) ([]PriceHistory, error)`
- `GetLatestPrice(variantID, countryID int64) (*PriceHistory, error)`
- `GetPriceHistoryByDateRange(variantID, countryID int64, startDate, endDate time.Time) ([]PriceHistory, error)`
- `GetAll() ([]PriceHistory, error)`
- `Update(id int64, req UpdatePriceHistoryRequest) (*PriceHistory, error)`
- `Delete(id int64) error`

## Dependencies

- `github.com/lib/pq` - PostgreSQL driver
- `github.com/golang-migrate/migrate/v4` - Database migration tool

## License

MIT
