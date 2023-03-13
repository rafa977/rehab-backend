package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/rehab-backend/internal/pkg/handlers"
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/rehab-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository repository.AccountsRepository
}

type Claims struct {
	Username   string `json:"username"`
	Authorized string `json:"authorized"`
	jwt.RegisteredClaims
}

func NewService() *service {
	return &service{repository: repository.NewAccountService()}
}

func (s *service) RegisterHandlers(route *mux.Router) {

	s.Handle(route)

}

func (s *service) Handle(route *mux.Router) {
	sub := route.PathPrefix("/user").Subrouter()

	sub.HandleFunc("/registerAccount", s.accountRegistration)
	sub.HandleFunc("/login", s.login)
	sub.HandleFunc("/getAccountById", s.getAccountById)
	sub.HandleFunc("/deleteAccountById", deleteAccountById)

}

func (s *service) accountRegistration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	var response models.Response

	account, err = s.repository.AddUser(account)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "users_company_email_key") {
			newerr = "user already exists!"
		} else {
			newerr = "Bad Request"
		}
		response.Status = "error"
		response.Message = newerr
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Fprintf(w, "Registration of Account - Successful")
}

func (s *service) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	// ctx := context.Background()

	retrievedAccount, err = s.repository.GetAccountByID(1)
	if err != nil {
		response.Status = "error"
		response.Message = "Unknown Username or Password"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println(account.Password)

	passwordValid := checkPasswordHash(account.Password, retrievedAccount.Password)
	fmt.Println(passwordValid)

	if passwordValid {

		retrievedAccount.LastLogin = time.Now()

		_, _ = s.repository.UpdateAccount(retrievedAccount)
		// _, err = s.dbConnection.NewUpdate().Model(&retrievedAccount).Column("last_login").Where("username = ?", account.Username).Exec(ctx)

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

		jsonRetrievedAccount, err := json.Marshal(retrievedAccount)
		if err != nil {
			fmt.Println(err)
			return
		}

		response.Response = token
		response.Status = "success"
		response.Message = string(jsonRetrievedAccount)
		json.NewEncoder(w).Encode(response)

		return
	} else {
		response.Status = "error"
		response.Message = "Authorization failed."
		json.NewEncoder(w).Encode(response)

		return
	}

}

func (s *service) getAccountById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// var account models.Account
	var response models.Response
	// var retrievedAccount models.Account

	// currentDate := time.Now().Format("2006-01-02 15:04:05")
	// response.Date = currentDate

	// id := r.URL.Query().Get("id")
	// if id == "" {
	// 	response.Status = "error"
	// 	response.Message = "Please input all required fields."
	// 	json.NewEncoder(w).Encode(response)

	// 	return
	// }

	// username, err := handlers.ValidateToken(w, r)
	// if err != nil {
	// 	response.Status = "error"
	// 	response.Message = err.Error()
	// 	json.NewEncoder(w).Encode(response)

	// 	return
	// }
	// username := gcontext.Get(r, "username").(string)

	// ctx := context.Background()

	// err := s.dbConnection.NewSelect().Model(&account).Where("user_id = ?", id).Scan(ctx, &retrievedAccount)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// if retrievedAccount.Username != username {
	// 	http.Error(w, "You are not authorized to view this data.", http.StatusBadRequest)
	// 	return
	// }

	// jsonRetrievedAccount, err := json.Marshal(retrievedAccount)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	response.Status = "success"
	response.Message = ""
	// response.Response = string(jsonRetrievedAccount)
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
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
