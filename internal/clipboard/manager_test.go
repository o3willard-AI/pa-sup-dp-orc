package clipboard

import (
	"context"
	"strings"
	"testing"

	clip "golang.design/x/clipboard"
)

func TestManager_CopyAndRead(t *testing.T) {
	// Skip test if clipboard is not available (e.g., in CI)
	if err := clip.Init(); err != nil {
		t.Skipf("Clipboard not available in this environment: %v", err)
	}

	m := NewManager()
	ctx := context.Background()
	text := "test clipboard content"

	// Copy
	err := m.CopyToTerminal(ctx, text, "")
	if err != nil {
		t.Fatalf("CopyToTerminal failed: %v", err)
	}

	// Read
	read, err := m.ReadFromClipboard(ctx)
	if err != nil {
		t.Fatalf("ReadFromClipboard failed: %v", err)
	}
	if read != text {
		t.Errorf("Expected %q, got %q", text, read)
	}
}

func TestCopyToTerminalEmptyString(t *testing.T) {
	if err := clip.Init(); err != nil {
		t.Skipf("Clipboard not available in this environment: %v", err)
	}

	m := NewManager()
	ctx := context.Background()

	err := m.CopyToTerminal(ctx, "", "")
	if err != nil {
		t.Errorf("CopyToTerminal with empty string failed: %v", err)
	}
}

func TestReadFromClipboardEmpty(t *testing.T) {
	if err := clip.Init(); err != nil {
		t.Skipf("Clipboard not available in this environment: %v", err)
	}

	m := NewManager()
	ctx := context.Background()

	// Ensure clipboard is empty (we cannot guarantee, but we can try to read)
	// If clipboard is empty, ReadFromClipboard should return an error.
	// However, we cannot guarantee emptiness, so we skip if clipboard contains data.
	// This test is best-effort.
	_, err := m.ReadFromClipboard(ctx)
	if err != nil && strings.Contains(err.Error(), "clipboard empty") {
		// Expected error for empty clipboard
		return
	}
	// If no error, clipboard had content; we cannot test empty case.
	t.Skip("Clipboard not empty; cannot test empty read")
}

func TestCopyToTerminalWithTerminalID(t *testing.T) {
	if err := clip.Init(); err != nil {
		t.Skipf("Clipboard not available in this environment: %v", err)
	}

	m := NewManager()
	ctx := context.Background()
	text := "test with terminal ID"

	err := m.CopyToTerminal(ctx, text, "terminal-123")
	if err != nil {
		t.Errorf("CopyToTerminal with terminal ID failed: %v", err)
	}
}
