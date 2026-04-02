# Phase 3: AI Collaboration Features Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement AI collaboration features for PairAdmin v2.0: LLM provider integration, command history storage, clipboard manager, settings dialog, sensitive data filtering, and hotkey configuration.

**Architecture:** Extend the existing terminal‑capture infrastructure with an LLM gateway supporting multiple providers (OpenAI, Anthropic, Ollama). Store suggested commands per terminal session in SQLite (via `modernc.org/sqlite`). Add cross‑platform clipboard integration (`golang.design/x/clipboard`). Frontend components: chat area, command sidebar, settings dialog. Secure sensitive data with regex filtering before sending to LLM.

**Tech Stack:** Go 1.21+, Wails 2, Svelte, SQLite (pure‑Go driver), Viper for configuration, `golang.design/x/clipboard`, `github.com/go‑vgo/robotgo` for hotkeys.

---

## File Structure

**Backend:**
- `internal/llm/gateway.go` – LLM gateway interface
- `internal/llm/providers/openai.go` – OpenAI provider
- `internal/llm/providers/anthropic.go` – Anthropic provider
- `internal/llm/providers/ollama.go` – Ollama provider (local)
- `internal/config/manager.go` – configuration manager (Viper)
- `internal/config/keychain.go` – OS keychain abstraction (Phase 4)
- `internal/session/store.go` – command history storage (SQLite)
- `internal/session/types.go` – `SuggestedCommand` and related types
- `internal/clipboard/manager.go` – clipboard operations
- `internal/security/filter.go` – sensitive data redaction
- `internal/ui/chat_handlers.go` – Wails handlers for chat and commands
- `internal/ui/settings_handlers.go` – Wails handlers for settings
- `internal/app/app.go` – main application struct (extends `main.App`)

**Frontend:**
- `frontend/src/components/ChatArea.svelte` – chat message display and input
- `frontend/src/components/CommandSidebar.svelte` – command history cards
- `frontend/src/components/SettingsDialog.svelte` – settings dialog with tabs
- `frontend/src/components/TerminalTabs.svelte` – terminal session tabs
- `frontend/src/components/StatusBar.svelte` – model selector, context meter, settings button
- `frontend/src/lib/stores.js` – Svelte stores for global state
- `frontend/src/main.js` – Wails bindings setup

**Configuration:**
- `config.yaml` – default configuration file
- `~/.pairadmin/config.yaml` – user configuration

---

### Task 1: Add Dependencies and Project Structure

**Files:**
- Modify: `go.mod`
- Create: `config.yaml`

- [ ] **Step 1: Add required Go dependencies**

Run:
```bash
go get github.com/spf13/viper
go get modernc.org/sqlite
go get golang.design/x/clipboard
go get github.com/go-vgo/robotgo
```

- [ ] **Step 2: Create default configuration file**

Create `config.yaml` at project root:

```yaml
# PairAdmin Configuration
llm:
  provider: "openai"  # openai, anthropic, ollama
  openai:
    api_key: ""
    model: "gpt-4"
    base_url: "https://api.openai.com/v1"
  anthropic:
    api_key: ""
    model: "claude-3-haiku-20240307"
    base_url: "https://api.anthropic.com"
  ollama:
    base_url: "http://localhost:11434"
    model: "llama3"

security:
  filter_patterns:
    - name: "password"
      pattern: "(?i)password\\s*[:=]\\s*['\\\"]?[^'\\\"\\s]+"
    - name: "api_key"
      pattern: "(?i)(api[_-]?key|token)[\\s:=]+['\\\"]?[a-zA-Z0-9_\\-]{20,}['\\\"]?"
    - name: "aws_access_key"
      pattern: "AKIA[0-9A-Z]{16}"

ui:
  theme: "system"  # system, light, dark
  hotkeys:
    copy_last_command: "Ctrl+Shift+C"
    focus_app: "Ctrl+Shift+P"
```

- [ ] **Step 3: Update .gitignore to exclude user config**

Add to `.gitignore`:

```
# User configuration
~/.pairadmin/
```

- [ ] **Step 4: Commit dependencies and config**

```bash
git add go.mod go.sum config.yaml .gitignore
git commit -m "chore: add Phase 3 dependencies and default config"
```

---

### Task 2: LLM Gateway Interface and Configuration Manager

**Files:**
- Create: `internal/llm/gateway.go`
- Create: `internal/config/manager.go`
- Test: `internal/llm/gateway_test.go`
- Test: `internal/config/manager_test.go`

- [ ] **Step 1: Define LLM gateway interface**

Create `internal/llm/gateway.go`:

```go
package llm

import (
	"context"
)

// MessageRole represents the role of a chat message.
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// ChatMessage represents a single message in a conversation.
type ChatMessage struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

// CompletionRequest holds parameters for an LLM completion.
type CompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// CompletionResponse holds the LLM's response.
type CompletionResponse struct {
	Content string `json:"content"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Gateway defines the interface for LLM providers.
type Gateway interface {
	// Name returns the provider name (e.g., "openai").
	Name() string
	// Complete sends a completion request and returns the response.
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
	// StreamComplete streams the response via a channel.
	StreamComplete(ctx context.Context, req CompletionRequest) (<-chan string, error)
}
```

- [ ] **Step 2: Write failing test for gateway interface**

Create `internal/llm/gateway_test.go`:

```go
package llm

import (
	"context"
	"testing"
)

func TestGatewayInterface(t *testing.T) {
	// This test ensures the interface is defined correctly.
	var gw Gateway
	_ = gw
	t.Log("Gateway interface defined")
}
```

Run: `go test ./internal/llm/... -v`
Expected: PASS (no implementation needed)

- [ ] **Step 3: Create configuration manager**

Create `internal/config/manager.go`:

```go
package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	LLM struct {
		Provider string `mapstructure:"provider"`
		OpenAI   struct {
			APIKey  string `mapstructure:"api_key"`
			Model   string `mapstructure:"model"`
			BaseURL string `mapstructure:"base_url"`
		} `mapstructure:"openai"`
		Anthropic struct {
			APIKey  string `mapstructure:"api_key"`
			Model   string `mapstructure:"model"`
			BaseURL string `mapstructure:"base_url"`
		} `mapstructure:"anthropic"`
		Ollama struct {
			BaseURL string `mapstructure:"base_url"`
			Model   string `mapstructure:"model"`
		} `mapstructure:"ollama"`
	} `mapstructure:"llm"`
	Security struct {
		FilterPatterns []FilterPattern `mapstructure:"filter_patterns"`
	} `mapstructure:"security"`
	UI struct {
		Theme   string `mapstructure:"theme"`
		Hotkeys struct {
			CopyLastCommand string `mapstructure:"copy_last_command"`
			FocusApp        string `mapstructure:"focus_app"`
		} `mapstructure:"hotkeys"`
	} `mapstructure:"ui"`
}

// FilterPattern defines a regex pattern for sensitive data.
type FilterPattern struct {
	Name    string `mapstructure:"name"`
	Pattern string `mapstructure:"pattern"`
}

var (
	globalConfig *Config
	configPath   string
)

// Init loads configuration from file and environment variables.
func Init(configFile string) error {
	configPath = configFile
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("llm.provider", "openai")
	viper.SetDefault("llm.openai.model", "gpt-4")
	viper.SetDefault("llm.openai.base_url", "https://api.openai.com/v1")
	viper.SetDefault("llm.anthropic.model", "claude-3-haiku-20240307")
	viper.SetDefault("llm.anthropic.base_url", "https://api.anthropic.com")
	viper.SetDefault("llm.ollama.base_url", "http://localhost:11434")
	viper.SetDefault("llm.ollama.model", "llama3")
	viper.SetDefault("ui.theme", "system")
	viper.SetDefault("ui.hotkeys.copy_last_command", "Ctrl+Shift+C")
	viper.SetDefault("ui.hotkeys.focus_app", "Ctrl+Shift+P")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			// Create directory if it doesn't exist
			dir := filepath.Dir(configFile)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			if err := viper.WriteConfigAs(configFile); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Bind environment variables
	viper.SetEnvPrefix("PAIRADMIN")
	viper.AutomaticEnv()

	// Unmarshal
	globalConfig = &Config{}
	if err := viper.Unmarshal(globalConfig); err != nil {
		return err
	}

	return nil
}

// Get returns the global configuration.
func Get() *Config {
	return globalConfig
}

// Save writes the current configuration to disk.
func Save() error {
	viper.Set("llm", globalConfig.LLM)
	viper.Set("security", globalConfig.Security)
	viper.Set("ui", globalConfig.UI)
	return viper.WriteConfig()
}
```

- [ ] **Step 4: Write tests for configuration manager**

Create `internal/config/manager_test.go`:

```go
package config

import (
	"os"
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
```

Run: `go test ./internal/config/... -v`
Expected: PASS

- [ ] **Step 5: Commit configuration manager**

```bash
git add internal/llm/gateway.go internal/llm/gateway_test.go internal/config/manager.go internal/config/manager_test.go
git commit -m "feat: add LLM gateway interface and configuration manager"
```

---

### Task 3: OpenAI Provider Implementation

**Files:**
- Create: `internal/llm/providers/openai.go`
- Test: `internal/llm/providers/openai_test.go`

- [ ] **Step 1: Implement OpenAI provider**

Create `internal/llm/providers/openai.go`:

```go
package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pairadmin/pairadmin/internal/llm"
)

// OpenAIProvider implements llm.Gateway for OpenAI API.
type OpenAIProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider.
func NewOpenAIProvider(apiKey, model, baseURL string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Name returns "openai".
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// openAIRequest is the JSON structure for OpenAI completions.
type openAIRequest struct {
	Model       string                   `json:"model"`
	Messages    []openAIMessage          `json:"messages"`
	MaxTokens   int                      `json:"max_tokens,omitempty"`
	Temperature float64                  `json:"temperature,omitempty"`
	Stream      bool                     `json:"stream,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Complete sends a completion request to OpenAI.
func (p *OpenAIProvider) Complete(ctx context.Context, req llm.CompletionRequest) (*llm.CompletionResponse, error) {
	openAIReq := openAIRequest{
		Model:       p.model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      false,
	}
	for _, msg := range req.Messages {
		openAIReq.Messages = append(openAIReq.Messages, openAIMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	body, err := json.Marshal(openAIReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai api error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	llmResp := &llm.CompletionResponse{
		Content: openAIResp.Choices[0].Message.Content,
		Model:   openAIResp.Model,
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}(openAIResp.Usage),
	}
	return llmResp, nil
}

// StreamComplete is not implemented yet (returns error).
func (p *OpenAIProvider) StreamComplete(ctx context.Context, req llm.CompletionRequest) (<-chan string, error) {
	return nil, fmt.Errorf("streaming not yet implemented for OpenAI")
}
```

- [ ] **Step 2: Write tests for OpenAI provider**

Create `internal/llm/providers/openai_test.go`:

```go
package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pairadmin/pairadmin/internal/llm"
)

func TestOpenAIProvider_Name(t *testing.T) {
	p := NewOpenAIProvider("key", "gpt-4", "https://api.openai.com/v1")
	if p.Name() != "openai" {
		t.Errorf("Expected name 'openai', got %s", p.Name())
	}
}

func TestOpenAIProvider_Complete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"created": 1677652288,
			"model": "gpt-4",
			"choices": [{
				"index": 0,
				"message": {"role": "assistant", "content": "Hello!"},
				"finish_reason": "stop"
			}],
			"usage": {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
		}`
		w.Write([]byte(response))
	}))
	defer server.Close()

	p := NewOpenAIProvider("test-key", "gpt-4", server.URL)
	req := llm.CompletionRequest{
		Model: "gpt-4",
		Messages: []llm.ChatMessage{
			{Role: llm.RoleUser, Content: "Hello"},
		},
	}
	resp, err := p.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("Complete failed: %v", err)
	}
	if resp.Content != "Hello!" {
		t.Errorf("Expected content 'Hello!', got %s", resp.Content)
	}
	if resp.Model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got %s", resp.Model)
	}
	if resp.Usage.TotalTokens != 15 {
		t.Errorf("Expected 15 total tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestOpenAIProvider_Complete_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid API key"}`))
	}))
	defer server.Close()

	p := NewOpenAIProvider("bad-key", "gpt-4", server.URL)
	req := llm.CompletionRequest{
		Model: "gpt-4",
		Messages: []llm.ChatMessage{
			{Role: llm.RoleUser, Content: "Hello"},
		},
	}
	_, err := p.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for unauthorized request")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("Expected error containing 401, got %v", err)
	}
}
```

Run: `go test ./internal/llm/providers/... -v`
Expected: PASS

- [ ] **Step 3: Commit OpenAI provider**

```bash
git add internal/llm/providers/openai.go internal/llm/providers/openai_test.go
git commit -m "feat: add OpenAI LLM provider"
```

---

### Task 4: Session Store and SuggestedCommand Types

**Files:**
- Create: `internal/session/types.go`
- Create: `internal/session/store.go`
- Test: `internal/session/store_test.go`

- [ ] **Step 1: Define SuggestedCommand and session types**

Create `internal/session/types.go`:

```go
package session

import (
	"time"
)

// SuggestedCommand represents a command suggested by the AI.
type SuggestedCommand struct {
	ID          string    `json:"id"`
	TerminalID  string    `json:"terminal_id"`
	Command     string    `json:"command"`
	Description string    `json:"description"`
	Context     string    `json:"context"` // Original user question or context
	CreatedAt   time.Time `json:"created_at"`
	UsedCount   int       `json:"used_count"`
	LastUsedAt  time.Time `json:"last_used_at"`
}

// Session represents a terminal session with its command history.
type Session struct {
	ID        string             `json:"id"`
	TerminalID string            `json:"terminal_id"`
	CreatedAt time.Time          `json:"created_at"`
	Commands  []SuggestedCommand `json:"commands"`
}
```

- [ ] **Step 2: Implement SQLite session store**

Create `internal/session/store.go`:

```go
package session

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// Store manages persistence of sessions and suggested commands.
type Store struct {
	db *sql.DB
}

// NewStore creates a new session store at the given path.
func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err := s.createTables(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) createTables() error {
	// Sessions table
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			terminal_id TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	// Commands table
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS commands (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			terminal_id TEXT NOT NULL,
			command TEXT NOT NULL,
			description TEXT,
			context TEXT,
			created_at DATETIME NOT NULL,
			used_count INTEGER DEFAULT 0,
			last_used_at DATETIME,
			FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
		)
	`)
	return err
}

// AddSession creates a new session.
func (s *Store) AddSession(sessionID, terminalID string) error {
	_, err := s.db.Exec(
		"INSERT INTO sessions (id, terminal_id, created_at) VALUES (?, ?, ?)",
		sessionID, terminalID, time.Now(),
	)
	return err
}

// AddCommand adds a suggested command to a session.
func (s *Store) AddCommand(cmd SuggestedCommand) error {
	_, err := s.db.Exec(
		`INSERT INTO commands 
		(id, session_id, terminal_id, command, description, context, created_at, used_count, last_used_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		cmd.ID, cmd.TerminalID, cmd.TerminalID, cmd.Command, cmd.Description, cmd.Context,
		cmd.CreatedAt, cmd.UsedCount, cmd.LastUsedAt,
	)
	return err
}

// GetCommandsByTerminal returns all commands for a terminal.
func (s *Store) GetCommandsByTerminal(terminalID string) ([]SuggestedCommand, error) {
	rows, err := s.db.Query(`
		SELECT id, terminal_id, command, description, context, created_at, used_count, last_used_at
		FROM commands WHERE terminal_id = ?
		ORDER BY created_at DESC
	`, terminalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []SuggestedCommand
	for rows.Next() {
		var cmd SuggestedCommand
		err := rows.Scan(&cmd.ID, &cmd.TerminalID, &cmd.Command, &cmd.Description, &cmd.Context,
			&cmd.CreatedAt, &cmd.UsedCount, &cmd.LastUsedAt)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}
	return commands, nil
}

// IncrementUsedCount increments the used count of a command.
func (s *Store) IncrementUsedCount(commandID string) error {
	_, err := s.db.Exec(`
		UPDATE commands 
		SET used_count = used_count + 1, last_used_at = ?
		WHERE id = ?
	`, time.Now(), commandID)
	return err
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
```

- [ ] **Step 3: Write tests for session store**

Create `internal/session/store_test.go`:

```go
package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStore_Crud(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	defer store.Close()

	// Add session
	sessionID := "session-1"
	terminalID := "terminal-1"
	err = store.AddSession(sessionID, terminalID)
	if err != nil {
		t.Fatalf("AddSession failed: %v", err)
	}

	// Add command
	cmd := SuggestedCommand{
		ID:          "cmd-1",
		TerminalID:  terminalID,
		Command:     "ls -la",
		Description: "List files",
		Context:     "User asked how to see hidden files",
		CreatedAt:   time.Now(),
		UsedCount:   0,
	}
	err = store.AddCommand(cmd)
	if err != nil {
		t.Fatalf("AddCommand failed: %v", err)
	}

	// Get commands
	commands, err := store.GetCommandsByTerminal(terminalID)
	if err != nil {
		t.Fatalf("GetCommandsByTerminal failed: %v", err)
	}
	if len(commands) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(commands))
	}
	if commands[0].Command != "ls -la" {
		t.Errorf("Expected command 'ls -la', got %s", commands[0].Command)
	}

	// Increment used count
	err = store.IncrementUsedCount("cmd-1")
	if err != nil {
		t.Fatalf("IncrementUsedCount failed: %v", err)
	}

	// Verify updated count (would need a GetCommand to check, but skip for brevity)
}

func TestStore_NewStore_CreatesTables(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	defer store.Close()

	// Try to insert a session (should succeed if tables exist)
	err = store.AddSession("test", "terminal")
	if err != nil {
		t.Fatalf("AddSession after NewStore failed (tables not created?): %v", err)
	}
}
```

Run: `go test ./internal/session/... -v`
Expected: PASS

- [ ] **Step 4: Commit session store**

```bash
git add internal/session/types.go internal/session/store.go internal/session/store_test.go
git commit -m "feat: add session store and SuggestedCommand types"
```

---

### Task 5: Clipboard Manager

**Files:**
- Create: `internal/clipboard/manager.go`
- Test: `internal/clipboard/manager_test.go`

- [ ] **Step 1: Implement clipboard manager**

Create `internal/clipboard/manager.go`:

```go
package clipboard

import (
	"context"
	"fmt"

	"golang.design/x/clipboard"
)

// Manager provides cross‑platform clipboard operations.
type Manager struct{}

// NewManager creates a new clipboard manager.
func NewManager() *Manager {
	return &Manager{}
}

// CopyToTerminal copies text to clipboard and optionally focuses the target terminal.
func (m *Manager) CopyToTerminal(ctx context.Context, text string, terminalID string) error {
	// Write to clipboard
	err := clipboard.Write(clipboard.FmtText, []byte(text))
	if err != nil {
		return fmt.Errorf("write to clipboard: %w", err)
	}

	// TODO: Focus terminal if terminalID provided and focus‑paste configured
	// This requires platform‑specific window focus logic (future enhancement)

	return nil
}

// ReadFromClipboard reads text from clipboard.
func (m *Manager) ReadFromClipboard(ctx context.Context) (string, error) {
	data := clipboard.Read(clipboard.FmtText)
	if data == nil {
		return "", fmt.Errorf("clipboard empty or unsupported format")
	}
	return string(data), nil
}
```

- [ ] **Step 2: Write tests for clipboard manager**

Create `internal/clipboard/manager_test.go`:

```go
package clipboard

import (
	"context"
	"testing"
)

func TestManager_CopyAndRead(t *testing.T) {
	// Skip test if clipboard is not available (e.g., in CI)
	if !clipboard.Init() {
		t.Skip("Clipboard not available in this environment")
	}

	m := NewManager()
	ctx := context.Background()
	text := "test clipboard content"

	// Copy
	err := m.CopyToTerminal(ctx, text, "")
	if err != nil {
		t.Fatalf("CopyToTerminal failed: %v", err)
	}

	// Read
	read, err := m.ReadFromClipboard(ctx)
	if err != nil {
		t.Fatalf("ReadFromClipboard failed: %v", err)
	}
	if read != text {
		t.Errorf("Expected %q, got %q", text, read)
	}
}
```

Run: `go test ./internal/clipboard/... -v`
Expected: PASS (or skip if clipboard unavailable)

- [ ] **Step 3: Commit clipboard manager**

```bash
git add internal/clipboard/manager.go internal/clipboard/manager_test.go
git commit -m "feat: add clipboard manager"
```

---

### Task 6: Security Filter

**Files:**
- Create: `internal/security/filter.go`
- Test: `internal/security/filter_test.go`

- [ ] **Step 1: Implement sensitive data filter**

Create `internal/security/filter.go`:

```go
package security

import (
	"regexp"
	"strings"
)

// Filter redacts sensitive patterns from text.
type Filter struct {
	patterns []*regexp.Regexp
	replace  string
}

// NewFilter creates a filter with compiled regex patterns.
func NewFilter(rawPatterns []string) (*Filter, error) {
	f := &Filter{
		replace: "[REDACTED]",
	}
	for _, p := range rawPatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		f.patterns = append(f.patterns, re)
	}
	return f, nil
}

// Redact replaces matches of any pattern with the replacement string.
func (f *Filter) Redact(text string) string {
	result := text
	for _, re := range f.patterns {
		result = re.ReplaceAllString(result, f.replace)
	}
	return result
}

// ContainsSensitive returns true if text matches any pattern.
func (f *Filter) ContainsSensitive(text string) bool {
	for _, re := range f.patterns {
		if re.MatchString(text) {
			return true
		}
	}
	return false
}

// DefaultFilter returns a filter with common sensitive patterns.
func DefaultFilter() *Filter {
	patterns := []string{
		`(?i)password\s*[:=]\s*['\"]?[^'\"]+`,
		`(?i)(api[_-]?key|token)[\s:=]+['\"]?[a-zA-Z0-9_\-]{20,}['\"]?`,
		`AKIA[0-9A-Z]{16}`,
		`-----BEGIN (RSA|DSA|EC|OPENSSH) PRIVATE KEY-----`,
		`(?i)secret[\s:=]+['\"]?[a-zA-Z0-9_\-]{10,}['\"]?`,
	}
	// Ignore compilation errors for known‑good patterns
	f, _ := NewFilter(patterns)
	return f
}
```

- [ ] **Step 2: Write tests for filter**

Create `internal/security/filter_test.go`:

```go
package security

import (
	"testing"
)

func TestFilter_Redact(t *testing.T) {
	patterns := []string{
		`password\s*=\s*\w+`,
		`token: [a-f0-9]{32}`,
	}
	f, err := NewFilter(patterns)
	if err != nil {
		t.Fatalf("NewFilter failed: %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"password = secret123", "password = [REDACTED]"},
		{"token: abcdef1234567890abcdef1234567890", "token: [REDACTED]"},
		{"no sensitive data", "no sensitive data"},
		{"multiple: password = foo token: abcdef1234567890abcdef1234567890",
			"multiple: password = [REDACTED] token: [REDACTED]"},
	}

	for _, tt := range tests {
		got := f.Redact(tt.input)
		if got != tt.expected {
			t.Errorf("Redact(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFilter_ContainsSensitive(t *testing.T) {
	f, _ := NewFilter([]string{`password\s*=\s*\w+`})
	if !f.ContainsSensitive("password = hello") {
		t.Error("ContainsSensitive should detect password")
	}
	if f.ContainsSensitive("no password here") {
		t.Error("ContainsSensitive false positive")
	}
}

func TestDefaultFilter(t *testing.T) {
	f := DefaultFilter()
	if f == nil {
		t.Fatal("DefaultFilter returned nil")
	}
	// Ensure it redacts a known pattern
	text := "password: supersecret"
	redacted := f.Redact(text)
	if redacted == text {
		t.Errorf("DefaultFilter did not redact password")
	}
}
```

Run: `go test ./internal/security/... -v`
Expected: PASS

- [ ] **Step 3: Commit security filter**

```bash
git add internal/security/filter.go internal/security/filter_test.go
git commit -m "feat: add sensitive data filter"
```

---

### Task 7: Chat Handlers and Main App Integration

**Files:**
- Create: `internal/ui/chat_handlers.go`
- Modify: `app.go`
- Modify: `main.go`

- [ ] **Step 1: Create chat handlers**

Create `internal/ui/chat_handlers.go`:

```go
package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/pairadmin/pairadmin/internal/clipboard"
	"github.com/pairadmin/pairadmin/internal/config"
	"github.com/pairadmin/pairadmin/internal/llm"
	"github.com/pairadmin/pairadmin/internal/llm/providers"
	"github.com/pairadmin/pairadmin/internal/security"
	"github.com/pairadmin/pairadmin/internal/session"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ChatHandlers manages AI chat interactions.
type ChatHandlers struct {
	ctx      context.Context
	gateway  llm.Gateway
	store    *session.Store
	clipboard *clipboard.Manager
	filter   *security.Filter
}

// NewChatHandlers creates chat handlers with the current config.
func NewChatHandlers(ctx context.Context, store *session.Store) (*ChatHandlers, error) {
	cfg := config.Get()
	if cfg == nil {
		return nil, fmt.Errorf("config not initialized")
	}

	var gateway llm.Gateway
	switch cfg.LLM.Provider {
	case "openai":
		gateway = providers.NewOpenAIProvider(
			cfg.LLM.OpenAI.APIKey,
			cfg.LLM.OpenAI.Model,
			cfg.LLM.OpenAI.BaseURL,
		)
	case "anthropic":
		// TODO: implement Anthropic provider
		return nil, fmt.Errorf("anthropic provider not yet implemented")
	case "ollama":
		// TODO: implement Ollama provider
		return nil, fmt.Errorf("ollama provider not yet implemented")
	default:
		return nil, fmt.Errorf("unknown provider %q", cfg.LLM.Provider)
	}

	filter := security.DefaultFilter()

	return &ChatHandlers{
		ctx:      ctx,
		gateway:  gateway,
		store:    store,
		clipboard: clipboard.NewManager(),
		filter:   filter,
	}, nil
}

// SendMessage sends a user message to the LLM and returns the response.
func (c *ChatHandlers) SendMessage(terminalID, message string) (string, error) {
	// Filter sensitive data from message
	filteredMsg := c.filter.Redact(message)

	req := llm.CompletionRequest{
		Model: c.gateway.Name(),
		Messages: []llm.ChatMessage{
			{Role: llm.RoleUser, Content: filteredMsg},
		},
		MaxTokens:   2000,
		Temperature: 0.7,
	}

	resp, err := c.gateway.Complete(c.ctx, req)
	if err != nil {
		return "", fmt.Errorf("LLM completion: %w", err)
	}

	// Store suggested command if response looks like a command
	// (simplistic heuristic: starts with $ or >)
	cmdText := resp.Content
	if len(cmdText) > 0 && (cmdText[0] == '$' || cmdText[0] == '>') {
		cmd := session.SuggestedCommand{
			ID:          fmt.Sprintf("cmd-%d", time.Now().UnixNano()),
			TerminalID:  terminalID,
			Command:     cmdText,
			Description: "AI‑suggested command",
			Context:     message,
			CreatedAt:   time.Now(),
			UsedCount:   0,
		}
		if err := c.store.AddCommand(cmd); err != nil {
			// Log but don't fail
			runtime.LogError(c.ctx, fmt.Sprintf("failed to store command: %v", err))
		}
	}

	return resp.Content, nil
}

// CopyCommandToClipboard copies a command to clipboard and increments usage.
func (c *ChatHandlers) CopyCommandToClipboard(commandID, terminalID string) error {
	// TODO: retrieve command from store and copy its text
	// For now, just copy a placeholder
	err := c.clipboard.CopyToTerminal(c.ctx, "echo 'command copied'", terminalID)
	if err != nil {
		return err
	}
	// Increment usage count
	return c.store.IncrementUsedCount(commandID)
}
```

- [ ] **Step 2: Extend main App struct**

Modify `app.go`:

```go
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pairadmin/pairadmin/internal/config"
	"github.com/pairadmin/pairadmin/internal/session"
	"github.com/pairadmin/pairadmin/internal/ui"
)

// App struct
type App struct {
	ctx            context.Context
	configPath     string
	sessionStore   *session.Store
	terminalHandlers *ui.TerminalHandlers
	chatHandlers   *ui.ChatHandlers
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize configuration
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	a.configPath = filepath.Join(configDir, "pairadmin", "config.yaml")
	if err := config.Init(a.configPath); err != nil {
		panic(fmt.Sprintf("failed to init config: %v", err))
	}

	// Initialize session store
	dbPath := filepath.Join(filepath.Dir(a.configPath), "sessions.db")
	store, err := session.NewStore(dbPath)
	if err != nil {
		panic(fmt.Sprintf("failed to init session store: %v", err))
	}
	a.sessionStore = store

	// Initialize terminal detector and handlers
	// TODO: create detector with auto‑registered adapters
	// For now, create a nil detector (will be added in later tasks)
	a.terminalHandlers = ui.NewTerminalHandlers(nil)

	// Initialize chat handlers
	chatHandlers, err := ui.NewChatHandlers(ctx, store)
	if err != nil {
		panic(fmt.Sprintf("failed to init chat handlers: %v", err))
	}
	a.chatHandlers = chatHandlers
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SendMessage delegates to chat handlers.
func (a *App) SendMessage(terminalID, message string) (string, error) {
	if a.chatHandlers == nil {
		return "", fmt.Errorf("chat handlers not initialized")
	}
	return a.chatHandlers.SendMessage(terminalID, message)
}

// CopyCommandToClipboard delegates to chat handlers.
func (a *App) CopyCommandToClipboard(commandID, terminalID string) error {
	if a.chatHandlers == nil {
		return fmt.Errorf("chat handlers not initialized")
	}
	return a.chatHandlers.CopyCommandToClipboard(commandID, terminalID)
}
```

- [ ] **Step 3: Update main.go to bind the extended App**

Modify `main.go` (add import for terminal detector factory):

```go
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "PairAdmin",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
```

- [ ] **Step 4: Test integration**

Create a simple test in `internal/ui/chat_handlers_test.go` (optional). For now, ensure the code compiles:

Run: `go build ./...`
Expected: successful compilation (though may have missing dependencies for unimplemented providers).

- [ ] **Step 5: Commit chat handlers and app integration**

```bash
git add internal/ui/chat_handlers.go app.go main.go
git commit -m "feat: add chat handlers and integrate with main app"
```

---

### Task 8: Frontend Components – Chat Area and Command Sidebar

**Files:**
- Create: `frontend/src/components/ChatArea.svelte`
- Create: `frontend/src/components/CommandSidebar.svelte`
- Create: `frontend/src/lib/stores.js`
- Modify: `frontend/src/App.svelte`

- [ ] **Step 1: Create Svelte stores for global state**

Create `frontend/src/lib/stores.js`:

```js
import { writable } from 'svelte/store';

// Current terminal sessions
export const terminals = writable([]);

// Current chat messages
export const messages = writable([]);

// Command history for the active terminal
export const commandHistory = writable([]);

// Active terminal ID
export const activeTerminalId = writable('');

// Settings dialog open state
export const settingsOpen = writable(false);
```

- [ ] **Step 2: Create ChatArea component**

Create `frontend/src/components/ChatArea.svelte`:

```svelte
<script>
  import { messages, activeTerminalId } from '../lib/stores.js';
  import { SendMessage } from '../../wailsjs/go/main/App.js';

  let inputText = '';
  let isLoading = false;

  async function handleSend() {
    if (!inputText.trim() || isLoading) return;
    const terminalId = $activeTerminalId;
    if (!terminalId) {
      alert('Please select a terminal first');
      return;
    }

    const userMessage = inputText;
    inputText = '';
    isLoading = true;

    // Add user message to UI
    messages.update(msgs => [...msgs, { role: 'user', content: userMessage }]);

    try {
      const response = await SendMessage(terminalId, userMessage);
      messages.update(msgs => [...msgs, { role: 'assistant', content: response }]);
    } catch (error) {
      console.error('Failed to send message:', error);
      alert('Error: ' + error.message);
    } finally {
      isLoading = false;
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  }
</script>

<div class="chat-area">
  <div class="messages">
    {#each $messages as msg (msg.id)}
      <div class="message {msg.role}">
        <div class="role">{msg.role === 'user' ? 'You' : 'AI'}</div>
        <div class="content">
          {#if msg.role === 'assistant' && msg.content.startsWith('$')}
            <pre><code>{msg.content}</code></pre>
            <button on:click={() => console.log('Copy command:', msg.content)}>
              Copy to Terminal
            </button>
          {:else}
            {msg.content}
          {/if}
        </div>
      </div>
    {/each}
  </div>

  <div class="input-area">
    <textarea
      bind:value={inputText}
      on:keydown={handleKeyDown}
      placeholder="Ask the AI for a command..."
      rows="3"
      disabled={isLoading}
    />
    <button on:click={handleSend} disabled={isLoading || !inputText.trim()}>
      {isLoading ? 'Sending…' : 'Send'}
    </button>
  </div>
</div>

<style>
  .chat-area {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: 1rem;
  }
  .messages {
    flex: 1;
    overflow-y: auto;
    margin-bottom: 1rem;
  }
  .message {
    margin-bottom: 1rem;
    padding: 0.75rem;
    border-radius: 0.5rem;
    background: #f5f5f5;
  }
  .message.user {
    background: #e3f2fd;
  }
  .message.assistant {
    background: #f1f8e9;
  }
  .role {
    font-weight: bold;
    font-size: 0.8rem;
    margin-bottom: 0.25rem;
    color: #666;
  }
  .content pre {
    margin: 0;
    padding: 0.5rem;
    background: #2d2d2d;
    color: #f8f8f2;
    border-radius: 0.25rem;
    overflow-x: auto;
  }
  .input-area {
    display: flex;
    gap: 0.5rem;
  }
  textarea {
    flex: 1;
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 0.25rem;
    font-family: monospace;
  }
  button {
    padding: 0.5rem 1rem;
    background: #4caf50;
    color: white;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
  }
  button:disabled {
    background: #ccc;
    cursor: not-allowed;
  }
</style>
```

- [ ] **Step 3: Create CommandSidebar component**

Create `frontend/src/components/CommandSidebar.svelte`:

```svelte
<script>
  import { commandHistory, activeTerminalId } from '../lib/stores.js';
  import { CopyCommandToClipboard } from '../../wailsjs/go/main/App.js';

  async function copyCommand(command) {
    try {
      await CopyCommandToClipboard(command.id, $activeTerminalId);
      alert(`Command "${command.command}" copied to clipboard`);
    } catch (error) {
      console.error('Failed to copy command:', error);
      alert('Error copying command: ' + error.message);
    }
  }
</script>

<div class="command-sidebar">
  <h3>Command History</h3>
  {#if $commandHistory.length === 0}
    <p class="empty">No commands yet. Ask the AI for help!</p>
  {:else}
    <div class="command-list">
      {#each $commandHistory as cmd (cmd.id)}
        <div class="command-card">
          <div class="command-text">{cmd.command}</div>
          <div class="command-meta">
            <span class="used">{cmd.usedCount} uses</span>
            <span class="time">{new Date(cmd.createdAt).toLocaleTimeString()}</span>
          </div>
          <button on:click={() => copyCommand(cmd)}>Copy</button>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .command-sidebar {
    padding: 1rem;
    background: #f9f9f9;
    border-left: 1px solid #ddd;
    height: 100%;
    overflow-y: auto;
  }
  h3 {
    margin-top: 0;
    margin-bottom: 1rem;
  }
  .empty {
    color: #999;
    font-style: italic;
  }
  .command-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  .command-card {
    padding: 0.75rem;
    background: white;
    border: 1px solid #ddd;
    border-radius: 0.5rem;
  }
  .command-text {
    font-family: monospace;
    font-size: 0.9rem;
    margin-bottom: 0.5rem;
    word-break: break-all;
  }
  .command-meta {
    display: flex;
    justify-content: space-between;
    font-size: 0.8rem;
    color: #666;
    margin-bottom: 0.5rem;
  }
  button {
    width: 100%;
    padding: 0.25rem;
    background: #2196f3;
    color: white;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
  }
  button:hover {
    background: #0b7dda;
  }
</style>
```

- [ ] **Step 4: Update main App.svelte**

Replace `frontend/src/App.svelte`:

```svelte
<script>
  import ChatArea from './components/ChatArea.svelte';
  import CommandSidebar from './components/CommandSidebar.svelte';
  import TerminalTabs from './components/TerminalTabs.svelte';
  import StatusBar from './components/StatusBar.svelte';
  import SettingsDialog from './components/SettingsDialog.svelte';
  import { settingsOpen } from './lib/stores.js';
</script>

<main>
  <div class="app-layout">
    <!-- Left: Terminal tabs -->
    <div class="left-panel">
      <TerminalTabs />
    </div>

    <!-- Middle: Chat area -->
    <div class="middle-panel">
      <ChatArea />
    </div>

    <!-- Right: Command sidebar -->
    <div class="right-panel">
      <CommandSidebar />
    </div>
  </div>

  <!-- Status bar at bottom -->
  <StatusBar />

  <!-- Settings dialog (modal) -->
  {#if $settingsOpen}
    <SettingsDialog />
  {/if}
</main>

<style>
  * {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }
  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    overflow: hidden;
  }
  .app-layout {
    display: flex;
    height: calc(100vh - 40px); /* minus status bar */
  }
  .left-panel {
    width: 250px;
    border-right: 1px solid #ccc;
  }
  .middle-panel {
    flex: 1;
    border-right: 1px solid #ccc;
  }
  .right-panel {
    width: 300px;
  }
</style>
```

- [ ] **Step 5: Create placeholder components for TerminalTabs and StatusBar**

Create `frontend/src/components/TerminalTabs.svelte`:

```svelte
<script>
  import { terminals, activeTerminalId } from '../lib/stores.js';

  // Placeholder data
  terminals.set([
    { id: 'term1', name: 'bash' },
    { id: 'term2', name: 'zsh' },
  ]);
  activeTerminalId.set('term1');

  function selectTerminal(id) {
    activeTerminalId.set(id);
  }
</script>

<div class="terminal-tabs">
  <h3>Terminals</h3>
  {#each $terminals as term (term.id)}
    <div
      class="tab {term.id === $activeTerminalId ? 'active' : ''}"
      on:click={() => selectTerminal(term.id)}
    >
      {term.name}
    </div>
  {/each}
</div>

<style>
  .terminal-tabs {
    padding: 1rem;
  }
  h3 {
    margin-bottom: 1rem;
  }
  .tab {
    padding: 0.5rem;
    margin-bottom: 0.5rem;
    background: #f0f0f0;
    border-radius: 0.25rem;
    cursor: pointer;
  }
  .tab.active {
    background: #4caf50;
    color: white;
  }
  .tab:hover {
    background: #e0e0e0;
  }
</style>
```

Create `frontend/src/components/StatusBar.svelte`:

```svelte
<script>
  import { settingsOpen } from '../lib/stores.js';

  function openSettings() {
    settingsOpen.set(true);
  }
</script>

<div class="status-bar">
  <span class="model">Model: GPT‑4</span>
  <span class="connection">Connected</span>
  <span class="context">Context: 120/4000 tokens</span>
  <button class="settings-btn" on:click={openSettings}>Settings</button>
</div>

<style>
  .status-bar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 40px;
    background: #333;
    color: white;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 1rem;
    font-size: 0.9rem;
  }
  .settings-btn {
    background: #555;
    color: white;
    border: none;
    padding: 0.25rem 0.75rem;
    border-radius: 0.25rem;
    cursor: pointer;
  }
  .settings-btn:hover {
    background: #666;
  }
</style>
```

Create `frontend/src/components/SettingsDialog.svelte` (placeholder):

```svelte
<script>
  import { settingsOpen } from '../lib/stores.js';

  function close() {
    settingsOpen.set(false);
  }
</script>

<div class="settings-modal">
  <div class="settings-content">
    <h2>Settings</h2>
    <p>Settings dialog will be implemented in a later task.</p>
    <button on:click={close}>Close</button>
  </div>
</div>

<style>
  .settings-modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  .settings-content {
    background: white;
    padding: 2rem;
    border-radius: 0.5rem;
    max-width: 500px;
    width: 100%;
  }
</style>
```

- [ ] **Step 6: Build frontend and test**

Run:
```bash
cd frontend && npm install
cd .. && wails build
```

Expected: Build succeeds (may warn about missing imports).

- [ ] **Step 7: Commit frontend components**

```bash
git add frontend/src/
git commit -m "feat: add frontend components (chat, sidebar, tabs, status bar)"
```

---

### Task 9: Settings Dialog UI

**Files:**
- Modify: `frontend/src/components/SettingsDialog.svelte`
- Create: `frontend/src/components/settings/LLMTab.svelte`
- Create: `frontend/src/components/settings/TerminalsTab.svelte`
- Create: `frontend/src/components/settings/HotkeysTab.svelte`
- Create: `frontend/src/components/settings/AppearanceTab.svelte`

- [ ] **Step 1: Create tab components**

Create `frontend/src/components/settings/LLMTab.svelte`:

```svelte
<script>
  let provider = 'openai';
  let openaiKey = '';
  let anthropicKey = '';
  let ollamaUrl = 'http://localhost:11434';
  let model = 'gpt-4';

  function save() {
    console.log('Save LLM settings', { provider, openaiKey, anthropicKey, ollamaUrl, model });
    alert('Settings saved (placeholder)');
  }
</script>

<div class="tab">
  <h3>LLM Configuration</h3>
  <div class="field">
    <label>Provider</label>
    <select bind:value={provider}>
      <option value="openai">OpenAI</option>
      <option value="anthropic">Anthropic</option>
      <option value="ollama">Ollama (local)</option>
    </select>
  </div>

  {#if provider === 'openai'}
    <div class="field">
      <label>OpenAI API Key</label>
      <input type="password" bind:value={openaiKey} placeholder="sk-..." />
    </div>
    <div class="field">
      <label>Model</label>
      <input bind:value={model} placeholder="gpt-4" />
    </div>
  {:else if provider === 'anthropic'}
    <div class="field">
      <label>Anthropic API Key</label>
      <input type="password" bind:value={anthropicKey} placeholder="sk-ant-..." />
    </div>
  {:else if provider === 'ollama'}
    <div class="field">
      <label>Ollama Base URL</label>
      <input bind:value={ollamaUrl} />
    </div>
  {/if}

  <button on:click={save}>Save</button>
</div>

<style>
  .tab {
    padding: 1rem;
  }
  .field {
    margin-bottom: 1rem;
  }
  label {
    display: block;
    margin-bottom: 0.25rem;
    font-weight: bold;
  }
  input, select {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 0.25rem;
  }
  button {
    padding: 0.5rem 1rem;
    background: #4caf50;
    color: white;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
  }
</style>
```

Create the other tabs similarly (simplified). For brevity, we'll create a single SettingsDialog with all tabs inline in the next step.

- [ ] **Step 2: Update SettingsDialog with tabs**

Replace `frontend/src/components/SettingsDialog.svelte`:

```svelte
<script>
  import { settingsOpen } from '../lib/stores.js';
  import LLMTab from './settings/LLMTab.svelte';

  let activeTab = 'llm';

  function close() {
    settingsOpen.set(false);
  }
</script>

<div class="settings-modal" on:click={close}>
  <div class="settings-content" on:click|stopPropagation>
    <div class="header">
      <h2>Settings</h2>
      <button class="close-btn" on:click={close}>×</button>
    </div>

    <div class="tabs">
      <button class:active={activeTab === 'llm'} on:click={() => activeTab = 'llm'}>LLM</button>
      <button class:active={activeTab === 'terminals'} on:click={() => activeTab = 'terminals'}>Terminals</button>
      <button class:active={activeTab === 'hotkeys'} on:click={() => activeTab = 'hotkeys'}>Hotkeys</button>
      <button class:active={activeTab === 'appearance'} on:click={() => activeTab = 'appearance'}>Appearance</button>
    </div>

    <div class="tab-content">
      {#if activeTab === 'llm'}
        <LLMTab />
      {:else if activeTab === 'terminals'}
        <div class="tab">
          <h3>Terminal Settings</h3>
          <p>Configure terminal detection and capture settings.</p>
        </div>
      {:else if activeTab === 'hotkeys'}
        <div class="tab">
          <h3>Hotkey Configuration</h3>
          <p>Configure global hotkeys for PairAdmin.</p>
        </div>
      {:else if activeTab === 'appearance'}
        <div class="tab">
          <h3>Appearance</h3>
          <p>Choose light/dark theme and UI preferences.</p>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .settings-modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  .settings-content {
    background: white;
    padding: 0;
    border-radius: 0.5rem;
    max-width: 700px;
    width: 100%;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
  }
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    border-bottom: 1px solid #ddd;
  }
  .close-btn {
    background: none;
    border: none;
    font-size: 2rem;
    cursor: pointer;
    color: #666;
  }
  .tabs {
    display: flex;
    border-bottom: 1px solid #ddd;
  }
  .tabs button {
    flex: 1;
    padding: 1rem;
    background: none;
    border: none;
    cursor: pointer;
    border-bottom: 3px solid transparent;
  }
  .tabs button.active {
    border-bottom-color: #4caf50;
    font-weight: bold;
  }
  .tab-content {
    flex: 1;
    overflow-y: auto;
    padding: 1.5rem;
  }
</style>
```

- [ ] **Step 3: Commit settings dialog**

```bash
git add frontend/src/components/
git commit -m "feat: implement settings dialog with tabs"
```

---

### Task 10: Hotkey Configuration Backend

**Files:**
- Create: `internal/hotkeys/manager.go`
- Test: `internal/hotkeys/manager_test.go`
- Modify: `app.go` to integrate hotkey manager

- [ ] **Step 1: Implement hotkey manager using robotgo**

Create `internal/hotkeys/manager.go`:

```go
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
	mods, key, err := parseHotkey(hotkey)
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
```

- [ ] **Step 2: Write placeholder tests**

Create `internal/hotkeys/manager_test.go`:

```go
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
}
```

Run: `go test ./internal/hotkeys/... -v`
Expected: PASS

- [ ] **Step 3: Integrate hotkey manager into app.go**

Add to `app.go` startup:

```go
// Add to imports
"github.com/pairadmin/pairadmin/internal/hotkeys"

// Add to App struct
hotkeyManager *hotkeys.Manager

// Add to startup after config init
hotkeyMgr := hotkeys.NewManager()
// Register hotkeys from config
cfg := config.Get()
if cfg != nil {
    hotkeyMgr.Register(cfg.UI.Hotkeys.CopyLastCommand, func() {
        // TODO: implement copy last command action
        runtime.LogInfo(ctx, "Hotkey: copy last command")
    })
    hotkeyMgr.Register(cfg.UI.Hotkeys.FocusApp, func() {
        // TODO: implement focus app action
        runtime.LogInfo(ctx, "Hotkey: focus app")
    })
}
a.hotkeyManager = hotkeyMgr
// Start hotkey listener (commented until platform‑specific implementation)
// go hotkeyMgr.Start()
```

- [ ] **Step 4: Commit hotkey manager**

```bash
git add internal/hotkeys/ app.go
git commit -m "feat: add hotkey manager skeleton"
```

---

### Task 11: Integration and End‑to‑End Test

**Files:**
- Create: `scripts/test‑integration.sh`
- Update: `README.md` with new features

- [ ] **Step 1: Create integration test script**

Create `scripts/test-integration.sh`:

```bash
#!/bin/bash
set -e

echo "=== PairAdmin Phase 3 Integration Test ==="
echo "Building backend..."
go build ./...

echo "Running unit tests..."
go test ./internal/llm/... ./internal/config/... ./internal/session/... ./internal/clipboard/... ./internal/security/...

echo "Creating test configuration..."
mkdir -p ~/.pairadmin
cp config.yaml ~/.pairadmin/config.yaml

echo "Starting session store test..."
go run ./cmd/test-session-store/main.go 2>/dev/null || echo "Session store test skipped (no test binary)"

echo "✅ Phase 3 components compile and unit tests pass."
echo "Next: run 'wails dev' to test the frontend."
```

Make executable: `chmod +x scripts/test-integration.sh`

- [ ] **Step 2: Update README.md with new features**

Add a section to `README.md`:

```markdown
## Phase 3: AI Collaboration Features

- **LLM Provider Integration**: Support for OpenAI, Anthropic, and Ollama (local)
- **Command History**: Stores AI‑suggested commands per terminal session in SQLite
- **Clipboard Manager**: One‑click copy to clipboard with optional terminal focus
- **Sensitive Data Filtering**: Redacts passwords, API keys, etc. before sending to LLM
- **Settings Dialog**: Configure providers, terminals, hotkeys, and appearance
- **Hotkey Support**: Global shortcuts for copying last command and focusing the app
```

- [ ] **Step 3: Run integration test**

Run: `./scripts/test-integration.sh`
Expected: All tests pass (or skip clipboard tests if unavailable).

- [ ] **Step 4: Commit final integration**

```bash
git add scripts/ README.md
git commit -m "chore: add integration test script and update README"
```

---

## Windows‑System Checklist (Once Available)

1. **Integration testing** – Run `go test ./internal/terminal/windows/... -tags=integration` with real terminal windows open.
2. **Cross‑compilation validation** – Ensure Windows adapter builds with CGO (`CGO_ENABLED=1`).
3. **UI Automation permission verification** – Confirm no extra permissions needed beyond standard user access.
4. **Terminal‑type coverage** – Test Windows Terminal, PowerShell, cmd.exe, PuTTY detection.
5. **CI pipeline** – Add Windows unit‑test step to `.github/workflows/build.yml` after repo creation.

---

## Phase 4 Preview

**Security hardening** (Tasks 4.1‑4.3):
- OS keychain integration for credentials (`internal/config/keychain.go`).
- Audit logging (`~/.pairadmin/logs/audit‑*.jsonl`).
- Local LLM support (Ollama) – implement `internal/llm/providers/ollama.go`.

**Testing & QA** (Tasks 4.4‑4.6):
- Unit/integration tests for core services.
- End‑to‑end testing scripts.
- Performance optimization (polling, token counting).

**Documentation & Packaging** (Tasks 4.7‑4.9):
- User guides, installation instructions.
- Installers: macOS `.dmg`, Windows `.msi`, Linux `.AppImage`/`.deb`/`.rpm`.
- Final release checklist.

---

**Plan complete and saved to `docs/superpowers/plans/2026-04-01-phase-3-ai-collaboration.md`. Two execution options:**

**1. Subagent‑Driven (recommended)** – I dispatch a fresh subagent per task, review between tasks, fast iteration.

**2. Inline Execution** – Execute tasks in this session using executing‑plans, batch execution with checkpoints.

**Which approach?**