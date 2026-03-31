package cli

import "github.com/spf13/cobra"

func (a *app) newCompletionCommand(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "bash",
			Short: "Generate Bash completion",
			RunE: func(cmd *cobra.Command, args []string) error {
				return root.GenBashCompletion(cmd.OutOrStdout())
			},
		},
		&cobra.Command{
			Use:   "zsh",
			Short: "Generate Zsh completion",
			RunE: func(cmd *cobra.Command, args []string) error {
				return root.GenZshCompletion(cmd.OutOrStdout())
			},
		},
		&cobra.Command{
			Use:   "fish",
			Short: "Generate Fish completion",
			RunE: func(cmd *cobra.Command, args []string) error {
				return root.GenFishCompletion(cmd.OutOrStdout(), true)
			},
		},
		&cobra.Command{
			Use:   "powershell",
			Short: "Generate PowerShell completion",
			RunE: func(cmd *cobra.Command, args []string) error {
				return root.GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			},
		},
	)

	return cmd
}
