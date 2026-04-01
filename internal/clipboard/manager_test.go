package clipboard

import (
	"context"
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
