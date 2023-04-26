package patients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rehab-backend/internal/middleware"
	"github.com/rehab-backend/internal/pkg/handlers"
	"github.com/rehab-backend/internal/pkg/models"
	"github.com/rehab-backend/internal/repository"
)

type clTestDysService struct {
	clinicalRepository repository.ClTestDysRepository
}

func NewClTestDysService() *clTestDysService {
	return &clTestDysService{
		clinicalRepository: repository.NewClTestDysService(),
	}
}

func (s *clTestDysService) RegisterHandlers(route *mux.Router) {
	s.Handle(route)
}

func (s *clTestDysService) Handle(route *mux.Router) {
	sub := route.PathPrefix("/clinicalTestDysfunction").Subrouter()

	sub.HandleFunc("/getClTestDys", middleware.AuthenticationMiddleware(s.getClTestDys))
	sub.HandleFunc("/deleteClTestDys", middleware.AuthenticationMiddleware(s.deleteClTestDys))
	sub.HandleFunc("/addClTestDys", middleware.AuthenticationMiddleware(s.addClTestDys))
	sub.HandleFunc("/updateClTestDys", middleware.AuthenticationMiddleware(s.updateClTestDys))
}

func (s *clTestDysService) addClTestDys(w http.ResponseWriter, r *http.Request) {
	var clinicalTest models.ClinicalTestDysfunction

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

	compIDs := gcontext.Get(r, "compIDs").([]uint)

	// TODO: check company ID if exists and if caller is related
	validCompanyID, validCompanyIDError := s.clinicalRepository.CheckCompany(compIDs, clinicalTest)
	if !validCompanyID {
		handlers.ProduceErrorResponse(validCompanyIDError, w, r)
		return
	}

	clinicalTest, err = s.clinicalRepository.AddClTestDys(clinicalTest)
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
	handlers.ProduceSuccessResponse("Clinical Test Registration - Successful", w, r)
}

func (s *clTestDysService) updateClTestDys(w http.ResponseWriter, r *http.Request) {
	var clinicalTest models.ClinicalTestDysfunction
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

	clinicalTest, err = s.clinicalRepository.UpdateClTestDys(clinicalTest)
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

func (s *clTestDysService) getClTestDys(w http.ResponseWriter, r *http.Request) {
	var clinicalTest models.ClinicalTestDysfunction
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

	clinicalTest, err = s.clinicalRepository.GetClTestDys(intID)
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

func (s *clTestDysService) deleteClTestDys(w http.ResponseWriter, r *http.Request) {
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

	_, err = s.clinicalRepository.DeleteClTestDys(intID)
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
