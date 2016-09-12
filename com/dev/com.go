package dev

import (
	"os"
	"path/filepath"

	"github.com/gliderlabs/pkg/com"
	"github.com/spf13/cobra"
)

func init() {
	com.Register("dev", &Component{})
}

type Component struct{}

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
