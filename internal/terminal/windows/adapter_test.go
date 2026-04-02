//go:build windows && !integration && mock
// +build windows,!integration,mock

package terminal

import (
	"context"
	"strings"
	"testing"
)

func TestAdapter_Name(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	if a.Name() != "windows" {
		t.Errorf("expected windows, got %s", a.Name())
	}
}

func TestAdapter_Available(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	ctx := context.Background()
	available := a.Available(ctx)
	if !available {
		t.Error("Expected adapter to be available with mocked C functions")
	}
}

func TestAdapter_ListSessions(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	ctx := context.Background()
	sessions, err := a.ListSessions(ctx)
	if err != nil {
		t.Errorf("ListSessions failed: %v", err)
	}
	if len(sessions) != 2 {
		t.Errorf("Expected 2 sessions, got %d", len(sessions))
	}

	// Verify at least one session has expected name
	foundPowerShell := false
	for _, session := range sessions {
		if strings.Contains(session.Name, "PowerShell") || strings.Contains(session.Name, "Command Prompt") {
			foundPowerShell = true
			break
		}
	}
	if !foundPowerShell {
		t.Error("Expected to find PowerShell or Command Prompt session")
	}
}

func TestAdapter_Capture(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	ctx := context.Background()
	sessions, err := a.ListSessions(ctx)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}
	if len(sessions) == 0 {
		t.Skip("No mock sessions available")
	}

	text, err := a.Capture(ctx, sessions[0].ID)
	if err != nil {
		t.Errorf("Capture failed: %v", err)
	}
	if !strings.Contains(text, "hello") && !strings.Contains(text, "dir") {
		t.Error("Capture returned unexpected content")
	}
}

func TestAdapter_GetDimensions(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	ctx := context.Background()
	sessions, err := a.ListSessions(ctx)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}
	if len(sessions) == 0 {
		t.Skip("No mock sessions available")
	}

	rows, cols, err := a.GetDimensions(ctx, sessions[0].ID)
	if err != nil {
		t.Error(err)
	}
	if rows != 24 || cols != 80 {
		t.Errorf("Expected dimensions 24x80, got %dx%d", rows, cols)
	}
}
