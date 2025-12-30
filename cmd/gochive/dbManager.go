package main

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/pressly/goose/v3"
)


func dbManager(db *sql.DB, q *database.Queries) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	if err := q.CreateSchema(ctx); err != nil {
		return err
	}

	if err := q.CreateArchiveTable(ctx); err != nil {
		return err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	if err := goose.Up(db, filepath.Join(cwd, "db", "migrations")); err != nil {
		return err
	}

	return nil
}
