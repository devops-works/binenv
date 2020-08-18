package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
func uninstallCmd() *cobra.Command {
	var bindir string

	a, err := app.New()
	if err != nil {
		panic(err)
	}
	cmd := &cobra.Command{
		Use:   "uninstall <distribution> [<version> [<distribution> <version>]]",
		Short: "Uninstall a distribution or a specific version for the distribution",
		Long: `When a version is not specified, only a single distribution is accepted.
In this case, all versions for the specified distribution will be removed (a confirmation is asked).

Multiple distribution / version pairs can be specified.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			a.SetBinDir(bindir)
			return a.Uninstall(args...)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) % 2 {
			case 0:
				// complete application name
				return a.GetPackagesListWithPrefix(toComplete), cobra.ShellCompDirectiveNoFileComp
			case 1:
				// complete application version
				return a.GetInstalledVersionsFor(args[len(args)-1]), cobra.ShellCompDirectiveNoFileComp
			default:
				// huh ?
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
		},
	}

	cmd.Flags().StringVarP(&bindir, "bindir", "b", app.GetDefaultBinDir(), "Binaries directory")

	// cmd.Flags().IntVarP(&a.Params.MinLength, "min-length", "m", 16, "Specify minimum password length, must not be less than 8")
	return cmd
}
