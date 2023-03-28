package companies

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rehab-backend/internal/middleware"
	"github.com/rehab-backend/internal/pkg/handlers"
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/rehab-backend/internal/repository"
)

const ADMIN = "admin"
const USER = "user"
const EMPLOYEE = "employee"
const EMPLOYER = "employer"

type service struct {
	repository repository.CompanyRepository
}

func NewService() *service {

	return &service{
		repository: repository.NewCompanyService(),
	}
}

func (s *service) RegisterHandlers(route *mux.Router) {

	s.Handle(route)

}

func (s *service) Handle(route *mux.Router) {

	sub := route.PathPrefix("/companies").Subrouter()

	sub.HandleFunc("/registerCompany", middleware.AuthenticationMiddleware(s.companyRegistration))
	sub.HandleFunc("/addRelation", middleware.AuthenticationMiddleware(s.addRelation))
	sub.HandleFunc("/getRelation", middleware.AuthenticationMiddleware(s.getRelation))
	// sub.HandleFunc("/getRelationByCompany", middleware.AuthenticationMiddleware(s.getRelationByCompany))

	sub.HandleFunc("/updateCompany", middleware.AuthenticationMiddleware(s.upateCompany))
	sub.HandleFunc("/getCompany", middleware.AuthenticationMiddleware(s.getCompanyData))
	// sub.HandleFunc("/addRelation", middleware.AuthenticationMiddleware(s.getAllPatients))

}

func (s *service) companyRegistration(w http.ResponseWriter, r *http.Request) {
	var company models.Company
	var response models.Response
	var relation models.Relation

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Unmarshal(data, &company)
	json.Unmarshal(data, &relation)

	isValid, errors := handlers.ValidateInputs(company)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	company, err = s.repository.RegisterCompany(company)
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

	id := gcontext.Get(r, "id").(uint)

	relation.AccountID = id
	relation.Companies = append(relation.Companies, company)
	relation.Title = "CEO"
	relation.Type = "admin"

	relation, err = s.repository.AddRelation(relation)
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

	fmt.Fprintf(w, "Registration of Company - Successful")
}

func (s *service) addRelation(w http.ResponseWriter, r *http.Request) {
	var relation models.Relation
	var response models.Response

	err := json.NewDecoder(r.Body).Decode(&relation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Validate user that is admin and can add relation

	// isValid, errors := handlers.ValidateInputs(company)
	// if !isValid {
	// 	for _, fieldError := range errors {
	// 		http.Error(w, fieldError, http.StatusBadRequest)
	// 		return
	// 	}
	// }

	relation, err = s.repository.AddRelation(relation)
	if err != nil {
		var newerr string
		response.Status = "error"
		response.Message = newerr
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Fprintf(w, "Registration of Company - Successful")
}

func (s *service) getCompanyData(w http.ResponseWriter, r *http.Request) {
	var company []models.Relation
	var response models.Response

	username := gcontext.Get(r, "username").(string)

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		response.Status = "error"
		response.Message = "Please input all required fields."
		json.NewEncoder(w).Encode(response)

		return
	}
	intID, err := strconv.Atoi(id)

	company, err = s.repository.GetCompanyByID(intID)
	if err != nil {
		response.Status = "error"
		response.Message = "Unknown Username or Password"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(company)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = username
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)
}

// func (s *service) getRelationByCompany(w http.ResponseWriter, r *http.Request) {
// 	var relation []models.Relation
// 	var response models.Response

// 	id := 2

// 	currentDate := time.Now().Format("2006-01-02 15:04:05")
// 	response.Date = currentDate

// 	relation, err := s.repository.GetRelationsByCompanyID(id)
// 	if err != nil {
// 		response.Status = "error"
// 		response.Message = "Unknown Username or Password"
// 		response.Response = ""
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	jsonRetrievedAccount, err := json.Marshal(relation)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	response.Status = "success"
// 	response.Message = ""
// 	response.Response = string(jsonRetrievedAccount)
// 	json.NewEncoder(w).Encode(response)
// }

func (s *service) getRelation(w http.ResponseWriter, r *http.Request) {
	var relation []models.Relation
	var response models.Response

	id := gcontext.Get(r, "id").(uint)

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	relation, err := s.repository.GetRelationsByAccountID(id)
	if err != nil {
		response.Status = "error"
		response.Message = "Unknown Username or Password"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(relation)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = ""
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)
}

func (s *service) upateCompany(w http.ResponseWriter, r *http.Request) {
	var company models.Company
	var response models.Response

	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// isValid, errors := handlers.ValidateInputs(company)
	// if !isValid {
	// 	for _, fieldError := range errors {
	// 		http.Error(w, fieldError, http.StatusBadRequest)
	// 		return
	// 	}
	// }

	company, err = s.repository.UpdateCompany(company)
	if err != nil {
		var newerr string
		// if strings.Contains(err.Error(), "users_company_email_key") {
		// 	newerr = "user already exists!"
		// } else {
		// 	newerr = "Bad Request"
		// }
		response.Status = "error"
		response.Message = newerr
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Fprintf(w, "Company Update - Successful")
}
