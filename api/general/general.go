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
}

func (s *service) getAllDrugs(w http.ResponseWriter, r *http.Request) {
	var drugs []models.Drug
	var response models.Response

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	drugs, err := s.generalRepository.GetAllDrugs()
	if err != nil {
		response.Status = "error"
		response.Message = "Something went wrong"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(drugs)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = ""
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)

}

func (s *service) getDrug(w http.ResponseWriter, r *http.Request) {
	var drug models.Drug
	var response models.Response

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

	drug, err = s.generalRepository.GetDrug(intID)
	if err != nil {
		response.Status = "error"
		response.Message = "Record not found"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(drug)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = ""
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)

}

func (s *service) deleteDrug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var response models.Response
	var drug models.Drug

	err := json.NewDecoder(r.Body).Decode(&drug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	_, err = s.generalRepository.DeleteDrug(drug.ID)
	if err != nil {
		response.Status = "error"
		response.Message = "Unknown Username or Password"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = "success"
	response.Message = "Deleted"
	response.Response = ""
	json.NewEncoder(w).Encode(response)

}

func (s *service) updateDrug(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(drug.DrugTitle)
	drug, err = s.generalRepository.UpdateDrug(drug)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "record not found") {
			newerr = "Record not Found!"
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
	var response models.Response

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	allergies, err := s.generalRepository.GetAllAllergies()
	if err != nil {
		response.Status = "error"
		response.Message = "Something went wrong"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(allergies)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = ""
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)

}

func (s *service) getAllergy(w http.ResponseWriter, r *http.Request) {
	var allergy models.Allergy
	var response models.Response

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

	allergy, err = s.generalRepository.GetAllergy(intID)
	if err != nil {
		response.Status = "error"
		response.Message = "Record not found"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(allergy)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = ""
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)

}

func (s *service) deleteAllergy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var response models.Response
	var allergy models.Allergy

	err := json.NewDecoder(r.Body).Decode(&allergy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	_, err = s.generalRepository.DeleteAllergy(allergy.ID)
	if err != nil {
		response.Status = "error"
		response.Message = "Something went wrong"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = "success"
	response.Message = "Deleted"
	response.Response = ""
	json.NewEncoder(w).Encode(response)

}

func (s *service) updateAllergy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var allergy models.Allergy
	var response models.Response
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
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
		response.Status = "error"
		response.Message = newerr
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Fprintf(w, "Allergy Udated - Successful")
}

func (s *service) addAllergy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var allergy models.Allergy
	var response models.Response
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
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
		response.Status = "error"
		response.Message = newerr
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Fprintf(w, "Registration of Allergy - Successful")
}
