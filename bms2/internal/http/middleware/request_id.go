package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := generateRequestID()

		ctx := context.WithValue(
			r.Context(),
			requestIDKey,
			id,
		)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetRequestID(ctx context.Context) string {

	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return id
}

func generateRequestID() string {

	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
