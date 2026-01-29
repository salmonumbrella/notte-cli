//go:build integration

package integration

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAgentsList(t *testing.T) {
	// List agents - should work even if empty
	result := runCLI(t, "agents", "list")
	requireSuccess(t, result)
	t.Log("Successfully listed agents")
}

func TestAgentsStartAndStatus(t *testing.T) {
	// Start an agent with a simple task
	result := runCLIWithTimeout(t, 120*time.Second, "agents", "start",
		"--task", "Navigate to example.com and report the page title",
		"--max-steps", "5",
	)
	requireSuccess(t, result)

	// Parse the response to get agent ID
	var startResp struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &startResp); err != nil {
		t.Fatalf("Failed to parse agent start response: %v", err)
	}
	agentID := startResp.AgentID
	if agentID == "" {
		t.Fatal("No agent ID returned from start command")
	}
	t.Logf("Started agent: %s", agentID)

	// Ensure cleanup
	defer cleanupAgent(t, agentID)

	// Get agent status
	result = runCLI(t, "agents", "status", "--id", agentID)
	requireSuccess(t, result)
	if !containsString(result.Stdout, agentID) {
		t.Error("Agent status did not contain agent ID")
	}
	t.Log("Successfully retrieved agent status")

	// List agents - should include our agent
	result = runCLI(t, "agents", "list")
	requireSuccess(t, result)
	// Note: Agent might complete quickly and not appear in list

	// Stop the agent (if still running)
	result = runCLI(t, "agents", "stop", "--id", agentID)
	// Don't require success - agent might have already completed
	t.Log("Agent stop attempted")
}

func TestAgentsStartWithReasoningModel(t *testing.T) {
	// Start an agent with a reasoning model specified
	result := runCLIWithTimeout(t, 120*time.Second, "agents", "start",
		"--task", "What is 2 + 2? Just answer with the number.",
		"--max-steps", "3",
		"--reasoning-model", "claude-sonnet-4-20250514",
	)
	requireSuccess(t, result)

	var startResp struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &startResp); err != nil {
		t.Fatalf("Failed to parse agent start response: %v", err)
	}
	agentID := startResp.AgentID
	if agentID == "" {
		t.Fatal("No agent ID returned from start command")
	}
	t.Logf("Started agent with reasoning model: %s", agentID)

	defer cleanupAgent(t, agentID)

	// Get status
	result = runCLI(t, "agents", "status", "--id", agentID)
	requireSuccess(t, result)

	// Stop agent
	result = runCLI(t, "agents", "stop", "--id", agentID)
	t.Log("Agent with reasoning model test completed")
}

func TestAgentsStartWithSession(t *testing.T) {
	// First start a session
	result := runCLI(t, "sessions", "start", "--headless")
	requireSuccess(t, result)

	var sessionResp struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &sessionResp); err != nil {
		t.Fatalf("Failed to parse session start response: %v", err)
	}
	sessionID := sessionResp.SessionID
	if sessionID == "" {
		t.Fatal("No session ID returned")
	}
	defer cleanupSession(t, sessionID)

	// Start an agent on the existing session
	result = runCLIWithTimeout(t, 120*time.Second, "agents", "start",
		"--task", "Report the current page URL",
		"--session", sessionID,
		"--max-steps", "3",
	)
	requireSuccess(t, result)

	var agentResp struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &agentResp); err != nil {
		t.Fatalf("Failed to parse agent start response: %v", err)
	}
	agentID := agentResp.AgentID
	if agentID == "" {
		t.Fatal("No agent ID returned")
	}
	t.Logf("Started agent on session: %s", agentID)

	defer cleanupAgent(t, agentID)

	// Get agent status
	result = runCLI(t, "agents", "status", "--id", agentID)
	requireSuccess(t, result)

	// Stop agent
	result = runCLI(t, "agents", "stop", "--id", agentID)
	t.Log("Agent on session test completed")
}

func TestAgentsStatusNonexistent(t *testing.T) {
	// Try to get status of a non-existent agent
	result := runCLI(t, "agents", "status", "--id", "nonexistent-agent-id-12345")
	requireFailure(t, result)
	t.Log("Correctly failed to get status of non-existent agent")
}

func TestAgentsStopNonexistent(t *testing.T) {
	// Try to stop a non-existent agent
	result := runCLI(t, "agents", "stop", "--id", "nonexistent-agent-id-12345")
	requireFailure(t, result)
	t.Log("Correctly failed to stop non-existent agent")
}
