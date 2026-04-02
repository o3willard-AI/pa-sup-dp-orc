// Copyright 2026 The PairAdmin Authors
// SPDX-License-Identifier: MIT

package terminal

import (
	"context"
	"testing"
	"time"
)

func TestNewDetector(t *testing.T) {
	config := DefaultDetectorConfig()
	det := NewDetector(config)

	if det == nil {
		t.Fatal("NewDetector returned nil")
	}
	if det.sessions == nil {
		t.Error("sessions map not initialized")
	}
	if det.events == nil {
		t.Error("events channel not initialized")
	}
	if det.config.PollInterval != 2*time.Second {
		t.Errorf("Expected PollInterval 2s, got %v", det.config.PollInterval)
	}
}

func TestNewDetectorWithAdapters(t *testing.T) {
	config := DefaultDetectorConfig()
	mockAdapter := &mockAdapter{name: "test"}
	det := NewDetector(config, mockAdapter)

	if len(det.adapters) != 1 {
		t.Errorf("Expected 1 adapter, got %d", len(det.adapters))
	}
}

func TestDetectorStartStop(t *testing.T) {
	det := NewDetector(DefaultDetectorConfig())

	if det.IsRunning() {
		t.Error("Detector should not be running initially")
	}

	if err := det.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if !det.IsRunning() {
		t.Error("Detector should be running after Start")
	}

	if err := det.Start(); err != nil {
		t.Errorf("Second Start should be no-op: %v", err)
	}

	if err := det.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	if det.IsRunning() {
		t.Error("Detector should not be running after Stop")
	}
}

func TestDetectorGetSessions(t *testing.T) {
	det := NewDetector(DefaultDetectorConfig())
	sessions := det.GetSessions()

	if len(sessions) != 0 {
		t.Errorf("Expected 0 sessions, got %d", len(sessions))
	}
}

func TestDetectorGetAdapterStatus(t *testing.T) {
	mockAdapter := &mockAdapter{name: "test", available: true}
	det := NewDetector(DefaultDetectorConfig(), mockAdapter)

	statuses := det.GetAdapterStatus()
	if len(statuses) != 1 {
		t.Errorf("Expected 1 status, got %d", len(statuses))
	}
	if statuses[0].Name != "test" {
		t.Errorf("Expected name 'test', got %s", statuses[0].Name)
	}
	if statuses[0].Type != "unknown" {
		t.Errorf("Expected type 'unknown', got %s", statuses[0].Type)
	}
}

func TestDetectorAddAdapter(t *testing.T) {
	det := NewDetector(DefaultDetectorConfig())
	mockAdapter := &mockAdapter{name: "test"}

	det.AddAdapter(mockAdapter)

	if len(det.adapters) != 1 {
		t.Errorf("Expected 1 adapter, got %d", len(det.adapters))
	}
}

func TestDetectorEvents(t *testing.T) {
	det := NewDetector(DefaultDetectorConfig())
	events := det.Events()

	if events == nil {
		t.Error("Events channel should not be nil")
	}
}

func TestDetectorGetAdapters(t *testing.T) {
	det := NewDetector(DefaultDetectorConfig())
	mockAdapter := &mockAdapter{name: "test"}
	det.AddAdapter(mockAdapter)

	adapters := det.GetAdapters()
	if len(adapters) != 1 {
		t.Errorf("Expected 1 adapter, got %d", len(adapters))
	}
}

type mockAdapter struct {
	name      string
	available bool
}

func (m *mockAdapter) Name() string                                             { return m.name }
func (m *mockAdapter) Available(context.Context) bool                           { return m.available }
func (m *mockAdapter) ListSessions(context.Context) ([]DetectedTerminal, error) { return nil, nil }
func (m *mockAdapter) Capture(context.Context, string) (string, error)          { return "", nil }
func (m *mockAdapter) Subscribe(context.Context, string) (<-chan TerminalEvent, error) {
	return nil, nil
}
func (m *mockAdapter) GetDimensions(context.Context, string) (int, int, error) { return 24, 80, nil }
