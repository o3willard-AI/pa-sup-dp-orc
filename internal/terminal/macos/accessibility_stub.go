//go:build !darwin
// +build !darwin

// Package macos provides stub implementations for non-macOS platforms
package macos

import "errors"

var (
	ErrAccessibilityNotEnabled = errors.New("accessibility not available on this platform")
	ErrNoTerminalWindows       = errors.New("terminal windows not available on this platform")
	ErrTextExtractionFailed    = errors.New("text extraction not available on this platform")
)

func IsAccessibilityEnabled() bool {
	return false
}

func RequestAccessibilityPermission() bool {
	return false
}

func GetTerminalWindowTitles() ([]string, error) {
	return nil, ErrNoTerminalWindows
}

func GetTerminalWindowPids() ([]int, error) {
	return nil, ErrNoTerminalWindows
}

func ExtractTextFromWindow(pid int) (string, error) {
	return "", ErrTextExtractionFailed
}

func ExtractTextFromFrontmostTerminal() (string, error) {
	return "", ErrTextExtractionFailed
}
