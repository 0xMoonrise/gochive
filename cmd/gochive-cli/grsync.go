package main

import (
	"fmt"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/spf13/cobra"
	"github.com/zloylos/grsync"
)

func newBackupCmd() *cobra.Command {
	app := core.NewApp()
	return &cobra.Command{
		Use:   "backup",
		Short: "Backup gochive data to the external backup location",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(core.StageConfig)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRsync(app.Config.Data+"/", app.Config.Backup+"/")
		},
	}
}

func newRestoreCmd() *cobra.Command {
	app := core.NewApp()
	return &cobra.Command{
		Use:   "restore",
		Short: "Restore gochive data from the external backup location",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(core.StageConfig)
		},
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
