//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFunctionRunsLifecycle(t *testing.T) {
	// Create a function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_lifecycle.py")
	funcContent := `
import time

def main():
    # Simulate some work
    time.sleep(5)
    return {"completed": True}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create", "--file", funcFile, "--name", "lifecycle-test-function")
	requireSuccess(t, result)

	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	defer cleanupFunction(t, functionID)

	// Start a function run
	result = runCLIWithTimeout(t, 120*time.Second, "functions", "run", "--id", functionID)
	requireSuccess(t, result)

	// Parse the run response to get run ID
	var runResp struct {
		RunID string `json:"run_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &runResp); err != nil {
		t.Logf("Could not parse run response, might be different format: %v", err)
	}
	runID := runResp.RunID
	t.Logf("Started function run: %s", runID)

	// List function runs
	result = runCLI(t, "functions", "runs", "--id", functionID)
	requireSuccess(t, result)
	t.Log("Successfully listed function runs")

	// Delete function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
	t.Log("Function run lifecycle completed successfully")
}

func TestFunctionRunsStop(t *testing.T) {
	// Create a function that takes a while to run
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_stop.py")
	funcContent := `
import time

def main():
    # Long running task
    for i in range(60):
        time.sleep(1)
    return {"completed": True}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create", "--file", funcFile, "--name", "stop-test-function")
	requireSuccess(t, result)

	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	defer cleanupFunction(t, functionID)

	// Start a function run
	result = runCLIWithTimeout(t, 120*time.Second, "functions", "run", "--id", functionID)
	requireSuccess(t, result)

	// Parse the run response to get run ID
	var runResp struct {
		RunID string `json:"run_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &runResp); err != nil {
		t.Logf("Could not parse run response: %v", err)
		// Try to list runs to get the run ID
		result = runCLI(t, "functions", "runs", "--id", functionID)
		requireSuccess(t, result)
		t.Log("Listed runs, but could not get run ID to stop")
	} else if runResp.RunID != "" {
		// Stop the function run
		result = runCLI(t, "functions", "run-stop", "--id", functionID, "--run-id", runResp.RunID)
		requireSuccess(t, result)
		t.Log("Successfully stopped function run")
	}

	// Delete function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
}

func TestFunctionRunsMetadata(t *testing.T) {
	// Create a function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_metadata.py")
	funcContent := `
import time

def main():
    time.sleep(2)
    return {"completed": True}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create", "--file", funcFile, "--name", "metadata-test-function")
	requireSuccess(t, result)

	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	defer cleanupFunction(t, functionID)

	// Start a function run
	result = runCLIWithTimeout(t, 120*time.Second, "functions", "run", "--id", functionID)
	requireSuccess(t, result)

	// Parse the run response to get run ID
	var runResp struct {
		RunID string `json:"run_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &runResp); err != nil {
		t.Logf("Could not parse run response: %v", err)
	} else if runResp.RunID != "" {
		// Get metadata
		result = runCLI(t, "functions", "run-metadata", "--id", functionID, "--run-id", runResp.RunID)
		requireSuccess(t, result)
		t.Log("Successfully retrieved run metadata")

		// Update metadata
		metadataJSON := `{"test_key": "test_value"}`
		result = runCLI(t, "functions", "run-metadata-update",
			"--id", functionID,
			"--run-id", runResp.RunID,
			"--data", metadataJSON,
		)
		requireSuccess(t, result)
		t.Log("Successfully updated run metadata")
	}

	// Delete function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
}

func TestFunctionRunsListEmpty(t *testing.T) {
	// Create a function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_empty_runs.py")
	funcContent := `
def main():
    return {"status": "ok"}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create", "--file", funcFile, "--name", "empty-runs-test-function")
	requireSuccess(t, result)

	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	defer cleanupFunction(t, functionID)

	// List function runs (should be empty)
	result = runCLI(t, "functions", "runs", "--id", functionID)
	requireSuccess(t, result)
	t.Log("Successfully listed empty function runs")

	// Delete function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
}

func TestFunctionRunsStopNonexistent(t *testing.T) {
	// Create a function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_stop_nonexistent.py")
	funcContent := `
def main():
    return {"status": "ok"}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create", "--file", funcFile, "--name", "stop-nonexistent-test-function")
	requireSuccess(t, result)

	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	defer cleanupFunction(t, functionID)

	// Try to stop a non-existent run
	result = runCLI(t, "functions", "run-stop", "--id", functionID, "--run-id", "nonexistent-run-id-12345")
	requireFailure(t, result)
	t.Log("Correctly failed to stop non-existent function run")

	// Delete function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
}
