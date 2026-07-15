package main

import (
	"log/slog"
	"net"
	"os"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/server"
)

func run() error {

	app := core.NewApp()

	err := app.Run(
		core.StageConfig,
		core.StageStorage,
		core.StageDB,
	)

	if err != nil {
		return err
	}

	defer app.Cleanup()

	server := server.NewServer(app)
	addr := net.JoinHostPort(app.Config.Host, app.Config.Port)
	if err := server.Run(addr); err != nil {
		slog.Error("Something went wrong while trying to run the server")
		return err
	}

	return nil

}

func main() {
	if err := run(); err != nil {
		slog.Error("gochive failed", "error", err)
		os.Exit(1)
	}
}
