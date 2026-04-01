package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// Name returns "openai".
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// openAIRequest is the JSON structure for OpenAI completions.
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
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

// normalizeBaseURL ensures the baseURL does not end with a slash.
func normalizeBaseURL(baseURL string) string {
	return strings.TrimSuffix(baseURL, "/")
}

// Complete sends a completion request to OpenAI.
func (p *OpenAIProvider) Complete(ctx context.Context, req llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	if p.model == "" {
		return nil, fmt.Errorf("model is required")
	}
	if p.baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
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

	httpReq, err := http.NewRequestWithContext(ctx, "POST", normalizeBaseURL(p.baseURL)+"/chat/completions", bytes.NewReader(body))
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
		limitedReader := io.LimitReader(resp.Body, 1<<20) // 1MB
		bodyBytes, _ := io.ReadAll(limitedReader)
		// Drain remaining body to allow connection reuse
		io.Copy(io.Discard, resp.Body)
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
