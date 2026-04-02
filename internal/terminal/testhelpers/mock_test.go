package testhelpers_test

import (
	"context"
	"testing"
	"time"

	"github.com/pairadmin/pairadmin/internal/terminal"
	"github.com/pairadmin/pairadmin/internal/terminal/testhelpers"
)

// TestAdapterInterface tests that all adapters implement the TerminalAdapter interface
func TestAdapterInterface(t *testing.T) {
	// Create mock adapter
	adapter := testhelpers.NewMockAdapter("test")

	// Verify it implements TerminalAdapter interface
	var _ terminal.TerminalAdapter = adapter
}

// TestAdapter_Name tests the Name method
func TestAdapter_Name(t *testing.T) {
	adapter := testhelpers.NewMockAdapter("my-adapter")

	name := adapter.Name()
	if name != "my-adapter" {
		t.Errorf("expected name 'my-adapter', got %q", name)
	}
}

// TestAdapter_Available tests the Available method
func TestAdapter_Available(t *testing.T) {
	ctx := context.Background()

	// Test available adapter
	adapter := testhelpers.NewMockAdapter("test")
	if !adapter.Available(ctx) {
		t.Error("expected adapter to be available")
	}

	// Test unavailable adapter
	adapter.SetAvailable(false)
	if adapter.Available(ctx) {
		t.Error("expected adapter to be unavailable")
	}
}

// TestAdapter_ListSessions tests the ListSessions method
func TestAdapter_ListSessions(t *testing.T) {
	ctx := context.Background()
	adapter := testhelpers.NewMockAdapter("test")

	// Empty list
	sessions, err := adapter.ListSessions(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions, got %d", len(sessions))
	}

	// Add terminals
	term1 := testhelpers.NewMockTerminal("term-1", "Terminal 1", "content1")
	term2 := testhelpers.NewMockTerminal("term-2", "Terminal 2", "content2")
	adapter.AddTerminal(term1)
	adapter.AddTerminal(term2)

	sessions, err = adapter.ListSessions(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}

	// Verify session details
	if sessions[0].ID != "term-1" {
		t.Errorf("expected first session ID 'term-1', got %q", sessions[0].ID)
	}
	if sessions[1].Name != "Terminal 2" {
		t.Errorf("expected second session name 'Terminal 2', got %q", sessions[1].Name)
	}
}

// TestAdapter_Capture tests the Capture method
func TestAdapter_Capture(t *testing.T) {
	ctx := context.Background()
	adapter := testhelpers.NewMockAdapter("test")

	term := testhelpers.NewMockTerminal("term-1", "Terminal 1", "initial content")
	adapter.AddTerminal(term)

	// Capture initial content
	content, err := adapter.Capture(ctx, "term-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if content != "initial content" {
		t.Errorf("expected 'initial content', got %q", content)
	}

	// Change content and capture again
	term.SimulateContentChange("updated content")
	content, err = adapter.Capture(ctx, "term-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if content != "updated content" {
		t.Errorf("expected 'updated content', got %q", content)
	}

	// Capture from non-existent terminal
	_, err = adapter.Capture(ctx, "non-existent")
	if err == nil {
		t.Error("expected error for non-existent terminal")
	}
}

// TestAdapter_GetDimensions tests the GetDimensions method
func TestAdapter_GetDimensions(t *testing.T) {
	ctx := context.Background()
	adapter := testhelpers.NewMockAdapter("test")

	term := testhelpers.NewMockTerminal("term-1", "Terminal 1", "content")
	adapter.AddTerminal(term)

	// Get initial dimensions
	rows, cols, err := adapter.GetDimensions(ctx, "term-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if rows != 24 || cols != 80 {
		t.Errorf("expected 24x80, got %dx%d", rows, cols)
	}

	// Change dimensions
	term.SimulateDimensionChange(40, 120)
	rows, cols, err = adapter.GetDimensions(ctx, "term-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if rows != 40 || cols != 120 {
		t.Errorf("expected 40x120, got %dx%d", rows, cols)
	}

	// Get dimensions from non-existent terminal
	_, _, err = adapter.GetDimensions(ctx, "non-existent")
	if err == nil {
		t.Error("expected error for non-existent terminal")
	}
}

// TestAdapter_Subscribe tests the Subscribe method
func TestAdapter_Subscribe(t *testing.T) {
	ctx := context.Background()
	adapter := testhelpers.NewMockAdapter("test")

	term := testhelpers.NewMockTerminal("term-1", "Terminal 1", "content")
	adapter.AddTerminal(term)

	// Subscribe to events
	events, err := adapter.Subscribe(ctx, "term-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Trigger content change
	term.SimulateContentChange("new content")

	// Wait for event
	select {
	case event := <-events:
		if event.Type != terminal.EventContentUpdate {
			t.Errorf("expected ContentUpdate event, got %v", event.Type)
		}
		if event.Data != "new content" {
			t.Errorf("expected 'new content', got %q", event.Data)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for event")
	}

	// Subscribe to non-existent terminal
	_, err = adapter.Subscribe(ctx, "non-existent")
	if err == nil {
		t.Error("expected error for non-existent terminal")
	}
}

// TestAdapter_Stop tests the Stop method
func TestAdapter_Stop(t *testing.T) {
	ctx := context.Background()
	adapter := testhelpers.NewMockAdapter("test")

	term := testhelpers.NewMockTerminal("term-1", "Terminal 1", "content")
	adapter.AddTerminal(term)

	// Subscribe
	events, _ := adapter.Subscribe(ctx, "term-1")

	// Stop adapter
	adapter.Stop()

	// Channel should be closed
	select {
	case _, ok := <-events:
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for channel close")
	}
}
