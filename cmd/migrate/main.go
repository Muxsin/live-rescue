package main

import (
	"live-rescue/internal/database"
	"live-rescue/internal/models"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	postgres, err := database.NewPostgres()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	pgDb, err := postgres.Db.DB()
	if err != nil {
		log.Fatalf("Failed to get *pgx.DB object: %v", err)
	}

	defer func() {
		if err := pgDb.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully.")
		}
	}()

	if err := postgres.Db.AutoMigrate(&models.Question{}); err != nil {
		log.Fatalf("Failed migration: %v", err)
	}

	log.Println("Database migration finished...")
}
