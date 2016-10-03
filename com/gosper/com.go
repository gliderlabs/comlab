package gosper

import (
	"github.com/gliderlabs/gosper/pkg/com"
	"github.com/spf13/cobra"
)

func init() {
	com.Register("gosper", &struct{}{})
}

type CommandProvider interface {
	RegisterCommands(root *cobra.Command)
}

func CommandProviders() []CommandProvider {
	var providers []CommandProvider
	for _, com := range com.Enabled(new(CommandProvider), nil) {
		providers = append(providers, com.(CommandProvider))
	}
	return providers
}
