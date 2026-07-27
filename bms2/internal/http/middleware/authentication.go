package middleware

import (
	"errors"
	"net/http"
	"strings"

	"bms2/internal/auth"
	httphelper "bms2/internal/httphelper"
)

func Authentication(jwt *auth.JWTManager) Middleware {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("Authorization")
			if header == "" {
				httphelper.WriteError(
					w,
					http.StatusUnauthorized,
					"missing authorization header",
				)
				return
			}

			const prefix = "Bearer "

			if !strings.HasPrefix(header, prefix) {
				httphelper.WriteError(
					w,
					http.StatusUnauthorized,
					"invalid authorization header",
				)
				return
			}

			tokenString := strings.TrimPrefix(header, prefix)

			claims, err := jwt.ValidateToken(tokenString)
			if err != nil {

				switch {

				case errors.Is(err, auth.ErrExpiredToken):
					httphelper.WriteError(
						w,
						http.StatusUnauthorized,
						"token expired",
					)

				case errors.Is(err, auth.ErrMalformedToken):
					httphelper.WriteError(
						w,
						http.StatusBadRequest,
						"malformed token",
					)

				default:
					httphelper.WriteError(
						w,
						http.StatusUnauthorized,
						"invalid token",
					)
				}

				return
			}

			user := auth.User{
				ID:    claims.UserID,
				Email: claims.Email,
				Role:  claims.Role,
			}

			ctx := auth.IntoContext(r.Context(), user)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
