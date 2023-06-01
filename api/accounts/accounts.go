package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rehab-backend/internal/middleware"
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
	sub.HandleFunc("/updateAccount", middleware.AuthenticationMiddleware(s.updateAccount))
	sub.HandleFunc("/updatePassword", middleware.AuthenticationMiddleware(s.updatePassword))
	sub.HandleFunc("/login", s.login)
	sub.HandleFunc("/getAccountById", middleware.AuthenticationMiddleware(s.getAccountById))
	sub.HandleFunc("/getCompaniesByAccountId", middleware.AuthenticationMiddleware(s.getCompaniesByAccountId))
	sub.HandleFunc("/deleteAccountById", deleteAccountById)
}

func (s *service) accountRegistration(w http.ResponseWriter, r *http.Request) {
	var account models.Account

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	isValid, errors := handlers.ValidateInputs(account)
	if !isValid {
		for _, fieldError := range errors {
			handlers.ProduceErrorResponse(fieldError, w, r)
			return
		}
	}

	hashedPassword, hashError := hashPassword(account.Password)
	if hashError != nil {
		handlers.ProduceErrorResponse(hashError.Error(), w, r)
		return
	}

	account.Password = hashedPassword

	// birthDate, timeParsError := time.Parse("2006-01-02", account.Birthdate.String())
	// if timeParsError != nil {
	// 	handlers.ProduceErrorResponse(timeParsError.Error(), w, r)
	// }

	account, err = s.repository.AddUser(account)
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "email_key") {
			msg = "User already exists!"
		} else if strings.Contains(err.Error(), "username_key") {
			msg = "User already exists!"
		} else {
			msg = "Bad Request"
		}
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Registration of Account - Successful", w, r)
}

func (s *service) updatePassword(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	var retrievedAccount models.Account

	id := gcontext.Get(r, "id").(uint)

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	isValid, errors := handlers.ValidateInputs(credentials)
	if !isValid {
		for _, fieldError := range errors {
			handlers.ProduceErrorResponse(fieldError, w, r)
			return
		}
	}

	retrievedAccount, err = s.repository.GetAccountByIDWithPassword(id)
	if err != nil {
		handlers.ProduceErrorResponse("You are not authorized to access this data", w, r)
		return
	}

	hashedPassword, hashError := hashPassword(credentials.NewPassword)
	if hashError != nil {
		handlers.ProduceErrorResponse(hashError.Error(), w, r)
		return
	}
	passwordValid := checkPasswordHash(credentials.OldPassword, retrievedAccount.Password)

	if passwordValid {
		retrievedAccount.Password = hashedPassword
	} else {
		handlers.ProduceErrorResponse("Old password does not match", w, r)
		return
	}

	_, err = s.repository.UpdateAccount(retrievedAccount)
	if err != nil {
		msg := "Bad Request"
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Update of Password - Successful", w, r)
}

func (s *service) updateAccount(w http.ResponseWriter, r *http.Request) {
	var updatedAccount models.AccountInt
	var retrievedAccount models.Account

	id := gcontext.Get(r, "id").(uint)

	err := json.NewDecoder(r.Body).Decode(&updatedAccount)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	isValid, errors := handlers.ValidateInputs(updatedAccount)
	if !isValid {
		for _, fieldError := range errors {
			handlers.ProduceErrorResponse(fieldError, w, r)
			return
		}
	}

	// account.Password = hashedPassword

	retrievedAccount, err = s.repository.GetAccountByID(id)
	if err != nil {
		handlers.ProduceErrorResponse("You are not authorized to access this data", w, r)
		return
	}

	retrievedAccount.Firstname = updatedAccount.Firstname
	retrievedAccount.Lastname = updatedAccount.Lastname
	retrievedAccount.Birthdate = updatedAccount.Birthdate
	retrievedAccount.Address = updatedAccount.Address
	retrievedAccount.Job = updatedAccount.Job

	_, err = s.repository.UpdateAccount(retrievedAccount)
	if err != nil {
		msg := "Bad Request"
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Update of Account - Successful", w, r)
}

func (s *service) login(w http.ResponseWriter, r *http.Request) {
	currentDate := time.Now().Format("2006-01-02 15:04:05")

	var account models.Account
	var retrievedAccount models.Account
	var response models.Response
	response.Date = currentDate

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if account.Username == "" {
		handlers.ProduceErrorResponse("Username cannot be empty", w, r)
		return
	}
	if account.Password == "" {
		handlers.ProduceErrorResponse("Password cannot be empty", w, r)
		return
	}

	// ctx := context.Background()

	retrievedAccount, err = s.repository.GetAccountByUsernameForLogin(account.Username)
	if err != nil {
		handlers.ProduceErrorResponse("Unknown Username or Password", w, r)
		return
	}

	fmt.Println(account.Password)

	passwordValid := checkPasswordHash(account.Password, retrievedAccount.Password)
	fmt.Println(passwordValid)

	if passwordValid {

		retrievedAccount.LastLogin = time.Now()

		_, _ = s.repository.UpdateAccount(retrievedAccount)
		// _, err = s.dbConnection.NewUpdate().Model(&retrievedAccount).Column("last_login").Where("username = ?", account.Username).Exec(ctx)

		retrievedCompanies := s.repository.GetCompaniesByAccountID(retrievedAccount.ID)

		token, expTime, hasError := handlers.GenerateJWT(account.Username, retrievedAccount.ID, retrievedCompanies)
		if hasError != nil {
			handlers.ProduceErrorResponse(hasError.Error(), w, r)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: expTime,
		})

		retrievedAccount.Password = ""

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

func (s *service) getCompaniesByAccountId(w http.ResponseWriter, r *http.Request) {
	var retrievedCompanies []uint

	id := gcontext.Get(r, "id").(uint)

	retrievedCompanies = s.repository.GetCompaniesByAccountID(id)

	jsonRetrievedAccount, err := json.Marshal(retrievedCompanies)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) getAccountById(w http.ResponseWriter, r *http.Request) {

	var response models.Response
	var retrievedAccount models.Account

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	uintID := uint(id)

	username := gcontext.Get(r, "username").(string)

	retrievedAccount, err = s.repository.GetAccountByID(uintID)
	if err != nil {
		handlers.ProduceErrorResponse("You are not authorized to access this data", w, r)
		return
	}

	if retrievedAccount.Username != username {
		handlers.ProduceErrorResponse("You are not authorized to view this data", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(retrievedAccount)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
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
