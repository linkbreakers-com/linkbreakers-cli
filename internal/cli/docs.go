package cli

import (
	"fmt"

	cliDocs "github.com/linkbreakers-com/linkbreakers-cli/internal/docs"
	"github.com/spf13/cobra"
)

func (a *app) newDocsCommand(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:    "gendocs",
		Short:  "Generate markdown command docs and llms.txt",
		Long:   "Maintainer command used in CI to regenerate CLI command docs and llms.txt.",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cliDocs.Generate(root, docsDir()); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "generated docs in %s\n", docsDir())
			return nil
		},
	}
}
