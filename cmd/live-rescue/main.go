package main

import (
	"html/template"
	"live-rescue/internal/database"
	"live-rescue/internal/handlers"
	"live-rescue/internal/repositories"
	"log"
	"net/http"
	"os"
	"path/filepath"

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

	temp, err := parseTemplates("resources/views")
	if err != nil {
		log.Fatalf("Failed parsing templates: %v", err)
	}

	mediaRepository := repositories.NewMedia()
	questionRepository := repositories.NewQuestion(postgres.Db)

	viewHandler := handlers.NewView(temp)
	questionHandler := handlers.NewQuestion(mediaRepository, questionRepository)

	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(os.Getenv("FLAG_STORAGE_PATH")))))

	http.HandleFunc("/", viewHandler.Index)
	http.HandleFunc("/questions/{id}", questionHandler.Get)
	http.HandleFunc("/questions/create", viewHandler.Create)
	http.HandleFunc("/questions", questionHandler.Handle)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseTemplates(dir string) (*template.Template, error) {
	t := template.New("")

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".html" {
			bytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = t.New(filepath.ToSlash(path)).Parse(string(bytes))
			if err != nil {
				return err
			}
		}

		return nil
	})

	return t, err
}
