package cli

import (
	"github.com/linkbreakers-com/linkbreakers-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *app) newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show CLI build version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return output.PrintJSON(map[string]string{
				"version": Version,
				"commit":  Commit,
				"date":    Date,
			})
		},
	}
}
