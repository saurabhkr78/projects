package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"redis-url-shortner/client"
	"redis-url-shortner/db"
	"redis-url-shortner/handler"
	"redis-url-shortner/repository"
	"redis-url-shortner/service"
)

func main() {

	// 1. Initialize Redis
	redisClient := client.NewClient()

	// 2. Initialize repository with fake DB
	repo := repository.NewDBRepository(db.URLs)

	// 3. Initialize service
	urlService := service.NewURLService(
		repo,
		redisClient,
	)

	// 4. Initialize handler
	urlHandler := handler.NewURLHandler(urlService)

	// 5. Initialize router
	r := mux.NewRouter()

	// Create short URL
	r.HandleFunc("/shorten", urlHandler.Shorten).Methods("POST")

	// Redirect using short ID
	r.HandleFunc("/{shortID}", urlHandler.Redirect).Methods("GET")

	// 6. Start server
	log.Println("Server running on :8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
