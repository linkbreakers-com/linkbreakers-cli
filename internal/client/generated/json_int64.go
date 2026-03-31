package linkbreakers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func decodeFlexibleInt64(data json.RawMessage) (*int64, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}

	var asInt int64
	if err := json.Unmarshal(trimmed, &asInt); err == nil {
		return &asInt, nil
	}

	var asString string
	if err := json.Unmarshal(trimmed, &asString); err == nil {
		parsed, err := strconv.ParseInt(asString, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse int64 from string %q: %w", asString, err)
		}
		return &parsed, nil
	}

	return nil, fmt.Errorf("unsupported int64 JSON value: %s", string(trimmed))
}
