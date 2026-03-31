package cli

import (
	"fmt"
	"os"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/config"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *app) newAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage CLI authentication",
	}

	cmd.AddCommand(
		a.newAuthSetTokenCommand(),
		a.newAuthStatusCommand(),
		a.newAuthClearCommand(),
	)

	return cmd
}

func (a *app) newAuthSetTokenCommand() *cobra.Command {
	var token string
	var baseURL string

	cmd := &cobra.Command{
		Use:   "set-token",
		Short: "Persist an API token for future commands",
		Example: "" +
			"  linkbreakers auth set-token --token <api-token>\n" +
			"  linkbreakers auth set-token --token <api-token> --base-url https://api.linkbreakers.com\n",
		RunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				return fmt.Errorf("--token is required")
			}

			fileCfg, err := config.LoadFileConfig()
			if err != nil {
				return err
			}

			fileCfg.Token = token
			if baseURL != "" {
				fileCfg.BaseURL = baseURL
			}

			if err := config.SaveFileConfig(fileCfg); err != nil {
				return err
			}

			path, _ := config.ConfigPath()
			return output.PrintJSON(map[string]any{
				"ok":          true,
				"config_path": path,
				"base_url":    fileCfg.BaseURL,
			})
		},
	}

	cmd.Flags().StringVar(&token, "token", "", "Linkbreakers API token.")
	cmd.Flags().StringVar(&baseURL, "base-url", "", "Optional API base URL to persist.")
	return cmd
}

func (a *app) newAuthStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show how authentication is currently resolved",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := a.runtimeConfig()
			if err != nil {
				return err
			}

			path, _ := config.ConfigPath()
			return output.PrintJSON(map[string]any{
				"has_token":          cfg.Token != "",
				"base_url":           cfg.BaseURL,
				"output":             cfg.Output,
				"config_path":        path,
				"env_token_override": os.Getenv("LINKBREAKERS_TOKEN") != "",
			})
		},
	}
}

func (a *app) newAuthClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Delete persisted CLI config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.DeleteFileConfig(); err != nil {
				return err
			}
			return output.PrintJSON(map[string]any{"ok": true})
		},
	}
}
