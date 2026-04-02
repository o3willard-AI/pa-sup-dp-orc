package terminal

import "fmt"

// Error types for terminal operations

// ErrTerminalNotFound is returned when a requested terminal session is not found
type ErrTerminalNotFound struct {
	TerminalID string
}

func (e ErrTerminalNotFound) Error() string {
	return fmt.Sprintf("terminal session not found: %s", e.TerminalID)
}

// ErrAdapterNotAvailable is returned when a terminal adapter is not available on the system
type ErrAdapterNotAvailable struct {
	AdapterName string
	Reason      string
}

func (e ErrAdapterNotAvailable) Error() string {
	return fmt.Sprintf("terminal adapter '%s' is not available: %s", e.AdapterName, e.Reason)
}

// ErrCaptureFailed is returned when capturing terminal content fails
type ErrCaptureFailed struct {
	TerminalID string
	Reason     string
}

func (e ErrCaptureFailed) Error() string {
	return fmt.Sprintf("failed to capture terminal %s: %s", e.TerminalID, e.Reason)
}

// ErrSubscriptionFailed is returned when subscribing to terminal events fails
type ErrSubscriptionFailed struct {
	TerminalID string
	Reason     string
}

func (e ErrSubscriptionFailed) Error() string {
	return fmt.Sprintf("failed to subscribe to terminal %s: %s", e.TerminalID, e.Reason)
}

// ErrNoTerminalWindows is returned when no terminal windows are found
type ErrNoTerminalWindows struct{}

func (e ErrNoTerminalWindows) Error() string {
	return "no terminal windows found"
}
