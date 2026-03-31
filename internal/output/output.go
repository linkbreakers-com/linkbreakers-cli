package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

func PrintJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(os.Stdout, string(data))
	return err
}

func PrintTable(headers []string, rows [][]string) error {
	w := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, strings.Join(headers, "\t")); err != nil {
		return err
	}
	for _, row := range rows {
		if _, err := fmt.Fprintln(w, strings.Join(row, "\t")); err != nil {
			return err
		}
	}
	return w.Flush()
}

func ReadBody(bodyArg, bodyFile string) ([]byte, error) {
	switch {
	case bodyArg != "" && bodyFile != "":
		return nil, fmt.Errorf("use either --body or --body-file, not both")
	case bodyArg != "":
		return []byte(bodyArg), nil
	case bodyFile != "":
		return os.ReadFile(bodyFile)
	default:
		stat, err := os.Stdin.Stat()
		if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
			return io.ReadAll(os.Stdin)
		}
		return nil, nil
	}
}
