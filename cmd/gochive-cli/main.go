package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
)

func Execute() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Something went wrong while loading the config", "config", err)
		return err
	}

	app := &core.App{Config: cfg}

	cmd := rootCmd(app)
	return cmd.Execute()
}

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
