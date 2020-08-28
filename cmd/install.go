package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
func installCmd(a *app.App) *cobra.Command {
	var bindir string
	var fromlock bool

	cmd := &cobra.Command{
		Use:   "install [--lock] [<distribution> <version> [<distribution> <version>]]",
		Short: "Install a version for the package",
		Long: `This command will install one or several distributions with the specified versions. 
If --lock is used, versions from the .binenv.lock file in the current directory will be installed.`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")

			a.SetVerbose(verbose)
			a.SetBinDir(bindir)
			if fromlock {
				a.InstallFromLock()
			}
			if len(args) > 0 {
				a.Install(args...)
			}
		},
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

	cmd.Flags().StringVarP(&bindir, "bindir", "b", app.GetDefaultBinDir(), "Binaries directory")
	cmd.Flags().BoolVarP(&fromlock, "lock", "l", false, "Install versions specified in ./.binenv.lock")

	return cmd
}
