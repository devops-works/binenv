package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
func updateCmd() *cobra.Command {
	app, err := app.New()
	if err != nil {
		panic(err)
	}

	var definitionsOnly, definitionsAlso bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update available software distributions",
		Long: `Available versions for every distribution will be updated.
`,
		// Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return app.Update(args[0], definitionsOnly, definitionsAlso)
			}
			return app.Update("", definitionsOnly, definitionsAlso)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// Remove already selected distributions from completion
			list := app.GetPackagesListWithPrefix(toComplete)
			list = removeFromSlice(list, args)
			return list, cobra.ShellCompDirectiveNoFileComp
		},
	}
	// verb, _ := cmd.Root().PersistentFlags().GetBool("verbose")

	// fmt.Printf("verbose is %v\n", verb)
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
