package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func ConnectDB() *bun.DB {

	var err error

	// Read configuration
	cfg, err := Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	psqlconn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)
	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(psqlconn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	errPing := db.Ping()
	fmt.Println(errPing)

	return db
}
