package controllers

import (
	"BMS/pkg/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func homepage(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome to the BMS API")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the BMS API"))
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	createdBook, err := models.CreateBook(&book)
	if err != nil {
		log.Printf("book creation failed: %v", err)
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdBook)
}
func GetBooks(w http.ResponseWriter, r *http.Request) {
	books := models.GetAllBooks()
	if len(books) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No books found"))
		return
	}
	// Serialize and return the books
	response, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error serializing books"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func GetBookById(w http.ResponseWriter, r *http.Request) {
	// Extract path variables from the URL.
	vars := mux.Vars(r)

	// Get the "bookId" value from the route: /book/{bookId}
	bookIDStr := vars["bookId"]

	// Convert the string ID to uint.
	bookID64, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		log.Printf("invalid book ID %q: %v", bookIDStr, err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	bookID := uint(bookID64)

	// Fetch the book from the repository.
	book, err := models.FindBookByID(bookID)
	if err != nil {
		log.Printf("failed to find book with ID %d: %v", bookID, err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Send a successful JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	// Extract the book ID from the URL.
	vars := mux.Vars(r)
	bookIDStr := vars["bookId"]

	// Convert the ID to uint.
	bookID64, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		log.Printf("invalid book ID %q: %v", bookIDStr, err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	bookID := uint(bookID64)

	// Decode the request body.
	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		log.Printf("invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple validation.
	if updatedBook.Name == "" {
		http.Error(w, "Book name is required", http.StatusBadRequest)
		return
	}

	// Delegate the database work.
	book, err := models.UpdateBookByID(bookID, &updatedBook)
	if err != nil {
		log.Printf("failed to update book %d: %v", bookID, err)
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	// Success response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	// Extract the path variable.
	vars := mux.Vars(r)

	// Get the book ID from the URL.
	bookIDStr := vars["bookId"]

	// Convert the ID from string to uint.
	bookID64, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		log.Printf("invalid book ID %q: %v", bookIDStr, err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	bookID := uint(bookID64)

	// Delete the book.
	err = models.DeleteBookByID(bookID)
	if err != nil {
		log.Printf("failed to delete book %d: %v", bookID, err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	log.Printf("Book %d deleted successfully", bookID)

	// Return 204 No Content.
	w.WriteHeader(http.StatusNoContent)
}
