package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/devopsworks/tools/binenv/internal/app"
)

// localCmd represents the local command
func installCmd() *cobra.Command {
	app := app.New()
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install a version for the package",
		Long:  `Better description here`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Install()
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0:
				// complete application name
				return app.GetPackagesListWithPrefix(toComplete), cobra.ShellCompDirectiveNoFileComp
			case 1:
				// complete application version
				return app.GetVersionsFromCacheFor(args[0]), cobra.ShellCompDirectiveNoFileComp
			default:
				// huh ?
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
		},
	}

	// cmd.Flags().IntVarP(&a.Params.MinLength, "min-length", "m", 16, "Specify minimum password length, must not be less than 8")
	return cmd
}
