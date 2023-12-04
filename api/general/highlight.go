package general

import (
	"encoding/json"
	"net/http"
	"rehab/internal/middleware"
	"rehab/internal/pkg/handlers"
	"rehab/internal/pkg/models"
	"rehab/internal/repository"

	"github.com/gorilla/mux"
)

type highlightService struct {
	highlightRepository repository.HighlightRepository
	response            models.Response
}

func NewHighlightService() *highlightService {
	return &highlightService{
		highlightRepository: repository.NewHighlightService(),
	}
}

func (s *highlightService) RegisterHandlers(route *mux.Router) {
	s.Handle(route)
}

func (s *highlightService) Handle(route *mux.Router) {

	sub := route.PathPrefix("/highlight").Subrouter()

	sub.HandleFunc("/addHighlight", middleware.AuthenticationMiddleware(s.addHighlight))

}

func (s *highlightService) addHighlight(w http.ResponseWriter, r *http.Request) {
	var highlight models.Highlight

	err := json.NewDecoder(r.Body).Decode(&highlight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: check company ID if exists and if caller is related
	compIDs := handlers.GetCompany(r)
	// userID := handlers.GetAccount(r)

	if len(compIDs) == 0 {
		handlers.ProduceErrorResponse("Please register your company", w, r)
		return
	}

	err = s.highlightRepository.AddHighlightDetails(highlight)
	if err != nil {
		handlers.ProduceErrorResponse(err.Error(), w, r)
		return
	}

	handlers.ProduceSuccessResponse("Highlight Added - Successful", "", w, r)
}
