package handlers

import (
	"encoding/json"
	"live-rescue/internal/models"
	"live-rescue/internal/repositories"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type questionHandler struct {
	mediaRepo    *repositories.MediaRepository
	questionRepo *repositories.QuestionRepository
}

func NewQuestion(
	mediaRepo *repositories.MediaRepository,
	quetionRepo *repositories.QuestionRepository,
) *questionHandler {
	return &questionHandler{
		mediaRepo:    mediaRepo,
		questionRepo: quetionRepo,
	}
}

func (h *questionHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	description := r.FormValue("description")

	filePath := filepath.Join(
		os.Getenv("FLAG_STORAGE_PATH"),
		strconv.Itoa(int(time.Now().UnixMilli()))+header.Filename,
	)

	if err := h.mediaRepo.StoreFile(file, header, filePath); err != nil {
		log.Printf("Failed to store file: %v", err.Error())

		http.Error(w, "Failed to store file", http.StatusInternalServerError)
		return
	}

	question := &models.Question{
		Title:       title,
		Description: description,
		ImagePath:   filePath,
	}

	if err := h.questionRepo.Create(question); err != nil {
		log.Printf("Failed to create question: %v", err.Error())

		h.mediaRepo.DeleteFile(filePath)
		http.Error(w, "Failed to create question", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}
