package repositories

import (
	"sneaker-price-tracker/internal/models"
	"testing"
)

func TestProductVariantRepository_Create(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	
	// Create a product first
	productReq := models.CreateProductRequest{
		BrandName: "Nike",
		ShoeName:  "Air Jordan 1",
		BaseURL:   "https://example.com/nike/aj1",
	}
	
	product, err := productRepo.Create(productReq)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	
	// Create variant
	variantReq := models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Black/White",
		ShoeSize:  "10.5",
	}
	
	variant, err := variantRepo.Create(variantReq)
	if err != nil {
		t.Fatalf("Failed to create variant: %v", err)
	}
	
	if variant.ID == 0 {
		t.Error("Expected variant ID to be set")
	}
	if variant.ProductID != product.ID {
		t.Errorf("Expected product ID %d, got %d", product.ID, variant.ProductID)
	}
	if variant.Color != variantReq.Color {
		t.Errorf("Expected color %s, got %s", variantReq.Color, variant.Color)
	}
}

func TestProductVariantRepository_GetByID(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Adidas",
		ShoeName:  "Yeezy",
		BaseURL:   "https://example.com/adidas/yeezy",
	})
	
	variantReq := models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Zebra",
		ShoeSize:  "9",
	}
	
	created, err := variantRepo.Create(variantReq)
	if err != nil {
		t.Fatalf("Failed to create variant: %v", err)
	}
	
	variant, err := variantRepo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get variant: %v", err)
	}
	
	if variant.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, variant.ID)
	}
}

func TestProductVariantRepository_GetByProductID(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Nike",
		ShoeName:  "Dunk",
		BaseURL:   "https://example.com/nike/dunk",
	})
	
	variants := []models.CreateProductVariantRequest{
		{ProductID: product.ID, Color: "Red", ShoeSize: "10"},
		{ProductID: product.ID, Color: "Blue", ShoeSize: "11"},
	}
	
	for _, req := range variants {
		_, err := variantRepo.Create(req)
		if err != nil {
			t.Fatalf("Failed to create variant: %v", err)
		}
	}
	
	productVariants, err := variantRepo.GetByProductID(product.ID)
	if err != nil {
		t.Fatalf("Failed to get variants: %v", err)
	}
	
	if len(productVariants) < len(variants) {
		t.Errorf("Expected at least %d variants, got %d", len(variants), len(productVariants))
	}
}

func TestProductVariantRepository_Update(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Puma",
		ShoeName:  "Suede",
		BaseURL:   "https://example.com/puma/suede",
	})
	
	variantReq := models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Black",
		ShoeSize:  "9.5",
	}
	
	created, err := variantRepo.Create(variantReq)
	if err != nil {
		t.Fatalf("Failed to create variant: %v", err)
	}
	
	newColor := "White"
	updateReq := models.UpdateProductVariantRequest{
		Color: &newColor,
	}
	
	updated, err := variantRepo.Update(created.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update variant: %v", err)
	}
	
	if updated.Color != newColor {
		t.Errorf("Expected color %s, got %s", newColor, updated.Color)
	}
}

func TestProductVariantRepository_Delete(t *testing.T) {
	productRepo := NewProductRepository(testDB)
	variantRepo := NewProductVariantRepository(testDB)
	
	product, _ := productRepo.Create(models.CreateProductRequest{
		BrandName: "Nike",
		ShoeName:  "Air Max",
		BaseURL:   "https://example.com/nike/airmax",
	})
	
	variantReq := models.CreateProductVariantRequest{
		ProductID: product.ID,
		Color:     "Grey",
		ShoeSize:  "10",
	}
	
	created, err := variantRepo.Create(variantReq)
	if err != nil {
		t.Fatalf("Failed to create variant: %v", err)
	}
	
	err = variantRepo.Delete(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete variant: %v", err)
	}
	
	_, err = variantRepo.GetByID(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted variant")
	}
}
