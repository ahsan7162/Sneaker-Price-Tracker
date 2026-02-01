package repositories

import (
	"database/sql"
	"os"
	"sneaker-price-tracker/internal/config"
	"sneaker-price-tracker/internal/db"
	"sneaker-price-tracker/internal/migrations"
	"sneaker-price-tracker/internal/models"
	"testing"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Setup test database
	cfg := config.LoadConfig()
	// Use test database
	cfg.Database.DBName = "sneaker_tracker_test"
	
	database, err := db.NewDB(cfg)
	if err != nil {
		panic(err)
	}
	testDB = database.DB
	
	// Run migrations
	if err := migrations.RunMigrations(testDB); err != nil {
		panic(err)
	}
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	migrations.DownMigrations(testDB)
	testDB.Close()
	
	os.Exit(code)
}

func TestProductRepository_Create(t *testing.T) {
	repo := NewProductRepository(testDB)
	
	req := models.CreateProductRequest{
		BrandName: "Nike",
		ShoeName:  "Air Jordan 1",
		BaseURL:   "https://example.com/nike/air-jordan-1",
	}
	
	product, err := repo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	
	if product.ID == 0 {
		t.Error("Expected product ID to be set")
	}
	if product.BrandName != req.BrandName {
		t.Errorf("Expected brand name %s, got %s", req.BrandName, product.BrandName)
	}
	if product.ShoeName != req.ShoeName {
		t.Errorf("Expected shoe name %s, got %s", req.ShoeName, product.ShoeName)
	}
}

func TestProductRepository_GetByID(t *testing.T) {
	repo := NewProductRepository(testDB)
	
	// Create a product first
	req := models.CreateProductRequest{
		BrandName: "Adidas",
		ShoeName:  "Yeezy Boost 350",
		BaseURL:   "https://example.com/adidas/yeezy",
	}
	
	created, err := repo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	
	// Get by ID
	product, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}
	
	if product.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, product.ID)
	}
	if product.BrandName != req.BrandName {
		t.Errorf("Expected brand name %s, got %s", req.BrandName, product.BrandName)
	}
}

func TestProductRepository_GetAll(t *testing.T) {
	repo := NewProductRepository(testDB)
	
	// Create multiple products
	products := []models.CreateProductRequest{
		{BrandName: "Nike", ShoeName: "Air Max", BaseURL: "https://example.com/nike/air-max"},
		{BrandName: "Puma", ShoeName: "Suede", BaseURL: "https://example.com/puma/suede"},
	}
	
	for _, req := range products {
		_, err := repo.Create(req)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}
	}
	
	// Get all products
	allProducts, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all products: %v", err)
	}
	
	if len(allProducts) < len(products) {
		t.Errorf("Expected at least %d products, got %d", len(products), len(allProducts))
	}
}

func TestProductRepository_Update(t *testing.T) {
	repo := NewProductRepository(testDB)
	
	// Create a product
	req := models.CreateProductRequest{
		BrandName: "Nike",
		ShoeName:  "Air Force 1",
		BaseURL:   "https://example.com/nike/af1",
	}
	
	created, err := repo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	
	// Update product
	newBrandName := "New Balance"
	updateReq := models.UpdateProductRequest{
		BrandName: &newBrandName,
	}
	
	updated, err := repo.Update(created.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}
	
	if updated.BrandName != newBrandName {
		t.Errorf("Expected brand name %s, got %s", newBrandName, updated.BrandName)
	}
	if updated.ShoeName != req.ShoeName {
		t.Errorf("Expected shoe name to remain %s, got %s", req.ShoeName, updated.ShoeName)
	}
}

func TestProductRepository_Delete(t *testing.T) {
	repo := NewProductRepository(testDB)
	
	// Create a product
	req := models.CreateProductRequest{
		BrandName: "Nike",
		ShoeName:  "Dunk Low",
		BaseURL:   "https://example.com/nike/dunk",
	}
	
	created, err := repo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	
	// Delete product
	err = repo.Delete(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}
	
	// Verify deletion
	_, err = repo.GetByID(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted product")
	}
}
