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

func TestManager_Register(t *testing.T) {
	m := NewManager()
	called := false
	err := m.Register("Ctrl+Shift+P", func() { called = true })
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if m.registered["Ctrl+Shift+P"] == nil {
		t.Error("Hotkey not stored in registered map")
	}
	_ = called // avoid unused variable warning
}