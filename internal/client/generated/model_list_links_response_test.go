package linkbreakers

import (
	"encoding/json"
	"testing"
)

func TestListLinksResponseUnmarshalJSONAcceptsStringTotalCount(t *testing.T) {
	var resp ListLinksResponse

	err := json.Unmarshal([]byte(`{
		"links": [],
		"nextPageToken": "cursor_123",
		"totalCount": "42"
	}`), &resp)
	if err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if resp.TotalCount == nil {
		t.Fatal("expected totalCount to be set")
	}

	if got := *resp.TotalCount; got != 42 {
		t.Fatalf("expected totalCount 42, got %d", got)
	}

	if got := resp.GetNextPageToken(); got != "cursor_123" {
		t.Fatalf("expected nextPageToken cursor_123, got %q", got)
	}
}

func TestListLinksResponseUnmarshalJSONAcceptsNumericTotalCount(t *testing.T) {
	var resp ListLinksResponse

	err := json.Unmarshal([]byte(`{
		"links": [],
		"totalCount": 7
	}`), &resp)
	if err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if resp.TotalCount == nil {
		t.Fatal("expected totalCount to be set")
	}

	if got := *resp.TotalCount; got != 7 {
		t.Fatalf("expected totalCount 7, got %d", got)
	}
}
