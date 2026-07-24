package middleware

import (
	"bms2/internal/requestid"
	"bms2/pkg/logger"
	"net/http"
	"time"
)

// refer line number 2465 in learning.md
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := requestid.FromContext(r.Context())
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		logger.Info(
			"[RequestID=%s] Started %s %s",
			requestID,
			r.Method,
			r.URL.Path,
		)

		next.ServeHTTP(rw, r)

		logger.Info(
			"[RequestID=%s] Completed %d in %v",
			requestID,
			rw.statusCode,
			time.Since(start),
		)
	})
}
