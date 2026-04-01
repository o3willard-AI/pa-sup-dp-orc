package ui

import (
	"context"
	"testing"

	"github.com/pairadmin/pairadmin/internal/session"
)

func TestNewChatHandlers_ConfigNotInitialized(t *testing.T) {
	ctx := context.Background()
	store, err := session.NewStore(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory store: %v", err)
	}
	defer store.Close()

	_, err = NewChatHandlers(ctx, store)
	if err == nil {
		t.Error("expected error when config not initialized, got nil")
	}
	if err.Error() != "config not initialized" {
		t.Errorf("expected error 'config not initialized', got: %v", err)
	}
}