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

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type therapyService struct {
	therapyRepository repository.TherapyRepository
}

func NewTherapyService() *therapyService {

	return &therapyService{
		therapyRepository: repository.NewTherapyService(),
	}
}

func (s *therapyService) RegisterHandlers(route *mux.Router) {

	s.Handle(route)

}

func (s *therapyService) Handle(route *mux.Router) {

	sub := route.PathPrefix("/therapy").Subrouter()

	sub.HandleFunc("/getTherapy", middleware.AuthenticationMiddleware(s.getTherapy))
	sub.HandleFunc("/deleteTherapy", middleware.AuthenticationMiddleware(s.deleteTherapy))
	sub.HandleFunc("/addTherapy", middleware.AuthenticationMiddleware(s.addTherapy))
	sub.HandleFunc("/updateTherapy", middleware.AuthenticationMiddleware(s.addTherapy))
}

func (s *therapyService) addTherapy(w http.ResponseWriter, r *http.Request) {
	var therapy models.Therapy
	var response models.Response

	err := json.NewDecoder(r.Body).Decode(&therapy)
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

	therapy, err = s.therapyRepository.AddTherapy(therapy)
	if err != nil {
		var newerr string
		// if strings.Contains(err.Error(), "users_company_email_key") {
		// 	newerr = "user already exists!"
		// } else {
		// 	newerr = "Bad Request"
		// }
		response.Status = "error"
		response.Message = newerr
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Fprintf(w, "Therapy Registration - Successful")
}

func (s *therapyService) updateTherapy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var therapy models.Therapy
	var response models.Response
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&therapy)
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

	therapy, err = s.therapyRepository.UpdateTherapy(therapy)
	if err != nil {
		var newerr string
		// if strings.Contains(err.Error(), "users_company_email_key") {
		// 	newerr = "user already exists!"
		// } else {
		// 	newerr = "Bad Request"
		// }
		response.Status = "error"
		response.Message = newerr
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Fprintf(w, "Therapy Update - Successful")
}

func (s *therapyService) getTherapy(w http.ResponseWriter, r *http.Request) {
	var therapy models.Therapy
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

	therapy, err = s.therapyRepository.GetTherapy(intID)
	if err != nil {
		response.Status = "error"
		response.Message = "Unknown Username or Password"
		response.Response = ""
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	jsonRetrievedAccount, err := json.Marshal(therapy)
	if err != nil {
		fmt.Println(err)
		return
	}

	response.Status = "success"
	response.Message = username
	response.Response = string(jsonRetrievedAccount)
	json.NewEncoder(w).Encode(response)
}

func (s *therapyService) deleteTherapy(w http.ResponseWriter, r *http.Request) {
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

	_, err = s.therapyRepository.DeleteTherapy(intID)
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
