package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestOpenAIProvider_Complete_Validation(t *testing.T) {
	p := NewOpenAIProvider("", "gpt-4", "https://api.openai.com/v1")
	req := llm.CompletionRequest{Messages: []llm.ChatMessage{{Role: llm.RoleUser, Content: "Hi"}}}
	_, err := p.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
	if !strings.Contains(err.Error(), "API key is required") {
		t.Errorf("expected error about API key, got %v", err)
	}

	p = NewOpenAIProvider("key", "", "https://api.openai.com/v1")
	_, err = p.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty model")
	}
	if !strings.Contains(err.Error(), "model is required") {
		t.Errorf("expected error about model, got %v", err)
	}

	p = NewOpenAIProvider("key", "gpt-4", "")
	_, err = p.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty base URL")
	}
	if !strings.Contains(err.Error(), "base URL is required") {
		t.Errorf("expected error about base URL, got %v", err)
	}
}

func TestOpenAIProvider_Complete_NetworkError(t *testing.T) {
	p := NewOpenAIProvider("key", "gpt-4", "http://localhost:9999")
	req := llm.CompletionRequest{Messages: []llm.ChatMessage{{Role: llm.RoleUser, Content: "Hi"}}}
	_, err := p.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected network error")
	}
}

func TestOpenAIProvider_Complete_ContextCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	p := NewOpenAIProvider("key", "gpt-4", server.URL)
	req := llm.CompletionRequest{Messages: []llm.ChatMessage{{Role: llm.RoleUser, Content: "Hi"}}}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := p.Complete(ctx, req)
	if err == nil {
		t.Fatal("expected error due to context timeout")
	}
}

func TestOpenAIProvider_Complete_MalformedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	p := NewOpenAIProvider("key", "gpt-4", server.URL)
	req := llm.CompletionRequest{Messages: []llm.ChatMessage{{Role: llm.RoleUser, Content: "Hi"}}}
	_, err := p.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
	if !strings.Contains(err.Error(), "decode response") {
		t.Errorf("expected decode error, got %v", err)
	}
}

func TestOpenAIProvider_Complete_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"created": 1677652288,
			"model": "gpt-4",
			"choices": [],
			"usage": {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
		}`
		w.Write([]byte(response))
	}))
	defer server.Close()

	p := NewOpenAIProvider("key", "gpt-4", server.URL)
	req := llm.CompletionRequest{Messages: []llm.ChatMessage{{Role: llm.RoleUser, Content: "Hi"}}}
	_, err := p.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty choices")
	}
	if !strings.Contains(err.Error(), "no choices") {
		t.Errorf("expected 'no choices' error, got %v", err)
	}
}

func TestOpenAIProvider_Complete_ErrorBodyLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		size := 1<<20 + 100
		data := make([]byte, size)
		for i := range data {
			data[i] = 'A'
		}
		w.Write(data)
	}))
	defer server.Close()

	p := NewOpenAIProvider("key", "gpt-4", server.URL)
	req := llm.CompletionRequest{Messages: []llm.ChatMessage{{Role: llm.RoleUser, Content: "Hi"}}}
	_, err := p.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for server error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error containing 500, got %v", err)
	}
}
