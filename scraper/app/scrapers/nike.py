"""Nike product scraper implementation."""

import json
import re
from typing import List, Set
from urllib.parse import urlparse

from app.models.product_variant import ProductVariantResponse
from app.models.price_history import PriceHistoryResponse
from app.scrapers.base import BaseScraper


class NikeScraper(BaseScraper):
    """Scraper for Nike product pages."""
    
    def __init__(self, timeout: int = 30):
        """Initialize Nike scraper."""
        super().__init__(timeout)
        self.base_url = "https://www.nike.com"
    
    async def _parse_variants(self, html: str, url: str) -> List[ProductVariantResponse]:
        """
        Parse product variants from Nike HTML.
        
        Extracts color and size combinations, along with unique identifiers.
        """
        soup = self._parse_html(html)
        variants = []
        
        # Try to extract from JSON-LD structured data first
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
        """Parse price and stock information from Nike HTML."""
        soup = self._parse_html(html)
        price_history = []
        
        # Try JSON-LD first
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
        """Parse product images from Nike HTML."""
        soup = self._parse_html(html)
        images = []
        
        # Try JSON-LD first
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
        """Extract JSON-LD structured data from HTML."""
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
        
        # Look for window.__INITIAL_STATE__ or similar React state
        state_patterns = [
            r'window\.__INITIAL_STATE__\s*=\s*({.*?});',
            r'window\.__NEXT_DATA__\s*=\s*({.*?});',
            r'"product":\s*({.*?})',
        ]
        
        for pattern in state_patterns:
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
        """Parse variants from JSON-LD or React state."""
        variants = []
        
        # Extract from offers or product variants
        offers = json_data.get("offers", {})
        if isinstance(offers, list):
            offers = offers[0] if offers else {}
        
        # Try to get SKU/identifier
        sku = json_data.get("sku") or json_data.get("productID") or json_data.get("mpn", "")
        
        # Extract colors and sizes from product data
        # Nike often stores this in nested structures
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
        
        # Extract style code from URL (e.g., IM6674-101)
        style_code = self._extract_style_code_from_url(url)
        
        # Find color swatches/options
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
                        unique_identifier=style_code or f"{color}-{size}"
                    ))
        elif colors:
            for color in colors:
                variants.append(ProductVariantResponse(
                    color=color,
                    shoe_size="N/A",
                    unique_identifier=style_code or color
                ))
        elif sizes:
            for size in sizes:
                variants.append(ProductVariantResponse(
                    color="N/A",
                    shoe_size=str(size),
                    unique_identifier=style_code or str(size)
                ))
        
        return variants
    
    def _create_default_variant(self, soup, url: str) -> ProductVariantResponse:
        """Create a default variant when parsing fails."""
        style_code = self._extract_style_code_from_url(url)
        
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
            unique_identifier=style_code or "UNKNOWN"
        )
    
    def _extract_style_code_from_url(self, url: str) -> str:
        """Extract style code from Nike URL."""
        # URL format: https://www.nike.com/t/{slug}/{style-code}
        parts = url.rstrip("/").split("/")
        if len(parts) >= 2:
            style_code = parts[-1]
            # Style codes are typically alphanumeric with dashes
            if re.match(r'^[A-Z0-9-]+$', style_code):
                return style_code
        return ""
    
    def _extract_colors_from_data(self, json_data: dict, product_data: dict) -> List[str]:
        """Extract color options from JSON data."""
        colors = []
        
        # Check various possible locations
        color_sources = [
            json_data.get("color"),
            product_data.get("color"),
            product_data.get("colors"),
            json_data.get("itemListElement", [{}])[0].get("item", {}).get("color"),
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
                color = variant.get("color") or variant.get("colorName")
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
        
        # Look for color swatches
        color_selectors = [
            'button[data-testid*="color"]',
            'button[aria-label*="color"]',
            '.color-selector button',
            '[data-color]',
            '.color-swatch',
        ]
        
        for selector in color_selectors:
            elements = soup.select(selector)
            for elem in elements:
                color = (
                    elem.get("aria-label", "") or
                    elem.get("data-color", "") or
                    elem.get("title", "") or
                    elem.get_text(strip=True)
                )
                if color and color not in colors:
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
            '[data-size]',
            '.size-option',
        ]
        
        for selector in size_selectors:
            elements = soup.select(selector)
            for elem in elements:
                size = (
                    elem.get("aria-label", "") or
                    elem.get("data-size", "") or
                    elem.get_text(strip=True)
                )
                if size and size not in sizes:
                    # Extract numeric size if present
                    size_match = re.search(r'(\d+(?:\.\d+)?)', size)
                    if size_match:
                        sizes.append(size_match.group(1))
                    elif size.strip():
                        sizes.append(size.strip())
        
        return sizes
    
    def _parse_price_from_json(self, json_data: dict) -> List[PriceHistoryResponse]:
        """Parse price from JSON data."""
        price_history = []
        
        offers = json_data.get("offers", {})
        if isinstance(offers, list):
            offers = offers[0] if offers else {}
        
        price_str = offers.get("price") or json_data.get("price")
        availability = offers.get("availability") or json_data.get("availability", "")
        
        if price_str:
            try:
                price = float(str(price_str).replace("$", "").replace(",", ""))
                is_in_stock = "in stock" in str(availability).lower() or availability == ""
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
            '[itemprop="price"]',
        ]
        
        price = None
        for selector in price_selectors:
            price_elem = soup.select_one(selector)
            if price_elem:
                price_text = price_elem.get_text(strip=True)
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
        ]
        
        is_in_stock = True  # Default to in stock
        for selector in stock_selectors:
            stock_elem = soup.select_one(selector)
            if stock_elem:
                stock_text = stock_elem.get_text(strip=True).lower()
                if "out of stock" in stock_text or "unavailable" in stock_text:
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
        ]
        
        for source in image_sources:
            if isinstance(source, str):
                images.append(source)
            elif isinstance(source, list):
                for img in source:
                    if isinstance(img, str):
                        images.append(img)
                    elif isinstance(img, dict):
                        url = img.get("url") or img.get("src") or img.get("image")
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
            '[data-image]',
        ]
        
        for selector in img_selectors:
            img_elements = soup.select(selector)
            for img in img_elements:
                # Try various attributes
                img_url = (
                    img.get("src") or
                    img.get("data-src") or
                    img.get("data-image") or
                    img.get("data-lazy-src")
                )
                if img_url:
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
