package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devops-works/binenv/internal/app"
)

// localCmd represents the local command
func localCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local <distribution> <version> [<distribution> <version>]",
		Short: "Sets local required versions for distributions.",
		Long: `This will write the specified version in ".binenv.lock" file in the current directory.
Any previously constraint used in this file for the distribution will be removed, and an exact match ('=') will be used.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			a.Local(args[0], args[1])
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

	return cmd
}
