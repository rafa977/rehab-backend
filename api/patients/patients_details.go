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

type detailsService struct {
	patientDetailsRepository repository.PatientDetailRepository
	patientRepository        repository.PatientRepository
}

func NewDetailsService() *detailsService {

	return &detailsService{
		patientDetailsRepository: repository.NewPatientDetailsService(),
		patientRepository:        repository.NewPatientService(),
	}
}

func (s *detailsService) RegisterDetailHandlers(route *mux.Router) {
	s.DetailHandle(route)
}

func (s *detailsService) DetailHandle(route *mux.Router) {

	sub := route.PathPrefix("/patientDetails").Subrouter()

	sub.HandleFunc("/patientDetailsRegistration", middleware.AuthenticationMiddleware(s.patientDetailsRegistration))
	sub.HandleFunc("/updatePatientDetails", middleware.AuthenticationMiddleware(s.updatePatientDetails))
	sub.HandleFunc("/getPatientDetails", middleware.AuthenticationMiddleware(s.getPatientDetailsFull))
	sub.HandleFunc("/getPatientsDetailsByCompanyID", middleware.AuthenticationMiddleware(s.getPatientsDetailsByCompanyID))
}

func (s *detailsService) patientDetailsRegistration(w http.ResponseWriter, r *http.Request) {
	var patientDetails models.PatientDetails

	err := json.NewDecoder(r.Body).Decode(&patientDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: check company ID if exists and if caller is related
	compIDs := gcontext.Get(r, "compIDs").([]uint)
	userID := gcontext.Get(r, "id").(uint)
	if len(compIDs) == 0 {
		handlers.ProduceErrorResponse("Please register your company", w, r)
		return
	}

	// check patient id exists and is under same company
	isPatientValid, validationError := s.patientRepository.CheckPatient(patientDetails.PatientID, compIDs)
	if !isPatientValid {
		handlers.ProduceErrorResponse(validationError, w, r)
		return
	}

	patientDetails.CreatedBy = userID

	patientDetails, err = s.patientDetailsRepository.AddPatientDetails(patientDetails)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	handlers.ProduceSuccessResponse("Registration of Details - Successful", w, r)
}

func (s *detailsService) getPatientsDetailsByCompanyID(w http.ResponseWriter, r *http.Request) {
	var patients []models.PatientDetails

	patients, err := s.patientDetailsRepository.GetPatientDetailsByCompanyID(1)
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

func (s *detailsService) getPatientDetailsFull(w http.ResponseWriter, r *http.Request) {
	var patientDetails models.PatientDetails
	var patient models.Patient
	var response models.Response

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	intID, err := strconv.Atoi(id)

	patientDetails, err = s.patientDetailsRepository.GetPatientDetailsFull(intID)
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "record not found") {
			msg = "You are not authorized to access these data!"
		} else {
			msg = "Bad Request"
		}
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}

	patient, err = s.patientRepository.GetPatient(int(patientDetails.PatientID))
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	ownsCompany, errMsg := handlers.ValidateCompany(patient.CompanyID, r)
	if !ownsCompany {
		handlers.ProduceErrorResponse(errMsg, w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(patientDetails)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

// func (s *detailsService) getPatientFull(w http.ResponseWriter, r *http.Request) {
// 	var patient models.Patient
// 	var response models.Response

// 	username := gcontext.Get(r, "username").(string)

// 	currentDate := time.Now().Format("2006-01-02 15:04:05")
// 	response.Date = currentDate

// 	id := r.URL.Query().Get("id")
// 	if id == "" {
// 		response.Status = "error"
// 		response.Message = "Please input all required fields."
// 		json.NewEncoder(w).Encode(response)

// 		return
// 	}
// 	intID, err := strconv.Atoi(id)

// 	patient, err = s.patientRepository.GetPatientFull(intID)
// 	if err != nil {
// 		response.Status = "error"
// 		response.Message = "Unknown Username or Password"
// 		response.Response = ""
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	// if account.Username != username {
// 	// 	http.Error(w, "You are not authorized to view this data.", http.StatusBadRequest)
// 	// 	return
// 	// }

// 	jsonRetrievedAccount, err := json.Marshal(patient)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	response.Status = "success"
// 	response.Message = username
// 	response.Response = string(jsonRetrievedAccount)
// 	json.NewEncoder(w).Encode(response)

// }

// func (s *detailsService) getPatientData(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	var patient models.Patient
// 	var response models.Response

// 	username := gcontext.Get(r, "username").(string)

// 	currentDate := time.Now().Format("2006-01-02 15:04:05")
// 	response.Date = currentDate

// 	id := r.URL.Query().Get("id")
// 	if id == "" {
// 		response.Status = "error"
// 		response.Message = "Please input all required fields."
// 		json.NewEncoder(w).Encode(response)

// 		return
// 	}
// 	intID, err := strconv.Atoi(id)

// 	patient, err = s.patientRepository.GetPatient(intID)
// 	if err != nil {
// 		response.Status = "error"
// 		response.Message = "Unknown Username or Password"
// 		response.Response = ""
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	// if account.Username != username {
// 	// 	http.Error(w, "You are not authorized to view this data.", http.StatusBadRequest)
// 	// 	return
// 	// }

// 	jsonRetrievedAccount, err := json.Marshal(patient)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	response.Status = "success"
// 	response.Message = username
// 	response.Response = string(jsonRetrievedAccount)
// 	json.NewEncoder(w).Encode(response)

// }

func (s *detailsService) updatePatientDetails(w http.ResponseWriter, r *http.Request) {
	var patient models.PatientDetails
	var response models.Response

	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// isValid, errors := handlers.ValidateInputs(patient)
	// if !isValid {
	// 	for _, fieldError := range errors {
	// 		http.Error(w, fieldError, http.StatusBadRequest)
	// 		return
	// 	}
	// }

	patient, err = s.patientDetailsRepository.UpdatePatientDetails(patient)
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
