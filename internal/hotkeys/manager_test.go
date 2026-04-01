package hotkeys

import (
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