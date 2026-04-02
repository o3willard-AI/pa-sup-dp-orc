//go:build windows && integration
// +build windows,integration

package terminal

import (
	"context"
	"testing"
)

func TestIntegration_WindowsAdapter_RealUI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	adapter := NewAdapter(DefaultConfig())
	ctx := context.Background()

	// Test Available with real UI Automation
	available := adapter.Available(ctx)
	if !available {
		t.Skip("UI Automation not available on this system")
	}

	// Test ListSessions with real terminals
	sessions, err := adapter.ListSessions(ctx)
	if err != nil {
		t.Errorf("ListSessions failed: %v", err)
	}

	t.Logf("Found %d terminal sessions", len(sessions))

	// If we have sessions, test Capture
	if len(sessions) > 0 {
		content, err := adapter.Capture(ctx, sessions[0].ID)
		if err != nil {
			t.Errorf("Capture failed: %v", err)
		}
		t.Logf("Captured %d characters from terminal", len(content))

		// Test GetDimensions
		rows, cols, err := adapter.GetDimensions(ctx, sessions[0].ID)
		if err != nil {
			t.Errorf("GetDimensions failed: %v", err)
		}
		if rows <= 0 || cols <= 0 {
			t.Errorf("Invalid dimensions: %dx%d", rows, cols)
		}
		t.Logf("Terminal dimensions: %dx%d", rows, cols)
	}

	// Cleanup
	adapter.Stop()
}

func TestIntegration_WindowsAdapter_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	adapter := NewAdapter(DefaultConfig())
	ctx := context.Background()

	// Test with invalid PID (should fallback to focused terminal)
	content, err := adapter.Capture(ctx, "windows-invalid-pid-test")
	if err != nil {
		t.Logf("Fallback failed with error (expected if no focused terminal): %v", err)
	} else {
		t.Logf("Fallback succeeded, captured %d characters from focused terminal", len(content))
	}

	adapter.Stop()
}
