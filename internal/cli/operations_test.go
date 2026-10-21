package cli

import (
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
)

func TestEveryOperationBecomesACommand(t *testing.T) {
	root := (&app{}).newRootCommand()
	for _, op := range generatedOperations {
		cmd, _, err := root.Find([]string{op.Group, op.Name})
		if err != nil || cmd.Name() != op.Name {
			t.Errorf("no command for %s %s", op.Group, op.Name)
		}
	}
	if cmd, _, err := root.Find([]string{"links", "create"}); err != nil || cmd.Flags().Lookup("destination") == nil {
		t.Fatal("links create has no --destination flag")
	}
}

func TestBuildBodyMergesDataAndFlags(t *testing.T) {
	op := operation{HasBody: true, Body: []field{
		{Name: "destination", Flag: "destination", Kind: "string"},
		{Name: "tags", Flag: "tags", Kind: "strings"},
		{Name: "qrcodeEnabled", Flag: "qrcode-enabled", Kind: "boolean"},
		{Name: "metadata", Flag: "metadata", Kind: "json"},
		{Name: "view", Flag: "view", Kind: "string", Enum: []string{"A", "B"}},
	}}
	cmd := &cobra.Command{}
	for _, f := range op.Body {
		addFieldFlag(cmd, f)
	}
	if err := cmd.ParseFlags([]string{"--destination", "https://example.com", "--tags", "a", "--tags", "b", "--qrcode-enabled=false", "--metadata", `{"k":"v"}`}); err != nil {
		t.Fatal(err)
	}
	raw, err := buildBody(cmd, op, `{"name":"from data","destination":"https://overridden.example"}`)
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	if body["destination"] != "https://example.com" || body["name"] != "from data" || body["qrcodeEnabled"] != false {
		t.Fatalf("body = %v", body)
	}
	if tags, _ := body["tags"].([]any); len(tags) != 2 {
		t.Fatalf("tags = %v", body["tags"])
	}
	if meta, _ := body["metadata"].(map[string]any); meta["k"] != "v" {
		t.Fatalf("metadata = %v", body["metadata"])
	}

	bad := &cobra.Command{}
	for _, f := range op.Body {
		addFieldFlag(bad, f)
	}
	_ = bad.ParseFlags([]string{"--view", "C"})
	if _, err := buildBody(bad, op, ""); err == nil {
		t.Fatal("a value outside the enum was accepted")
	}
}

func TestAPIErrorKeepsTheExplanation(t *testing.T) {
	err := apiError(403, []byte(`{"codeType":"WORKSPACE_UNVERIFIED","message":"This workspace is paused.","metadata":{"pausedAt":"2026-10-28T09:00:00Z"}}`), nil)
	if err == nil || err.Error() != "403 WORKSPACE_UNVERIFIED: This workspace is paused. (pausedAt=2026-10-28T09:00:00Z)" {
		t.Fatalf("err = %v", err)
	}
}
