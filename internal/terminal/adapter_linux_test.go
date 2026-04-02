//go:build linux

package terminal

import (
	"context"
	"testing"
)

func TestAdapter_Name(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	if a.Name() != "linux" {
		t.Errorf("expected linux, got %s", a.Name())
	}
}

func TestAdapter_Available(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	ctx := context.Background()
	_ = a.Available(ctx) // may be true/false based on env
}

func TestAdapter_ListSessions(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	ctx := context.Background()
	sessions, err := a.ListSessions(ctx)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Found %d sessions", len(sessions))
}

func TestAdapter_Capture(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	ctx := context.Background()
	sessions, _ := a.ListSessions(ctx)
	if len(sessions) == 0 {
		t.Skip("No terminals")
	}
	text, err := a.Capture(ctx, sessions[0].ID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Captured: %s...", text[:min(100, len(text))])
}

func TestAdapter_GetDimensions(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	ctx := context.Background()
	sessions, _ := a.ListSessions(ctx)
	if len(sessions) == 0 {
		t.Skip("No terminals")
	}
	rows, cols, err := a.GetDimensions(ctx, sessions[0].ID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Dims: %dx%d", rows, cols)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
