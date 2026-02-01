"""Base scraper interface for product scraping."""

from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import List
import httpx
from bs4 import BeautifulSoup

from app.models.product_variant import ProductVariantResponse
from app.models.price_history import PriceHistoryResponse


@dataclass
class ScrapeResult:
    """Result container for scraped data."""
    
    variants: List[ProductVariantResponse]
    price_history: List[PriceHistoryResponse]
    images: List[str]


class BaseScraper(ABC):
    """Abstract base class for product scrapers."""
    
    def __init__(self, timeout: int = 30):
        """Initialize the scraper with HTTP client."""
        self.timeout = timeout
        self.headers = {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
            "Accept-Language": "en-US,en;q=0.5",
            "Accept-Encoding": "gzip, deflate, br",
            "Connection": "keep-alive",
            "Upgrade-Insecure-Requests": "1",
        }
    
    async def scrape(self, url: str) -> ScrapeResult:
        """
        Main scraping method that orchestrates the scraping process.
        
        Args:
            url: Product page URL to scrape
            
        Returns:
            ScrapeResult containing variants, price history, and images
        """
        html = await self._fetch_page(url)
        variants = await self._parse_variants(html, url)
        price_history = await self._parse_price_history(html, url)
        images = await self._parse_images(html, url)
        
        return ScrapeResult(
            variants=variants,
            price_history=price_history,
            images=images
        )
    
    async def _fetch_page(self, url: str) -> str:
        """
        Fetch HTML content from the given URL.
        
        Args:
            url: URL to fetch
            
        Returns:
            HTML content as string
            
        Raises:
            httpx.HTTPError: If request fails
        """
        async with httpx.AsyncClient(timeout=self.timeout, follow_redirects=True) as client:
            response = await client.get(url, headers=self.headers)
            response.raise_for_status()
            return response.text
    
    @abstractmethod
    async def _parse_variants(self, html: str, url: str) -> List[ProductVariantResponse]:
        """
        Parse product variants from HTML.
        
        Args:
            html: HTML content of the product page
            url: Original URL for context
            
        Returns:
            List of product variants
        """
        pass
    
    @abstractmethod
    async def _parse_price_history(self, html: str, url: str) -> List[PriceHistoryResponse]:
        """
        Parse price history data from HTML.
        
        Args:
            html: HTML content of the product page
            url: Original URL for context
            
        Returns:
            List of price history entries
        """
        pass
    
    @abstractmethod
    async def _parse_images(self, html: str, url: str) -> List[str]:
        """
        Parse product images from HTML.
        
        Args:
            html: HTML content of the product page
            url: Original URL for context
            
        Returns:
            List of image URLs
        """
        pass
    
    def _parse_html(self, html: str) -> BeautifulSoup:
        """
        Parse HTML string into BeautifulSoup object.
        
        Args:
            html: HTML content string
            
        Returns:
            BeautifulSoup object
        """
        return BeautifulSoup(html, "lxml")
