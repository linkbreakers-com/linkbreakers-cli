package cli

import (
	"fmt"
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
	root := newRootCommand()
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if shouldCheckForUpdates(os.Args[1:]) {
		update.MaybeNotify(Version, os.Stderr)
	}
	return 0
}

func newRootCommand() *cobra.Command {
	a := &app{}

	cmd := &cobra.Command{
		Use:           "linkbreakers",
		Short:         "Official CLI for the Linkbreakers API",
		Long:          "The official Linkbreakers CLI for authentication, workspace resources, raw API access, shell completion, and generated command documentation.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: "" +
			"  linkbreakers auth set-token --token <api-token>\n" +
			"  linkbreakers links list --page-size 20\n" +
			"  linkbreakers links create --destination https://example.com --name \"My link\"\n" +
			"  linkbreakers raw GET /v1/links?pageSize=5\n",
	}

	cmd.PersistentFlags().StringVar(&a.tokenFlag, "token", "", "API token. Overrides config and LINKBREAKERS_TOKEN.")
	cmd.PersistentFlags().StringVar(&a.baseURLFlag, "base-url", "", "API base URL. Overrides config and LINKBREAKERS_BASE_URL.")
	cmd.PersistentFlags().StringVarP(&a.outputFlag, "output", "o", "", "Output format: json or table.")

	cmd.AddCommand(
		a.newAuthCommand(),
		a.newCompletionCommand(cmd),
		a.newDirectoriesCommand(),
		a.newCustomDomainsCommand(),
		a.newDocsCommand(cmd),
		a.newLinksCommand(),
		a.newRawCommand(),
		a.newSelfUpdateCommand(),
		a.newVersionCommand(),
	)

	return cmd
}

func (a *app) runtimeConfig() (config.RuntimeConfig, error) {
	return config.Resolve(a.tokenFlag, a.baseURLFlag, a.outputFlag)
}

func (a *app) requireClient() (*api.Client, config.RuntimeConfig, error) {
	cfg, err := a.runtimeConfig()
	if err != nil {
		return nil, config.RuntimeConfig{}, err
	}
	if err := cfg.RequireToken(); err != nil {
		return nil, config.RuntimeConfig{}, err
	}
	return api.NewClient(cfg), cfg, nil
}

func docsDir() string {
	return filepath.Join("docs", "commands")
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
