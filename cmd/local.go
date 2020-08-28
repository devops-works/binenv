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
			verbose, _ := cmd.Flags().GetBool("verbose")
			a.SetVerbose(verbose)
			a.Local(args[0], args[1])
		},
	}

	return cmd
}
