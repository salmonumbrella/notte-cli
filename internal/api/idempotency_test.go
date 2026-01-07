// internal/api/idempotency_test.go
package api

import (
	"net/http"
	"testing"
)

func TestGenerateIdempotencyKey(t *testing.T) {
	key1, err := GenerateIdempotencyKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	key2, err := GenerateIdempotencyKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

func TestAddIdempotencyKey_MutatingMethods(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, _ := http.NewRequest(method, "http://example.com", nil)
			AddIdempotencyKey(req)

			key := req.Header.Get(IdempotencyKeyHeader)
			if key == "" {
				t.Errorf("%s request should have Idempotency-Key header", method)
			}
		})
	}
}

func TestAddIdempotencyKey_NonMutatingMethods(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodHead, http.MethodOptions}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, _ := http.NewRequest(method, "http://example.com", nil)
			AddIdempotencyKey(req)

			key := req.Header.Get(IdempotencyKeyHeader)
			if key != "" {
				t.Errorf("%s request should NOT have Idempotency-Key header, got %q", method, key)
			}
		})
	}
}

func TestAddIdempotencyKey_PreservesExisting(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "http://example.com", nil)
	req.Header.Set(IdempotencyKeyHeader, "existing-key")

	AddIdempotencyKey(req)

	key := req.Header.Get(IdempotencyKeyHeader)
	if key != "existing-key" {
		t.Errorf("should preserve existing key, got %q", key)
	}
}
