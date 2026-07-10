package core

import (
	"database/sql"
	"embed"
	"path"

	"github.com/0xMoonrise/gochive/internal/database"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed db/migrations/*
var migrationsFS embed.FS

func migrations(dialect string, db *sql.DB) error {

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	goose.SetBaseFS(migrationsFS)
	if err := goose.Up(db, "db/migrations"); err != nil {
		return err
	}

	return nil
}

func BootDatabase(app *App) (func() error, error) {

	db, err := sql.Open("sqlite3", path.Join(app.Config.Data, "gochive.db"))
	if err != nil {
		return nil, err
	}

	if err := migrations("sqlite3", db); err != nil {
		return nil, err
	}

	app.DB = Database{
		DB:      db,
		Queries: database.New(db),
	}
	return db.Close, nil
}
