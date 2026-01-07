package cmd

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/notte-cli/internal/testutil"
)

func TestRunWorkflowsList_Success(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/workflows", 200, `{"items":[{"workflow_id":"wf_1","latest_version":"1","status":"active","created_at":"2020-01-01T00:00:00Z","updated_at":"2020-01-01T00:00:00Z","versions":["1"]}]}`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowsList(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowsList_Empty(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/workflows", 200, `{"items":[]}`)

	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowsList(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(stdout, "No workflows found.") {
		t.Errorf("expected empty message, got %q", stdout)
	}
}

func TestRunWorkflowsCreate_Success(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/workflows", 200, `{"workflow_id":"wf_1","latest_version":"1","status":"active","created_at":"2020-01-01T00:00:00Z","updated_at":"2020-01-01T00:00:00Z","versions":["1"]}`)

	tmpFile, err := os.CreateTemp("", "workflow-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := tmpFile.WriteString(`{"steps":[]}`); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(tmpFile.Name()) })

	origFile := workflowsCreateFile
	origName := workflowsCreateName
	origDesc := workflowsCreateDescription
	origShared := workflowsCreateShared
	t.Cleanup(func() {
		workflowsCreateFile = origFile
		workflowsCreateName = origName
		workflowsCreateDescription = origDesc
		workflowsCreateShared = origShared
	})

	workflowsCreateFile = tmpFile.Name()
	workflowsCreateName = "Test Workflow"
	workflowsCreateDescription = "Test description"

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.Flags().BoolVar(&workflowsCreateShared, "shared", false, "")
	_ = cmd.Flags().Set("shared", "true")
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowsCreate(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowsCreate_Minimal(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/workflows", 200, `{"workflow_id":"wf_2","latest_version":"1","status":"active","created_at":"2020-01-01T00:00:00Z","updated_at":"2020-01-01T00:00:00Z","versions":["1"]}`)

	tmpFile, err := os.CreateTemp("", "workflow-min-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := tmpFile.WriteString(`{"steps":[]}`); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(tmpFile.Name()) })

	origFile := workflowsCreateFile
	origName := workflowsCreateName
	origDesc := workflowsCreateDescription
	origShared := workflowsCreateShared
	t.Cleanup(func() {
		workflowsCreateFile = origFile
		workflowsCreateName = origName
		workflowsCreateDescription = origDesc
		workflowsCreateShared = origShared
	})

	workflowsCreateFile = tmpFile.Name()
	workflowsCreateName = ""
	workflowsCreateDescription = ""
	workflowsCreateShared = false

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowsCreate(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowsCreate_MissingFile(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	origFile := workflowsCreateFile
	workflowsCreateFile = "missing-workflow.json"
	t.Cleanup(func() { workflowsCreateFile = origFile })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	err := runWorkflowsCreate(cmd, nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "failed to open file") {
		t.Fatalf("unexpected error: %v", err)
	}
}
