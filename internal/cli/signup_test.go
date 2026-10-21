package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/config"
)

// fakeAgentHost answers self-registration, the claim and the claim grant the
// way the MCP host does; claimed flips the grant from pending to a token.
func fakeAgentHost(t *testing.T, claimed *bool) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/agent/identity":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "lb_agent_token", "workspace_id": "ws-1", "claim_token": "clm_1",
				"claim_token_expires": time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
				"pre_claim_scopes":    []string{"links:write"},
			})
		case "/agent/identity/claim":
			if body["email"] == "" || body["claim_token"] != "clm_1" {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_request"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"claim_attempt": map[string]any{
				"user_code": "123456", "verification_uri": "https://app.example/agent/claim?claim_attempt_token=cat_1", "interval": 5,
			}})
		case "/oauth/token":
			if !*claimed {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "authorization_pending"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"access_token": "lb_claimed_token", "scope": "links:write media:write"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func useTempConfig(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("LINKBREAKERS_TOKEN", "")
}

func runCLI(t *testing.T, args ...string) error {
	t.Helper()
	a := &app{}
	root := a.newRootCommand()
	root.SetArgs(args)
	root.SetContext(context.Background())
	return root.Execute()
}

func TestSignupSavesTokenAndClaimUpgradesIt(t *testing.T) {
	useTempConfig(t)
	claimed := false
	host := fakeAgentHost(t, &claimed)
	defer host.Close()
	t.Setenv("LINKBREAKERS_AGENT_BASE_URL", host.URL)

	if err := runCLI(t, "signup"); err == nil || !strings.Contains(err.Error(), "--email is required") {
		t.Fatalf("signup without email = %v", err)
	}
	if err := runCLI(t, "signup", "--email", "you@example.com"); err != nil {
		t.Fatalf("signup: %v", err)
	}
	cfg, _ := config.LoadFileConfig()
	if cfg.Token != "lb_agent_token" || cfg.Claim == nil || cfg.Claim.Email != "you@example.com" || cfg.Claim.ClaimToken != "clm_1" {
		t.Fatalf("config after signup = %+v", cfg)
	}
	if err := runCLI(t, "signup", "--email", "you@example.com"); err == nil {
		t.Fatal("a second signup replaced the token without --force")
	}

	// The reminder learns about the claim on its own and upgrades the token.
	claimed = true
	(&app{}).claimReminder(context.Background())
	cfg, _ = config.LoadFileConfig()
	if cfg.Token != "lb_claimed_token" || cfg.Claim != nil {
		t.Fatalf("config after the claim = %+v", cfg)
	}
}

func TestClaimWaitUpgradesTheToken(t *testing.T) {
	useTempConfig(t)
	claimed := true
	host := fakeAgentHost(t, &claimed)
	defer host.Close()
	t.Setenv("LINKBREAKERS_AGENT_BASE_URL", host.URL)

	if err := config.SaveFileConfig(config.FileConfig{Token: "lb_agent_token", Claim: &config.Claim{ClaimToken: "clm_1", Email: "you@example.com", ExpiresAt: time.Now().Add(time.Hour)}}); err != nil {
		t.Fatal(err)
	}
	if err := runCLI(t, "claim", "--wait"); err != nil {
		t.Fatalf("claim --wait: %v", err)
	}
	cfg, _ := config.LoadFileConfig()
	if cfg.Token != "lb_claimed_token" || cfg.Claim != nil {
		t.Fatalf("config after claim --wait = %+v", cfg)
	}
}

func TestMissingTokenOffersBothRoutes(t *testing.T) {
	useTempConfig(t)
	err := runCLI(t, "links", "list")
	if err == nil || !strings.Contains(err.Error(), "auth set-token") || !strings.Contains(err.Error(), "linkbreakers signup") {
		t.Fatalf("missing token error = %v", err)
	}
}
