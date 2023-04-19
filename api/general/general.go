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
}

func (s *service) getAllDrugs(w http.ResponseWriter, r *http.Request) {
	var drugs []models.Drug

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	drugs, err := s.generalRepository.GetAllDrugs()
	if err != nil {
		s.response.Status = "error"
		s.response.Message = "Something went wrong"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(drugs)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.response.Status = "success"
	s.response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(s.response)
}

func (s *service) getDrug(w http.ResponseWriter, r *http.Request) {
	var drug models.Drug

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		s.response.Status = "error"
		s.response.Message = "Please input all required fields."
		json.NewEncoder(w).Encode(s.response)

		return
	}
	intID, err := strconv.Atoi(id)

	drug, err = s.generalRepository.GetDrug(intID)
	if err != nil {
		s.response.Status = "error"
		s.response.Message = "Record not found"
		s.response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(drug)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.response.Status = "success"
	s.response.Message = ""
	s.response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(s.response)

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
		s.response.Status = "error"
		s.response.Message = "Unknown Username or Password"
		s.response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.response)
		return
	}

	s.response.Status = "success"
	s.response.Message = "Deleted"
	s.response.Response = ""
	json.NewEncoder(w).Encode(s.response)

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
		s.response.Status = "error"
		s.response.Message = newerr
		s.response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.response)
		return
	}

	fmt.Fprintf(w, "Drug Udated - Successful")
}

func (s *service) addDrug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var drug models.Drug
	var response models.Response
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
		response.Status = "error"
		response.Message = newerr
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Fprintf(w, "Registration of Drug - Successful")
}

////////////////////// ############## ALLERGIES ################# /////////////////

func (s *service) getAllAllergies(w http.ResponseWriter, r *http.Request) {
	var allergies []models.Allergy

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	allergies, err := s.generalRepository.GetAllAllergies()
	if err != nil {
		s.response.Status = "error"
		s.response.Message = "Something went wrong"
		s.response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(allergies)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.response.Status = "success"
	s.response.Message = ""
	s.response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(s.response)

}

func (s *service) getAllergy(w http.ResponseWriter, r *http.Request) {
	var allergy models.Allergy

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	s.response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		s.response.Status = "error"
		s.response.Message = "Please input all required fields."
		json.NewEncoder(w).Encode(s.response)

		return
	}
	intID, err := strconv.Atoi(id)

	allergy, err = s.generalRepository.GetAllergy(intID)
	if err != nil {
		s.response.Status = "error"
		s.response.Message = "Record not found"
		s.response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(allergy)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.response.Status = "success"
	s.response.Message = ""
	s.response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(s.response)

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
		s.response.Status = "error"
		s.response.Message = "Something went wrong"
		s.response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.response)
		return
	}

	s.response.Status = "success"
	s.response.Message = "Deleted"
	s.response.Response = ""
	json.NewEncoder(w).Encode(s.response)

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
		s.response.Status = "error"
		s.response.Message = newerr
		s.response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.response)
		return
	}

	fmt.Fprintf(w, "Allergy Udated - Successful")
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
		s.response.Status = "error"
		s.response.Message = newerr
		s.response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.response)
		return
	}

	fmt.Fprintf(w, "Registration of Allergy - Successful")
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

	category, err = s.generalRepository.AddClinicalTestCategory(category)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "users_company_email_key") {
			newerr = "Category already registered!"
		} else {
			newerr = "Bad Request"
		}
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Registration of Category - Successful", w, r)
}
