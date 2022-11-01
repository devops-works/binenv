package cmd

import (
	"os"

	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
func installCmd(a *app.App) *cobra.Command {
	var fromlock, dryrun bool

	cmd := &cobra.Command{
		Use:   "install [--lock] [--dry-run] [<distribution> <version> [<distribution> <version>]]",
		Short: "Install a version for the package",
		Long: `This command will install one or several distributions with the specified versions. 
If --lock is used, versions from the .binenv.lock file in the current directory will be installed.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 && !fromlock {
				cmd.Help()
				os.Exit(1)
			}
			a.SetDryRun(dryrun)

			if fromlock {
				a.InstallFromLock()
				return
			}

			a.Install(args...)
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

	cmd.Flags().BoolVarP(&fromlock, "lock", "l", false, "Install versions specified in ./.binenv.lock")
	cmd.Flags().BoolVarP(&dryrun, "dry-run", "n", false, "Do not install, just simulate")

	return cmd
}
