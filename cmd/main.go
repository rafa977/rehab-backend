package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/rehab-backend/api/accounts"
	patients "github.com/rehab-backend/api/patients"
	config "github.com/rehab-backend/config/database"
	"github.com/rehab-backend/internal/pkg/handlers"
)

func main() {

	// Read configuration
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(cfg.Postgres.Host)

	// psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	// // open database
	// db, err := sql.Open("postgres", psqlconn)
	// CheckError(err)

	// // close database
	// defer db.Close()

	// // check db
	// err = db.Ping()
	// CheckError(err)

	fmt.Println("Connected!")
	// Instantiates the database
	// postgres, err := database.NewPostgres(cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	router := handlers.NewRouter()

	// queries := database.New(postgres.DB)
	// authorService := accounts.NewService()

	// authorService.RegisterHandlers()

	// queries := database.New(postgres.DB)
	patientService := patients.NewService()
	patientService.RegisterHandlers(router)

	accountService := accounts.NewService()
	accountService.RegisterHandlers(router)

	handlers.ListenRoute(router)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
