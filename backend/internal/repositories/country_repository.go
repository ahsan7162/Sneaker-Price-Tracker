package repositories

import (
	"database/sql"
	"fmt"
	"sneaker-price-tracker/internal/models"
)

// CountryRepository handles country database operations
type CountryRepository struct {
	db *sql.DB
}

// NewCountryRepository creates a new country repository
func NewCountryRepository(db *sql.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

// Create creates a new country
func (r *CountryRepository) Create(req models.CreateCountryRequest) (*models.Country, error) {
	query := `
		INSERT INTO countries (country_code, currency)
		VALUES ($1, $2)
		RETURNING id, country_code, currency
	`

	var country models.Country
	err := r.db.QueryRow(query, req.CountryCode, req.Currency).Scan(
		&country.ID,
		&country.CountryCode,
		&country.Currency,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create country: %w", err)
	}

	return &country, nil
}

// GetByID retrieves a country by ID
func (r *CountryRepository) GetByID(id int64) (*models.Country, error) {
	query := `
		SELECT id, country_code, currency
		FROM countries
		WHERE id = $1
	`

	var country models.Country
	err := r.db.QueryRow(query, id).Scan(
		&country.ID,
		&country.CountryCode,
		&country.Currency,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to get country: %w", err)
	}

	return &country, nil
}

// GetByCountryCode retrieves a country by country code
func (r *CountryRepository) GetByCountryCode(countryCode string) (*models.Country, error) {
	query := `
		SELECT id, country_code, currency
		FROM countries
		WHERE country_code = $1
	`

	var country models.Country
	err := r.db.QueryRow(query, countryCode).Scan(
		&country.ID,
		&country.CountryCode,
		&country.Currency,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to get country: %w", err)
	}

	return &country, nil
}

// GetAll retrieves all countries
func (r *CountryRepository) GetAll() ([]models.Country, error) {
	query := `
		SELECT id, country_code, currency
		FROM countries
		ORDER BY country_code
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get countries: %w", err)
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		if err := rows.Scan(
			&country.ID,
			&country.CountryCode,
			&country.Currency,
		); err != nil {
			return nil, fmt.Errorf("failed to scan country: %w", err)
		}
		countries = append(countries, country)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating countries: %w", err)
	}

	return countries, nil
}

// Update updates a country
func (r *CountryRepository) Update(id int64, req models.UpdateCountryRequest) (*models.Country, error) {
	// Build dynamic update query
	query := "UPDATE countries SET "
	args := []interface{}{}
	argPos := 1

	if req.CountryCode != nil {
		query += fmt.Sprintf("country_code = $%d, ", argPos)
		args = append(args, *req.CountryCode)
		argPos++
	}
	if req.Currency != nil {
		query += fmt.Sprintf("currency = $%d, ", argPos)
		args = append(args, *req.Currency)
		argPos++
	}

	// Remove trailing comma and add WHERE clause
	query = query[:len(query)-2] + " WHERE id = $" + fmt.Sprintf("%d", argPos)
	args = append(args, id)

	// Add RETURNING clause
	query += " RETURNING id, country_code, currency"

	var country models.Country
	err := r.db.QueryRow(query, args...).Scan(
		&country.ID,
		&country.CountryCode,
		&country.Currency,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to update country: %w", err)
	}

	return &country, nil
}

// Delete deletes a country
func (r *CountryRepository) Delete(id int64) error {
	query := "DELETE FROM countries WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete country: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("country not found")
	}

	return nil
}
