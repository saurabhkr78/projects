package main

import (
	"book-api/internal/book"
	"book-api/pkg/config"
	"book-api/pkg/database"
	"log"
	"net/http"
)

func main() {
	//laod the config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	//connect to the database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Step 3
	repo := book.NewRepository(db)

	// Step 4
	service := book.NewService(repo)

	// Step 5
	handler := book.NewHandler(service)

	// Step 6
	mux := http.NewServeMux()

	book.RegisterRoutes(mux, handler)

	log.Println("Server running on :", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}

}
