"""FastAPI application for product scraping."""

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field, HttpUrl

from app.models.product_variant import ProductVariantResponse
from app.models.price_history import PriceHistoryResponse
from app.scrapers.nike import NikeScraper
from app.scrapers.adidas import AdidasScraper
from app.utils.url_validator import is_nike_url, is_adidas_url, validate_url


app = FastAPI(
    title="Sneaker Price Tracker Scraper",
    description="API for scraping Nike and Adidas product pages",
    version="1.0.0"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class ScrapeRequest(BaseModel):
    """Request model for scrape endpoints."""
    
    url: HttpUrl = Field(..., description="Product page URL to scrape")
    
    class Config:
        json_schema_extra = {
            "example": {
                "url": "https://www.nike.com/t/pegasus-41-mens-road-running-shoes-LMhfRGdO/IM6674-101"
            }
        }


class ScrapeResponse(BaseModel):
    """Response model for scrape endpoints."""
    
    variants: list[ProductVariantResponse] = Field(..., description="List of product variants")
    price_history: list[PriceHistoryResponse] = Field(..., description="List of price history entries")
    images: list[str] = Field(..., description="List of product image URLs")
    
    class Config:
        json_schema_extra = {
            "example": {
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
                        "is_in_stock": True
                    }
                ],
                "images": [
                    "https://static.nike.com/a/images/t_PDP_1728_v1/..."
                ]
            }
        }


@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "message": "Sneaker Price Tracker Scraper API",
        "version": "1.0.0",
        "endpoints": {
            "nike": "/scrape/nike",
            "adidas": "/scrape/adidas",
            "health": "/health"
        }
    }


@app.get("/health")
async def health():
    """Health check endpoint."""
    return {"status": "healthy"}


@app.post("/scrape/nike", response_model=ScrapeResponse)
async def scrape_nike(request: ScrapeRequest):
    """
    Scrape Nike product page.
    
    Args:
        request: ScrapeRequest containing the Nike product URL
        
    Returns:
        ScrapeResponse with variants, price history, and images
        
    Raises:
        HTTPException: If URL is invalid or scraping fails
    """
    url = str(request.url)
    
    # Validate URL
    if not validate_url(url):
        raise HTTPException(status_code=400, detail="Invalid URL format")
    
    # Verify it's a Nike URL
    if not is_nike_url(url):
        raise HTTPException(
            status_code=400,
            detail="URL does not appear to be a Nike product page"
        )
    
    try:
        scraper = NikeScraper()
        result = await scraper.scrape(url)
        
        return ScrapeResponse(
            variants=result.variants,
            price_history=result.price_history,
            images=result.images
        )
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Failed to scrape Nike page: {str(e)}"
        )


@app.post("/scrape/adidas", response_model=ScrapeResponse)
async def scrape_adidas(request: ScrapeRequest):
    """
    Scrape Adidas product page.
    
    Args:
        request: ScrapeRequest containing the Adidas product URL
        
    Returns:
        ScrapeResponse with variants, price history, and images
        
    Raises:
        HTTPException: If URL is invalid or scraping fails
    """
    url = str(request.url)
    
    # Validate URL
    if not validate_url(url):
        raise HTTPException(status_code=400, detail="Invalid URL format")
    
    # Verify it's an Adidas URL
    if not is_adidas_url(url):
        raise HTTPException(
            status_code=400,
            detail="URL does not appear to be an Adidas product page"
        )
    
    try:
        scraper = AdidasScraper()
        result = await scraper.scrape(url)
        
        return ScrapeResponse(
            variants=result.variants,
            price_history=result.price_history,
            images=result.images
        )
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Failed to scrape Adidas page: {str(e)}"
        )


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
