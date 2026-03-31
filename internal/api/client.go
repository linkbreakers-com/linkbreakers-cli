package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/client/generated"
	"github.com/linkbreakers-com/linkbreakers-cli/internal/config"
)

type Client struct {
	cfg        config.RuntimeConfig
	generated  *linkbreakers.APIClient
	httpClient *http.Client
}

func NewClient(cfg config.RuntimeConfig) *Client {
	genCfg := linkbreakers.NewConfiguration()
	genCfg.Servers = linkbreakers.ServerConfigurations{
		{URL: cfg.BaseURL, Description: "resolved by linkbreakers CLI"},
	}
	genCfg.HTTPClient = &http.Client{Timeout: cfg.Timeout}
	genCfg.UserAgent = "linkbreakers-cli"

	return &Client{
		cfg:        cfg,
		generated:  linkbreakers.NewAPIClient(genCfg),
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

func (c *Client) Generated() *linkbreakers.APIClient {
	return c.generated
}

func (c *Client) Context() context.Context {
	return context.WithValue(context.Background(), linkbreakers.ContextAccessToken, c.cfg.Token)
}

func (c *Client) RawRequest(ctx context.Context, method, path string, headers map[string]string, query map[string]string, body []byte) (int, []byte, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	base, err := url.Parse(c.cfg.BaseURL)
	if err != nil {
		return 0, nil, fmt.Errorf("parse base url: %w", err)
	}

	rel, err := url.Parse(path)
	if err != nil {
		return 0, nil, fmt.Errorf("parse path: %w", err)
	}

	reqURL := base.ResolveReference(rel)
	values := reqURL.Query()
	for k, v := range query {
		values.Set(k, v)
	}
	reqURL.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bytes.NewReader(body))
	if err != nil {
		return 0, nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if c.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return resp.StatusCode, respBody, fmt.Errorf("api returned %s", resp.Status)
	}

	return resp.StatusCode, respBody, nil
}

func NormalizeJSON(body []byte) ([]byte, error) {
	if len(bytes.TrimSpace(body)) == 0 {
		return []byte("{}"), nil
	}

	var v any
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, err
	}

	return json.Marshal(v)
}
