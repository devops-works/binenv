package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version of current binary
	Version string
	// BuildDate of current binary
	BuildDate string
)

// versionCmd lists installable versions 'as seen from cache
func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show binenv version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("binenv version %s (built %s)\n", Version, BuildDate)
		},
	}

	return cmd
}
