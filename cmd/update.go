package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
func updateCmd(a *app.App) *cobra.Command {
	var (
		distributionsOnly, distributionsAlso, noCache bool
		concurrency                                   int
	)

	cmd := &cobra.Command{
		Use:   "update [--all|--distributions] [--cache]",
		Short: "Update available software distributions",
		Long: `Available versions listed distribution will be updated.
If not distribution is specified, versions for all distributions will be updated.`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			a.SetVerbose(verbose)
			a.SetConcurrency(concurrency)

			a.Initialize()

			if len(args) >= 1 {
				a.Update(distributionsOnly, distributionsAlso, noCache, args...)
				return
			}
			a.Update(distributionsOnly, distributionsAlso, noCache)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// Remove already selected distributions from completion
			list := a.GetPackagesListWithPrefix(toComplete)
			list = removeFromSlice(list, args)
			return list, cobra.ShellCompDirectiveNoFileComp
		},
	}

	cmd.Flags().BoolVarP(&distributionsOnly, "distributions", "d", false, "Update only distributions")
	cmd.Flags().BoolVarP(&distributionsAlso, "all", "a", false, "Update distributions and distributions versions")
	cmd.Flags().BoolVarP(&noCache, "nocache", "f", false, "Distributions versions will be updated from each release and not from github cache")
	cmd.Flags().IntVarP(&concurrency, "concurrency", "c", 8, "Concurrency for cache update")
	return cmd
}

func removeFromSlice(orig, remove []string) []string {
	for i := 0; i < len(orig); i++ {
		url := orig[i]
		for _, rem := range remove {
			if url == rem {
				orig = append(orig[:i], orig[i+1:]...)
				i--
				break
			}
		}
	}

	return orig
}
