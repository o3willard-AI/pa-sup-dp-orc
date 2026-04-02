//go:build darwin
// +build darwin

// Package macos provides macOS Terminal.app adapter implementation
package macos

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/pairadmin/pairadmin/internal/terminal"
)

// Adapter implements terminal.TerminalAdapter for macOS Terminal.app
type Adapter struct {
	config     Config
	mu         sync.RWMutex
	running    bool
	cancelFunc context.CancelFunc
}

// Config holds macOS adapter configuration
type Config struct {
	PollIntervalMs int
}

// DefaultConfig returns default adapter configuration
func DefaultConfig() Config {
	return Config{
		PollIntervalMs: 500,
	}
}

// NewAdapter creates a new macOS Terminal adapter
func NewAdapter(config Config) *Adapter {
	return &Adapter{
		config: config,
	}
}

// Name returns the adapter name
func (a *Adapter) Name() string {
	return "macos"
}

// Available checks if the adapter is available (macOS + accessibility permission)
func (a *Adapter) Available(ctx context.Context) bool {
	if runtime.GOOS != "darwin" {
		return false
	}
	return IsAccessibilityEnabled()
}

// ListSessions returns all detected Terminal.app windows
func (a *Adapter) ListSessions(ctx context.Context) ([]terminal.DetectedTerminal, error) {
	titles, err := GetTerminalWindowTitles()
	if err != nil {
		return nil, fmt.Errorf("get titles: %w", err)
	}

	pids, err := GetTerminalWindowPids()
	if err != nil {
		return nil, fmt.Errorf("get pids: %w", err)
	}

	var sessions []terminal.DetectedTerminal
	for i, title := range titles {
		pid := 0
		if i < len(pids) {
			pid = pids[i]
		}

		sessions = append(sessions, terminal.DetectedTerminal{
			ID:       fmt.Sprintf("macos-%d-%s", pid, title),
			Name:     title,
			Adapter:  "macos",
			Type:     "Terminal.app",
			IsActive: true,
			Command:  fmt.Sprintf("pid:%d", pid),
		})
	}

	if len(sessions) == 0 {
		return nil, terminal.ErrNoTerminalWindows{}
	}

	return sessions, nil
}

// Capture captures the current content of a Terminal.app window
func (a *Adapter) Capture(ctx context.Context, terminalID string) (string, error) {
	content, err := ExtractTextFromFrontmostTerminal()
	if err == nil && content != "" {
		return content, nil
	}

	var pid int
	_, err = fmt.Sscanf(terminalID, "macos-%d-", &pid)
	if err != nil || pid == 0 {
		return "", fmt.Errorf("parse terminal ID: %w", err)
	}

	content, err = ExtractTextFromWindow(pid)
	if err != nil {
		return "", fmt.Errorf("extract text: %w", err)
	}

	return content, nil
}

// Subscribe starts streaming terminal events for a session
func (a *Adapter) Subscribe(ctx context.Context, terminalID string) (<-chan terminal.TerminalEvent, error) {
	events := make(chan terminal.TerminalEvent)

	ctx, cancel := context.WithCancel(ctx)
	a.mu.Lock()
	a.cancelFunc = cancel
	a.running = true
	a.mu.Unlock()

	go func() {
		defer close(events)
		ticker := time.NewTicker(time.Duration(a.config.PollIntervalMs) * time.Millisecond)
		defer ticker.Stop()

		var lastContent string

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				content, err := a.Capture(ctx, terminalID)
				if err != nil {
					events <- terminal.TerminalEvent{
						Type:       terminal.EventActivityChange,
						TerminalID: terminalID,
						Timestamp:  time.Now(),
						Data:       fmt.Sprintf("capture error: %v", err),
					}
					continue
				}

				if content != lastContent {
					events <- terminal.TerminalEvent{
						Type:       terminal.EventContentUpdate,
						TerminalID: terminalID,
						Timestamp:  time.Now(),
						Data:       content,
					}
					lastContent = content
				}
			}
		}
	}()

	return events, nil
}

// GetDimensions returns the terminal dimensions (rows, cols)
func (a *Adapter) GetDimensions(ctx context.Context, terminalID string) (int, int, error) {
	content, err := a.Capture(ctx, terminalID)
	if err != nil {
		return 0, 0, err
	}

	lines := 0
	maxCols := 0
	for _, line := range splitLines(content) {
		lines++
		if len(line) > maxCols {
			maxCols = len(line)
		}
	}

	if lines == 0 {
		lines = 24
		maxCols = 80
	}

	return lines, maxCols, nil
}

func splitLines(content string) []string {
	var lines []string
	current := ""
	for _, r := range content {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// Stop stops the adapter
func (a *Adapter) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cancelFunc != nil {
		a.cancelFunc()
	}
	a.running = false
}

// IsRunning returns whether the adapter is currently running
func (a *Adapter) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.running
}
