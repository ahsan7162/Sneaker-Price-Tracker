package models

// ProductVariant represents the product_variants table
type ProductVariant struct {
	ID               int64  `json:"id" db:"id"`
	ProductID        int64  `json:"product_id" db:"product_id"`
	Color            string `json:"color" db:"color"`
	ShoeSize         string `json:"shoe_size" db:"shoe_size"`
	UniqueIdentifier string `json:"unique_identifier" db:"unique_identifier"`
}

// CreateProductVariantRequest represents the data needed to create a product variant
type CreateProductVariantRequest struct {
	ProductID int64  `json:"product_id" validate:"required"`
	Color     string `json:"color" validate:"required"`
	ShoeSize  string `json:"shoe_size" validate:"required"`
}

// UpdateProductVariantRequest represents the data needed to update a product variant
type UpdateProductVariantRequest struct {
	Color    *string `json:"color,omitempty"`
	ShoeSize *string `json:"shoe_size,omitempty"`
}
