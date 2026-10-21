package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultBaseURL      = "https://api.linkbreakers.com"
	defaultAgentBaseURL = "https://mcp.linkbreakers.com"
	defaultOutput       = "json"
)

type FileConfig struct {
	Token   string `json:"token,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
	Output  string `json:"output,omitempty"`
	// AgentBaseURL is where `linkbreakers signup` registers a workspace.
	AgentBaseURL string `json:"agent_base_url,omitempty"`
	// Claim is set while a workspace created by `linkbreakers signup` waits
	// for its person to claim it, and cleared once they have.
	Claim *Claim `json:"claim,omitempty"`
}

// Claim is what the CLI keeps about a workspace it signed up for until a
// person claims it.
type Claim struct {
	WorkspaceID   string    `json:"workspace_id"`
	Email         string    `json:"email"`
	ClaimToken    string    `json:"claim_token"`
	ExpiresAt     time.Time `json:"expires_at"`
	LastCheckedAt time.Time `json:"last_checked_at,omitempty"`
}

type RuntimeConfig struct {
	Token        string
	BaseURL      string
	AgentBaseURL string
	Output       string
	Timeout      time.Duration
}

func DefaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		BaseURL:      defaultBaseURL,
		AgentBaseURL: defaultAgentBaseURL,
		Output:       defaultOutput,
		Timeout:      30 * time.Second,
	}
}

func ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve config dir: %w", err)
	}

	return filepath.Join(dir, "linkbreakers", "config.json"), nil
}

func LoadFileConfig() (FileConfig, error) {
	path, err := ConfigPath()
	if err != nil {
		return FileConfig{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return FileConfig{}, nil
		}
		return FileConfig{}, fmt.Errorf("read config: %w", err)
	}

	var cfg FileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return FileConfig{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}

func SaveFileConfig(cfg FileConfig) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize config: %w", err)
	}

	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

func DeleteFileConfig() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove config: %w", err)
	}

	return nil
}

func Resolve(tokenFlag, baseURLFlag, outputFlag string) (RuntimeConfig, error) {
	fileCfg, err := LoadFileConfig()
	if err != nil {
		return RuntimeConfig{}, err
	}

	cfg := DefaultRuntimeConfig()

	if fileCfg.Token != "" {
		cfg.Token = strings.TrimSpace(fileCfg.Token)
	}
	if fileCfg.BaseURL != "" {
		cfg.BaseURL = strings.TrimRight(strings.TrimSpace(fileCfg.BaseURL), "/")
	}
	if fileCfg.Output != "" {
		cfg.Output = strings.TrimSpace(fileCfg.Output)
	}
	if fileCfg.AgentBaseURL != "" {
		cfg.AgentBaseURL = strings.TrimRight(strings.TrimSpace(fileCfg.AgentBaseURL), "/")
	}
	if env := strings.TrimSpace(os.Getenv("LINKBREAKERS_AGENT_BASE_URL")); env != "" {
		cfg.AgentBaseURL = strings.TrimRight(env, "/")
	}

	if env := strings.TrimSpace(os.Getenv("LINKBREAKERS_TOKEN")); env != "" {
		cfg.Token = env
	}
	if env := strings.TrimSpace(os.Getenv("LINKBREAKERS_BASE_URL")); env != "" {
		cfg.BaseURL = strings.TrimRight(env, "/")
	}
	if env := strings.TrimSpace(os.Getenv("LINKBREAKERS_OUTPUT")); env != "" {
		cfg.Output = env
	}

	if tokenFlag != "" {
		cfg.Token = strings.TrimSpace(tokenFlag)
	}
	if baseURLFlag != "" {
		cfg.BaseURL = strings.TrimRight(strings.TrimSpace(baseURLFlag), "/")
	}
	if outputFlag != "" {
		cfg.Output = strings.TrimSpace(outputFlag)
	}

	cfg.Output = strings.ToLower(cfg.Output)
	if cfg.Output != "json" {
		return RuntimeConfig{}, fmt.Errorf("unsupported output %q, expected json", cfg.Output)
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}

	return cfg, nil
}

var ErrMissingToken = errors.New("no API token configured.\n" +
	"  Already have a Linkbreakers account? Create a token at https://app.linkbreakers.com/workspace/dashboard/api-tokens, then run `linkbreakers auth set-token --token <api-token>` (or set LINKBREAKERS_TOKEN).\n" +
	"  New to Linkbreakers? Run `linkbreakers signup --email <your email>` to get a workspace and a token right away.")

func (c RuntimeConfig) RequireToken() error {
	if strings.TrimSpace(c.Token) == "" {
		return ErrMissingToken
	}
	return nil
}
