package cmd

import (
	"github.com/devops-works/binenv/internal/app"
	"github.com/spf13/cobra"
)

func expandCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expand [term]",
		Short: "Return full distribution binary path.",
		Long:  `Return full installed distribution binary path.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			a.SetFlag("justExpand", true)
			a.Execute(args)
			return nil
		},
	}

	return cmd

}
