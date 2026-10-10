package output

import (
	"strings"
	"testing"
)

func TestReadBodyIgnoresStdinUnlessDash(t *testing.T) {
	body, err := readBody("", "", strings.NewReader(`{"a":1}`))
	if err != nil || body != nil {
		t.Fatalf("expected no body, got %q, %v", body, err)
	}

	body, err = readBody("", "-", strings.NewReader(`{"a":1}`))
	if err != nil || string(body) != `{"a":1}` {
		t.Fatalf("expected stdin body, got %q, %v", body, err)
	}
}

func TestReadBodyRejectsBothSources(t *testing.T) {
	if _, err := readBody("{}", "x.json", nil); err == nil {
		t.Fatal("expected error when both --body and --body-file are set")
	}
}
