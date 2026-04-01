package session

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// isConstraintError returns true if err appears to be a SQLite constraint violation.
func isConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE") ||
		strings.Contains(msg, "FOREIGN KEY") ||
		strings.Contains(msg, "constraint") ||
		strings.Contains(msg, "duplicate")
}

// newTestStore creates a temporary store for testing.
func newTestStore(t *testing.T) (*Store, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	return store, func() { store.Close() }
}

func TestStore_Crud(t *testing.T) {
	store, cleanup := newTestStore(t)
	defer cleanup()
	var err error

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
		SessionID:   sessionID,
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
	store, cleanup := newTestStore(t)
	defer cleanup()
	var err error

	// Try to insert a session (should succeed if tables exist)
	err = store.AddSession("test", "terminal")
	if err != nil {
		t.Fatalf("AddSession after NewStore failed (tables not created?): %v", err)
	}
}

func TestStore_AddSession_DuplicateID(t *testing.T) {
	store, cleanup := newTestStore(t)
	defer cleanup()
	var err error

	// Add first session
	err = store.AddSession("session-1", "terminal-1")
	if err != nil {
		t.Fatalf("AddSession failed: %v", err)
	}

	// Try to add duplicate session ID
	err = store.AddSession("session-1", "terminal-2")
	if err == nil {
		t.Error("Expected error for duplicate session ID, got nil")
	} else if !isConstraintError(err) {
		t.Errorf("Expected constraint violation error, got: %v", err)
	}
}

func TestStore_AddCommand_DuplicateID(t *testing.T) {
	store, cleanup := newTestStore(t)
	defer cleanup()
	var err error

	// Need a session first
	err = store.AddSession("session-1", "terminal-1")
	if err != nil {
		t.Fatalf("AddSession failed: %v", err)
	}

	cmd := SuggestedCommand{
		ID:          "cmd-1",
		SessionID:   "session-1",
		TerminalID:  "terminal-1",
		Command:     "ls -la",
		Description: "List files",
		Context:     "test",
		CreatedAt:   time.Now(),
		UsedCount:   0,
	}
	err = store.AddCommand(cmd)
	if err != nil {
		t.Fatalf("AddCommand failed: %v", err)
	}

	// Try duplicate command ID
	err = store.AddCommand(cmd)
	if err == nil {
		t.Error("Expected error for duplicate command ID, got nil")
	} else if !isConstraintError(err) {
		t.Errorf("Expected constraint violation error, got: %v", err)
	}
}

func TestStore_AddCommand_ForeignKeyViolation(t *testing.T) {
	store, cleanup := newTestStore(t)
	defer cleanup()
	var err error

	// Do NOT create a session, try to add command with non-existent session ID
	cmd := SuggestedCommand{
		ID:          "cmd-1",
		SessionID:   "non-existent-session",
		TerminalID:  "terminal-1",
		Command:     "ls -la",
		Description: "List files",
		Context:     "test",
		CreatedAt:   time.Now(),
		UsedCount:   0,
	}
	err = store.AddCommand(cmd)
	if err == nil {
		t.Error("Expected foreign key violation error, got nil")
	} else if !isConstraintError(err) {
		t.Errorf("Expected constraint violation error, got: %v", err)
	}
}

func TestStore_DeleteSession_Cascade(t *testing.T) {
	store, cleanup := newTestStore(t)
	defer cleanup()
	var err error

	// Create session
	sessionID := "session-1"
	terminalID := "terminal-1"
	err = store.AddSession(sessionID, terminalID)
	if err != nil {
		t.Fatalf("AddSession failed: %v", err)
	}

	// Add command to session
	cmd := SuggestedCommand{
		ID:          "cmd-1",
		SessionID:   sessionID,
		TerminalID:  terminalID,
		Command:     "ls -la",
		Description: "List files",
		Context:     "test",
		CreatedAt:   time.Now(),
		UsedCount:   0,
	}
	err = store.AddCommand(cmd)
	if err != nil {
		t.Fatalf("AddCommand failed: %v", err)
	}

	// Verify command exists
	commands, err := store.GetCommandsByTerminal(terminalID)
	if err != nil {
		t.Fatalf("GetCommandsByTerminal failed: %v", err)
	}
	if len(commands) != 1 {
		t.Fatalf("Expected 1 command before delete, got %d", len(commands))
	}

	// Delete session
	err = store.DeleteSession(sessionID)
	if err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}

	// Verify command is also deleted
	commands, err = store.GetCommandsByTerminal(terminalID)
	if err != nil {
		t.Fatalf("GetCommandsByTerminal after delete failed: %v", err)
	}
	if len(commands) != 0 {
		t.Fatalf("Expected 0 commands after cascade delete, got %d", len(commands))
	}
}
