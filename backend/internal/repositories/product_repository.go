package repositories

import (
	"database/sql"
	"fmt"
	"sneaker-price-tracker/internal/models"
)

// ProductRepository handles product database operations
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product
func (r *ProductRepository) Create(req models.CreateProductRequest) (*models.Product, error) {
	query := `
		INSERT INTO products (brand_name, shoe_name, base_url, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, brand_name, shoe_name, base_url, created_at
	`

	var product models.Product
	err := r.db.QueryRow(query, req.BrandName, req.ShoeName, req.BaseURL).Scan(
		&product.ID,
		&product.BrandName,
		&product.ShoeName,
		&product.BaseURL,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &product, nil
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(id int64) (*models.Product, error) {
	query := `
		SELECT id, brand_name, shoe_name, base_url, created_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.BrandName,
		&product.ShoeName,
		&product.BaseURL,
		&product.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// GetAll retrieves all products
func (r *ProductRepository) GetAll() ([]models.Product, error) {
	query := `
		SELECT id, brand_name, shoe_name, base_url, created_at
		FROM products
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.BrandName,
			&product.ShoeName,
			&product.BaseURL,
			&product.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// Update updates a product
func (r *ProductRepository) Update(id int64, req models.UpdateProductRequest) (*models.Product, error) {
	// Build dynamic update query
	query := "UPDATE products SET "
	args := []interface{}{}
	argPos := 1

	if req.BrandName != nil {
		query += fmt.Sprintf("brand_name = $%d, ", argPos)
		args = append(args, *req.BrandName)
		argPos++
	}
	if req.ShoeName != nil {
		query += fmt.Sprintf("shoe_name = $%d, ", argPos)
		args = append(args, *req.ShoeName)
		argPos++
	}
	if req.BaseURL != nil {
		query += fmt.Sprintf("base_url = $%d, ", argPos)
		args = append(args, *req.BaseURL)
		argPos++
	}

	// Remove trailing comma and add WHERE clause
	query = query[:len(query)-2] + " WHERE id = $" + fmt.Sprintf("%d", argPos)
	args = append(args, id)

	// Add RETURNING clause
	query += " RETURNING id, brand_name, shoe_name, base_url, created_at"

	var product models.Product
	err := r.db.QueryRow(query, args...).Scan(
		&product.ID,
		&product.BrandName,
		&product.ShoeName,
		&product.BaseURL,
		&product.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &product, nil
}

// Delete deletes a product
func (r *ProductRepository) Delete(id int64) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
