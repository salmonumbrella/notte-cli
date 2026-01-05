package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestReadJSONInput_FromValue(t *testing.T) {
	cmd := &cobra.Command{}
	data, err := readJSONInput(cmd, `{"ok":true}`, "action")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"ok":true}` {
		t.Fatalf("unexpected data: %s", string(data))
	}
}

func TestReadJSONInput_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "payload.json")
	if err := os.WriteFile(path, []byte(`{"file":1}`), 0o600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	cmd := &cobra.Command{}
	data, err := readJSONInput(cmd, "@"+path, "data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"file":1}` {
		t.Fatalf("unexpected data: %s", string(data))
	}
}

func TestReadJSONInput_FromStdin(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetIn(bytes.NewBufferString(`{"stdin":true}`))
	data, err := readJSONInput(cmd, "-", "data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"stdin":true}` {
		t.Fatalf("unexpected data: %s", string(data))
	}
}

func TestReadJSONInput_EmptyStdin(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetIn(bytes.NewBufferString(""))
	_, err := readJSONInput(cmd, "-", "data")
	if err == nil {
		t.Fatalf("expected error for empty stdin")
	}
}
