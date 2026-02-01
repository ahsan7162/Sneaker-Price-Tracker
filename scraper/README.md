# Sneaker Price Tracker Scraper

A FastAPI-based web scraping service for extracting product information from Nike and Adidas product pages. The service extracts product variants (color, size, unique identifier), price history (price, stock status), and product images.

## Features

- **Nike Scraper**: Extracts product data from Nike product pages
- **Adidas Scraper**: Extracts product data from Adidas product pages
- **Generalized Architecture**: Easy to extend with additional brand scrapers
- **FastAPI REST API**: Clean REST endpoints for scraping operations
- **Docker Support**: Containerized deployment with Docker
- **Product Images**: Extracts all product images from pages
- **Comprehensive Data**: Extracts variants, prices, and stock information

## Requirements

- Python 3.14 (or Python 3.13/3.12 if 3.14 is not available)
- Docker and Docker Compose (optional, for containerized deployment)

## Installation

### Option 1: Docker (Recommended)

1. **Build and run with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

2. **View logs:**
   ```bash
   docker-compose logs -f scraper
   ```

3. **Stop the service:**
   ```bash
   docker-compose down
   ```

### Option 2: Local Development

1. **Create a virtual environment:**
   ```bash
   python3.14 -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```

2. **Install dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

3. **Run the application:**
   ```bash
   uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
   ```

## API Endpoints

### Root Endpoint

- **GET** `/` - API information and available endpoints

### Health Check

- **GET** `/health` - Health check endpoint

### Scrape Nike Product

- **POST** `/scrape/nike`
- **Request Body:**
  ```json
  {
    "url": "https://www.nike.com/t/pegasus-41-mens-road-running-shoes-LMhfRGdO/IM6674-101"
  }
  ```
- **Response:**
  ```json
  {
    "variants": [
      {
        "color": "White/White/Hyper Pink/Black",
        "shoe_size": "10",
        "unique_identifier": "IM6674-101"
      }
    ],
    "price_history": [
      {
        "price": 145.00,
        "is_in_stock": true
      }
    ],
    "images": [
      "https://static.nike.com/a/images/t_PDP_1728_v1/..."
    ]
  }
  ```

### Scrape Adidas Product

- **POST** `/scrape/adidas`
- **Request Body:**
  ```json
  {
    "url": "https://www.adidas.com/us/adizero-evo-sl-shoes/KJ1363.html"
  }
  ```
- **Response:** Same format as Nike endpoint

## Usage Examples

### Using cURL

**Scrape Nike product:**
```bash
curl -X POST "http://localhost:8000/scrape/nike" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.nike.com/t/pegasus-41-mens-road-running-shoes-LMhfRGdO/IM6674-101"}'
```

**Scrape Adidas product:**
```bash
curl -X POST "http://localhost:8000/scrape/adidas" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.adidas.com/us/adizero-evo-sl-shoes/KJ1363.html"}'
```

### Using Python

```python
import httpx

async def scrape_nike():
    async with httpx.AsyncClient() as client:
        response = await client.post(
            "http://localhost:8000/scrape/nike",
            json={
                "url": "https://www.nike.com/t/pegasus-41-mens-road-running-shoes-LMhfRGdO/IM6674-101"
            }
        )
        return response.json()

# Run with asyncio
import asyncio
result = asyncio.run(scrape_nike())
print(result)
```

### Using JavaScript/Node.js

```javascript
const fetch = require('node-fetch');

async function scrapeNike() {
  const response = await fetch('http://localhost:8000/scrape/nike', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      url: 'https://www.nike.com/t/pegasus-41-mens-road-running-shoes-LMhfRGdO/IM6674-101'
    })
  });
  return await response.json();
}

scrapeNike().then(console.log);
```

## API Documentation

When the service is running, interactive API documentation is available at:

- **Swagger UI**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc

## Project Structure

```
scraper/
├── app/
│   ├── __init__.py
│   ├── main.py                 # FastAPI application
│   ├── models/
│   │   ├── __init__.py
│   │   ├── product_variant.py  # Product variant models
│   │   └── price_history.py     # Price history models
│   ├── scrapers/
│   │   ├── __init__.py
│   │   ├── base.py             # Base scraper class
│   │   ├── nike.py             # Nike scraper
│   │   └── adidas.py            # Adidas scraper
│   └── utils/
│       ├── __init__.py
│       └── url_validator.py     # URL validation utilities
├── requirements.txt
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```

## Architecture

The scraper service uses a generalized architecture:

1. **Base Scraper** (`BaseScraper`): Abstract base class defining the scraping interface
2. **Brand Scrapers**: Implementations for each brand (Nike, Adidas)
3. **Models**: Pydantic models matching the Go backend structs
4. **API Layer**: FastAPI endpoints for each brand

### Adding a New Brand Scraper

To add a new brand scraper:

1. Create a new scraper class in `app/scrapers/` that extends `BaseScraper`
2. Implement the abstract methods:
   - `_parse_variants()`: Extract product variants
   - `_parse_price_history()`: Extract price and stock information
   - `_parse_images()`: Extract product images
3. Add a new endpoint in `app/main.py`
4. Update URL validator if needed

## Error Handling

The API returns appropriate HTTP status codes:

- **200**: Success
- **400**: Bad Request (invalid URL or wrong brand)
- **500**: Internal Server Error (scraping failure)

## Limitations

- Websites may change their HTML structure, requiring scraper updates
- Some pages may require JavaScript rendering (not currently supported)
- Rate limiting may apply from target websites
- Stock status may not always be accurate

## Development

### Running Tests

```bash
# Install test dependencies
pip install pytest pytest-asyncio httpx

# Run tests
pytest
```

### Code Style

The project follows PEP 8 style guidelines. Consider using:
- `black` for code formatting
- `flake8` for linting
- `mypy` for type checking

## Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
HOST=0.0.0.0
PORT=8000
SCRAPER_TIMEOUT=30
```

## Notes

- **Python Version**: The Dockerfile uses Python 3.14. If this version is not available, update the Dockerfile to use `python:3.13-slim` or `python:3.12-slim`
- **Image URLs**: Extracted image URLs are normalized to full URLs
- **Variants**: The scraper attempts to extract all color/size combinations, but may fall back to default values if parsing fails
- **Price**: Prices are extracted as floats, currency symbols are removed

## License

This project is part of the Sneaker Price Tracker system.
