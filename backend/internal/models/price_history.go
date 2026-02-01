package models

import "time"

// PriceHistory represents the price_history table
type PriceHistory struct {
	ID          int64     `json:"id" db:"id"`
	VariantID   int64     `json:"variant_id" db:"variant_id"`
	CountryID   int64     `json:"country_id" db:"country_id"`
	Price       float64   `json:"price" db:"price"`
	IsInStock   bool      `json:"is_in_stock" db:"is_in_stock"`
	CapturedAt  time.Time `json:"captured_at" db:"captured_at"`
}

// CreatePriceHistoryRequest represents the data needed to create a price history entry
type CreatePriceHistoryRequest struct {
	VariantID int64   `json:"variant_id" validate:"required"`
	CountryID int64   `json:"country_id" validate:"required"`
	Price     float64 `json:"price" validate:"required,min=0"`
	IsInStock bool    `json:"is_in_stock"`
}

// UpdatePriceHistoryRequest represents the data needed to update a price history entry
type UpdatePriceHistoryRequest struct {
	Price     *float64 `json:"price,omitempty"`
	IsInStock *bool    `json:"is_in_stock,omitempty"`
}
