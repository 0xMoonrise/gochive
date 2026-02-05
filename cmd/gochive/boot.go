package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/0xMoonrise/gochive/internal/app"
	"github.com/0xMoonrise/gochive/internal/database"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

const maxRetries = 3

func migrations(dialect string, db *sql.DB) error {

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	cwd, _ := os.Getwd()

	if err := goose.Up(db, filepath.Join(cwd, "db", "migrations")); err != nil {
		return err
	}

	return nil
}

func newPG() (db *sql.DB, err error) {

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
	for i := range maxRetries {
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
			"max", maxRetries,
			"error", err,
		)

		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("database unavailable after %d attempts: %w", maxRetries, err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

	slog.Info("Database connected successfully")

	return db, nil
}

func newSQLITE() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "/opt/gochive/gochive.db")
	return
}

func bootDatabase(app *app.App) (func() error, error) {
	db, err := newSQLITE()

	if err != nil {
		return nil, err
	}

	schema, err := os.ReadFile("db/sql/schema.sql")
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(string(schema)); err != nil {
		return nil, err
	}

	database := database.New(db)

	// if err := migrations("sqlite", db); err != nil {
	// 	return nil, err
	// }

	app.Db = database

	return db.Close, nil
}
