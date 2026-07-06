package book

import (
	"encoding/json"
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

	//MAKE AN EMPTY BOOK OBJECT
	var book Book
	//DECODE THE JSON BODY INTO THE BOOK OBJECT
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//GET THE CONTEXT FROM THE REQUEST
	ctx := r.Context()
	//CALL THE SERVICE TO CREATE THE BOOK
	createdBook, err := h.service.CreateBook(ctx, &book)
	if err != nil {
		http.Error(w, "Failed to create book: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//SET THE RESPONSE HEADER AND STATUS CODE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//ENCODE THE CREATED BOOK AS JSON AND WRITE IT TO THE RESPONSE
	if err := json.NewEncoder(w).Encode(createdBook); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

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
