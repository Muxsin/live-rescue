package main

import (
	"live-rescue/internal/database"
	"live-rescue/internal/handlers"
	"live-rescue/internal/repositories"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgres, err := database.NewPostgres()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	mediaRepository := repositories.NewMedia()
	questionRepository := repositories.NewQuestion(postgres.Db)

	indexHandler := handlers.NewView()
	questionHandler := handlers.NewQuestion(mediaRepository, questionRepository)

	http.HandleFunc("/", indexHandler.Index)
	http.HandleFunc("/questions", questionHandler.Create)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
