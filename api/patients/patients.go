package patients

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/middleware"
	"github.com/rehab-backend/internal/pkg/handlers"
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/uptrace/bun"
)

type Service struct {
	*sql.Tx

	dbConnection *bun.DB
}

func NewService() *Service {
	dbConnection := config.ConnectDB()

	// return &Service{}
	return &Service{dbConnection: dbConnection}
}

func (s *Service) RegisterHandlers(route *mux.Router) {

	s.Handle(route)

}

func (s *Service) Handle(route *mux.Router) {

	sub := route.PathPrefix("/patient").Subrouter()

	sub.HandleFunc("/registerPatient", s.patientRegistration)
	sub.HandleFunc("/getPatient", middleware.AuthenticationMiddleware(s.getPatientData)).Methods("GET")
}

func (s *Service) getPatientData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var account models.Account
	var response models.Response

	username := gcontext.Get(r, "username").(string)

	var retrievedAccount models.Account

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		response.Status = "error"
		response.Message = "Please input all required fields."
		json.NewEncoder(w).Encode(response)

		return
	}

	ctx := context.Background()

	err = s.dbConnection.NewSelect().Model(&account).Where("user_id = ?", id).Scan(ctx, &retrievedAccount)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if retrievedAccount.Username != username {
		http.Error(w, "You are not authorized to view this data.", http.StatusBadRequest)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(retrievedAccount)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = username
	response.Response = string("Hello")
	json.NewEncoder(w).Encode(response)

}

func (s *Service) patientRegistration(w http.ResponseWriter, r *http.Request) {
	// currentDate := time.Now().Format("2006-01-02 15:04:05")

	var patient models.PatientPersonal

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

	patient.CreatedOn = time.Now()
	ctx := context.Background()

	tx, err := s.dbConnection.BeginTx(ctx, &sql.TxOptions{})

	_, err = tx.NewInsert().Model(&patient.Account).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		err = tx.Rollback()

		return
	}

	updateField(&patient, "UserID", patient.Account.UserID)

	_, err = tx.NewInsert().Model(&patient.Injury).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		err = tx.Rollback()

		return
	}

	_, err = tx.NewInsert().Model(&patient.Therapy).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		err = tx.Rollback()

		return
	}

	_, err = tx.NewInsert().Model(&patient.DrugTreatment).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		err = tx.Rollback()

		return
	}

	_, err = tx.NewInsert().Model(&patient.PersonalAllergy).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		err = tx.Rollback()

		return
	}

	_, err = tx.NewInsert().Model(&patient.MedicalTherapy).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		err = tx.Rollback()

		return
	}

	_, err = tx.NewInsert().Model(&patient.PersonalDisorder).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		err = tx.Rollback()

		return
	}

	// Commit the transaction if all parts of it succeed
	err = tx.Commit()
	if err != nil {
		panic(err)
	}

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
