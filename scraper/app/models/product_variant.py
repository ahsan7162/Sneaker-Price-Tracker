"""Product variant models matching Go structs."""

from pydantic import BaseModel, Field


class ProductVariantResponse(BaseModel):
    """Product variant response matching Go ProductVariant struct."""
    
    color: str = Field(..., description="Product color")
    shoe_size: str = Field(..., description="Shoe size")
    unique_identifier: str = Field(..., description="Unique product identifier (SKU/style code)")
    
    class Config:
        json_schema_extra = {
            "example": {
                "color": "White/White/Hyper Pink/Black",
                "shoe_size": "10",
                "unique_identifier": "IM6674-101"
            }
        }
