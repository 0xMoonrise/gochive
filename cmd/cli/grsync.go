package main

import (
	"fmt"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/spf13/cobra"
	"github.com/zloylos/grsync"
)

func newBackupCmd(app *core.App) *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Backup gochive data to the external backup location",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRsync(app.Config.Data+"/", app.Config.Backup+"/")
		},
	}
}

func newRestoreCmd(app *core.App) *cobra.Command {
	return &cobra.Command{
		Use:   "restore",
		Short: "Restore gochive data from the external backup location",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRsync(app.Config.Backup+"/", app.Config.Data+"/")
		},
	}
}

func runRsync(src, dst string) error {
	task := grsync.NewTask(src, dst, grsync.RsyncOptions{
		Archive: true,
		Update:  true,
		Verbose: true,
	})

	if err := task.Run(); err != nil {
		return err
	}

	fmt.Println(task.Log().Stdout)
	return nil
}
