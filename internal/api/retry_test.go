// internal/api/retry_test.go
package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/salmonumbrella/notte-cli/internal/testutil"
)

func TestRetryConfig_ShouldRetry(t *testing.T) {
	cfg := DefaultRetryConfig()

	tests := []struct {
		name        string
		statusCode  int
		method      string
		attempt     int
		shouldRetry bool
	}{
		{"429 first attempt", 429, "GET", 0, true},
		{"429 second attempt", 429, "GET", 1, true},
		{"429 third attempt", 429, "GET", 2, true},
		{"429 fourth attempt", 429, "GET", 3, false}, // Max retries exceeded
		{"500 GET", 500, "GET", 0, true},
		{"500 POST", 500, "POST", 0, false}, // Non-idempotent
		{"502 GET", 502, "GET", 0, true},
		{"400 any", 400, "GET", 0, false}, // Client error
		{"401 any", 401, "GET", 0, false}, // Auth error
		{"200 any", 200, "GET", 0, false}, // Success
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.ShouldRetry(tt.statusCode, tt.method, tt.attempt)
			if got != tt.shouldRetry {
				t.Errorf("ShouldRetry(%d, %s, %d) = %v, want %v",
					tt.statusCode, tt.method, tt.attempt, got, tt.shouldRetry)
			}
		})
	}
}

func TestRetryConfig_Backoff(t *testing.T) {
	cfg := DefaultRetryConfig()

	// First backoff should be around initial (1s)
	b0 := cfg.Backoff(0)
	if b0 < 500*time.Millisecond || b0 > 2*time.Second {
		t.Errorf("Backoff(0) = %v, want ~1s", b0)
	}

	// Second backoff should be around 2s
	b1 := cfg.Backoff(1)
	if b1 < time.Second || b1 > 4*time.Second {
		t.Errorf("Backoff(1) = %v, want ~2s", b1)
	}

	// Should cap at max
	b10 := cfg.Backoff(10)
	if b10 > cfg.MaxBackoff+time.Second {
		t.Errorf("Backoff(10) = %v, should be capped at %v", b10, cfg.MaxBackoff)
	}
}

func TestDoWithRetry_Success(t *testing.T) {
	server := testutil.NewMockServer()
	defer server.Close()

	server.AddResponse("/test", http.StatusOK, `{"status": "ok"}`)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL()+"/test", nil)

	resp, err := DoWithRetry(context.Background(), client, req, DefaultRetryConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
}

func TestDoWithRetry_ContextCancellation(t *testing.T) {
	server := testutil.NewMockServer()
	defer server.Close()

	server.AddResponse("/slow", http.StatusOK, `{}`)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := &http.Client{}
	req, _ := http.NewRequestWithContext(ctx, "GET", server.URL()+"/slow", nil)

	_, err := DoWithRetry(ctx, client, req, DefaultRetryConfig())
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestDoWithRetry_NilConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	resp, err := DoWithRetry(context.Background(), http.DefaultClient, req, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDoWithRetry_NonIdempotentNoRetry(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPost, server.URL, nil)
	cfg := &RetryConfig{MaxRetries: 3, InitialBackoff: time.Millisecond}

	resp, _ := DoWithRetry(context.Background(), http.DefaultClient, req, cfg)
	if resp != nil {
		_ = resp.Body.Close()
	}

	if callCount != 1 {
		t.Errorf("expected 1 call (no retry for POST), got %d", callCount)
	}
}

func TestSleepWithContext_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	start := time.Now()
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	sleepWithContext(ctx, time.Second)
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("sleep should have been cancelled, took %v", elapsed)
	}
}
