package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorDetail struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

func WriteError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Message: message,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}
