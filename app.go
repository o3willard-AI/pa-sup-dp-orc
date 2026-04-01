package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pairadmin/pairadmin/internal/config"
	"github.com/pairadmin/pairadmin/internal/session"
	"github.com/pairadmin/pairadmin/internal/ui"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SendMessageResponse is the response from SendMessage.
type SendMessageResponse struct {
	Content   string `json:"content"`
	CommandID string `json:"commandID"`
}

// App struct
type App struct {
	ctx              context.Context
	configPath       string
	sessionStore     *session.Store
	terminalHandlers *ui.TerminalHandlers
	chatHandlers     *ui.ChatHandlers
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize configuration
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	a.configPath = filepath.Join(configDir, "pairadmin", "config.yaml")
	if err := config.Init(a.configPath); err != nil {
		runtime.LogError(ctx, fmt.Sprintf("failed to init config: %v", err))
		panic(fmt.Sprintf("failed to initialize configuration: %v", err))
	}

	// Initialize session store
	dbPath := filepath.Join(filepath.Dir(a.configPath), "sessions.db")
	store, err := session.NewStore(dbPath)
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("failed to init session store: %v", err))
		panic(fmt.Sprintf("failed to initialize session store: %v", err))
	}
	a.sessionStore = store

	// Initialize terminal detector and handlers
	// TODO: create detector with auto‑registered adapters
	// For now, create a nil detector (will be added in later tasks)
	a.terminalHandlers = ui.NewTerminalHandlers(nil)

	// Initialize chat handlers
	chatHandlers, err := ui.NewChatHandlers(ctx, store)
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("failed to init chat handlers: %v", err))
		panic(fmt.Sprintf("failed to initialize chat handlers: %v", err))
	}
	a.chatHandlers = chatHandlers
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SendMessage delegates to chat handlers.
func (a *App) SendMessage(terminalID, message string) (SendMessageResponse, error) {
	if a.chatHandlers == nil {
		return SendMessageResponse{}, fmt.Errorf("chat handlers not initialized")
	}
	content, commandID, err := a.chatHandlers.SendMessage(terminalID, message)
	if err != nil {
		return SendMessageResponse{}, err
	}
	return SendMessageResponse{Content: content, CommandID: commandID}, nil
}

// GetCommandsByTerminal returns all commands for a terminal.
func (a *App) GetCommandsByTerminal(terminalID string) ([]session.SuggestedCommand, error) {
	if a.sessionStore == nil {
		return nil, fmt.Errorf("session store not initialized")
	}
	return a.sessionStore.GetCommandsByTerminal(terminalID)
}

// CopyCommandToClipboard delegates to chat handlers.
func (a *App) CopyCommandToClipboard(commandID, terminalID string) error {
	if a.chatHandlers == nil {
		return fmt.Errorf("chat handlers not initialized")
	}
	return a.chatHandlers.CopyCommandToClipboard(commandID, terminalID)
}
