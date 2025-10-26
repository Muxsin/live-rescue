package handlers

import (
	"html/template"
	"net/http"
)

type viewHandler struct {
	template *template.Template
}

func NewView(template *template.Template) *viewHandler {
	return &viewHandler{
		template: template,
	}
}

func (h *viewHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := h.template.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *viewHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := h.template.ExecuteTemplate(w, "create", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
