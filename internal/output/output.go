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

// ReadBody never reads stdin implicitly: agent and CI shells keep a non-TTY
// stdin open forever, which made every bodyless call block. Use --body-file -.
func ReadBody(bodyArg, bodyFile string) ([]byte, error) {
	return readBody(bodyArg, bodyFile, os.Stdin)
}

func readBody(bodyArg, bodyFile string, stdin io.Reader) ([]byte, error) {
	switch {
	case bodyArg != "" && bodyFile != "":
		return nil, fmt.Errorf("use either --body or --body-file, not both")
	case bodyArg != "":
		return []byte(bodyArg), nil
	case bodyFile == "-":
		return io.ReadAll(stdin)
	case bodyFile != "":
		return os.ReadFile(bodyFile)
	default:
		return nil, nil
	}
}
