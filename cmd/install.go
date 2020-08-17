package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/devopsworks/tools/binenv/internal/app"
)

// localCmd represents the local command
func installCmd() *cobra.Command {
	var bindir string

	a, err := app.New()
	if err != nil {
		panic(err)
	}
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install a version for the package",
		Long:  `Better description here`,
		RunE: func(cmd *cobra.Command, args []string) error {
			a.SetBinDir(bindir)
			if len(args) > 1 {
				return a.Install(args[0], args[1])
			}
			return a.Install(args[0], "")
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0:
				// complete application name
				return a.GetPackagesListWithPrefix(toComplete), cobra.ShellCompDirectiveNoFileComp
			case 1:
				// complete application version
				return a.GetVersionsFromCacheFor(args[0]), cobra.ShellCompDirectiveNoFileComp
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
