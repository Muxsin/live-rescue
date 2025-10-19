package handlers

import "net/http"

type viewHandler struct {
}

func NewView() *viewHandler {
	return &viewHandler{}
}

func (h *viewHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "./resources/views/index.html")
}
