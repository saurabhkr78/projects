package main

import (
	"bms2/internal/auth"
	"bms2/internal/book"
	"bms2/internal/config"
	"bms2/internal/database"
	"bms2/internal/http/middleware"
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	} else if r.Method != http.MethodGet {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	} else {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "world"
		}
		fmt.Fprintf(w, "hello %s", name)
	}
}
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/form" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	} else if r.Method != http.MethodPost {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form data", http.StatusBadRequest)
			return
		}
		name := r.FormValue("name")
		email := r.FormValue("email")

		fmt.Fprintf(w, "Hello %s, your email is %s", name, email)
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := book.NewRepository(db)
	service := book.NewService(repo)
	handler := book.NewHandler(service)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	router := http.NewServeMux()

	book.RegisterRoutes(router, handler)

	app := middleware.Chain(
		router,
		middleware.RequestID,
		middleware.Authentication(jwtManager),
		middleware.Logging,
		middleware.Recovery,
	)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: app,
	}

	log.Printf("Server listening on :%s", cfg.Port)

	log.Fatal(server.ListenAndServe())
}
