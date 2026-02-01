"""Price history models matching Go structs."""

from pydantic import BaseModel, Field


class PriceHistoryResponse(BaseModel):
    """Price history response matching Go PriceHistory struct."""
    
    price: float = Field(..., description="Product price", ge=0)
    is_in_stock: bool = Field(..., description="Stock availability status")
    
    class Config:
        json_schema_extra = {
            "example": {
                "price": 145.00,
                "is_in_stock": True
            }
        }
