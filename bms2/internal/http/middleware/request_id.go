package middleware

import (
	"net/http"

	"bms2/internal/requestid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := requestid.GenerateNewId()

		ctx := requestid.IntoContext(
			r.Context(),
			id,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}
