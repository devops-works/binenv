package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// RootCmd returns the root cobra command
func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "binenv",
		Short: "Install binary distributions easily",
		Long: `binenv lets you install binary-distributed applications
		(e.g. terraform, kubectl, ...) easily and switch between any version.
		
		If your directory has a '.binenv.lock', proper evrsions will always be
		selected.`,
	}

	rootCmd.AddCommand(
		completionCmd(),
		installCmd(),
		localCmd(),
		updateCmd(),
		versionsCmd(),
	)

	d, err := homedir.Dir()
	if err != nil {
		d = "~"
	}

	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose operation")
	rootCmd.Flags().StringP("bindir", "b", d+"/.binenv/", "Binaries directory")

	return rootCmd
}
