package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"redis-url-shortner/service"
)

type URLHandler struct {
	URLserv *service.URLService
}

func NewURLHandler(urlServ *service.URLService) *URLHandler {
	return &URLHandler{
		URLserv: urlServ,
	}
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	shortID := vars["shortID"]

	url, err := h.URLserv.GetOriginalURL(ctx, shortID)
	if err != nil {
		http.Error(w, "Failed to get original URL", http.StatusInternalServerError)
		return
	}
	if url == "" {
		http.Error(w, "Failed to get original URL", http.StatusNotFound)
		return
	}
	//if found then redirect to the original url
	http.Redirect(w, r, url, http.StatusFound)
}
func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		OriginalURL string `json:"original_url"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OriginalURL == "" {
		http.Error(w, "Original URL is required", http.StatusBadRequest)
		return
	}

	shortURL, err := h.URLserv.SaveURL(ctx, req.OriginalURL)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message":   "URL shortened successfully",
		"short_url": shortURL,
	})
}
