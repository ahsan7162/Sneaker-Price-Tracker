package main

import (
	"fmt"
	"log"
	"os"
	"sneaker-price-tracker/internal/config"
	"sneaker-price-tracker/internal/db"
	"sneaker-price-tracker/internal/migrations"
)

func main() {
	cfg := config.LoadConfig()
	
	database, err := db.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	
	switch command {
	case "up":
		fmt.Println("Running migrations...")
		if err := migrations.RunMigrations(database.DB); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully!")
	case "down":
		fmt.Println("Rolling back migrations...")
		if err := migrations.DownMigrations(database.DB); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		fmt.Println("Migrations rolled back successfully!")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage: go run cmd/migrate/main.go [up|down]")
		os.Exit(1)
	}
}
