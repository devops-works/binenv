package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// upgradeCmd upgrade all installed distributions
func upgradeCmd(a *app.App) *cobra.Command {
	var bindir string
	var ignoreInstallErrors bool

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade all installed distributions",
		Long:  `Upgrade all installed distributions to the last version available on cache.`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")

			a.SetVerbose(verbose)
			a.SetBinDir(bindir)
			a.Upgrade(ignoreInstallErrors)
		},
	}

	cmd.Flags().StringVarP(&bindir, "bindir", "b", app.GetDefaultBinDir(), "Binaries directory")
	cmd.Flags().BoolVarP(&ignoreInstallErrors, "ignore-install-errors", "i", true, "Ignore install errors during upgrade")

	return cmd
}
