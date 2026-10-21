package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/config"
)

type Client struct {
	cfg        config.RuntimeConfig
	httpClient *http.Client
	userAgent  string
}

func NewClient(cfg config.RuntimeConfig, userAgent string) *Client {
	if cfg.Timeout <= 0 {
		cfg.Timeout = config.DefaultRuntimeConfig().Timeout
	}
	return &Client{cfg: cfg, httpClient: &http.Client{Timeout: cfg.Timeout}, userAgent: userAgent}
}

// Do sends one API request and returns the status and body. A status of 400
// or above is an error carrying the body, which holds the API's explanation.
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, headers map[string]string, body []byte) (int, []byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	base, err := url.Parse(c.cfg.BaseURL)
	if err != nil {
		return 0, nil, fmt.Errorf("parse base url: %w", err)
	}
	reqURL := base.ResolveReference(&url.URL{Path: path})
	if len(query) > 0 {
		reqURL.RawQuery = query.Encode()
	}

	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), reader)
	if err != nil {
		return 0, nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
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
