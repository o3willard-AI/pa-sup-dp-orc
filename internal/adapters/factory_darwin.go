//go:build darwin

package adapters

import (
	"github.com/pairadmin/pairadmin/internal/terminal"
	"github.com/pairadmin/pairadmin/internal/terminal/macos"
	"github.com/pairadmin/pairadmin/internal/terminal/tmux"
)

// NewDetectorWithAutoAdapters creates detector with tmux + macos adapters (darwin).
func NewDetectorWithAutoAdapters(config terminal.DetectorConfig) *terminal.Detector {
	adapters := autoRegisterAdapters()
	return terminal.NewDetector(config, adapters...)
}

func autoRegisterAdapters() []terminal.TerminalAdapter {
	var adapters []terminal.TerminalAdapter

	tmuxAdapter := tmux.NewAdapter(tmux.DefaultConfig())
	adapters = append(adapters, tmuxAdapter)

	macosAdapter := macos.NewAdapter(macos.DefaultConfig())
	adapters = append(adapters, macosAdapter)

	return adapters
}

// RequestAccessibilityPermission delegates to macos.
func RequestAccessibilityPermission() bool {
	return macos.RequestAccessibilityPermission()
}

// IsAccessibilityEnabled delegates to macos.
func IsAccessibilityEnabled() bool {
	return macos.IsAccessibilityEnabled()
}
