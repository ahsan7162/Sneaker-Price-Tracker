package models

import "time"

// Product represents the products table
type Product struct {
	ID        int64     `json:"id" db:"id"`
	BrandName string    `json:"brand_name" db:"brand_name"`
	ShoeName  string    `json:"shoe_name" db:"shoe_name"`
	BaseURL   string    `json:"base_url" db:"base_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateProductRequest represents the data needed to create a product
type CreateProductRequest struct {
	BrandName string `json:"brand_name" validate:"required"`
	ShoeName  string `json:"shoe_name" validate:"required"`
	BaseURL   string `json:"base_url" validate:"required,url"`
}

// UpdateProductRequest represents the data needed to update a product
type UpdateProductRequest struct {
	BrandName *string `json:"brand_name,omitempty"`
	ShoeName  *string `json:"shoe_name,omitempty"`
	BaseURL   *string `json:"base_url,omitempty"`
}
