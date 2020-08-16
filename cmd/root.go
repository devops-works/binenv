package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/devopsworks/tools/binenv/internal/app"
)

// RootCmd returns the root cobra command
func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "binenv",
		Short: "Install binary distributions easily",
		Long: `binenv lets you install binary-distributed applications
(e.g. terraform, kubectl, ...) easily and switch between any version.
		
If your directory has a '.binenv.lock', proper versions will always be
selected.`,
	}

	var bindir string
	// var verbose
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose operation")
	rootCmd.Flags().StringVarP(&bindir, "bindir", "b", app.GetDefaultBinDir(), "Binaries directory")

	// Run as a shim for installed distributions
	if !strings.HasSuffix(os.Args[0], "binenv") {
		// fmt.Printf("called as %s\n", os.Args[0])
		// fmt.Printf("bindir is %s\n", bindir)

		a, err := app.New(
			app.WithDiscard(),
			app.WithBinDir(bindir),
		)
		if err != nil {
			fmt.Printf("got error %v\n", err)
			panic(err)
		}

		fmt.Printf("calling execute for %s\n", os.Args[0])
		a.Execute(os.Args)
	}

	rootCmd.AddCommand(
		completionCmd(),
		installCmd(),
		localCmd(),
		updateCmd(),
		versionsCmd(),
	)

	return rootCmd
}
