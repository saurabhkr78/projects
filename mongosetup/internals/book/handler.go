package book

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

	ctx := r.Context()

	books, err := h.service.GetBooks(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET /books/{id}
func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {

	// Get the book ID from the URL.
	id := mux.Vars(r)["id"]

	ctx := r.Context()

	book, err := h.service.GetBookByID(ctx, id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// PUT /books/{id}
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	updatedBook, err := h.service.UpdateBook(ctx, id, &book)
	if err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(updatedBook); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DELETE /books/{id}
func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	ctx := r.Context()

	if err := h.service.DeleteBook(ctx, id); err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
