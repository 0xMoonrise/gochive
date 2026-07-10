package main

import (
	"log/slog"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/spf13/cobra"
)

func rootCmd(app *core.App) *cobra.Command {
	var closeDB func() error

	root := &cobra.Command{
		Use:   "gochive",
		Short: "gochive-cli to interact with the Gochive backend",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "backup" || cmd.Name() == "restore" {
				return nil
			}

			client, err := app.SetMode()
			if err != nil {
				slog.Error("Something went wrong while trying to create a storage client", "store client", err)
				return err
			}
			app.Storage = client

			cdb, err := core.BootDatabase(app)
			if err != nil {
				slog.Error("Something went wrong while trying booting the database", "database", err)
				return err
			}
			closeDB = cdb
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if closeDB != nil {
				return closeDB()
			}
			return nil
		},
	}

	root.AddCommand(newBackupCmd(app))
	root.AddCommand(newRestoreCmd(app))
	root.AddCommand(version())
	root.AddCommand(newGenerateThumbnail(app))
	root.AddCommand(newUploadArchive(app))

	return root
}
