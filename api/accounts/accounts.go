package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	dbConnection *bun.DB
}

func NewService() *Service {
	dbConnection := config.ConnectDB()

	return &Service{dbConnection: dbConnection}
}

func (s *Service) RegisterHandlers(route *mux.Router) {

	s.Handle(route)

}

func (s *Service) Handle(route *mux.Router) {
	sub := route.PathPrefix("/user").Subrouter()

	sub.HandleFunc("/registerAccount", s.accountRegistration)
	sub.HandleFunc("/login", s.login)
	sub.HandleFunc("/getAccountById", getAccountById)
	sub.HandleFunc("/deleteAccountById", deleteAccountById)

}

func (s *Service) accountRegistration(w http.ResponseWriter, r *http.Request) {
	currentDate := time.Now().Format("2006-01-02 15:04:05")

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
	if account.Password == "" {
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	}

	hashedPassword, hashError := hashPassword(account.Password)
	if hashError != nil {
		http.Error(w, hashError.Error(), http.StatusBadRequest)
		return
	}

	account.Password = hashedPassword
	account.CreatedOn = currentDate

	ctx := context.Background()

	_, err = s.dbConnection.NewInsert().Model(&account).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Registration of Account - Successful")
}

func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	currentDate := time.Now().Format("2006-01-02 15:04:05")

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
	if account.Password == "" {
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	}

	var retrievedAccount models.Account

	ctx := context.Background()

	err = s.dbConnection.NewSelect().Model(&account).Where("username = ?", account.Username).Scan(ctx, &retrievedAccount)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordValid := checkPasswordHash(account.Password, retrievedAccount.Password)

	if passwordValid {

		retrievedAccount.LastLogin = currentDate

		_, err = s.dbConnection.NewUpdate().Model(&retrievedAccount).Column("last_login").Where("username = ?", account.Username).Exec(ctx)

		fmt.Fprintf(w, "Login")

	} else {
		fmt.Fprintf(w, "Login Failed")
	}
}

func getAccountById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get account data!")
}

func deleteAccountById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete account!")
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
