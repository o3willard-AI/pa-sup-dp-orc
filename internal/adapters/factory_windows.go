//go:build windows
// +build windows

// Package adapters provides terminal adapter registration and factory functions.
package adapters

import (
	"github.com/pairadmin/pairadmin/internal/terminal"
	"github.com/pairadmin/pairadmin/internal/terminal/tmux"
	windows "github.com/pairadmin/pairadmin/internal/terminal/windows"
)

// NewDetectorWithAutoAdapters creates a detector with auto-registered adapters.
// On Windows, both tmux and Windows UI Automation adapters are registered.
func NewDetectorWithAutoAdapters(config terminal.DetectorConfig) *terminal.Detector {
	adapters := autoRegisterAdapters()
	return terminal.NewDetector(config, adapters...)
}

func autoRegisterAdapters() []terminal.TerminalAdapter {
	var adapters []terminal.TerminalAdapter

	tmuxAdapter := tmux.NewAdapter(tmux.DefaultConfig())
	adapters = append(adapters, tmuxAdapter)

	windowsAdapter := windows.NewAdapter(windows.DefaultConfig())
	adapters = append(adapters, windowsAdapter)

	return adapters
}

// RequestAccessibilityPermission always returns false on Windows.
func RequestAccessibilityPermission() bool {
	return false
}

// IsAccessibilityEnabled always returns false on Windows.
func IsAccessibilityEnabled() bool {
	return false
}
