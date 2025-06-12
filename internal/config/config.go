package config

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
)

func LoadConfig() (*sql.DB, error) {

	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")

	slog.Info("Loading db config")

	db_uri := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", db_user, db_pass, db_host, db_port, db_name)
	db, err := sql.Open("postgres", db_uri)

	// Should I shut down the web server if it can't connect to the database?
	if err != nil {
		return nil, err
	}

	return db, nil
}
