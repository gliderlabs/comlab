package dev

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func (c *Component) RegisterCommands(root *cobra.Command) {
	root.AddCommand(&cobra.Command{
		Use:   "dev",
		Short: "Dev runner",
		Run: func(cmd *cobra.Command, args []string) {
			cwd, _ := os.Getwd()
			cmdName := filepath.Base(cwd)
			Run(cmdName)
		},
	})
}
