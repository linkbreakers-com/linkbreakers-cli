package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/api"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/config"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/output"
	"github.com/spf13/cobra"
)

// Self-registration lives on the MCP host, outside the REST API, so these
// commands are the only ones not generated from the OpenAPI document.

type registration struct {
	AccessToken       string   `json:"access_token"`
	WorkspaceID       string   `json:"workspace_id"`
	ClaimToken        string   `json:"claim_token"`
	ClaimTokenExpires string   `json:"claim_token_expires"`
	PreClaimScopes    []string `json:"pre_claim_scopes"`
	PostClaimScopes   []string `json:"post_claim_scopes"`
}

type claimStarted struct {
	ExpiresAt    string `json:"expires_at"`
	ClaimAttempt struct {
		Interval int `json:"interval"`
	} `json:"claim_attempt"`
}

type agentToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

type agentError struct {
	Status      int
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

func (e *agentError) Error() string {
	if e.Description == "" {
		return fmt.Sprintf("%d %s", e.Status, e.Code)
	}
	return fmt.Sprintf("%d %s: %s", e.Status, e.Code, e.Description)
}

func (a *app) newSignupCommand() *cobra.Command {
	var email string
	var force bool

	cmd := &cobra.Command{
		Use:   "signup",
		Short: "Create a Linkbreakers workspace and save its API token, in one step",
		Long: "Creates a free Linkbreakers workspace and saves an API token for it, so every other command works right away.\n\n" +
			"The workspace is handed to the person whose email you pass: they receive a link and a 6-digit code to claim it, by email only. " +
			"Until they do, the token can create and manage links, QR codes, workflows and page themes and read analytics. " +
			"An unclaimed workspace is deleted with everything in it 24 hours after signup.\n\n" +
			"Already have an account? Create a token in the dashboard and run `linkbreakers auth set-token` instead.",
		Example: "" +
			"  linkbreakers signup --email you@example.com\n" +
			"  linkbreakers links create --destination https://example.com --name \"My first link\"\n",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			email = strings.TrimSpace(email)
			if email == "" {
				return fmt.Errorf("--email is required: the address of the person the workspace is for. They receive the link to claim it")
			}
			fileCfg, err := config.LoadFileConfig()
			if err != nil {
				return err
			}
			if fileCfg.Token != "" && !force {
				return fmt.Errorf("a token is already configured. Run `linkbreakers auth status` to see it, or pass --force to replace it with a new workspace")
			}
			cfg, err := a.runtimeConfig()
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), cfg.Timeout)
			defer cancel()
			var reg registration
			if err := postAgent(ctx, cfg.AgentBaseURL+"/agent/identity", map[string]string{"type": "anonymous"}, &reg); err != nil {
				return fmt.Errorf("could not create a workspace: %w", err)
			}
			expires, _ := time.Parse(time.RFC3339, reg.ClaimTokenExpires)

			fileCfg.Token = reg.AccessToken
			fileCfg.Claim = &config.Claim{WorkspaceID: reg.WorkspaceID, Email: email, ClaimToken: reg.ClaimToken, ExpiresAt: expires}
			if err := config.SaveFileConfig(fileCfg); err != nil {
				return err
			}

			result := map[string]any{
				"workspace_id": reg.WorkspaceID,
				"token":        "saved",
				"scopes":       reg.PreClaimScopes,
				"next": []string{
					"linkbreakers links create --destination https://example.com --name \"My first link\"",
					"linkbreakers links list",
					"linkbreakers --help",
				},
				"limits": "Free plan. See them with `linkbreakers allowances list-allowances`.",
			}
			var started claimStarted
			if err := postAgent(ctx, cfg.AgentBaseURL+"/agent/identity/claim", map[string]string{"claim_token": reg.ClaimToken, "email": email}, &started); err != nil {
				result["claim"] = map[string]any{
					"status":       "not_started",
					"error":        err.Error(),
					"instructions": fmt.Sprintf("The workspace works, but the claim email could not be sent. Run `linkbreakers claim --email %s` before %s or the workspace is deleted.", email, formatDeadline(expires)),
				}
			} else {
				result["claim"] = claimSummary(email, expires)
			}
			return output.PrintJSON(result)
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "Email of the person the workspace is for. Only they receive the link and code that claim it, so ask them for their real address. Required.")
	cmd.Flags().BoolVar(&force, "force", false, "Replace the token already configured with the new workspace's.")
	return cmd
}

func (a *app) newClaimCommand() *cobra.Command {
	var email string
	var wait bool

	cmd := &cobra.Command{
		Use:   "claim",
		Short: "Hand the workspace created by `linkbreakers signup` to a person",
		Long: "Emails the person a link and a 6-digit code to claim the workspace `linkbreakers signup` created. They go only to the inbox, so claiming also verifies the address. " +
			"Each call replaces the previous email's link. With --wait, it waits until the workspace is claimed and then upgrades the saved token to the claimed workspace's permissions.",
		Example: "" +
			"  linkbreakers claim --email you@example.com\n" +
			"  linkbreakers claim --wait\n",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fileCfg, err := config.LoadFileConfig()
			if err != nil {
				return err
			}
			if fileCfg.Claim == nil {
				return fmt.Errorf("no workspace is waiting to be claimed: `linkbreakers claim` follows `linkbreakers signup`")
			}
			cfg, err := a.runtimeConfig()
			if err != nil {
				return err
			}
			claim := fileCfg.Claim
			if email = strings.TrimSpace(email); email != "" {
				claim.Email = email
			}

			if email != "" || !wait {
				ctx, cancel := context.WithTimeout(cmd.Context(), cfg.Timeout)
				var started claimStarted
				err := postAgent(ctx, cfg.AgentBaseURL+"/agent/identity/claim", map[string]string{"claim_token": claim.ClaimToken, "email": claim.Email}, &started)
				cancel()
				if err != nil {
					return fmt.Errorf("could not start the claim: %w", err)
				}
				if err := config.SaveFileConfig(fileCfg); err != nil {
					return err
				}
				if !wait {
					return output.PrintJSON(map[string]any{"claim": claimSummary(claim.Email, claim.ExpiresAt)})
				}
				fmt.Fprintf(os.Stderr, "Waiting for %s to open the link we emailed and enter its code...\n", claim.Email)
			}

			for {
				token, err := pollClaim(cmd.Context(), cfg, claim.ClaimToken)
				if err == nil {
					if err := saveClaimed(fileCfg, token); err != nil {
						return err
					}
					return output.PrintJSON(map[string]any{"claimed": true, "workspace_id": claim.WorkspaceID, "scopes": strings.Fields(token.Scope)})
				}
				var agentErr *agentError
				if !asAgentError(err, &agentErr) || (agentErr.Code != "authorization_pending" && agentErr.Code != "slow_down") {
					return fmt.Errorf("the workspace was not claimed: %w", err)
				}
				select {
				case <-cmd.Context().Done():
					return cmd.Context().Err()
				case <-time.After(5 * time.Second):
				}
			}
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "Email of the person to hand the workspace to. Defaults to the one given at signup.")
	cmd.Flags().BoolVar(&wait, "wait", false, "Wait until the workspace is claimed, then upgrade the saved token.")
	return cmd
}

func claimSummary(email string, expires time.Time) map[string]any {
	return map[string]any{
		"status":      "emailed",
		"email":       email,
		"required_by": expires.UTC().Format(time.RFC3339),
		"instructions": fmt.Sprintf("We emailed %s a link to claim this workspace. Ask them to open it before %s; only they have it. If it did not arrive, check the address and run `linkbreakers claim --email <address>` to send a new one. An unclaimed workspace is deleted with its links and QR codes at that time.",
			email, formatDeadline(expires)),
	}
}

// claimReminder runs after every command while a signed-up workspace waits to
// be claimed. At most once a minute it asks whether the claim happened, and
// upgrades the saved token when it has; otherwise it says on stderr what
// happens if nobody claims it.
func (a *app) claimReminder(ctx context.Context) {
	fileCfg, err := config.LoadFileConfig()
	if err != nil || fileCfg.Claim == nil {
		return
	}
	claim := fileCfg.Claim
	if time.Since(claim.LastCheckedAt) > time.Minute {
		cfg, err := a.runtimeConfig()
		if err != nil {
			return
		}
		pollCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		token, err := pollClaim(pollCtx, cfg, claim.ClaimToken)
		cancel()
		if err == nil {
			if saveClaimed(fileCfg, token) == nil {
				fmt.Fprintf(os.Stderr, "\nThe workspace was claimed by %s. The saved token now has the claimed workspace's permissions.\n", claim.Email)
			}
			return
		}
		var agentErr *agentError
		if asAgentError(err, &agentErr) && (agentErr.Code == "expired_token" || agentErr.Code == "invalid_grant") && time.Now().After(claim.ExpiresAt) {
			fileCfg.Claim = nil
			fileCfg.Token = ""
			_ = config.SaveFileConfig(fileCfg)
			fmt.Fprintln(os.Stderr, "\nThe workspace created by `linkbreakers signup` was never claimed and has been deleted. Run `linkbreakers signup --email <email>` to start again.")
			return
		}
		claim.LastCheckedAt = time.Now()
		_ = config.SaveFileConfig(fileCfg)
	}
	fmt.Fprintf(os.Stderr, "\nUnclaimed workspace: it will be deleted on %s unless %s claims it with the link we emailed. Resend the link with `linkbreakers claim --email %s`.\n",
		formatDeadline(claim.ExpiresAt), claim.Email, claim.Email)
}

func pollClaim(ctx context.Context, cfg config.RuntimeConfig, claimToken string) (agentToken, error) {
	var token agentToken
	err := postAgent(ctx, cfg.AgentBaseURL+"/oauth/token", map[string]string{
		"grant_type":  "urn:workos:agent-auth:grant-type:claim",
		"claim_token": claimToken,
	}, &token)
	return token, err
}

func saveClaimed(fileCfg config.FileConfig, token agentToken) error {
	if token.AccessToken != "" {
		fileCfg.Token = token.AccessToken
	}
	fileCfg.Claim = nil
	return config.SaveFileConfig(fileCfg)
}

func postAgent(ctx context.Context, endpoint string, body any, target any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent())
	req.Header.Set(api.ClientHeader, clientHeader())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		agentErr := &agentError{Status: resp.StatusCode}
		_ = json.Unmarshal(data, agentErr)
		if agentErr.Code == "" {
			agentErr.Code = http.StatusText(resp.StatusCode)
		}
		return agentErr
	}
	return json.Unmarshal(data, target)
}

func asAgentError(err error, target **agentError) bool {
	e, ok := err.(*agentError)
	if ok {
		*target = e
	}
	return ok
}

func formatDeadline(t time.Time) string {
	if t.IsZero() {
		return "24 hours after signup"
	}
	return t.UTC().Format("2006-01-02 15:04 UTC")
}
