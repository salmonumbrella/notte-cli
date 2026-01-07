package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/notte-cli/internal/testutil"
)

const sessionIDTest = "sess_123"

func setupSessionTest(t *testing.T) *testutil.MockServer {
	t.Helper()
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	t.Cleanup(func() { server.Close() })
	env.SetEnv("NOTTE_API_URL", server.URL())

	origID := sessionID
	sessionID = sessionIDTest
	t.Cleanup(func() { sessionID = origID })

	return server
}

func sessionJSON() string {
	return `{"session_id":"` + sessionIDTest + `","status":"ACTIVE","created_at":"2020-01-01T00:00:00Z","last_accessed_at":"2020-01-01T00:00:00Z","timeout_minutes":0}`
}

func TestRunSessionStatus(t *testing.T) {
	server := setupSessionTest(t)
	server.AddResponse("/sessions/"+sessionIDTest, 200, sessionJSON())

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionStatus(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionStop(t *testing.T) {
	server := setupSessionTest(t)
	server.AddResponse("/sessions/"+sessionIDTest+"/stop", 200, sessionJSON())

	SetSkipConfirmation(true)
	t.Cleanup(func() { SetSkipConfirmation(false) })

	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionStop(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(stdout, "stopped") {
		t.Errorf("expected stop message, got %q", stdout)
	}
}

func TestRunSessionObserve(t *testing.T) {
	server := setupSessionTest(t)
	observeResp := fmt.Sprintf(`{"metadata":{"tabs":[{"tab_id":1,"title":"Tab","url":"https://example.com"}],"title":"Tab","url":"https://example.com"},"screenshot":{"raw":"aGVsbG8="},"session":%s,"space":{"category":"page","description":"desc","interaction_actions":[]}}`, sessionJSON())
	server.AddResponse("/sessions/"+sessionIDTest+"/page/observe", 200, observeResp)

	origURL := sessionObserveURL
	sessionObserveURL = "https://example.com"
	t.Cleanup(func() { sessionObserveURL = origURL })

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionObserve(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionExecute(t *testing.T) {
	server := setupSessionTest(t)
	execResp := fmt.Sprintf(`{"action":{"type":"noop"},"data":{},"message":"ok","session":%s,"success":true}`, sessionJSON())
	server.AddResponse("/sessions/"+sessionIDTest+"/page/execute", 200, execResp)

	origAction := sessionExecuteAction
	sessionExecuteAction = `{"action":"noop"}`
	t.Cleanup(func() { sessionExecuteAction = origAction })

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionExecute(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionScrape(t *testing.T) {
	server := setupSessionTest(t)
	scrapeResp := fmt.Sprintf(`{"markdown":"hi","structured":{},"session":%s}`, sessionJSON())
	server.AddResponse("/sessions/"+sessionIDTest+"/page/scrape", 200, scrapeResp)

	origInstructions := sessionScrapeInstructions
	origOnlyMain := sessionScrapeOnlyMain
	sessionScrapeInstructions = "extract"
	sessionScrapeOnlyMain = true
	t.Cleanup(func() {
		sessionScrapeInstructions = origInstructions
		sessionScrapeOnlyMain = origOnlyMain
	})

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionScrape(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionCookies(t *testing.T) {
	server := setupSessionTest(t)
	server.AddResponse("/sessions/"+sessionIDTest+"/cookies", 200, `{"cookies":[{"domain":"example.com","httpOnly":true,"name":"a","path":"/","value":"b"}]}`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionCookies(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionCookiesSet(t *testing.T) {
	server := setupSessionTest(t)
	server.AddResponse("/sessions/"+sessionIDTest+"/cookies", 200, `{"message":"ok","success":true}`)

	tmpFile, err := os.CreateTemp("", "cookies-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := tmpFile.WriteString(`{"cookies":[{"domain":"example.com","httpOnly":true,"name":"a","path":"/","value":"b"}]}`); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(tmpFile.Name()) })

	origFile := sessionCookiesSetFile
	sessionCookiesSetFile = tmpFile.Name()
	t.Cleanup(func() { sessionCookiesSetFile = origFile })

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionCookiesSet(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionDebug(t *testing.T) {
	server := setupSessionTest(t)
	server.AddResponse("/sessions/"+sessionIDTest+"/debug", 200, `{"debug_url":"http://debug","tabs":[{"debug_url":"http://debug/tab","ws_url":"ws://tab","metadata":{"tab_id":1,"title":"t","url":"u"}}],"ws":{"cdp":"ws://cdp","logs":"ws://logs","recording":"ws://rec"}}`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionDebug(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionNetwork(t *testing.T) {
	server := setupSessionTest(t)
	server.AddResponse("/sessions/"+sessionIDTest+"/network/logs", 200, `{"requests":[],"responses":[],"session_id":"`+sessionIDTest+`","total_count":0}`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionNetwork(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionReplay(t *testing.T) {
	server := setupSessionTest(t)
	server.AddResponse("/sessions/"+sessionIDTest+"/replay", 200, `replay-data`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionReplay(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionOffset(t *testing.T) {
	server := setupSessionTest(t)
	server.AddResponse("/sessions/"+sessionIDTest+"/offset", 200, `{"offset":3}`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionOffset(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionWorkflowCode(t *testing.T) {
	server := setupSessionTest(t)
	server.AddResponse("/sessions/"+sessionIDTest+"/workflow/code", 200, `{"json_actions":[{"type":"noop"}],"python_script":"print('hi')"}`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionWorkflowCode(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}
