package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/spf13/viper"
)

func resetViper(t *testing.T) {
	t.Helper()
	viper.Reset()
	globalConfig = nil
	configPath = ""
}

func TestConfigInitAndSave(t *testing.T) {
	resetViper(t)
	t.Setenv("PAIRADMIN_LLM_OPENAI_API_KEY", "sk-test")
	t.Setenv("PAIRADMIN_LLM_ANTHROPIC_API_KEY", "sk-test")

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Write minimal config with API keys
	configContent := `
llm:
  provider: openai
  openai:
    api_key: sk-test
  anthropic:
    api_key: sk-test
`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

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

func TestConfigValidation(t *testing.T) {
	resetViper(t)
	tests := []struct {
		name         string
		provider     string
		openaiKey    string
		anthropicKey string
		ollamaURL    string
		wantErr      bool
	}{
		{
			name:      "valid openai",
			provider:  "openai",
			openaiKey: "sk-test",
			wantErr:   false,
		},
		{
			name:      "openai missing api key",
			provider:  "openai",
			openaiKey: "",
			wantErr:   true,
		},
		{
			name:         "valid anthropic",
			provider:     "anthropic",
			anthropicKey: "sk-test",
			wantErr:      false,
		},
		{
			name:         "anthropic missing api key",
			provider:     "anthropic",
			anthropicKey: "",
			wantErr:      true,
		},
		{
			name:      "valid ollama",
			provider:  "ollama",
			ollamaURL: "http://localhost:11434",
			wantErr:   false,
		},
		{
			name:      "ollama missing base_url",
			provider:  "ollama",
			ollamaURL: "",
			wantErr:   true,
		},
		{
			name:     "unknown provider",
			provider: "unknown",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			cfg.LLM.Provider = tt.provider
			cfg.LLM.OpenAI.APIKey = tt.openaiKey
			cfg.LLM.Anthropic.APIKey = tt.anthropicKey
			cfg.LLM.Ollama.BaseURL = tt.ollamaURL

			err := cfg.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigInvalidYAML(t *testing.T) {
	resetViper(t)
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Write invalid YAML
	if err := os.WriteFile(configFile, []byte("invalid: yaml: :"), 0644); err != nil {
		t.Fatal(err)
	}

	err := Init(configFile)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestConfigPermissionError(t *testing.T) {
	resetViper(t)
	if os.Geteuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	// Create subdirectory without write permissions
	if err := os.MkdirAll(subDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(subDir, 0755)
	configFile := filepath.Join(subDir, "config.yaml")

	err := Init(configFile)
	if err == nil {
		t.Error("Expected permission error, got nil")
	}
}

func TestConfigEmptyConfig(t *testing.T) {
	resetViper(t)
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Write empty file
	if err := os.WriteFile(configFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	// Should fail because openai provider requires api_key (config: llm.openai.api_key, env: PAIRADMIN_LLM_OPENAI_API_KEY)
	err := Init(configFile)
	if err == nil {
		cfg := Get()
		t.Errorf("Expected validation error for empty config, got nil. Config: provider=%s, openai.api_key='%s'", cfg.LLM.Provider, cfg.LLM.OpenAI.APIKey)
	} else {
		t.Logf("Got error: %v", err)
		if err.Error() != "openai provider requires api_key (config: llm.openai.api_key, env: PAIRADMIN_LLM_OPENAI_API_KEY)" {
			t.Errorf("Unexpected error message: %v", err)
		}
	}
}
func TestConcurrentAccess(t *testing.T) {
	resetViper(t)
	t.Setenv("PAIRADMIN_LLM_OPENAI_API_KEY", "sk-test-concurrent")
	t.Setenv("PAIRADMIN_LLM_ANTHROPIC_API_KEY", "sk-test-concurrent")

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
llm:
  provider: openai
  openai:
    api_key: sk-test-concurrent
  anthropic:
    api_key: sk-test-concurrent
`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	err := Init(configFile)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	const readers = 10
	const writers = 2
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(readers + writers)

	// Reader goroutines
	for i := 0; i < readers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				cfg := Get()
				if cfg == nil {
					t.Errorf("Reader %d iteration %d: Get returned nil", id, j)
					return
				}
				_ = cfg.LLM.Provider
			}
		}(i)
	}

	// Writer goroutines
	for i := 0; i < writers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				if err := Save(); err != nil {
					t.Errorf("Writer %d iteration %d: Save failed: %v", id, j, err)
				}
			}
		}(i)
	}

	wg.Wait()
}
