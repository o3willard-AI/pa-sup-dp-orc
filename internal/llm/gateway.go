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
