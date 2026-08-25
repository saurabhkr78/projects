package handler

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"net/http"
)

// inject the redis client into handler struct
type Handler struct {
	//redis client instance
	redisClient *redis.Client
}

func NewHandler(redisClient *redis.Client) *Handler {
	return &Handler{
		redisClient: redisClient,
	}
}

//we need to use the context if client cancels the request or if the request times out. The context will be used to cancel the operation if the client disconnects or if the request takes too long.

func (h *Handler) Increment(w http.ResponseWriter, r *http.Request) {

	value, err := h.redisClient.Incr(r.Context(), "counter").Result()
	if err != nil {
		http.Error(w, "Failed to increment counter", http.StatusInternalServerError)
		return
	}
	// Do something with the incremented value, e.g., send it in the response
	//using writer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"counter": value})

}
func (h *Handler) Decrement(w http.ResponseWriter, r *http.Request) {
	val, err := h.redisClient.Decr(r.Context(), "counter").Result()
	if err != nil {
		http.Error(w, "Failed to decrement counter", http.StatusInternalServerError)
		return
	}
	// Do something with the decremented value, e.g., send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"counter": val})

}
func (h *Handler) GetValue(w http.ResponseWriter, r *http.Request) {
	val, err := h.redisClient.Get(r.Context(), "counter").Result()
	if err != nil {
		http.Error(w, "Failed to get counter value", http.StatusInternalServerError)
		return
	}
	// Do something with the value, e.g., send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"counter": val})

}
