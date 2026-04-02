//go:build !darwin && !windows
// +build !darwin,!windows

// Package adapters provides terminal adapter registration and factory functions.
package adapters

import (
	"github.com/pairadmin/pairadmin/internal/terminal"
	"github.com/pairadmin/pairadmin/internal/terminal/tmux"
)

// NewDetectorWithAutoAdapters creates a detector with auto-registered adapters.
// On non-macOS platforms, only tmux adapter is registered.
func NewDetectorWithAutoAdapters(config terminal.DetectorConfig) *terminal.Detector {
	adapters := autoRegisterAdapters()
	return terminal.NewDetector(config, adapters...)
}

func autoRegisterAdapters() []terminal.TerminalAdapter {
	var adapters []terminal.TerminalAdapter

	tmuxAdapter := tmux.NewAdapter(tmux.DefaultConfig())
	adapters = append(adapters, tmuxAdapter)

	return adapters
}

// RequestAccessibilityPermission always returns false on non-macOS.
func RequestAccessibilityPermission() bool {
	return false
}

// IsAccessibilityEnabled always returns false on non-macOS.
func IsAccessibilityEnabled() bool {
	return false
}
