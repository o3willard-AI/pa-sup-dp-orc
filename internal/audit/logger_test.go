package audit

import (
	"os"
	"sync"
	"testing"
	"time"
)

func TestLogger_Log(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	event := Event{
		Timestamp: time.Now(),
		Action:    "config_save",
		User:      "test-user",
		Success:   true,
		Details:   "LLM provider changed to OpenAI",
	}

	err = logger.Log(event)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}
}

func TestLogger_LogAfterClose(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	// Close immediately
	if err := logger.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	// Attempt to log after close
	event := Event{
		Timestamp: time.Now(),
		Action:    "config_save",
		Success:   true,
	}
	err = logger.Log(event)
	if err == nil {
		t.Error("Expected error when logging after close, got nil")
	}
}

func TestLogger_ConcurrentLog(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	const workers = 10
	const eventsPerWorker = 5
	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < eventsPerWorker; i++ {
				event := Event{
					Timestamp: time.Now(),
					Action:    "concurrent_action",
					User:      "user",
					Success:   true,
					Details:   "concurrent test",
				}
				if err := logger.Log(event); err != nil {
					t.Errorf("Worker %d event %d failed: %v", id, i, err)
				}
			}
		}(w)
	}
	wg.Wait()
}
