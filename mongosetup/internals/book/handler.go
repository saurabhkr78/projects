package book

import (
	"net/http"
)

/*
Handler handles all HTTP requests related to books.

Responsibilities:
- Read request
- Decode JSON
- Read URL params
- Call Service
- Write HTTP response

It should NOT:
- Talk to MongoDB
- Contain business logic
*/
type Handler struct {
	service *Service
}

// NewHandler creates a new Handler.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// POST /books
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {

}

// GET /books
func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {

}

// GET /books/{id}
func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {

}

// PUT /books/{id}
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {

}

// DELETE /books/{id}
func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {

}
