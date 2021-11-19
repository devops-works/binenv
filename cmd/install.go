package cmd

import (
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

			verbose, _ := cmd.Flags().GetBool("verbose")
			// bindir, _ := cmd.Flags().GetString("bindir")
			// cachedir, _ := cmd.Flags().GetString("cachedir")
			// distdir, _ := cmd.Flags().GetString("distdir")
			a.SetVerbose(verbose)
			a.SetDryRun(dryrun)
			// a.SetBindir(bindir)
			// a.SetCachedir(cachedir)
			// a.SetDistdir(distdir)

			a.Initialize()

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

	// cmd.Flags().StringVarP(&bindir, "bindir", "b", app.GetDefaultBinDir(), "Binaries directory")
	cmd.Flags().BoolVarP(&fromlock, "lock", "l", false, "Install versions specified in ./.binenv.lock")
	cmd.Flags().BoolVarP(&dryrun, "dry-run", "n", false, "Do not install, just simulate")
	// cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose operations")

	return cmd
}
