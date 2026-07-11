package book

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /books", handler.CreateBook)
	mux.HandleFunc("GET /books", handler.GetAllBooks)
	mux.HandleFunc("GET /books/{id}", handler.GetBookByID)
	mux.HandleFunc("PUT /books/{id}", handler.UpdateBook)
	mux.HandleFunc("DELETE /books/{id}", handler.DeleteBook)
}
