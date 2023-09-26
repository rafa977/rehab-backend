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
	sub.HandleFunc("/getDiseasesPDID", middleware.AuthenticationMiddleware(s.getDiseasesPDID))
	sub.HandleFunc("/deleteDisease", middleware.AuthenticationMiddleware(s.deleteDisease))
	// sub.HandleFunc("/updateDysfunction", middleware.AuthenticationMiddleware(s.updateDysfunction))
}

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

func (s *diseaseService) getDiseasesPDID(w http.ResponseWriter, r *http.Request) {
	var disease []models.Disease

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	// compIDs := gcontext.Get(r, "compIDs").([]uint)
	patientDetailsID, err := strconv.Atoi(id)

	// isValidPatient, validationError := s.diseaseRepository.CheckPatientDetailsOwning(patientDetailsID, compIDs)
	// if !isValidPatient {
	// 	handlers.ProduceErrorResponse(validationError, w, r)
	// 	return
	// }

	disease, err = s.diseaseRepository.GetAllDiseasesPatientDetailsID(patientDetailsID)
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

func (s *diseaseService) deleteDisease(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	intID, err := strconv.Atoi(id)

	_, err = s.diseaseRepository.DeleteDisease(intID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}
	handlers.ProduceSuccessResponse("Disease Delete - Succesfull", w, r)
}
