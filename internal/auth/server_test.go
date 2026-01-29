package auth

import (
	"strings"
	"testing"

	"github.com/salmonumbrella/notte-cli/internal/config"
)

// TestNewSetupServerReturnsError verifies that NewSetupServer returns
// (*SetupServer, error) - this is a compile-time check that the signature
// is correct. If the function doesn't return an error, this test will
// fail to compile.
func TestNewSetupServerReturnsError(t *testing.T) {
	server, err := NewSetupServer()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if server == nil {
		t.Fatal("expected non-nil server")
	}
	if server.csrfToken == "" {
		t.Error("expected non-empty CSRF token")
	}
	if server.oauthState == "" {
		t.Error("expected non-empty OAuth state")
	}
}

func TestValidateConsoleURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "default URL is valid",
			url:     config.DefaultConsoleURL,
			wantErr: false,
		},
		{
			name:    "custom HTTPS URL is valid",
			url:     "https://staging.notte.cc",
			wantErr: false,
		},
		{
			name:      "HTTP URL is rejected",
			url:       "http://console.notte.cc",
			wantErr:   true,
			errSubstr: "must use HTTPS",
		},
		{
			name:      "HTTP localhost is rejected",
			url:       "http://localhost:8080",
			wantErr:   true,
			errSubstr: "must use HTTPS",
		},
		{
			name:      "invalid URL is rejected",
			url:       "://invalid",
			wantErr:   true,
			errSubstr: "invalid console URL",
		},
		{
			name:      "empty scheme is rejected",
			url:       "console.notte.cc",
			wantErr:   true,
			errSubstr: "must use HTTPS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConsoleURL(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateConsoleURL(%q) expected error, got nil", tt.url)
					return
				}
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("ValidateConsoleURL(%q) error = %q, want substring %q", tt.url, err.Error(), tt.errSubstr)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateConsoleURL(%q) unexpected error: %v", tt.url, err)
				}
			}
		})
	}
}
