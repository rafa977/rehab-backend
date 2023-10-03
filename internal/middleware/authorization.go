package middleware

import (
	"encoding/json"
	"net/http"
	"rehab/internal/pkg/handlers"
	"rehab/internal/pkg/models"
	"strings"
	"time"

	"github.com/gorilla/context"
)

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var response models.Response
		currentDate := time.Now().Format("2006-01-02 15:04:05")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			response.Status = "error"
			response.Message = "Authorization Required"
			response.Response = ""
			response.Date = currentDate
			json.NewEncoder(w).Encode(response)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			response.Status = "error"
			response.Message = "Bearer Token Required"
			response.Response = ""
			response.Date = currentDate
			json.NewEncoder(w).Encode(response)
			return
		}

		username, id, compIDs, roleID, err := handlers.ValidateToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			response.Status = "error"
			response.Message = "Authorization Failed"
			response.Response = ""
			response.Date = currentDate
			json.NewEncoder(w).Encode(response)
			return
		} else {

			context.Set(r, "username", username)
			context.Set(r, "id", id)
			context.Set(r, "compIDs", compIDs)
			context.Set(r, "roleID", roleID)
			next.ServeHTTP(w, r)
		}
	}
}
