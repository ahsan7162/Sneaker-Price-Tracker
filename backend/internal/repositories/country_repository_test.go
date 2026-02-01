package repositories

import (
	"sneaker-price-tracker/internal/models"
	"testing"
)

func TestCountryRepository_Create(t *testing.T) {
	repo := NewCountryRepository(testDB)
	
	req := models.CreateCountryRequest{
		CountryCode: "US",
		Currency:    "USD",
	}
	
	country, err := repo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create country: %v", err)
	}
	
	if country.ID == 0 {
		t.Error("Expected country ID to be set")
	}
	if country.CountryCode != req.CountryCode {
		t.Errorf("Expected country code %s, got %s", req.CountryCode, country.CountryCode)
	}
	if country.Currency != req.Currency {
		t.Errorf("Expected currency %s, got %s", req.Currency, country.Currency)
	}
}

func TestCountryRepository_GetByID(t *testing.T) {
	repo := NewCountryRepository(testDB)
	
	req := models.CreateCountryRequest{
		CountryCode: "UK",
		Currency:    "GBP",
	}
	
	created, err := repo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create country: %v", err)
	}
	
	country, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get country: %v", err)
	}
	
	if country.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, country.ID)
	}
}

func TestCountryRepository_GetByCountryCode(t *testing.T) {
	repo := NewCountryRepository(testDB)
	
	req := models.CreateCountryRequest{
		CountryCode: "IN",
		Currency:    "INR",
	}
	
	created, err := repo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create country: %v", err)
	}
	
	country, err := repo.GetByCountryCode(req.CountryCode)
	if err != nil {
		t.Fatalf("Failed to get country by code: %v", err)
	}
	
	if country.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, country.ID)
	}
}

func TestCountryRepository_GetAll(t *testing.T) {
	repo := NewCountryRepository(testDB)
	
	countries := []models.CreateCountryRequest{
		{CountryCode: "CA", Currency: "CAD"},
		{CountryCode: "AU", Currency: "AUD"},
	}
	
	for _, req := range countries {
		_, err := repo.Create(req)
		if err != nil {
			t.Fatalf("Failed to create country: %v", err)
		}
	}
	
	allCountries, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all countries: %v", err)
	}
	
	if len(allCountries) < len(countries) {
		t.Errorf("Expected at least %d countries, got %d", len(countries), len(allCountries))
	}
}

func TestCountryRepository_Update(t *testing.T) {
	repo := NewCountryRepository(testDB)
	
	req := models.CreateCountryRequest{
		CountryCode: "DE",
		Currency:    "EUR",
	}
	
	created, err := repo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create country: %v", err)
	}
	
	newCurrency := "EUR"
	updateReq := models.UpdateCountryRequest{
		Currency: &newCurrency,
	}
	
	updated, err := repo.Update(created.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update country: %v", err)
	}
	
	if updated.Currency != newCurrency {
		t.Errorf("Expected currency %s, got %s", newCurrency, updated.Currency)
	}
}

func TestCountryRepository_Delete(t *testing.T) {
	repo := NewCountryRepository(testDB)
	
	req := models.CreateCountryRequest{
		CountryCode: "FR",
		Currency:    "EUR",
	}
	
	created, err := repo.Create(req)
	if err != nil {
		t.Fatalf("Failed to create country: %v", err)
	}
	
	err = repo.Delete(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete country: %v", err)
	}
	
	_, err = repo.GetByID(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted country")
	}
}
