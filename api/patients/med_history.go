package patients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rehab/internal/pkg/handlers"
	"rehab/internal/pkg/models"
	"rehab/internal/repository"
	"strings"

	"rehab/internal/middleware"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type medHistoryService struct {
	medHistoryRepository repository.MedHistoryRepository
	patientRepository    repository.PatientRepository
}

func NewMedHistoryService() *medHistoryService {

	return &medHistoryService{
		medHistoryRepository: repository.NewMedHistoryService(),
		patientRepository:    repository.NewPatientService(),
	}
}

func (s *medHistoryService) RegisterMedHistoryHandlers(route *mux.Router) {
	s.MedHistoryHandle(route)
}

func (s *medHistoryService) MedHistoryHandle(route *mux.Router) {

	sub := route.PathPrefix("/medHistory").Subrouter()

	sub.HandleFunc("/addMedHistory", middleware.AuthenticationMiddleware(s.addMedHistory))
	sub.HandleFunc("/updateMedHistory", middleware.AuthenticationMiddleware(s.updateMedHistory))
	sub.HandleFunc("/getMedHistory/{id}/{type}", middleware.AuthenticationMiddleware(s.getMedHistory))
	// sub.HandleFunc("/getAllMedHistoryCards/{id}", middleware.AuthenticationMiddleware(s.getAllMedHistoryCards))

	// sub.HandleFunc("/deleteMedHistory", middleware.AuthenticationMiddleware(s.deleteMedHistory))

	//TODO:
	//Give access to a specific account for patient details
}

func (s *medHistoryService) addMedHistory(w http.ResponseWriter, r *http.Request) {
	var medHistory models.MedHistory

	err := json.NewDecoder(r.Body).Decode(&medHistory)
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
	isPatientValid, validationError := s.patientRepository.CheckPatient(medHistory.PatientID, compIDs)
	if !isPatientValid {
		handlers.ProduceErrorResponse(validationError, w, r)
		return
	}

	medHistory.AddedByID = userID
	// medHistory.LastUpdatedByID = userID

	medHistory, err = s.medHistoryRepository.AddMedHistory(medHistory)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	// Add Access Rights to Account who creates the record
	var permissions models.MedHistoryPermission
	permissions.AccountID = userID
	permissions.MedHistoryID = medHistory.ID
	permissions.Access = true

	permissions, err = s.medHistoryRepository.AddMedHistoryPermission(permissions)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	handlers.ProduceSuccessResponse("Registration of Medical History - Successful", "", w, r)
}

func (s *medHistoryService) getMedHistory(w http.ResponseWriter, r *http.Request) {
	var medHistory models.MedHistory
	var patient models.Patient

	// Current account id
	accountId := gcontext.Get(r, "id").(uint)

	params := mux.Vars(r)

	id := params["id"]
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	historyType := params["type"]

	// Convert string parameter to uint
	medHistoryID, err := handlers.ConvertStrToUint(id)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	roleID := handlers.GetRole(r)

	if roleID != 1 {
		// Check permissions
		permissions, err := s.medHistoryRepository.GetMedHistoryPermission(medHistoryID, accountId)
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

	if historyType == "a" {
		medHistory, err = s.medHistoryRepository.GetMedicalHistoryFull(medHistoryID)
		if err != nil {
			var msg string
			if strings.Contains(err.Error(), "record not found") {
				msg = "No medical history created"
			} else {
				msg = "Bad Request"
			}
			handlers.ProduceErrorResponse(msg, w, r)
			return
		}
	} else {
		var actualType string

		switch historyType {
		case "allergy":
			actualType = "PersonalAllergies.Allergy"
			break
		case "medtherapy":
			actualType = "MedicalTherapies"
			break
		case "injury":
			actualType = "Injuries"
			break
		case "therapy":
			actualType = "Therapies"
			break
		case "drug":
			actualType = "DrugTreatments.Drug"
			break
		case "disorder":
			actualType = "PersonalDisorders.Disorder"
			break
		case "surgery":
			actualType = "Surgeries"
			break
		}

		medHistory, err = s.medHistoryRepository.GetMedicalHistorySpecific(medHistoryID, actualType)
		if err != nil {
			var msg string
			if strings.Contains(err.Error(), "record not found") {
				msg = "No medical history created"
			} else {
				msg = "Bad Request"
			}
			handlers.ProduceErrorResponse(msg, w, r)
			return
		}
	}

	patient, err = s.patientRepository.GetPatient(medHistory.PatientID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	ownsCompany, errMsg := handlers.ValidateCompany(patient.CompanyID, r)
	if !ownsCompany {
		handlers.ProduceErrorResponse(errMsg, w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(medHistory)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
}

// // Get all Patient Detail Cards based on Patient ID
// func (s *medHistoryService) getAllPatientDetailsCards(w http.ResponseWriter, r *http.Request) {
// 	var patientDetails []models.PatientDetails
// 	var patient models.Patient

// 	params := mux.Vars(r)

// 	id := params["id"]
// 	if id == "" {
// 		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
// 		return
// 	}

// 	// Current account id
// 	accountId := gcontext.Get(r, "id").(uint)

// 	// Convert string parameter to uint
// 	patientID, err := handlers.ConvertStrToUint(id)
// 	if err != nil {
// 		handlers.ProduceErrorResponse(err.Error(), w, r)
// 		return
// 	}

// 	roleID := handlers.GetRole(r)

// 	if roleID != 1 {
// 		// Check permissions

// 		// Bring all cards that current account has access

// 		permissions, err := s.patientDetailsRepository.GetPatientDetailsForEmployeeID(patientID, accountId)
// 		if err != nil {
// 			var msg string
// 			if strings.Contains(err.Error(), "record not found") {
// 				msg = "You are not authorized to access these data!"
// 			} else {
// 				msg = "Bad Request"
// 			}
// 			handlers.ProduceErrorResponse(msg, w, r)
// 			return
// 		}

// 		jsonRetrievedAccount, err := json.Marshal(permissions)
// 		if err != nil {
// 			fmt.Println(err)

// 			return
// 		}

// 		handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
// 		return
// 	}

// 	patient, err = s.patientRepository.GetPatient(patientID)
// 	if err != nil {
// 		handlers.ProduceErrorResponse(err.Error(), w, r)
// 		return
// 	}

// 	ownsCompany, errMsg := handlers.ValidateCompany(patient.CompanyID, r)
// 	if !ownsCompany {
// 		handlers.ProduceErrorResponse(errMsg, w, r)
// 		return
// 	}

// 	patientDetails, err = s.patientDetailsRepository.GetPatientDetailsByPatientID(patientID)
// 	if err != nil {
// 		var msg string
// 		if strings.Contains(err.Error(), "record not found") {
// 			msg = "You are not authorized to access these data!"
// 		} else {
// 			msg = "Bad Request"
// 		}
// 		handlers.ProduceErrorResponse(msg, w, r)
// 		return
// 	}

// 	jsonRetrievedAccount, err := json.Marshal(patientDetails)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
// }

func (s *medHistoryService) updateMedHistory(w http.ResponseWriter, r *http.Request) {
	var medHistory models.MedHistory

	err := json.NewDecoder(r.Body).Decode(&medHistory)
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

	medHistory, err = s.medHistoryRepository.UpdateMedicalHistory(medHistory)
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

	handlers.ProduceSuccessResponse("Update of Medical History - Successful", "", w, r)
}

// func (s *medHistoryService) deletePatientDetails(w http.ResponseWriter, r *http.Request) {

// 	id := r.URL.Query().Get("id")
// 	if id == "" {
// 		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
// 		return
// 	}

// 	intID, err := strconv.Atoi(id)

// 	_, err = s.patientDetailsRepository.DeletePatientDetails(intID)
// 	if err != nil {
// 		handlers.ProduceErrorResponse(err.Error(), w, r)
// 		return
// 	}
// 	handlers.ProduceSuccessResponse("Patient Details Delete - Succesfull", "", w, r)
// }
