package repositories

import (
	"database/sql"
	"fmt"
	"sneaker-price-tracker/internal/models"
)

// ProductVariantRepository handles product variant database operations
type ProductVariantRepository struct {
	db *sql.DB
}

// NewProductVariantRepository creates a new product variant repository
func NewProductVariantRepository(db *sql.DB) *ProductVariantRepository {
	return &ProductVariantRepository{db: db}
}

// Create creates a new product variant
func (r *ProductVariantRepository) Create(req models.CreateProductVariantRequest) (*models.ProductVariant, error) {
	// Generate unique identifier
	uniqueIdentifier := fmt.Sprintf("%d_%s_%s", req.ProductID, req.Color, req.ShoeSize)

	query := `
		INSERT INTO product_variants (product_id, color, shoe_size, unique_identifier)
		VALUES ($1, $2, $3, $4)
		RETURNING id, product_id, color, shoe_size, unique_identifier
	`

	var variant models.ProductVariant
	err := r.db.QueryRow(query, req.ProductID, req.Color, req.ShoeSize, uniqueIdentifier).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.Color,
		&variant.ShoeSize,
		&variant.UniqueIdentifier,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product variant: %w", err)
	}

	return &variant, nil
}

// GetByID retrieves a product variant by ID
func (r *ProductVariantRepository) GetByID(id int64) (*models.ProductVariant, error) {
	query := `
		SELECT id, product_id, color, shoe_size, unique_identifier
		FROM product_variants
		WHERE id = $1
	`

	var variant models.ProductVariant
	err := r.db.QueryRow(query, id).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.Color,
		&variant.ShoeSize,
		&variant.UniqueIdentifier,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product variant not found")
		}
		return nil, fmt.Errorf("failed to get product variant: %w", err)
	}

	return &variant, nil
}

// GetByProductID retrieves all variants for a product
func (r *ProductVariantRepository) GetByProductID(productID int64) ([]models.ProductVariant, error) {
	query := `
		SELECT id, product_id, color, shoe_size, unique_identifier
		FROM product_variants
		WHERE product_id = $1
		ORDER BY color, shoe_size
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}
	defer rows.Close()

	var variants []models.ProductVariant
	for rows.Next() {
		var variant models.ProductVariant
		if err := rows.Scan(
			&variant.ID,
			&variant.ProductID,
			&variant.Color,
			&variant.ShoeSize,
			&variant.UniqueIdentifier,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product variant: %w", err)
		}
		variants = append(variants, variant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product variants: %w", err)
	}

	return variants, nil
}

// GetAll retrieves all product variants
func (r *ProductVariantRepository) GetAll() ([]models.ProductVariant, error) {
	query := `
		SELECT id, product_id, color, shoe_size, unique_identifier
		FROM product_variants
		ORDER BY product_id, color, shoe_size
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}
	defer rows.Close()

	var variants []models.ProductVariant
	for rows.Next() {
		var variant models.ProductVariant
		if err := rows.Scan(
			&variant.ID,
			&variant.ProductID,
			&variant.Color,
			&variant.ShoeSize,
			&variant.UniqueIdentifier,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product variant: %w", err)
		}
		variants = append(variants, variant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product variants: %w", err)
	}

	return variants, nil
}

// Update updates a product variant
func (r *ProductVariantRepository) Update(id int64, req models.UpdateProductVariantRequest) (*models.ProductVariant, error) {
	// First get the current variant to rebuild unique_identifier if needed
	currentVariant, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Build dynamic update query
	query := "UPDATE product_variants SET "
	args := []interface{}{}
	argPos := 1

	color := currentVariant.Color
	shoeSize := currentVariant.ShoeSize

	if req.Color != nil {
		query += fmt.Sprintf("color = $%d, ", argPos)
		args = append(args, *req.Color)
		color = *req.Color
		argPos++
	}
	if req.ShoeSize != nil {
		query += fmt.Sprintf("shoe_size = $%d, ", argPos)
		args = append(args, *req.ShoeSize)
		shoeSize = *req.ShoeSize
		argPos++
	}

	// Update unique_identifier if color or size changed
	uniqueIdentifier := fmt.Sprintf("%d_%s_%s", currentVariant.ProductID, color, shoeSize)
	query += fmt.Sprintf("unique_identifier = $%d, ", argPos)
	args = append(args, uniqueIdentifier)
	argPos++

	// Remove trailing comma and add WHERE clause
	query = query[:len(query)-2] + " WHERE id = $" + fmt.Sprintf("%d", argPos)
	args = append(args, id)

	// Add RETURNING clause
	query += " RETURNING id, product_id, color, shoe_size, unique_identifier"

	var variant models.ProductVariant
	err = r.db.QueryRow(query, args...).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.Color,
		&variant.ShoeSize,
		&variant.UniqueIdentifier,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product variant not found")
		}
		return nil, fmt.Errorf("failed to update product variant: %w", err)
	}

	return &variant, nil
}

// Delete deletes a product variant
func (r *ProductVariantRepository) Delete(id int64) error {
	query := "DELETE FROM product_variants WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product variant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product variant not found")
	}

	return nil
}
