package book

import (
	"encoding/json"
	"net/http"

	httphelper "bms2/internal/http"
)

type Handler struct {
	service *BookService
}

func NewHandler(service *BookService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var newBook CreateBookRequest
	//call decoder ,decode it into var and write if any error occurs
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		httphelper.WriteError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}
	//call the service layer to create the book
	if err := h.service.Create(r.Context(), &newBook); err != nil {
		httphelper.WriteError(
			w,
			http.StatusInternalServerError,
			err.Error(), //it's a inbuilt interface that returns the error message as a string
		)
		return
	}
	//write the response
	httphelper.WriteJSON(
		w,
		http.StatusCreated,
		map[string]string{
			"message": "Book created successfully",
		},
	)

}

func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.ParseIntPathParam(r, "id")
	if err != nil {
		httphelper.WriteError(
			w,
			http.StatusBadRequest,
			"invalid book id",
		)
		return
	}

	book, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httphelper.WriteError(
			w,
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	httphelper.WriteJSON(
		w,
		http.StatusOK,
		book,
	)
}

func (h *Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetAll(r.Context())
	if err != nil {
		httphelper.WriteError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	httphelper.WriteJSON(
		w,
		http.StatusOK,
		books,
	)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.ParseIntPathParam(r, "id")
	if err != nil {
		httphelper.WriteError(
			w,
			http.StatusBadRequest,
			"invalid book id",
		)
		return
	}

	var req UpdateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelper.WriteError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if err := h.service.Update(r.Context(), id, &req); err != nil {
		httphelper.WriteError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	httphelper.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "book updated successfully",
		},
	)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := httphelper.ParseIntPathParam(r, "id")
	if err != nil {
		httphelper.WriteError(
			w,
			http.StatusBadRequest,
			"invalid book id",
		)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httphelper.WriteError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
