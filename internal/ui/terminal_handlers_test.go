package ui

import (
	"context"
	"testing"

	"github.com/pairadmin/pairadmin/internal/terminal"
)

func TestNewTerminalHandlers(t *testing.T) {
	det := terminal.NewDetector(terminal.DefaultDetectorConfig())
	h := NewTerminalHandlers(det)
	if h == nil {
		t.Fatal("NewTerminalHandlers returned nil")
	}
	if h.sessions == nil {
		t.Fatal("sessions map not initialized")
	}
}

func TestListSessions(t *testing.T) {
	det := terminal.NewDetector(terminal.DefaultDetectorConfig())
	h := NewTerminalHandlers(det)
	sessions := h.ListSessions(context.Background())
	if len(sessions) != 0 {
		t.Errorf("Expected 0 sessions, got %d", len(sessions))
	}
}
