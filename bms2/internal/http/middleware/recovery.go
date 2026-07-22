package middleware

import (
	"log"
	"net/http"

	httphelper "bms2/internal/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		defer func() {

			if err := recover(); err != nil {

				log.Println(err)

				httphelper.WriteError(
					w,
					http.StatusInternalServerError,
					"internal server error",
				)

			}

		}()

		next.ServeHTTP(w, r)

	})
}
