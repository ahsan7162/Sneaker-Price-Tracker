"""Scrapers package exports."""

from .base import BaseScraper, ScrapeResult
from .nike import NikeScraper
from .adidas import AdidasScraper

__all__ = ["BaseScraper", "ScrapeResult", "NikeScraper", "AdidasScraper"]
