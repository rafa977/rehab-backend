package patients

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Service struct {
	// queries *database.Queries
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) RegisterHandlers(route *mux.Router) {

	Handle(route)

}

func Handle(route *mux.Router) {
	route.HandleFunc("/registerUser", homeLink)
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home222!")
}
