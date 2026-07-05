// Which URL calls which handler?
package book

import (
	"github.com/gorilla/mux"
	"net/http"
)

/*

Client
   │
   ▼
Router (Receptionist)
   │
   ├── POST /books  ─────► CreateBook()
   ├── GET  /books  ─────► GetBooks()
   ├── GET  /books/{id} ─► GetBookByID()
   ├── PUT  /books/{id} ─► UpdateBook()
   └── DELETE /books/{id} ─► DeleteBook()
*/
/*Instead of registering routes directly on the root router, group them using a subrouter*/
func RegisterRoutes(r *mux.Router, h *Handler) {
	bookRouter := r.PathPrefix("/books").Subrouter()

	bookRouter.HandleFunc("", h.CreateBook).Methods("POST")
	bookRouter.HandleFunc("", h.GetBooks).Methods("GET")
	bookRouter.HandleFunc("/{id}", h.GetBookByID).Methods("GET")
	bookRouter.HandleFunc("/{id}", h.UpdateBook).Methods("PUT")
	bookRouter.HandleFunc("/{id}", h.DeleteBook).Methods("DELETE")
}
