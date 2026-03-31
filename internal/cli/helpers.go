package cli

import (
	"encoding/json"
	"fmt"
	"strings"
)

func parseJSON(data []byte, target any) error {
	return json.Unmarshal(data, target)
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func parseStringMap(pairs []string) (map[string]string, error) {
	out := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid key/value %q, expected key=value", pair)
		}
		out[parts[0]] = parts[1]
	}
	return out, nil
}
