package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeLink)
	return r
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func ListenRoute(r *mux.Router) {

	r.HandleFunc("/", homeLink)
	log.Fatal(http.ListenAndServe(":8080", r))

}
