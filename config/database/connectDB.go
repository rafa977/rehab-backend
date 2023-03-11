package config

import (
	"fmt"
	"log"

	"github.com/rehab-backend/internal/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {

	var err error

	// Read configuration
	cfg, err := Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&models.Account{}, &models.Patient{}, &models.Therapy{})

	// psqlconn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)
	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	// sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(psqlconn)))

	// db := bun.NewDB(sqldb, pgdialect.New())

	// errPing := db.Ping()
	// fmt.Println(errPing)

	return db
}
