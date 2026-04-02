//go:build darwin
// +build darwin

package macos

import (
	"context"
	"testing"
	"time"
)

func TestAdapterName(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	if a.Name() != "macos" {
		t.Errorf("Expected 'macos', got %s", a.Name())
	}
}

func TestAdapterAvailable(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	available := a.Available(context.Background())
	t.Logf("Adapter available: %v", available)
}

func TestAdapterListSessions(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	sessions, err := a.ListSessions(context.Background())
	if err != nil {
		t.Logf("ListSessions error (expected if no Terminal.app): %v", err)
		return
	}
	t.Logf("Found %d sessions", len(sessions))
	for _, s := range sessions {
		t.Logf("  - %s (%s)", s.Name, s.ID)
	}
}

func TestAdapterCapture(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	sessions, err := a.ListSessions(context.Background())
	if err != nil || len(sessions) == 0 {
		t.Skip("No Terminal.app sessions available")
	}

	content, err := a.Capture(context.Background(), sessions[0].ID)
	if err != nil {
		t.Logf("Capture error: %v", err)
		return
	}
	t.Logf("Captured %d characters", len(content))
}

func TestAdapterGetDimensions(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	sessions, err := a.ListSessions(context.Background())
	if err != nil || len(sessions) == 0 {
		t.Skip("No Terminal.app sessions available")
	}

	rows, cols, err := a.GetDimensions(context.Background(), sessions[0].ID)
	if err != nil {
		t.Logf("GetDimensions error: %v", err)
		return
	}
	t.Logf("Dimensions: %d rows x %d cols", rows, cols)
}

func TestAdapterSubscribe(t *testing.T) {
	a := NewAdapter(DefaultConfig())
	sessions, err := a.ListSessions(context.Background())
	if err != nil || len(sessions) == 0 {
		t.Skip("No Terminal.app sessions available")
	}

	ctx := context.Background()
	events, err := a.Subscribe(ctx, sessions[0].ID)
	if err != nil {
		t.Fatalf("Subscribe error: %v", err)
	}

	timeout := time.After(2 * time.Second)
	select {
	case event := <-events:
		t.Logf("Received event: %s", event.Type)
	case <-timeout:
		t.Log("No events received within timeout")
	}

	a.Stop()
}
