package ph_therapies

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rehab-backend/internal/middleware"
	"github.com/rehab-backend/internal/pkg/handlers"
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
	sub.HandleFunc("/getPhTherapy", middleware.AuthenticationMiddleware(s.getPhTherapy))
	sub.HandleFunc("/getPhTherapiesByCompID", middleware.AuthenticationMiddleware(s.getPhTherapiesByCompanyID))
}

func (s *phTherapyService) addDysfunction(w http.ResponseWriter, r *http.Request) {
	var dysfunction models.Dysfunction
	var patientDetails models.PatientDetails
	var isOwner = false

	err := json.NewDecoder(r.Body).Decode(&dysfunction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	compIDs := gcontext.Get(r, "compIDs").([]uint)
	for _, v := range compIDs {
		if v == dysfunction.CompanyID {
			isOwner = true
		}
	}
	if !isOwner {
		handlers.ProduceErrorResponse("You do not have permissions to add these data", w, r)
		return
	}

	patientDetails, err = s.patientDetailsRepository.GetPatientDetails(int(dysfunction.PatientDetailsID))
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}
	fmt.Println(patientDetails.Patient.CompanyID)

	dysfunction, err = s.phTherapyRepository.AddDysfunction(dysfunction)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
	} else {
		handlers.ProduceSuccessResponse("Registration of Dysfunction - Successful", w, r)
	}
}

func (s *phTherapyService) getDysfunction(w http.ResponseWriter, r *http.Request) {
	var dysfnuction models.Dysfunction

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}
	intID, err := strconv.Atoi(id)

	dysfnuction, err = s.phTherapyRepository.GetDysfunction(intID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(dysfnuction)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *phTherapyService) addPhTherapy(w http.ResponseWriter, r *http.Request) {
	var phTherapy models.PhTherapy
	var dysfunction models.Dysfunction
	var isOwner = false

	err := json.NewDecoder(r.Body).Decode(&phTherapy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dysfunction, err = s.phTherapyRepository.GetDysfunction(int(phTherapy.DysfunctionID))
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
		handlers.ProduceErrorResponse("You do not have permissions to access these data", w, r)
		return
	}

	phTherapy.PatientDetailsID = dysfunction.PatientDetailsID

	phTherapy, err = s.phTherapyRepository.AddPhTherapy(phTherapy)
	if err != nil {
		handlers.ProduceErrorResponse("Something went wrong", w, r)
		return
	}
	handlers.ProduceSuccessResponse("Registration of Therapy - Successful", w, r)
}

func (s *phTherapyService) getPhTherapy(w http.ResponseWriter, r *http.Request) {
	var phTherapy models.PhTherapy

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}
	intID, err := strconv.Atoi(id)

	phTherapy, err = s.phTherapyRepository.GetPhTherapy(intID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(phTherapy)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *phTherapyService) getPhTherapiesByCompanyID(w http.ResponseWriter, r *http.Request) {
	var phTherapies []models.PhTherapy

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}
	intID, err := strconv.Atoi(id)

	phTherapies, err = s.phTherapyRepository.GetPhTherapiesByCompanyID(intID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(phTherapies)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

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
