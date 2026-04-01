package ui

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/pairadmin/pairadmin/internal/config"
	"github.com/pairadmin/pairadmin/internal/llm"
	"github.com/pairadmin/pairadmin/internal/security"
	"github.com/pairadmin/pairadmin/internal/session"
)

func testContext(t *testing.T) context.Context {
	return context.WithValue(context.Background(), testModeKey, true)
}

type mockGateway struct {
	completeFunc func(ctx context.Context, req llm.CompletionRequest) (*llm.CompletionResponse, error)
	capturedReq  *llm.CompletionRequest
}

func (m *mockGateway) Name() string { return "mock" }
func (m *mockGateway) Complete(ctx context.Context, req llm.CompletionRequest) (*llm.CompletionResponse, error) {
	m.capturedReq = &req
	if m.completeFunc != nil {
		return m.completeFunc(ctx, req)
	}
	return &llm.CompletionResponse{Content: "test response"}, nil
}
func (m *mockGateway) StreamComplete(ctx context.Context, req llm.CompletionRequest) (<-chan string, error) {
	ch := make(chan string)
	close(ch)
	return ch, nil
}

type mockStore struct {
	addSessionFunc         func(sessionID, terminalID string) error
	getSessionFunc         func(sessionID string) (*session.Session, error)
	addCommandFunc         func(cmd session.SuggestedCommand) error
	getCommandByIDFunc     func(commandID string) (*session.SuggestedCommand, error)
	incrementUsedCountFunc func(commandID string) error
}

func (m *mockStore) AddSession(sessionID, terminalID string) error {
	if m.addSessionFunc != nil {
		return m.addSessionFunc(sessionID, terminalID)
	}
	return nil
}
func (m *mockStore) GetSession(sessionID string) (*session.Session, error) {
	if m.getSessionFunc != nil {
		return m.getSessionFunc(sessionID)
	}
	return nil, nil
}
func (m *mockStore) AddCommand(cmd session.SuggestedCommand) error {
	if m.addCommandFunc != nil {
		return m.addCommandFunc(cmd)
	}
	return nil
}
func (m *mockStore) GetCommandByID(commandID string) (*session.SuggestedCommand, error) {
	if m.getCommandByIDFunc != nil {
		return m.getCommandByIDFunc(commandID)
	}
	return nil, nil
}
func (m *mockStore) IncrementUsedCount(commandID string) error {
	if m.incrementUsedCountFunc != nil {
		return m.incrementUsedCountFunc(commandID)
	}
	return nil
}

type mockClipboard struct {
	copyToTerminalFunc func(ctx context.Context, text string, terminalID string) error
}

func (m *mockClipboard) CopyToTerminal(ctx context.Context, text string, terminalID string) error {
	if m.copyToTerminalFunc != nil {
		return m.copyToTerminalFunc(ctx, text, terminalID)
	}
	return nil
}

func setupTestConfig(t *testing.T) string {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configContent := `llm:
  provider: openai
  openai:
    api_key: dummy
    model: gpt-4
    base_url: http://localhost`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = config.Init(configPath)
	if err != nil {
		t.Fatal(err)
	}
	return configPath
}

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

func TestChatHandlers_SendMessage_CommandStored(t *testing.T) {
	setupTestConfig(t)
	ctx := context.Background()

	store := &mockStore{}
	clipboard := &mockClipboard{}
	gateway := &mockGateway{
		completeFunc: func(ctx context.Context, req llm.CompletionRequest) (*llm.CompletionResponse, error) {
			return &llm.CompletionResponse{Content: "$ ls -la"}, nil
		},
	}

	handlers := &ChatHandlers{
		ctx:       ctx,
		gateway:   gateway,
		store:     store,
		clipboard: clipboard,
		filter:    security.DefaultFilter(),
	}

	store.addSessionFunc = func(sessionID, terminalID string) error {
		return nil
	}
	store.getSessionFunc = func(sessionID string) (*session.Session, error) {
		return nil, sql.ErrNoRows
	}
	var capturedCmd session.SuggestedCommand
	store.addCommandFunc = func(cmd session.SuggestedCommand) error {
		capturedCmd = cmd
		return nil
	}

	terminalID := "term-1"
	message := "list files"
	resp, err := handlers.SendMessage(terminalID, message)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	if resp != "$ ls -la" {
		t.Errorf("expected response '$ ls -la', got %q", resp)
	}
	if capturedCmd.Command != "$ ls -la" {
		t.Errorf("stored command mismatch: got %q", capturedCmd.Command)
	}
	if capturedCmd.SessionID != terminalID || capturedCmd.TerminalID != terminalID {
		t.Errorf("session/terminal IDs mismatch")
	}
}

func TestChatHandlers_CopyCommandToClipboard_Success(t *testing.T) {
	setupTestConfig(t)
	ctx := context.Background()

	store := &mockStore{}
	clipboard := &mockClipboard{}
	gateway := &mockGateway{}

	handlers := &ChatHandlers{
		ctx:       ctx,
		gateway:   gateway,
		store:     store,
		clipboard: clipboard,
		filter:    security.DefaultFilter(),
	}

	cmd := session.SuggestedCommand{
		ID:      "cmd-123",
		Command: "ls -la",
	}
	store.getCommandByIDFunc = func(commandID string) (*session.SuggestedCommand, error) {
		if commandID == "cmd-123" {
			return &cmd, nil
		}
		return nil, sql.ErrNoRows
	}
	var copiedText string
	clipboard.copyToTerminalFunc = func(ctx context.Context, text string, terminalID string) error {
		copiedText = text
		return nil
	}
	store.incrementUsedCountFunc = func(commandID string) error {
		if commandID != "cmd-123" {
			t.Errorf("unexpected command ID %s", commandID)
		}
		return nil
	}

	err := handlers.CopyCommandToClipboard("cmd-123", "term-1")
	if err != nil {
		t.Fatalf("CopyCommandToClipboard failed: %v", err)
	}
	if copiedText != "ls -la" {
		t.Errorf("expected copied text 'ls -la', got %q", copiedText)
	}
}

func TestChatHandlers_CopyCommandToClipboard_CommandNotFound(t *testing.T) {
	setupTestConfig(t)
	ctx := context.Background()
	store := &mockStore{}
	clipboard := &mockClipboard{}
	gateway := &mockGateway{}
	handlers := &ChatHandlers{
		ctx:       ctx,
		gateway:   gateway,
		store:     store,
		clipboard: clipboard,
		filter:    security.DefaultFilter(),
	}
	store.getCommandByIDFunc = func(commandID string) (*session.SuggestedCommand, error) {
		return nil, sql.ErrNoRows
	}
	err := handlers.CopyCommandToClipboard("cmd-unknown", "term-1")
	if err == nil {
		t.Error("expected error when command not found")
	}
	if !strings.Contains(err.Error(), "retrieve command") {
		t.Errorf("expected error about retrieve command, got %v", err)
	}
}

func TestCreateFilterFromConfig(t *testing.T) {
	ctx := testContext(t)

	// Test nil config -> default filter
	filter := createFilterFromConfig(ctx, nil)
	defaultFilter := security.DefaultFilter()
	text := "password = secret"
	if filter.Redact(text) != defaultFilter.Redact(text) {
		t.Errorf("nil config should return default filter")
	}

	// Test empty config -> default filter
	cfg := &config.Config{}
	filter = createFilterFromConfig(ctx, cfg)
	if filter.Redact(text) != defaultFilter.Redact(text) {
		t.Errorf("empty patterns should return default filter")
	}

	// Test custom pattern
	cfg.Security.FilterPatterns = []config.FilterPattern{
		{Pattern: "FOO_\\d+"},
	}
	filter = createFilterFromConfig(ctx, cfg)
	redacted := filter.Redact("secret FOO_123 bar")
	if redacted != "secret [REDACTED] bar" {
		t.Errorf("custom pattern not applied, got %q", redacted)
	}

	// Test invalid regex -> fallback with logging (cannot test logging easily)
	cfg.Security.FilterPatterns = []config.FilterPattern{
		{Pattern: "["}, // invalid regex
	}
	filter = createFilterFromConfig(ctx, cfg)
	if filter.Redact(text) != defaultFilter.Redact(text) {
		t.Errorf("invalid regex should fall back to default filter")
	}
}

func TestNewChatHandlers_CustomFilterPatterns(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configContent := `llm:
  provider: openai
  openai:
    api_key: dummy
    model: gpt-4
    base_url: http://localhost
security:
  filter_patterns:
    - name: "test"
      pattern: "SECRET_\\w+"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = config.Init(configPath)
	if err != nil {
		t.Fatal(err)
	}
	defer config.Init(configPath) // reset

	ctx := testContext(t)
	store, err := session.NewStore(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory store: %v", err)
	}
	defer store.Close()

	handlers, err := NewChatHandlers(ctx, store)
	if err != nil {
		t.Fatalf("NewChatHandlers failed: %v", err)
	}
	// Replace gateway with mock to capture request
	mockGW := &mockGateway{}
	gatewayField := reflect.ValueOf(handlers).Elem().FieldByName("gateway")
	if !gatewayField.IsValid() {
		t.Fatal("gateway field not found")
	}
	// Use unsafe to set unexported field
	ptr := unsafe.Pointer(gatewayField.UnsafeAddr())
	rf := reflect.NewAt(gatewayField.Type(), ptr).Elem()
	rf.Set(reflect.ValueOf(mockGW))
	// Call SendMessage with a secret pattern
	_, err = handlers.SendMessage("term-1", "hello SECRET_KEY world")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	// Check captured request
	if mockGW.capturedReq == nil {
		t.Fatal("expected captured request")
	}
	// The filtered message should have redacted the secret
	if len(mockGW.capturedReq.Messages) == 0 {
		t.Fatal("no messages captured")
	}
	content := mockGW.capturedReq.Messages[0].Content
	if content != "hello [REDACTED] world" {
		t.Errorf("expected redacted content, got %q", content)
	}
}
