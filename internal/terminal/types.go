package terminal

import (
	"context"
	"time"
)

// TerminalAdapter defines the contract for terminal multiplexer implementations.
// All terminal backends (tmux, screen, Windows Terminal) must implement this interface.
type TerminalAdapter interface {
	// Name returns the adapter name (e.g., "tmux", "screen")
	Name() string

	// Available checks if the terminal multiplexer is installed and accessible
	Available(ctx context.Context) bool

	// ListSessions returns all active terminal sessions
	ListSessions(ctx context.Context) ([]DetectedTerminal, error)

	// Capture captures the current content of a terminal session
	Capture(ctx context.Context, terminalID string) (string, error)

	// Subscribe starts streaming terminal events
	Subscribe(ctx context.Context, terminalID string) (<-chan TerminalEvent, error)

	// GetDimensions returns the terminal dimensions (rows, cols)
	GetDimensions(ctx context.Context, terminalID string) (int, int, error)
}

// DetectedTerminal represents a detected terminal session or window
type DetectedTerminal struct {
	// ID is the unique identifier for this terminal session
	ID string `json:"id"`

	// Name is the human-readable name (e.g., tmux session name, window title)
	Name string `json:"name"`

	// Adapter is the name of the adapter that detected this terminal
	Adapter string `json:"adapter"`

	// Type is the terminal type (e.g., "tmux", "screen", "Terminal.app", "PuTTY")
	Type string `json:"type"`

	// Created is when the session was created
	Created time.Time `json:"created"`

	// IsActive indicates if the terminal is currently active
	IsActive bool `json:"is_active"`

	// WorkingDir is the current working directory if available
	WorkingDir string `json:"working_dir,omitempty"`

	// Command is the running command if available
	Command string `json:"command,omitempty"`
}

// TerminalEvent represents an event from a terminal session
type TerminalEvent struct {
	// Type is the event type
	Type TerminalEventType `json:"type"`

	// TerminalID is the ID of the terminal that generated this event
	TerminalID string `json:"terminal_id"`

	// Timestamp is when the event occurred
	Timestamp time.Time `json:"timestamp"`

	// Data contains event-specific data (e.g., new content, dimension changes)
	Data string `json:"data,omitempty"`

	// Rows is the new number of rows (for dimension change events)
	Rows int `json:"rows,omitempty"`

	// Cols is the new number of columns (for dimension change events)
	Cols int `json:"cols,omitempty"`
}

// TerminalEventType defines the types of terminal events
type TerminalEventType string

const (
	// EventContentUpdate is emitted when terminal content changes
	EventContentUpdate TerminalEventType = "content_update"

	// EventDimensionChange is emitted when terminal dimensions change
	EventDimensionChange TerminalEventType = "dimension_change"

	// EventSessionCreated is emitted when a new terminal session is detected
	EventSessionCreated TerminalEventType = "session_created"

	// EventSessionClosed is emitted when a terminal session is closed
	EventSessionClosed TerminalEventType = "session_closed"

	// EventActivityChange is emitted when terminal activity state changes
	EventActivityChange TerminalEventType = "activity_change"
)
