package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAPIKey_EnvVar(t *testing.T) {
	// Set env var
	os.Setenv(EnvAPIKey, "env_test_key")
	defer os.Unsetenv(EnvAPIKey)

	key, source, err := GetAPIKey("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "env_test_key" {
		t.Errorf("got %q, want 'env_test_key'", key)
	}
	if source != SourceEnv {
		t.Errorf("got source %q, want %q", source, SourceEnv)
	}
}

func TestGetAPIKey_ConfigFile(t *testing.T) {
	// Ensure no env var
	os.Unsetenv(EnvAPIKey)

	// Create temp config
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")
	content := `{"api_key": "config_test_key"}`
	if err := os.WriteFile(cfgPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	key, source, err := GetAPIKey(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "config_test_key" {
		t.Errorf("got %q, want 'config_test_key'", key)
	}
	if source != SourceConfig {
		t.Errorf("got source %q, want %q", source, SourceConfig)
	}
}

func TestGetAPIKey_NotFound(t *testing.T) {
	os.Unsetenv(EnvAPIKey)

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")

	_, _, err := GetAPIKey(cfgPath)
	if err == nil {
		t.Error("expected error when no API key found")
	}
}
