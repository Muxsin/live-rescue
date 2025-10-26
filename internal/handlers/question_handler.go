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

func (h *questionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
}

func (h *questionHandler) filePathToURL(filePath string) string {
	filename := filepath.Base(filePath)
	return "/media/" + filename
}

func (h *questionHandler) List(w http.ResponseWriter, r *http.Request) {
	questions, err := h.questionRepo.GetAll()
	if err != nil {
		log.Printf("Listing questions failed: %v", err.Error())
		http.Error(w, "Listing questions failed", http.StatusInternalServerError)
		return
	}

	type QuestionResponse struct {
		Title       string `json:"Title"`
		Description string `json:"Description"`
		ImageURL    string `json:"ImageURL"`
	}

	var response []QuestionResponse
	for _, q := range questions {
		imageURL := h.filePathToURL(q.ImagePath)

		response = append(response, QuestionResponse{
			Title:       q.Title,
			Description: q.Description,
			ImageURL:    imageURL,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *questionHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
