package cli

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/config"
)

func isolateConfig(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("LINKBREAKERS_TOKEN", "")
	t.Setenv("LINKBREAKERS_BASE_URL", "")
	t.Setenv("LINKBREAKERS_OUTPUT", "")
}

// holdStdinOpen mimics agent and CI shells, whose stdin is a pipe that never closes.
func holdStdinOpen(t *testing.T) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdin
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = orig
		_ = w.Close()
		_ = r.Close()
	})
}

func runWithDeadline(t *testing.T, args ...string) error {
	t.Helper()
	cmd := newRootCommand()
	cmd.SetArgs(args)
	done := make(chan error, 1)
	go func() { done <- cmd.Execute() }()
	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		t.Fatalf("linkbreakers %v did not return within 5s", args)
		return nil
	}
}

func TestRawWithoutTokenFailsFast(t *testing.T) {
	isolateConfig(t)
	holdStdinOpen(t)

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
	}))
	defer srv.Close()

	err := runWithDeadline(t, "raw", "GET", "/v1/links", "--base-url", srv.URL)
	if !errors.Is(err, config.ErrMissingToken) {
		t.Fatalf("expected ErrMissingToken, got %v", err)
	}
	if hits.Load() != 0 {
		t.Fatalf("expected no API call without a token, got %d", hits.Load())
	}
}

func TestRawGetDoesNotWaitOnOpenStdin(t *testing.T) {
	isolateConfig(t)
	holdStdinOpen(t)

	var gotAuth atomic.Value
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth.Store(r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"links":[]}`))
	}))
	defer srv.Close()

	if err := runWithDeadline(t, "raw", "GET", "/v1/links", "--base-url", srv.URL, "--token", "tok"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth.Load() != "Bearer tok" {
		t.Fatalf("expected bearer token, got %v", gotAuth.Load())
	}
}
