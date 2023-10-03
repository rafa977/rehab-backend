package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"rehab/internal/middleware"
	"rehab/internal/pkg/handlers"
	"rehab/internal/pkg/models"
	"rehab/internal/repository"
	"rehab/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository        repository.AccountsRepository
	companyRepository repository.CompanyRepository
}

type Claims struct {
	Username   string `json:"username"`
	Authorized string `json:"authorized"`
	jwt.RegisteredClaims
}

func NewService() *service {
	return &service{
		repository:        repository.NewAccountService(),
		companyRepository: repository.NewCompanyService(),
	}
}

func (s *service) RegisterHandlers(route *mux.Router) {
	s.Handle(route)
}

func (s *service) Handle(route *mux.Router) {
	sub := route.PathPrefix("/user").Subrouter()

	sub.HandleFunc("/adminRegistration", s.adminRegistration)
	sub.HandleFunc("/registerAccount", s.accountRegistration)
	sub.HandleFunc("/userInvitation", middleware.AuthenticationMiddleware(s.userInvitation))
	sub.HandleFunc("/updateAccount", middleware.AuthenticationMiddleware(s.updateAccount))
	sub.HandleFunc("/updatePassword", middleware.AuthenticationMiddleware(s.updatePassword))
	sub.HandleFunc("/login", s.login)
	sub.HandleFunc("/getAccountById", middleware.AuthenticationMiddleware(s.getAccountById))
	sub.HandleFunc("/getAccount", middleware.AuthenticationMiddleware(s.getAccount))

	sub.HandleFunc("/getCompaniesByAccountId", middleware.AuthenticationMiddleware(s.getCompaniesByAccountId))
	sub.HandleFunc("/deleteAccountById", deleteAccountById)
}

func (s *service) adminRegistration(w http.ResponseWriter, r *http.Request) {
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
	account.RoleID = 1

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

	token, expTime, hasError := handlers.GenerateJWT(account.Username, account.ID, nil, 1)
	if hasError != nil {
		handlers.ProduceErrorResponse(hasError.Error(), w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expTime,
	})

	account.Password = ""

	jsonRetrievedAccount, err := json.Marshal(account)
	if err != nil {
		fmt.Println(err)
		return
	}
	currentDate := time.Now().Format("2006-01-02 15:04:05")

	var response models.Response
	response.Date = currentDate

	response.Response = token
	response.Status = "success"
	response.Message = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)

	return
}

func (s *service) accountRegistration(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	var relation models.Relation

	verifToken := r.URL.Query().Get("token")
	if verifToken == "" {
		handlers.ProduceErrorResponse("Verification token invalid", w, r)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	retrievedSmartRegisteredLink, err := s.companyRepository.GetInvitationToken(verifToken)
	if err != nil {
		handlers.ProduceErrorResponse("Invalid Token", w, r)
		return
	}

	if retrievedSmartRegisteredLink.Email != account.Email {
		handlers.ProduceErrorResponse("Invalid Token.", w, r)
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
	account.RoleID = 2

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

	retrievedCompany, err := s.companyRepository.GetCompanyByID(int(retrievedSmartRegisteredLink.CompanyID))
	if err != nil {
		handlers.ProduceErrorResponse("Invalid Token", w, r)
		return
	}

	compArr := []uint{retrievedSmartRegisteredLink.CompanyID}

	relation.AccountID = account.ID
	relation.Companies = append(relation.Companies, retrievedCompany)
	relation.Title = "Employee"
	relation.Type = "user"
	relation.AddedByID = retrievedSmartRegisteredLink.AddedByID

	relation, err = s.companyRepository.AddRelation(relation)
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "users_company_email_key") {
			msg = "user already exists!"
		} else {
			msg = "Bad Request"
		}
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}

	// Delete invitation record after registration successful
	_, err = s.companyRepository.DeleteInvitation(retrievedSmartRegisteredLink.ID)
	if err != nil {
		handlers.ProduceErrorResponse("Error on deleting invitation", w, r)
		return
	}

	token, expTime, hasError := handlers.GenerateJWT(account.Username, account.ID, compArr, 2)
	if hasError != nil {
		handlers.ProduceErrorResponse(hasError.Error(), w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expTime,
	})

	account.Password = ""

	jsonRetrievedAccount, err := json.Marshal(account)
	if err != nil {
		fmt.Println(err)
		return
	}
	currentDate := time.Now().Format("2006-01-02 15:04:05")

	var response models.Response
	response.Date = currentDate

	response.Response = token
	response.Status = "success"
	response.Message = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)

	return
}

func (s *service) userInvitation(w http.ResponseWriter, r *http.Request) {
	var account models.SmartRegisterLink

	roleID := gcontext.Get(r, "roleID").(uint)

	if roleID != 1 {
		handlers.ProduceErrorResponse("You are not authorized to do this action.", w, r)
		return
	}

	id := gcontext.Get(r, "id").(uint)
	if roleID != 1 {
		handlers.ProduceErrorResponse("You are not authorized to do this action.", w, r)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	ownsCompany, errMsg := handlers.ValidateCompany(account.CompanyID, r)
	if !ownsCompany {
		handlers.ProduceErrorResponse(errMsg, w, r)
		return
	}

	newAccount := models.NewSmartRegisterLink()
	newAccount.CompanyID = account.CompanyID
	newAccount.Email = account.Email
	newAccount.AddedByID = id

	newAccount, err = s.companyRepository.AddInvitationUser(newAccount)
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "email_key") {
			msg = "User already has been invited!"
		} else {
			msg = "Bad Request"
		}
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}

	utils.NewredisRepository().Add_mail_to_Redis_queue(account.Email, os.Getenv("CLIENT_ORIGIN")+"/join/"+newAccount.Token, "Your account verification code", "Registration")

	handlers.ProduceSuccessResponse("User has been invited", "", w, r)
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
	handlers.ProduceSuccessResponse("Update of Password - Successful", "", w, r)
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
	handlers.ProduceSuccessResponse("Update of Account - Successful", "", w, r)
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

		token, expTime, hasError := handlers.GenerateJWT(account.Username, retrievedAccount.ID, retrievedCompanies, retrievedAccount.RoleID)
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

func (s *service) getAccount(w http.ResponseWriter, r *http.Request) {

	var retrievedAccount models.Account

	username := gcontext.Get(r, "username").(string)

	retrievedAccount, err := s.repository.GetAccountByUsername(username)
	if err != nil {
		handlers.ProduceErrorResponse("You are not authorized to access this data", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(retrievedAccount)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
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

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
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

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
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
