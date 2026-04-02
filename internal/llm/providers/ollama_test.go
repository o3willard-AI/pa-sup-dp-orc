package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pairadmin/pairadmin/internal/llm"
)

func TestOllamaProvider_Name(t *testing.T) {
	p := NewOllamaProvider("http://localhost:11434", "llama3")
	if p.Name() != "ollama" {
		t.Errorf("Expected name 'ollama', got %s", p.Name())
	}
}

func TestOllamaProvider_Complete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"model": "llama3",
			"created_at": "2024-01-01T00:00:00Z",
			"response": "Hello!",
			"done": true,
			"done_reason": "stop"
		}`
		w.Write([]byte(response))
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "llama3")
	req := llm.CompletionRequest{
		Model: "llama3",
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
	if resp.Model != "llama3" {
		t.Errorf("Expected model 'llama3', got %s", resp.Model)
	}
}

func TestOllamaProvider_Complete_Validation(t *testing.T) {
	p := NewOllamaProvider("", "llama3")
	req := llm.CompletionRequest{Messages: []llm.ChatMessage{{Role: llm.RoleUser, Content: "Hi"}}}
	_, err := p.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty base URL")
	}
	if !strings.Contains(err.Error(), "base URL is required") {
		t.Errorf("expected error about base URL, got %v", err)
	}

	p = NewOllamaProvider("http://localhost:11434", "")
	_, err = p.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty model")
	}
	if !strings.Contains(err.Error(), "model is required") {
		t.Errorf("expected error about model, got %v", err)
	}
}
