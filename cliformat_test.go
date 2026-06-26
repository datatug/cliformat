package cliformat_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/datatug/cliformat"
	"github.com/spf13/cobra"
)

type item struct {
	Name  string `json:"name"  yaml:"name"`
	Value int    `json:"value" yaml:"value"`
}

var testItems = []item{
	{Name: "alpha", Value: 1},
	{Name: "beta", Value: 2},
}

func nameOf(i item) string { return i.Name }

func TestWriteList_Name(t *testing.T) {
	var buf bytes.Buffer
	if err := cliformat.WriteList(&buf, "name", testItems, nameOf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 || lines[0] != "alpha" || lines[1] != "beta" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWriteList_EmptyFormat_FallsBackToName(t *testing.T) {
	var buf bytes.Buffer
	if err := cliformat.WriteList(&buf, "", testItems, nameOf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "alpha") {
		t.Errorf("expected name output, got %q", buf.String())
	}
}

func TestWriteList_YAML(t *testing.T) {
	var buf bytes.Buffer
	if err := cliformat.WriteList(&buf, "yaml", testItems, nameOf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "name: alpha") || !strings.Contains(out, "value: 2") {
		t.Errorf("unexpected yaml output: %q", out)
	}
}

func TestWriteList_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := cliformat.WriteList(&buf, "json", testItems, nameOf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"name": "alpha"`) || !strings.Contains(out, `"value": 2`) {
		t.Errorf("unexpected json output: %q", out)
	}
	// Must be an array
	trimmed := strings.TrimSpace(out)
	if trimmed[0] != '[' || trimmed[len(trimmed)-1] != ']' {
		t.Errorf("json output is not an array: %q", out)
	}
}

func TestWriteList_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := cliformat.WriteList(&buf, "xml", testItems, nameOf)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("error should mention the bad format, got: %v", err)
	}
}

func TestAddFlag_And_GetFormat(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cliformat.AddFlag(cmd, "name")

	got := cliformat.GetFormat(cmd)
	if got != "name" {
		t.Errorf("expected default 'name', got %q", got)
	}
}
