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

func TestPrintListOrEmpty_NonSlice(t *testing.T) {
	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	printed, err := PrintListOrEmpty("not a slice", "No files.")
	if err == nil {
		t.Fatalf("expected error for non-slice type")
	}
	if printed {
		t.Fatalf("expected printed=false for non-slice type")
	}
	if !strings.Contains(err.Error(), "expected slice") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestPrintListOrEmpty_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		items     any
		emptyMsg  string
		wantPrint bool
		wantErr   bool
	}{
		{
			name:      "nil slice",
			items:     nil,
			emptyMsg:  "No items",
			wantPrint: true,
			wantErr:   false,
		},
		{
			name:      "empty string slice",
			items:     []string{},
			emptyMsg:  "No strings",
			wantPrint: true,
			wantErr:   false,
		},
		{
			name:      "non-empty slice",
			items:     []string{"a", "b"},
			emptyMsg:  "No items",
			wantPrint: false,
			wantErr:   false,
		},
		{
			name:      "not a slice",
			items:     "string",
			emptyMsg:  "Error",
			wantPrint: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origFormat := outputFormat
			outputFormat = "text"
			t.Cleanup(func() { outputFormat = origFormat })

			printed, err := PrintListOrEmpty(tt.items, tt.emptyMsg)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if printed != tt.wantPrint {
				t.Errorf("printed = %v, want %v", printed, tt.wantPrint)
			}
		})
	}
}
