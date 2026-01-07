package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/pressly/goose/v3"
)

func loadDB(db *sql.DB, q *database.Queries) error {
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


func InitDataBase() (*sql.DB, error) {

  databaseUri := fmt.Sprintf(
    "postgresql://%s:%s@%s:%s/%s?sslmode=disable",
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASS"),
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_NAME"),
  )

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

    time.Sleep(2 * time.Second)
  }

  if err != nil {
    return nil, fmt.Errorf("database unavailable after %d attempts: %w", MAX_RETRIES, err)
  }

  conn.SetMaxOpenConns(5)
  conn.SetMaxIdleConns(2)

  slog.Info("Database connected successfully")

  return conn, nil
}
