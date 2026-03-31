package cli

import (
	"fmt"

	linkbreakers "github.com/linkbreakers-com/linkbreakers-cli/internal/client/generated"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *app) newCustomDomainsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "custom-domains",
		Aliases: []string{"domains"},
		Short:   "Manage custom domains",
	}

	cmd.AddCommand(
		a.newCustomDomainsListCommand(),
		a.newCustomDomainsGetCommand(),
		a.newCustomDomainsCreateCommand(),
		a.newCustomDomainsCheckCommand(),
		a.newCustomDomainsDeleteCommand(),
	)

	return cmd
}

func (a *app) newCustomDomainsListCommand() *cobra.Command {
	var status string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List custom domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, cfg, err := a.requireClient()
			if err != nil {
				return err
			}

			req := client.Generated().CustomDomainsAPI.CustomDomainsServiceList(client.Context())
			if status != "" {
				req = req.Status(status)
			}

			resp, _, err := req.Execute()
			if err != nil {
				return err
			}

			if cfg.Output == "table" {
				rows := make([][]string, 0, len(resp.GetCustomDomains()))
				for _, item := range resp.GetCustomDomains() {
					rows = append(rows, []string{
						item.GetId(),
						item.GetName(),
						item.GetStatus(),
						item.GetRootDestination(),
					})
				}
				return output.PrintTable([]string{"ID", "NAME", "STATUS", "ROOT_DESTINATION"}, rows)
			}

			return output.PrintJSON(resp)
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "Filter by custom domain status.")
	return cmd
}

func (a *app) newCustomDomainsGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get ID",
		Short: "Get custom domain details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			resp, _, err := client.Generated().CustomDomainsAPI.CustomDomainsServiceGet(client.Context(), args[0]).Execute()
			if err != nil {
				return err
			}

			return output.PrintJSON(resp)
		},
	}
}

func (a *app) newCustomDomainsCreateCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a custom domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}

			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			body := linkbreakers.NewCreateCustomDomainRequest()
			body.SetName(name)

			resp, _, err := client.Generated().CustomDomainsAPI.CustomDomainsServiceCreate(client.Context()).CreateCustomDomainRequest(*body).Execute()
			if err != nil {
				return err
			}

			return output.PrintJSON(resp)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Domain name to register.")
	return cmd
}

func (a *app) newCustomDomainsCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check ID",
		Short: "Re-check DNS and TLS status for a custom domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			resp, _, err := client.Generated().CustomDomainsAPI.CustomDomainsServiceCheck(client.Context(), args[0]).Execute()
			if err != nil {
				return err
			}

			return output.PrintJSON(resp)
		},
	}
}

func (a *app) newCustomDomainsDeleteCommand() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete ID",
		Short: "Delete a custom domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				return fmt.Errorf("refusing to delete without --yes")
			}

			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			resp, _, err := client.Generated().CustomDomainsAPI.CustomDomainsServiceDelete(client.Context(), args[0]).Execute()
			if err != nil {
				return err
			}

			return output.PrintJSON(resp)
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm deletion.")
	return cmd
}
