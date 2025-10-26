package main

import (
	"live-rescue/internal/database"
	"live-rescue/internal/handlers"
	"live-rescue/internal/repositories"
	"log"
	"net/http"
	"os"

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

	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(os.Getenv("FLAG_STORAGE_PATH")))))

	http.HandleFunc("/", indexHandler.Index)
	http.HandleFunc("/questions/create", indexHandler.Create)
	http.HandleFunc("/questions", questionHandler.Handle)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
