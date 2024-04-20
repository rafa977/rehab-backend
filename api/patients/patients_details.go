package patients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rehab/internal/pkg/handlers"
	"rehab/internal/pkg/models"
	"rehab/internal/repository"
	"strconv"
	"strings"

	"rehab/internal/middleware"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
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

	sub.HandleFunc("/addPatientDetails", middleware.AuthenticationMiddleware(s.addPatientDetails))
	sub.HandleFunc("/updatePatientDetails", middleware.AuthenticationMiddleware(s.updatePatientDetails))
	sub.HandleFunc("/getPatientDetails/{id}", middleware.AuthenticationMiddleware(s.getPatientDetails))
	sub.HandleFunc("/getAllPatientDetailsCards/{id}", middleware.AuthenticationMiddleware(s.getAllPatientDetailsCards))

	sub.HandleFunc("/deletePatientDetails", middleware.AuthenticationMiddleware(s.deletePatientDetails))

	//TODO:
	//Give access to a specific account for patient details
}

func (s *detailsService) addPatientDetails(w http.ResponseWriter, r *http.Request) {
	var patientDetails models.PatientDetails

	err := json.NewDecoder(r.Body).Decode(&patientDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: check company ID if exists and if caller is related
	compIDs := handlers.GetCompany(r)
	userID := handlers.GetAccount(r)

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

	patientDetails.AddedByID = userID
	patientDetails.LastUpdatedByID = userID

	patientDetails, err = s.patientDetailsRepository.AddPatientDetails(patientDetails)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	// Add Access Rights to Account who creates the record
	var permissions models.PatientDetailsPermission
	permissions.AccountID = userID
	permissions.PatientDetailsID = patientDetails.ID
	permissions.Access = true

	permissions, err = s.patientDetailsRepository.AddPatientDetailsPermission(permissions)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	handlers.ProduceSuccessResponse("Registration of Details - Successful", "", w, r)
}

func (s *detailsService) getPatientDetails(w http.ResponseWriter, r *http.Request) {
	var patientDetails models.PatientDetails
	var patient models.Patient

	// Current account id
	accountId := gcontext.Get(r, "id").(uint)

	params := mux.Vars(r)

	id := params["id"]
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	// Convert string parameter to uint
	patientCardID, err := handlers.ConvertStrToUint(id)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	roleID := handlers.GetRole(r)

	if roleID != 1 {
		// Check permissions
		permissions, err := s.patientDetailsRepository.GetPatientDetailsPermission(patientCardID, accountId)
		if err != nil {
			handlers.ProduceErrorResponse(err.Error(), w, r)
			return
		}

		if permissions.ID == 0 {
			handlers.ProduceErrorResponse("You do not have access to these data", w, r)
			return
		}

		if permissions.Access == false {
			handlers.ProduceErrorResponse("You do not have access to these data.", w, r)
			return
		}
	}

	patientDetails, err = s.patientDetailsRepository.GetPatientDetailsFull(patientCardID)
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "record not found") {
			msg = "No patient cards created"
		} else {
			msg = "Bad Request"
		}
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}

	patient, err = s.patientRepository.GetPatient(patientDetails.PatientID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	ownsCompany, errMsg := handlers.ValidateCompany(patient.CompanyID, r)
	if !ownsCompany {
		handlers.ProduceErrorResponse(errMsg, w, r)
		return
	}

	handlers.ProduceJsonSuccessResponse(patientDetails, "", w, r)
}

// Get all Patient Detail Cards based on Patient ID
func (s *detailsService) getAllPatientDetailsCards(w http.ResponseWriter, r *http.Request) {
	var patientDetails []models.PatientDetails
	var patient models.Patient

	params := mux.Vars(r)

	id := params["id"]
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	// Current account id
	accountId := gcontext.Get(r, "id").(uint)

	// Convert string parameter to uint
	patientID, err := handlers.ConvertStrToUint(id)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	roleID := handlers.GetRole(r)

	if roleID != 1 {
		// Check permissions

		// Bring all cards that current account has access

		permissions, err := s.patientDetailsRepository.GetPatientDetailsForEmployeeID(patientID, accountId)
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

		jsonRetrievedAccount, err := json.Marshal(permissions)
		if err != nil {
			fmt.Println(err)

			return
		}

		handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
		return
	}

	patient, err = s.patientRepository.GetPatient(patientID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	ownsCompany, errMsg := handlers.ValidateCompany(patient.CompanyID, r)
	if !ownsCompany {
		handlers.ProduceErrorResponse(errMsg, w, r)
		return
	}

	patientDetails, err = s.patientDetailsRepository.GetPatientDetailsByPatientID(patientID)
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

	// jsonRetrievedAccount, err := json.Marshal(patientDetails)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
	handlers.ProduceJsonSuccessResponse(patientDetails, "", w, r)
}

func (s *detailsService) updatePatientDetails(w http.ResponseWriter, r *http.Request) {
	var patient models.PatientDetails

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
		var errMsg string
		if strings.Contains(err.Error(), "users_company_email_key") {
			errMsg = "user already exists!"
		} else {
			errMsg = "Bad Request"
		}

		handlers.ProduceErrorResponse(errMsg, w, r)
		return
	}

	handlers.ProduceSuccessResponse("Update of Patient Details - Successful", "", w, r)
}

func (s *detailsService) deletePatientDetails(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	intID, err := strconv.Atoi(id)

	_, err = s.patientDetailsRepository.DeletePatientDetails(intID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}
	handlers.ProduceSuccessResponse("Patient Details Delete - Succesfull", "", w, r)
}
