package config

import (
	"path/filepath"
	"testing"
)

func TestConfigInitAndSave(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	err := Init(configFile)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	cfg := Get()
	if cfg == nil {
		t.Fatal("Get returned nil")
	}
	if cfg.LLM.Provider != "openai" {
		t.Errorf("Expected provider openai, got %s", cfg.LLM.Provider)
	}
	if cfg.UI.Hotkeys.CopyLastCommand != "Ctrl+Shift+C" {
		t.Errorf("Expected hotkey Ctrl+Shift+C, got %s", cfg.UI.Hotkeys.CopyLastCommand)
	}

	// Modify and save
	cfg.LLM.Provider = "anthropic"
	err = Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Re‑init and verify
	err = Init(configFile)
	if err != nil {
		t.Fatalf("Re‑init failed: %v", err)
	}
	cfg2 := Get()
	if cfg2.LLM.Provider != "anthropic" {
		t.Errorf("After save, expected provider anthropic, got %s", cfg2.LLM.Provider)
	}
}