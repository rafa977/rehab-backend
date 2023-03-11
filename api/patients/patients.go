package patients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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
	repository repository.PatientRepository
}

func NewService() *service {

	return &service{repository: repository.NewPatientService()}
}

func (s *service) RegisterHandlers(route *mux.Router) {

	s.Handle(route)

}

func (s *service) Handle(route *mux.Router) {

	sub := route.PathPrefix("/patient").Subrouter()

	sub.HandleFunc("/registerPatient", middleware.AuthenticationMiddleware(s.patientRegistration)).Methods("POST")
	sub.HandleFunc("/getPatient", middleware.AuthenticationMiddleware(s.getPatientData)).Methods("GET")
	sub.HandleFunc("/getAllPatients", middleware.AuthenticationMiddleware(s.getAllPatients)).Methods("GET")

}

func (s *service) getAllPatients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var patients []models.Patient
	var response models.Response

	username := gcontext.Get(r, "username").(string)

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	patients, err := s.repository.GetAllPatients()
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

	patient, err := s.repository.GetPatient(1)
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

func (s *service) patientRegistration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var patient models.Patient
	var response models.Response
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
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

	patient, err = s.repository.AddPatient(patient)
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

	// tx, err := s.dbConnection.BeginTx(ctx, &sql.TxOptions{})

	// _, err = tx.NewInsert().Model(&patient.Account).Exec(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	err = tx.Rollback()

	// 	return
	// }

	// updateField(&patient, "UserID", patient.Account.UserID)

	// _, err = tx.NewInsert().Model(&patient.Injury).Exec(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	err = tx.Rollback()

	// 	return
	// }

	// _, err = tx.NewInsert().Model(&patient.Therapy).Exec(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	err = tx.Rollback()

	// 	return
	// }

	// _, err = tx.NewInsert().Model(&patient.DrugTreatment).Exec(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	err = tx.Rollback()

	// 	return
	// }

	// _, err = tx.NewInsert().Model(&patient.PersonalAllergy).Exec(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	err = tx.Rollback()

	// 	return
	// }

	// _, err = tx.NewInsert().Model(&patient.MedicalTherapy).Exec(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	err = tx.Rollback()

	// 	return
	// }

	// _, err = tx.NewInsert().Model(&patient.PersonalDisorder).Exec(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	err = tx.Rollback()

	// 	return
	// }

	// // Commit the transaction if all parts of it succeed
	// err = tx.Commit()
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Fprintf(w, "Registration of Account - Successful")
}

// Function to update the UserID field in the PatientPersonal struct
func updateField(p *models.PatientPersonal, key string, value string) {
	// Get the reflect value of the PatientPersonal struct
	reflectValue := reflect.ValueOf(p).Elem()

	// Loop through each field in the struct
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)

		// Check if the field is a slice
		if field.Kind() == reflect.Slice {
			// Loop through each element in the slice
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)

				// Check if the element has a UserID field
				userIDField := elem.FieldByName(key)
				// if userIDField.IsValid() && userIDField.String() == key {
				if userIDField.IsValid() {
					// Update the UserID field with the new value
					userIDField.SetString(value)
				}
			}
		} else {
			// Check if the field has a UserID field
			userIDField := field.FieldByName(key)
			// if userIDField.IsValid() && userIDField.String() == key {
			if userIDField.IsValid() {
				// Update the UserID field with the new value
				userIDField.SetString(value)
			}
		}
	}
}
