//go:build integration

package integration

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSessionsLifecycle(t *testing.T) {
	// Start a new session
	result := runCLI(t, "sessions", "start", "--headless")
	requireSuccess(t, result)

	// Parse the response to get session ID
	var startResp struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &startResp); err != nil {
		t.Fatalf("Failed to parse session start response: %v", err)
	}
	sessionID := startResp.SessionID
	if sessionID == "" {
		t.Fatal("No session ID returned from start command")
	}
	t.Logf("Started session: %s", sessionID)

	// Ensure cleanup
	defer cleanupSession(t, sessionID)

	// Get session status
	result = runCLI(t, "sessions", "status", "--id", sessionID)
	requireSuccess(t, result)
	if !containsString(result.Stdout, sessionID) {
		t.Error("Session status did not contain session ID")
	}

	// List sessions - should include our session
	result = runCLI(t, "sessions", "list")
	requireSuccess(t, result)
	if !containsString(result.Stdout, sessionID) {
		t.Error("Session list did not contain our session")
	}

	// Stop the session
	result = runCLI(t, "sessions", "stop", "--id", sessionID)
	requireSuccess(t, result)
	t.Log("Session stopped successfully")
}

func TestSessionsStartWithOptions(t *testing.T) {
	// Start session with custom options
	result := runCLI(t, "sessions", "start",
		"--headless",
		"--browser", "chromium",
		"--idle-timeout", "5",
		"--max-duration", "10",
	)
	requireSuccess(t, result)

	var startResp struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &startResp); err != nil {
		t.Fatalf("Failed to parse session start response: %v", err)
	}
	sessionID := startResp.SessionID
	if sessionID == "" {
		t.Fatal("No session ID returned from start command")
	}
	t.Logf("Started session with options: %s", sessionID)

	defer cleanupSession(t, sessionID)

	// Verify session is running
	result = runCLI(t, "sessions", "status", "--id", sessionID)
	requireSuccess(t, result)

	// Stop the session
	result = runCLI(t, "sessions", "stop", "--id", sessionID)
	requireSuccess(t, result)
}

func TestSessionsCookies(t *testing.T) {
	// Start a session
	result := runCLI(t, "sessions", "start", "--headless")
	requireSuccess(t, result)

	var startResp struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &startResp); err != nil {
		t.Fatalf("Failed to parse session start response: %v", err)
	}
	sessionID := startResp.SessionID
	defer cleanupSession(t, sessionID)

	// Get cookies (should be empty or minimal initially)
	result = runCLI(t, "sessions", "cookies", "--id", sessionID)
	requireSuccess(t, result)
	t.Log("Successfully retrieved session cookies")

	// Stop session
	result = runCLI(t, "sessions", "stop", "--id", sessionID)
	requireSuccess(t, result)
}

func TestSessionsObserve(t *testing.T) {
	// Start a session
	result := runCLI(t, "sessions", "start", "--headless")
	requireSuccess(t, result)

	var startResp struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &startResp); err != nil {
		t.Fatalf("Failed to parse session start response: %v", err)
	}
	sessionID := startResp.SessionID
	defer cleanupSession(t, sessionID)

	// Wait a moment for the session to be fully ready
	time.Sleep(2 * time.Second)

	// Observe the page (navigate to a URL)
	result = runCLIWithTimeout(t, 120*time.Second, "sessions", "observe", "--id", sessionID, "--url", "https://example.com")
	requireSuccess(t, result)
	t.Log("Successfully observed page")

	// Stop session
	result = runCLI(t, "sessions", "stop", "--id", sessionID)
	requireSuccess(t, result)
}

func TestSessionsScrape(t *testing.T) {
	// Start a session
	result := runCLI(t, "sessions", "start", "--headless")
	requireSuccess(t, result)

	var startResp struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &startResp); err != nil {
		t.Fatalf("Failed to parse session start response: %v", err)
	}
	sessionID := startResp.SessionID
	defer cleanupSession(t, sessionID)

	// Wait a moment for the session to be fully ready
	time.Sleep(2 * time.Second)

	// First navigate to a page
	result = runCLIWithTimeout(t, 120*time.Second, "sessions", "observe", "--id", sessionID, "--url", "https://example.com")
	requireSuccess(t, result)

	// Scrape the page content
	result = runCLIWithTimeout(t, 120*time.Second, "sessions", "scrape", "--id", sessionID)
	requireSuccess(t, result)
	t.Log("Successfully scraped page content")

	// Stop session
	result = runCLI(t, "sessions", "stop", "--id", sessionID)
	requireSuccess(t, result)
}

func TestSessionsNetwork(t *testing.T) {
	// Start a session
	result := runCLI(t, "sessions", "start", "--headless")
	requireSuccess(t, result)

	var startResp struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &startResp); err != nil {
		t.Fatalf("Failed to parse session start response: %v", err)
	}
	sessionID := startResp.SessionID
	defer cleanupSession(t, sessionID)

	// Wait a moment for the session to be fully ready
	time.Sleep(2 * time.Second)

	// Navigate to generate some network activity
	result = runCLIWithTimeout(t, 120*time.Second, "sessions", "observe", "--id", sessionID, "--url", "https://example.com")
	requireSuccess(t, result)

	// Get network logs
	result = runCLI(t, "sessions", "network", "--id", sessionID)
	requireSuccess(t, result)
	t.Log("Successfully retrieved network logs")

	// Stop session
	result = runCLI(t, "sessions", "stop", "--id", sessionID)
	requireSuccess(t, result)
}

func TestSessionsList(t *testing.T) {
	// List sessions - this should always work, even if empty
	result := runCLI(t, "sessions", "list")
	requireSuccess(t, result)
	t.Log("Successfully listed sessions")
}

func TestSessionsStatusNonexistent(t *testing.T) {
	// Try to get status of a non-existent session
	result := runCLI(t, "sessions", "status", "--id", "nonexistent-session-id-12345")
	requireFailure(t, result)
	t.Log("Correctly failed to get status of non-existent session")
}
