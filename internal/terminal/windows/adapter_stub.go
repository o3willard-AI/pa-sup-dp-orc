//go:build !windows

package terminal

import (
	"context"
	"fmt"

	"github.com/pairadmin/pairadmin/internal/terminal"
)

// Adapter implements TerminalAdapter for Windows terminals (stub for non-Windows)
type Adapter struct{}

// Config holds Windows adapter configuration
type Config struct{}

// DefaultConfig returns default adapter configuration
func DefaultConfig() Config {
	return Config{}
}

// NewAdapter creates a new Windows Terminal adapter (stub)
func NewAdapter(config Config) *Adapter {
	return &Adapter{}
}

// Name returns the adapter name
func (a *Adapter) Name() string { return "windows" }

// Available checks if UI Automation is available (always false on non-Windows)
func (a *Adapter) Available(ctx context.Context) bool {
	return false
}

// ListSessions returns all detected terminal windows (stub)
func (a *Adapter) ListSessions(ctx context.Context) ([]terminal.DetectedTerminal, error) {
	return nil, fmt.Errorf("windows adapter not available on this platform")
}

// Capture captures the current content of a terminal window (stub)
func (a *Adapter) Capture(ctx context.Context, terminalID string) (string, error) {
	return "", fmt.Errorf("windows adapter not available on this platform")
}

// Subscribe starts streaming terminal events (stub)
func (a *Adapter) Subscribe(ctx context.Context, terminalID string) (<-chan terminal.TerminalEvent, error) {
	return nil, fmt.Errorf("windows adapter not available on this platform")
}

// GetDimensions returns terminal dimensions (stub)
func (a *Adapter) GetDimensions(ctx context.Context, terminalID string) (int, int, error) {
	return 0, 0, fmt.Errorf("windows adapter not available on this platform")
}

// Stop stops the adapter (stub)
func (a *Adapter) Stop() {}
