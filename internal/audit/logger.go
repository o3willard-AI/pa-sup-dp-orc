package audit

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Event represents an auditable system event.
// Fields:
//
//	Timestamp: when the event occurred (required)
//	Action: type of action (e.g., "config_save", "auth_attempt") (required)
//	User: user identifier, if applicable
//	Terminal: terminal session ID, if applicable
//	Details: human-readable description of the event
//	Success: whether the action succeeded
//	Error: error message if the action failed
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	User      string    `json:"user,omitempty"`
	Terminal  string    `json:"terminal,omitempty"`
	Details   string    `json:"details,omitempty"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes audit events to JSONL files.
type Logger struct {
	mu     sync.Mutex
	file   *os.File
	writer *json.Encoder
	closed bool
}

// NewLogger creates an audit logger writing to ~/.pairadmin/logs/audit-YYYY-MM-DD.jsonl
func NewLogger() (*Logger, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	logDir := filepath.Join(configDir, "pairadmin", "logs")
	if err := os.MkdirAll(logDir, 0700); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	filename := fmt.Sprintf("audit-%s.jsonl", time.Now().Format("2006-01-02"))
	path := filepath.Join(logDir, filename)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open audit log: %w", err)
	}

	return &Logger{
		file:   file,
		writer: json.NewEncoder(file),
	}, nil
}

// Log writes an audit event.
func (l *Logger) Log(event Event) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return errors.New("audit logger closed")
	}
	if err := l.writer.Encode(event); err != nil {
		return fmt.Errorf("encode audit event: %w", err)
	}
	return nil
}

// Close closes the audit log file.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return nil // already closed
	}
	err := l.file.Close()
	if err == nil {
		l.closed = true
	}
	return err
}
