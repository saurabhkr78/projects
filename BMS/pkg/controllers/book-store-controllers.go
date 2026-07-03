package controllers

import (
	"BMS/pkg/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func homepage(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome to the BMS API")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the BMS API"))
}

var NewBook models.Book

func CreateBook(w http.ResponseWriter, r *http.Request) {}
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
func GetBookById(w http.ResponseWriter, r *http.Request) {}
func UpdateBook(w http.ResponseWriter, r *http.Request)  {}
func DeleteBook(w http.ResponseWriter, r *http.Request)  {}
