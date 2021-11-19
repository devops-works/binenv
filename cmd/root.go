package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// RootCmd returns the root cobra command
func RootCmd() *cobra.Command {
	a, err := app.New()
	if err != nil {
		fmt.Printf("got error %v\n", err)
		panic(err)
	}

	rootCmd := &cobra.Command{
		Use:   "binenv",
		Short: "Install binary distributions easily",
		Long: `binenv lets you install binary-distributed applications
(e.g. terraform, kubectl, ...) easily and switch between any version.
		
If your directory has a '.binenv.lock', proper versions will always be
selected.`,
		SilenceUsage: true,
	}

	if !strings.HasSuffix(os.Args[0], "binenv") {
		a.Execute(os.Args)
	}

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose operation")
	rootCmd.PersistentFlags().BoolP("system", "s", false, "Use system locations for directories")
	rootCmd.PersistentFlags().StringP("bindir", "B", app.GetDefaultBinDir(), "Distribution binaries path")
	rootCmd.PersistentFlags().StringP("distdir", "D", app.GetDefaultDistDir(), "Distributions path")
	rootCmd.PersistentFlags().StringP("cachedir", "C", app.GetDefaultCacheDir(), "Cache path")

	rootCmd.AddCommand(
		completionCmd(),
		installCmd(a),
		localCmd(a),
		searchCmd(a),
		uninstallCmd(a),
		updateCmd(a),
		versionCmd(),
		versionsCmd(a),
	)

	return rootCmd
}
