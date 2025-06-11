package main

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func Init() error {

	slog.Info("Reading .env file...")
	err := godotenv.Load()

	if err != nil {
		slog.Warn("Something went wrong while trying to read the .env file")
		return err
	}
	
	slog.Info("Success! .env has been loaded")

	return nil
}
