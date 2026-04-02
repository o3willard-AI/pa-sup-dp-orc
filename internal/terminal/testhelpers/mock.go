// Package testhelpers provides testing utilities for terminal adapter testing
package testhelpers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pairadmin/pairadmin/internal/terminal"
)

// MockTerminal simulates a terminal for testing
type MockTerminal struct {
	ID      string
	Name    string
	Content string
	Rows    int
	Cols    int
	Events  chan terminal.TerminalEvent
	mu      sync.RWMutex
}

// NewMockTerminal creates a mock terminal with configurable behavior
func NewMockTerminal(id, name string, content string) *MockTerminal {
	return &MockTerminal{
		ID:      id,
		Name:    name,
		Content: content,
		Rows:    24,
		Cols:    80,
		Events:  make(chan terminal.TerminalEvent, 10),
	}
}

// SimulateContentChange updates content and emits event
func (m *MockTerminal) SimulateContentChange(newContent string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Content = newContent
	m.Events <- terminal.TerminalEvent{
		Type:       terminal.EventContentUpdate,
		TerminalID: m.ID,
		Timestamp:  time.Now(),
		Data:       newContent,
	}
}

// SimulateDimensionChange updates dimensions and emits event
func (m *MockTerminal) SimulateDimensionChange(rows, cols int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Rows = rows
	m.Cols = cols
	m.Events <- terminal.TerminalEvent{
		Type:       terminal.EventDimensionChange,
		TerminalID: m.ID,
		Timestamp:  time.Now(),
		Data:       fmt.Sprintf("%dx%d", rows, cols),
	}
}

// MockAdapter implements TerminalAdapter for testing
type MockAdapter struct {
	name        string
	available   bool
	terminals   []*MockTerminal
	mu          sync.RWMutex
	subscribers map[string]chan terminal.TerminalEvent
}

// NewMockAdapter creates a new mock adapter for testing
func NewMockAdapter(name string) *MockAdapter {
	return &MockAdapter{
		name:        name,
		available:   true,
		terminals:   make([]*MockTerminal, 0),
		subscribers: make(map[string]chan terminal.TerminalEvent),
	}
}

// AddTerminal adds a mock terminal to the adapter
func (m *MockAdapter) AddTerminal(t *MockTerminal) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.terminals = append(m.terminals, t)
}

// SetAvailable sets the adapter availability
func (m *MockAdapter) SetAvailable(available bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.available = available
}

// Name returns the adapter name
func (m *MockAdapter) Name() string { return m.name }

// Available checks if the adapter is available
func (m *MockAdapter) Available(ctx context.Context) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.available
}

// ListSessions returns all mock terminals
func (m *MockAdapter) ListSessions(ctx context.Context) ([]terminal.DetectedTerminal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]terminal.DetectedTerminal, len(m.terminals))
	for i, t := range m.terminals {
		result[i] = terminal.DetectedTerminal{
			ID:       t.ID,
			Name:     t.Name,
			Adapter:  m.name,
			Type:     "mock",
			IsActive: true,
		}
	}
	return result, nil
}

// Capture returns the content of a mock terminal
func (m *MockAdapter) Capture(ctx context.Context, terminalID string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, t := range m.terminals {
		if t.ID == terminalID {
			t.mu.RLock()
			content := t.Content
			t.mu.RUnlock()
			return content, nil
		}
	}
	return "", fmt.Errorf("terminal not found: %s", terminalID)
}

// Subscribe starts streaming terminal events
func (m *MockAdapter) Subscribe(ctx context.Context, terminalID string) (<-chan terminal.TerminalEvent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var targetTerminal *MockTerminal
	for _, t := range m.terminals {
		if t.ID == terminalID {
			targetTerminal = t
			break
		}
	}
	if targetTerminal == nil {
		return nil, fmt.Errorf("terminal not found: %s", terminalID)
	}

	events := make(chan terminal.TerminalEvent, 10)
	m.subscribers[terminalID] = events

	go func() {
		for event := range targetTerminal.Events {
			select {
			case <-ctx.Done():
				return
			case events <- event:
			}
		}
	}()

	return events, nil
}

// GetDimensions returns terminal dimensions
func (m *MockAdapter) GetDimensions(ctx context.Context, terminalID string) (int, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, t := range m.terminals {
		if t.ID == terminalID {
			t.mu.RLock()
			rows := t.Rows
			cols := t.Cols
			t.mu.RUnlock()
			return rows, cols, nil
		}
	}
	return 0, 0, fmt.Errorf("terminal not found: %s", terminalID)
}

// Stop stops the adapter
func (m *MockAdapter) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, ch := range m.subscribers {
		close(ch)
	}
	m.subscribers = make(map[string]chan terminal.TerminalEvent)
}
