package main

import (
	"log/slog"
	"net"
	"net/http"
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

	addr := net.JoinHostPort(app.Config.Host, app.Config.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("failed to bind port", "addr", addr, "error", err)
		os.Exit(1)
	}

	slog.Info("server ready", "addr", addr)

	server := server.NewServer(app)
	if err := http.Serve(ln, server); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}

	return nil

}

func main() {
	if err := run(); err != nil {
		slog.Error("gochive failed", "error", err)
		os.Exit(1)
	}
}
