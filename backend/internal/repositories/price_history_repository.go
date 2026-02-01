package repositories

import (
	"database/sql"
	"fmt"
	"sneaker-price-tracker/internal/models"
	"time"
)

// PriceHistoryRepository handles price history database operations
type PriceHistoryRepository struct {
	db *sql.DB
}

// NewPriceHistoryRepository creates a new price history repository
func NewPriceHistoryRepository(db *sql.DB) *PriceHistoryRepository {
	return &PriceHistoryRepository{db: db}
}

// Create creates a new price history entry
func (r *PriceHistoryRepository) Create(req models.CreatePriceHistoryRequest) (*models.PriceHistory, error) {
	query := `
		INSERT INTO price_history (variant_id, country_id, price, is_in_stock, captured_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, variant_id, country_id, price, is_in_stock, captured_at
	`

	var history models.PriceHistory
	err := r.db.QueryRow(query, req.VariantID, req.CountryID, req.Price, req.IsInStock).Scan(
		&history.ID,
		&history.VariantID,
		&history.CountryID,
		&history.Price,
		&history.IsInStock,
		&history.CapturedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create price history: %w", err)
	}

	return &history, nil
}

// GetByID retrieves a price history entry by ID
func (r *PriceHistoryRepository) GetByID(id int64) (*models.PriceHistory, error) {
	query := `
		SELECT id, variant_id, country_id, price, is_in_stock, captured_at
		FROM price_history
		WHERE id = $1
	`

	var history models.PriceHistory
	err := r.db.QueryRow(query, id).Scan(
		&history.ID,
		&history.VariantID,
		&history.CountryID,
		&history.Price,
		&history.IsInStock,
		&history.CapturedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("price history not found")
		}
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}

	return &history, nil
}

// GetByVariantID retrieves all price history entries for a variant
func (r *PriceHistoryRepository) GetByVariantID(variantID int64) ([]models.PriceHistory, error) {
	query := `
		SELECT id, variant_id, country_id, price, is_in_stock, captured_at
		FROM price_history
		WHERE variant_id = $1
		ORDER BY captured_at DESC
	`

	rows, err := r.db.Query(query, variantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}
	defer rows.Close()

	var histories []models.PriceHistory
	for rows.Next() {
		var history models.PriceHistory
		if err := rows.Scan(
			&history.ID,
			&history.VariantID,
			&history.CountryID,
			&history.Price,
			&history.IsInStock,
			&history.CapturedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan price history: %w", err)
		}
		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating price history: %w", err)
	}

	return histories, nil
}

// GetByVariantAndCountry retrieves price history for a specific variant and country
func (r *PriceHistoryRepository) GetByVariantAndCountry(variantID, countryID int64) ([]models.PriceHistory, error) {
	query := `
		SELECT id, variant_id, country_id, price, is_in_stock, captured_at
		FROM price_history
		WHERE variant_id = $1 AND country_id = $2
		ORDER BY captured_at DESC
	`

	rows, err := r.db.Query(query, variantID, countryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}
	defer rows.Close()

	var histories []models.PriceHistory
	for rows.Next() {
		var history models.PriceHistory
		if err := rows.Scan(
			&history.ID,
			&history.VariantID,
			&history.CountryID,
			&history.Price,
			&history.IsInStock,
			&history.CapturedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan price history: %w", err)
		}
		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating price history: %w", err)
	}

	return histories, nil
}

// GetLatestPrice retrieves the latest price for a variant and country
func (r *PriceHistoryRepository) GetLatestPrice(variantID, countryID int64) (*models.PriceHistory, error) {
	query := `
		SELECT id, variant_id, country_id, price, is_in_stock, captured_at
		FROM price_history
		WHERE variant_id = $1 AND country_id = $2
		ORDER BY captured_at DESC
		LIMIT 1
	`

	var history models.PriceHistory
	err := r.db.QueryRow(query, variantID, countryID).Scan(
		&history.ID,
		&history.VariantID,
		&history.CountryID,
		&history.Price,
		&history.IsInStock,
		&history.CapturedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("price history not found")
		}
		return nil, fmt.Errorf("failed to get latest price: %w", err)
	}

	return &history, nil
}

// GetAll retrieves all price history entries
func (r *PriceHistoryRepository) GetAll() ([]models.PriceHistory, error) {
	query := `
		SELECT id, variant_id, country_id, price, is_in_stock, captured_at
		FROM price_history
		ORDER BY captured_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}
	defer rows.Close()

	var histories []models.PriceHistory
	for rows.Next() {
		var history models.PriceHistory
		if err := rows.Scan(
			&history.ID,
			&history.VariantID,
			&history.CountryID,
			&history.Price,
			&history.IsInStock,
			&history.CapturedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan price history: %w", err)
		}
		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating price history: %w", err)
	}

	return histories, nil
}

// Update updates a price history entry
func (r *PriceHistoryRepository) Update(id int64, req models.UpdatePriceHistoryRequest) (*models.PriceHistory, error) {
	// Build dynamic update query
	query := "UPDATE price_history SET "
	args := []interface{}{}
	argPos := 1

	if req.Price != nil {
		query += fmt.Sprintf("price = $%d, ", argPos)
		args = append(args, *req.Price)
		argPos++
	}
	if req.IsInStock != nil {
		query += fmt.Sprintf("is_in_stock = $%d, ", argPos)
		args = append(args, *req.IsInStock)
		argPos++
	}

	// Remove trailing comma and add WHERE clause
	query = query[:len(query)-2] + " WHERE id = $" + fmt.Sprintf("%d", argPos)
	args = append(args, id)

	// Add RETURNING clause
	query += " RETURNING id, variant_id, country_id, price, is_in_stock, captured_at"

	var history models.PriceHistory
	err := r.db.QueryRow(query, args...).Scan(
		&history.ID,
		&history.VariantID,
		&history.CountryID,
		&history.Price,
		&history.IsInStock,
		&history.CapturedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("price history not found")
		}
		return nil, fmt.Errorf("failed to update price history: %w", err)
	}

	return &history, nil
}

// Delete deletes a price history entry
func (r *PriceHistoryRepository) Delete(id int64) error {
	query := "DELETE FROM price_history WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete price history: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("price history not found")
	}

	return nil
}

// GetPriceHistoryByDateRange retrieves price history within a date range
func (r *PriceHistoryRepository) GetPriceHistoryByDateRange(variantID, countryID int64, startDate, endDate time.Time) ([]models.PriceHistory, error) {
	query := `
		SELECT id, variant_id, country_id, price, is_in_stock, captured_at
		FROM price_history
		WHERE variant_id = $1 AND country_id = $2 AND captured_at BETWEEN $3 AND $4
		ORDER BY captured_at DESC
	`

	rows, err := r.db.Query(query, variantID, countryID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}
	defer rows.Close()

	var histories []models.PriceHistory
	for rows.Next() {
		var history models.PriceHistory
		if err := rows.Scan(
			&history.ID,
			&history.VariantID,
			&history.CountryID,
			&history.Price,
			&history.IsInStock,
			&history.CapturedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan price history: %w", err)
		}
		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating price history: %w", err)
	}

	return histories, nil
}
