# Product Requirements Document
## PairAdmin v2.0
### Cross-Platform AI-Assisted Terminal Administration

**Version:** 2.0  
**Date:** February 2026  
**Status:** CONFIDENTIAL  
**Revision Note:** Complete architectural pivot from embedded PuTTY extension to standalone cross-platform application

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [User Stories and Acceptance Criteria](#2-user-stories-and-acceptance-criteria)
3. [UX/UI Wireframe Descriptions](#3-uxui-wireframe-descriptions)
4. [Technical Architecture](#4-technical-architecture)
5. [Terminal Integration Specifications](#5-terminal-integration-specifications)
6. [Security Constraints and Considerations](#6-security-constraints-and-considerations)
7. [Slash Commands Reference](#7-slash-commands-reference)
8. [Implementation Roadmap](#8-implementation-roadmap)
9. [Appendix](#9-appendix)

---

## 1. Executive Summary

### 1.1 Product Vision

PairAdmin enables **"Pair Administration"**—a collaboration model where human system administrators work alongside AI agents to manage Linux, Unix, macOS, and Windows systems via terminal interfaces. The AI observes terminal sessions in real-time and provides contextual command suggestions that users can execute with a single hotkey or click.

### 1.2 Pivot Rationale

The original v1.0 design embedded AI functionality directly into PuTTY through source modifications. After extensive prototyping, this approach proved untenable due to:

- Complex UI integration failures between PuTTY's Win32 rendering and modern UI frameworks
- Unexpected terminal I/O edge cases creating unreliable capture
- Single-platform focus (Windows/PuTTY only) limiting market reach
- Maintenance burden of tracking upstream PuTTY changes

**v2.0 adopts a fundamentally different architecture:** a standalone cross-platform application that integrates with existing terminals through well-defined system APIs rather than embedding within them.

### 1.3 Solution Overview

PairAdmin v2.0 is a **standalone desktop application** written in Go that:

- **Reads terminal buffers** automatically from supported terminals via platform-native APIs
- **Presents an AI chat interface** with full context of the observed terminal session
- **Provides one-click command execution** via clipboard injection into the target terminal
- **Supports multiple concurrent terminals** through a tabbed interface with isolated context per session
- **Works cross-platform** on Linux, macOS, and Windows

### 1.4 Supported Terminal Integrations (v1.0 Scope)

| Platform | Terminals |
|----------|-----------|
| **Linux/Unix** | tmux, GNOME Terminal, Konsole |
| **macOS** | tmux, Terminal.app |
| **Windows** | PuTTY |

---

## 2. User Stories and Acceptance Criteria

### 2.1 User Personas

*(Unchanged from v1.0)*

#### Persona 1: Junior System Administrator
| Attribute | Details |
|-----------|---------|
| **Name** | Alex Chen |
| **Experience** | 1-2 years in IT operations |
| **Goals** | Learn best practices, reduce errors, build confidence |
| **Pain Points** | Unfamiliar commands, hesitant without verification, documentation-dependent |

#### Persona 2: Senior DevOps Engineer
| Attribute | Details |
|-----------|---------|
| **Name** | Maria Santos |
| **Experience** | 8+ years managing production infrastructure |
| **Goals** | Maximize efficiency, automate tasks, maintain audit trails |
| **Pain Points** | Routine task overhead, context switching, documentation burden |

### 2.2 Core User Stories

#### US-001: Automatic Terminal Context Capture

> **As a** system administrator, **I want** PairAdmin to automatically capture and display the contents of my active terminal session, **so that** I don't have to manually copy/paste terminal output to get AI assistance.

**Acceptance Criteria:**

- [ ] Given PairAdmin is running and a supported terminal is active, when the terminal content changes, then PairAdmin captures the new content within 500ms
- [ ] Given a terminal session is captured, when the user views the PairAdmin tab for that terminal, then they see a synchronized representation of recent terminal output
- [ ] Given multiple terminals are open, when the user switches between PairAdmin tabs, then each tab shows context isolated to its associated terminal
- [ ] Given a terminal buffer exceeds the context limit, when captured, then the most recent content is prioritized

---

#### US-002: Contextual AI Assistance

> **As a** system administrator, **I want** to ask the AI questions about what I see in my terminal, **so that** I can understand errors, plan next steps, and get command suggestions without leaving my workflow.

**Acceptance Criteria:**

- [ ] Given a terminal tab is active in PairAdmin, when the user types a question, then the AI response incorporates the current terminal context
- [ ] Given the terminal shows an error message, when the user asks "what went wrong?", then the AI explains the error with reference to the specific output
- [ ] Given the AI suggests a command, when displayed in chat, then it appears in a code block with a prominent "Copy to Terminal" button
- [ ] Given the user clicks "Copy to Terminal", then the command is placed in the system clipboard AND optionally auto-pasted to the target terminal

---

#### US-003: Command History Sidebar

> **As a** power user, **I want** quick access to all commands the AI has suggested during my session, **so that** I can rapidly re-execute or modify previous suggestions.

**Acceptance Criteria:**

- [ ] Given the AI has suggested commands, when the user views the sidebar, then all suggested commands appear in reverse-chronological order
- [ ] Given a command in the sidebar, when the user clicks it, then it is copied to clipboard and optionally pasted to the terminal
- [ ] Given a command in the sidebar, when the user hovers over it, then they see the context/question that generated it
- [ ] Given the sidebar contains many commands, when the user scrolls, then smooth scrolling maintains position and performance

---

#### US-004: Multi-Terminal Tab Management

> **As a** DevOps engineer managing multiple systems, **I want** to maintain separate AI conversations for each terminal session, **so that** context from one server doesn't confuse assistance for another.

**Acceptance Criteria:**

- [ ] Given multiple terminals are detected, when the user views PairAdmin, then each terminal has its own tab
- [ ] Given a terminal tab, when selected, then the chat history and command sidebar show only that terminal's context
- [ ] Given a new terminal session starts, when detected, then a new tab is automatically created
- [ ] Given a terminal session ends, when detected, then the user is prompted to close or archive the tab

---

#### US-005: One-Click Command Execution

> **As a** system administrator, **I want** to execute AI-suggested commands with a single click or hotkey, **so that** I can work rapidly without manual copy/paste operations.

**Acceptance Criteria:**

- [ ] Given an AI-suggested command, when the user clicks "Copy to Terminal", then the command is injected into the clipboard
- [ ] Given clipboard injection, when the target terminal is focused, then the user can paste with standard Ctrl+V / Cmd+V
- [ ] Given a global hotkey is configured, when pressed, then the most recent suggested command is copied and the target terminal is focused
- [ ] Given command execution, when complete, then PairAdmin captures the resulting terminal output for continued conversation

---

#### US-006: LLM Configuration

> **As a** user, **I want** to configure which AI model PairAdmin uses and customize the system prompt, **so that** I can use my preferred provider and tailor the AI's behavior to my needs.

**Acceptance Criteria:**

- [ ] Given the settings menu, when opened, then the user can select from available LLM providers (OpenAI, Anthropic, Ollama, etc.)
- [ ] Given a provider is selected, when API credentials are entered, then connection status is validated and displayed
- [ ] Given the system prompt settings, when the user adds custom instructions, then these are appended to the built-in agent prompt
- [ ] Given custom prompt extensions, when the user chats, then the AI behavior reflects the custom instructions

---

## 3. UX/UI Wireframe Descriptions

### 3.1 Application Window Layout

PairAdmin uses a **three-column layout** optimized for terminal administration workflows:

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  PairAdmin                                                    [─] [□] [×]       │
├─────────────────────────────────────────────────────────────────────────────────┤
│  [File]  [Settings]  [Help]                                                     │
├────────────────┬────────────────────────────────────────┬───────────────────────┤
│                │                                        │                       │
│   TERMINAL     │              CHAT AREA                 │      COMMAND          │
│    TABS        │                                        │      SIDEBAR          │
│                │  ┌──────────────────────────────────┐  │                       │
│  ┌──────────┐  │  │ 🤖 AI:                           │  │  ┌─────────────────┐  │
│  │ prod-web │◄─┤  │ I see a "Permission denied"     │  │  │ sudo systemctl  │  │
│  └──────────┘  │  │ error on line 47. This occurs   │  │  │ restart nginx   │  │
│  ┌──────────┐  │  │ because the nginx user cannot   │  │  │                 │  │
│  │ prod-db  │  │  │ write to /var/log/custom.log.   │  │  │ 2 min ago       │  │
│  └──────────┘  │  │                                  │  │  └─────────────────┘  │
│  ┌──────────┐  │  │ Suggested fix:                   │  │  ┌─────────────────┐  │
│  │ staging  │  │  │ ┌────────────────────────────┐   │  │  │ chown nginx:    │  │
│  └──────────┘  │  │ │ sudo chown nginx:nginx     │   │  │  │ nginx /var/log/ │  │
│  ┌──────────┐  │  │ │ /var/log/custom.log        │   │  │  │ custom.log      │  │
│  │ + New    │  │  │ │            [Copy to Term ▶]│   │  │  │                 │  │
│  └──────────┘  │  │ └────────────────────────────┘   │  │  │ 5 min ago       │  │
│                │  └──────────────────────────────────┘  │  └─────────────────┘  │
│                │                                        │  ┌─────────────────┐  │
│  ───────────   │  ┌──────────────────────────────────┐  │  │ tail -f /var/   │  │
│  TERMINAL      │  │ 👤 You:                          │  │  │ log/nginx/      │  │
│  PREVIEW       │  │ What's causing this nginx error? │  │  │ error.log       │  │
│  ───────────   │  └──────────────────────────────────┘  │  │                 │  │
│                │                                        │  │ 8 min ago       │  │
│  user@prod:~$  │                                        │  └─────────────────┘  │
│  nginx: error  │                                        │                       │
│  Permission    │  ┌──────────────────────────────────┐  │  ···                  │
│  denied...     │  │ Ask about this terminal...   [➤] │  │                       │
│  [see full ▼]  │  └──────────────────────────────────┘  │  [Clear History]      │
│                │                                        │                       │
├────────────────┴────────────────────────────────────────┴───────────────────────┤
│  [GPT-4 ▼]  │  Connected: prod-web (tmux)  │  Context: 12.4K/32K  │  [⚙ Settings] │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Layout Components

#### 3.2.1 Terminal Tabs (Left Column)

| Element | Specification |
|---------|---------------|
| **Width** | Fixed 160px, non-resizable |
| **Tab Content** | Terminal identifier (hostname, session name, or user-assigned label) |
| **Active Indicator** | Left border highlight + background color change |
| **New Tab Button** | Opens terminal detection/selection dialog |
| **Terminal Preview** | Collapsible panel showing last ~10 lines of terminal output |
| **Preview Toggle** | "See full" expands to scrollable view of captured buffer |

#### 3.2.2 Chat Area (Center Column)

| Element | Specification |
|---------|---------------|
| **Width** | Flexible, fills available space (min 400px) |
| **Message Layout** | User messages right-aligned (light blue bg), AI messages left-aligned (neutral bg) |
| **Code Blocks** | Monospace font, syntax highlighting, "Copy to Terminal" button |
| **Copy to Terminal Button** | Primary action button, prominent color (green), includes paste icon |
| **Input Field** | Fixed at bottom, auto-expanding (1-4 lines), placeholder: "Ask about this terminal..." |
| **Send Button** | Arrow icon, Enter to send, Shift+Enter for newline |

#### 3.2.3 Command Sidebar (Right Column)

| Element | Specification |
|---------|---------------|
| **Width** | Fixed 220px, collapsible |
| **Command Cards** | Each suggested command in a card with: command text (truncated), timestamp, click-to-copy |
| **Hover Detail** | Tooltip shows full command + original question context |
| **Ordering** | Reverse chronological (newest at top) |
| **Clear Button** | Bottom of sidebar, clears command history for current tab |
| **Scroll Behavior** | Smooth scroll, maintains position when new commands added |

#### 3.2.4 Status Bar (Bottom)

| Element | Position | Content |
|---------|----------|---------|
| **Model Selector** | Left | Dropdown: current LLM provider/model |
| **Connection Status** | Center-left | Terminal type + identifier (e.g., "Connected: prod-web (tmux)") |
| **Context Meter** | Center-right | Token usage bar + numeric (e.g., "12.4K/32K") |
| **Settings Button** | Right | Gear icon, opens settings dialog |

### 3.3 Settings Dialog

```
┌─────────────────────────────────────────────────────────────┐
│  Settings                                              [×]  │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐                                            │
│  │ LLM Config  │◄── Active Tab                              │
│  ├─────────────┤                                            │
│  │ Prompts     │    ┌─────────────────────────────────────┐ │
│  ├─────────────┤    │ Provider: [OpenAI          ▼]       │ │
│  │ Terminals   │    │                                     │ │
│  ├─────────────┤    │ Model:    [gpt-4-turbo      ▼]      │ │
│  │ Hotkeys     │    │                                     │ │
│  ├─────────────┤    │ API Key:  [••••••••••••••••••]      │ │
│  │ Appearance  │    │                                     │ │
│  └─────────────┘    │ Status:   🟢 Connected               │ │
│                     │                                     │ │
│                     │ [Test Connection]                   │ │
│                     └─────────────────────────────────────┘ │
│                                                             │
│                     ┌─────────────────────────────────────┐ │
│                     │          [Save]  [Cancel]           │ │
│                     └─────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

**Settings Tabs:**

| Tab | Contents |
|-----|----------|
| **LLM Config** | Provider selection, model selection, API key, connection test |
| **Prompts** | View built-in system prompt (read-only), custom prompt extension text area |
| **Terminals** | Auto-detection settings, refresh interval, per-terminal capture settings |
| **Hotkeys** | Global hotkey configuration (copy last command, focus PairAdmin, etc.) |
| **Appearance** | Theme (dark/light), font sizes, sidebar visibility defaults |

### 3.4 Visual Design Guidelines

| Aspect | Specification |
|--------|---------------|
| **Color Scheme** | Dark theme default; light theme available |
| **Primary Accent** | Green (#22c55e) for "Copy to Terminal" and positive actions |
| **Typography** | System UI font (SF Pro, Segoe UI, Ubuntu); Monospace for code (JetBrains Mono, Consolas) |
| **Spacing** | 8px grid system |
| **Border Radius** | 6px for cards/buttons, 4px for inputs |
| **Contrast** | WCAG AA compliance |

---

## 4. Technical Architecture

### 4.1 System Overview

PairAdmin is a **standalone Go application** using a cross-platform GUI framework. It integrates with terminals through platform-native APIs rather than embedding within terminal applications.

### 4.2 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         PairAdmin Application (Go)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         UI Layer (Fyne / Wails)                      │   │
│  │  ┌──────────────┐  ┌──────────────────┐  ┌────────────────────────┐  │   │
│  │  │ Terminal     │  │   Chat View      │  │  Command Sidebar       │  │   │
│  │  │ Tab Manager  │  │                  │  │                        │  │   │
│  │  └──────────────┘  └──────────────────┘  └────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Core Services Layer                             │   │
│  │  ┌──────────────┐  ┌──────────────────┐  ┌────────────────────────┐  │   │
│  │  │ Session      │  │  LLM Gateway     │  │  Clipboard Manager     │  │   │
│  │  │ Manager      │  │                  │  │                        │  │   │
│  │  └──────────────┘  └──────────────────┘  └────────────────────────┘  │   │
│  │  ┌──────────────┐  ┌──────────────────┐  ┌────────────────────────┐  │   │
│  │  │ Context      │  │  Command         │  │  Settings Manager      │  │   │
│  │  │ Manager      │  │  History         │  │                        │  │   │
│  │  └──────────────┘  └──────────────────┘  └────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                  Terminal Integration Layer                          │   │
│  │  ┌──────────────┐  ┌──────────────────┐  ┌────────────────────────┐  │   │
│  │  │ tmux         │  │  macOS           │  │  Windows               │  │   │
│  │  │ Adapter      │  │  Accessibility   │  │  Adapter               │  │   │
│  │  │              │  │  Adapter         │  │  (PuTTY)               │  │   │
│  │  └──────────────┘  └──────────────────┘  └────────────────────────┘  │   │
│  │  ┌──────────────┐                                                    │   │
│  │  │ Linux        │                                                    │   │
│  │  │ Accessibility│                                                    │   │
│  │  │ Adapter      │                                                    │   │
│  │  └──────────────┘                                                    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          External Systems                                    │
│  ┌──────────────┐  ┌──────────────────┐  ┌────────────────────────────────┐ │
│  │ Terminal     │  │  LLM APIs        │  │  System Clipboard              │ │
│  │ Applications │  │  (Cloud/Local)   │  │                                │ │
│  └──────────────┘  └──────────────────┘  └────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.3 Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Language** | Go 1.21+ | Cross-platform compilation, strong concurrency, single binary distribution |
| **GUI Framework** | Wails v2 or Fyne | Wails: native webview UI with Go backend; Fyne: pure Go native UI |
| **HTTP Client** | net/http + resty | LLM API communication |
| **Accessibility** | Platform-specific CGO bindings | macOS: Accessibility API; Linux: AT-SPI2; Windows: UI Automation |
| **Clipboard** | golang.design/x/clipboard | Cross-platform clipboard access |
| **Configuration** | Viper | Settings persistence (YAML/JSON) |
| **Logging** | zerolog | Structured logging |

### 4.4 Core Services

#### 4.4.1 Session Manager

```go
type SessionManager interface {
    // Terminal session lifecycle
    ListSessions() []Session
    GetSession(id string) (*Session, error)
    CreateSession(terminalID string) (*Session, error)
    CloseSession(id string) error
    
    // Active session
    GetActiveSession() *Session
    SetActiveSession(id string) error
}

type Session struct {
    ID           string
    TerminalID   string
    TerminalType TerminalType  // tmux, gnome-terminal, putty, etc.
    Label        string
    CreatedAt    time.Time
    Context      *ContextManager
    ChatHistory  []ChatMessage
    CommandHistory []SuggestedCommand
}
```

#### 4.4.2 Context Manager

```go
type ContextManager interface {
    // Terminal content management
    UpdateContent(content string) error
    GetContent() string
    GetContentWithinTokenLimit(maxTokens int) string
    
    // Token tracking
    GetTokenCount() int
    GetTokenLimit() int
}

type contextManager struct {
    buffer       *CircularBuffer
    maxLines     int
    tokenCounter TokenCounter
}
```

#### 4.4.3 LLM Gateway

```go
type LLMGateway interface {
    // Provider management
    ListProviders() []Provider
    SetProvider(providerID string) error
    GetActiveProvider() Provider
    TestConnection() error
    
    // Completion
    Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
    StreamComplete(ctx context.Context, req CompletionRequest) (<-chan CompletionChunk, error)
}

type CompletionRequest struct {
    SystemPrompt    string
    CustomPrompt    string  // User's prompt extension
    TerminalContext string
    ChatHistory     []ChatMessage
    MaxTokens       int
}
```

#### 4.4.4 Clipboard Manager

```go
type ClipboardManager interface {
    // Basic operations
    Copy(text string) error
    Paste() (string, error)
    
    // Terminal-aware operations
    CopyToTerminal(command string, terminalID string) error
    GetLastCopiedCommand() string
}
```

### 4.5 Data Flow

#### Terminal Capture Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Terminal   │────►│  Adapter    │────►│  Context    │────►│  UI Update  │
│  (tmux,     │     │  (polling   │     │  Manager    │     │  (terminal  │
│  Terminal.  │     │  or event)  │     │             │     │  preview)   │
│  app, etc.) │     │             │     │             │     │             │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
     500ms polling interval (configurable)
```

#### Chat Completion Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  User       │────►│  Context    │────►│  LLM        │────►│  Response   │
│  Message    │     │  Assembly   │     │  Gateway    │     │  Parser     │
└─────────────┘     └─────────────┘     └─────────────┘     └──────┬──────┘
                                                                   │
                    ┌─────────────┐     ┌─────────────┐            │
                    │  Command    │◄────│  Chat       │◄───────────┘
                    │  Sidebar    │     │  Display    │
                    └─────────────┘     └─────────────┘
```

#### Command Execution Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  User       │────►│  Clipboard  │────►│  System     │────►│  User       │
│  Clicks     │     │  Manager    │     │  Clipboard  │     │  Pastes     │
│  "Copy to   │     │             │     │             │     │  in         │
│  Terminal"  │     │             │     │             │     │  Terminal   │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
```

---

## 5. Terminal Integration Specifications

### 5.1 Integration Strategy Overview

PairAdmin uses **platform-native APIs** to capture terminal content automatically. The user does not need to manually copy terminal content—PairAdmin observes it in real-time.

| Terminal | Platform | Capture Method | Polling Interval |
|----------|----------|----------------|------------------|
| tmux | Linux/macOS | `tmux capture-pane` command | 500ms |
| Terminal.app | macOS | Accessibility API (AXUIElement) | 500ms |
| GNOME Terminal | Linux | AT-SPI2 Accessibility API | 500ms |
| Konsole | Linux | AT-SPI2 Accessibility API | 500ms |
| PuTTY | Windows | UI Automation API | 500ms |

### 5.2 tmux Integration

**Method:** Execute `tmux capture-pane` via subprocess

```go
type TmuxAdapter struct {
    pollInterval time.Duration
    sessions     map[string]*TmuxSession
}

func (t *TmuxAdapter) CapturePane(sessionID string, paneID string) (string, error) {
    // tmux capture-pane -t {session}:{pane} -p -S -1000
    cmd := exec.Command("tmux", "capture-pane", 
        "-t", fmt.Sprintf("%s:%s", sessionID, paneID),
        "-p",      // Print to stdout
        "-S", "-1000",  // Start 1000 lines back
    )
    output, err := cmd.Output()
    return string(output), err
}

func (t *TmuxAdapter) ListSessions() ([]TmuxSession, error) {
    // tmux list-sessions -F "#{session_id}:#{session_name}"
    cmd := exec.Command("tmux", "list-sessions", 
        "-F", "#{session_id}:#{session_name}")
    // Parse output...
}

func (t *TmuxAdapter) ListPanes(sessionID string) ([]TmuxPane, error) {
    // tmux list-panes -t {session} -F "#{pane_id}:#{pane_title}"
}
```

**Advantages:**
- No special permissions required
- Works over SSH (for remote tmux sessions)
- Reliable, well-documented API

**Limitations:**
- Requires tmux to be installed and sessions to be running in tmux
- Cannot capture terminals not running in tmux

### 5.3 macOS Terminal.app Integration

**Method:** Accessibility API via CGO bindings

```go
// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa -framework ApplicationServices
// #include "accessibility_darwin.h"
import "C"

type MacOSAccessibilityAdapter struct {
    pollInterval time.Duration
}

func (m *MacOSAccessibilityAdapter) RequestPermission() (bool, error) {
    // Prompt user for Accessibility permission
    granted := C.AXIsProcessTrusted()
    return bool(granted), nil
}

func (m *MacOSAccessibilityAdapter) GetTerminalContent(bundleID string) (string, error) {
    // 1. Get AXUIElement for Terminal.app
    // 2. Navigate to text area
    // 3. Extract AXValue attribute
    cContent := C.GetTerminalTextContent(C.CString(bundleID))
    return C.GoString(cContent), nil
}
```

**Permission Flow:**
1. On first launch, detect if Accessibility permission is granted
2. If not, display dialog explaining why permission is needed
3. Open System Preferences → Security & Privacy → Accessibility
4. User grants permission to PairAdmin
5. PairAdmin can now read terminal content

**Supported Bundle IDs:**
- `com.apple.Terminal` (Terminal.app)

### 5.4 Linux GNOME Terminal / Konsole Integration

**Method:** AT-SPI2 (Assistive Technology Service Provider Interface)

```go
// #cgo pkg-config: atspi-2
// #include <atspi/atspi.h>
import "C"

type LinuxAccessibilityAdapter struct {
    pollInterval time.Duration
}

func (l *LinuxAccessibilityAdapter) Init() error {
    // Initialize AT-SPI2
    C.atspi_init()
    return nil
}

func (l *LinuxAccessibilityAdapter) GetTerminalContent(appName string) (string, error) {
    // 1. Get accessible desktop
    // 2. Find application by name (gnome-terminal, konsole)
    // 3. Navigate to terminal text area
    // 4. Extract text content
}
```

**Dependencies:**
- `libatspi2.0-dev` (Debian/Ubuntu)
- `at-spi2-core` (Fedora/RHEL)

**Permission Notes:**
- AT-SPI2 typically enabled by default on modern Linux desktops
- No special permission prompts required (unlike macOS)

### 5.5 Windows PuTTY Integration

**Method:** UI Automation API

```go
// #cgo LDFLAGS: -lole32 -loleaut32 -luuid
// #include "uiautomation_windows.h"
import "C"

type WindowsUIAutomationAdapter struct {
    pollInterval time.Duration
}

func (w *WindowsUIAutomationAdapter) FindPuTTYWindows() ([]PuTTYWindow, error) {
    // Use UI Automation to find windows with class "PuTTY"
    // Return list of window handles
}

func (w *WindowsUIAutomationAdapter) GetWindowContent(hwnd uintptr) (string, error) {
    // 1. Get IUIAutomationElement for window
    // 2. Find text pattern provider
    // 3. Extract text content
    // Note: PuTTY may require screen scraping as fallback
}
```

**Fallback Strategy:**

If UI Automation cannot extract PuTTY text (PuTTY uses custom rendering), implement screen capture + OCR:

```go
func (w *WindowsUIAutomationAdapter) GetWindowContentOCR(hwnd uintptr) (string, error) {
    // 1. Capture window as bitmap
    // 2. Run through Tesseract OCR
    // 3. Return extracted text
}
```

**Note:** OCR fallback adds latency and reduces accuracy. Document this as a known limitation.

### 5.6 Terminal Detection & Auto-Discovery

```go
type TerminalDetector interface {
    // Discover available terminals
    Scan() ([]DetectedTerminal, error)
    
    // Watch for new terminals
    Watch(ctx context.Context) (<-chan TerminalEvent, error)
}

type DetectedTerminal struct {
    ID           string
    Type         TerminalType
    Name         string         // Session name, window title, etc.
    PID          int            // Process ID (if applicable)
    Accessible   bool           // Can we capture content?
    ErrorMessage string         // If not accessible, why?
}

type TerminalEvent struct {
    Type     EventType  // Created, Closed, TitleChanged
    Terminal DetectedTerminal
}
```

---

## 6. Security Constraints and Considerations

### 6.1 Threat Model

PairAdmin has **read access to terminal content** which may contain sensitive information. The security model must protect against:

1. **Credential Exposure:** Passwords, API keys, tokens visible in terminal
2. **Data Exfiltration:** Sensitive terminal content sent to external LLM APIs
3. **Command Injection:** Malicious commands suggested by AI
4. **Unauthorized Access:** Other users/processes reading PairAdmin's data

### 6.2 Sensitive Data Handling

#### Pre-LLM Filtering Pipeline

All terminal content passes through a filtering pipeline before being sent to the LLM:

```go
type SensitiveDataFilter interface {
    Filter(content string) string
    AddPattern(pattern FilterPattern) error
    RemovePattern(id string) error
    ListPatterns() []FilterPattern
}

type FilterPattern struct {
    ID          string
    Name        string
    Regex       *regexp.Regexp
    Action      FilterAction  // Redact, Exclude, Mask
    Replacement string        // For Redact action
}
```

**Built-in Patterns:**

| Pattern Name | Regex | Action |
|--------------|-------|--------|
| Password Prompts | `(?i)(password|passphrase|passwd|pin):\s*$` | Exclude line |
| AWS Access Keys | `AKIA[0-9A-Z]{16}` | Redact → `[AWS_KEY_REDACTED]` |
| AWS Secret Keys | `[A-Za-z0-9/+=]{40}` (in context) | Redact → `[AWS_SECRET_REDACTED]` |
| Private Keys | `-----BEGIN.*PRIVATE KEY-----` | Exclude block |
| Generic API Keys | `(api[_-]?key|apikey|access[_-]?token)\s*[:=]\s*['"]?[A-Za-z0-9]{20,}` | Redact |
| Database URLs | `(mysql|postgres|mongodb)://[^:]+:[^@]+@` | Mask password |

#### User-Configurable Patterns

Users can add custom patterns via `/filter add` command or Settings dialog:

```
/filter add "internal-hosts" "internal-[a-z]+-[0-9]+\.corp\.example\.com" redact "[INTERNAL_HOST]"
```

### 6.3 Local Model Support

For organizations with strict data residency requirements, PairAdmin supports local LLM execution:

| Provider | Setup | Data Location |
|----------|-------|---------------|
| **Ollama** | `ollama serve` on localhost | All data stays on local machine |
| **LM Studio** | Download model, start server | All data stays on local machine |
| **llama.cpp** | Self-hosted server | Self-managed infrastructure |

When local models are configured, **no terminal data leaves the machine**.

### 6.4 Clipboard Security

| Risk | Mitigation |
|------|------------|
| Clipboard history exposure | Clear PairAdmin-copied commands from clipboard after 60 seconds (configurable) |
| Cross-application clipboard sniffing | Document limitation; recommend clipboard managers with app-specific history |

### 6.5 Credential Storage

| Credential Type | Storage Method |
|-----------------|----------------|
| LLM API Keys | OS keychain (macOS Keychain, Windows Credential Manager, libsecret on Linux) |
| User preferences | Plain YAML/JSON (no secrets) |
| Chat history | Local SQLite database (optional encryption) |

### 6.6 Audit Logging

All AI interactions are logged locally for compliance review:

```go
type AuditEntry struct {
    Timestamp     time.Time
    SessionID     string
    TerminalID    string
    EventType     string    // "user_message", "ai_response", "command_copied"
    Content       string    // Sanitized (filtered) content
    CommandCopied string    // If applicable
}
```

Logs stored in: `~/.pairadmin/logs/audit-YYYY-MM-DD.jsonl`

---

## 7. Slash Commands Reference

| Command | Parameters | Description |
|---------|------------|-------------|
| `/help` | `[command]` | Display all commands or details for specific command |
| `/clear` | — | Clear chat history for current tab |
| `/context` | `<lines>` | Set terminal context window size (default: 100) |
| `/refresh` | — | Force re-capture of terminal content |
| `/model` | `<model-id>` | Switch active LLM provider/model |
| `/filter` | `<add\|remove\|list>` | Manage sensitive data filter patterns |
| `/export` | `<json\|txt>` | Export current session chat history |
| `/rename` | `<new-label>` | Rename current terminal tab |
| `/theme` | `<dark\|light>` | Switch UI color scheme |
| `/hotkey` | `<action> <key>` | Configure global hotkeys |

### Command Examples

```bash
# Set context to last 200 lines of terminal output
/context 200

# Switch to local Ollama model
/model ollama:llama2

# Add custom filter for company secrets
/filter add "vault-tokens" "hvs\.[A-Za-z0-9]{24}" redact "[VAULT_TOKEN]"

# Export session for incident review
/export json

# Rename tab for clarity
/rename prod-database-primary
```

---

## 8. Implementation Roadmap

### Phase 1: Foundation (Weeks 1-6)

| Week | Milestone | Deliverables |
|------|-----------|--------------|
| 1-2 | Project Setup | Go project structure, CI/CD pipeline, GUI framework evaluation (Wails vs Fyne) |
| 3-4 | Core UI | Three-column layout, tab management, chat display |
| 5-6 | LLM Integration | OpenAI integration, basic chat completion, context assembly |

**Exit Criteria:** User can chat with AI about manually-pasted terminal content.

### Phase 2: Terminal Integration (Weeks 7-12)

| Week | Milestone | Deliverables |
|------|-----------|--------------|
| 7-8 | tmux Adapter | Full tmux capture, session detection, auto-polling |
| 9-10 | macOS Adapter | Accessibility API integration, permission flow, Terminal.app capture |
| 11-12 | Linux Adapter | AT-SPI2 integration, GNOME Terminal + Konsole capture |

**Exit Criteria:** Auto-capture working on tmux (all platforms), Terminal.app (macOS), GNOME Terminal (Linux).

### Phase 3: Windows & Polish (Weeks 13-18)

| Week | Milestone | Deliverables |
|------|-----------|--------------|
| 13-14 | Windows Adapter | UI Automation integration, PuTTY capture (with OCR fallback if needed) |
| 15-16 | Command Sidebar | History tracking, click-to-copy, hover details |
| 17-18 | Settings & Filtering | Full settings dialog, sensitive data filtering, multi-provider support |

**Exit Criteria:** All v1.0 terminal integrations working. Settings fully functional.

### Phase 4: Hardening & Launch (Weeks 19-24)

| Week | Milestone | Deliverables |
|------|-----------|--------------|
| 19-20 | Security Hardening | Credential storage, audit logging, security review |
| 21-22 | Testing & QA | Cross-platform testing, edge case handling, performance optimization |
| 23-24 | Documentation & Launch | User guide, installation packages (DMG, MSI, AppImage), public release |

**Exit Criteria:** Stable release with installers for all platforms.

---

## 9. Appendix

### A. Glossary

| Term | Definition |
|------|------------|
| **LLM** | Large Language Model - AI system for natural language understanding and generation |
| **tmux** | Terminal multiplexer allowing multiple terminal sessions within a single window |
| **AT-SPI2** | Assistive Technology Service Provider Interface - Linux accessibility framework |
| **UI Automation** | Windows accessibility framework for programmatic UI interaction |
| **Context Window** | Maximum text an LLM can process in one request (measured in tokens) |
| **Wails** | Go framework for building desktop apps with web frontends |
| **Fyne** | Pure Go cross-platform GUI toolkit |

### B. GUI Framework Decision Matrix

| Criterion | Wails | Fyne |
|-----------|-------|------|
| **Native Look** | Uses system webview | Custom rendering |
| **Bundle Size** | ~10MB | ~15MB |
| **Development Speed** | Fast (HTML/CSS/JS) | Moderate (Go only) |
| **Cross-Platform** | Excellent | Excellent |
| **Accessibility** | Via web standards | Limited |
| **CGO Dependency** | Yes | Yes |

**Recommendation:** Wails v2 for faster UI iteration and better accessibility support.

### C. References

1. Go Programming Language: https://go.dev/
2. Wails Framework: https://wails.io/
3. Fyne Toolkit: https://fyne.io/
4. tmux Manual: https://man7.org/linux/man-pages/man1/tmux.1.html
5. macOS Accessibility API: https://developer.apple.com/documentation/applicationservices/axuielement_h
6. AT-SPI2 Documentation: https://docs.gtk.org/atspi2/
7. Windows UI Automation: https://docs.microsoft.com/en-us/windows/win32/winauto/entry-uiauto-win32

### D. Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | January 2026 | Product Team | Initial PRD (PuTTY-embedded approach) |
| 2.0 | February 2026 | Product Team | Complete pivot to standalone Go application |

---

*This document is confidential and intended for internal use by the PairAdmin development team.*
