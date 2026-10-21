package main

import (
	"encoding/json"
	"testing"
)

func TestKebab(t *testing.T) {
	for in, want := range map[string]string{
		"Create":              "create",
		"ResolveQrcodeImages": "resolve-qrcode-images",
		"GetMCPServer":        "get-mcp-server",
		"Idempotency-Key":     "idempotency-key",
		"qrcodeEnabled":       "qrcode-enabled",
		"linkId":              "link-id",
		"QrcodeTemplates":     "qrcode-templates",
	} {
		if got := kebab(in); got != want {
			t.Errorf("kebab(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBuildReadsParametersAndBody(t *testing.T) {
	var doc document
	if err := json.Unmarshal([]byte(`{
	  "components": {"schemas": {
	    "Create": {"required": ["destination"], "properties": {
	      "destination": {"type": "string", "description": "Where it goes"},
	      "tags": {"type": "array", "items": {"type": "string"}},
	      "qrcodeEnabled": {"type": "boolean"},
	      "metadata": {"type": "object", "additionalProperties": {"type": "string"}},
	      "view": {"$ref": "#/components/schemas/View"},
	      "data": {"type": "string"}
	    }},
	    "View": {"enum": ["A", "B"]}
	  }}
	}`), &doc); err != nil {
		t.Fatal(err)
	}
	op, err := build(doc, "/v1/links/{linkId}", "POST", operationDoc{
		OperationID: "LinksService_Create",
		Summary:     "Create a link.",
		Security:    []any{map[string]any{}},
		Parameters:  []parameter{{Name: "linkId", In: "path", Schema: schema{Type: json.RawMessage(`"string"`)}}},
		RequestBody: &requestBody{Required: true, Content: map[string]struct {
			Schema schema `json:"schema"`
		}{"application/json": {Schema: schema{Ref: "#/components/schemas/Create"}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if op.Group != "links" || op.Name != "create" || op.Summary != "Create a link" || !op.Auth || !op.HasBody {
		t.Fatalf("operation = %+v", op)
	}
	if len(op.Params) != 1 || op.Params[0].Flag != "link-id" || !op.Params[0].Required {
		t.Fatalf("params = %+v", op.Params)
	}
	kinds := map[string]field{}
	for _, f := range op.Body {
		kinds[f.Name] = f
	}
	for name, want := range map[string]string{"destination": "string", "tags": "strings", "qrcodeEnabled": "boolean", "metadata": "json", "view": "string"} {
		if kinds[name].Kind != want {
			t.Errorf("%s kind = %q, want %q", name, kinds[name].Kind, want)
		}
	}
	if !kinds["destination"].Required || len(kinds["view"].Enum) != 2 {
		t.Errorf("destination required = %v, view enum = %v", kinds["destination"].Required, kinds["view"].Enum)
	}
	// A body field named like a global flag does not shadow it.
	if kinds["data"].Flag != "body-data" {
		t.Errorf("data flag = %q, want body-data", kinds["data"].Flag)
	}
}
