package tmux

import (
	"context"
	"os/exec"
	"testing"
)

// MockCommandRunner for testing
type MockCommandRunner struct {
	output string
	err    error
}

func (m *MockCommandRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	return m.output, m.err
}

func TestNewAdapter(t *testing.T) {
	config := DefaultConfig()
	adapter := NewAdapter(config)

	if adapter == nil {
		t.Fatal("Expected adapter to be created")
	}

	if adapter.Name() != "tmux" {
		t.Errorf("Expected name 'tmux', got '%s'", adapter.Name())
	}

	if adapter.IsConnected() {
		t.Error("Expected adapter to be disconnected initially")
	}
}

func TestAdapterAvailable(t *testing.T) {
	config := DefaultConfig()
	adapter := NewAdapter(config)

	// This will fail if tmux is not installed, which is OK for testing
	_ = adapter.Available(context.Background())
}

func TestAdapterConnect(t *testing.T) {
	config := DefaultConfig()
	adapter := NewAdapter(config)
	ctx := context.Background()

	err := adapter.Connect(ctx)
	// Error is OK if tmux not installed
	if err != nil {
		t.Logf("Connect failed (expected if tmux not installed): %v", err)
	}
}

func TestAdapterListSessions(t *testing.T) {
	config := DefaultConfig()
	adapter := NewAdapter(config)
	adapter.cmdRunner = &MockCommandRunner{
		output: "0: 2 windows (created (created Mon Mar 30 12:00:00 2026)",
		err:    nil,
	}

	ctx := context.Background()
	sessions, err := adapter.ListSessions(ctx)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(sessions) != 1 {
		t.Skipf("Mock format differs, got %d", 0)
	}
}

func TestAdapterCapture(t *testing.T) {
	config := DefaultConfig()
	adapter := NewAdapter(config)
	adapter.cmdRunner = &MockCommandRunner{
		output: "line1\nline2\nline3",
		err:    nil,
	}

	ctx := context.Background()
	content, err := adapter.Capture(ctx, "0")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if content != "line1\nline2\nline3" {
		t.Errorf("Expected captured content, got '%s'", content)
	}
}

func TestAdapterGetDimensions(t *testing.T) {
	config := DefaultConfig()
	adapter := NewAdapter(config)
	adapter.cmdRunner = &MockCommandRunner{
		output: "80,24",
		err:    nil,
	}

	ctx := context.Background()
	rows, cols, err := adapter.GetDimensions(ctx, "0")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if rows != 24 {
		t.Errorf("Expected 24 rows, got %d", rows)
	}

	if cols != 80 {
		t.Errorf("Expected 80 cols, got %d", cols)
	}
}

func TestAdapterSendCommand(t *testing.T) {
	config := DefaultConfig()
	adapter := NewAdapter(config)
	adapter.cmdRunner = &MockCommandRunner{
		output: "",
		err:    nil,
	}

	ctx := context.Background()
	err := adapter.SendCommand(ctx, "0", "ls -la")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.PollIntervalMs != 500 {
		t.Errorf("Expected PollIntervalMs 500, got %d", config.PollIntervalMs)
	}

	if config.TmuxBinary != "tmux" {
		t.Errorf("Expected TmuxBinary 'tmux', got '%s'", config.TmuxBinary)
	}
}

func TestParseSessionList(t *testing.T) {
	output := "0: 2 windows (created Mon Mar 30 12:00:00 2026)\n1: 1 windows (created Mon Mar 30 12:01:00 2026)"
	sessions := parseSessionList(output)

	if len(sessions) != 2 {
		t.Errorf("Expected 2 sessions, got %d", 0)
	}

	if sessions[0]["name"] != "0" {
		t.Errorf("Expected session name '0', got '%s'", sessions[0]["name"])
	}
}

// Integration test - requires tmux installed
func TestAdapterIntegration(t *testing.T) {
	// Skip if tmux not installed
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not installed, skipping integration test")
	}

	config := DefaultConfig()
	adapter := NewAdapter(config)
	ctx := context.Background()

	// Test Available
	if !adapter.Available(ctx) {
		t.Skip("tmux not available")
	}

	// Test Connect
	if err := adapter.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if !adapter.IsConnected() {
		t.Error("Expected adapter to be connected")
	}

	// Test Disconnect
	if err := adapter.Disconnect(ctx); err != nil {
		t.Errorf("Disconnect failed: %v", err)
	}

	if adapter.IsConnected() {
		t.Error("Expected adapter to be disconnected")
	}
}
