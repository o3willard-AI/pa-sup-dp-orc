# PairAdmin v2.0 Implementation Plan

## Overview
This document provides a detailed, atomized implementation plan for PairAdmin v2.0, a cross-platform AI-assisted terminal administration tool. The plan is structured to be executable by an open locally run model with a 32K context window, breaking down tasks into small, independent units.

## Technology Stack
- **Language**: Go 1.21+
- **GUI Framework**: Wails v2 (recommended for faster UI iteration and accessibility)
- **Frontend**: HTML/CSS/JavaScript (Vue.js or React optional; PRD doesn't specify)
- **Package Management**: Go modules, npm/yarn
- **Build Tools**: Wails CLI, platform-specific packaging (DMG, MSI, AppImage)

## Directory Structure
```
pairadmin/
├── cmd/
│   └── pairadmin/
│       └── main.go                 # Application entry point
├── internal/
│   ├── session/                    # Session management
│   │   ├── manager.go
│   │   ├── session.go
│   │   └── store.go
│   ├── context/                    # Terminal context management
│   │   ├── manager.go
│   │   ├── buffer.go
│   │   └── tokenizer.go
│   ├── llm/                        # LLM gateway and providers
│   │   ├── gateway.go
│   │   ├── providers/
│   │   │   ├── openai.go
│   │   │   ├── anthropic.go
│   │   │   ├── ollama.go
│   │   │   └── local.go
│   │   └── completion.go
│   ├── clipboard/                  # Cross-platform clipboard operations
│   │   ├── manager.go
│   │   └── platform/
│   │       ├── darwin.go
│   │       ├── linux.go
│   │       └── windows.go
│   ├── terminal/                   # Terminal adapters
│   │   ├── detector.go
│   │   ├── types.go
│   │   ├── tmux/
│   │   │   ├── adapter.go
│   │   │   └── session.go
│   │   ├── macos/
│   │   │   ├── adapter.go
│   │   │   └── accessibility.go
│   │   ├── linux/
│   │   │   ├── adapter.go
│   │   │   └── atspi.go
│   │   └── windows/
│   │       ├── adapter.go
│   │       └── ui_automation.go
│   ├── config/                     # Configuration management
│   │   ├── manager.go
│   │   ├── defaults.go
│   │   └── keychain.go
│   ├── security/                   # Sensitive data filtering
│   │   ├── filter.go
│   │   ├── patterns.go
│   │   └── redactor.go
│   ├── audit/                      # Audit logging
│   │   ├── logger.go
│   │   └── entry.go
│   └── ui/                         # Backend UI bindings (Wails)
│       ├── app.go
│       ├── handlers.go
│       └── events.go
├── frontend/                       # Wails frontend assets
│   ├── src/
│   │   ├── components/             # Vue/React components
│   │   │   ├── TerminalTabs.vue
│   │   │   ├── ChatArea.vue
│   │   │   ├── CommandSidebar.vue
│   │   │   ├── StatusBar.vue
│   │   │   └── SettingsDialog.vue
│   │   ├── views/
│   │   │   └── MainWindow.vue
│   │   ├── stores/                 # State management (Pinia/Vuex)
│   │   │   ├── session.js
│   │   │   ├── terminal.js
│   │   │   └── settings.js
│   │   ├── assets/
│   │   │   └── styles/
│   │   │       └── main.css
│   │   └── main.js                 # Frontend entry point
│   ├── index.html
│   ├── package.json
│   └── wails.js                    # Wails runtime integration
├── build/                          # Build artifacts
│   ├── darwin/
│   ├── linux/
│   └── windows/
├── scripts/                        # Helper scripts
│   ├── install-deps.sh
│   └── build-all.sh
├── go.mod
├── go.sum
├── wails.json                      # Wails configuration
├── .gitignore
└── README.md
```

## Phase 1: Foundation (Weeks 1-6)
### Week 1-2: Project Setup
**Goal**: Establish a working Go + Wails project with basic CI/CD.

#### Task 1.1: Initialize Go module and Wails project
- Create `go.mod` with Go 1.21+.
- Run `wails init` to scaffold project with Wails v2.
- Configure `wails.json` with appropriate project settings.
- Verify project builds and runs a blank window.

#### Task 1.2: Set up development environment
- Create `.gitignore` for Go, Wails, and OS-specific files.
- Create `scripts/install-deps.sh` to install Wails CLI, npm dependencies.
- Set up pre-commit hooks for Go formatting and linting.

#### Task 1.3: Establish CI/CD pipeline
- Create GitHub Actions workflow (or similar) for building on Linux, macOS, Windows.
- Add steps to run tests, linting, and building binaries.
- Store workflow in `.github/workflows/build.yml`.

### Week 3-4: Core UI Implementation
**Goal**: Implement the three-column layout and basic tab management.

#### Task 1.4: Create main window layout (HTML/CSS)
- Create `frontend/src/views/MainWindow.vue` with three-column flex layout.
- Implement left column (terminal tabs), center (chat area), right (command sidebar).
- Ensure responsive design with minimum widths.

#### Task 1.5: Terminal tabs component
- Create `frontend/src/components/TerminalTabs.vue` with static list of tabs.
- Implement active tab highlighting and click-to-switch behavior.
- Add "+ New" button placeholder.

#### Task 1.6: Chat area component
- Create `frontend/src/components/ChatArea.vue` with message bubbles.
- Implement user messages (right-aligned) and AI messages (left-aligned).
- Add input field at bottom with send button.

#### Task 1.7: Command sidebar component
- Create `frontend/src/components/CommandSidebar.vue` with static command cards.
- Implement reverse chronological ordering.
- Add hover tooltip placeholder.

#### Task 1.8: Status bar component
- Create `frontend/src/components/StatusBar.vue` with model selector, connection status, context meter, settings button.
- Use placeholder data.

### Week 5-6: LLM Integration (Basic)
**Goal**: Integrate with OpenAI API for basic chat completions.

#### Task 1.9: Create LLM gateway interface
- Define `internal/llm/gateway.go` with `LLMGateway` interface.
- Implement provider-agnostic completion request/response structures.

#### Task 1.10: Implement OpenAI provider
- Create `internal/llm/providers/openai.go` with API client using `net/http`.
- Support chat completions endpoint with system prompt, user message, terminal context.
- Handle API errors and timeouts.

#### Task 1.11: Create configuration manager
- Create `internal/config/manager.go` to load/save settings via Viper.
- Support YAML/JSON config files in `~/.pairadmin/config.yaml`.
- Store LLM provider, API key, model selection.

#### Task 1.12: Connect frontend to backend
- Create `internal/ui/app.go` with Wails bindings for sending messages to LLM.
- Expose `SendMessage(terminalID, message)` function to frontend.
- Return AI response and display in chat area.

**Exit Criteria**: User can manually paste terminal content into UI, ask questions, and receive AI responses via OpenAI.

## Phase 2: Terminal Integration (Weeks 7-12)
### Week 7-8: tmux Adapter
**Goal**: Capture terminal content from tmux sessions automatically.

#### Task 2.1: Define terminal adapter interface
- Create `internal/terminal/types.go` with `TerminalAdapter` interface (`Capture`, `ListSessions`, etc.).
- Define `DetectedTerminal` and `TerminalEvent` structures.

#### Task 2.2: Implement tmux adapter
- Create `internal/terminal/tmux/adapter.go` using `exec.Command` to run `tmux capture-pane`.
- Parse `tmux list-sessions` and `tmux list-panes` output.
- Poll each pane at configurable interval (default 500ms).

#### Task 2.3: Terminal detector service
- Create `internal/terminal/detector.go` that scans for tmux sessions periodically.
- Emit events when sessions are created/closed.
- Map terminal IDs to internal session IDs.

#### Task 2.4: Integrate terminal capture with UI
- Update `internal/ui/handlers.go` to push terminal content updates to frontend.
- Display captured content in terminal preview pane (left column).
- Switch terminal tabs when user selects a detected tmux session.

### Week 9-10: macOS Adapter (Terminal.app)
**Goal**: Capture content from macOS Terminal.app via Accessibility API.

#### Task 2.5: macOS accessibility CGO bindings
- Create `internal/terminal/macos/accessibility.h` and `.m` files for Objective-C bridging.
- Implement functions to request permissions, get window list, extract text from Terminal.app.
- Use CGO to call from Go.

#### Task 2.6: macOS adapter implementation
- Create `internal/terminal/macos/adapter.go` that uses CGO bindings.
- Handle permission flow: if not granted, show dialog and open System Preferences.
- Poll Terminal.app windows for content.

#### Task 2.7: Extend terminal detector for macOS
- Update detector to scan for Terminal.app windows via Accessibility API.
- Differentiate between tmux sessions inside Terminal.app and raw Terminal.app windows.

### Week 11-12: Linux Adapter (GNOME Terminal, Konsole)
**Goal**: Capture content from Linux terminals via AT-SPI2.

#### Task 2.8: Linux AT-SPI2 CGO bindings
- Create `internal/terminal/linux/atspi.h` with C bindings to AT-SPI2 library.
- Implement functions to initialize AT-SPI2, find terminal applications, extract text.

#### Task 2.9: Linux adapter implementation
- Create `internal/terminal/linux/adapter.go` using AT-SPI2 bindings.
- Support GNOME Terminal and Konsole (detect by application name).
- Poll each terminal window.

#### Task 2.10: Permission handling on Linux
- Document that AT-SPI2 may require `gsettings set org.gnome.desktop.interface accessibility true`.
- Provide troubleshooting guide.

**Exit Criteria**: Auto-capture working on tmux (all platforms), Terminal.app (macOS), GNOME Terminal (Linux). UI updates with live terminal content.

## Phase 3: Windows & Polish (Weeks 13-18)
### Week 13-14: Windows Adapter (PuTTY)
**Goal**: Capture content from PuTTY via UI Automation API.

#### Task 3.1: Windows UI Automation CGO bindings
- Create `internal/terminal/windows/ui_automation.h` with C bindings to UI Automation COM API.
- Implement functions to find PuTTY windows, extract text via `IUIAutomationTextPattern`.

#### Task 3.2: Windows adapter implementation
- Create `internal/terminal/windows/adapter.go` using UI Automation.
- Fallback to screen capture + OCR if text extraction fails (optional later).

#### Task 3.3: Terminal detection on Windows
- Extend detector to scan for PuTTY windows and other terminals (Windows Terminal may require different approach).

### Week 15-16: Command Sidebar & Clipboard Integration
**Goal**: Implement command history sidebar and one-click command execution.

#### Task 3.4: Command history storage
- Extend `internal/session/` to store `SuggestedCommand` slices.
- Persist commands per session across app restarts (SQLite optional).

#### Task 3.5: Clipboard manager
- Create `internal/clipboard/manager.go` using `golang.design/x/clipboard`.
- Implement `CopyToTerminal` that copies command to clipboard and optionally focuses target terminal.

#### Task 3.6: Frontend command sidebar integration
- Update `CommandSidebar.vue` to fetch command history from backend.
- Implement click-to-copy: when a command card is clicked, call backend's `CopyToTerminal`.
- Show tooltip with original question context on hover.

#### Task 3.7: "Copy to Terminal" button in chat
- Add button to AI-suggested command code blocks in chat.
- On click, copy command to clipboard and optionally auto-paste (if configured).

### Week 17-18: Settings & Filtering
**Goal**: Implement settings dialog and sensitive data filtering.

#### Task 3.8: Settings dialog UI
- Create `frontend/src/components/SettingsDialog.vue` with tabs: LLM Config, Prompts, Terminals, Hotkeys, Appearance.
- Connect each tab to backend configuration.

#### Task 3.9: LLM provider configuration
- Extend `internal/config/` to support multiple providers (OpenAI, Anthropic, Ollama, etc.).
- Add UI to select provider, enter API key, test connection.

#### Task 3.10: Sensitive data filter
- Create `internal/security/filter.go` with regex patterns for passwords, API keys, etc.
- Integrate filter into context manager: filter terminal content before sending to LLM.
- Provide UI to add custom patterns via `/filter` command or settings.

#### Task 3.11: Hotkey configuration
- Use global hotkey library (e.g., `github.com/micmonay/keybd_event` for Linux, `github.com/go-vgo/robotgo` cross-platform).
- Allow user to configure hotkey for copying last command, focusing PairAdmin.

**Exit Criteria**: All v1.0 terminal integrations working. Settings fully functional.

## Phase 4: Hardening & Launch (Weeks 19-24)
### Week 19-20: Security Hardening
**Goal**: Implement credential storage, audit logging, and security review.

#### Task 4.1: Secure credential storage
- Use OS keychain: macOS Keychain, Windows Credential Manager, libsecret on Linux.
- Create `internal/config/keychain.go` abstraction.

#### Task 4.2: Audit logging
- Create `internal/audit/logger.go` that writes sanitized logs to `~/.pairadmin/logs/audit-*.jsonl`.
- Log events: user message, AI response, command copied.

#### Task 4.3: Local LLM support (Ollama)
- Implement `internal/llm/providers/ollama.go` for local Ollama server.
- Ensure no data leaves machine when local model is selected.

### Week 21-22: Testing & QA
**Goal**: Cross-platform testing, edge case handling, performance optimization.

#### Task 4.4: Unit and integration tests
- Write Go tests for core services (session, context, LLM, terminal adapters).
- Mock external dependencies (tmux, Accessibility APIs).

#### Task 4.5: End-to-end testing
- Create script to simulate user interactions (using UI automation tools).
- Test on each platform via CI if possible (macOS requires real GUI).

#### Task 4.6: Performance optimization
- Profile terminal capture polling to avoid high CPU usage.
- Optimize context token counting and truncation.

### Week 23-24: Documentation & Launch
**Goal**: Create user guide, installation packages, and public release.

#### Task 4.7: User documentation
- Write `README.md` with installation instructions for each platform.
- Create `docs/` directory with advanced configuration, troubleshooting.

#### Task 4.8: Packaging
- Use Wails build system to generate:
  - macOS: `.dmg` and `.app`
  - Windows: `.msi` installer
  - Linux: `.AppImage` and `.deb`/`.rpm` packages

#### Task 4.9: Release checklist
- Verify all acceptance criteria from PRD.
- Conduct final security review.
- Publish binaries to GitHub Releases.

**Exit Criteria**: Stable release with installers for all platforms.

## Atomic Task List for 32K Context Window
Each task below is designed to be independently described and implemented within a 32K context window. They correspond to the tasks above but are further decomposed where needed.

### Phase 1 Tasks
1. **Initialize Go module** – Create `go.mod` with Go 1.21+.
2. **Scaffold Wails project** – Run `wails init`, configure `wails.json`.
3. **Set up .gitignore** – Add patterns for Go, Wails, OS files.
4. **Create install-deps script** – Install Wails CLI and npm dependencies.
5. **Create CI workflow** – GitHub Actions for building on three OSes.
6. **Main window layout (HTML/CSS)** – Three-column flex layout.
7. **Terminal tabs component (static)** – List of tabs with active highlighting.
8. **Chat area component** – Message bubbles, input field, send button.
9. **Command sidebar component (static)** – Cards with placeholder commands.
10. **Status bar component** – Model selector, connection status, context meter, settings button.
11. **LLM gateway interface** – Define `LLMGateway` and related structs.
12. **OpenAI provider implementation** – HTTP client for completions.
13. **Configuration manager** – Load/save YAML config with Viper.
14. **Wails bindings for chat** – Expose `SendMessage` to frontend.
15. **Connect frontend to backend** – Call `SendMessage` and display response.

### Phase 2 Tasks
16. **Terminal adapter interface** – Define `TerminalAdapter` and data types.
17. **tmux adapter: capture pane** – Execute `tmux capture-pane` and parse output.
18. **tmux adapter: list sessions** – Parse `tmux list-sessions`.
19. **tmux adapter: list panes** – Parse `tmux list-panes`.
20. **Terminal detector service** – Scan for tmux sessions periodically.
21. **Integrate terminal capture with UI** – Push updates to frontend preview.
22. **macOS accessibility CGO bindings (header)** – Objective‑C functions.
23. **macOS accessibility CGO bindings (implementation)** – Request permission, get window text.
24. **macOS adapter implementation** – Use CGO bindings to poll Terminal.app.
25. **Extend detector for macOS** – Scan Terminal.app windows.
26. **Linux AT‑SPI2 CGO bindings** – C functions to interact with AT‑SPI2.
27. **Linux adapter implementation** – Use AT‑SPI2 to poll GNOME Terminal/Konsole.
28. **Permission documentation for Linux** – Guide for enabling accessibility.

### Phase 3 Tasks
29. **Windows UI Automation CGO bindings** – C functions for UI Automation API.
30. **Windows adapter implementation** – Use UI Automation to capture PuTTY.
31. **Extend detector for Windows** – Scan for PuTTY windows.
32. **Command history storage** – Store `SuggestedCommand` per session.
33. **Clipboard manager** – Cross‑platform copy/paste using `golang.design/x/clipboard`.
34. **Frontend command sidebar integration** – Fetch history, click‑to‑copy.
35. **“Copy to Terminal” button in chat** – Add button to code blocks.
36. **Settings dialog UI (tabs)** – Create Vue component with tab navigation.
37. **LLM provider configuration UI** – Dropdowns, API key field, test button.
38. **Sensitive data filter implementation** – Regex patterns, redaction.
39. **Hotkey configuration** – Global hotkey library integration.

### Phase 4 Tasks
40. **OS keychain integration** – Store API keys securely.
41. **Audit logging** – Write JSONL logs for compliance.
42. **Ollama provider implementation** – Local LLM support.
43. **Unit tests for core services** – Session, context, LLM, terminal adapters.
44. **End‑to‑end testing script** – Simulate user interactions.
45. **Performance profiling** – Optimize polling and token counting.
46. **User documentation** – README and advanced guides.
47. **Packaging for macOS** – Generate `.dmg` and `.app`.
48. **Packaging for Windows** – Generate `.msi`.
49. **Packaging for Linux** – Generate `.AppImage`, `.deb`, `.rpm`.
50. **Final release checklist** – Verify acceptance criteria, security review.

## Notes for Execution
- Each atomic task should be assigned to a single LLM invocation with clear input/output specifications.
- The frontend can be built with any framework (Vue.js recommended for Wails), but the plan is framework‑agnostic.
- Terminal adapters rely on platform‑specific APIs; CGO bindings must be carefully tested on each OS.
- Sensitive data filtering is critical; ensure redaction occurs before any content leaves the machine (except for local LLMs).
- The roadmap assumes a 24‑week timeline but can be adjusted based on team size and priorities.

This plan provides a granular decomposition of the PairAdmin v2.0 implementation, enabling an open locally run model with a 32K context window to execute each task independently.