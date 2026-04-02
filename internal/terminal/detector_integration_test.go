package terminal_test

import (
	"testing"
	"time"

	"github.com/pairadmin/pairadmin/internal/terminal"
	"github.com/pairadmin/pairadmin/internal/terminal/testhelpers"
)

// TestDetector_WithMockAdapters tests detector with multiple mock adapters
func TestDetector_WithMockAdapters(t *testing.T) {
	mock1 := testhelpers.NewMockAdapter("mock1")
	mock2 := testhelpers.NewMockAdapter("mock2")

	term1 := testhelpers.NewMockTerminal("term-1", "Terminal 1", "content1")
	term2 := testhelpers.NewMockTerminal("term-2", "Terminal 2", "content2")
	mock1.AddTerminal(term1)
	mock2.AddTerminal(term2)

	config := terminal.DetectorConfig{
		PollInterval: 100 * time.Millisecond,
	}
	detector := terminal.NewDetector(config, mock1, mock2)

	if err := detector.Start(); err != nil {
		t.Fatalf("failed to start detector: %v", err)
	}
	defer detector.Stop()

	time.Sleep(150 * time.Millisecond)

	sessions := detector.GetSessions()
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}
}

// TestDetector_AdapterFallback tests fallback when primary adapter fails
func TestDetector_AdapterFallback(t *testing.T) {
	primary := testhelpers.NewMockAdapter("primary")
	primary.SetAvailable(false)

	backup := testhelpers.NewMockAdapter("backup")
	term := testhelpers.NewMockTerminal("term-1", "Terminal 1", "content")
	backup.AddTerminal(term)

	config := terminal.DetectorConfig{
		PollInterval: 100 * time.Millisecond,
	}
	detector := terminal.NewDetector(config, primary, backup)

	if err := detector.Start(); err != nil {
		t.Fatalf("failed to start detector: %v", err)
	}
	defer detector.Stop()

	time.Sleep(150 * time.Millisecond)

	sessions := detector.GetSessions()
	if len(sessions) != 1 {
		t.Errorf("expected 1 session from backup, got %d", len(sessions))
	}
}

// TestDetector_AddAdapterDynamic tests adding adapters after creation
func TestDetector_AddAdapterDynamic(t *testing.T) {
	config := terminal.DetectorConfig{
		PollInterval: 100 * time.Millisecond,
	}
	detector := terminal.NewDetector(config)

	if err := detector.Start(); err != nil {
		t.Fatalf("failed to start detector: %v", err)
	}
	defer detector.Stop()

	time.Sleep(150 * time.Millisecond)

	sessions := detector.GetSessions()
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions initially, got %d", len(sessions))
	}

	adapter := testhelpers.NewMockAdapter("dynamic")
	term := testhelpers.NewMockTerminal("term-1", "Terminal 1", "content")
	adapter.AddTerminal(term)
	detector.AddAdapter(adapter)

	time.Sleep(150 * time.Millisecond)

	sessions = detector.GetSessions()
	if len(sessions) != 1 {
		t.Errorf("expected 1 session after adding adapter, got %d", len(sessions))
	}
}

// TestDetector_ConcurrentAccess tests thread-safe concurrent access
func TestDetector_ConcurrentAccess(t *testing.T) {
	adapter := testhelpers.NewMockAdapter("test")
	for i := 0; i < 10; i++ {
		term := testhelpers.NewMockTerminal(
			"term-"+string(rune('0'+i)),
			"Terminal "+string(rune('0'+i)),
			"content",
		)
		adapter.AddTerminal(term)
	}

	config := terminal.DetectorConfig{
		PollInterval: 100 * time.Millisecond,
	}
	detector := terminal.NewDetector(config, adapter)

	if err := detector.Start(); err != nil {
		t.Fatalf("failed to start detector: %v", err)
	}
	defer detector.Stop()

	time.Sleep(150 * time.Millisecond)

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_ = detector.GetSessions()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
