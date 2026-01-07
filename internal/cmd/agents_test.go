package cmd

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/notte-cli/internal/testutil"
)

func TestRunAgentsList_Success(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/agents", 200, `{"items":[{"agent_id":"agent_1","session_id":"sess_1","status":"RUNNING","created_at":"2020-01-01T00:00:00Z"}]}`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runAgentsList(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunAgentsList_Empty(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/agents", 200, `{"items":[]}`)

	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runAgentsList(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(stdout, "No running agents.") {
		t.Errorf("expected empty message, got %q", stdout)
	}
}

func TestRunAgentsStart_Success(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/agents/start", 200, `{"agent_id":"agent_1","session_id":"sess_1","status":"RUNNING","created_at":"2020-01-01T00:00:00Z"}`)

	origTask := agentsStartTask
	origSession := agentsStartSession
	origVault := agentsStartVault
	origPersona := agentsStartPersona
	origMaxSteps := agentsStartMaxSteps
	origReasoning := agentsStartReasoningModel
	t.Cleanup(func() {
		agentsStartTask = origTask
		agentsStartSession = origSession
		agentsStartVault = origVault
		agentsStartPersona = origPersona
		agentsStartMaxSteps = origMaxSteps
		agentsStartReasoningModel = origReasoning
	})

	agentsStartTask = "do the thing"
	agentsStartSession = "sess_123"
	agentsStartVault = "vault_123"
	agentsStartPersona = "persona_123"
	agentsStartMaxSteps = 5
	agentsStartReasoningModel = "custom-model"

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runAgentsStart(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}
