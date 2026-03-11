package postgres

import (
	"database/sql"
	"log"

	"github.com/UGRORF/price-tracker/internal/config"
	_ "github.com/lib/pq"
)

func NewDB(cfg *config.DatabaseConfig) (*sql.DB, error) {
	connString := cfg.ConnectionString()

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}
