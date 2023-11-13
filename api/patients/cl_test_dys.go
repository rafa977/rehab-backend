package patients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rehab/internal/pkg/models"
	"rehab/internal/repository"
	"strconv"
	"time"

	"rehab/internal/middleware"
	"rehab/internal/pkg/handlers"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type clTestDisService struct {
	clinicalRepository repository.ClTestDisRepository
}

func NewClTestDisService() *clTestDisService {
	return &clTestDisService{
		clinicalRepository: repository.NewClTestDisService(),
	}
}

func (s *clTestDisService) RegisterHandlers(route *mux.Router) {
	s.Handle(route)
}

func (s *clTestDisService) Handle(route *mux.Router) {
	sub := route.PathPrefix("/clinicalTestDisease").Subrouter()

	sub.HandleFunc("/getClTestDis", middleware.AuthenticationMiddleware(s.getClTestDis))
	sub.HandleFunc("/getClTestDisByDisID/{id}", middleware.AuthenticationMiddleware(s.getClTestDisByDiseaseID))
	sub.HandleFunc("/deleteClTestDis", middleware.AuthenticationMiddleware(s.deleteClTestDis))
	sub.HandleFunc("/addClTestDis", middleware.AuthenticationMiddleware(s.addClTestDis))
	sub.HandleFunc("/updateClTestDis", middleware.AuthenticationMiddleware(s.updateClTestDis))
}

func (s *clTestDisService) addClTestDis(w http.ResponseWriter, r *http.Request) {
	var clinicalTest models.ClinicalTestDisease

	err := json.NewDecoder(r.Body).Decode(&clinicalTest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errors := handlers.ValidateInputs(clinicalTest)
	if !isValid {
		for _, fieldError := range errors {
			http.Error(w, fieldError, http.StatusBadRequest)
			return
		}
	}

	compIDs := handlers.GetCompany(r)

	// TODO: check company ID if exists and if caller is related
	validCompanyID, validCompanyIDError := s.clinicalRepository.CheckCompany(compIDs, clinicalTest)
	if !validCompanyID {
		handlers.ProduceErrorResponse(validCompanyIDError, w, r)
		return
	}

	clinicalTest, err = s.clinicalRepository.AddClTestDis(clinicalTest)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}
	handlers.ProduceSuccessResponse("Clinical Test Registration - Successful", "", w, r)
}

func (s *clTestDisService) getClTestDisByDiseaseID(w http.ResponseWriter, r *http.Request) {
	var clinicalTest []models.ClinicalTestDisease

	// // Current account id
	// accountId := gcontext.Get(r, "id").(uint)

	params := mux.Vars(r)

	id := params["id"]
	if id == "" {
		handlers.ProduceErrorResponse("Please input all required fields.", w, r)
		return
	}

	// Convert string parameter to uint
	diseaseID, err := handlers.ConvertStrToUint(id)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	compIDs := handlers.GetCompany(r)

	// TODO: check company ID if exists and if caller is related
	validCompanyID, validCompanyIDError := s.clinicalRepository.CheckDiseaseCompanyClinical(compIDs, diseaseID)
	if !validCompanyID {
		handlers.ProduceErrorResponse(validCompanyIDError, w, r)
		return
	}

	clinicalTest, err = s.clinicalRepository.GetClTestDisByDisID(diseaseID)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	jsonRetrieved, err := json.Marshal(clinicalTest)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.ProduceSuccessResponse(string(jsonRetrieved), "", w, r)
}

func (s *clTestDisService) updateClTestDis(w http.ResponseWriter, r *http.Request) {
	var clinicalTest models.ClinicalTestDisease
	var response models.Response
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&clinicalTest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// isValid, errors := handlers.ValidateInputs(therapy)
	// if !isValid {
	// 	for _, fieldError := range errors {
	// 		http.Error(w, fieldError, http.StatusBadRequest)
	// 		return
	// 	}
	// }

	clinicalTest, err = s.clinicalRepository.UpdateClTestDis(clinicalTest)
	if err != nil {
		var msg string
		// if strings.Contains(err.Error(), "users_company_email_key") {
		// 	newerr = "user already exists!"
		// } else {
		// 	newerr = "Bad Request"
		// }
		response.Status = "error"
		response.Message = msg
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Fprintf(w, "Therapy Update - Successful")
}

func (s *clTestDisService) getClTestDis(w http.ResponseWriter, r *http.Request) {
	var clinicalTest models.ClinicalTestDisease
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

	intID, err := strconv.Atoi(id)

	clinicalTest, err = s.clinicalRepository.GetClTestDis(intID)
	if err != nil {
		response.Status = "error"
		response.Message = "Unknown Username or Password"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(clinicalTest)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = username
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)
}

func (s *clTestDisService) deleteClTestDis(w http.ResponseWriter, r *http.Request) {
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

	intID, err := strconv.Atoi(id)

	_, err = s.clinicalRepository.DeleteClTestDis(intID)
	if err != nil {
		response.Status = "error"
		response.Message = "Unknown Username or Password"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = "success"
	response.Message = username
	response.Response = ""
	json.NewEncoder(w).Encode(response)
}
