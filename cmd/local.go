package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
func localCmd() *cobra.Command {
	app, err := app.New()
	if err != nil {
		panic(err)
	}
	cmd := &cobra.Command{
		Use:   "local <distribution> <version> [<distribution> <version>]",
		Short: "Sets local required versions for distributions.",
		Long: `This will write the specified version in ".binenv.lock" file in the current directory.
Any previously constraint used in this file for the distribution will be removed, and an exact match ('=') will be used.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Local(args[0], args[1])
		},
	}

	return cmd
}
