package hotkeys

import (
	"fmt"
	"strings"

	"github.com/go-vgo/robotgo"
)

// Manager handles global hotkey registration.
type Manager struct {
	registered map[string]func()
}

// NewManager creates a new hotkey manager.
func NewManager() *Manager {
	return &Manager{
		registered: make(map[string]func()),
	}
}

// Register binds a hotkey string (e.g., "Ctrl+Shift+C") to a callback.
func (m *Manager) Register(hotkey string, callback func()) error {
	_, _, err := parseHotkey(hotkey)
	if err != nil {
		return err
	}

	// robotgo.AddHotkey expects modifiers as separate arguments
	// For simplicity, we'll use robotgo.EventHook for now (alternative approach)
	// This is a placeholder implementation; real registration requires platform‑specific hooks
	// We'll log and store the mapping for now.
	m.registered[hotkey] = callback
	return nil
}

// Start listening for hotkeys (platform‑specific).
func (m *Manager) Start() error {
	// TODO: implement platform‑specific global hotkey listening
	// This may require separate implementations for Windows/macOS/Linux
	_ = robotgo.Version // keep import for now
	return nil
}

// Stop listening for hotkeys.
func (m *Manager) Stop() {
	// TODO: clean up hooks
}

func parseHotkey(hotkey string) (modifiers []string, key string, err error) {
	parts := strings.Split(hotkey, "+")
	if len(parts) < 2 {
		return nil, "", fmt.Errorf("hotkey must include modifier and key")
	}
	key = parts[len(parts)-1]
	modifiers = parts[:len(parts)-1]
	// Normalize modifier names
	for i, mod := range modifiers {
		modifiers[i] = strings.Title(strings.ToLower(mod))
	}
	return modifiers, key, nil
}