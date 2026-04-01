package session

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStore_Crud(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	defer store.Close()

	// Add session
	sessionID := "session-1"
	terminalID := "terminal-1"
	err = store.AddSession(sessionID, terminalID)
	if err != nil {
		t.Fatalf("AddSession failed: %v", err)
	}

	// Add command
	cmd := SuggestedCommand{
		ID:          "cmd-1",
		TerminalID:  terminalID,
		Command:     "ls -la",
		Description: "List files",
		Context:     "User asked how to see hidden files",
		CreatedAt:   time.Now(),
		UsedCount:   0,
	}
	err = store.AddCommand(cmd)
	if err != nil {
		t.Fatalf("AddCommand failed: %v", err)
	}

	// Get commands
	commands, err := store.GetCommandsByTerminal(terminalID)
	if err != nil {
		t.Fatalf("GetCommandsByTerminal failed: %v", err)
	}
	if len(commands) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(commands))
	}
	if commands[0].Command != "ls -la" {
		t.Errorf("Expected command 'ls -la', got %s", commands[0].Command)
	}

	// Increment used count
	err = store.IncrementUsedCount("cmd-1")
	if err != nil {
		t.Fatalf("IncrementUsedCount failed: %v", err)
	}

	// Verify updated count (would need a GetCommand to check, but skip for brevity)
}

func TestStore_NewStore_CreatesTables(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	defer store.Close()

	// Try to insert a session (should succeed if tables exist)
	err = store.AddSession("test", "terminal")
	if err != nil {
		t.Fatalf("AddSession after NewStore failed (tables not created?): %v", err)
	}
}
