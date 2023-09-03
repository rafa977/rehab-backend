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
	sub.HandleFunc("/getRelationIDsByAccountID", middleware.AuthenticationMiddleware(s.getRelationIDsByAccountID))

	sub.HandleFunc("/updateCompany", middleware.AuthenticationMiddleware(s.upateCompany))
	sub.HandleFunc("/getCompany", middleware.AuthenticationMiddleware(s.getCompanyData))
	sub.HandleFunc("/getCompaniesDetails", middleware.AuthenticationMiddleware(s.getCompaniesDetailsDataByAccountID))

	// sub.HandleFunc("/addRelation", middleware.AuthenticationMiddleware(s.getAllPatients))

}

func (s *service) companyRegistration(w http.ResponseWriter, r *http.Request) {
	var company models.Company
	var relation models.Relation

	roleID := gcontext.Get(r, "roleID").(uint)
	if roleID != 1 {
		handlers.ProduceErrorResponse("You are not authorized to do this action.", w, r)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	json.Unmarshal(data, &company)
	json.Unmarshal(data, &relation)

	isValid, errors := handlers.ValidateInputs(company)
	if !isValid {
		for _, fieldError := range errors {
			handlers.ProduceErrorResponse(fieldError, w, r)
			return
		}
	}

	company, err = s.repository.RegisterCompany(company)
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "companies_tax_id") {
			msg = "Company already registered!"
		} else {
			msg = "Bad Request"
		}
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}

	id := gcontext.Get(r, "id").(uint)
	username := gcontext.Get(r, "username").(string)

	relation.AccountID = id
	relation.Companies = append(relation.Companies, company)
	relation.Title = "CEO"
	relation.Type = "admin"
	relation.AddedByID = id

	relation, err = s.repository.AddRelation(relation)
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

	retrievedCompanies := s.repository.GetCompaniesByAccountID(id)

	token, _, hasError := handlers.GenerateJWT(username, id, retrievedCompanies, 1)
	if hasError != nil {
		http.Error(w, hasError.Error(), http.StatusBadRequest)
		return
	}

	handlers.ProduceSuccessResponse(token, w, r)
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
	var company models.Company

	intID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	isOwner, errMsg := handlers.ValidateCompany(uint(intID), r)
	if !isOwner {
		handlers.ProduceErrorResponse(errMsg, w, r)
		return
	}

	company, err = s.repository.GetCompanyByID(intID)
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}

	jsonRetrievedCompany, err := json.Marshal(company)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedCompany), w, r)
	return
}

func (s *service) getCompaniesDetailsDataByAccountID(w http.ResponseWriter, r *http.Request) {

	id := gcontext.Get(r, "id").(uint)

	retrievedCompanies := s.repository.GetCompaniesDetailsByAccountID(id)

	jsonRetrieved, err := json.Marshal(retrievedCompanies)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrieved), w, r)
}

func (s *service) getRelationIDsByAccountID(w http.ResponseWriter, r *http.Request) {
	var relationIDs []models.Relation
	var response models.Response

	id := r.URL.Query().Get("id")
	if id == "" {
		response.Status = "error"
		response.Message = "Please input all required fields."
		json.NewEncoder(w).Encode(response)

		return
	}
	intID, err := strconv.Atoi(id)

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	relationIDs, err = s.repository.GetRelationIDsByAccountID(intID)
	if err != nil {
		response.Status = "error"
		response.Message = "Unknown Username or Password"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println(relationIDs)
	jsonRetrievedAccount, err := json.Marshal(relationIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = ""
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)
}

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
