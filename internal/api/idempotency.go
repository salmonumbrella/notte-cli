// internal/api/idempotency.go
package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// IdempotencyKeyHeader is the header name for idempotency keys
const IdempotencyKeyHeader = "Idempotency-Key"

// GenerateIdempotencyKey generates a cryptographically random idempotency key
func GenerateIdempotencyKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based key if crypto/rand fails
		panic("crypto/rand failed: " + err.Error())
	}
	return hex.EncodeToString(bytes)
}

// IsMutatingMethod returns true for HTTP methods that modify state
func IsMutatingMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// AddIdempotencyKey adds an idempotency key header to mutating requests
func AddIdempotencyKey(req *http.Request) {
	if IsMutatingMethod(req.Method) {
		if req.Header.Get(IdempotencyKeyHeader) == "" {
			req.Header.Set(IdempotencyKeyHeader, GenerateIdempotencyKey())
		}
	}
}
