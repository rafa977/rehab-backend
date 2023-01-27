package accounts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	currentDate := time.Now().Format("2006-01-02 15:04:05")

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

var sampleSecretKey = []byte("SecretYouShouldHide")

func generateJWT(username string) (string, time.Time, error) {

	expTime := time.Now().Add(time.Minute * 60)

	claims := &Claims{
		Username:   username,
		Authorized: "authorized",
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(claims)

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		fmt.Println("here is the error")
		return "", expTime, err
	}

	return tokenString, expTime, nil

}

func validateToken(w http.ResponseWriter, r *http.Request) (string, error) {

	if r.Header["Authorization"] == nil {
		fmt.Fprintf(w, "Authorization is required")
		return "", errors.New("Authorization is required")
	}

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		// Error: Bearer token not in proper format
	}

	reqToken = strings.TrimSpace(splitToken[1])

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	username := claims.Username
	fmt.Println(username)

	fmt.Println(err)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return "", errors.New("Invalid signature")
		}
		w.WriteHeader(http.StatusBadRequest)
		return "", errors.New("Token is expired")
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return "", errors.New("Token is invalid")
	}

	return username, nil
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

		retrievedAccount.LastLogin = currentDate

		_, err = s.dbConnection.NewUpdate().Model(&retrievedAccount).Column("last_login").Where("username = ?", account.Username).Exec(ctx)

		token, expTime, hasError := generateJWT(account.Username)
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

	username, err := validateToken(w, r)
	if err != nil {
		response.Status = "error"
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)

		return
	}

	ctx := context.Background()

	err = s.dbConnection.NewSelect().Model(&account).Where("user_id = ?", id).Scan(ctx, &retrievedAccount)
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
