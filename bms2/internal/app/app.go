package app

import (
	"bms2/internal/book"
	"bms2/pkg/config"
	"bms2/pkg/database"
	"log"
	"net/http"
)

type App struct {
	server *http.Server
	db     *database.DB
}

// New builds the entire application and returns it.
func New() *App {

	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database.
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create application dependencies.
	repo := book.NewRepository(db)
	service := book.NewService(repo)
	handler := book.NewHandler(service)

	// Register routes.
	router := http.NewServeMux()
	book.RegisterRoutes(router, handler)

	// Create HTTP server.
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	return &App{
		server: server,
	}
}

// Run starts the HTTP server.
func (a *App) Run() error {
	log.Printf("Server listening on %s", a.server.Addr)

	return a.server.ListenAndServe()
}
