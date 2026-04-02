package tmux

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pairadmin/pairadmin/internal/terminal"
)

// Adapter implements the TerminalAdapter interface for tmux
type Adapter struct {
	config    Config
	connected bool
	mu        sync.RWMutex
	cmdRunner CommandRunner
}

// Config holds tmux adapter configuration
type Config struct {
	PollIntervalMs int
	TmuxBinary     string
	SocketName     string
}

// DefaultConfig returns default tmux configuration
func DefaultConfig() Config {
	return Config{
		PollIntervalMs: 500,
		TmuxBinary:     "tmux",
		SocketName:     "",
	}
}

// CommandRunner interface for testability
type CommandRunner interface {
	Run(ctx context.Context, name string, args ...string) (string, error)
}

// DefaultCommandRunner is the default command runner implementation
type DefaultCommandRunner struct{}

func (d *DefaultCommandRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// NewAdapter creates a new tmux adapter
func NewAdapter(config Config) *Adapter {
	return &Adapter{
		config:    config,
		connected: false,
		cmdRunner: &DefaultCommandRunner{},
	}
}

// Name returns the adapter name
func (a *Adapter) Name() string {
	return "tmux"
}

// Available checks if tmux is installed and accessible
func (a *Adapter) Available(ctx context.Context) bool {
	_, err := a.runCommand(ctx, "-V")
	return err == nil
}

// Connect establishes connection to tmux
func (a *Adapter) Connect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.Available(ctx) {
		return terminal.ErrAdapterNotAvailable{
			AdapterName: "tmux",
			Reason:      "tmux binary not found in PATH",
		}
	}

	a.connected = true
	return nil
}

// Disconnect closes tmux connection
func (a *Adapter) Disconnect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.connected = false
	return nil
}

// IsConnected returns connection status
func (a *Adapter) IsConnected() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.connected
}

// ListSessions returns all active tmux sessions
func (a *Adapter) ListSessions(ctx context.Context) ([]terminal.DetectedTerminal, error) {
	output, err := a.runCommand(ctx, "list-sessions", "-F", "#{session_name},#{session_created},#{session_activity}")
	if err != nil {
		return nil, err
	}

	var sessions []terminal.DetectedTerminal
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}

		created, _ := strconv.ParseInt(parts[1], 10, 64)
		sessions = append(sessions, terminal.DetectedTerminal{
			ID:       parts[0],
			Name:     parts[0],
			Adapter:  "tmux",
			Type:     "tmux",
			Created:  time.Unix(created, 0),
			IsActive: true,
		})
	}

	return sessions, nil
}

// ListPanes returns all panes in a session
func (a *Adapter) ListPanes(ctx context.Context, sessionID string) ([]terminal.DetectedTerminal, error) {
	format := "#{pane_id},#{session_name},#{window_index},#{pane_index},#{pane_current_command},#{pane_current_path}"
	output, err := a.runCommand(ctx, "list-panes", "-t", sessionID, "-F", format)
	if err != nil {
		return nil, err
	}

	var panes []terminal.DetectedTerminal
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 6 {
			continue
		}

		panes = append(panes, terminal.DetectedTerminal{
			ID:         parts[0],
			Name:       fmt.Sprintf("%s:%s.%s", parts[1], parts[2], parts[3]),
			Adapter:    "tmux",
			Type:       "tmux-pane",
			Command:    parts[4],
			WorkingDir: parts[5],
			IsActive:   true,
		})
	}

	return panes, nil
}

// Capture captures the current content of a tmux pane
func (a *Adapter) Capture(ctx context.Context, paneID string) (string, error) {
	return a.runCommand(ctx, "capture-pane", "-p", "-t", paneID)
}

// Subscribe starts streaming terminal events (polling-based for tmux)
func (a *Adapter) Subscribe(ctx context.Context, paneID string) (<-chan terminal.TerminalEvent, error) {
	events := make(chan terminal.TerminalEvent)

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
				content, err := a.Capture(ctx, paneID)
				if err != nil {
					events <- terminal.TerminalEvent{
						Type:       terminal.EventActivityChange,
						TerminalID: paneID,
						Timestamp:  time.Now(),
						Data:       fmt.Sprintf("capture error: %v", err),
					}
					continue
				}

				if content != lastContent {
					events <- terminal.TerminalEvent{
						Type:       terminal.EventContentUpdate,
						TerminalID: paneID,
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
func (a *Adapter) GetDimensions(ctx context.Context, paneID string) (int, int, error) {
	output, err := a.runCommand(ctx, "display-message", "-p", "-t", paneID, "#{pane_width},#{pane_height}")
	if err != nil {
		return 0, 0, err
	}

	parts := strings.Split(output, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid dimensions output: %s", output)
	}

	cols, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	rows, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return rows, cols, nil
}

// SendCommand sends a command to a tmux pane
func (a *Adapter) SendCommand(ctx context.Context, paneID string, command string) error {
	_, err := a.runCommand(ctx, "send-keys", "-t", paneID, command, "Enter")
	return err
}

// runCommand executes a tmux command and returns output
func (a *Adapter) runCommand(ctx context.Context, args ...string) (string, error) {
	allArgs := []string{}
	if a.config.SocketName != "" {
		allArgs = append(allArgs, "-L", a.config.SocketName)
	}
	allArgs = append(allArgs, args...)

	return a.cmdRunner.Run(ctx, a.config.TmuxBinary, allArgs...)
}

// parseSessionList parses tmux list-sessions output
func parseSessionList(output string) []map[string]string {
	var sessions []map[string]string
	re := regexp.MustCompile(`(\w+): (\d+) windows`)

	for _, line := range strings.Split(output, "\n") {
		if matches := re.FindStringSubmatch(line); matches != nil {
			sessions = append(sessions, map[string]string{
				"name":    matches[1],
				"windows": matches[2],
			})
		}
	}

	return sessions
}
