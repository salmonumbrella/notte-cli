// internal/testutil/httpserver_test.go
package testutil

import (
	"io"
	"net/http"
	"testing"
)

func TestMockServer_ReturnsResponse(t *testing.T) {
	server := NewMockServer()
	defer server.Close()

	server.AddResponse("/test", http.StatusOK, `{"message": "hello"}`)

	resp, err := http.Get(server.URL() + "/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"message": "hello"}` {
		t.Errorf("got body %q, want %q", string(body), `{"message": "hello"}`)
	}
}

func TestMockServer_TracksRequests(t *testing.T) {
	server := NewMockServer()
	defer server.Close()

	server.AddResponse("/api/sessions", http.StatusOK, `{}`)

	http.Get(server.URL() + "/api/sessions")
	http.Get(server.URL() + "/api/sessions")

	requests := server.Requests("/api/sessions")
	if len(requests) != 2 {
		t.Errorf("got %d requests, want 2", len(requests))
	}
}

func TestMockServer_Returns404ForUnknownPaths(t *testing.T) {
	server := NewMockServer()
	defer server.Close()

	resp, err := http.Get(server.URL() + "/unknown")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
