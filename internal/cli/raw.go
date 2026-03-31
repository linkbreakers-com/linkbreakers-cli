package cli

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/api"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *app) newRawCommand() *cobra.Command {
	var body string
	var bodyFile string
	var headers []string

	cmd := &cobra.Command{
		Use:   "raw METHOD PATH",
		Short: "Call any Linkbreakers API endpoint directly",
		Long:  "Fallback command for endpoints that do not yet have first-class CLI coverage.",
		Example: "" +
			"  linkbreakers raw GET /v1/links?pageSize=5\n" +
			"  linkbreakers raw POST /v1/links --body '{\"destination\":\"https://example.com\"}'\n" +
			"  linkbreakers raw PATCH /v1/links/<id> --body-file link.json\n",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, err := a.requireClient()
			if err != nil {
				return err
			}

			rawBody, err := output.ReadBody(body, bodyFile)
			if err != nil {
				return err
			}
			if len(rawBody) > 0 {
				rawBody, err = api.NormalizeJSON(rawBody)
				if err != nil {
					return fmt.Errorf("invalid JSON body: %w", err)
				}
			}

			method := strings.ToUpper(args[0])
			pathArg := args[1]
			parsed, err := url.Parse(pathArg)
			if err != nil {
				return fmt.Errorf("parse path: %w", err)
			}

			query := map[string]string{}
			for key, values := range parsed.Query() {
				if len(values) > 0 {
					query[key] = values[len(values)-1]
				}
			}

			headerMap := map[string]string{}
			for _, header := range headers {
				parts := strings.SplitN(header, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid --header %q, expected key=value", header)
				}
				headerMap[parts[0]] = parts[1]
			}

			status, respBody, err := client.RawRequest(context.Background(), method, parsed.Path, headerMap, query, rawBody)
			if err != nil {
				return err
			}

			if len(respBody) == 0 {
				return output.PrintJSON(map[string]any{"status": status, "ok": true})
			}

			var decoded any
			if err := parseJSON(respBody, &decoded); err != nil {
				return output.PrintJSON(map[string]any{
					"status": status,
					"body":   string(respBody),
				})
			}

			return output.PrintJSON(map[string]any{
				"status": status,
				"body":   decoded,
			})
		},
	}

	cmd.Flags().StringVar(&body, "body", "", "Inline JSON request body.")
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to a JSON request body file.")
	cmd.Flags().StringArrayVar(&headers, "header", nil, "Extra header in key=value form. Repeatable.")
	return cmd
}
