package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/linkbreakers-com/linkbreakers-cli/internal/output"
	"github.com/spf13/cobra"
)

// operation is one API operation, read from the OpenAPI document by
// internal/gen/commands. Every API command is built from these, so the CLI
// gains, changes and loses commands exactly as the API does.
type operation struct {
	Group, GroupSummary, Name, Summary, Description, Method, Path string
	Params, Body                                                  []field
	HasBody, BodyRequired, BodyJSONOnly, Auth                     bool
	Deprecated                                                    bool
}

// field is a path, query or header parameter, or a top-level property of the
// request body. Kind is how the flag is read: string, integer, number,
// boolean, strings (repeatable) or json.
type field struct {
	Name, Flag, In, Kind, Description string
	Enum                              []string
	Required                          bool
}

// addOperationCommands adds one command group per API service, and attaches
// to groups that already exist, such as auth.
func (a *app) addOperationCommands(root *cobra.Command) {
	groups := map[string]*cobra.Command{}
	for _, existing := range root.Commands() {
		groups[existing.Name()] = existing
	}
	for i := range generatedOperations {
		op := generatedOperations[i]
		group, ok := groups[op.Group]
		if !ok {
			short := op.GroupSummary
			if short == "" {
				short = fmt.Sprintf("%s API operations", strings.ReplaceAll(op.Group, "-", " "))
			}
			group = &cobra.Command{Use: op.Group, Short: short}
			groups[op.Group] = group
			root.AddCommand(group)
		}
		group.AddCommand(a.newOperationCommand(op))
	}
}

func (a *app) newOperationCommand(op operation) *cobra.Command {
	long := op.Summary
	if op.Description != "" && op.Description != op.Summary {
		long = op.Description
	}
	long += fmt.Sprintf("\n\nAPI: %s %s", op.Method, op.Path)

	cmd := &cobra.Command{
		Use:        op.Name,
		Short:      op.Summary,
		Long:       long,
		Args:       cobra.NoArgs,
		Deprecated: deprecation(op),
	}

	var data string
	if op.HasBody {
		usage := "Request body as JSON. Use @file.json to read a file or - for stdin. Flags set fields on top of it."
		if op.BodyJSONOnly {
			usage = "Request body as JSON (required). Use @file.json to read a file or - for stdin."
		}
		cmd.Flags().StringVar(&data, "data", "", usage)
	}
	for _, f := range append(append([]field{}, op.Params...), op.Body...) {
		addFieldFlag(cmd, f)
	}

	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		client, err := a.client(op.Auth)
		if err != nil {
			return err
		}

		path := op.Path
		query := url.Values{}
		headers := map[string]string{}
		for _, p := range op.Params {
			if !cmd.Flags().Changed(p.Flag) {
				if p.In == "path" {
					return fmt.Errorf("--%s is required", p.Flag)
				}
				continue
			}
			values, err := flagStrings(cmd, p)
			if err != nil {
				return err
			}
			switch p.In {
			case "path":
				path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(values[0]))
			case "query":
				for _, v := range values {
					query.Add(p.Name, v)
				}
			case "header":
				headers[p.Name] = values[0]
			}
		}

		var body []byte
		if op.HasBody {
			body, err = buildBody(cmd, op, data)
			if err != nil {
				return err
			}
		}

		status, response, err := client.Do(cmd.Context(), op.Method, path, query, headers, body)
		if err != nil {
			return apiError(status, response, err)
		}
		return output.PrintBody(response)
	}
	return cmd
}

func deprecation(op operation) string {
	if op.Deprecated {
		return "this API operation is deprecated and may be removed"
	}
	return ""
}

func addFieldFlag(cmd *cobra.Command, f field) {
	usage := f.Description
	if len(f.Enum) > 0 {
		usage = strings.TrimSpace(usage + " One of: " + strings.Join(f.Enum, ", ") + ".")
	}
	if f.Kind == "json" {
		usage = strings.TrimSpace(usage + " JSON value.")
	}
	if f.Required && !strings.Contains(strings.ToLower(usage), "required") {
		usage = strings.TrimSpace(usage + " Required.")
	}
	flags := cmd.Flags()
	switch f.Kind {
	case "boolean":
		flags.Bool(f.Flag, false, usage)
	case "integer":
		flags.Int64(f.Flag, 0, usage)
	case "number":
		flags.Float64(f.Flag, 0, usage)
	case "strings":
		flags.StringArray(f.Flag, nil, usage+" Repeat the flag for several values.")
	default:
		flags.String(f.Flag, "", usage)
	}
}

// flagStrings reads a parameter as the strings it is sent as.
func flagStrings(cmd *cobra.Command, f field) ([]string, error) {
	flags := cmd.Flags()
	switch f.Kind {
	case "boolean":
		v, err := flags.GetBool(f.Flag)
		return []string{strconv.FormatBool(v)}, err
	case "integer":
		v, err := flags.GetInt64(f.Flag)
		return []string{strconv.FormatInt(v, 10)}, err
	case "number":
		v, err := flags.GetFloat64(f.Flag)
		return []string{strconv.FormatFloat(v, 'f', -1, 64)}, err
	case "strings":
		return flags.GetStringArray(f.Flag)
	default:
		v, err := flags.GetString(f.Flag)
		if err == nil && len(f.Enum) > 0 && !contains(f.Enum, v) {
			return nil, fmt.Errorf("--%s must be one of %s", f.Flag, strings.Join(f.Enum, ", "))
		}
		return []string{v}, err
	}
}

// buildBody starts from --data and sets every field passed as a flag on top.
func buildBody(cmd *cobra.Command, op operation, data string) ([]byte, error) {
	raw, err := readData(data)
	if err != nil {
		return nil, err
	}
	if op.BodyJSONOnly {
		if raw == nil {
			return nil, fmt.Errorf("--data is required: this operation takes a JSON body")
		}
		return raw, nil
	}

	body := map[string]any{}
	if raw != nil {
		if err := json.Unmarshal(raw, &body); err != nil {
			return nil, fmt.Errorf("--data must be a JSON object: %w", err)
		}
	}
	flags := cmd.Flags()
	for _, f := range op.Body {
		if !flags.Changed(f.Flag) {
			continue
		}
		var value any
		switch f.Kind {
		case "boolean":
			value, err = flags.GetBool(f.Flag)
		case "integer":
			value, err = flags.GetInt64(f.Flag)
		case "number":
			value, err = flags.GetFloat64(f.Flag)
		case "strings":
			var values []string
			values, err = flags.GetStringArray(f.Flag)
			for _, v := range values {
				if len(f.Enum) > 0 && !contains(f.Enum, v) {
					return nil, fmt.Errorf("--%s must be one of %s", f.Flag, strings.Join(f.Enum, ", "))
				}
			}
			value = values
		case "json":
			var s string
			s, err = flags.GetString(f.Flag)
			if err == nil {
				if jsonErr := json.Unmarshal([]byte(s), &value); jsonErr != nil {
					return nil, fmt.Errorf("--%s must be JSON: %w", f.Flag, jsonErr)
				}
			}
		default:
			var s string
			s, err = flags.GetString(f.Flag)
			if err == nil && len(f.Enum) > 0 && !contains(f.Enum, s) {
				return nil, fmt.Errorf("--%s must be one of %s", f.Flag, strings.Join(f.Enum, ", "))
			}
			value = s
		}
		if err != nil {
			return nil, err
		}
		body[f.Name] = value
	}
	if len(body) == 0 && !op.BodyRequired {
		return nil, nil
	}
	return json.Marshal(body)
}

// readData never reads stdin implicitly: agent and CI shells keep a non-TTY
// stdin open forever, which would make every command block. Use --data -.
func readData(data string) ([]byte, error) {
	switch {
	case data == "":
		return nil, nil
	case data == "-":
		return readAll(os.Stdin)
	case strings.HasPrefix(data, "@"):
		return os.ReadFile(strings.TrimPrefix(data, "@"))
	default:
		return []byte(data), nil
	}
}

// apiError keeps the API's own explanation: agents act on the message and
// the code, not on the status line.
func apiError(status int, response []byte, err error) error {
	var body struct {
		CodeType string            `json:"codeType"`
		Message  string            `json:"message"`
		Metadata map[string]string `json:"metadata"`
	}
	if status == 0 || json.Unmarshal(response, &body) != nil || body.Message == "" {
		return err
	}
	msg := fmt.Sprintf("%d %s: %s", status, body.CodeType, body.Message)
	if len(body.Metadata) > 0 {
		keys := make([]string, 0, len(body.Metadata))
		for k := range body.Metadata {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, k+"="+body.Metadata[k])
		}
		msg += " (" + strings.Join(parts, ", ") + ")"
	}
	return fmt.Errorf("%s", msg)
}

func contains(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}
