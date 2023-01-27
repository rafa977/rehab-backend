package main

import (
	_ "github.com/lib/pq"

	"github.com/rehab-backend/api/accounts"
	"github.com/rehab-backend/internal/pkg/handlers"
)

func main() {

	// dbConnection := config.ConnectDB()

	router := handlers.NewRouter()

	// queries := database.New(postgres.DB)
	// authorService := accounts.NewService()

	// authorService.RegisterHandlers()

	// queries := database.New(postgres.DB)
	// patientService := patients.NewService()
	// patientService.RegisterHandlers(router)

	accountService := accounts.NewService()
	accountService.RegisterHandlers(router)

	handlers.ListenRoute(router)
}
