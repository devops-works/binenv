package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// versionsCmd lists installable versions as seen from cache
func versionsCmd(a *app.App) *cobra.Command {
	var bindir string

	cmd := &cobra.Command{
		Use:   "versions [distribution...]",
		Short: "List installable versions",
		Long: `List all installable versions for a distribution.
If the distribution is not specified, lists all available version for all distributions.

Version currenyly in used has a '*' next to it.
Versions installed locally have a '+'.

Use 'binenv update' to update the list of available versions.`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			a.SetVerbose(verbose)
			a.SetBinDir(bindir)
			a.Versions(args...)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return a.GetPackagesListWithPrefix(toComplete), cobra.ShellCompDirectiveNoFileComp
		},
	}

	cmd.Flags().StringVarP(&bindir, "bindir", "b", app.GetDefaultBinDir(), "Binaries directory")

	return cmd
}
