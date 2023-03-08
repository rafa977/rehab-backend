package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/rehab-backend/internal/pkg/handlers"
)

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		username, err := handlers.ValidateToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {

			context.Set(r, "username", username)
			next.ServeHTTP(w, r)
		}
	}
}
