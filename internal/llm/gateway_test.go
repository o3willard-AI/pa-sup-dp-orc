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

// mockGateway is a test implementation of Gateway
type mockGateway struct{}

func (m *mockGateway) Name() string {
	return "mock"
}

func (m *mockGateway) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	return &CompletionResponse{
		Content: "mock response",
		Model:   req.Model,
	}, nil
}

func (m *mockGateway) StreamComplete(ctx context.Context, req CompletionRequest) (<-chan string, error) {
	ch := make(chan string, 1)
	ch <- "mock stream"
	close(ch)
	return ch, nil
}

func TestMockGatewayImplementsInterface(t *testing.T) {
	var _ Gateway = (*mockGateway)(nil)
	// If this compiles, the interface is satisfied
}

func TestMockGatewayMethods(t *testing.T) {
	ctx := context.Background()
	mock := &mockGateway{}

	if name := mock.Name(); name != "mock" {
		t.Errorf("Expected name 'mock', got %s", name)
	}

	req := CompletionRequest{
		Model: "test-model",
	}
	resp, err := mock.Complete(ctx, req)
	if err != nil {
		t.Errorf("Complete returned error: %v", err)
	}
	if resp.Content != "mock response" {
		t.Errorf("Expected content 'mock response', got %s", resp.Content)
	}
	if resp.Model != req.Model {
		t.Errorf("Expected model %s, got %s", req.Model, resp.Model)
	}

	streamCh, err := mock.StreamComplete(ctx, req)
	if err != nil {
		t.Errorf("StreamComplete returned error: %v", err)
	}
	msg, ok := <-streamCh
	if !ok {
		t.Error("Stream channel closed unexpectedly")
	}
	if msg != "mock stream" {
		t.Errorf("Expected stream message 'mock stream', got %s", msg)
	}
}

func TestGatewayEdgeCases(t *testing.T) {
	// Test that nil Gateway can't be used (compile-time)
	// This is just a placeholder for future edge case tests
	// For example, testing with empty request, cancelled context, etc.
	// Since we don't have concrete implementations, we can't test those.
	t.Log("Edge case tests would require concrete implementations")
}