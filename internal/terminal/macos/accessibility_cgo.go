//go:build darwin
// +build darwin

// Package macos provides CGO bindings for macOS Accessibility API
package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework ApplicationServices
#include "accessibility.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

var (
	// ErrAccessibilityNotEnabled is returned when accessibility permission is not granted
	ErrAccessibilityNotEnabled = errors.New("accessibility permission not enabled")

	// ErrNoTerminalWindows is returned when no Terminal.app windows are found
	ErrNoTerminalWindows = errors.New("no Terminal.app windows found")

	// ErrTextExtractionFailed is returned when text extraction fails
	ErrTextExtractionFailed = errors.New("failed to extract text from terminal")
)

// IsAccessibilityEnabled checks if the application has accessibility permission
func IsAccessibilityEnabled() bool {
	return C.accessibility_is_enabled() == 1
}

// RequestAccessibilityPermission requests accessibility permission from the user
// Returns true if permission was granted, false if denied
func RequestAccessibilityPermission() bool {
	return C.accessibility_request_permission() == 1
}

// GetTerminalWindowTitles returns the titles of all Terminal.app windows
func GetTerminalWindowTitles() ([]string, error) {
	var count C.int
	cArray := C.accessibility_get_terminal_window_titles(&count)

	if cArray == nil || count == 0 {
		return nil, ErrNoTerminalWindows
	}

	defer C.accessibility_free_string_array(cArray, count)

	titles := make([]string, 0, count)
	slice := unsafe.Slice(cArray, count)

	for i := 0; i < int(count); i++ {
		if slice[i] != nil {
			titles = append(titles, C.GoString(slice[i]))
		}
	}

	return titles, nil
}

// GetTerminalWindowPids returns the PIDs of Terminal.app instances
func GetTerminalWindowPids() ([]int, error) {
	var count C.int
	cArray := C.accessibility_get_terminal_window_pids(&count)

	if cArray == nil || count == 0 {
		return nil, ErrNoTerminalWindows
	}

	defer C.accessibility_free_int_array(cArray)

	slice := unsafe.Slice(cArray, count)
	pids := make([]int, count)

	for i := 0; i < int(count); i++ {
		pids[i] = int(slice[i])
	}

	return pids, nil
}

// ExtractTextFromWindow extracts text content from a Terminal.app window
func ExtractTextFromWindow(pid int) (string, error) {
	cStr := C.accessibility_extract_text_from_window(C.int(pid))

	if cStr == nil {
		return "", ErrTextExtractionFailed
	}

	defer C.accessibility_free_string(cStr)

	return C.GoString(cStr), nil
}

// ExtractTextFromFrontmostTerminal extracts text from the frontmost Terminal.app window
func ExtractTextFromFrontmostTerminal() (string, error) {
	cStr := C.accessibility_extract_text_from_frontmost_terminal()

	if cStr == nil {
		return "", ErrTextExtractionFailed
	}

	defer C.accessibility_free_string(cStr)

	return C.GoString(cStr), nil
}
