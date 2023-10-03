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
	sub.HandleFunc("/deletePatient/{id}", middleware.AuthenticationMiddleware(s.deletePatient))
	sub.HandleFunc("/getPatient/{id}", middleware.AuthenticationMiddleware(s.getPatientData))

	sub.HandleFunc("/searchPatient", middleware.AuthenticationMiddleware(s.getPatientDataKeyword))
	sub.HandleFunc("/getAllPatients", middleware.AuthenticationMiddleware(s.getAllPatients))
	sub.HandleFunc("/getAllPatientsCompanyID/{id}", middleware.AuthenticationMiddleware(s.getAllPatientsByCompanyId))
	sub.HandleFunc("/getAllPatientsDetails", middleware.AuthenticationMiddleware(s.getAllPatientsDetails))
}

func (s *service) patientRegistration(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient
	// var signature models.Signature

	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	// err = json.NewDecoder(r.Body).Decode(&signature)
	// if err != nil {
	// 	handlers.ProduceErrorResponse(err.Error(), w, r)
	// 	return
	// }

	// TODO: check company ID if exists and if caller is related
	isOwner, ownerError := handlers.ValidateCompany(patient.CompanyID, r)
	if !isOwner {
		handlers.ProduceErrorResponse(ownerError, w, r)
		return
	}

	id := gcontext.Get(r, "id").(uint)
	patient.AddedByID = id

	isValid, errors := handlers.ValidateInputs(patient)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	patient, err = s.patientRepository.AddPatient(patient)
	if err != nil {
		var msg string
		if strings.Contains(err.Error(), "idx_companyid_amka") {
			msg = "Patient with same AMKA already exists!"
		} else if strings.Contains(err.Error(), "fk_patients_company") {
			msg = "Please register your company"
		} else {
			msg = "Bad Request"
		}
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}

	handlers.ProduceSuccessResponse("Registration of Account - Successful", fmt.Sprint(patient.ID), w, r)
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

	// TODO: check company ID if exists and if caller is related
	compIDs := handlers.GetCompany(r)

	if len(compIDs) == 0 {
		handlers.ProduceErrorResponse("Please register your company", w, r)
		return
	}

	// check patient id exists and is under same company
	isPatientValid, validationError := s.patientRepository.CheckPatient(patient.ID, compIDs)
	if !isPatientValid {
		handlers.ProduceErrorResponse(validationError, w, r)
		return
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
	handlers.ProduceSuccessResponse("Update of Account - Successful", "", w, r)
}

func (s *service) deletePatient(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id := params["id"]
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	// Convert string parameter to uint
	patientID, err := handlers.ConvertStrToUint(id)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	// TODO: check company ID if exists and if caller is related
	compIDs := handlers.GetCompany(r)
	if len(compIDs) == 0 {
		handlers.ProduceErrorResponse("Please register your company", w, r)
		return
	}

	// check patient id exists and is under same company
	isPatientValid, validationError := s.patientRepository.CheckPatient(patientID, compIDs)
	if !isPatientValid {
		handlers.ProduceErrorResponse(validationError, w, r)
		return
	}

	_, err = s.patientRepository.DeletePatient(patientID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	handlers.ProduceSuccessResponse("Delete of Patient - Successful", "", w, r)
}

func (s *service) getPatientData(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient

	params := mux.Vars(r)

	id := params["id"]
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	// Convert string parameter to uint
	patientID, err := handlers.ConvertStrToUint(id)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	patient, err = s.patientRepository.GetPatient(patientID)
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

	ownsCompany, errMsg := handlers.ValidateCompany(patient.CompanyID, r)
	if !ownsCompany {
		handlers.ProduceErrorResponse(errMsg, w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(patient)
	if err != nil {
		fmt.Println(err)
		return
	}

	roleID := handlers.GetRole(r)

	if roleID == 2 {
		var patientEmployee models.PatientEmployee
		err := json.Unmarshal(jsonRetrievedAccount, &patientEmployee)
		if err != nil {
			handlers.ProduceErrorResponse(err.Error(), w, r)
			return
		}

		newJsonPatient, err := json.Marshal(patientEmployee)
		if err != nil {
			fmt.Println(err)
			return
		}

		handlers.ProduceSuccessResponse(string(newJsonPatient), "", w, r)
	} else if roleID == 1 {
		handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
	}
}

func (s *service) getAllPatientsByCompanyId(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	roleID := gcontext.Get(r, "roleID").(uint)

	id := params["id"]
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	// Convert string parameter to uint
	companyID, err := handlers.ConvertStrToUint(id)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	// Check if caller belongs to company
	ownsCompany, errMsg := handlers.ValidateCompany(companyID, r)
	if !ownsCompany {
		handlers.ProduceErrorResponse(errMsg, w, r)
		return
	}

	patients, err := s.patientRepository.GetAllPatientsByCompanyId(companyID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(patients)
	if err != nil {
		handlers.ProduceErrorResponse("Error on converting retrived data.", w, r)
		return
	}

	if roleID == 2 {
		var patientEmployees []models.PatientEmployee
		err := json.Unmarshal(jsonRetrievedAccount, &patientEmployees)
		if err != nil {
			handlers.ProduceErrorResponse(err.Error(), w, r)
			return
		}

		newJsonPatients, err := json.Marshal(patientEmployees)
		if err != nil {
			fmt.Println(err)
			return
		}

		handlers.ProduceSuccessResponse(string(newJsonPatients), "", w, r)
	} else if roleID == 1 {
		handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
	}
}

func (s *service) getAllPatients(w http.ResponseWriter, r *http.Request) {

	compIDs := handlers.GetCompany(r)

	patients, err := s.patientRepository.GetAllPatients(compIDs)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(patients)
	if err != nil {
		handlers.ProduceErrorResponse("Error on converting retrived data.", w, r)
		return
	}

	roleID := handlers.GetRole(r)

	if roleID == 2 {
		var patientEmployees []models.PatientEmployee
		err := json.Unmarshal(jsonRetrievedAccount, &patientEmployees)
		if err != nil {
			handlers.ProduceErrorResponse(err.Error(), w, r)
			return
		}

		newJsonPatients, err := json.Marshal(patientEmployees)
		if err != nil {
			fmt.Println(err)
			return
		}

		handlers.ProduceSuccessResponse(string(newJsonPatients), "", w, r)
	} else if roleID == 1 {
		handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
	}
}

func (s *service) getAllPatientsDetails(w http.ResponseWriter, r *http.Request) {

	compIDs := handlers.GetCompany(r)

	patients, err := s.patientRepository.GetAllPatientsDetails(compIDs)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(patients)
	if err != nil {
		handlers.ProduceErrorResponse("Error on converting retrived data.", w, r)
		return
	}

	roleID := handlers.GetRole(r)

	if roleID == 2 {
		var patientEmployees []models.PatientEmployee
		err := json.Unmarshal(jsonRetrievedAccount, &patientEmployees)
		if err != nil {
			handlers.ProduceErrorResponse(err.Error(), w, r)
			return
		}

		newJsonPatients, err := json.Marshal(patientEmployees)
		if err != nil {
			fmt.Println(err)
			return
		}

		handlers.ProduceSuccessResponse(string(newJsonPatients), "", w, r)
	} else if roleID == 1 {
		handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
	}
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

	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), "", w, r)
}
