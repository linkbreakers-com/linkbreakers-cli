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
	defaultBaseURL = "https://api.linkbreakers.com"
	defaultOutput  = "json"
)

type FileConfig struct {
	Token   string `json:"token,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
	Output  string `json:"output,omitempty"`
}

type RuntimeConfig struct {
	Token   string
	BaseURL string
	Output  string
	Timeout time.Duration
}

func DefaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		BaseURL: defaultBaseURL,
		Output:  defaultOutput,
		Timeout: 30 * time.Second,
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
	switch cfg.Output {
	case "json", "table":
	default:
		return RuntimeConfig{}, fmt.Errorf("unsupported output %q, expected json or table", cfg.Output)
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}

	return cfg, nil
}

func (c RuntimeConfig) RequireToken() error {
	if strings.TrimSpace(c.Token) == "" {
		return errors.New("missing API token: set LINKBREAKERS_TOKEN or run `linkbreakers auth set-token`")
	}
	return nil
}
