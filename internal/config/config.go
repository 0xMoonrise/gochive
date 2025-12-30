package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"
)

const MAX_RETRIES = 3

func InitDataBase() (*sql.DB, error) {
	databaseUri := os.Getenv("DB_URI")

	var (
		conn *sql.DB
		err error
	)
	
	if databaseUri == "" {
		return nil, errors.New("DB_URI is not set")
	}

	slog.Info("Initializing database connection")

	for i := range MAX_RETRIES {
		conn, err = sql.Open("postgres", databaseUri)
		if err == nil {
			err = conn.Ping()
		}

		if err == nil {
			break
		}

		slog.Warn(
			"Cannot connect to database, retrying",
			"attempt", i,
			"max", MAX_RETRIES,
			"error", err,
		)

		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("database unavailable after %d attempts: %w", MAX_RETRIES, err)
	}

	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(2)

	slog.Info("Database connected successfully")

	return conn, nil
}
