"""URL validation and brand detection utilities."""

import re
from urllib.parse import urlparse


def is_nike_url(url: str) -> bool:
    """Check if URL is a Nike product page."""
    parsed = urlparse(url)
    return (
        parsed.netloc.endswith("nike.com") or
        parsed.netloc.endswith("nike.com")
    ) and "/t/" in parsed.path


def is_adidas_url(url: str) -> bool:
    """Check if URL is an Adidas product page."""
    parsed = urlparse(url)
    return (
        parsed.netloc.endswith("adidas.com") or
        parsed.netloc.endswith("adidas.com")
    ) and parsed.path.endswith(".html")


def validate_url(url: str) -> bool:
    """Validate that URL is a valid HTTP/HTTPS URL."""
    try:
        parsed = urlparse(url)
        return parsed.scheme in ("http", "https") and bool(parsed.netloc)
    except Exception:
        return False