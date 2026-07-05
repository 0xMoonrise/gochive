package main

import (
	"database/sql"
	"os"
	"path"
	"path/filepath"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
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
	path := filepath.Join(cwd, "db", "migrations")
	if err := goose.Up(db, path); err != nil {
		return err
	}

	return nil
}

func bootDatabase(app *core.App) (func() error, error) {

	db, err := sql.Open("sqlite3", path.Join(config.ROOT, "gochive.db"))
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

	app.DB = core.Database{
		DB:      db,
		Queries: database.New(db),
	}
	return db.Close, nil
}
