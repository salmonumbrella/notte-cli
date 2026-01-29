//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFunctionsList(t *testing.T) {
	// List functions - should work even if empty
	result := runCLI(t, "functions", "list")
	requireSuccess(t, result)
	t.Log("Successfully listed functions")
}

func TestFunctionsCreateAndDelete(t *testing.T) {
	// Create a temporary Python file for the function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function.py")
	funcContent := `
# Test function for integration testing
def main():
    print("Hello from integration test!")
    return {"status": "success"}

if __name__ == "__main__":
    main()
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	// Create a new function
	result := runCLI(t, "functions", "create",
		"--file", funcFile,
		"--name", "integration-test-function",
		"--description", "Function for integration testing",
	)
	requireSuccess(t, result)

	// Parse the response to get function ID
	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	if functionID == "" {
		t.Fatal("No function ID returned from create command")
	}
	t.Logf("Created function: %s", functionID)

	// Ensure cleanup
	defer cleanupFunction(t, functionID)

	// Show function details (download URL)
	result = runCLI(t, "functions", "show", "--id", functionID)
	requireSuccess(t, result)
	t.Log("Successfully retrieved function details")

	// List functions - should include our function
	result = runCLI(t, "functions", "list")
	requireSuccess(t, result)
	if !containsString(result.Stdout, functionID) {
		t.Error("Function list did not contain our function")
	}

	// Delete the function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
	t.Log("Function deleted successfully")
}

func TestFunctionsUpdate(t *testing.T) {
	// Create initial function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_v1.py")
	funcContent := `
def main():
    return {"version": 1}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create", "--file", funcFile, "--name", "update-test-function")
	requireSuccess(t, result)

	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	defer cleanupFunction(t, functionID)

	// Create updated function file
	funcFileV2 := filepath.Join(tmpDir, "test_function_v2.py")
	funcContentV2 := `
def main():
    return {"version": 2}
`
	if err := os.WriteFile(funcFileV2, []byte(funcContentV2), 0o644); err != nil {
		t.Fatalf("Failed to create updated function file: %v", err)
	}

	// Update the function
	result = runCLI(t, "functions", "update", "--id", functionID, "--file", funcFileV2)
	requireSuccess(t, result)
	t.Log("Successfully updated function")

	// Delete function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
}

func TestFunctionsFork(t *testing.T) {
	// Create initial function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_fork.py")
	funcContent := `
def main():
    return {"message": "original"}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create", "--file", funcFile, "--name", "fork-test-function")
	requireSuccess(t, result)

	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	defer cleanupFunction(t, functionID)

	// Fork the function
	result = runCLI(t, "functions", "fork", "--id", functionID)
	requireSuccess(t, result)

	// Parse forked function ID
	var forkResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &forkResp); err != nil {
		t.Logf("Could not parse fork response: %v", err)
	} else if forkResp.FunctionID != "" {
		defer cleanupFunction(t, forkResp.FunctionID)
	}
	t.Log("Successfully forked function")

	// Delete original function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
}

func TestFunctionsShowNonexistent(t *testing.T) {
	// Try to show a non-existent function
	result := runCLI(t, "functions", "show", "--id", "nonexistent-function-id-12345")
	requireFailure(t, result)
	t.Log("Correctly failed to show non-existent function")
}

func TestFunctionsDeleteNonexistent(t *testing.T) {
	// Try to delete a non-existent function
	result := runCLI(t, "functions", "delete", "--id", "nonexistent-function-id-12345")
	requireFailure(t, result)
	t.Log("Correctly failed to delete non-existent function")
}

func TestFunctionsSchedule(t *testing.T) {
	// Create a function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_schedule.py")
	funcContent := `
def main():
    return {"scheduled": True}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create", "--file", funcFile, "--name", "schedule-test-function")
	requireSuccess(t, result)

	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	defer cleanupFunction(t, functionID)

	// Schedule the function (every hour)
	result = runCLI(t, "functions", "schedule", "--id", functionID, "--cron", "0 * * * *")
	requireSuccess(t, result)
	t.Log("Successfully scheduled function")

	// Unschedule the function
	result = runCLI(t, "functions", "unschedule", "--id", functionID)
	requireSuccess(t, result)
	t.Log("Successfully unscheduled function")

	// Delete function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
}

func TestFunctionsCreateShared(t *testing.T) {
	// Create a shared function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_shared.py")
	funcContent := `
def main():
    return {"public": True}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create",
		"--file", funcFile,
		"--name", "shared-test-function",
		"--shared",
	)
	requireSuccess(t, result)

	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	defer cleanupFunction(t, functionID)

	t.Log("Successfully created shared function")

	// Delete function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
}

func TestFunctionsRuns(t *testing.T) {
	// Create a function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_runs.py")
	funcContent := `
def main():
    return {"executed": True}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create", "--file", funcFile, "--name", "runs-test-function")
	requireSuccess(t, result)

	var createResp struct {
		FunctionID string `json:"function_id"`
	}
	if err := json.Unmarshal([]byte(result.Stdout), &createResp); err != nil {
		t.Fatalf("Failed to parse function create response: %v", err)
	}
	functionID := createResp.FunctionID
	defer cleanupFunction(t, functionID)

	// List function runs (should be empty initially)
	result = runCLI(t, "functions", "runs", "--id", functionID)
	requireSuccess(t, result)
	t.Log("Successfully listed function runs")

	// Delete function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
}

func TestFunctionsRunStart(t *testing.T) {
	// Create a function
	tmpDir := t.TempDir()
	funcFile := filepath.Join(tmpDir, "test_function_run.py")
	funcContent := `
def main():
    import time
    time.sleep(1)
    return {"executed": True}
`
	if err := os.WriteFile(funcFile, []byte(funcContent), 0o644); err != nil {
		t.Fatalf("Failed to create test function file: %v", err)
	}

	result := runCLI(t, "functions", "create", "--file", funcFile, "--name", "run-start-test-function")
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
	t.Log("Successfully started function run")

	// Delete function
	result = runCLI(t, "functions", "delete", "--id", functionID)
	requireSuccess(t, result)
}
