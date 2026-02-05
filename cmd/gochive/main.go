package main

import (
	"log/slog"
	"net"
	"os"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/server"
	"github.com/joho/godotenv"
)

func run() error {

	if err := godotenv.Load(); err != nil {
		slog.Info(".env not loaded")
	}

	app := &core.App{}

	// app.Storage, err = application.NewS3Client()
	client, err := core.NewfsClient()
	if err != nil {
		slog.Error("Something went wrong while trying to create a storage client",
			"error",
			err,
		)
		return err
	}

	app.Storage = client
	closeDB, err := bootDatabase(app)
	if err != nil {
		slog.Error("Something went wrong while trying booting the database",
			"error",
			err,
		)
		return err
	}
	defer closeDB()

	server := server.NewServer(app)
	addr := net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT"))
	if err := server.Run(addr); err != nil {
		slog.Error("Something went wrong while trying to run the server",
			"error",
			err,
		)
		return err
	}

	return nil

}

func main() {

	if err := run(); err != nil {
		slog.Error("fatal:", "error", err)
		os.Exit(1)
	}

}
