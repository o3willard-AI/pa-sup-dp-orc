//go:build linux
// +build linux

// Package terminal provides Linux terminal adapter using AT-SPI2
package terminal

/*
#cgo pkg-config: atspi-2
#include <stdlib.h>
#include "atspi2.h"
*/
import "C"
import (
	"context"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// Adapter implements TerminalAdapter for Linux terminals via AT-SPI2
type Adapter struct {
	config     Config
	mu         sync.RWMutex
	running    bool
	cancelFunc context.CancelFunc
}

// Config holds Linux adapter configuration
type Config struct {
	PollIntervalMs int
}

// DefaultConfig returns default adapter configuration
func DefaultConfig() Config {
	return Config{PollIntervalMs: 500}
}

// NewAdapter creates a new Linux Terminal adapter
func NewAdapter(config Config) *Adapter {
	return &Adapter{config: config}
}

// Name returns the adapter name
func (a *Adapter) Name() string { return "linux" }

// Available checks if AT-SPI2 is available
func (a *Adapter) Available(ctx context.Context) bool {
	return bool(C.pa_atspi_available())
}

// ListSessions returns all detected terminal windows
func (a *Adapter) ListSessions(ctx context.Context) ([]DetectedTerminal, error) {
	var count C.int
	cSessions := C.pa_atspi_list_sessions(&count)
	if cSessions == nil {
		return nil, nil
	}
	defer C.pa_atspi_free_sessions(cSessions)

	sessions := unsafe.Slice(cSessions, int(count))
	var terminals []DetectedTerminal

	for i := 0; i < int(count); i++ {
		name := C.GoString(sessions[i])
		terminals = append(terminals, DetectedTerminal{
			ID:       fmt.Sprintf("linux-%d-%s", i, name),
			Name:     name,
			Adapter:  "linux",
			Type:     "AT-SPI2",
			IsActive: true,
		})
	}
	return terminals, nil
}

// Capture captures the current content of a terminal window
func (a *Adapter) Capture(ctx context.Context, terminalID string) (string, error) {
	// Pass terminal ID as object path (simplified)
	cPath := C.CString(terminalID)
	defer C.free(unsafe.Pointer(cPath))

	cText := C.pa_atspi_capture(cPath)
	if cText == nil {
		return "", fmt.Errorf("failed to capture text")
	}
	defer C.pa_atspi_free_capture(cText)

	return C.GoString(cText), nil
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
	cPath := C.CString(terminalID)
	defer C.free(unsafe.Pointer(cPath))

	width := int(C.pa_atspi_get_width(cPath))
	height := int(C.pa_atspi_get_height(cPath))

	// Default if zero
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 24
	}

	return height, width, nil
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
