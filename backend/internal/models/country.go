package models

// Country represents the countries table
type Country struct {
	ID          int64  `json:"id" db:"id"`
	CountryCode string `json:"country_code" db:"country_code"`
	Currency    string `json:"currency" db:"currency"`
}

// CreateCountryRequest represents the data needed to create a country
type CreateCountryRequest struct {
	CountryCode string `json:"country_code" validate:"required,len=2"`
	Currency    string `json:"currency" validate:"required,len=3"`
}

// UpdateCountryRequest represents the data needed to update a country
type UpdateCountryRequest struct {
	CountryCode *string `json:"country_code,omitempty"`
	Currency    *string `json:"currency,omitempty"`
}
