package clipboard

import (
	"context"
	"fmt"

	"golang.design/x/clipboard"
)

// Manager provides cross‑platform clipboard operations.
type Manager struct {
	initErr error
}

// NewManager creates a new clipboard manager.
func NewManager() *Manager {
	m := &Manager{}
	m.initErr = clipboard.Init()
	return m
}

// CopyToTerminal copies text to clipboard and optionally focuses the target terminal.
func (m *Manager) CopyToTerminal(ctx context.Context, text string, terminalID string) error {
	if m.initErr != nil {
		return fmt.Errorf("clipboard init failed: %w", m.initErr)
	}
	// Write to clipboard
	clipboard.Write(clipboard.FmtText, []byte(text))
	// The returned channel signals when clipboard is overwritten; we ignore it.

	// TODO: Focus terminal if terminalID provided and focus‑paste configured
	// This requires platform‑specific window focus logic (future enhancement)

	return nil
}

// ReadFromClipboard reads text from clipboard.
func (m *Manager) ReadFromClipboard(ctx context.Context) (string, error) {
	if m.initErr != nil {
		return "", fmt.Errorf("clipboard init failed: %w", m.initErr)
	}
	data := clipboard.Read(clipboard.FmtText)
	if data == nil {
		return "", fmt.Errorf("clipboard empty or unsupported format")
	}
	return string(data), nil
}