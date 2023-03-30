package patients

import (
	"encoding/json"
	"fmt"
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

type service struct {
	patientRepository repository.PatientRepository
}

func NewService() *service {

	return &service{
		patientRepository: repository.NewPatientService(),
	}
}

func (s *service) RegisterHandlers(route *mux.Router) {

	s.Handle(route)

}

func (s *service) Handle(route *mux.Router) {

	sub := route.PathPrefix("/patient").Subrouter()

	sub.HandleFunc("/registerPatient", middleware.AuthenticationMiddleware(s.patientRegistration))
	sub.HandleFunc("/updatePatient", middleware.AuthenticationMiddleware(s.updatePatient))
	sub.HandleFunc("/getPatient", middleware.AuthenticationMiddleware(s.getPatientData))
	sub.HandleFunc("/searchPatient", middleware.AuthenticationMiddleware(s.getPatientDataKeyword))
	sub.HandleFunc("/getAllPatients", middleware.AuthenticationMiddleware(s.getAllPatients))
}

func (s *service) getAllPatients(w http.ResponseWriter, r *http.Request) {
	var patients []models.Patient
	var response models.Response

	username := gcontext.Get(r, "username").(string)

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	patients, err := s.patientRepository.GetAllPatients()
	if err != nil {
		response.Status = "error"
		response.Message = "Unknown Username or Password"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(patients)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = username
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)

}

func (s *service) getPatientDataKeyword(w http.ResponseWriter, r *http.Request) {
	var patients []models.Patient
	var response models.Response
	var err error

	username := gcontext.Get(r, "username").(string)

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		response.Status = "error"
		response.Message = "Please input all required fields."
		json.NewEncoder(w).Encode(response)

		return
	}

	numberKeyword := keyword
	if numberKeyword, errConvert := strconv.Atoi(numberKeyword); errConvert == nil {
		patients, err = s.patientRepository.GetPatientAmka(numberKeyword)
	} else {
		keyword = "%" + keyword + "%"
		patients, err = s.patientRepository.GetPatientKeyword(keyword)
	}
	if err != nil {
		response.Status = "error"
		response.Message = err.Error()
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// if account.Username != username {
	// 	http.Error(w, "You are not authorized to view this data.", http.StatusBadRequest)
	// 	return
	// }

	jsonRetrievedAccount, err := json.Marshal(patients)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = username
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)

}

func (s *service) getPatientData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var patient models.Patient
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

	patient, err = s.patientRepository.GetPatient(intID)
	if err != nil {
		response.Status = "error"
		response.Message = "Unknown Username or Password"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// if account.Username != username {
	// 	http.Error(w, "You are not authorized to view this data.", http.StatusBadRequest)
	// 	return
	// }

	jsonRetrievedAccount, err := json.Marshal(patient)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = username
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)

}

func (s *service) updatePatient(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient
	var response models.Response

	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(patient)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	patient, err = s.patientRepository.UpdatePatient(patient)
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

func (s *service) patientRegistration(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient
	var response models.Response

	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: check company ID if exists and if caller is related

	isValid, errors := handlers.ValidateInputs(patient)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	patient, err = s.patientRepository.AddPatient(patient)
	if err != nil {
		var newerr string
		if strings.Contains(err.Error(), "idx_patients_amka") {
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

	response.Status = "success"
	response.Message = "Registration of Account - Successful"
	response.Response = ""
	json.NewEncoder(w).Encode(response)
	return
}
