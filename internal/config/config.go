package config

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"time"
)

func InitDataBase() (*sql.DB, error) {
	databaseUri := os.Getenv("DB_URI")

	if databaseUri == "" {
		return nil, errors.New("DB_URI is not set")
	}
	slog.Info("Initializing database connection")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := sql.Open("postgres", databaseUri)
	if err != nil {
		return nil, err
	}

	if err := conn.PingContext(ctx); err != nil {
		return nil, err
	}
	
	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(2)
	
	slog.Info("Database connected successfully")

	return conn, nil
}
