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

	app := &core.App{}
	log.Println(config.MODE)
	client, err := setMode(config.MODE)
	if err != nil {
		slog.Error("Something went wrong while trying to create a storage client",
			"store client",
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
	if err := cli(app); err != nil {
		return nil
	}
	server := server.NewServer(app)
	addr := net.JoinHostPort(config.HOST, config.PORT)
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
		log.Fatal(err)
	}
}
