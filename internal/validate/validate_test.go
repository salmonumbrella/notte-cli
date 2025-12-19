// internal/validate/validate_test.go
package validate

import (
	"testing"
)

func TestURL(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"https://example.com", false},
		{"http://localhost:8080", false},
		{"https://api.notte.cc/v1", false},
		{"", true},
		{"not-a-url", true},
		{"ftp://example.com", true}, // Only http/https allowed
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := URL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{`{}`, false},
		{`{"key": "value"}`, false},
		{`[1, 2, 3]`, false},
		{`"string"`, false},
		{``, true},
		{`{invalid}`, true},
		{`{"key": }`, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := JSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestBrowser(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"chromium", false},
		{"firefox", false},
		{"webkit", false},
		{"chrome", true}, // Not valid - should be chromium
		{"safari", true}, // Not valid - should be webkit
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := Browser(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Browser(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestPositiveInt(t *testing.T) {
	tests := []struct {
		input   int
		wantErr bool
	}{
		{1, false},
		{100, false},
		{0, true},
		{-1, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := PositiveInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PositiveInt(%d) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"30s", false},
		{"5m", false},
		{"1h", false},
		{"1h30m", false},
		{"invalid", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := Duration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Duration(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestOutputFormat(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"text", false},
		{"json", false},
		{"yaml", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := OutputFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("OutputFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
