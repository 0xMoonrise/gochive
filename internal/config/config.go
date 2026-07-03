package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var MODE int

var (
	// http server
	PORT,
	HOST,
	// db config
	DB_USER,
	DB_PASS,
	DB_HOST,
	DB_NAME,
	DB_PORT,
	// S3 client
	BUCKET,
	ACCESS_KEY,
	SECRET_KEY,
	S3_ENDPOINT,
	REGION string
)

func mustGetEnv(arg string) string {
	v := os.Getenv(arg)
	if v == "" {
		slog.Error("missing required env var", "name", arg)
		os.Exit(1)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func init() {

	var err error
	if err = godotenv.Load(); err != nil {
		slog.Info(".env not loaded")
	}

	MODE, err = strconv.Atoi(os.Getenv("MODE"))
	if err != nil {
		slog.Error("Cannot be possible to conver Mode to string")
		os.Exit(1)
	}

	PORT = getEnv("PORT", "8080")
	HOST = getEnv("HOST", LOCAL)

	DB_HOST = getEnv("DB_HOST", LOCAL)
	DB_NAME = getEnv("DB_NAME", "gochive")
	DB_PORT = getEnv("DB_PORT", "5432")

	if MODE == S3 {
		BUCKET = mustGetEnv("BUCKET")
		ACCESS_KEY = mustGetEnv("ACCESS_KEY")
		SECRET_KEY = mustGetEnv("SECRET_KEY")

		REGION = getEnv("REGION", "us-east-1")
		S3_ENDPOINT = getEnv("S3_ENDPOINT", LOCAL)
	}
}
