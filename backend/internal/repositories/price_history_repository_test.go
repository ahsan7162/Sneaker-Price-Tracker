package repositories

import (
	"sneaker-price-tracker/internal/models"
	"testing"
	"time"
)

func TestPriceHistoryRepository_Create(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	countryRepo := NewCountryRepository(testDB)
	historyRepo := NewPriceHistoryRepository(testDB)
	
	// Setup dependencies
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Nike",
		ShoeName:  "Air Jordan 1",
		BaseURL:   "https://example.com/nike/aj1",
	})
	
	variant, _ := variantRepo.Create(models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Black",
		ShoeSize:  "10",
	})
	
	country, _ := countryRepo.Create(models.CreateCountryRequest{
		CountryCode: "US",
		Currency:    "USD",
	})
	
	// Create price history
	req := models.CreatePriceHistoryRequest{
		VariantID: variant.ID,
		CountryID: country.ID,
		Price:     199.99,
		IsInStock: true,
	}
	
	history, err := historyRepo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create price history: %v", err)
	}
	
	if history.ID == 0 {
		t.Error("Expected history ID to be set")
	}
	if history.Price != req.Price {
		t.Errorf("Expected price %.2f, got %.2f", req.Price, history.Price)
	}
	if history.IsInStock != req.IsInStock {
		t.Errorf("Expected is_in_stock %v, got %v", req.IsInStock, history.IsInStock)
	}
}

func TestPriceHistoryRepository_GetByID(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	countryRepo := NewCountryRepository(testDB)
	historyRepo := NewPriceHistoryRepository(testDB)
	
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Adidas",
		ShoeName:  "Yeezy",
		BaseURL:   "https://example.com/adidas/yeezy",
	})
	
	variant, _ := variantRepo.Create(models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Zebra",
		ShoeSize:  "9",
	})
	
	country, _ := countryRepo.Create(models.CreateCountryRequest{
		CountryCode: "UK",
		Currency:    "GBP",
	})
	
	req := models.CreatePriceHistoryRequest{
		VariantID: variant.ID,
		CountryID: country.ID,
		Price:     150.00,
		IsInStock: true,
	}
	
	created, err := historyRepo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create price history: %v", err)
	}
	
	history, err := historyRepo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get price history: %v", err)
	}
	
	if history.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, history.ID)
	}
}

func TestPriceHistoryRepository_GetByVariantID(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	countryRepo := NewCountryRepository(testDB)
	historyRepo := NewPriceHistoryRepository(testDB)
	
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Nike",
		ShoeName:  "Dunk",
		BaseURL:   "https://example.com/nike/dunk",
	})
	
	variant, _ := variantRepo.Create(models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Red",
		ShoeSize:  "10",
	})
	
	country, _ := countryRepo.Create(models.CreateCountryRequest{
		CountryCode: "US",
		Currency:    "USD",
	})
	
	// Create multiple price history entries
	histories := []models.CreatePriceHistoryRequest{
		{VariantID: variant.ID, CountryID: country.ID, Price: 100.00, IsInStock: true},
		{VariantID: variant.ID, CountryID: country.ID, Price: 110.00, IsInStock: true},
	}
	
	for _, req := range histories {
		_, err := historyRepo.Create(req)
		if err != nil {
			t.Fatalf("Failed to create price history: %v", err)
		}
		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}
	
	variantHistories, err := historyRepo.GetByVariantID(variant.ID)
	if err != nil {
		t.Fatalf("Failed to get price history: %v", err)
	}
	
	if len(variantHistories) < len(histories) {
		t.Errorf("Expected at least %d histories, got %d", len(histories), len(variantHistories))
	}
}

func TestPriceHistoryRepository_GetLatestPrice(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	countryRepo := NewCountryRepository(testDB)
	historyRepo := NewPriceHistoryRepository(testDB)
	
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Puma",
		ShoeName:  "Suede",
		BaseURL:   "https://example.com/puma/suede",
	})
	
	variant, _ := variantRepo.Create(models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Black",
		ShoeSize:  "9.5",
	})
	
	country, _ := countryRepo.Create(models.CreateCountryRequest{
		CountryCode: "IN",
		Currency:    "INR",
	})
	
	// Create multiple price entries
	prices := []float64{5000.00, 5500.00, 5200.00}
	for _, price := range prices {
		_, err := historyRepo.Create(models.CreatePriceHistoryRequest{
			VariantID: variant.ID,
			CountryID: country.ID,
			Price:     price,
			IsInStock: true,
		})
		if err != nil {
			t.Fatalf("Failed to create price history: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	
	latest, err := historyRepo.GetLatestPrice(variant.ID, country.ID)
	if err != nil {
		t.Fatalf("Failed to get latest price: %v", err)
	}
	
	// Latest should be the last one created (5200.00)
	if latest.Price != 5200.00 {
		t.Errorf("Expected latest price 5200.00, got %.2f", latest.Price)
	}
}

func TestPriceHistoryRepository_Update(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	countryRepo := NewCountryRepository(testDB)
	historyRepo := NewPriceHistoryRepository(testDB)
	
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Nike",
		ShoeName:  "Air Max",
		BaseURL:   "https://example.com/nike/airmax",
	})
	
	variant, _ := variantRepo.Create(models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Grey",
		ShoeSize:  "10",
	})
	
	country, _ := countryRepo.Create(models.CreateCountryRequest{
		CountryCode: "US",
		Currency:    "USD",
	})
	
	req := models.CreatePriceHistoryRequest{
		VariantID: variant.ID,
		CountryID: country.ID,
		Price:     120.00,
		IsInStock: true,
	}
	
	created, err := historyRepo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create price history: %v", err)
	}
	
	newPrice := 130.00
	newStockStatus := false
	updateReq := models.UpdatePriceHistoryRequest{
		Price:     &newPrice,
		IsInStock: &newStockStatus,
	}
	
	updated, err := historyRepo.Update(created.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update price history: %v", err)
	}
	
	if updated.Price != newPrice {
		t.Errorf("Expected price %.2f, got %.2f", newPrice, updated.Price)
	}
	if updated.IsInStock != newStockStatus {
		t.Errorf("Expected is_in_stock %v, got %v", newStockStatus, updated.IsInStock)
	}
}

func TestPriceHistoryRepository_Delete(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	countryRepo := NewCountryRepository(testDB)
	historyRepo := NewPriceHistoryRepository(testDB)
	
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Adidas",
		ShoeName:  "Yeezy",
		BaseURL:   "https://example.com/adidas/yeezy",
	})
	
	variant, _ := variantRepo.Create(models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Zebra",
		ShoeSize:  "9",
	})
	
	country, _ := countryRepo.Create(models.CreateCountryRequest{
		CountryCode: "UK",
		Currency:    "GBP",
	})
	
	req := models.CreatePriceHistoryRequest{
		VariantID: variant.ID,
		CountryID: country.ID,
		Price:     150.00,
		IsInStock: true,
	}
	
	created, err := historyRepo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create price history: %v", err)
	}
	
	err = historyRepo.Delete(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete price history: %v", err)
	}
	
	_, err = historyRepo.GetByID(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted price history")
	}
}

func TestPriceHistoryRepository_GetPriceHistoryByDateRange(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	countryRepo := NewCountryRepository(testDB)
	historyRepo := NewPriceHistoryRepository(testDB)
	
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Nike",
		ShoeName:  "Dunk",
		BaseURL:   "https://example.com/nike/dunk",
	})
	
	variant, _ := variantRepo.Create(models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Blue",
		ShoeSize:  "11",
	})
	
	country, _ := countryRepo.Create(models.CreateCountryRequest{
		CountryCode: "US",
		Currency:    "USD",
	})
	
	startDate := time.Now()
	
	// Create entries
	_, _ = historyRepo.Create(models.CreatePriceHistoryRequest{
		VariantID: variant.ID,
		CountryID: country.ID,
		Price:     100.00,
		IsInStock: true,
	})
	time.Sleep(10 * time.Millisecond)
	
	_, _ = historyRepo.Create(models.CreatePriceHistoryRequest{
		VariantID: variant.ID,
		CountryID: country.ID,
		Price:     110.00,
		IsInStock: true,
	})
	
	endDate := time.Now().Add(1 * time.Second)
	
	histories, err := historyRepo.GetPriceHistoryByDateRange(variant.ID, country.ID, startDate, endDate)
	if err != nil {
		t.Fatalf("Failed to get price history by date range: %v", err)
	}
	
	if len(histories) < 2 {
		t.Errorf("Expected at least 2 histories, got %d", len(histories))
	}
}
