package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/config"
)

func TestNewClientAlwaysHasTimeout(t *testing.T) {
	c := NewClient(config.RuntimeConfig{BaseURL: "http://example.invalid"})
	if c.httpClient.Timeout <= 0 {
		t.Fatal("expected a default HTTP timeout")
	}
}

func TestRawRequestTimesOut(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
	}))
	defer srv.Close()
	defer close(release)

	c := NewClient(config.RuntimeConfig{BaseURL: srv.URL, Token: "tok", Timeout: 100 * time.Millisecond})
	start := time.Now()
	_, _, err := c.RawRequest(context.Background(), http.MethodGet, "/v1/links", nil, nil, nil)
	if err == nil {
		t.Fatal("expected a timeout error")
	}
	if time.Since(start) > 2*time.Second {
		t.Fatalf("request took %s, timeout not applied", time.Since(start))
	}
}
