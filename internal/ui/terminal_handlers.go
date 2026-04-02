// Package ui provides Wails backend handlers for UI features like terminals.
package ui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/pairadmin/pairadmin/internal/terminal"
)

// TerminalHandlers manages multiple terminal sessions with real-time streaming to frontend.
type TerminalHandlers struct {
	detector *terminal.Detector
	sessions map[string]*TerminalSession
	mu       sync.RWMutex
}

// TerminalSession represents a single terminal session
type TerminalSession struct {
	ID       string
	Terminal *terminal.DetectedTerminal
	Cmd      *exec.Cmd
	Stdin    io.WriteCloser
}

// NewTerminalHandlers creates and returns handlers bound for Wails
func NewTerminalHandlers(d *terminal.Detector) *TerminalHandlers {
	return &TerminalHandlers{
		detector: d,
		sessions: make(map[string]*TerminalSession),
	}
}

// ListSessions returns active session IDs for frontend tabs
func (t *TerminalHandlers) ListSessions(ctx context.Context) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	ids := make([]string, 0, len(t.sessions))
	for id := range t.sessions {
		ids = append(ids, id)
	}
	return ids
}

// GetDetectedTerminals returns detected terminals from the detector
func (t *TerminalHandlers) GetDetectedTerminals() map[string]*terminal.DetectedTerminal {
	return t.detector.GetSessions()
}

// StartSession starts a bash shell session for a terminal
func (t *TerminalHandlers) StartSession(terminalID string, ctx context.Context) (string, error) {
	sessions := t.detector.GetSessions()
	term, exists := sessions[terminalID]
	if !exists {
		return "", fmt.Errorf("terminal %q not found", terminalID)
	}

	sessionID := fmt.Sprintf("%s-session", terminalID)

	t.mu.Lock()
	if _, exists := t.sessions[sessionID]; exists {
		t.mu.Unlock()
		return "", fmt.Errorf("session %q already exists", sessionID)
	}
	t.mu.Unlock()

	cmd := exec.Command("bash")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("start cmd: %w", err)
	}

	session := &TerminalSession{
		ID:       sessionID,
		Terminal: term,
		Cmd:      cmd,
		Stdin:    stdin,
	}

	t.mu.Lock()
	t.sessions[sessionID] = session
	t.mu.Unlock()

	go t.streamPipe(sessionID, stdout, "stdout")
	go t.streamPipe(sessionID, stderr, "stderr")

	runtime.EventsEmit(ctx, "terminal-session-started", map[string]any{
		"id":         sessionID,
		"terminalID": terminalID,
		"name":       term.Name,
	})
	return sessionID, nil
}

// streamPipe reads lines from pipe and emits to frontend
func (t *TerminalHandlers) streamPipe(sessionID string, r io.ReadCloser, streamType string) {
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		runtime.EventsEmit(context.Background(), "terminal-output", map[string]string{
			"sessionID": sessionID,
			"stream":    streamType,
			"line":      line,
		})
	}
}

// SendInput forwards user input to session stdin
func (t *TerminalHandlers) SendInput(sessionID, input string, ctx context.Context) error {
	t.mu.RLock()
	session, exists := t.sessions[sessionID]
	t.mu.RUnlock()
	if !exists {
		return fmt.Errorf("session %q not found", sessionID)
	}
	_, err := session.Stdin.Write([]byte(input + "\n"))
	if err != nil {
		return fmt.Errorf("write input: %w", err)
	}
	return nil
}

// SelectSession selects a session for focus
func (t *TerminalHandlers) SelectSession(sessionID string, ctx context.Context) error {
	t.mu.RLock()
	_, exists := t.sessions[sessionID]
	t.mu.RUnlock()
	if !exists {
		return fmt.Errorf("session %q not found", sessionID)
	}
	runtime.EventsEmit(ctx, "terminal-session-selected", sessionID)
	return nil
}

// SwitchTab emits event for tab switching
func (t *TerminalHandlers) SwitchTab(tabID string, ctx context.Context) {
	runtime.EventsEmit(ctx, "terminal-tab-switched", tabID)
}

// CloseSession kills a session and cleans up
func (t *TerminalHandlers) CloseSession(sessionID string, ctx context.Context) error {
	t.mu.Lock()
	session, exists := t.sessions[sessionID]
	delete(t.sessions, sessionID)
	t.mu.Unlock()
	if !exists {
		return fmt.Errorf("session %q not found", sessionID)
	}
	session.Cmd.Process.Kill()
	runtime.EventsEmit(ctx, "terminal-session-closed", sessionID)
	return nil
}
