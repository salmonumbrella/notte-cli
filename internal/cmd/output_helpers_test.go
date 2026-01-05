package cmd

import (
	"strings"
	"testing"

	"github.com/salmonumbrella/notte-cli/internal/testutil"
)

func TestPrintListOrEmpty_JSON(t *testing.T) {
	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	stdout, _ := testutil.CaptureOutput(func() {
		printed, err := PrintListOrEmpty([]string{}, "No files.")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !printed {
			t.Fatalf("expected printed=true")
		}
	})

	if strings.TrimSpace(stdout) != "[]" {
		t.Fatalf("unexpected stdout: %q", stdout)
	}
}

func TestPrintListOrEmpty_Text(t *testing.T) {
	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	stdout, _ := testutil.CaptureOutput(func() {
		printed, err := PrintListOrEmpty([]string{}, "No files.")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !printed {
			t.Fatalf("expected printed=true")
		}
	})

	if strings.TrimSpace(stdout) != "No files." {
		t.Fatalf("unexpected stdout: %q", stdout)
	}
}
