// Package cliformat provides a shared --format flag and output writer
// for Go CLI tools. Supported formats: name, yaml, json.
//
// Usage:
//
//	cliformat.AddFlag(cmd, "name")         // register --format on a cobra command
//	format := cliformat.GetFormat(cmd)     // read the resolved value
//	cliformat.WriteList(w, format, items, func(d *MyType) string { return d.Name })
package cliformat

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// ValidFormats is the set of formats supported by WriteList.
var ValidFormats = []string{"name", "yaml", "json"}

// AddFlag registers a --format flag on cmd with the given default value.
// The flag description lists all ValidFormats.
func AddFlag(cmd *cobra.Command, defaultFormat string) {
	cmd.Flags().String("format", defaultFormat,
		fmt.Sprintf("output format: %s", joinFormats(ValidFormats)))
}

// GetFormat returns the resolved --format flag value from cmd (lower-cased).
// Returns an empty string if the flag was not registered.
func GetFormat(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("format")
	return v
}

// WriteList renders items to w in the requested format.
//
//   - "name" (or ""): one name per line, using nameOf to extract each name.
//   - "yaml":         YAML document stream, one --- block per item.
//   - "json":         JSON array.
//
// T may be any type — struct, map, pointer, etc.
func WriteList[T any](w io.Writer, format string, items []T, nameOf func(T) string) error {
	switch format {
	case "name", "":
		for _, item := range items {
			if _, err := fmt.Fprintln(w, nameOf(item)); err != nil {
				return err
			}
		}
		return nil

	case "yaml":
		enc := yaml.NewEncoder(w)
		enc.SetIndent(2)
		for _, item := range items {
			if err := enc.Encode(item); err != nil {
				return fmt.Errorf("yaml encode: %w", err)
			}
		}
		return enc.Close()

	case "json":
		out, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("json encode: %w", err)
		}
		_, err = fmt.Fprintf(w, "%s\n", out)
		return err

	default:
		return fmt.Errorf("unknown format %q — valid formats: %s", format, joinFormats(ValidFormats))
	}
}

func joinFormats(formats []string) string {
	out := ""
	for i, f := range formats {
		if i > 0 {
			out += ", "
		}
		out += f
	}
	return out
}
