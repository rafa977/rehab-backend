package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/rehab-backend/internal/pkg/handlers"
	"github.com/rehab-backend/internal/pkg/models"
)

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var response models.Response

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			response.Status = "error"
			response.Message = "Authorization Required"
			response.Response = ""
			json.NewEncoder(w).Encode(response)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			response.Status = "error"
			response.Message = "Bearer Token Required"
			response.Response = ""
			json.NewEncoder(w).Encode(response)
			return
		}

		username, err := handlers.ValidateToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			response.Status = "error"
			response.Message = "Authorization Failed"
			response.Response = ""
			json.NewEncoder(w).Encode(response)
			return
		} else {

			context.Set(r, "username", username)
			next.ServeHTTP(w, r)
		}
	}
}
