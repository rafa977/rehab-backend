package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rehab-backend/internal/pkg/models"
)

type Service struct {
	// queries *database.Queries
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) RegisterHandlers(route *mux.Router) {

	Handle(route)

}

func Handle(route *mux.Router) {
	sub := route.PathPrefix("/user").Subrouter()

	sub.HandleFunc("/registerAccount", accountRegistration)
	sub.HandleFunc("/login", login)
	sub.HandleFunc("/getAccountById", getAccountById)
	sub.HandleFunc("/deleteAccountById", deleteAccountById)

}

func accountRegistration(w http.ResponseWriter, r *http.Request) {

	var account models.Account

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if account.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	if account.ID == "" {
		http.Error(w, "ID cannot be empty", http.StatusBadRequest)
		return
	}

	// Do something with the Person struct...
	fmt.Fprintf(w, "Account: %+v", account.Username)
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login!")
}

func getAccountById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get account data!")
}

func deleteAccountById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete account!")
}
