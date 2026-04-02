package hotkeys

import (
	"fmt"
	"testing"
)

func TestParseHotkey(t *testing.T) {
	mods, key, err := parseHotkey("Ctrl+Shift+C")
	if err != nil {
		t.Fatalf("parseHotkey failed: %v", err)
	}
	if key != "C" {
		t.Errorf("Expected key 'C', got %s", key)
	}
	if len(mods) != 2 || mods[0] != "Ctrl" || mods[1] != "Shift" {
		t.Errorf("Expected modifiers [Ctrl Shift], got %v", mods)
	}
}

func TestParseHotkey_Error(t *testing.T) {
	tests := []struct {
		name    string
		hotkey  string
		wantErr bool
	}{
		{"empty", "", true},
		{"single key", "C", true},
		{"no modifier", "C+", true},
		{"valid", "Ctrl+C", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseHotkey(tt.hotkey)
			if tt.wantErr && err == nil {
				t.Errorf("parseHotkey(%q) expected error, got nil", tt.hotkey)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("parseHotkey(%q) unexpected error: %v", tt.hotkey, err)
			}
		})
	}
}

func TestManager_Register(t *testing.T) {
	m := NewManager()
	called := false
	callback := func() { called = true }
	err := m.Register("Ctrl+Shift+P", callback)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if m.registered["Ctrl+Shift+P"] == nil {
		t.Error("Hotkey not stored in registered map")
	}
	// Verify stored callback is the same by invoking it
	m.registered["Ctrl+Shift+P"]()
	if !called {
		t.Error("Stored callback did not update called flag")
	}
}

func TestManager_StartStop(t *testing.T) {
	m := NewManager()
	if err := m.Start(); err != nil {
		t.Errorf("Start returned error: %v", err)
	}
	m.Stop()
	// Should not panic
}

func TestManager_RegisterInvalidHotkey(t *testing.T) {
	m := NewManager()
	err := m.Register("Invalid", nil)
	if err == nil {
		t.Error("Register with invalid hotkey should return error")
	}
}

func TestManager_RegisterDuplicate(t *testing.T) {
	m := NewManager()
	callCount := 0
	callback1 := func() { callCount += 1 }
	callback2 := func() { callCount += 2 }

	err := m.Register("Ctrl+A", callback1)
	if err != nil {
		t.Fatalf("First register failed: %v", err)
	}
	err = m.Register("Ctrl+A", callback2)
	if err != nil {
		t.Fatalf("Second register failed: %v", err)
	}
	// Invoke the stored callback (should be callback2)
	m.registered["Ctrl+A"]()
	if callCount != 2 {
		t.Errorf("Expected callCount 2, got %d", callCount)
	}
}

func TestManager_ConcurrentRegister(t *testing.T) {
	m := NewManager()
	const goroutines = 10
	errCh := make(chan error, goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			hotkey := fmt.Sprintf("Ctrl+%d", idx)
			err := m.Register(hotkey, func() {})
			if err != nil {
				errCh <- fmt.Errorf("Register failed for %s: %v", hotkey, err)
			} else {
				errCh <- nil
			}
		}(i)
	}
	for i := 0; i < goroutines; i++ {
		if err := <-errCh; err != nil {
			t.Error(err)
		}
	}
}
