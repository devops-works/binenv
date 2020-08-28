package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
func updateCmd(a *app.App) *cobra.Command {
	var definitionsOnly, definitionsAlso bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update available software distributions",
		Long: `Available versions listed distribution will be updated.
If not distribution is specified, versions for all distributions will be updated.`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			a.SetVerbose(verbose)

			if len(args) >= 1 {
				a.Update(definitionsOnly, definitionsAlso, args...)
				return
			}
			a.Update(definitionsOnly, definitionsAlso)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// Remove already selected distributions from completion
			list := a.GetPackagesListWithPrefix(toComplete)
			list = removeFromSlice(list, args)
			return list, cobra.ShellCompDirectiveNoFileComp
		},
	}

	cmd.Flags().BoolVarP(&definitionsOnly, "definitions", "d", false, "Update only distributions definitions")
	cmd.Flags().BoolVarP(&definitionsAlso, "all", "a", false, "Update distributions definitions and distributions versions")
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
