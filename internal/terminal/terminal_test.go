package terminal

import "testing"

func TestDetectedTerminal(t *testing.T) {
	term := DetectedTerminal{
		ID:         "test-1",
		Name:       "Test Session",
		Adapter:    "tmux",
		Type:       "tmux",
		IsActive:   true,
		WorkingDir: "/home/user",
		Command:    "bash",
	}

	if term.ID != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", term.ID)
	}
	if !term.IsActive {
		t.Error("Expected IsActive to be true")
	}
}

func TestTerminalEventType(t *testing.T) {
	tests := []struct {
		name     string
		event    TerminalEventType
		expected string
	}{
		{"Content Update", EventContentUpdate, "content_update"},
		{"Dimension Change", EventDimensionChange, "dimension_change"},
		{"Session Created", EventSessionCreated, "session_created"},
		{"Session Closed", EventSessionClosed, "session_closed"},
		{"Activity Change", EventActivityChange, "activity_change"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.event) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.event)
			}
		})
	}
}
