package database

import (
	"database/sql"
	"fmt"
	"log"

	config "rehab/config/database"
)

func Connection() {

	// Read configuration
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(cfg.Postgres.Host)

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")
	// Instantiates the database
	// postgres, err := database.NewPostgres(cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
