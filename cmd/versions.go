package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/devopsworks/tools/binenv/internal/app"
)

// versionsCmd lists installable versions 'as seen from cache
func versionsCmd() *cobra.Command {
	app := app.New()
	cmd := &cobra.Command{
		Use:   "versions",
		Short: "List installable versions",
		Long:  `List all installable versions for this application.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return app.Versions(args[0])
			}
			return app.Versions("")
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return app.GetPackagesListWithPrefix(toComplete), cobra.ShellCompDirectiveNoFileComp
		},
	}

	// cmd.Flags().IntVarP(&a.Params.MinLength, "min-length", "m", 16, "Specify minimum password length, must not be less than 8")
	return cmd
}
