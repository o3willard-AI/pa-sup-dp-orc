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
	ctx       context.Context
	gateway   llm.Gateway
	store     *session.Store
	clipboard *clipboard.Manager
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
		return nil, fmt.Errorf("anthropic provider not yet implemented")
	case "ollama":
		// TODO: implement Ollama provider
		return nil, fmt.Errorf("ollama provider not yet implemented")
	default:
		return nil, fmt.Errorf("unknown provider %q", cfg.LLM.Provider)
	}

	filter := security.DefaultFilter()

	return &ChatHandlers{
		ctx:       ctx,
		gateway:   gateway,
		store:     store,
		clipboard: clipboard.NewManager(),
		filter:    filter,
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
			SessionID:   terminalID, // Using terminalID as session ID for now
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