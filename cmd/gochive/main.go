package main

import (
	"log"
	"log/slog"
	"net"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/server"
)

func run() error {

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Something went wrong while loading the config",
			"store client",
			err,
		)
		return err
	}

	app := &core.App{
		Config: cfg,
	}
	client, err := app.SetMode()
	if err != nil {
		slog.Error("Something went wrong while trying to create a storage client",
			"store client",
			err,
		)
		return err
	}

	app.Storage = client
	closeDB, err := core.BootDatabase(app)
	if err != nil {
		slog.Error("Something went wrong while trying booting the database",
			"database error",
			err,
		)
		return err
	}

	defer closeDB()

	server := server.NewServer(app)
	addr := net.JoinHostPort(app.Config.Host, app.Config.Port)
	if err := server.Run(addr); err != nil {
		slog.Error("Something went wrong while trying to run the server",
			"server error",
			err,
		)
		return err
	}

	return nil

}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
