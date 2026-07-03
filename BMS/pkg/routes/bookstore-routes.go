package routes

import (
	"BMS/pkg/controllers"
	"github.com/gorilla/mux"
)

func RegisterBookStoreRoutes(r *mux.Router) {
	r.HandleFunc("/book/", controllers.CreateBook).Methods("POST")
	r.HandleFunc("/book/", controllers.GetBooks).Methods("GET")
	r.HandleFunc("/book/{bookId}", controllers.GetBookById).Methods("GET")
	r.HandleFunc("/book/{bookId}", controllers.UpdateBook).Methods("PUT")
	r.HandleFunc("/book/{id}", controllers.DeleteBook).Methods("DELETE")
}
