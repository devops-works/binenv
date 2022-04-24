package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// versionsCmd lists installable versions as seen from cache
func versionsCmd(a *app.App) *cobra.Command {
	var freeze bool

	cmd := &cobra.Command{
		Use:   "versions [distribution...] [--freeze]",
		Short: "List installable versions",
		Long: `List all installable versions for a distribution.
If the distribution is not specified, lists all available version for all distributions.

Version currenyly in used has a '*' next to it.
Versions installed locally have a '+'.

The --freeze (-f) argument will output a list of currently selected distribution versions on stdout.

Use 'binenv update' to update the list of available versions.`,
		Run: func(cmd *cobra.Command, args []string) {
			freeze, _ := cmd.Flags().GetBool("freeze")
			a.Versions(freeze, args...)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return a.GetPackagesListWithPrefix(toComplete), cobra.ShellCompDirectiveNoFileComp
		},
	}

	cmd.Flags().BoolVarP(&freeze, "freeze", "f", false, "Write a .binenv.lock file to stdout containing currently selected versions")

	return cmd
}
