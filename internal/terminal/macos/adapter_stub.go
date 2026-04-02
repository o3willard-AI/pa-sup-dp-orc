//go:build !darwin
// +build !darwin

// Package macos provides stub implementation for non-macOS platforms
package macos

import (
	"context"
	"errors"

	"github.com/pairadmin/pairadmin/internal/terminal"
)

// ErrNotOnMacOS is returned when macOS adapter is used on non-macOS platforms
var ErrNotOnMacOS = errors.New("macOS adapter only available on macOS")

// Adapter stub for non-macOS platforms
type Adapter struct{}

// Config holds adapter configuration
type Config struct {
	PollIntervalMs int
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{PollIntervalMs: 500}
}

// NewAdapter creates a stub adapter
func NewAdapter(config Config) *Adapter {
	return &Adapter{}
}

// Name returns the adapter name
func (a *Adapter) Name() string {
	return "macos"
}

// Available always returns false on non-macOS
func (a *Adapter) Available(ctx context.Context) bool {
	return false
}

// ListSessions returns error on non-macOS
func (a *Adapter) ListSessions(ctx context.Context) ([]terminal.DetectedTerminal, error) {
	return nil, ErrNotOnMacOS
}

// Capture returns error on non-macOS
func (a *Adapter) Capture(ctx context.Context, terminalID string) (string, error) {
	return "", ErrNotOnMacOS
}

// Subscribe returns error on non-macOS
func (a *Adapter) Subscribe(ctx context.Context, terminalID string) (<-chan terminal.TerminalEvent, error) {
	return nil, ErrNotOnMacOS
}

// GetDimensions returns error on non-macOS
func (a *Adapter) GetDimensions(ctx context.Context, terminalID string) (int, int, error) {
	return 0, 0, ErrNotOnMacOS
}

// Stop is a no-op on stub
func (a *Adapter) Stop() {}

// IsRunning always returns false on stub
func (a *Adapter) IsRunning() bool {
	return false
}
