package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// upgradeCmd upgrade all installed distributions
func upgradeCmd(a *app.App) *cobra.Command {
	var ignoreInstallErrors bool

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade all installed distributions",
		Long:  `Upgrade all installed distributions to the last version available on cache.`,
		Run: func(cmd *cobra.Command, args []string) {
			a.Upgrade(ignoreInstallErrors)
		},
	}

	cmd.Flags().BoolVarP(&ignoreInstallErrors, "ignore-install-errors", "i", true, "Ignore install errors during upgrade")

	return cmd
}
