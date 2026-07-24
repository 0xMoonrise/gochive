package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/server"
	"github.com/spf13/cobra"
)

func newServerCmd() *cobra.Command {
	app := core.NewApp()
	return &cobra.Command{
		Use:                   "serve",
		Short:                 "start http server",
		DisableFlagsInUseLine: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(
				core.StageConfig,
				core.StageStorage,
				core.StageDB,
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {

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
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return app.Cleanup()
		},
	}
}
