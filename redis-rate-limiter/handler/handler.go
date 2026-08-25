package handler

import (
	"encoding/json"
	"net/http"
	"redis-rate-limiter/service"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// handler dont return anything it just writes the response to the response writer
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]

	user := h.UserService.Profile(ctx, id)

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
