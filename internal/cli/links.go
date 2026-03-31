package cli

import (
	"fmt"

	linkbreakers "github.com/linkbreakers-com/linkbreakers-cli/internal/client/generated"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *app) newLinksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links",
		Short: "Manage links",
	}

	cmd.AddCommand(
		a.newLinksListCommand(),
		a.newLinksGetCommand(),
		a.newLinksCreateCommand(),
		a.newLinksDeleteCommand(),
	)

	return cmd
}

func (a *app) newLinksListCommand() *cobra.Command {
	var pageSize int64
	var pageToken string
	var search string
	var include string
	var tags string
	var sortBy string
	var sortDirection string
	var directoryID string
	var includeAllDirectories bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List links",
		Example: "" +
			"  linkbreakers links list --page-size 20\n" +
			"  linkbreakers links list --search summer --tags campaign-a,campaign-b\n" +
			"  linkbreakers links list --output table --include qrcodeSignedUrl\n",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, cfg, err := a.requireClient()
			if err != nil {
				return err
			}

			req := client.Generated().LinksAPI.LinksServiceList(client.Context()).PageSize(pageSize)
			if pageToken != "" {
				req = req.PageToken(pageToken)
			}
			if search != "" {
				req = req.Search(search)
			}
			if include != "" {
				req = req.Include(splitCSV(include))
			}
			if tags != "" {
				req = req.Tags(splitCSV(tags))
			}
			if sortBy != "" {
				req = req.SortBy(sortBy)
			}
			if sortDirection != "" {
				req = req.SortDirection(sortDirection)
			}
			if directoryID != "" {
				req = req.DirectoryId(directoryID)
			}
			if includeAllDirectories {
				req = req.IncludeAllDirectories(true)
			}

			resp, _, err := req.Execute()
			if err != nil {
				return err
			}

			if cfg.Output == "table" {
				rows := make([][]string, 0, len(resp.GetLinks()))
				for _, item := range resp.GetLinks() {
					rows = append(rows, []string{
						item.GetId(),
						item.GetName(),
						item.GetShortlink(),
						item.GetEntrypoint(),
						item.GetEventCount(),
					})
				}
				return output.PrintTable([]string{"ID", "NAME", "SHORTLINK", "ENTRYPOINT", "EVENTS"}, rows)
			}

			return output.PrintJSON(resp)
		},
	}

	cmd.Flags().Int64Var(&pageSize, "page-size", 20, "Maximum number of links to return.")
	cmd.Flags().StringVar(&pageToken, "page-token", "", "Pagination cursor.")
	cmd.Flags().StringVar(&search, "search", "", "Search by name or shortlink.")
	cmd.Flags().StringVar(&include, "include", "", "Comma-separated related resources to include.")
	cmd.Flags().StringVar(&tags, "tags", "", "Comma-separated tags filter.")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort field.")
	cmd.Flags().StringVar(&sortDirection, "sort-direction", "", "Sort direction.")
	cmd.Flags().StringVar(&directoryID, "directory-id", "", "Filter by directory ID.")
	cmd.Flags().BoolVar(&includeAllDirectories, "include-all-directories", false, "Include links from all directories.")
	return cmd
}

func (a *app) newLinksGetCommand() *cobra.Command {
	var include string

	cmd := &cobra.Command{
		Use:   "get ID",
		Short: "Get link details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			req := client.Generated().LinksAPI.LinksServiceGet(client.Context(), args[0])
			if include != "" {
				req = req.Include(splitCSV(include))
			}

			resp, _, err := req.Execute()
			if err != nil {
				return err
			}

			return output.PrintJSON(resp)
		},
	}

	cmd.Flags().StringVar(&include, "include", "", "Comma-separated related resources to include.")
	return cmd
}

func (a *app) newLinksCreateCommand() *cobra.Command {
	var destination string
	var name string
	var shortlink string
	var directoryID string
	var fallbackDestination string
	var customDomainID string
	var qrCodeDesignID string
	var qrCodeTemplateID string
	var waitForQRCode bool
	var conversionTracking bool
	var tags string
	var metadata []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a link",
		Example: "" +
			"  linkbreakers links create --destination https://example.com\n" +
			"  linkbreakers links create --destination https://example.com --name \"Launch\" --shortlink launch\n" +
			"  linkbreakers links create --destination https://example.com --tags campaign-a,campaign-b --metadata owner=marketing\n",
		RunE: func(cmd *cobra.Command, args []string) error {
			if destination == "" {
				return fmt.Errorf("--destination is required")
			}

			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			body := linkbreakers.NewCreateLinkRequest(destination)
			if name != "" {
				body.SetName(name)
			}
			if shortlink != "" {
				body.SetShortlink(shortlink)
			}
			if directoryID != "" {
				body.SetDirectoryId(directoryID)
			}
			if fallbackDestination != "" {
				body.SetFallbackDestination(fallbackDestination)
			}
			if customDomainID != "" {
				body.SetCustomDomainId(customDomainID)
			}
			if qrCodeDesignID != "" {
				body.SetQrcodeDesignId(qrCodeDesignID)
			}
			if qrCodeTemplateID != "" {
				body.SetQrcodeTemplateId(qrCodeTemplateID)
			}
			if waitForQRCode {
				body.SetWaitForQrcode(true)
			}
			if conversionTracking {
				body.SetConversionTracking(true)
			}
			if tags != "" {
				body.Tags = splitCSV(tags)
			}
			if len(metadata) > 0 {
				values, err := parseStringMap(metadata)
				if err != nil {
					return err
				}
				body.Metadata = values
			}

			resp, _, err := client.Generated().LinksAPI.LinksServiceCreate(client.Context()).CreateLinkRequest(*body).Execute()
			if err != nil {
				return err
			}

			return output.PrintJSON(resp)
		},
	}

	cmd.Flags().StringVar(&destination, "destination", "", "Destination URL.")
	cmd.Flags().StringVar(&name, "name", "", "Optional display name.")
	cmd.Flags().StringVar(&shortlink, "shortlink", "", "Optional shortlink path.")
	cmd.Flags().StringVar(&directoryID, "directory-id", "", "Directory ID.")
	cmd.Flags().StringVar(&fallbackDestination, "fallback-destination", "", "Fallback destination URL.")
	cmd.Flags().StringVar(&customDomainID, "custom-domain-id", "", "Custom domain ID.")
	cmd.Flags().StringVar(&qrCodeDesignID, "qrcode-design-id", "", "QR code design ID.")
	cmd.Flags().StringVar(&qrCodeTemplateID, "qrcode-template-id", "", "QR code template ID.")
	cmd.Flags().BoolVar(&waitForQRCode, "wait-for-qrcode", false, "Wait for QR code generation.")
	cmd.Flags().BoolVar(&conversionTracking, "conversion-tracking", false, "Preserve lbid during redirects.")
	cmd.Flags().StringVar(&tags, "tags", "", "Comma-separated tags.")
	cmd.Flags().StringArrayVar(&metadata, "metadata", nil, "Metadata in key=value form. Repeatable.")
	return cmd
}

func (a *app) newLinksDeleteCommand() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete ID",
		Short: "Delete a link",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				return fmt.Errorf("refusing to delete without --yes")
			}

			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			resp, _, err := client.Generated().LinksAPI.LinksServiceDelete(client.Context(), args[0]).Execute()
			if err != nil {
				return err
			}

			return output.PrintJSON(resp)
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm deletion.")
	return cmd
}
