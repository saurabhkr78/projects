package handler

import (
	"net/http"

	"github.com/saurabhkr78/redis-recent-views/service"
)

type Handler struct {
	svc service.ProductService
}

func NewHandler(s service.ProductService) *Handler {
	return &Handler{
		svc: s,
	}
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// 1. Get product ID from URL
	// 2. Parse it
	// 3. Call h.svc.GetProduct()
	// 4. Return JSON
}

func (h *Handler) SetRecentView(w http.ResponseWriter, r *http.Request) {
	// 1. Get userID
	// 2. Get productID
	// 3. Call h.svc.SetRecentView()
	// 4. Return response
}

func (h *Handler) GetRecentViews(w http.ResponseWriter, r *http.Request) {
	// 1. Get userID
	// 2. Call h.svc.GetRecentViews()
	// 3. Return JSON
}
