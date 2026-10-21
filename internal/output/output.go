package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func PrintJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(os.Stdout, string(data))
	return err
}

// PrintBody prints an API response, indented when it is JSON and as it came
// otherwise. An empty body prints {} so the output is always JSON.
func PrintBody(body []byte) error {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		trimmed = []byte("{}")
	}
	var indented bytes.Buffer
	if err := json.Indent(&indented, trimmed, "", "  "); err != nil {
		_, err = fmt.Fprintln(os.Stdout, string(trimmed))
		return err
	}
	_, err := fmt.Fprintln(os.Stdout, indented.String())
	return err
}
