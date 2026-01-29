package auth

import "testing"

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
