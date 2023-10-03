package main

import (
	"net/http"
	"rehab/api/accounts"
	"rehab/api/companies"
	"rehab/api/general"
	"rehab/api/patients"
	"rehab/api/ph_therapies"
	"rehab/internal/pkg/handlers"
)

func main() {

	// dbConnection := config.ConnectDB()
	router := handlers.NewRouter()
	router.Use(corsMiddleware)

	// queries := database.New(postgres.DB)
	// authorService := accounts.NewService()

	// authorService.RegisterHandlers()

	// queries := database.New(postgres.DB)
	patientService := patients.NewService()
	patientService.RegisterHandlers(router)

	patientDetailsService := patients.NewDetailsService()
	patientDetailsService.RegisterDetailHandlers(router)

	therapyService := patients.NewTherapyService()
	therapyService.RegisterHandlers(router)

	phTherapyService := ph_therapies.NewPhTherapiesService()
	phTherapyService.RegisterPhTherapiesHandlers(router)

	companyService := companies.NewService()
	companyService.RegisterHandlers(router)

	accountService := accounts.NewService()
	accountService.RegisterHandlers(router)

	generalService := general.NewService()
	generalService.RegisterHandlers(router)

	diseaseService := patients.NewDiseaseService()
	diseaseService.RegisterHandlers(router)

	clTestDysService := patients.NewClTestDisService()
	clTestDysService.RegisterHandlers(router)

	handlers.ListenRoute(router)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		w.Header().Set("content-type", "application/json;charset=UTF-8")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
