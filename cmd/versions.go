package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// versionsCmd lists installable versions as seen from cache
func versionsCmd() *cobra.Command {
	app, err := app.New()
	if err != nil {
		panic(err)
	}
	cmd := &cobra.Command{
		Use:   "versions [distribution...]",
		Short: "List installable versions",
		Long: `List all installable versions for a distribution.
If the distribution is not specified, lists all available version for all distributions.

Version currenyly in used has a '*' next to it.
Versions installed locally have a '+'.

Use 'binenv update' to update the list of available versions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Versions(args...)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return app.GetPackagesListWithPrefix(toComplete), cobra.ShellCompDirectiveNoFileComp
		},
	}

	// cmd.Flags().IntVarP(&a.Params.MinLength, "min-length", "m", 16, "Specify minimum password length, must not be less than 8")
	return cmd
}
