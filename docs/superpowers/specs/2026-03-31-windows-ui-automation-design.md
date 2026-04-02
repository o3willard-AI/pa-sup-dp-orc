# Windows UI Automation CGO Bindings Design

**Version:** 1.0  
**Date:** March 31, 2026  
**Status:** APPROVED  
**Author:** AI Development Team  
**Task:** 3.1 - Windows UI Automation CGO bindings (Phase 3)

---

## 1. Overview

This document specifies the design for Task 3.1: Windows UI Automation CGO bindings, enabling automatic terminal content capture on Windows via the UI Automation COM API. This is a full implementation (not stub) of the Windows terminal adapter for PairAdmin v2.0.

### 1.1 Scope

**In Scope:**
- UI Automation COM API integration for Windows Terminal, PowerShell, cmd.exe, and PuTTY
- CGO bindings for Go → COM communication
- Cross-platform compilation with proper build tags
- Comprehensive error handling and logging
- Unit testing with mocked C functions

**Out of Scope:**
- Screen capture + OCR fallback (deferred to future)
- Support for non-standard terminal emulators
- Alternative capture methods (Win32 API, etc.)

### 1.2 Design Decisions

| Decision | Rationale |
|----------|-----------|
| **Apartment-threaded COM model** | Standard for UI Automation, safer with STA objects |
| **All four terminal types** | Comprehensive support: Windows Terminal, PowerShell, cmd.exe, PuTTY |
| **Pure UI Automation approach** | Simpler than hybrid Win32 API, maintains single tech stack |
| **Mixed error handling** | Critical errors propagated to Go, non-critical logged internally |
| **Mocked C functions for unit tests** | Enables cross-platform testing without Windows SDK |

---

## 2. Architecture

### 2.1 Component Architecture

```
Go Application Layer
├── adapter_windows.go          # TerminalAdapter implementation
│   └── calls CGO functions
│
CGO Binding Layer
├── ui_automation.h            # C header (interface)
├── ui_automation.c            # C implementation (COM UI Automation)
│   └── singleton COM instance
│
COM UI Automation Layer
├── IUIAutomation              # Main UI Automation interface
├── ITextPattern              # Text extraction interface
└── Windows OS APIs           # Terminal window discovery
```

### 2.2 Build Configuration

```go
//go:build windows
// +build windows

package terminal

/*
#cgo windows LDFLAGS: -lole32 -lcombase
#include "ui_automation.h"
*/
import "C"
```

**Non-Windows Build:**
- `adapter_windows_stub.go` with `//go:build !windows`
- Returns appropriate errors for cross-platform compilation

---

## 3. Components

### 3.1 C Header (`ui_automation.h`)

**Updates to existing stub header:**

```c
// Error codes
typedef enum {
    UIA_OK = 0,
    UIA_INIT_FAILED = 1,
    UIA_COM_FAILURE = 2,
    UIA_WINDOW_NOT_FOUND = 3,
    UIA_TEXT_PATTERN_UNAVAILABLE = 4,
    UIA_MEMORY_ERROR = 5
} UIA_ErrorCode;

// Initialize UI Automation (returns error code)
int uia_initialize(void);

// Check availability
int uia_is_available(void);

// Cleanup resources
void uia_cleanup(void);

// Get terminal window titles with error reporting
wchar_t** get_terminal_window_titles(int* count, int* error_code, wchar_t** error_message);

// Get terminal window PIDs
int* get_terminal_window_pids(int* count, int* error_code, wchar_t** error_message);

// Extract text from specific window
wchar_t* extract_text_from_window(int pid, int* error_code, wchar_t** error_message);

// Extract text from focused terminal
wchar_t* extract_text_from_focused_terminal(int* error_code, wchar_t** error_message);

// Memory management
void free_wide_string_array(wchar_t** array, int count);
void free_wide_string(wchar_t* str);
```

### 3.2 C Implementation (`ui_automation.c`)

**Key Implementation Details:**

```c
// Global COM instance
static IUIAutomation* g_pAutomation = NULL;
static int g_initialized = 0;

int uia_initialize(void) {
    HRESULT hr = CoInitializeEx(NULL, COINIT_APARTMENTTHREADED);
    if (FAILED(hr)) return UIA_INIT_FAILED;
    
    hr = CoCreateInstance(&CLSID_CUIAutomation, NULL, CLSCTX_INPROC_SERVER,
                          &IID_IUIAutomation, (void**)&g_pAutomation);
    if (FAILED(hr)) {
        CoUninitialize();
        return UIA_COM_FAILURE;
    }
    
    g_initialized = 1;
    return UIA_OK;
}

// Window discovery using UI Automation conditions
wchar_t** get_terminal_window_titles(int* count, int* error_code, wchar_t** error_message) {
    // Create condition array for terminal window classes
    VARIANT classNames[4] = {
        {VT_BSTR, .bstrVal = L"CASCADIA_HOSTING_WINDOW_CLASS"},  // Windows Terminal
        {VT_BSTR, .bstrVal = L"ConsoleWindowClass"},             // PowerShell/cmd.exe
        {VT_BSTR, .bstrVal = L"PuTTY"},                          // PuTTY
        {VT_BSTR, .bstrVal = L"ConsoleWindowClass"}              // Additional for PowerShell
    };
    
    // Use FindAll with TreeScope_Children to locate terminal windows
    // Return array of window titles
}
```

### 3.3 Go Adapter (`adapter_windows.go`)

**Enhanced with proper error handling:**

```go
// Adapter implements TerminalAdapter for Windows terminals
type Adapter struct {
    config     Config
    mu         sync.RWMutex
    running    bool
    cancelFunc context.CancelFunc
}

// Name returns "windows"
func (a *Adapter) Name() string { return "windows" }

// Available checks UI Automation availability
func (a *Adapter) Available(ctx context.Context) bool {
    return C.uia_is_available() != 0
}

// ListSessions returns detected terminal windows with error handling
func (a *Adapter) ListSessions(ctx context.Context) ([]DetectedTerminal, error) {
    var cCount C.int
    var cError C.int
    var cErrorMessage *C.wchar_t
    defer func() {
        if cErrorMessage != nil {
            C.free_wide_string(cErrorMessage)
        }
    }()
    
    cTitles := C.get_terminal_window_titles(&cCount, &cError, &cErrorMessage)
    if cError != 0 {
        errorMsg := "unknown error"
        if cErrorMessage != nil {
            errorMsg = C.GoString((*C.char)(unsafe.Pointer(cErrorMessage)))
        }
        return nil, fmt.Errorf("UI Automation error %d: %s", int(cError), errorMsg)
    }
    if cTitles == nil {
        return nil, nil
    }
    defer C.free_wide_string_array(cTitles, cCount)
    
    // Convert C array to Go slice
    // Format terminal IDs: windows-{pid}-{window-title}
}
```

### 3.4 Go Stub (`adapter_windows_stub.go`)

```go
//go:build !windows
// +build !windows

package terminal

// Stub returns appropriate errors for non-Windows platforms
func (a *Adapter) Available(ctx context.Context) bool {
    return false
}

func (a *Adapter) ListSessions(ctx context.Context) ([]DetectedTerminal, error) {
    return nil, errors.New("Windows adapter not available on this platform")
}
```

---

## 4. Data Flow

### 4.1 Terminal Detection Flow

1. **Initialize COM:** `CoInitializeEx(NULL, COINIT_APARTMENTTHREADED)`
2. **Create UI Automation:** `CoCreateInstance(CLSID_CUIAutomation)`
3. **Build Condition Array:** Terminal window class names
4. **Find Windows:** `IUIAutomation::FindAll(TreeScope_Children, condition)`
5. **Extract Properties:** `GetCurrentName()`, `GetCurrentProcessId()`
6. **Format Results:** Convert to `DetectedTerminal` structs

### 4.2 Text Extraction Flow

1. **Locate Window:** Find `IUIAutomationElement` by PID
2. **Query Text Pattern:** `GetCurrentPattern(UIA_TextPatternId)`
3. **Extract Text:** `ITextPattern::GetSelection()` or `GetText()`
4. **Convert Encoding:** Wide char to UTF-8 for Go
5. **Error Handling:** Check HRESULT, propagate or log

### 4.3 Terminal ID Format

```
windows-{pid}-{window-title}
Examples:
- windows-12345-Administrator:Windows PowerShell
- windows-67890-C:\Windows\System32\cmd.exe
- windows-54321-user@host:~
```

---

## 5. Error Handling

### 5.1 Error Classification

| Error Type | Handling | Examples |
|------------|----------|----------|
| **Critical** | Propagate to Go | COM initialization failure, memory allocation errors |
| **Non-critical** | Log internally, return empty | Window not found, text pattern unavailable |
| **Recoverable** | Retry or use defaults | Temporary COM failures, focus changes |

**Note:** Non-critical errors are logged internally using Go's standard `log` package with appropriate log levels. Critical errors propagate to Go callers with descriptive error messages.

### 5.2 Error Propagation

```c
// C functions return error codes with optional messages
wchar_t* extract_text_from_window(int pid, int* error_code, wchar_t** error_message) {
    if (!g_initialized) {
        *error_code = UIA_INIT_FAILED;
        *error_message = L"UI Automation not initialized";
        return NULL;
    }
    
    // Implementation with HRESULT checking
    HRESULT hr = S_OK; // Actual COM call implementation would go here
    if (FAILED(hr)) {
        *error_code = UIA_COM_FAILURE;
        *error_message = L"COM operation failed";
        return NULL;
    }
    
    *error_code = UIA_OK;
    return extracted_text;
}
```

### 5.3 Go Error Wrapping

```go
func ExtractTextFromWindow(pid int) (string, error) {
    var cError C.int
    var cErrorMessage *C.wchar_t
    defer func() {
        if cErrorMessage != nil {
            C.free_wide_string(cErrorMessage)
        }
    }()
    
    cText := C.extract_text_from_window(C.int(pid), &cError, &cErrorMessage)
    if cError != 0 {
        errorMsg := "unknown error"
        if cErrorMessage != nil {
            errorMsg = C.GoString((*C.char)(unsafe.Pointer(cErrorMessage)))
        }
        return "", fmt.Errorf("UI Automation error %d: %s", int(cError), errorMsg)
    }
    
    defer C.free_wide_string(cText)
    return C.GoString((*C.char)(unsafe.Pointer(cText))), nil
}
```

---

## 6. Testing Strategy

### 6.1 Unit Tests (Cross-Platform)

**Mocked C Functions:**

```go
//go:build windows
// +build windows

package terminal

/*
// Mock implementations for testing
int mock_uia_initialize(void) { return 0; }
int mock_uia_is_available(void) { return 1; }

wchar_t* mock_extract_text_from_window(int pid, int* error_code, wchar_t** error_msg) {
    static wchar_t text[] = L"mock terminal content\n$ echo hello\nhello";
    *error_code = 0;
    return _wcsdup(text);
}
*/
import "C"

func TestWindowsAdapter_Mocked(t *testing.T) {
    adapter := NewAdapter(DefaultConfig())
    
    // Test Available()
    if !adapter.Available(context.Background()) {
        t.Error("Expected adapter to be available")
    }
    
    // Test ListSessions()
    sessions, err := adapter.ListSessions(context.Background())
    if err != nil {
        t.Errorf("ListSessions failed: %v", err)
    }
    
    // Test Capture()
    if len(sessions) > 0 {
        content, err := adapter.Capture(context.Background(), sessions[0].ID)
        if err != nil {
            t.Errorf("Capture failed: %v", err)
        }
        if !strings.Contains(content, "mock terminal") {
            t.Error("Capture returned unexpected content")
        }
    }
}
```

**Mocking Approach:** Mock C functions are defined in test files with the same signatures as real functions. During unit testing, the test file's C declarations override the real implementations via compile-time substitution (test build tags ensure mock functions are used).

### 6.2 Integration Tests (Windows-Only)

**Prerequisites:**
- Windows Terminal or PowerShell running with content
- UI Automation enabled (default on Windows)
- Build with CGO_ENABLED=1

**Test Cases:**
1. **Basic Detection:** Verify terminal windows are detected
2. **Text Extraction:** Capture content matches terminal display
3. **Error Handling:** Test error paths (e.g., invalid PID)
4. **Concurrent Access:** Multiple goroutines using same adapter
5. **Cleanup:** Verify COM resources are released

### 6.3 Cross-Platform Build Verification

```bash
# Linux: should use stub
GOOS=linux go build ./internal/terminal/windows/...

# Windows: should build with CGO
GOOS=windows CGO_ENABLED=1 go build ./internal/terminal/windows/...

# Test with race detector
go test -race ./internal/terminal/windows/...
```

---

## 7. Dependencies

### 7.1 System Requirements

**Windows SDK:**
- UI Automation headers (`uiautomation.h`, `uiautomationclient.h`)
- COM libraries (`ole32.lib`, `combase.lib`)
- Minimum: Windows 7 (UI Automation introduced in Windows 7)

**Build Tools:**
- MSVC or MinGW toolchain for CGO
- `pkg-config` not required (Windows uses direct linking)

### 7.2 Go Dependencies

```go
// Internal dependencies
import (
    "context"
    "fmt"
    "sync"
    "time"
    "unsafe"
)

// No external Go dependencies for CGO bindings
```

### 7.3 Development Setup

```bash
# Install Windows SDK (if not already present)
# Typically included with Visual Studio Build Tools

# Verify COM headers exist
ls "C:\Program Files (x86)\Windows Kits\10\Include\*\um\uiautomation.h"

# Build with CGO
set CGO_ENABLED=1
go build ./internal/terminal/windows/...
```

---

## 8. Security Considerations

### 8.1 Permission Requirements

**UI Automation:**
- Requires no special permissions (standard user access)
- Uses COM security with default impersonation level
- May require UAC elevation for elevated terminals

**Privacy:**
- Only captures terminal windows (not other applications)
- No screen capture or keystroke logging
- Text extraction limited to terminal content

### 8.2 COM Security

**Threading Model:** `COINIT_APARTMENTTHREADED` ensures UI Automation runs in STA
**Error Handling:** Proper `HRESULT` checking prevents security issues
**Resource Cleanup:** `CoUninitialize()` releases COM resources

### 8.3 Memory Safety

**C Memory:** `CoTaskMemAlloc()` for COM strings, `free()` for C allocations
**Go Finalizers:** Use `defer` for C string cleanup
**Buffer Overflow:** Safe string conversion with length checking

---

## 9. References

### 9.1 Microsoft Documentation

- [UI Automation Overview](https://docs.microsoft.com/en-us/windows/win32/winauto/uiauto-overview)
- [IUIAutomation Interface](https://docs.microsoft.com/en-us/windows/win32/api/uiautomationcore/nn-uiautomationcore-iuiautomation)
- [ITextPattern Interface](https://docs.microsoft.com/en-us/windows/win32/api/uiautomationclient/nn-uiautomationclient-itextpattern)
- [COM Programming Guide](https://docs.microsoft.com/en-us/windows/win32/com/com-programming-guide)

### 9.2 Go CGO References

- [CGO with Windows](https://github.com/golang/go/wiki/cgo#windows)
- [Go and COM](https://github.com/go-ole/go-ole) (reference implementation)
- [Windows Syscall Package](https://pkg.go.dev/golang.org/x/sys/windows)

### 9.3 Related PairAdmin Components

- **Task 2.5:** macOS Accessibility Adapter (similar CGO pattern)
- **Task 2.8:** Linux AT-SPI2 Adapter (similar platform-specific approach)
- **Task 2.9:** Windows Adapter Stub (existing stub to be replaced)
- **Task 3.2:** Windows Adapter Implementation (next task, builds on these bindings)

---

## 10. Next Steps

### 10.1 Implementation Plan

1. **Update C Header:** Add error codes and message parameters
2. **Implement C Functions:** COM UI Automation integration
3. **Enhance Go Adapter:** Add error handling and session management
4. **Write Unit Tests:** Mocked C functions for cross-platform testing
5. **Integration Testing:** Windows-only tests with real terminals
6. **Documentation:** Update permissions guide for Windows

### 10.2 Success Criteria

✅ CGO bindings compile on Windows with COM libraries  
✅ Terminal detection finds Windows Terminal/PowerShell/cmd.exe/PuTTY  
✅ Text extraction returns accurate terminal content  
✅ Error handling propagates critical errors to Go  
✅ Unit tests pass on all platforms (using mocks)  
✅ Integration tests pass on Windows (with real terminals)  
✅ Cross-compilation works (stub on non-Windows)  

### 10.3 Exit Criteria

**Phase 3.1 Complete When:**
- Windows UI Automation CGO bindings are fully implemented
- All acceptance criteria met
- Tests pass with >80% coverage
- Documentation updated
- Ready for integration with Task 3.2 (Windows adapter implementation)

---

**Design Approved:** March 31, 2026  
**Implementation Start:** Immediate  
**Expected Completion:** 2-3 development cycles