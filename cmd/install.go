package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
func installCmd() *cobra.Command {
	// var bindir string

	a, err := app.New()
	if err != nil {
		panic(err)
	}
	cmd := &cobra.Command{
		Use:   "install <distribution> <version> [<distribution> <version>]",
		Short: "Install a version for the package",
		Long:  `This command will install one or several distributions with the specified versions. `,
		RunE: func(cmd *cobra.Command, args []string) error {
			// a.SetBinDir(bindir)
			return a.Install(args...)
		},
		Args: cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) % 2 {
			case 0:
				// complete application name
				return a.GetPackagesListWithPrefix(toComplete), cobra.ShellCompDirectiveNoFileComp
			case 1:
				// complete application version
				return a.GetAvailableVersionsFor(args[len(args)-1]), cobra.ShellCompDirectiveNoFileComp
			default:
				// huh ?
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
		},
	}

	// cmd.Flags().StringVarP(&bindir, "bindir", "b", app.GetDefaultBinDir(), "Binaries directory")

	return cmd
}
