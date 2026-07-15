package main

import (
	"os"

	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gochive-cli",
		Short: "gochive-cli to interact with the Gochive backend",
	}

	root.AddCommand(version())
	root.AddCommand(status())
	root.AddCommand(newBackupCmd())
	root.AddCommand(newRestoreCmd())
	root.AddCommand(newGenerateThumbnail())
	root.AddCommand(newUploadArchive())

	return root
}

func main() {
	cmd := rootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
