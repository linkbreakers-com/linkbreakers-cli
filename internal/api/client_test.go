package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/config"
)

func TestNewClientAlwaysHasTimeout(t *testing.T) {
	c := NewClient(config.RuntimeConfig{BaseURL: "http://example.invalid"}, "test")
	if c.httpClient.Timeout <= 0 {
		t.Fatal("expected a default HTTP timeout")
	}
}

func TestDoTimesOut(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
	}))
	defer srv.Close()
	defer close(release)

	c := NewClient(config.RuntimeConfig{BaseURL: srv.URL, Token: "tok", Timeout: 100 * time.Millisecond}, "test")
	start := time.Now()
	if _, _, err := c.Do(context.Background(), http.MethodGet, "/v1/links", nil, nil, nil); err == nil {
		t.Fatal("expected a timeout error")
	}
	if time.Since(start) > 2*time.Second {
		t.Fatalf("request took %s, timeout not applied", time.Since(start))
	}
}

func TestDoSendsTokenQueryAndBody(t *testing.T) {
	var got *http.Request
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := NewClient(config.RuntimeConfig{BaseURL: srv.URL, Token: "tok"}, "test")
	status, body, err := c.Do(context.Background(), http.MethodPost, "/v1/links", url.Values{"tags": {"a", "b"}}, nil, []byte(`{"destination":"https://example.com"}`))
	if err != nil || status != 200 || string(body) != `{"ok":true}` {
		t.Fatalf("Do = %d %s %v", status, body, err)
	}
	if got.Header.Get("Authorization") != "Bearer tok" || got.URL.RawQuery != "tags=a&tags=b" || gotBody != `{"destination":"https://example.com"}` {
		t.Fatalf("request = %s %s %q", got.Header.Get("Authorization"), got.URL.RawQuery, gotBody)
	}
}
