//go:build windows && mock
// +build windows,mock

package terminal

/*
#include <stdlib.h>
#include <wchar.h>

// Mock error codes matching ui_automation.h (see ui_automation.h for authoritative definition)
typedef enum {
    UIA_OK = 0,
    UIA_INIT_FAILED = 1,
    UIA_COM_FAILURE = 2,
    UIA_WINDOW_NOT_FOUND = 3,
    UIA_TEXT_PATTERN_UNAVAILABLE = 4,
    UIA_MEMORY_ERROR = 5,
    UIA_INVALID_ARGUMENT = 6,
    UIA_UNKNOWN_ERROR = 7
} UIA_ErrorCode;

// Mock implementations for testing
static int mock_initialized = 0;
static UIA_ErrorCode mock_error = UIA_OK;

// Set mock error for testing (call from Go tests)
void set_mock_error(UIA_ErrorCode err) {
    mock_error = err;
}

// Reset mock error to UIA_OK
void reset_mock_error(void) {
    mock_error = UIA_OK;
}

UIA_ErrorCode uia_initialize(void) {
    if (mock_error != UIA_OK) {
        return mock_error;
    }
    mock_initialized = 1;
    return UIA_OK;
}

int uia_is_available(void) {
    return mock_initialized;
}

void uia_cleanup(void) {
    mock_initialized = 0;
}

wchar_t** get_terminal_window_titles(int* count, UIA_ErrorCode* error_code, wchar_t** error_message) {
    if (mock_error != UIA_OK) {
        *error_code = mock_error;
        *error_message = (wchar_t*)malloc((wcslen(L"Mock error injected") + 1) * sizeof(wchar_t));
        if (*error_message) wcscpy(*error_message, L"Mock error injected");
        *count = 0;
        return NULL;
    }

    *error_code = UIA_OK;
    *error_message = NULL;
    *count = 2;

    wchar_t** titles = (wchar_t**)malloc(2 * sizeof(wchar_t*));
    if (!titles) return NULL;

    // Duplicate strings (caller must free with free_wide_string)
    titles[0] = (wchar_t*)malloc((wcslen(L"Windows PowerShell") + 1) * sizeof(wchar_t));
    if (!titles[0]) {
        free(titles);
        return NULL;
    }
    wcscpy(titles[0], L"Windows PowerShell");

    titles[1] = (wchar_t*)malloc((wcslen(L"Command Prompt") + 1) * sizeof(wchar_t));
    if (!titles[1]) {
        free(titles[0]);
        free(titles);
        return NULL;
    }
    wcscpy(titles[1], L"Command Prompt");

    return titles;
}

int* get_terminal_window_pids(int* count, UIA_ErrorCode* error_code, wchar_t** error_message) {
    if (mock_error != UIA_OK) {
        *error_code = mock_error;
        *error_message = (wchar_t*)malloc((wcslen(L"Mock error injected") + 1) * sizeof(wchar_t));
        if (*error_message) wcscpy(*error_message, L"Mock error injected");
        *count = 0;
        return NULL;
    }

    *error_code = UIA_OK;
    *error_message = NULL;
    *count = 2;

    int* pids = (int*)malloc(2 * sizeof(int));
    if (!pids) return NULL;

    pids[0] = 12345;
    pids[1] = 67890;

    return pids;
}

wchar_t* extract_text_from_window(int pid, UIA_ErrorCode* error_code, wchar_t** error_message) {
    if (mock_error != UIA_OK) {
        *error_code = mock_error;
        *error_message = (wchar_t*)malloc((wcslen(L"Mock error injected") + 1) * sizeof(wchar_t));
        if (*error_message) wcscpy(*error_message, L"Mock error injected");
        return NULL;
    }

    *error_code = UIA_OK;
    *error_message = NULL;

    if (pid == 12345) {
        const wchar_t* content = L"PS C:\\Users\\test> echo hello\nhello\nPS C:\\Users\\test>";
        wchar_t* copy = (wchar_t*)malloc((wcslen(content) + 1) * sizeof(wchar_t));
        if (copy) wcscpy(copy, content);
        return copy;
    } else if (pid == 67890) {
        const wchar_t* content = L"C:\\Users\\test> dir\n Volume in drive C is OS\n Directory of C:\\Users\\test";
        wchar_t* copy = (wchar_t*)malloc((wcslen(content) + 1) * sizeof(wchar_t));
        if (copy) wcscpy(copy, content);
        return copy;
    }

    const wchar_t* content = L"mock terminal content";
    wchar_t* copy = (wchar_t*)malloc((wcslen(content) + 1) * sizeof(wchar_t));
    if (copy) wcscpy(copy, content);
    return copy;
}

wchar_t* extract_text_from_focused_terminal(UIA_ErrorCode* error_code, wchar_t** error_message) {
    if (mock_error != UIA_OK) {
        *error_code = mock_error;
        *error_message = (wchar_t*)malloc((wcslen(L"Mock error injected") + 1) * sizeof(wchar_t));
        if (*error_message) wcscpy(*error_message, L"Mock error injected");
        return NULL;
    }

    *error_code = UIA_OK;
    *error_message = NULL;
    const wchar_t* content = L"Focused terminal content\n$ pwd\n/home/user";
    wchar_t* copy = (wchar_t*)malloc((wcslen(content) + 1) * sizeof(wchar_t));
    if (copy) wcscpy(copy, content);
    return copy;
}

void free_wide_string_array(wchar_t** array, int count) {
    if (!array) return;
    for (int i = 0; i < count; i++) {
        if (array[i]) free(array[i]);
    }
    free(array);
}

void free_wide_string(wchar_t* str) {
    if (str) free(str);
}

void free_int_array(int* array) {
    if (array) free(array);
}
*/
import "C"
import (
	"context"
	"strings"
	"testing"
)

func TestWindowsAdapter_Mocked(t *testing.T) {
	t.Cleanup(func() { C.reset_mock_error() })

	t.Run("HappyPath", func(t *testing.T) {
		adapter := NewAdapter(DefaultConfig())

		// Test Available() with mocked initialization
		if !adapter.Available(context.Background()) {
			t.Error("Expected adapter to be available with mocked C functions")
		}

		// Test ListSessions()
		sessions, err := adapter.ListSessions(context.Background())
		if err != nil {
			t.Errorf("ListSessions failed: %v", err)
		}
		if len(sessions) != 2 {
			t.Errorf("Expected 2 sessions, got %d", len(sessions))
		}

		// Verify session IDs
		for i, session := range sessions {
			if !strings.HasPrefix(session.ID, "windows-") {
				t.Errorf("Session ID doesn't start with 'windows-': %s", session.ID)
			}
			if session.Adapter != "windows" {
				t.Errorf("Expected adapter 'windows', got %s", session.Adapter)
			}
			if i == 0 && session.Name != "Windows PowerShell" {
				t.Errorf("Expected name 'Windows PowerShell', got %s", session.Name)
			}
		}

		// Test Capture()
		if len(sessions) > 0 {
			content, err := adapter.Capture(context.Background(), sessions[0].ID)
			if err != nil {
				t.Errorf("Capture failed: %v", err)
			}
			if !strings.Contains(content, "echo hello") {
				t.Error("Capture returned unexpected content")
			}
		}

		// Test Stop()
		adapter.Stop()
	})

	t.Run("ErrorInjection", func(t *testing.T) {
		// Test initialization error
		C.set_mock_error(C.UIA_INIT_FAILED)
		defer C.reset_mock_error()

		adapter := NewAdapter(DefaultConfig())
		if adapter.Available(context.Background()) {
			t.Error("Expected adapter to be unavailable with mock init error")
		}

		// Test ListSessions error
		C.set_mock_error(C.UIA_COM_FAILURE)
		sessions, err := adapter.ListSessions(context.Background())
		if err == nil {
			t.Error("Expected ListSessions to fail with mock error")
		}
		if sessions != nil {
			t.Error("Expected nil sessions on error")
		}

		// Test Capture error
		C.set_mock_error(C.UIA_WINDOW_NOT_FOUND)
		content, err := adapter.Capture(context.Background(), "windows-12345")
		if err == nil {
			t.Error("Expected Capture to fail with mock error")
		}
		if content != "" {
			t.Error("Expected empty content on error")
		}

		C.reset_mock_error()
	})

	t.Run("GetDimensions", func(t *testing.T) {
		adapter := NewAdapter(DefaultConfig())
		rows, cols, err := adapter.GetDimensions(context.Background(), "windows-12345")
		if err != nil {
			t.Errorf("GetDimensions failed: %v", err)
		}
		if rows != 24 || cols != 80 {
			t.Errorf("Expected dimensions 24x80, got %dx%d", rows, cols)
		}
	})

	t.Run("SubscribeNoPanic", func(t *testing.T) {
		adapter := NewAdapter(DefaultConfig())
		events, err := adapter.Subscribe(context.Background(), "windows-12345")
		if err != nil {
			t.Errorf("Subscribe failed: %v", err)
		}
		// Close the adapter to stop the goroutine
		adapter.Stop()
		// Drain any pending events (should be closed)
		for range events {
		}
	})
}
