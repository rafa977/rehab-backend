package patients

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	config "github.com/rehab-backend/config/database"
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

	fmt.Println(patient.Account.UserID)

	// var account models.Account
	// account = patient.Account

	// var injury models.Injury
	// injury = patient.Injury[0]

	// var medTherapy models.MedicalTherapy
	// medTherapy = patient.MedicalTherapy[0]

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

	fmt.Println(patient.Account)

	for i := range patient.Injury {
		patient.Injury[i].UserID = patient.Account.UserID
	}
	for i := range patient.MedicalTherapy {
		patient.MedicalTherapy[i].UserID = patient.Account.UserID
	}
	for i := range patient.PersonalDisorder {
		patient.PersonalDisorder[i].UserID = patient.Account.UserID
	}
	for i := range patient.DrugTreatment {
		patient.DrugTreatment[i].UserID = patient.Account.UserID
	}

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

	err = tx.Commit()

	fmt.Println(patient)

	fmt.Fprintf(w, "Registration of Account - Successful")
}
