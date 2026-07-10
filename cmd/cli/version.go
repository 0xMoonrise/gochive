package main

import (
	"fmt"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/spf13/cobra"
)

func version() *cobra.Command {
	return &cobra.Command{
		Use:                   "version",
		Short:                 "Prints version",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.VERSION)
		},
	}
}
