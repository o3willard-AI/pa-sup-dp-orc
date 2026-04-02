# PairAdmin v2.0 QA Checkpoints

## Overview
This document defines QA checkpoints at natural completion points in the PairAdmin implementation. Each checkpoint validates that a set of atomic tasks has been completed successfully and that the system is ready for downstream development. Checkpoints are designed to catch upstream bugs before they propagate.

## Checkpoint Structure
Each checkpoint includes:
- **Milestone ID**: Reference to the milestone
- **Tasks Included**: Which atomic tasks are covered
- **Validation Criteria**: Specific tests and checks to perform
- **Required Artifacts**: Files and outputs that must exist
- **Success Criteria**: Conditions that must be met to pass
- **Blocking Issues**: Common failures that would block progress

---

## Milestone 1: Project Setup Complete
**Tasks Included**: 1–5 (Initialize Go module, Scaffold Wails project, Set up .gitignore, Create install‑deps script, Create CI workflow)

### Validation Criteria
1. **Go module**: `go.mod` exists with `go 1.21+` and correct module path.
2. **Wails project**: `wails.json` exists with appropriate settings; `wails build` succeeds without errors.
3. **Directory structure**: `frontend/`, `internal/`, `cmd/` directories exist (or equivalent Wails layout).
4. **Git ignore**: `.gitignore` includes patterns for Go binaries, OS‑specific files, node_modules, etc.
5. **Install script**: `scripts/install-deps.sh` (or equivalent) installs Wails CLI and npm dependencies.
6. **CI workflow**: `.github/workflows/build.yml` (or other CI config) builds on at least Linux; no syntax errors.

### Required Artifacts
- `go.mod`, `wails.json`, `.gitignore`
- `scripts/install-deps.sh`
- CI configuration file
- Successful `go mod tidy` output

### Success Criteria
- Project compiles on the development machine (`wails build`).
- CI pipeline passes (or at least runs without configuration errors).
- All files are committed to version control.

### Blocking Issues
- Wails CLI not installed globally (should be covered by install script).
- Go version < 1.21.
- Missing `wails.json` or misconfigured project name.

---

## Milestone 2: Core UI Complete
**Tasks Included**: 6–10 (Main window layout, Terminal tabs component, Chat area component, Command sidebar component, Status bar component)

### Validation Criteria
1. **Three‑column layout**: Browser inspect shows flex/grid layout with left, center, right columns.
2. **Terminal tabs**: Left column displays list of tabs (static) with active highlighting; "+ New" button present.
3. **Chat area**: Center column shows message bubbles (user right, AI left), input field, send button.
4. **Command sidebar**: Right column shows cards with placeholder commands, reverse chronological order.
5. **Status bar**: Bottom bar contains model selector, connection status, context meter, settings button.
6. **Responsive design**: Layout adapts to window resizing (minimum widths respected).

### Required Artifacts
- `frontend/src/views/MainWindow.vue` (or equivalent)
- `frontend/src/components/TerminalTabs.vue`
- `frontend/src/components/ChatArea.vue`
- `frontend/src/components/CommandSidebar.vue`
- `frontend/src/components/StatusBar.vue`
- CSS files defining the visual design (dark theme default)

### Success Criteria
- Application runs (`wails dev`) and displays the three‑column UI without JavaScript errors.
- All UI components render; no missing images or broken styles.
- Basic interactivity: clicking tabs changes active state, input field accepts text.

### Blocking Issues
- Vue/React component syntax errors.
- Missing CSS causing layout collapse.
- JavaScript runtime errors in browser console.

---

## Milestone 3: Basic LLM Integration
**Tasks Included**: 11–15 (LLM gateway interface, OpenAI provider, Configuration manager, Wails bindings, Frontend‑backend connection)

### Validation Criteria
1. **LLM gateway interface**: `internal/llm/gateway.go` defines `LLMGateway` interface and `CompletionRequest`/`Response` structs.
2. **OpenAI provider**: `internal/llm/providers/openai.go` implements the interface, makes HTTP calls to OpenAI API.
3. **Configuration manager**: `internal/config/manager.go` loads/saves YAML config; stores provider, API key, model.
4. **Wails bindings**: `internal/ui/app.go` exports `SendMessage(terminalID, message)` to frontend.
5. **Frontend integration**: Chat input calls backend `SendMessage`, displays AI response in chat area.
6. **API key handling**: Configuration reads API key from file; provider uses it for authentication.

### Required Artifacts
- All Go files listed above.
- `~/.pairadmin/config.yaml` (or equivalent) with placeholder API key.
- Frontend JavaScript that calls `window.backend.SendMessage`.

### Success Criteria
- With a valid OpenAI API key, user can type a question in chat and receive an AI response.
- Response appears in chat area with proper formatting (code blocks if suggested command).
- Error handling: invalid API key shows appropriate error message.

### Blocking Issues
- Missing API key causes crash.
- HTTP requests fail due to network/timeout.
- Frontend‑backend binding not working (method not exposed).

---

## Milestone 4: tmux Integration
**Tasks Included**: 16–21 (Terminal adapter interface, tmux capture, list sessions, list panes, Terminal detector, UI integration)

### Validation Criteria
1. **Terminal adapter interface**: `internal/terminal/types.go` defines `TerminalAdapter`, `DetectedTerminal`, etc.
2. **tmux adapter capture**: `internal/terminal/tmux/adapter.go` executes `tmux capture-pane` and returns text.
3. **tmux list sessions/panes**: Adapter parses `tmux list-sessions` and `tmux list-panes` output.
4. **Terminal detector**: `internal/terminal/detector.go` scans for tmux sessions periodically (configurable interval).
5. **UI integration**: Detected tmux sessions appear as tabs in left column; terminal preview shows captured content.
6. **Polling**: Content updates every 500ms (configurable) without high CPU usage.

### Required Artifacts
- All terminal adapter Go files.
- At least one active tmux session on the test machine.
- Terminal preview pane in UI showing live terminal output.

### Success Criteria
- PairAdmin detects running tmux sessions and creates corresponding tabs.
- Terminal preview displays the last ~10 lines of the tmux pane.
- Content updates when user types in the tmux pane (within 500ms).
- Selecting a different tmux tab switches the preview to that pane.

### Blocking Issues
- tmux not installed or not in PATH.
- Permission issues (tmux server not accessible).
- Parsing errors in `tmux` output formats.

---

## Milestone 5: macOS Integration (Terminal.app)
**Tasks Included**: 22–25 (macOS accessibility CGO bindings, macOS adapter, Extend detector)

### Validation Criteria
1. **CGO bindings**: `internal/terminal/macos/accessibility.h/.m` implement Objective‑C functions to access Terminal.app via Accessibility API.
2. **Permission flow**: Adapter checks for Accessibility permission; if not granted, guides user to System Preferences.
3. **macOS adapter**: `internal/terminal/macos/adapter.go` uses CGO bindings to poll Terminal.app windows.
4. **Detector extension**: Terminal detector includes Terminal.app windows in scan results.
5. **Differentiation**: Detector distinguishes between tmux sessions inside Terminal.app and raw Terminal.app windows.

### Required Artifacts
- CGO source files and Go adapter.
- Terminal.app running with at least one window.
- Accessibility permission granted (may require manual approval).

### Success Criteria
- PairAdmin detects Terminal.app windows and creates tabs for them.
- Terminal preview shows content from Terminal.app windows.
- Content updates as user types in Terminal.app.
- If permission not granted, user is informed with clear instructions.

### Blocking Issues
- macOS Accessibility permission not granted (cannot capture content).
- CGO compilation errors on non‑macOS platforms.
- Terminal.app bundle ID changed or not found.

---

## Milestone 6: Linux Integration (GNOME Terminal, Konsole)
**Tasks Included**: 26–28 (Linux AT‑SPI2 CGO bindings, Linux adapter, Permission documentation)

### Validation Criteria
1. **AT‑SPI2 bindings**: `internal/terminal/linux/atspi.h` provides C functions to interact with AT‑SPI2.
2. **Linux adapter**: `internal/terminal/linux/adapter.go` uses bindings to poll GNOME Terminal and Konsole.
3. **Detection**: Adapter identifies both GNOME Terminal and Konsole by application name.
4. **Permission documentation**: README or docs explain how to enable AT‑SPI2 (`gsettings set org.gnome.desktop.interface accessibility true`).
5. **Integration**: Detected Linux terminals appear in tabs and preview.

### Required Artifacts
- CGO source files and Go adapter.
- GNOME Terminal or Konsole running on a Linux desktop with AT‑SPI2 enabled.
- Documentation in `docs/linux-setup.md`.

### Success Criteria
- PairAdmin detects GNOME Terminal/Konsole windows and creates tabs.
- Terminal preview shows live content.
- Content updates within polling interval.
- Documentation exists for enabling accessibility on Linux.

### Blocking Issues
- AT‑SPI2 not installed or disabled.
- Desktop environment not GNOME/KDE (may not support AT‑SPI2).
- CGO compilation requires `libatspi2.0-dev`.

---

## Milestone 7: Windows Integration (PuTTY)
**Tasks Included**: 29–31 (Windows UI Automation CGO bindings, Windows adapter, Extend detector)

### Validation Criteria
1. **UI Automation bindings**: `internal/terminal/windows/ui_automation.h` provides C functions for Windows UI Automation API.
2. **Windows adapter**: `internal/terminal/windows/adapter.go` uses UI Automation to capture PuTTY window text.
3. **Fallback noted**: Documentation mentions OCR fallback for PuTTY if text extraction fails (optional).
4. **Detector extension**: Terminal detector scans for PuTTY windows (class "PuTTY").
5. **Integration**: Detected PuTTY windows appear in tabs and preview.

### Required Artifacts
- CGO source files and Go adapter.
- PuTTY running on Windows test machine.
- UI Automation API available (Windows 7+).

### Success Criteria
- PairAdmin detects PuTTY windows and creates tabs.
- Terminal preview shows content from PuTTY (if UI Automation works).
- Application runs on Windows without crashes.

### Blocking Issues
- UI Automation cannot extract PuTTY text (may require OCR fallback).
- Windows‑specific CGO code fails to compile on non‑Windows.
- PuTTY window class name differs.

---

## Milestone 8: Command Execution
**Tasks Included**: 32–35 (Command history storage, Clipboard manager, Frontend sidebar integration, “Copy to Terminal” button)

### Validation Criteria
1. **Command history storage**: `internal/session/` stores `SuggestedCommand` slices per session; persists across restarts (SQLite optional).
2. **Clipboard manager**: `internal/clipboard/manager.go` uses `golang.design/x/clipboard` for cross‑platform copy/paste.
3. **Sidebar integration**: `CommandSidebar.vue` fetches history from backend; click‑to‑copy works.
4. **Copy button in chat**: AI‑suggested commands in chat have a “Copy to Terminal” button that copies to clipboard.
5. **Tooltip hover**: Hover over command in sidebar shows original question context.

### Required Artifacts
- Updated session manager and clipboard manager.
- Updated frontend components.
- SQLite database file (if persistence implemented).

### Success Criteria
- AI‑suggested commands appear in sidebar.
- Clicking a command in sidebar copies it to system clipboard.
- “Copy to Terminal” button in chat copies command.
- Clipboard content can be pasted into a terminal.

### Blocking Issues
- Clipboard access permissions (Linux may require `xclip` or `xsel`).
- SQLite driver not imported.
- Frontend not updating sidebar when new commands generated.

---

## Milestone 9: Settings & Filtering
**Tasks Included**: 36–39 (Settings dialog UI, LLM provider configuration, Sensitive data filter, Hotkey configuration)

### Validation Criteria
1. **Settings dialog UI**: `SettingsDialog.vue` with tabs: LLM Config, Prompts, Terminals, Hotkeys, Appearance.
2. **LLM provider config**: UI to select provider (OpenAI, Anthropic, Ollama), enter API key, test connection.
3. **Sensitive data filter**: `internal/security/filter.go` applies regex patterns to redact passwords, API keys, etc. before sending to LLM.
4. **Filter UI**: Users can add custom patterns via `/filter` command or settings dialog.
5. **Hotkey configuration**: UI to set global hotkeys (copy last command, focus PairAdmin); backend registers hotkeys.

### Required Artifacts
- Settings dialog component.
- Security filter implementation.
- Hotkey library integrated.
- Configuration saved to disk.

### Success Criteria
- Settings dialog opens from status bar gear icon.
- Changing LLM provider and saving works; test connection succeeds with valid key.
- Filter redacts a test password (e.g., `password: secret` → `password: [REDACTED]`).
- Hotkey can be configured and triggers action (e.g., Ctrl+Shift+C copies last command).

### Blocking Issues
- Configuration not persisting between app launches.
- Filter regex patterns too aggressive (breaking legitimate content).
- Hotkey conflicts with OS or other applications.

---

## Milestone 10: Security Hardening
**Tasks Included**: 40–42 (OS keychain integration, Audit logging, Local LLM support)

### Validation Criteria
1. **OS keychain**: `internal/config/keychain.go` abstracts macOS Keychain, Windows Credential Manager, libsecret on Linux.
2. **API key storage**: LLM API keys stored in OS keychain, not plain config file.
3. **Audit logging**: `internal/audit/logger.go` writes sanitized JSONL logs to `~/.pairadmin/logs/audit-*.jsonl`.
4. **Log events**: User messages, AI responses, commands copied (with sensitive data redacted).
5. **Local LLM support**: `internal/llm/providers/ollama.go` connects to local Ollama server; no data leaves machine.

### Required Artifacts
- Keychain abstraction code.
- Audit log files.
- Ollama provider implementation.

### Success Criteria
- API keys are retrieved from OS keychain (not visible in config file).
- Audit logs are created and contain expected events (with redaction).
- Selecting Ollama provider works when Ollama is running locally; no external network calls.
- All data filtering occurs before any content leaves the machine for cloud LLMs.

### Blocking Issues
- OS keychain APIs not available on some Linux distributions.
- Audit logs fail to write due to permission issues.
- Ollama provider cannot connect to local server.

---

## Milestone 11: Testing Infrastructure
**Tasks Included**: 43–45 (Unit tests for core services, End‑to‑end testing script, Performance profiling)

### Validation Criteria
1. **Unit tests**: Go tests for session, context, LLM, terminal adapters (mocked external dependencies).
2. **Test coverage**: Critical paths have >70% line coverage (measured by `go test -cover`).
3. **E2E script**: Script that simulates user interactions (opens app, sends message, copies command) using UI automation tools.
4. **Performance profiling**: Terminal capture polling uses <5% CPU when idle; token counting efficient.
5. **Memory leaks**: No obvious leaks after prolonged operation (checked with Go pprof).

### Required Artifacts
- `*_test.go` files in relevant packages.
- `scripts/e2e-test.sh` (or similar).
- Performance profiling reports.

### Success Criteria
- `go test ./...` passes all tests.
- E2E script runs without errors (may require manual intervention for GUI).
- CPU usage remains low during idle polling.
- Memory usage stable over time.

### Blocking Issues
- Unit tests fail due to missing mocks.
- E2E script requires GUI automation tools not installed.
- Performance issues with terminal capture polling.

---

## Milestone 12: Packaging
**Tasks Included**: 46–49 (User documentation, Packaging for macOS, Windows, Linux)

### Validation Criteria
1. **User documentation**: `README.md` includes installation instructions for each platform, basic usage.
2. **Advanced docs**: `docs/` directory with configuration, troubleshooting, security notes.
3. **macOS package**: Wails generates `.dmg` and `.app` bundle; codesigned (if possible).
4. **Windows package**: Wails generates `.msi` installer; runs on Windows 10/11.
5. **Linux package**: Wails generates `.AppImage`; optionally `.deb` and `.rpm` packages.
6. **Installation**: Packages install successfully on clean VMs of each OS.

### Required Artifacts
- README.md and docs/ files.
- Generated installer files (`.dmg`, `.msi`, `.AppImage`).
- Installation test results.

### Success Criteria
- Documentation is clear and covers common use cases.
- Packages install without errors on target OS.
- Installed application launches and passes smoke test (UI loads, can open settings).

### Blocking Issues
- Wails packaging fails due to missing dependencies.
- Installers not signed (may trigger OS warnings).
- Linux package depends on libraries not present on target systems.

---

## Comprehensive QA Execution

Each checkpoint should be executed before moving to tasks that depend on its components. For example, Milestone 4 (tmux Integration) must pass before tasks that assume terminal capture works (e.g., command execution).

### QA Process
1. **Automated checks**: Run `go test`, `wails build`, `npm run lint` (if configured).
2. **Manual verification**: Follow validation criteria manually; document results.
3. **Integration testing**: Test interactions between newly implemented components and existing ones.
4. **Regression testing**: Ensure previous milestones still pass.

### Failure Handling
If a checkpoint fails:
- Identify the specific validation criterion that failed.
- Determine root cause (upstream bug, missing artifact, environment issue).
- Fix the issue before proceeding with dependent tasks.
- Re‑run the checkpoint validation.

### Sign‑off
Each milestone should be signed off by the developer (or QA) before downstream tasks are assigned. Sign‑off includes:
- All validation criteria met.
- Required artifacts present and functional.
- No blocking issues outstanding.

---

## Notes
- Some checkpoints require platform‑specific hardware/OS (macOS, Windows, Linux). Plan testing accordingly.
- The CI pipeline should run automated tests for the current platform; cross‑platform tests may require separate runners.
- Security‑related checkpoints (Milestone 10) are critical and should involve dedicated security review.