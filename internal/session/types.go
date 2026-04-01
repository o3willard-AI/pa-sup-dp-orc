package session

import (
	"time"
)

// SuggestedCommand represents a command suggested by the AI.
type SuggestedCommand struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	TerminalID  string    `json:"terminal_id"`
	Command     string    `json:"command"`
	Description string    `json:"description"`
	Context     string    `json:"context"` // Original user question or context
	CreatedAt   time.Time `json:"created_at"`
	UsedCount   int       `json:"used_count"`
	LastUsedAt  time.Time `json:"last_used_at"`
}

// Session represents a terminal session with its command history.
type Session struct {
	ID         string             `json:"id"`
	TerminalID string             `json:"terminal_id"`
	CreatedAt  time.Time          `json:"created_at"`
	Commands   []SuggestedCommand `json:"commands"`
}
