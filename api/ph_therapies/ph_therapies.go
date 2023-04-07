package ph_therapies

import (
	"encoding/json"
	"fmt"
	"net/http"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rehab-backend/internal/middleware"
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/rehab-backend/internal/repository"
)

type phTherapyService struct {
	phTherapyRepository      repository.PhTherapyRepository
	patientDetailsRepository repository.PatientDetailRepository
}

func NewPhTherapiesService() *phTherapyService {

	return &phTherapyService{
		phTherapyRepository:      repository.NewPhTherapyService(),
		patientDetailsRepository: repository.NewPatientDetailsService(),
	}
}

func (s *phTherapyService) RegisterPhTherapiesHandlers(route *mux.Router) {

	s.DetailHandle(route)

}

func (s *phTherapyService) DetailHandle(route *mux.Router) {

	sub := route.PathPrefix("/ph_therapy").Subrouter()

	sub.HandleFunc("/addDysfunction", middleware.AuthenticationMiddleware(s.addDysfunction))
	// sub.HandleFunc("/updatePatientDetails", middleware.AuthenticationMiddleware(s.updatePatientDetails))
	// sub.HandleFunc("/getPatientDetailsFull", middleware.AuthenticationMiddleware(s.getPatientDetailsFull))
	// sub.HandleFunc("/getPatientsDetailsByCompanyID", middleware.AuthenticationMiddleware(s.getPatientsDetailsByCompanyID))

	sub.HandleFunc("/addPhTherapy", middleware.AuthenticationMiddleware(s.addPhTherapy))
}

func (s *phTherapyService) addDysfunction(w http.ResponseWriter, r *http.Request) {
	var dysfunction models.Dysfunction
	var response models.Response

	err := json.NewDecoder(r.Body).Decode(&dysfunction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: check if patient ID has been registered under the company
	// get patient details by ID
	var patientDetails models.PatientDetails

	// TODO: check company ID if exists and if caller is related
	compIDs := gcontext.Get(r, "compIDs").([]uint)

	isOwner := false

	for _, v := range compIDs {
		if v == dysfunction.CompanyID {
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

	patientDetails = s.patientDetailsRepository.GetPatientDetailsByIdAndCompanyID(int(dysfunction.PatientDetailsID), int(dysfunction.CompanyID))
	fmt.Println(patientDetails.CompanyID)

	if patientDetails.ID == 0 {
		response.Status = "error"
		response.Message = "You must first register your patient to continue to therapies"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	dysfunction, err = s.phTherapyRepository.AddDysfunction(dysfunction)
	if err != nil {
		var newerr string
		response.Status = "error"
		response.Message = newerr
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	} else {
		response.Status = "success"
		response.Message = ""
		response.Response = "Registration Successful"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func (s *phTherapyService) addPhTherapy(w http.ResponseWriter, r *http.Request) {
	var phTherapy models.PhTherapy
	var response models.Response
	var dysfunction models.Dysfunction

	err := json.NewDecoder(r.Body).Decode(&phTherapy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	isOwner := false

	dysfunction, err = s.phTherapyRepository.GetDysfunction(phTherapy.DysfunctionID)
	if dysfunction.CompanyID > 0 {
		// TODO: check company ID if exists and if caller is related
		compIDs := gcontext.Get(r, "compIDs").([]uint)

		for _, v := range compIDs {
			if v == dysfunction.CompanyID {
				isOwner = true
			}
		}
	}

	if !isOwner {
		response.Status = "error"
		response.Message = "You do not have permissions to access these data"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	phTherapy.PatientDetailsID = dysfunction.PatientDetailsID

	phTherapy, err = s.phTherapyRepository.AddPhTherapy(phTherapy)
	if err != nil {
		response.Status = "error"
		response.Message = "Somethign went wrong"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Status = "success"
	response.Message = "Registration of Therapy Successful"
	response.Response = ""
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}

// func (s *detailsService) getPatientsDetailsByCompanyID(w http.ResponseWriter, r *http.Request) {
// 	var patients []models.PatientDetails
// 	var response models.Response

// 	username := gcontext.Get(r, "username").(string)

// 	currentDate := time.Now().Format("2006-01-02 15:04:05")
// 	response.Date = currentDate

// 	patients, err := s.patientDetailsRepository.GetPatientDetailsByCompanyID(1)
// 	if err != nil {
// 		response.Status = "error"
// 		response.Message = "Unknown Username or Password"
// 		response.Response = ""
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	jsonRetrievedAccount, err := json.Marshal(patients)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	response.Status = "success"
// 	response.Message = username
// 	response.Response = string(jsonRetrievedAccount)
// 	json.NewEncoder(w).Encode(response)

// }

// func (s *detailsService) getPatientDetailsFull(w http.ResponseWriter, r *http.Request) {
// 	var patient models.PatientDetails
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

// 	patient, err = s.patientDetailsRepository.GetPatientDetailsFull(intID)
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

// func (s *detailsService) updatePatientDetails(w http.ResponseWriter, r *http.Request) {
// 	var patient models.PatientDetails
// 	var response models.Response

// 	err := json.NewDecoder(r.Body).Decode(&patient)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// isValid, errors := handlers.ValidateInputs(patient)
// 	// if !isValid {
// 	// 	for _, fieldError := range errors {
// 	// 		http.Error(w, fieldError, http.StatusBadRequest)
// 	// 		return
// 	// 	}
// 	// }

// 	patient, err = s.patientDetailsRepository.UpdatePatientDetails(patient)
// 	if err != nil {
// 		var newerr string
// 		if strings.Contains(err.Error(), "users_company_email_key") {
// 			newerr = "user already exists!"
// 		} else {
// 			newerr = "Bad Request"
// 		}
// 		response.Status = "error"
// 		response.Message = newerr
// 		response.Response = ""
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	fmt.Fprintf(w, "Registration of Account - Successful")
// }
