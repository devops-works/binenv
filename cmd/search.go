package cmd

import (
	"fmt"

	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
func searchCmd(a *app.App) *cobra.Command {
	var (
		wide bool
	)

	cmd := &cobra.Command{
		Use:   "search [term]",
		Short: "Search term in software distributions",
		Long:  `Search a term in distribution names or descriptions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			global, _ := cmd.Flags().GetBool("global")
			bindir, _ := cmd.Flags().GetString("bindir")

			a.SetVerbose(verbose)
			a.SetBinDir(bindir)
			a.SetGlobal(global)

			switch {
			case len(args) > 1:
				return fmt.Errorf("command requires a single argument")
			case len(args) == 0:
				a.Search("", wide)
			default:
				a.Search(args[0], wide)

			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&wide, "wide", "w", false, "Wide output")

	return cmd
}
