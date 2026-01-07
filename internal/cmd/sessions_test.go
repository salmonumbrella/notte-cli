package cmd

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/notte-cli/internal/testutil"
)

func TestRunSessionsList_Success(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/sessions", 200, `{"items": [{"session_id": "sess_123", "status": "ACTIVE"}]}`)

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionsList(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionsList_Empty(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/sessions", 200, `{"items": []}`)

	origFormat := outputFormat
	outputFormat = "text"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionsList(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" || !strings.Contains(stdout, "No active sessions.") {
		t.Errorf("expected empty message, got %q", stdout)
	}
}

func TestRunSessionsStart(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/sessions/start", 200, `{"session_id":"sess_123","status":"ACTIVE","created_at":"2020-01-01T00:00:00Z","last_accessed_at":"2020-01-01T00:00:00Z","timeout_minutes":5}`)

	origHeadless := sessionsStartHeadless
	origBrowser := sessionsStartBrowser
	origTimeout := sessionsStartTimeout
	origProxies := sessionsStartProxies
	origSolve := sessionsStartSolveCaptchas
	origVW := sessionsStartViewportW
	origVH := sessionsStartViewportH
	origUA := sessionsStartUserAgent
	origCDP := sessionsStartCdpURL
	t.Cleanup(func() {
		sessionsStartHeadless = origHeadless
		sessionsStartBrowser = origBrowser
		sessionsStartTimeout = origTimeout
		sessionsStartProxies = origProxies
		sessionsStartSolveCaptchas = origSolve
		sessionsStartViewportW = origVW
		sessionsStartViewportH = origVH
		sessionsStartUserAgent = origUA
		sessionsStartCdpURL = origCDP
	})

	sessionsStartHeadless = false
	sessionsStartBrowser = "firefox"
	sessionsStartTimeout = 5
	sessionsStartProxies = true
	sessionsStartSolveCaptchas = true
	sessionsStartViewportW = 1280
	sessionsStartViewportH = 720
	sessionsStartUserAgent = "test-agent"
	sessionsStartCdpURL = "ws://cdp"

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.Flags().BoolVar(&sessionsStartHeadless, "headless", true, "")
	cmd.Flags().BoolVar(&sessionsStartProxies, "proxies", false, "")
	cmd.Flags().BoolVar(&sessionsStartSolveCaptchas, "solve-captchas", false, "")
	_ = cmd.Flags().Set("headless", "false")
	_ = cmd.Flags().Set("proxies", "true")
	_ = cmd.Flags().Set("solve-captchas", "true")
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionsStart(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}

func TestRunSessionsStart_Minimal(t *testing.T) {
	env := testutil.SetupTestEnv(t)
	env.SetEnv("NOTTE_API_KEY", "test-key")

	server := testutil.NewMockServer()
	defer server.Close()
	env.SetEnv("NOTTE_API_URL", server.URL())

	server.AddResponse("/sessions/start", 200, `{"session_id":"sess_456","status":"ACTIVE","created_at":"2020-01-01T00:00:00Z","last_accessed_at":"2020-01-01T00:00:00Z","timeout_minutes":3}`)

	origHeadless := sessionsStartHeadless
	origBrowser := sessionsStartBrowser
	origTimeout := sessionsStartTimeout
	origProxies := sessionsStartProxies
	origSolve := sessionsStartSolveCaptchas
	origVW := sessionsStartViewportW
	origVH := sessionsStartViewportH
	origUA := sessionsStartUserAgent
	origCDP := sessionsStartCdpURL
	t.Cleanup(func() {
		sessionsStartHeadless = origHeadless
		sessionsStartBrowser = origBrowser
		sessionsStartTimeout = origTimeout
		sessionsStartProxies = origProxies
		sessionsStartSolveCaptchas = origSolve
		sessionsStartViewportW = origVW
		sessionsStartViewportH = origVH
		sessionsStartUserAgent = origUA
		sessionsStartCdpURL = origCDP
	})

	sessionsStartHeadless = true
	sessionsStartBrowser = ""
	sessionsStartTimeout = 0
	sessionsStartProxies = false
	sessionsStartSolveCaptchas = false
	sessionsStartViewportW = 0
	sessionsStartViewportH = 0
	sessionsStartUserAgent = ""
	sessionsStartCdpURL = ""

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := &cobra.Command{}
	cmd.Flags().BoolVar(&sessionsStartHeadless, "headless", true, "")
	cmd.Flags().BoolVar(&sessionsStartProxies, "proxies", false, "")
	cmd.Flags().BoolVar(&sessionsStartSolveCaptchas, "solve-captchas", false, "")
	cmd.SetContext(context.Background())

	stdout, _ := testutil.CaptureOutput(func() {
		err := runSessionsStart(cmd, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if stdout == "" {
		t.Error("expected output, got empty string")
	}
}
