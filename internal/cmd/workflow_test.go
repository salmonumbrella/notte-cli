package cmd

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/notte-cli/internal/testutil"
)

const (
	workflowIDTest = "wf_123"
	workflowRunID  = "run_123"
)

func setupWorkflowTest(t *testing.T) *testutil.MockServer {
	t.Helper()
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	t.Cleanup(func() { server.Close() })
	env.SetEnv("NOTTE_API_URL", server.URL())

	origWorkflowID := workflowID
	origRunID := runID
	workflowID = workflowIDTest
	runID = workflowRunID
	t.Cleanup(func() {
		workflowID = origWorkflowID
		runID = origRunID
	})

	return server
}

func workflowJSON() string {
	return `{"workflow_id":"` + workflowIDTest + `","latest_version":"1","status":"active","created_at":"2020-01-01T00:00:00Z","updated_at":"2020-01-01T00:00:00Z","versions":["1"]}`
}

func workflowWithLinkJSON() string {
	return `{"workflow_id":"` + workflowIDTest + `","latest_version":"1","status":"active","created_at":"2020-01-01T00:00:00Z","updated_at":"2020-01-01T00:00:00Z","versions":["1"],"url":"https://example.com/workflow.json"}`
}

func workflowRunJSON() string {
	return `{"workflow_id":"` + workflowIDTest + `","workflow_run_id":"` + workflowRunID + `","status":"RUNNING","created_at":"2020-01-01T00:00:00Z","updated_at":"2020-01-01T00:00:00Z"}`
}

func updateWorkflowRunJSON() string {
	return `{"workflow_id":"` + workflowIDTest + `","workflow_run_id":"` + workflowRunID + `","updated_at":"2020-01-01T00:00:00Z","status":"STOPPED"}`
}

func TestRunWorkflowShow(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest, 200, workflowWithLinkJSON())

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowShow(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowUpdate(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest, 200, workflowJSON())

	tmpFile, err := os.CreateTemp("", "workflow-update-*.json")
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

	origFile := workflowUpdateFile
	workflowUpdateFile = tmpFile.Name()
	t.Cleanup(func() { workflowUpdateFile = origFile })

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowUpdate(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowUpdate_MissingFile(t *testing.T) {
	_ = setupWorkflowTest(t)

	origFile := workflowUpdateFile
	workflowUpdateFile = "missing-workflow.json"
	t.Cleanup(func() { workflowUpdateFile = origFile })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	err := runWorkflowUpdate(cmd, nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "failed to open file") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWorkflowDelete(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest, 200, `{"message":"deleted","status":"deleted"}`)

	SetSkipConfirmation(true)
	t.Cleanup(func() { SetSkipConfirmation(false) })

	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowDelete(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(stdout, "deleted") {
		t.Errorf("expected delete message, got %q", stdout)
	}
}

func TestRunWorkflowDeleteCancelled(t *testing.T) {
	_ = setupWorkflowTest(t)

	origSkip := skipConfirmation
	t.Cleanup(func() { skipConfirmation = origSkip })
	skipConfirmation = false

	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	_, _ = w.WriteString("n\n")
	_ = w.Close()

	origStdin := os.Stdin
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = origStdin
		_ = r.Close()
	})

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowDelete(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(stdout, "Cancelled.") {
		t.Errorf("expected cancel message, got %q", stdout)
	}
}

func TestRunWorkflowRun(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest+"/runs/start", 200, `{"run_id":"`+workflowRunID+`"}`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowRun(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowRuns(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest+"/runs", 200, `{"items":[`+workflowRunJSON()+`]}`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowRuns(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowRuns_Empty(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest+"/runs", 200, `{"items":[]}`)

	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowRuns(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(stdout, "No workflow runs found.") {
		t.Errorf("expected empty message, got %q", stdout)
	}
}

func TestRunWorkflowFork(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest+"/fork", 200, workflowJSON())

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowFork(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowRunStop(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest+"/runs/"+workflowRunID, 200, updateWorkflowRunJSON())

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowRunStop(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowRunMetadata(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest+"/runs/"+workflowRunID, 200, workflowRunJSON())

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowRunMetadata(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowRunMetadataUpdate(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest+"/runs/"+workflowRunID, 200, updateWorkflowRunJSON())

	origMetadata := metadataJSON
	metadataJSON = `{"result":{"ok":true},"status":"DONE"}`
	t.Cleanup(func() { metadataJSON = origMetadata })

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowRunMetadataUpdate(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunWorkflowRunMetadataUpdate_InvalidJSON(t *testing.T) {
	_ = setupWorkflowTest(t)

	origMetadata := metadataJSON
	metadataJSON = "{"
	t.Cleanup(func() { metadataJSON = origMetadata })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	err := runWorkflowRunMetadataUpdate(cmd, nil)
	if err == nil {
		t.Fatal("expected error for invalid metadata JSON")
	}
	if !strings.Contains(err.Error(), "failed to parse JSON metadata") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWorkflowSchedule(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest+"/schedule", 200, `{"status":"scheduled"}`)

	origCron := cronExpression
	cronExpression = "0 0 * * *"
	t.Cleanup(func() { cronExpression = origCron })

	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowSchedule(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(stdout, "scheduled") {
		t.Errorf("expected schedule message, got %q", stdout)
	}
}

func TestRunWorkflowUnschedule(t *testing.T) {
	server := setupWorkflowTest(t)
	server.AddResponse("/workflows/"+workflowIDTest+"/schedule", 200, `{"status":"removed"}`)

	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runWorkflowUnschedule(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(stdout, "schedule removed") {
		t.Errorf("expected unschedule message, got %q", stdout)
	}
}
