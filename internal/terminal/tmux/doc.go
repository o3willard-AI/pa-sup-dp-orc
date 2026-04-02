// Package tmux provides a tmux terminal adapter for PairAdmin.
//
// This package implements the TerminalAdapter interface to enable
// PairAdmin to connect to and interact with tmux sessions.
//
// Usage:
//
//	config := tmux.DefaultConfig()
//	adapter := tmux.NewAdapter(config)
//
//	if err := adapter.Connect(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	defer adapter.Disconnect(ctx)
//
//	sessions, err := adapter.ListSessions(ctx)
//	content, err := adapter.Capture(ctx, "session-name")
package tmux
