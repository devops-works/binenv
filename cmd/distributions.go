package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/devopsworks/tools/binenv/internal/app"
)

// localCmd represents the local command
func distributionsCmd() *cobra.Command {
	app, err := app.New()
	if err != nil {
		panic(err)
	}
	cmd := &cobra.Command{
		Use:   "distributions",
		Short: "Show or update distributions list",
		Long: `Show or update what distributions can be installed. It does not
update distribution versions (see "binenv update [distribution]" for this), 
only the list of distributions (e.g. applications) that can be installed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Distributions()
		},
	}

	// cmd.Flags().IntVarP(&a.Params.MinLength, "min-length", "m", 16, "Specify minimum password length, must not be less than 8")
	return cmd
}
