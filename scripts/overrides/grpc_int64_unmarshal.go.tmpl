// grpc_int64_unmarshal.go provides custom UnmarshalJSON methods for generated
// models whose proto definitions use int64 fields. gRPC-gateway serialises
// int64/uint64 values as JSON strings per the proto3 JSON spec, but the
// OpenAPI-generated Go structs expect JSON numbers. These overrides accept
// both representations so the CLI works against the real API.
//
// This file is copied into internal/client/generated/ by generate-client.sh
// after code-generation runs, so it survives re-generation.

package linkbreakers

import (
	"encoding/json"
	"strconv"
)

// unmarshalFlexibleInt64 parses a json.RawMessage that may be a JSON number
// or a quoted string (as produced by gRPC-gateway for proto int64 fields)
// and returns the int64 value. Returns nil when raw is empty or JSON null.
func unmarshalFlexibleInt64(raw json.RawMessage) (*int64, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var n int64
	if err := json.Unmarshal(raw, &n); err == nil {
		return &n, nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// --- ListLinksResponse ---

func (o *ListLinksResponse) UnmarshalJSON(data []byte) error {
	type Alias ListLinksResponse
	aux := &struct {
		TotalCount json.RawMessage `json:"totalCount,omitempty"`
		*Alias
	}{Alias: (*Alias)(o)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	v, err := unmarshalFlexibleInt64(aux.TotalCount)
	if err != nil {
		return err
	}
	o.TotalCount = v
	return nil
}

// --- ListVisitorsJsonResponse ---

func (o *ListVisitorsJsonResponse) UnmarshalJSON(data []byte) error {
	type Alias ListVisitorsJsonResponse
	aux := &struct {
		TotalCount json.RawMessage `json:"totalCount,omitempty"`
		*Alias
	}{Alias: (*Alias)(o)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	v, err := unmarshalFlexibleInt64(aux.TotalCount)
	if err != nil {
		return err
	}
	o.TotalCount = v
	return nil
}

// --- ListEventsJsonResponse ---

func (o *ListEventsJsonResponse) UnmarshalJSON(data []byte) error {
	type Alias ListEventsJsonResponse
	aux := &struct {
		TotalCount json.RawMessage `json:"totalCount,omitempty"`
		*Alias
	}{Alias: (*Alias)(o)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	v, err := unmarshalFlexibleInt64(aux.TotalCount)
	if err != nil {
		return err
	}
	o.TotalCount = v
	return nil
}

// --- Device ---

func (o *Device) UnmarshalJSON(data []byte) error {
	type Alias Device
	aux := &struct {
		Asn json.RawMessage `json:"asn,omitempty"`
		*Alias
	}{Alias: (*Alias)(o)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	v, err := unmarshalFlexibleInt64(aux.Asn)
	if err != nil {
		return err
	}
	o.Asn = v
	return nil
}
