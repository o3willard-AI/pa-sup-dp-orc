// Copyright 2026 The PairAdmin Authors
// SPDX-License-Identifier: MIT

package terminal

import (
	"context"
	"sync"
	"time"
)

// Detector scans for terminal sessions and emits lifecycle events.
type Detector struct {
	config DetectorConfig

	mu      sync.RWMutex
	running bool
	ctx     context.Context
	cancel  context.CancelFunc

	sessions map[string]*DetectedTerminal
	events   chan TerminalEvent

	adapters []TerminalAdapter
}

// DetectorConfig holds detector configuration.
type DetectorConfig struct {
	PollInterval time.Duration
	MaxRetries   int
}

// DefaultDetectorConfig returns default detector configuration.
func DefaultDetectorConfig() DetectorConfig {
	return DetectorConfig{
		PollInterval: 2 * time.Second,
		MaxRetries:   3,
	}
}

// AdapterStatus holds status information for an adapter.
type AdapterStatus struct {
	Name      string
	Available bool
	Type      string // "multiplexer" or "native"
}

// NewDetector creates a new terminal detector with the provided adapters.
func NewDetector(config DetectorConfig, adapters ...TerminalAdapter) *Detector {
	ctx, cancel := context.WithCancel(context.Background())

	return &Detector{
		config:   config,
		ctx:      ctx,
		cancel:   cancel,
		sessions: make(map[string]*DetectedTerminal),
		events:   make(chan TerminalEvent, 100),
		adapters: adapters,
	}
}

// Start begins the detection loop.
func (d *Detector) Start() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.running {
		return nil
	}

	d.running = true
	go d.runDetectionLoop()

	return nil
}

// Stop halts the detection loop.
func (d *Detector) Stop() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.running {
		return nil
	}

	d.running = false
	d.cancel()

	return nil
}

// IsRunning returns whether the detector is active.
func (d *Detector) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.running
}

// GetSessions returns a copy of the current session map.
func (d *Detector) GetSessions() map[string]*DetectedTerminal {
	d.mu.RLock()
	defer d.mu.RUnlock()

	copy := make(map[string]*DetectedTerminal, len(d.sessions))
	for k, v := range d.sessions {
		copy[k] = v
	}
	return copy
}

// Events returns the event channel for subscription.
func (d *Detector) Events() <-chan TerminalEvent {
	return d.events
}

// GetAdapterStatus returns status of each registered adapter.
func (d *Detector) GetAdapterStatus() []AdapterStatus {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var statuses []AdapterStatus
	for _, adapter := range d.adapters {
		statuses = append(statuses, AdapterStatus{
			Name:      adapter.Name(),
			Available: adapter.Available(d.ctx),
			Type:      getAdapterType(adapter),
		})
	}
	return statuses
}

// GetAdapters returns the registered adapters.
func (d *Detector) GetAdapters() []TerminalAdapter {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make([]TerminalAdapter, len(d.adapters))
	copy(result, d.adapters)
	return result
}

// AddAdapter adds an adapter to the detector.
func (d *Detector) AddAdapter(adapter TerminalAdapter) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.adapters = append(d.adapters, adapter)
}

// getAdapterType returns the type of adapter.
func getAdapterType(adapter TerminalAdapter) string {
	switch adapter.Name() {
	case "tmux", "screen":
		return "multiplexer"
	case "macos", "linux", "windows":
		return "native"
	default:
		return "unknown"
	}
}

// runDetectionLoop polls adapters and emits events.
func (d *Detector) runDetectionLoop() {
	ticker := time.NewTicker(d.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			d.pollAdapters()
		}
	}
}

// pollAdapters queries all adapters for sessions.
func (d *Detector) pollAdapters() {
	d.mu.RLock()
	adapters := make([]TerminalAdapter, len(d.adapters))
	copy(adapters, d.adapters)
	d.mu.RUnlock()

	var allSessions []DetectedTerminal

	for _, adapter := range adapters {
		if !adapter.Available(d.ctx) {
			continue
		}

		sessions, err := adapter.ListSessions(d.ctx)
		if err != nil {
			continue
		}

		allSessions = append(allSessions, sessions...)
	}

	d.updateSessions(allSessions)
}

// updateSessions updates session cache and emits events.
func (d *Detector) updateSessions(sessions []DetectedTerminal) {
	d.mu.Lock()
	defer d.mu.Unlock()

	currentIDs := make(map[string]bool)

	for _, s := range sessions {
		currentIDs[s.ID] = true

		if _, exists := d.sessions[s.ID]; !exists {
			d.sessions[s.ID] = &s
			d.emitEvent(TerminalEvent{
				Type:       EventSessionCreated,
				TerminalID: s.ID,
				Timestamp:  time.Now(),
			})
		}
	}

	for id := range d.sessions {
		if !currentIDs[id] {
			d.emitEvent(TerminalEvent{
				Type:       EventSessionClosed,
				TerminalID: id,
				Timestamp:  time.Now(),
			})
			delete(d.sessions, id)
		}
	}
}

// emitEvent emits a terminal event.
func (d *Detector) emitEvent(event TerminalEvent) {
	select {
	case d.events <- event:
	default:
	}
}
