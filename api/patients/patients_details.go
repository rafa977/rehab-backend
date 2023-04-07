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
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/rehab-backend/internal/repository"
)

type detailsService struct {
	patientDetailsRepository repository.PatientDetailRepository
}

func NewDetailsService() *detailsService {

	return &detailsService{
		patientDetailsRepository: repository.NewPatientDetailsService(),
	}
}

func (s *detailsService) RegisterDetailHandlers(route *mux.Router) {

	s.DetailHandle(route)

}

func (s *detailsService) DetailHandle(route *mux.Router) {

	sub := route.PathPrefix("/patientDetails").Subrouter()

	sub.HandleFunc("/patientDetailsRegistration", middleware.AuthenticationMiddleware(s.patientDetailsRegistration))
	sub.HandleFunc("/updatePatientDetails", middleware.AuthenticationMiddleware(s.updatePatientDetails))
	sub.HandleFunc("/getPatientDetailsFull", middleware.AuthenticationMiddleware(s.getPatientDetailsFull))
	sub.HandleFunc("/getPatientsDetailsByCompanyID", middleware.AuthenticationMiddleware(s.getPatientsDetailsByCompanyID))
}

func (s *detailsService) patientDetailsRegistration(w http.ResponseWriter, r *http.Request) {
	var patient models.PatientDetails
	var response models.Response
	isOwner := false

	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: check company ID if exists and if caller is related
	compIDs := gcontext.Get(r, "compIDs").([]uint)
	userID := gcontext.Get(r, "id").(uint)

	if len(compIDs) == 0 {
		response.Status = "error"
		response.Message = "Please register your company"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	for _, v := range compIDs {
		if v == patient.CompanyID {
			isOwner = true
		}
	}

	if !isOwner {
		response.Status = "error"
		response.Message = "You do not have permissions to add these data"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	// isValid, errors := handlers.ValidateInputs(patient)
	// if !isValid {
	// 	for _, fieldError := range errors {
	// 		http.Error(w, fieldError, http.StatusBadRequest)
	// 		return
	// 	}
	// }

	patient.CreatedBy = userID

	patient, err = s.patientDetailsRepository.AddPatientDetails(patient)
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

	fmt.Fprintf(w, "Registration of Details - Successful")
}

func (s *detailsService) getPatientsDetailsByCompanyID(w http.ResponseWriter, r *http.Request) {
	var patients []models.PatientDetails
	var response models.Response

	username := gcontext.Get(r, "username").(string)

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	patients, err := s.patientDetailsRepository.GetPatientDetailsByCompanyID(1)
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

func (s *detailsService) getPatientDetailsFull(w http.ResponseWriter, r *http.Request) {
	var patient models.PatientDetails
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

	patient, err = s.patientDetailsRepository.GetPatientDetailsFull(intID)
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
