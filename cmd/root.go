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
	a, err := app.New(
	// app.WithDiscard(),
	// app.WithBinDir(bindir),
	)
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

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose operation (env: BINENV_VERBOSE)")
	rootCmd.PersistentFlags().BoolP("global", "g", false, "Global mode (env: BINENV_GLOBAL)")
	rootCmd.PersistentFlags().StringP("bindir", "b", "~/.binenv", "Binaries directory")

	if !strings.HasSuffix(os.Args[0], "binenv") {
		verbose := truthify(os.Getenv("BINENV_VERBOSE"))
		global := truthify(os.Getenv("BINENV_GLOBAL"))

		a.SetGlobal(global)
		a.SetVerbose(verbose)
		
		a.Execute(os.Args)
	}

	rootCmd.AddCommand(
		completionCmd(),
		expandCmd(a),
		installCmd(a),
		localCmd(a),
		searchCmd(a),
		uninstallCmd(a),
		updateCmd(a),
		versionCmd(),
		versionsCmd(a),
		upgradeCmd(a),
	)

	return rootCmd
}

func truthify(s string) bool {
	s = strings.ToLower(s)
	// trueness suggestions courtesy of Github co-pilot
	return s == "true" || s == "1" || s == "yes" || s == "y" || s == "on" || s == "enable" || s == "enabled" || s == "active"
}
