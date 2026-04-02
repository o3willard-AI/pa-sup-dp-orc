// Package testhelpers provides testing utilities for terminal adapter testing
package testhelpers

import (
	"strings"
	"testing"
	"time"

	"github.com/pairadmin/pairadmin/internal/terminal"
)

// AssertTerminalContent verifies captured content matches expected with optional tolerance
func AssertTerminalContent(t *testing.T, actual, expected string, tolerance float64) {
	t.Helper()
	if actual == expected {
		return
	}

	// Check if within tolerance (for partial matches)
	if tolerance > 0 {
		matchLen := 0
		for i := 0; i < len(actual) && i < len(expected); i++ {
			if actual[i] == expected[i] {
				matchLen++
			}
		}
		matchPct := float64(matchLen) / float64(max(len(actual), len(expected)))
		if matchPct >= tolerance {
			return
		}
	}

	t.Errorf("content mismatch:\nexpected: %q\nactual:   %q", expected, actual)
}

// AssertTerminalDetected verifies a terminal with given ID was detected
func AssertTerminalDetected(t *testing.T, terminals []terminal.DetectedTerminal, id string) {
	t.Helper()
	for _, term := range terminals {
		if term.ID == id {
			return
		}
	}
	t.Errorf("terminal %q not found in detected terminals", id)
}

// AssertTerminalNotDetected verifies a terminal with given ID was NOT detected
func AssertTerminalNotDetected(t *testing.T, terminals []terminal.DetectedTerminal, id string) {
	t.Helper()
	for _, term := range terminals {
		if term.ID == id {
			t.Errorf("terminal %q should not be detected", id)
			return
		}
	}
}

// AssertEventReceived verifies an event was received within timeout
func AssertEventReceived(t *testing.T, events <-chan terminal.TerminalEvent, eventType terminal.TerminalEventType, timeout time.Duration) terminal.TerminalEvent {
	t.Helper()
	select {
	case event := <-events:
		if event.Type != eventType {
			t.Errorf("expected event type %v, got %v", eventType, event.Type)
		}
		return event
	case <-time.After(timeout):
		t.Errorf("timeout waiting for event type %v", eventType)
		return terminal.TerminalEvent{}
	}
}

// AssertEventNotReceived verifies no event was received within timeout
func AssertEventNotReceived(t *testing.T, events <-chan terminal.TerminalEvent, timeout time.Duration) {
	t.Helper()
	select {
	case event := <-events:
		t.Errorf("unexpected event received: %v", event)
	case <-time.After(timeout):
		// Expected - no event
	}
}

// AssertContains verifies string contains substring
func AssertContains(t *testing.T, str, substr string) {
	t.Helper()
	if !strings.Contains(str, substr) {
		t.Errorf("string %q does not contain %q", str, substr)
	}
}

// AssertNotEmpty verifies string is not empty
func AssertNotEmpty(t *testing.T, str string, name string) {
	t.Helper()
	if str == "" {
		t.Errorf("%s should not be empty", name)
	}
}

// AssertDimensions verifies terminal dimensions
func AssertDimensions(t *testing.T, rows, cols int, expectedRows, expectedCols int) {
	t.Helper()
	if rows != expectedRows {
		t.Errorf("expected %d rows, got %d", expectedRows, rows)
	}
	if cols != expectedCols {
		t.Errorf("expected %d cols, got %d", expectedCols, cols)
	}
}

// AssertError verifies an error occurred (or didn't)
func AssertError(t *testing.T, err error, shouldExist bool) {
	t.Helper()
	if shouldExist && err == nil {
		t.Error("expected error, got nil")
	}
	if !shouldExist && err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
