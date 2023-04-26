package patients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rehab-backend/internal/middleware"
	"github.com/rehab-backend/internal/pkg/handlers"
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/rehab-backend/internal/repository"
)

type dysfunctionService struct {
	clinicalRepository repository.DysfunctionRepository
}

func NewDysfunctionService() *dysfunctionService {
	return &dysfunctionService{
		clinicalRepository: repository.NewDysfunctionService(),
	}
}

func (s *dysfunctionService) RegisterHandlers(route *mux.Router) {
	s.Handle(route)
}

func (s *dysfunctionService) Handle(route *mux.Router) {
	sub := route.PathPrefix("/dysfunction").Subrouter()

	sub.HandleFunc("/getDysfunction", middleware.AuthenticationMiddleware(s.getDysfunction))
	sub.HandleFunc("/deleteDysfunction", middleware.AuthenticationMiddleware(s.deleteDysfunction))
	sub.HandleFunc("/addDysfunction", middleware.AuthenticationMiddleware(s.addDysfunction))
	sub.HandleFunc("/updateDysfunction", middleware.AuthenticationMiddleware(s.updateDysfunction))
}

func (s *dysfunctionService) addDysfunction(w http.ResponseWriter, r *http.Request) {
	var dysfunction models.Dysfunction

	err := json.NewDecoder(r.Body).Decode(&dysfunction)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	isValid, errors := handlers.ValidateInputs(dysfunction)
	if !isValid {
		for _, fieldError := range errors {
			handlers.ProduceErrorResponse(fieldError, w, r)
			return
		}
	}

	// TODO: check company ID if exists and if caller is related
	isOwner, ownerError := handlers.ValidateCompany(dysfunction.CompanyID, r)
	if !isOwner {
		handlers.ProduceErrorResponse(ownerError, w, r)
		return
	}

	isValidPatient, validationError := s.clinicalRepository.CheckPatientDetails(dysfunction.PatientDetailsID, dysfunction.CompanyID)
	if !isValidPatient {
		handlers.ProduceErrorResponse(validationError, w, r)
		return
	}

	dysfunction, err = s.clinicalRepository.AddDysfunction(dysfunction)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}
	handlers.ProduceSuccessResponse("Dysfunction Registration - Successful", w, r)
}

func (s *dysfunctionService) updateDysfunction(w http.ResponseWriter, r *http.Request) {
	var dysfunction models.Dysfunction

	err := json.NewDecoder(r.Body).Decode(&dysfunction)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	isValid, errors := handlers.ValidateInputs(dysfunction)
	if !isValid {
		for _, fieldError := range errors {
			handlers.ProduceErrorResponse(fieldError, w, r)
			return
		}
	}

	dysfunction, err = s.clinicalRepository.UpdateDysfunction(dysfunction)
	if err != nil {
		var msg string
		// if strings.Contains(err.Error(), "users_company_email_key") {
		// 	newerr = "user already exists!"
		// } else {
		// 	newerr = "Bad Request"
		// }
		handlers.ProduceErrorResponse(msg, w, r)
		return
	}
	handlers.ProduceSuccessResponse("Dysfunction Update - Successful", w, r)
}

func (s *dysfunctionService) getDysfunction(w http.ResponseWriter, r *http.Request) {
	var dysfunction models.Dysfunction

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	intID, err := strconv.Atoi(id)

	dysfunction, err = s.clinicalRepository.GetDysfunction(intID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(dysfunction)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrievedAccount), w, r)
}

func (s *dysfunctionService) deleteDysfunction(w http.ResponseWriter, r *http.Request) {
	var response models.Response

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	intID, err := strconv.Atoi(id)

	_, err = s.clinicalRepository.DeleteDysfunction(intID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}
	handlers.ProduceSuccessResponse("Dysfunction Delete - Succesfull", w, r)
}
