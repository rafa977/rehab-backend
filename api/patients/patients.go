package patients

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

func (s *service) patientRegistration(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient

	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: check company ID if exists and if caller is related
	isOwner, ownerError := handlers.ValidateCompany(patient.CompanyID, r)
	if !isOwner {
		handlers.ProduceErrorResponse(ownerError, w, r)
		return
	}

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
		if strings.Contains(err.Error(), "patients_amka_key") {
			newerr = "Patient with same AMKA already exists!"
		} else {
			newerr = "Bad Request"
		}
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}

	handlers.ProduceSuccessResponse("Registration of Account - Successful", w, r)
}

func (s *service) getAllPatients(w http.ResponseWriter, r *http.Request) {
	var patients []models.Patient
	var response models.Response

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	patients, err := s.patientRepository.GetAllPatients()
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(patients)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) getPatientDataKeyword(w http.ResponseWriter, r *http.Request) {
	var patients []models.Patient
	var err error

	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
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
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(patients)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) getPatientData(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}
	intID, err := strconv.Atoi(id)

	patient, err = s.patientRepository.GetPatient(intID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(patient)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *service) updatePatient(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient

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
		handlers.ProduceErrorResponse(newerr, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Update of Account - Successful", w, r)
}
