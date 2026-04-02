package session

import (
	"path/filepath"
	"testing"
)

func TestStore_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	dbPath := filepath.Join(t.TempDir(), "integration.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	defer store.Close()

	// Test concurrent access
	err = store.AddSession("session-1", "terminal-1")
	if err != nil {
		t.Errorf("AddSession failed: %v", err)
	}

	commands, err := store.GetCommandsByTerminal("terminal-1")
	if err != nil {
		t.Errorf("GetCommandsByTerminal failed: %v", err)
	}
	if len(commands) != 0 {
		t.Errorf("Expected 0 commands, got %d", len(commands))
	}
}
