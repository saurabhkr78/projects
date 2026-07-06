//How do I assemble the whole application?
//assemble the application (also called composition root or dependency injection).

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"mongosetup/config"
	"mongosetup/database"
	"mongosetup/internals/book"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to MongoDB
	client, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(nil)

	// Get database
	db := client.Database(cfg.DatabaseName)

	// Get collection
	booksCollection := db.Collection("books")

	// Create repository
	repository := book.NewRepository(booksCollection)

	// Create service
	service := book.NewService(repository)

	// Create handler
	handler := book.NewHandler(service)

	// Create router
	router := mux.NewRouter()

	// Register routes
	book.RegisterRoutes(router, handler)

	log.Printf("Server running on port %s", cfg.Port)

	// Start server
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
