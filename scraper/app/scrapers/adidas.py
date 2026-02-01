"""Adidas product scraper implementation."""

import json
import re
from typing import List
from urllib.parse import urlparse

from app.models.product_variant import ProductVariantResponse
from app.models.price_history import PriceHistoryResponse
from app.scrapers.base import BaseScraper


class AdidasScraper(BaseScraper):
    """Scraper for Adidas product pages."""
    
    def __init__(self, timeout: int = 30):
        """Initialize Adidas scraper."""
        super().__init__(timeout)
        self.base_url = "https://www.adidas.com"
    
    async def _parse_variants(self, html: str, url: str) -> List[ProductVariantResponse]:
        """
        Parse product variants from Adidas HTML.
        
        Extracts color and size combinations, along with unique identifiers.
        """
        soup = self._parse_html(html)
        variants = []
        
        # Try to extract from JSON-LD or embedded JSON data
        json_data = self._extract_json_data(html)
        
        if json_data:
            variants.extend(self._parse_variants_from_json(json_data))
        
        # Fallback to HTML parsing
        if not variants:
            variants.extend(self._parse_variants_from_html(soup, url))
        
        # If still no variants, create a default one from the page
        if not variants:
            variants.append(self._create_default_variant(soup, url))
        
        return variants
    
    async def _parse_price_history(self, html: str, url: str) -> List[PriceHistoryResponse]:
        """Parse price and stock information from Adidas HTML."""
        soup = self._parse_html(html)
        price_history = []
        
        # Try JSON data first
        json_data = self._extract_json_data(html)
        
        if json_data:
            price_history.extend(self._parse_price_from_json(json_data))
        
        # Fallback to HTML parsing
        if not price_history:
            price_history.extend(self._parse_price_from_html(soup))
        
        # If still no price, create default entry
        if not price_history:
            price_history.append(PriceHistoryResponse(price=0.0, is_in_stock=False))
        
        return price_history
    
    async def _parse_images(self, html: str, url: str) -> List[str]:
        """Parse product images from Adidas HTML."""
        soup = self._parse_html(html)
        images = []
        
        # Try JSON data first
        json_data = self._extract_json_data(html)
        
        if json_data:
            images.extend(self._parse_images_from_json(json_data))
        
        # Parse from HTML img tags and data attributes
        images.extend(self._parse_images_from_html(soup))
        
        # Remove duplicates and normalize URLs
        unique_images = list(dict.fromkeys(images))  # Preserves order
        normalized_images = [self._normalize_image_url(img) for img in unique_images if img]
        
        return normalized_images
    
    def _extract_json_data(self, html: str) -> dict:
        """Extract JSON-LD structured data or embedded JSON from HTML."""
        # Look for JSON-LD script tags
        json_ld_pattern = r'<script[^>]*type=["\']application/ld\+json["\'][^>]*>(.*?)</script>'
        matches = re.findall(json_ld_pattern, html, re.DOTALL | re.IGNORECASE)
        
        for match in matches:
            try:
                data = json.loads(match.strip())
                if isinstance(data, dict) and data.get("@type") == "Product":
                    return data
                elif isinstance(data, list):
                    for item in data:
                        if isinstance(item, dict) and item.get("@type") == "Product":
                            return item
            except json.JSONDecodeError:
                continue
        
        # Look for Adidas-specific data attributes or script tags
        # Adidas often uses data attributes or embedded JSON
        data_patterns = [
            r'window\.__PRELOADED_STATE__\s*=\s*({.*?});',
            r'window\.__INITIAL_STATE__\s*=\s*({.*?});',
            r'"product":\s*({.*?})',
            r'data-product=\'({.*?})\'',
        ]
        
        for pattern in data_patterns:
            match = re.search(pattern, html, re.DOTALL)
            if match:
                try:
                    data = json.loads(match.group(1))
                    if isinstance(data, dict):
                        return data
                except json.JSONDecodeError:
                    continue
        
        return {}
    
    def _parse_variants_from_json(self, json_data: dict) -> List[ProductVariantResponse]:
        """Parse variants from JSON data."""
        variants = []
        
        # Extract product code/SKU
        sku = (
            json_data.get("sku") or
            json_data.get("productID") or
            json_data.get("mpn") or
            json_data.get("product", {}).get("articleNumber", "") if isinstance(json_data.get("product"), dict) else ""
        )
        
        # Extract product data
        product_data = json_data.get("product", {}) if isinstance(json_data.get("product"), dict) else {}
        
        colors = self._extract_colors_from_data(json_data, product_data)
        sizes = self._extract_sizes_from_data(json_data, product_data)
        
        # Create variants for each color/size combination
        if colors and sizes:
            for color in colors:
                for size in sizes:
                    variants.append(ProductVariantResponse(
                        color=color,
                        shoe_size=str(size),
                        unique_identifier=sku or f"{color}-{size}"
                    ))
        elif colors:
            for color in colors:
                variants.append(ProductVariantResponse(
                    color=color,
                    shoe_size="N/A",
                    unique_identifier=sku or color
                ))
        elif sizes:
            for size in sizes:
                variants.append(ProductVariantResponse(
                    color="N/A",
                    shoe_size=str(size),
                    unique_identifier=sku or str(size)
                ))
        
        return variants
    
    def _parse_variants_from_html(self, soup, url: str) -> List[ProductVariantResponse]:
        """Parse variants from HTML structure."""
        variants = []
        
        # Extract product code from URL (e.g., KJ1363)
        product_code = self._extract_product_code_from_url(url)
        
        # Find color options
        colors = self._extract_colors_from_html(soup)
        
        # Find size options
        sizes = self._extract_sizes_from_html(soup)
        
        # Create variants
        if colors and sizes:
            for color in colors:
                for size in sizes:
                    variants.append(ProductVariantResponse(
                        color=color,
                        shoe_size=str(size),
                        unique_identifier=product_code or f"{color}-{size}"
                    ))
        elif colors:
            for color in colors:
                variants.append(ProductVariantResponse(
                    color=color,
                    shoe_size="N/A",
                    unique_identifier=product_code or color
                ))
        elif sizes:
            for size in sizes:
                variants.append(ProductVariantResponse(
                    color="N/A",
                    shoe_size=str(size),
                    unique_identifier=product_code or str(size)
                ))
        
        return variants
    
    def _create_default_variant(self, soup, url: str) -> ProductVariantResponse:
        """Create a default variant when parsing fails."""
        product_code = self._extract_product_code_from_url(url)
        
        # Try to get product name/title
        title_elem = soup.find("h1") or soup.find("title")
        color = "N/A"
        if title_elem:
            title_text = title_elem.get_text(strip=True)
            # Try to extract color from title
            color_match = re.search(r'-\s*([^-]+)$', title_text)
            if color_match:
                color = color_match.group(1).strip()
        
        return ProductVariantResponse(
            color=color,
            shoe_size="N/A",
            unique_identifier=product_code or "UNKNOWN"
        )
    
    def _extract_product_code_from_url(self, url: str) -> str:
        """Extract product code from Adidas URL."""
        # URL format: https://www.adidas.com/{country}/{slug}/{product-code}.html
        parts = url.rstrip("/").split("/")
        if len(parts) >= 2:
            last_part = parts[-1]
            # Remove .html extension
            product_code = last_part.replace(".html", "")
            # Product codes are typically alphanumeric
            if re.match(r'^[A-Z0-9-]+$', product_code):
                return product_code
        return ""
    
    def _extract_colors_from_data(self, json_data: dict, product_data: dict) -> List[str]:
        """Extract color options from JSON data."""
        colors = []
        
        # Check various possible locations
        color_sources = [
            json_data.get("color"),
            product_data.get("color"),
            product_data.get("colorName"),
            product_data.get("colors"),
            product_data.get("availableColors"),
        ]
        
        for source in color_sources:
            if isinstance(source, str):
                colors.append(source)
            elif isinstance(source, list):
                colors.extend([str(c) for c in source if c])
        
        # Also check for color variants
        variants = product_data.get("variants", []) or json_data.get("variants", [])
        for variant in variants:
            if isinstance(variant, dict):
                color = variant.get("color") or variant.get("colorName") or variant.get("colorway")
                if color:
                    colors.append(str(color))
        
        return list(dict.fromkeys(colors))  # Remove duplicates
    
    def _extract_sizes_from_data(self, json_data: dict, product_data: dict) -> List[str]:
        """Extract size options from JSON data."""
        sizes = []
        
        # Check various possible locations
        size_sources = [
            product_data.get("sizes"),
            product_data.get("availableSizes"),
            product_data.get("sizeOptions"),
            json_data.get("sizes"),
        ]
        
        for source in size_sources:
            if isinstance(source, list):
                sizes.extend([str(s) for s in source if s])
            elif isinstance(source, str):
                sizes.append(source)
        
        return list(dict.fromkeys(sizes))  # Remove duplicates
    
    def _extract_colors_from_html(self, soup) -> List[str]:
        """Extract color options from HTML."""
        colors = []
        
        # Look for color swatches/options
        color_selectors = [
            'button[data-testid*="color"]',
            'button[aria-label*="color"]',
            '.color-selector button',
            '.color-picker button',
            '[data-color]',
            '.color-swatch',
            '[class*="color"] button',
        ]
        
        for selector in color_selectors:
            elements = soup.select(selector)
            for elem in elements:
                color = (
                    elem.get("aria-label", "") or
                    elem.get("data-color", "") or
                    elem.get("data-colorway", "") or
                    elem.get("title", "") or
                    elem.get_text(strip=True)
                )
                if color and color not in colors and len(color) > 1:
                    colors.append(color)
        
        return colors
    
    def _extract_sizes_from_html(self, soup) -> List[str]:
        """Extract size options from HTML."""
        sizes = []
        
        # Look for size buttons/options
        size_selectors = [
            'button[data-testid*="size"]',
            'button[aria-label*="size"]',
            '.size-selector button',
            '.size-picker button',
            '[data-size]',
            '.size-option',
            '[class*="size"] button',
        ]
        
        for selector in size_selectors:
            elements = soup.select(selector)
            for elem in elements:
                size = (
                    elem.get("aria-label", "") or
                    elem.get("data-size", "") or
                    elem.get_text(strip=True)
                )
                if size:
                    # Extract numeric size if present
                    size_match = re.search(r'(\d+(?:\.\d+)?)', size)
                    if size_match:
                        sizes.append(size_match.group(1))
                    elif size.strip() and size.strip() not in sizes:
                        sizes.append(size.strip())
        
        return sizes
    
    def _parse_price_from_json(self, json_data: dict) -> List[PriceHistoryResponse]:
        """Parse price from JSON data."""
        price_history = []
        
        offers = json_data.get("offers", {})
        if isinstance(offers, list):
            offers = offers[0] if offers else {}
        
        # Try various price fields
        price_str = (
            offers.get("price") or
            json_data.get("price") or
            json_data.get("product", {}).get("price", "") if isinstance(json_data.get("product"), dict) else ""
        )
        
        availability = (
            offers.get("availability") or
            json_data.get("availability", "") or
            json_data.get("product", {}).get("availability", "") if isinstance(json_data.get("product"), dict) else ""
        )
        
        if price_str:
            try:
                price = float(str(price_str).replace("$", "").replace(",", "").replace("€", "").replace("£", ""))
                is_in_stock = (
                    "in stock" in str(availability).lower() or
                    availability == "" or
                    "available" in str(availability).lower()
                )
                price_history.append(PriceHistoryResponse(
                    price=price,
                    is_in_stock=is_in_stock
                ))
            except (ValueError, TypeError):
                pass
        
        return price_history
    
    def _parse_price_from_html(self, soup) -> List[PriceHistoryResponse]:
        """Parse price from HTML."""
        price_history = []
        
        # Look for price elements
        price_selectors = [
            '[data-testid*="price"]',
            '.product-price',
            '.price',
            '.gl-price',
            '[itemprop="price"]',
            '[class*="price"]',
        ]
        
        price = None
        for selector in price_selectors:
            price_elem = soup.select_one(selector)
            if price_elem:
                price_text = price_elem.get_text(strip=True)
                # Remove currency symbols and extract number
                price_match = re.search(r'[\d,]+\.?\d*', price_text.replace(",", ""))
                if price_match:
                    try:
                        price = float(price_match.group(0))
                        break
                    except ValueError:
                        continue
        
        # Check stock status
        stock_selectors = [
            '[data-testid*="stock"]',
            '.stock-status',
            '[aria-label*="stock"]',
            '[class*="stock"]',
            '.availability',
        ]
        
        is_in_stock = True  # Default to in stock
        for selector in stock_selectors:
            stock_elem = soup.select_one(selector)
            if stock_elem:
                stock_text = stock_elem.get_text(strip=True).lower()
                if "out of stock" in stock_text or "unavailable" in stock_text or "sold out" in stock_text:
                    is_in_stock = False
                    break
        
        if price is not None:
            price_history.append(PriceHistoryResponse(
                price=price,
                is_in_stock=is_in_stock
            ))
        
        return price_history
    
    def _parse_images_from_json(self, json_data: dict) -> List[str]:
        """Parse images from JSON data."""
        images = []
        
        # Check various image fields
        image_sources = [
            json_data.get("image"),
            json_data.get("images"),
            json_data.get("product", {}).get("images", []) if isinstance(json_data.get("product"), dict) else [],
            json_data.get("product", {}).get("imageUrls", []) if isinstance(json_data.get("product"), dict) else [],
        ]
        
        for source in image_sources:
            if isinstance(source, str):
                images.append(source)
            elif isinstance(source, list):
                for img in source:
                    if isinstance(img, str):
                        images.append(img)
                    elif isinstance(img, dict):
                        url = img.get("url") or img.get("src") or img.get("image") or img.get("href")
                        if url:
                            images.append(url)
        
        return images
    
    def _parse_images_from_html(self, soup) -> List[str]:
        """Parse images from HTML."""
        images = []
        
        # Look for product images
        img_selectors = [
            'img[data-testid*="product"]',
            '.product-image img',
            '.gallery img',
            '.image-carousel img',
            '[data-image]',
            '[data-src]',
            '.gl-image img',
        ]
        
        for selector in img_selectors:
            img_elements = soup.select(selector)
            for img in img_elements:
                # Try various attributes
                img_url = (
                    img.get("src") or
                    img.get("data-src") or
                    img.get("data-image") or
                    img.get("data-lazy-src") or
                    img.get("data-srcset")
                )
                if img_url:
                    # Handle srcset (can contain multiple URLs)
                    if "," in img_url and " " in img_url:
                        urls = [u.strip().split()[0] for u in img_url.split(",")]
                        images.extend(urls)
                    else:
                        images.append(img_url)
        
        # Also check for picture/source elements
        picture_elements = soup.find_all("picture")
        for picture in picture_elements:
            source = picture.find("source") or picture.find("img")
            if source:
                img_url = source.get("srcset") or source.get("src")
                if img_url:
                    # srcset can contain multiple URLs
                    if "," in img_url:
                        urls = [u.strip().split()[0] for u in img_url.split(",")]
                        images.extend(urls)
                    else:
                        images.append(img_url)
        
        # Check for data attributes on containers
        image_containers = soup.select('[data-image-url], [data-img-url]')
        for container in image_containers:
            img_url = container.get("data-image-url") or container.get("data-img-url")
            if img_url:
                images.append(img_url)
        
        return images
    
    def _normalize_image_url(self, url: str) -> str:
        """Normalize image URL to full URL."""
        if not url:
            return ""
        
        # If already a full URL, return as is
        if url.startswith("http://") or url.startswith("https://"):
            return url
        
        # If relative URL, make it absolute
        if url.startswith("//"):
            return f"https:{url}"
        
        if url.startswith("/"):
            return f"{self.base_url}{url}"
        
        # Otherwise assume it's a relative path
        return f"{self.base_url}/{url.lstrip('/')}"
