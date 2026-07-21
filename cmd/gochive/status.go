package main

import (
	"fmt"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/spf13/cobra"
)

func status() *cobra.Command {
	app := core.NewApp()
	return &cobra.Command{
		Use:                   "status",
		Short:                 "Show Gochive configuration and status summary",
		DisableFlagsInUseLine: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(core.StageConfig)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Port: ", app.Config.Port)
			fmt.Println("Host: ", app.Config.Host)
			fmt.Println("Data: ", app.Config.Data)
			fmt.Printf("Mode: %v -> %v\n", app.Config.Mode, app.ModeToString())
			switch app.Config.Mode {
			case config.FS:
				fmt.Printf("File system root: %s\n", app.Config.FS.Root)
			case config.S3:
				fmt.Printf("Bucket: %s\n", app.Config.S3.Bucket)
				fmt.Printf("S3 Endpoint: %s\n", app.Config.S3.S3Endpoint)
				fmt.Printf("Region: %s\n", app.Config.S3.Region)
			}
		},
	}
}
