package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pairadmin/pairadmin/internal/llm"
)

// OllamaProvider implements llm.Gateway for local Ollama API.
type OllamaProvider struct {
	baseURL string
	model   string
	client  *http.Client
}

// NewOllamaProvider creates a new Ollama provider.
func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	return &OllamaProvider{
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// Name returns "ollama".
func (p *OllamaProvider) Name() string {
	return "ollama"
}

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Model      string `json:"model"`
	Response   string `json:"response"`
	Done       bool   `json:"done"`
	DoneReason string `json:"done_reason"`
}

// Complete sends a completion request to Ollama.
func (p *OllamaProvider) Complete(ctx context.Context, req llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if p.baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if p.model == "" {
		return nil, fmt.Errorf("model is required")
	}
	ollamaReq := ollamaRequest{
		Model:  p.model,
		Stream: false,
	}
	for _, msg := range req.Messages {
		ollamaReq.Messages = append(ollamaReq.Messages, ollamaMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", normalizeBaseURL(p.baseURL)+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limitedReader := io.LimitReader(resp.Body, 1<<20) // 1MB
		bodyBytes, _ := io.ReadAll(limitedReader)
		// Drain remaining body to allow connection reuse
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, fmt.Errorf("ollama api error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	llmResp := &llm.CompletionResponse{
		Content: ollamaResp.Response,
		Model:   ollamaResp.Model,
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{}, // Ollama doesn't provide token counts
	}
	return llmResp, nil
}

// StreamComplete is not implemented yet.
func (p *OllamaProvider) StreamComplete(ctx context.Context, req llm.CompletionRequest) (<-chan string, error) {
	return nil, fmt.Errorf("streaming not yet implemented for Ollama")
}
