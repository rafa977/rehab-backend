package patients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rehab-backend/internal/middleware"
	"github.com/rehab-backend/internal/pkg/handlers"
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/rehab-backend/internal/repository"
)

type diseaseService struct {
	diseaseRepository repository.DiseaseRepository
}

func NewDiseaseService() *diseaseService {
	return &diseaseService{
		diseaseRepository: repository.NewDiseaseService(),
	}
}

func (s *diseaseService) RegisterHandlers(route *mux.Router) {
	s.Handle(route)
}

func (s *diseaseService) Handle(route *mux.Router) {
	sub := route.PathPrefix("/disease").Subrouter()

	sub.HandleFunc("/getDisease", middleware.AuthenticationMiddleware(s.getDisease))
	// sub.HandleFunc("/getAllDysfunctionsPatientID", middleware.AuthenticationMiddleware(s.getAllDysfunctionsPatientID))
	// sub.HandleFunc("/deleteDysfunction", middleware.AuthenticationMiddleware(s.deleteDysfunction))
	// sub.HandleFunc("/updateDysfunction", middleware.AuthenticationMiddleware(s.updateDysfunction))
}

// func (s *dysfunctionService) getAllDysfunctionsPatientID(w http.ResponseWriter, r *http.Request) {
// 	var dysfunctions []models.Dysfunction

// 	id := r.URL.Query().Get("id")
// 	if id == "" {
// 		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
// 		return
// 	}

// 	compIDs := gcontext.Get(r, "compIDs").([]uint)
// 	patientDetailsID, err := strconv.Atoi(id)

// 	isValidPatient, validationError := s.dysfunctionRepository.CheckPatientDetailsOwning(patientDetailsID, compIDs)
// 	if !isValidPatient {
// 		handlers.ProduceErrorResponse(validationError, w, r)
// 		return
// 	}

// 	dysfunctions, err = s.dysfunctionRepository.GetAllDysfunctionsPatientDetailsID(patientDetailsID)
// 	if err != nil {
// 		handlers.ProduceErrorResponse(err.Error(), w, r)
// 		return
// 	}

// 	jsonRetrieved, err := json.Marshal(dysfunctions)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	handlers.ProduceSuccessResponse(string(jsonRetrieved), w, r)
// }

// func (s *dysfunctionService) updateDysfunction(w http.ResponseWriter, r *http.Request) {
// 	var dysfunction models.Dysfunction

// 	err := json.NewDecoder(r.Body).Decode(&dysfunction)
// 	if err != nil {
// 		handlers.ProduceErrorResponse(err.Error(), w, r)
// 		return
// 	}

// 	isValid, errors := handlers.ValidateInputs(dysfunction)
// 	if !isValid {
// 		for _, fieldError := range errors {
// 			handlers.ProduceErrorResponse(fieldError, w, r)
// 			return
// 		}
// 	}

// 	dysfunction, err = s.dysfunctionRepository.UpdateDysfunction(dysfunction)
// 	if err != nil {
// 		var msg string
// 		// if strings.Contains(err.Error(), "users_company_email_key") {
// 		// 	newerr = "user already exists!"
// 		// } else {
// 		// 	newerr = "Bad Request"
// 		// }
// 		handlers.ProduceErrorResponse(msg, w, r)
// 		return
// 	}
// 	handlers.ProduceSuccessResponse("Dysfunction Update - Successful", w, r)
// }

func (s *diseaseService) getDisease(w http.ResponseWriter, r *http.Request) {
	var disease models.Disease

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	intID, err := strconv.Atoi(id)
	// compIDs := gcontext.Get(r, "compIDs").([]uint)

	// validateCompany, validationError := s.diseaseRepository.CheckDysfunctionCompany(compIDs, intID)
	// if !validateCompany {
	// 	handlers.ProduceErrorResponse(validationError, w, r)
	// 	return
	// }

	disease, err = s.diseaseRepository.GetDisease(intID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrieved, err := json.Marshal(disease)
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.ProduceSuccessResponse(string(jsonRetrieved), w, r)
}

// func (s *dysfunctionService) deleteDysfunction(w http.ResponseWriter, r *http.Request) {
// 	var response models.Response

// 	currentDate := time.Now().Format("2006-01-02 15:04:05")
// 	response.Date = currentDate

// 	id := r.URL.Query().Get("id")
// 	if id == "" {
// 		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
// 		return
// 	}

// 	intID, err := strconv.Atoi(id)

// 	_, err = s.dysfunctionRepository.DeleteDysfunction(intID)
// 	if err != nil {
// 		handlers.ProduceErrorResponse(err.Error(), w, r)
// 		return
// 	}
// 	handlers.ProduceSuccessResponse("Dysfunction Delete - Succesfull", w, r)
// }
