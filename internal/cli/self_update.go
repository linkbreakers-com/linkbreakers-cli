package cli

import (
	"fmt"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/output"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/update"
	"github.com/spf13/cobra"
)

func (a *app) newSelfUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "self-update",
		Short: "Download and install the latest CLI release",
		Long:  "Updates the current linkbreakers binary in place when supported by the current platform.",
		RunE: func(cmd *cobra.Command, args []string) error {
			release, err := update.SelfUpdate(Version)
			if err != nil {
				return err
			}

			message := fmt.Sprintf("linkbreakers is already up to date (%s)", release.Version)
			if Version == "" || Version == "dev" || release.Version != Version {
				message = fmt.Sprintf("updated linkbreakers to %s", release.Version)
			}

			return output.PrintJSON(map[string]any{
				"ok":          true,
				"version":     release.Version,
				"release_url": release.ReleaseURL,
				"message":     message,
			})
		},
	}
}
