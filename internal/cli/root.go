package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/api"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/config"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/update"
	"github.com/spf13/cobra"
)

type app struct {
	tokenFlag   string
	baseURLFlag string
	outputFlag  string
}

func Execute() int {
	a := &app{}
	root := a.newRootCommand()
	err := root.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
	}
	if shouldRemindClaim(os.Args[1:]) {
		a.claimReminder(context.Background())
	}
	if err != nil {
		return 1
	}

	if shouldCheckForUpdates(os.Args[1:]) {
		update.MaybeNotify(Version, os.Stderr)
	}
	return 0
}

func (a *app) newRootCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:           "linkbreakers",
		Short:         "Official CLI for the Linkbreakers API",
		Long:          "The official Linkbreakers CLI. Every API operation is a command, generated from the API's OpenAPI document, so the CLI always matches the API.\n\nNo account yet? `linkbreakers signup --email <email>` creates a workspace and saves its token in one step.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: "" +
			"  linkbreakers signup --email you@example.com\n" +
			"  linkbreakers links create --destination https://example.com --name \"My link\"\n" +
			"  linkbreakers links list --page-size 20\n" +
			"  linkbreakers analytics-data get-top-links\n",
	}

	cmd.PersistentFlags().StringVar(&a.tokenFlag, "token", "", "API token. Overrides config and LINKBREAKERS_TOKEN.")
	cmd.PersistentFlags().StringVar(&a.baseURLFlag, "base-url", "", "API base URL. Overrides config and LINKBREAKERS_BASE_URL.")
	cmd.PersistentFlags().StringVarP(&a.outputFlag, "output", "o", "", "Output format: json.")

	cmd.AddCommand(
		a.newAuthCommand(),
		a.newClaimCommand(),
		a.newCompletionCommand(cmd),
		a.newDocsCommand(cmd),
		a.newSelfUpdateCommand(),
		a.newSignupCommand(),
		a.newVersionCommand(),
	)
	a.addOperationCommands(cmd)

	return cmd
}

func (a *app) runtimeConfig() (config.RuntimeConfig, error) {
	return config.Resolve(a.tokenFlag, a.baseURLFlag, a.outputFlag)
}

// client is the API client for a command; authenticated operations need a
// token, public ones do not.
func (a *app) client(auth bool) (*api.Client, error) {
	cfg, err := a.runtimeConfig()
	if err != nil {
		return nil, err
	}
	if auth {
		if err := cfg.RequireToken(); err != nil {
			return nil, err
		}
	}
	return api.NewClient(cfg, userAgent()), nil
}

func userAgent() string {
	return "linkbreakers-cli/" + Version
}

func readAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func docsDir() string {
	return filepath.Join("docs", "commands")
}

// shouldRemindClaim leaves out the commands that are about the claim itself
// or say nothing about the workspace.
func shouldRemindClaim(args []string) bool {
	if len(args) == 0 {
		return false
	}
	for _, arg := range args {
		switch arg {
		case "help", "--help", "-h", "completion", "gendocs", "self-update", "version", "signup", "claim", "auth":
			return false
		}
	}
	return true
}

func shouldCheckForUpdates(args []string) bool {
	if len(args) == 0 {
		return false
	}

	for _, arg := range args {
		switch arg {
		case "help", "--help", "-h", "completion", "gendocs", "self-update", "version":
			return false
		}
	}

	return true
}
