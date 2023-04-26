package general

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rehab-backend/internal/middleware"
	"github.com/rehab-backend/internal/pkg/handlers"
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/rehab-backend/internal/repository"
)

type service struct {
	generalRepository repository.GeneralRepository
	response          models.Response
}

func NewService() *service {
	return &service{
		generalRepository: repository.NewGeneralRepoService(),
	}
}

func (s *service) RegisterHandlers(route *mux.Router) {
	s.Handle(route)
}

func (s *service) Handle(route *mux.Router) {

	sub := route.PathPrefix("/general").Subrouter()

	sub.HandleFunc("/addDrug", middleware.AuthenticationMiddleware(s.addDrug))
	sub.HandleFunc("/updateDrug", middleware.AuthenticationMiddleware(s.updateDrug))
	sub.HandleFunc("/getDrug", middleware.AuthenticationMiddleware(s.getDrug))
	sub.HandleFunc("/deleteDrug", middleware.AuthenticationMiddleware(s.deleteDrug))
	sub.HandleFunc("/getAllDrugs", middleware.AuthenticationMiddleware(s.getAllDrugs))

	sub.HandleFunc("/addAllergy", middleware.AuthenticationMiddleware(s.addAllergy))
	sub.HandleFunc("/updateAllergy", middleware.AuthenticationMiddleware(s.updateAllergy))
	sub.HandleFunc("/getAllergy", middleware.AuthenticationMiddleware(s.getAllergy))
	sub.HandleFunc("/deleteAllergy", middleware.AuthenticationMiddleware(s.deleteAllergy))
	sub.HandleFunc("/getAllAllergies", middleware.AuthenticationMiddleware(s.getAllAllergies))

	sub.HandleFunc("/addDisorder", middleware.AuthenticationMiddleware(s.addDisorder))
	sub.HandleFunc("/updateDisorder", middleware.AuthenticationMiddleware(s.updateDisorder))
	sub.HandleFunc("/getDisorder", middleware.AuthenticationMiddleware(s.getDisorder))
	sub.HandleFunc("/deleteDisorder", middleware.AuthenticationMiddleware(s.deleteDisorder))
	sub.HandleFunc("/getAllDisorders", middleware.AuthenticationMiddleware(s.getAllDisorders))

	sub.HandleFunc("/addClinicalTestCategory", middleware.AuthenticationMiddleware(s.addClinicalTestCategory))
	sub.HandleFunc("/updateClinicalTestCategory", middleware.AuthenticationMiddleware(s.updateClinicalTestCategory))
	sub.HandleFunc("/getClinicalTestCategory", middleware.AuthenticationMiddleware(s.getClinicalTestCategory))
	sub.HandleFunc("/deleteClinicalTestCategory", middleware.AuthenticationMiddleware(s.deleteClinicalTestCategory))
	sub.HandleFunc("/getAllClinicalTestCategories", middleware.AuthenticationMiddleware(s.getAllClinicalTestCategories))

	sub.HandleFunc("/addClinicalTest", middleware.AuthenticationMiddleware(s.addClinicalTest))
	sub.HandleFunc("/updateClinicalTest", middleware.AuthenticationMiddleware(s.updateClinicalTest))
	sub.HandleFunc("/getClinicalTest", middleware.AuthenticationMiddleware(s.getClinicalTest))
	sub.HandleFunc("/deleteClinicalTest", middleware.AuthenticationMiddleware(s.deleteClinicalTest))
	sub.HandleFunc("/getAllClinicalTests", middleware.AuthenticationMiddleware(s.getAllClinicalTests))
}

////////////////////// ############## DRUGS ################# /////////////////

func (s *service) getAllDrugs(w http.ResponseWriter, r *http.Request) {
	var drugs []models.Drug

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	drugs, err := s.generalRepository.GetAllDrugs()
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(drugs)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) getDrug(w http.ResponseWriter, r *http.Request) {
	var drug models.Drug

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}
	intID, err := strconv.Atoi(id)

	drug, err = s.generalRepository.GetDrug(intID)
	if err != nil {
		handlers.ProduceErrorResponse("Record not found", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(drug)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) deleteDrug(w http.ResponseWriter, r *http.Request) {
	var drug models.Drug

	err := json.NewDecoder(r.Body).Decode(&drug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	_, err = s.generalRepository.DeleteDrug(drug.ID)
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}

	handlers.ProduceSuccessResponse("Drug Deleted - Successful", w, r)
}

func (s *service) updateDrug(w http.ResponseWriter, r *http.Request) {
	var drug models.Drug

	err := json.NewDecoder(r.Body).Decode(&drug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(drug)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}
	fmt.Println(drug.DrugTitle)
	drug, err = s.generalRepository.UpdateDrug(drug)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "record not found") {
			newerr = "Record not Found!"
		} else {
			newerr = "Bad Request"
		}
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Drug Udated - Successful", w, r)
}

func (s *service) addDrug(w http.ResponseWriter, r *http.Request) {
	var drug models.Drug
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&drug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(drug)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	drug, err = s.generalRepository.AddDrug(drug)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "users_company_email_key") {
			newerr = "Drug already registered!"
		} else {
			newerr = "Bad Request"
		}
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Registration of Drug - Successful", w, r)
}

////////////////////// ############## ALLERGIES ################# /////////////////

func (s *service) getAllAllergies(w http.ResponseWriter, r *http.Request) {
	var allergies []models.Allergy

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	allergies, err := s.generalRepository.GetAllAllergies()
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(allergies)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) getAllergy(w http.ResponseWriter, r *http.Request) {
	var allergy models.Allergy

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}
	intID, err := strconv.Atoi(id)

	allergy, err = s.generalRepository.GetAllergy(intID)
	if err != nil {
		handlers.ProduceErrorResponse("Record not found", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(allergy)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) deleteAllergy(w http.ResponseWriter, r *http.Request) {
	var allergy models.Allergy

	err := json.NewDecoder(r.Body).Decode(&allergy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	_, err = s.generalRepository.DeleteAllergy(allergy.ID)
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}

	handlers.ProduceSuccessResponse("Allergy Deleted - Successful", w, r)
}

func (s *service) updateAllergy(w http.ResponseWriter, r *http.Request) {
	var allergy models.Allergy

	err := json.NewDecoder(r.Body).Decode(&allergy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(allergy)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}
	allergy, err = s.generalRepository.UpdateAllergy(allergy)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "record not found") {
			newerr = "Record not Found!"
		} else {
			newerr = "Bad Request"
		}
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Allergy Udated - Successful", w, r)
}

func (s *service) addAllergy(w http.ResponseWriter, r *http.Request) {
	var allergy models.Allergy

	err := json.NewDecoder(r.Body).Decode(&allergy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(allergy)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	allergy, err = s.generalRepository.AddAllergy(allergy)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "users_company_email_key") {
			newerr = "Allergy already registered!"
		} else {
			newerr = "Bad Request"
		}
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Registration of Allergy - Successful", w, r)
}

////////////////////// ############## DISORDERS ################# /////////////////
func (s *service) getAllDisorders(w http.ResponseWriter, r *http.Request) {
	var disorders []models.Disorder

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	disorders, err := s.generalRepository.GetAllDisorders()
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(disorders)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) getDisorder(w http.ResponseWriter, r *http.Request) {
	var disorder models.Disorder

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}
	intID, err := strconv.Atoi(id)

	disorder, err = s.generalRepository.GetDisorder(intID)
	if err != nil {
		handlers.ProduceErrorResponse("Record not found", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(disorder)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) deleteDisorder(w http.ResponseWriter, r *http.Request) {
	var disorder models.Disorder

	err := json.NewDecoder(r.Body).Decode(&disorder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	_, err = s.generalRepository.DeleteDisorder(disorder.ID)
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}

	handlers.ProduceSuccessResponse("Record Deleted", w, r)
}

func (s *service) updateDisorder(w http.ResponseWriter, r *http.Request) {
	var disorder models.Disorder

	err := json.NewDecoder(r.Body).Decode(&disorder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(disorder)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}
	disorder, err = s.generalRepository.UpdateDisorder(disorder)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "record not found") {
			newerr = "Record not Found!"
		} else {
			newerr = "Bad Request"
		}
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Disorder Udated - Successful", w, r)
}

func (s *service) addDisorder(w http.ResponseWriter, r *http.Request) {
	var disorder models.Disorder

	err := json.NewDecoder(r.Body).Decode(&disorder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(disorder)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	disorder, err = s.generalRepository.AddDisorder(disorder)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "users_company_email_key") {
			newerr = "Allergy already registered!"
		} else {
			newerr = "Bad Request"
		}
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Registration of Disorder - Successful", w, r)
}

////////////////////// ############## CLINICAL TEST CATEGORIES ################# /////////////////
func (s *service) getAllClinicalTestCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.ClinicalTestCategory

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	categories, err := s.generalRepository.GetAllClinicalTestCategories()
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(categories)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) getClinicalTestCategory(w http.ResponseWriter, r *http.Request) {
	var category models.ClinicalTestCategory

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}
	intID, err := strconv.Atoi(id)

	category, err = s.generalRepository.GetClinicalTestCategory(intID)
	if err != nil {
		handlers.ProduceErrorResponse("Record not found", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(category)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) deleteClinicalTestCategory(w http.ResponseWriter, r *http.Request) {
	var category models.ClinicalTestCategory

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	_, err = s.generalRepository.DeleteClinicalTestCategory(category.ID)
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}
	handlers.ProduceSuccessResponse("Record Deleted", w, r)
}

func (s *service) updateClinicalTestCategory(w http.ResponseWriter, r *http.Request) {
	var category models.ClinicalTestCategory

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(category)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}
	category, err = s.generalRepository.UpdateClinicalTestCategory(category)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "record not found") {
			newerr = "Record not Found!"
		} else {
			newerr = "Bad Request"
		}
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Category Udated - Successful", w, r)
}

func (s *service) addClinicalTestCategory(w http.ResponseWriter, r *http.Request) {
	var category models.ClinicalTestCategory

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(category)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	// TODO: check company ID if exists and if caller is related
	isOwner, ownerError := handlers.ValidateCompany(category.CompanyID, r)
	if !isOwner {
		handlers.ProduceErrorResponse(ownerError, w, r)
		return
	}

	category, err = s.generalRepository.AddClinicalTestCategory(category)
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "duplicate key value violates") {
			msg = "Clinical test category already registered!"
		} else {
			msg = "Bad Request"
		}
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Registration of Category - Successful", w, r)
}

////////////////////// ############## CLINICAL TESTS ################# /////////////////
func (s *service) getAllClinicalTests(w http.ResponseWriter, r *http.Request) {
	var clinical []models.ClinicalTests

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	clinical, err := s.generalRepository.GetAllClinicalTests()
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(clinical)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) getClinicalTest(w http.ResponseWriter, r *http.Request) {
	var clinical models.ClinicalTests

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}
	intID, err := strconv.Atoi(id)

	clinical, err = s.generalRepository.GetClinicalTest(intID)
	if err != nil {
		handlers.ProduceErrorResponse("Record not found", w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(clinical)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) deleteClinicalTest(w http.ResponseWriter, r *http.Request) {
	var clinical models.ClinicalTests

	err := json.NewDecoder(r.Body).Decode(&clinical)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	// TODO: check company ID if exists and if caller is related
	isOwner, ownerError := handlers.ValidateCompany(clinical.CompanyID, r)
	if !isOwner {
		handlers.ProduceErrorResponse(ownerError, w, r)
		return
	}

	_, err = s.generalRepository.DeleteClinicalTest(clinical.ID)
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}
	handlers.ProduceSuccessResponse("Record Deleted", w, r)
}

func (s *service) updateClinicalTest(w http.ResponseWriter, r *http.Request) {
	var clinical models.ClinicalTests

	err := json.NewDecoder(r.Body).Decode(&clinical)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(clinical)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	// TODO: check company ID if exists and if caller is related
	isOwner, ownerError := handlers.ValidateCompany(clinical.CompanyID, r)
	if !isOwner {
		handlers.ProduceErrorResponse(ownerError, w, r)
		return
	}

	clinical, err = s.generalRepository.UpdateClinicalTest(clinical)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "record not found") {
			newerr = "Record not Found!"
		} else {
			newerr = "Bad Request"
		}
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Category Udated - Successful", w, r)
}

func (s *service) addClinicalTest(w http.ResponseWriter, r *http.Request) {
	var clinical models.ClinicalTests

	err := json.NewDecoder(r.Body).Decode(&clinical)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(clinical)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	// TODO: check company ID if exists and if caller is related
	isOwner, ownerError := handlers.ValidateCompany(clinical.CompanyID, r)
	if !isOwner {
		handlers.ProduceErrorResponse(ownerError, w, r)
		return
	}

	category, err := s.generalRepository.GetClinicalTestCategory(int(clinical.ClinicalTestCategoryID))
	if err != nil {
		handlers.ProduceErrorResponse("Record not found", w, r)
		return
	}

	if category.CompanyID != clinical.CompanyID {
		handlers.ProduceErrorResponse("Company does not match", w, r)
		return
	}

	clinical, err = s.generalRepository.AddClinicalTest(clinical)
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "title") {
			msg = "Clinical test already registered!"
		} else if strings.Contains(err.Error(), "fk_") {
			msg = "Clinical test category does not exist!"
		} else {
			msg = "Bad Request"
		}
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Registration of Category - Successful", w, r)
}
