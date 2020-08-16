package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/devopsworks/tools/binenv/internal/app"
)

// localCmd represents the local command
func localCmd() *cobra.Command {
	app := app.New()
	cmd := &cobra.Command{
		Use:   "local",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Local()
		},
	}

	// cmd.Flags().IntVarP(&a.Params.MinLength, "min-length", "m", 16, "Specify minimum password length, must not be less than 8")
	return cmd
}
