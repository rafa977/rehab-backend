package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/handlers"
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	dbConnection *bun.DB
}

type Claims struct {
	Username   string `json:"username"`
	Authorized string `json:"authorized"`
	jwt.RegisteredClaims
}

func NewService() *Service {
	dbConnection := config.ConnectDB()

	// return &Service{}
	return &Service{dbConnection: dbConnection}

}

func (s *Service) RegisterHandlers(route *mux.Router) {

	s.Handle(route)

}

func (s *Service) Handle(route *mux.Router) {
	sub := route.PathPrefix("/user").Subrouter()

	sub.HandleFunc("/registerAccount", s.accountRegistration)
	sub.HandleFunc("/login", s.login)
	sub.HandleFunc("/getAccountById", s.getAccountById)
	sub.HandleFunc("/deleteAccountById", deleteAccountById)

}

func (s *Service) accountRegistration(w http.ResponseWriter, r *http.Request) {
	// currentDate := time.Now().Format("2006-01-02 15:04:05")

	var account models.Account

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(account)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	hashedPassword, hashError := hashPassword(account.Password)
	if hashError != nil {
		http.Error(w, hashError.Error(), http.StatusBadRequest)
		return
	}

	account.Password = hashedPassword
	account.CreatedOn = time.Now()

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
	w.Header().Set("Content-Type", "application/json")

	currentDate := time.Now().Format("2006-01-02 15:04:05")

	var account models.Account
	var retrievedAccount models.Account
	var response models.Response
	response.Date = currentDate

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

	ctx := context.Background()

	err = s.dbConnection.NewSelect().Model(&account).Where("username = ?", account.Username).Scan(ctx, &retrievedAccount)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordValid := checkPasswordHash(account.Password, retrievedAccount.Password)

	if passwordValid {

		retrievedAccount.LastLogin = time.Now()

		_, err = s.dbConnection.NewUpdate().Model(&retrievedAccount).Column("last_login").Where("username = ?", account.Username).Exec(ctx)

		token, expTime, hasError := handlers.GenerateJWT(account.Username)
		if hasError != nil {
			http.Error(w, hasError.Error(), http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: expTime,
		})

		response.Response = token
		response.Status = "success"
		response.Message = ""
		json.NewEncoder(w).Encode(response)

		return
	} else {
		response.Status = "error"
		response.Message = "Authorization failed."
		json.NewEncoder(w).Encode(response)

		return
	}

}

func (s *Service) getAccountById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var account models.Account
	var response models.Response
	var retrievedAccount models.Account

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		response.Status = "error"
		response.Message = "Please input all required fields."
		json.NewEncoder(w).Encode(response)

		return
	}

	// username, err := handlers.ValidateToken(w, r)
	// if err != nil {
	// 	response.Status = "error"
	// 	response.Message = err.Error()
	// 	json.NewEncoder(w).Encode(response)

	// 	return
	// }
	username := gcontext.Get(r, "username").(string)

	ctx := context.Background()

	err := s.dbConnection.NewSelect().Model(&account).Where("user_id = ?", id).Scan(ctx, &retrievedAccount)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if retrievedAccount.Username != username {
		http.Error(w, "You are not authorized to view this data.", http.StatusBadRequest)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(retrievedAccount)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = ""
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)

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
