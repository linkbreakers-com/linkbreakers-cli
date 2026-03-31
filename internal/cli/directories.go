package cli

import (
	"fmt"

	linkbreakers "github.com/linkbreakers-com/linkbreakers-cli/internal/client/generated"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *app) newDirectoriesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "directories",
		Short: "Manage directories",
	}

	cmd.AddCommand(
		a.newDirectoriesListCommand(),
		a.newDirectoriesGetCommand(),
		a.newDirectoriesCreateCommand(),
		a.newDirectoriesDeleteCommand(),
	)

	return cmd
}

func (a *app) newDirectoriesListCommand() *cobra.Command {
	var pageSize int64
	var pageToken string
	var parentDirectoryID string
	var search string
	var includeRoot bool
	var recursive bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List directories",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, cfg, err := a.requireClient()
			if err != nil {
				return err
			}

			req := client.Generated().DirectoriesAPI.DirectoriesServiceList(client.Context()).PageSize(pageSize)
			if pageToken != "" {
				req = req.PageToken(pageToken)
			}
			if parentDirectoryID != "" {
				req = req.ParentDirectoryId(parentDirectoryID)
			}
			if search != "" {
				req = req.Search(search)
			}
			if includeRoot {
				req = req.IncludeRoot(true)
			}
			if recursive {
				req = req.Recursive(true)
			}

			resp, _, err := req.Execute()
			if err != nil {
				return err
			}

			if cfg.Output == "table" {
				rows := make([][]string, 0, len(resp.GetDirectories()))
				for _, item := range resp.GetDirectories() {
					rows = append(rows, []string{
						item.GetId(),
						item.GetName(),
						item.GetPath(),
						item.GetParentDirectoryId(),
					})
				}
				return output.PrintTable([]string{"ID", "NAME", "PATH", "PARENT_ID"}, rows)
			}

			return output.PrintJSON(resp)
		},
	}

	cmd.Flags().Int64Var(&pageSize, "page-size", 100, "Maximum number of directories to return.")
	cmd.Flags().StringVar(&pageToken, "page-token", "", "Pagination cursor.")
	cmd.Flags().StringVar(&parentDirectoryID, "parent-directory-id", "", "Filter by parent directory ID.")
	cmd.Flags().StringVar(&search, "search", "", "Search by directory name.")
	cmd.Flags().BoolVar(&includeRoot, "include-root", false, "Include root directories when filtering by parent.")
	cmd.Flags().BoolVar(&recursive, "recursive", false, "Return descendants recursively.")
	return cmd
}

func (a *app) newDirectoriesGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get ID",
		Short: "Get directory details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			resp, _, err := client.Generated().DirectoriesAPI.DirectoriesServiceGet(client.Context(), args[0]).Execute()
			if err != nil {
				return err
			}

			return output.PrintJSON(resp)
		},
	}
}

func (a *app) newDirectoriesCreateCommand() *cobra.Command {
	var name string
	var parentDirectoryID string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}

			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			body := linkbreakers.NewCreateDirectoryRequest(name)
			if parentDirectoryID != "" {
				body.SetParentDirectoryId(parentDirectoryID)
			}

			resp, _, err := client.Generated().DirectoriesAPI.DirectoriesServiceCreate(client.Context()).CreateDirectoryRequest(*body).Execute()
			if err != nil {
				return err
			}

			return output.PrintJSON(resp)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Directory name.")
	cmd.Flags().StringVar(&parentDirectoryID, "parent-directory-id", "", "Optional parent directory ID.")
	return cmd
}

func (a *app) newDirectoriesDeleteCommand() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete ID",
		Short: "Delete a directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				return fmt.Errorf("refusing to delete without --yes")
			}

			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			resp, _, err := client.Generated().DirectoriesAPI.DirectoriesServiceDelete(client.Context(), args[0]).Execute()
			if err != nil {
				return err
			}

			return output.PrintJSON(resp)
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm deletion.")
	return cmd
}
