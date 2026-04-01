package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
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
