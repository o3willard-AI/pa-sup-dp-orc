//go:build darwin
// +build darwin

package macos

import (
	"testing"
)

func TestIsAccessibilityEnabled(t *testing.T) {
	enabled := IsAccessibilityEnabled()
	t.Logf("Accessibility enabled: %v", enabled)
}

func TestRequestAccessibilityPermission(t *testing.T) {
	t.Skip("Manual: Call RequestAccessibilityPermission() and check prompt/grant")
}

func TestGetTerminalWindowPids(t *testing.T) {
	pids, err := GetTerminalWindowPids()
	if err != nil {
		t.Logf("No Terminal.app windows found: %v", err)
		return
	}
	t.Logf("Found %d Terminal.app PIDs: %v", len(pids), pids)
}

func TestGetTerminalWindowTitles(t *testing.T) {
	titles, err := GetTerminalWindowTitles()
	if err != nil {
		t.Logf("No Terminal.app window titles found: %v", err)
		return
	}
	t.Logf("Found %d window titles: %v", len(titles), titles)
}

func TestExtractTextFromFrontmostTerminal(t *testing.T) {
	text, err := ExtractTextFromFrontmostTerminal()
	if err != nil {
		t.Logf("Failed to extract text: %v", err)
		return
	}
	t.Logf("Extracted %d characters", len(text))
}
