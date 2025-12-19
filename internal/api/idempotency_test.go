// internal/api/idempotency_test.go
package api

import (
	"testing"
)

func TestGenerateIdempotencyKey(t *testing.T) {
	key1 := GenerateIdempotencyKey()
	key2 := GenerateIdempotencyKey()

	// Should be 64 characters (32 bytes hex encoded)
	if len(key1) != 64 {
		t.Errorf("key length = %d, want 64", len(key1))
	}

	// Should be unique
	if key1 == key2 {
		t.Error("keys should be unique")
	}
}

func TestIsMutatingMethod(t *testing.T) {
	tests := []struct {
		method   string
		mutating bool
	}{
		{"GET", false},
		{"HEAD", false},
		{"OPTIONS", false},
		{"POST", true},
		{"PUT", true},
		{"PATCH", true},
		{"DELETE", true},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			got := IsMutatingMethod(tt.method)
			if got != tt.mutating {
				t.Errorf("IsMutatingMethod(%q) = %v, want %v", tt.method, got, tt.mutating)
			}
		})
	}
}
