package ui

import (
	"context"
	"database/sql"
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

type contextKey int

const (
	testModeKey contextKey = iota
)

type store interface {
	AddSession(sessionID, terminalID string) error
	GetSession(sessionID string) (*session.Session, error)
	AddCommand(cmd session.SuggestedCommand) error
	GetCommandByID(commandID string) (*session.SuggestedCommand, error)
	GetCommandsByTerminal(terminalID string) ([]session.SuggestedCommand, error)
	IncrementUsedCount(commandID string) error
}

type clipboardManager interface {
	CopyToTerminal(ctx context.Context, text string, terminalID string) error
}

type stubProvider struct {
	name string
}

func (s *stubProvider) Name() string { return s.name }

func (s *stubProvider) Complete(ctx context.Context, req llm.CompletionRequest) (*llm.CompletionResponse, error) {
	return nil, fmt.Errorf("%s provider not yet implemented", s.name)
}

func (s *stubProvider) StreamComplete(ctx context.Context, req llm.CompletionRequest) (<-chan string, error) {
	return nil, fmt.Errorf("%s provider not yet implemented", s.name)
}

func createFilterFromConfig(ctx context.Context, cfg *config.Config) *security.Filter {
	if cfg == nil || len(cfg.Security.FilterPatterns) == 0 {
		return security.DefaultFilter()
	}
	var rawPatterns []string
	for _, fp := range cfg.Security.FilterPatterns {
		rawPatterns = append(rawPatterns, fp.Pattern)
	}
	filter, err := security.NewFilter(rawPatterns)
	if err != nil {
		// Skip logging in test mode
		if ctx.Value(testModeKey) == nil {
			runtime.LogError(ctx, fmt.Sprintf("failed to compile custom filter patterns: %v, falling back to default", err))
		}
		return security.DefaultFilter()
	}
	return filter
}

// ChatHandlers manages AI chat interactions.
type ChatHandlers struct {
	ctx       context.Context
	gateway   llm.Gateway
	store     store
	clipboard clipboardManager
	filter    *security.Filter
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
		runtime.LogWarning(ctx, fmt.Sprintf("Provider %q not yet implemented, using stub provider", cfg.LLM.Provider))
		gateway = &stubProvider{name: cfg.LLM.Provider}
	case "ollama":
		gateway = providers.NewOllamaProvider(
			cfg.LLM.Ollama.BaseURL,
			cfg.LLM.Ollama.Model,
		)
	default:
		return nil, fmt.Errorf("unknown provider %q", cfg.LLM.Provider)
	}

	filter := createFilterFromConfig(ctx, cfg)

	return &ChatHandlers{
		ctx:       ctx,
		gateway:   gateway,
		store:     store,
		clipboard: clipboard.NewManager(),
		filter:    filter,
	}, nil
}

// SendMessage sends a user message to the LLM and returns the response.
func (c *ChatHandlers) SendMessage(terminalID, message string) (string, string, error) {
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

	ctx, cancel := context.WithTimeout(c.ctx, 30*time.Second)
	defer cancel()

	resp, err := c.gateway.Complete(ctx, req)
	if err != nil {
		return "", "", fmt.Errorf("LLM completion: %w", err)
	}

	// Store suggested command if response looks like a command
	// (simplistic heuristic: starts with $ or >)
	cmdText := resp.Content
	commandID := ""
	if len(cmdText) > 0 && (cmdText[0] == '$' || cmdText[0] == '>') {
		// Ensure a session exists for this terminal (session ID = terminalID)
		_, sessionErr := c.store.GetSession(terminalID)
		if sessionErr != nil && sessionErr != sql.ErrNoRows {
			runtime.LogError(c.ctx, fmt.Sprintf("failed to check session: %v", sessionErr))
		} else if sessionErr == sql.ErrNoRows {
			// Create session
			if err := c.store.AddSession(terminalID, terminalID); err != nil {
				runtime.LogError(c.ctx, fmt.Sprintf("failed to create session: %v", err))
			}
		}

		cmd := session.SuggestedCommand{
			ID:          fmt.Sprintf("cmd-%d", time.Now().UnixNano()),
			SessionID:   terminalID,
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
		} else {
			commandID = cmd.ID
		}
	}

	return resp.Content, commandID, nil
}

// CopyCommandToClipboard copies a command to clipboard and increments usage.
func (c *ChatHandlers) CopyCommandToClipboard(commandID, terminalID string) error {
	cmd, err := c.store.GetCommandByID(commandID)
	if err != nil {
		return fmt.Errorf("retrieve command: %w", err)
	}
	err = c.clipboard.CopyToTerminal(c.ctx, cmd.Command, terminalID)
	if err != nil {
		return err
	}
	// Increment usage count
	return c.store.IncrementUsedCount(commandID)
}
