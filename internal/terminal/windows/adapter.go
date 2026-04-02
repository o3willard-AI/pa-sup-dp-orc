//go:build windows

package terminal

/*
#cgo windows LDFLAGS: -lole32 -lcombase
#include "ui_automation.h"
*/
import "C"
import (
	"context"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// wcharToString converts a C wide string (UTF-16) to a Go string.
func wcharToString(wstr *C.wchar_t) string {
	if wstr == nil {
		return ""
	}
	// C.wchar_t is uint16 on Windows
	return windows.UTF16PtrToString((*uint16)(unsafe.Pointer(wstr)))
}

// Adapter implements TerminalAdapter for Windows terminals via UI Automation
type Adapter struct {
	config     Config
	mu         sync.RWMutex
	running    bool
	cancelFunc context.CancelFunc
}

// Config holds Windows adapter configuration
type Config struct {
	PollIntervalMs int
}

// DefaultConfig returns default adapter configuration
func DefaultConfig() Config {
	return Config{PollIntervalMs: 500}
}

// NewAdapter creates a new Windows Terminal adapter
func NewAdapter(config Config) *Adapter {
	return &Adapter{config: config}
}

// Name returns the adapter name
func (a *Adapter) Name() string { return "windows" }

// Available checks if UI Automation is available
func (a *Adapter) Available(ctx context.Context) bool {
	// Try to initialize UI Automation
	result := C.uia_initialize()
	if result != 0 {
		return false
	}

	available := C.uia_is_available() != 0
	if !available {
		C.uia_cleanup()
	}
	return available
}

// ListSessions returns all detected terminal windows
func (a *Adapter) ListSessions(ctx context.Context) ([]DetectedTerminal, error) {
	// Get window titles
	var titleCount C.int
	var titleError C.UIA_ErrorCode
	var titleErrorMessage *C.wchar_t
	defer func() {
		if titleErrorMessage != nil {
			C.free_wide_string(titleErrorMessage)
		}
	}()

	cTitles := C.get_terminal_window_titles(&titleCount, &titleError, &titleErrorMessage)
	if titleError != 0 {
		errorMsg := "unknown error"
		if titleErrorMessage != nil {
			errorMsg = wcharToString(titleErrorMessage)
		}
		return nil, fmt.Errorf("UI Automation error %d: %s", int(titleError), errorMsg)
	}
	if cTitles == nil {
		return nil, nil
	}
	defer C.free_wide_string_array(cTitles, titleCount)

	// Get window PIDs
	var pidCount C.int
	var pidError C.UIA_ErrorCode
	var pidErrorMessage *C.wchar_t
	defer func() {
		if pidErrorMessage != nil {
			C.free_wide_string(pidErrorMessage)
		}
	}()

	cPids := C.get_terminal_window_pids(&pidCount, &pidError, &pidErrorMessage)
	if pidError != 0 {
		errorMsg := "unknown error"
		if pidErrorMessage != nil {
			errorMsg = wcharToString(pidErrorMessage)
		}
		return nil, fmt.Errorf("UI Automation error %d: %s", int(pidError), errorMsg)
	}
	if cPids == nil {
		return nil, nil
	}
	defer C.free_int_array(cPids)

	// Guard against invalid counts (should never happen if error codes are zero)
	if titleCount < 0 || pidCount < 0 {
		return nil, fmt.Errorf("invalid count from UI Automation")
	}

	// Use minimum count in case arrays don't match (shouldn't happen)
	count := int(titleCount)
	if int(pidCount) < count {
		count = int(pidCount)
	}

	titles := unsafe.Slice(cTitles, int(titleCount))
	pids := unsafe.Slice(cPids, int(pidCount))
	var terminals []DetectedTerminal

	for i := 0; i < count; i++ {
		name := wcharToString(titles[i])
		pid := int(pids[i])
		terminals = append(terminals, DetectedTerminal{
			ID:       fmt.Sprintf("windows-%d-%s", pid, name),
			Name:     name,
			Adapter:  "windows",
			Type:     "UIA",
			IsActive: true,
		})
	}
	return terminals, nil
}

// Capture captures the current content of a terminal window
func (a *Adapter) Capture(ctx context.Context, terminalID string) (string, error) {
	// Extract PID from terminal ID (format: windows-pid-name)
	// Parse format: windows-{pid}-{name}
	var pid int
	n, err := fmt.Sscanf(terminalID, "windows-%d-", &pid)
	if err != nil || n != 1 {
		// Fallback: try to extract PID from focused terminal
		return a.captureFocusedTerminal()
	}

	var cError C.UIA_ErrorCode
	var cErrorMessage *C.wchar_t
	defer func() {
		if cErrorMessage != nil {
			C.free_wide_string(cErrorMessage)
		}
	}()

	cText := C.extract_text_from_window(C.int(pid), &cError, &cErrorMessage)
	if cError != 0 {
		errorMsg := "unknown error"
		if cErrorMessage != nil {
			errorMsg = wcharToString(cErrorMessage)
		}
		return "", fmt.Errorf("UI Automation error %d: %s", int(cError), errorMsg)
	}
	if cText == nil {
		return "", fmt.Errorf("failed to capture text")
	}
	defer C.free_wide_string(cText)

	return wcharToString(cText), nil
}

// Helper: Capture from focused terminal
func (a *Adapter) captureFocusedTerminal() (string, error) {
	var cError C.UIA_ErrorCode
	var cErrorMessage *C.wchar_t
	defer func() {
		if cErrorMessage != nil {
			C.free_wide_string(cErrorMessage)
		}
	}()

	cText := C.extract_text_from_focused_terminal(&cError, &cErrorMessage)
	if cError != 0 {
		errorMsg := "unknown error"
		if cErrorMessage != nil {
			errorMsg = wcharToString(cErrorMessage)
		}
		return "", fmt.Errorf("UI Automation error %d: %s", int(cError), errorMsg)
	}
	if cText == nil {
		return "", fmt.Errorf("failed to capture text from focused terminal")
	}
	defer C.free_wide_string(cText)

	return wcharToString(cText), nil
}

// Subscribe starts streaming terminal events
func (a *Adapter) Subscribe(ctx context.Context, terminalID string) (<-chan TerminalEvent, error) {
	events := make(chan TerminalEvent)

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
					continue
				}
				if content != lastContent {
					events <- TerminalEvent{
						Type:       EventContentUpdate,
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

// GetDimensions returns terminal dimensions (rows, cols)
func (a *Adapter) GetDimensions(ctx context.Context, terminalID string) (int, int, error) {
	// Default terminal dimensions
	return 24, 80, nil
}

// Stop stops the adapter
func (a *Adapter) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cancelFunc != nil {
		a.cancelFunc()
	}
	a.running = false

	// Cleanup UI Automation
	C.uia_cleanup()
}
