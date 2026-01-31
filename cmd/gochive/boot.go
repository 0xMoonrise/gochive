package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/0xMoonrise/gochive/internal/app"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/pressly/goose/v3"
)

const MAX_RETRIES = 3

func schemas(q *database.Queries) error {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if err := q.CreateSchema(ctx); err != nil {
		return err
	}

	if err := q.CreateArchiveTable(ctx); err != nil {
		return err
	}

	return nil
}

func migrations(db *sql.DB) error {

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	cwd, _ := os.Getwd()

	if err := goose.Up(db, filepath.Join(cwd, "db", "migrations")); err != nil {
		return err
	}

	return nil
}

func connectDB() (db *sql.DB, err error) {

	u := &url.URL{
		Scheme: "postgresql",
		User: url.UserPassword(
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
		),
		Host: os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Path: os.Getenv("DB_NAME"),
	}

	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()

	slog.Info("Initializing database connection")
	for i := range MAX_RETRIES {
		db, err = sql.Open("postgres", u.String())
		if err == nil {
			err = db.Ping()
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

		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("database unavailable after %d attempts: %w", MAX_RETRIES, err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

	slog.Info("Database connected successfully")

	return db, nil
}

func bootDatabase(app *app.App) (func() error, error) {
	db, err := connectDB()

	if err != nil {
		return nil, err
	}

	database := database.New(db)

	if err := schemas(database); err != nil {
		return nil, err
	}

	if err := migrations(db); err != nil {
		return nil, err
	}

	app.Db = database

	return db.Close, nil
}
