# PairAdmin v2.0 – Project History and Maintenance Guide

## Executive Summary

**PairAdmin** is a cross-platform desktop application that enables **“Pair Administration”**—a collaboration model where human system administrators work alongside AI agents to manage Linux, Unix, macOS, and Windows systems via terminal interfaces. The AI observes terminal sessions in real-time and provides contextual command suggestions that users can execute with a single hotkey or click.

This document is the single source of truth for the project’s history, architecture, build process, and maintenance guidelines. It is intended to give future maintainers—human or AI-assisted—full context to extend, fix, and evolve the codebase. All other Markdown files (except the PRD and user documentation) may be deleted after this guide is verified.

---

## Table of Contents

1. [Project History and Evolution](#1-project-history-and-evolution)
2. [Architecture and Design Decisions](#2-architecture-and-design-decisions)
3. [Development Environment and Setup](#3-development-environment-and-setup)
4. [Build and Packaging](#4-build-and-packaging)
5. [Release Process](#5-release-process)
6. [Known Issues and Troubleshooting](#6-known-issues-and-troubleshooting)
7. [Extension Points and Maintenance Guidelines](#7-extension-points-and-maintenance-guidelines)
8. [Development Methodology (Orchestrator)](#8-development-methodology-orchestrator)
9. [Future Considerations](#9-future-considerations)
10. [References](#10-references)
11. [License](#11-license)

## 1. Project History and Evolution

### 1.1 Origin and Pivot

The project began as **PairAdmin v1.0**, an embedded AI extension for the PuTTY terminal emulator on Windows. After extensive prototyping, the embedded approach proved untenable due to:

- Complex UI integration failures between PuTTY’s Win32 rendering and modern UI frameworks.
- Unexpected terminal I/O edge cases creating unreliable capture.
- Single-platform focus (Windows/PuTTY only) limiting market reach.
- Maintenance burden of tracking upstream PuTTY changes.

**v2.0 adopted a fundamentally different architecture:** a standalone cross-platform application that integrates with existing terminals through well-defined system APIs rather than embedding within them.

### 1.2 Development Phases

The implementation followed a four-phase roadmap:

| Phase | Focus | Key Deliverables |
|-------|-------|------------------|
| **Phase 1** | Foundation | Wails project scaffold, basic UI, terminal detection skeleton |
| **Phase 2** | Core AI Integration | LLM gateway, multi-provider support (OpenAI, Anthropic, Ollama), command history, clipboard manager |
| **Phase 3** | Windows & macOS Terminal Adapters | Windows UI Automation (CGo bindings), macOS Accessibility API, tmux adapter, integration tests |
| **Phase 4** | Security Hardening & Packaging | OS keychain integration, audit logging, configuration race fixes, GitHub Actions CI/CD, multi-platform packaging |

### 1.3 Release v2.0.0

The **v2.0.0 release** (April 2, 2026) marks the completion of Phase 4. It includes:

- **Security enhancements**: OS keychain integration (macOS Keychain, Windows Credential Manager, Linux Secret Service), tamper-evident audit logging, atomic configuration operations.
- **Local AI support**: Ollama provider for running models locally, token estimation, connection pooling.
- **Comprehensive testing**: 100% unit-test coverage of core modules, integration and end-to-end test suites.
- **Production packaging**: Automated GitHub Actions workflows that build macOS `.app` bundles, Windows `.exe` binaries, and Linux executables on tag push.
- **CI/CD pipeline**: GitHub Actions workflow (`release.yml`) with matrix strategy (`fail-fast: false`), Linux dependency resolution (`libatspi2.0-dev`, `webkit2gtk-4.1-dev`), Electron version fix (`electron@^28.0.0`), and platform-specific build constraints (`//go:build linux` for `atspi2.c`).
- **Release artifacts**: Three binary assets uploaded to GitHub release:
  - `pairadmin-linux` (Linux executable)
  - `PairAdmin-mac.zip` (macOS .app bundle)
  - `pairadmin-windows.exe` (Windows executable)

All binaries are available on the [GitHub releases page](https://github.com/o3willard-AI/pa-sup-dp-orc/releases/tag/v2.0.0).

**Release checklist**: The release process followed a detailed checklist (`RELEASE_CHECKLIST.md`) covering dependency validation, build testing, artifact verification, and publication steps. Future releases should replicate this structured approach.

### 1.4 Development Methodology

The project was built using a **multi-tier LLM orchestration system** (the “Orchestrator”) that routes development tasks through a hierarchy of AI models based on complexity and cost. This approach enabled rapid iteration, rigorous code review, and systematic escalation of difficult problems. The Orchestrator’s templates and logging infrastructure reside in `orchestrator/` and `docs/workflow/`; they are not required for building or running PairAdmin but are preserved for future AI-assisted development.

---

## 2. Architecture and Design Decisions

### 2.1 High-Level Structure

PairAdmin is a **Wails v2** application with a Go backend and a Svelte/Vite frontend. The backend is organized into domain-focused packages under `internal/`:

```
internal/
├── adapters/          # Terminal adapter factory (platform-specific registration)
├── audit/             # Tamper-evident audit logging
├── clipboard/         # Cross-platform clipboard operations
├── config/            # Configuration management with OS keychain integration
├── hotkeys/           # Global hotkey manager (platform-specific implementations pending)
├── llm/               # LLM gateway and provider implementations (OpenAI, Anthropic, Ollama)
├── security/          # Sensitive-data filtering and redaction
├── session/           # SQLite store for command history and session state
├── terminal/          # Terminal detection and adapter interfaces
│   ├── windows/       # Windows UI Automation CGo bindings
│   ├── macos/         # macOS Accessibility API stubs
│   ├── tmux/          # tmux adapter (pure Go)
│   └── testhelpers/   # Mock adapters for testing
└── ui/                # Wails-exposed handlers for terminal and chat operations
```

### 2.2 Key Design Decisions

#### 2.2.1 Terminal Integration Strategy

Instead of embedding inside terminal emulators, PairAdmin uses **platform-native accessibility APIs**:

- **Windows**: UI Automation (UIA) via CGo bindings (`ui_automation.c`/`.h`). The adapter enumerates terminal windows and captures text through the `IUIAutomation` COM interface.
- **macOS**: Accessibility API (stubbed; full implementation requires entitlements and user permission).
- **Linux**: AT-SPI2 (stubbed; full implementation requires D-Bus integration).
- **tmux**: Pure-Go adapter that attaches to existing tmux sessions via the `tmux` CLI.

This approach ensures broad compatibility without modifying terminal source code.

#### 2.2.2 Configuration and Secrets

- **Configuration file**: `~/.pairadmin/config.yaml` (auto-created with defaults).
- **Secrets storage**: OS-native keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service) via the `config.Keychain` interface. Secrets are never written to disk in plaintext.
- **Atomic operations**: The configuration manager uses file locking and atomic writes to prevent race conditions and partial writes.

#### 2.2.3 LLM Gateway

- **Multi-provider architecture**: Plug-in providers for OpenAI, Anthropic, and Ollama (local). Providers implement a simple `Provider` interface.
- **Token estimation**: Accurate token counting using rune-based heuristics (provider-specific tokenizers are not required).
- **Connection pooling**: Reusable HTTP clients with timeouts and retries.
- **Sensitive-data filtering**: The `security` package redacts passwords, API keys, and other sensitive patterns before sending content to LLMs.

#### 2.2.4 Audit Logging

- **Tamper-evident logs**: `~/.pairadmin/audit.log` records all security-sensitive operations (keychain access, configuration changes, LLM requests) with SHA-256 hashes chained across entries.
- **Immutable records**: Once written, log entries cannot be altered without breaking the hash chain.

#### 2.2.5 Frontend-Backend Communication

- **Wails bindings**: The `App` struct in `app.go` exposes methods to the frontend (e.g., `SendMessage`, `GetCommandsByTerminal`).
- **UI handlers**: The `ui` package contains `TerminalHandlers` and `ChatHandlers` that orchestrate terminal detection, LLM requests, and clipboard operations.
- **Session store**: SQLite database (`sessions.db`) persists command history, usage counts, and terminal-session associations.

### 2.3 Platform-Specific Notes

#### Windows
- Requires `UI Automation` COM APIs (available since Windows Vista).
- CGo source files (`ui_automation.c`, `ui_automation.h`) are guarded by `//go:build windows && !mock`.
- The Windows adapter is the only fully implemented native adapter; it works with Windows Terminal, cmd.exe, and PowerShell.

#### macOS
- The Accessibility API requires the `com.apple.security.accessibility` entitlement and user approval in System Settings.
- Current implementation is a stub (`accessibility_stub.go`); a full CGo implementation is planned for a future release.

#### Linux
- AT-SPI2 integration requires D-Bus and the `atspi-2` development package.
- The stub (`atspi2.c`, `atspi2.h`) is guarded by `//go:build linux`.
- The tmux adapter provides full functionality for tmux sessions.

### 2.4 Dependencies

| Package | Purpose | License |
|---------|---------|---------|
| **Wails v2** | GUI framework | MIT |
| **go-vgo/robotgo** | Keyboard/mouse simulation (hotkeys) | MIT |
| **golang.design/x/clipboard** | Cross-platform clipboard | MIT |
| **samber/lo** | Go generics utilities | MIT |
| **testify** | Testing framework | MIT |
| **modernc.org/sqlite** | SQLite driver (pure Go) | BSD-3-Clause |

All dependencies are vendored via `go.mod`; the frontend uses npm packages managed in `frontend/package.json`.

---

## 3. Development Environment and Setup

### 3.1 Prerequisites

- **Go** 1.21 or later
- **Node.js** 20.x and npm
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Platform-specific build tools**:
  - **Linux**: `libgtk-3-dev`, `libwebkit2gtk-4.1-dev`, `libatspi2.0-dev`, `libx11-dev`, `pkg-config`
  - **macOS**: Xcode Command Line Tools
  - **Windows**: MSYS2 or Visual Studio Build Tools

### 3.2 One-Command Setup

Run the dependency installer:

```bash
./scripts/install-deps.sh
```

This script:
1. Installs Wails CLI.
2. Downloads Go dependencies (`go mod download`).
3. Installs frontend packages (`npm install` in `frontend/`).
4. On Linux, installs system packages via `apt-get` and creates a `webkit2gtk-4.0.pc` symlink for compatibility with Wails.
5. Installs root-level npm devDependencies (Electron-builder fallback).

### 3.3 Development Workflow

```bash
# Live development (hot-reload)
wails dev

# Production build for current platform
wails build

# Platform-specific build
wails build -platform linux   # or windows, darwin
```

The built binary appears in `build/bin/` (Linux/Windows) or `build/bin/PairAdmin.app` (macOS).

### 3.4 Testing

```bash
# Unit tests
go test ./internal/... -short

# Integration tests
./scripts/test-integration.sh

# End-to-end tests (requires webkit2gtk-4.1)
./scripts/test-e2e.sh
```

**Coverage**: Run `go test ./internal/... -coverprofile=coverage.out && go tool cover -html=coverage.out`.

---

## 4. Build and Packaging

### 4.1 Wails Configuration

`wails.json` defines:
- Application metadata (name, version, author).
- Platform-specific packaging options (Debian dependencies, NSIS installer directory, macOS bundle identifier).
- Frontend build commands.

### 4.2 GitHub Actions CI/CD

Two workflows are defined in `.github/workflows/`:

1. **`release.yml`** – Triggered on tag pushes (`v*`). Builds on Ubuntu, Windows, and macOS runners using a matrix strategy with `fail-fast: false`. Uploads artifacts named `pairadmin-<os>-latest`.
2. **`build.yml`** – Manual workflow for testing builds on the `master` branch.

### 4.3 Artifact Processing

The release workflow produces:
- **Linux**: A single executable `pairadmin`.
- **macOS**: A bundled `.app` inside a zip archive.
- **Windows**: A single `pairadmin.exe`.

These artifacts are uploaded to the GitHub release via the API (see Section 5).

### 4.4 Electron-builder Fallback

A `package.json` at the project root provides an Electron-builder configuration as a fallback packaging method. It is not used by the standard Wails build but is kept for alternative distribution strategies.

---

## 5. Release Process

### 5.1 Pre-Release Checklist

1. All unit tests pass (`go test ./internal/...`).
2. Integration tests pass (`./scripts/test-integration.sh`).
3. End-to-end tests pass (`./scripts/test-e2e.sh`).
4. No known critical bugs.
5. Documentation updated (user guide, installation guide).
6. Version number updated in `wails.json`.
7. Release notes written (`RELEASE_NOTES_v2.0.0.md`).
8. Git tag `v2.0.0` created and pushed.

### 5.2 Automated Packaging

1. Push the tag to GitHub: `git push origin v2.0.0`.
2. GitHub Actions runs the `release.yml` workflow, building binaries for all three platforms.
3. Download artifacts from the workflow run via the GitHub API.

### 5.3 Creating the GitHub Release

```bash
# Create release with notes
curl -X POST -H "Authorization: token <PAT>" \
  -H "Content-Type: application/json" \
  https://api.github.com/repos/o3willard-AI/pa-sup-dp-orc/releases \
  -d '{"tag_name":"v2.0.0","name":"PairAdmin v2.0.0","body":"$(cat RELEASE_NOTES_v2.0.0.md)","draft":false,"prerelease":false}'

# Upload each binary
curl -X POST -H "Authorization: token <PAT>" \
  -H "Content-Type: application/octet-stream" \
  --data-binary @pairadmin-linux \
  "https://uploads.github.com/.../assets?name=pairadmin-linux"
```

### 5.4 Post-Release Tasks

- Update website/download page (if applicable).
- Announce on social channels.
- Monitor error logs and gather user feedback.
- Plan next iteration.

### 5.5 CI/CD Configuration Notes

- **Repository authentication**: The CI workflow uses a GitHub Personal Access Token (PAT) embedded in the remote URL (`https://<PAT>@github.com/o3willard-AI/pa-sup-dp-orc.git`). This allows pushing tags and creating releases without manual authentication. For security, rotate the PAT periodically and use fine-grained tokens with minimal permissions.
- **Workflow file**: `.github/workflows/release.yml` defines the matrix build strategy with `fail-fast: false` to ensure all platform builds complete even if one fails.
- **Dependency installation**: The `scripts/install-deps.sh` script installs system packages and npm dependencies; it also creates symlinks for webkit2gtk-4.0.pc on Ubuntu 24.04.
- **Artifact naming**: The workflow uploads artifacts named `pairadmin-<platform>-latest`; the release step renames them to user-friendly names.

---

## 6. Known Issues and Troubleshooting

### 6.1 Build Issues

| Symptom | Cause | Solution |
|---------|-------|----------|
| **Linux build fails with “Package atspi-2 not found”** | Missing `libatspi2.0-dev` | Install via `apt-get install libatspi2.0-dev`. |
| **Linux build fails with “X11/Xlib.h not found”** | Missing X11 development packages | Install `libx11-dev`, `libxrandr-dev`, `libxinerama-dev`, `libxcursor-dev`, `libxi-dev`, `libxtst-dev`. |
| **Linux build fails with “webkit2gtk-4.0.pc not found”** | Ubuntu 24.04 provides `webkit2gtk-4.1.pc` | The install script creates a symlink; manually run `sudo ln -sf /usr/lib/x86_64-linux-gnu/pkgconfig/webkit2gtk-4.1.pc /usr/lib/x86_64-linux-gnu/pkgconfig/webkit2gtk-4.0.pc`. |
| **npm install fails with “Cannot find package 'electron'”** | Root `package.json` specifies `"electron": "^latest"` which npm cannot resolve | Change `package.json` to use a specific version like `"electron": "^28.0.0"`. |
| **macOS build fails with “C source files not allowed when not using cgo”** | `atspi2.c` included in non-Linux builds | Ensure `//go:build linux` is present at the top of `atspi2.c`. |
| **Windows build fails with UI Automation errors** | Missing Windows SDK or COM headers | Install Visual Studio Build Tools with C++ and Windows SDK components. |

### 6.2 Runtime Issues

| Symptom | Cause | Solution |
|---------|-------|----------|
| **Terminal not detected** | Adapter for current platform not implemented or permission denied. | On Windows, ensure UI Automation is enabled. On macOS, grant accessibility permission in System Settings. Use tmux adapter on Linux. |
| **AI responses slow** | Network latency or Ollama not running. | Check internet connection; for Ollama, verify `ollama serve` is running. |
| **Configuration not saving** | Write permissions for `~/.pairadmin/`. | Fix directory permissions or run with appropriate user privileges. |
| **Hotkeys not working** | Platform-specific implementation pending. | Hotkey manager is registered but not started (`hotkeyMgr.Start()` is commented out). Uncomment after implementing platform-specific listeners. |

### 6.3 False Positives

- **LSP errors in terminal adapter types**: Some Go language servers report type mismatches in adapter interfaces; unit tests pass, so these are false positives.
- **Coverage gaps in platform-specific CGo code**: Go coverage tools cannot instrument C files; focus on integration tests for validation.

---

## 7. Extension Points and Maintenance Guidelines

### 7.1 Adding a New Terminal Adapter

1. Create a new package under `internal/terminal/<platform>`.
2. Implement the `TerminalAdapter` interface:
   ```go
   type TerminalAdapter interface {
       Name() string
       Available(ctx context.Context) bool
       ListSessions(ctx context.Context) ([]DetectedTerminal, error)
       Capture(ctx context.Context, terminalID string) (string, error)
       Subscribe(ctx context.Context, terminalID string) (<-chan TerminalEvent, error)
       GetDimensions(ctx context.Context, terminalID string) (int, int, error)
   }
   ```
3. Register the adapter in the appropriate factory file (`factory.go`, `factory_darwin.go`, `factory_windows.go`).
4. Write unit and integration tests (see `windows/adapter_test.go` for an example).

### 7.2 Adding a New LLM Provider

1. Create a new file in `internal/llm/providers/`.
2. Implement the `Provider` interface:
   ```go
   type Provider interface {
       Name() string
       Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
       EstimateTokens(text string) int
   }
   ```
3. Register the provider in `llm/gateway.go` (`defaultProviders` map).
4. Add configuration fields in `config.Config.LLM`.

### 7.3 Modifying the Configuration Schema

1. Update `config.Config` struct in `config/manager.go`.
2. Provide sensible defaults in `defaultConfig()`.
3. Ensure backward compatibility: old configurations should be automatically migrated (see `migrateConfig` function).
4. Update the keychain integration if new secret fields are added.

### 7.4 Enhancing Security

- **New sensitive patterns**: Add regex patterns to `security/redact.go`.
- **Audit log events**: Use `audit.Logger.LogEvent()` for new security-sensitive operations.
- **Keychain operations**: Always use `config.Keychain` methods for storing/retrieving secrets.

### 7.5 Platform-Specific Code

- Guard platform-specific files with build tags (`//go:build windows`, `//go:build darwin`, `//go:build linux`).
- Keep CGo code isolated in `.c`/`.h` files with minimal exported Go functions.
- Provide stub implementations for other platforms to avoid compilation errors.

### 7.6 Testing Strategy

- **Unit tests**: Every exported function should have a test.
- **Integration tests**: Test adapters against real terminals (use Docker for isolation where possible).
- **End-to-end tests**: Simulate full user workflows (terminal detection → AI request → clipboard copy).
- **Platform-specific tests**: Use build tags to run tests only on relevant platforms.

---

## 8. Development Methodology (Orchestrator)

### 8.1 Overview

The **Multi-Tier LLM Orchestrator** is a meta-development tool that was used to build PairAdmin. It routes tasks through a cascade of AI models:

- **L0-Planner** (Qwen3.5 397B): Creates detailed task specifications.
- **L0-Coder** (Qwen3-Coder local): First implementation attempt.
- **L1-Coder** (Grok 4.1 Fast): Re-implementation if L0 output is rejected.
- **L2-Coder** (MiniMax M2.7): Complex task escalation.
- **L3-Architect** (Claude 3.7 Sonnet): Architectural review and final approval.

### 8.2 Artifacts

The Orchestrator produces:

- **Task specifications** in `docs/tasks/`.
- **Handoff logs** in `docs/workflow/handoffs/`.
- **Escalation records** in `docs/workflow/escalations/`.
- **Learning documents** in `docs/workflow/learnings/`.

### 8.3 Usage

While not required for maintaining PairAdmin, the Orchestrator can be reused for future AI-assisted development. The templates are in `orchestrator/templates/` and `docs/workflow/templates/`.

---

## 9. Future Considerations

### 9.1 Immediate Next Steps

- **macOS Accessibility implementation**: Complete the CGo bindings for the macOS Accessibility API.
- **Linux AT-SPI2 implementation**: Integrate with D-Bus for native Linux terminal detection.
- **Hotkey platform implementations**: Implement `hotkeys.Manager.Start()` for Windows, macOS, and Linux.
- **Additional terminal adapters**: Support for iTerm2, GNOME Terminal, Konsole, etc.

### 9.2 Long-Term Vision

- **Plugin system**: Allow third-party terminal adapters and LLM providers.
- **Collaboration features**: Multi-user session sharing and real-time collaboration.
- **Advanced AI features**: Automated remediation, predictive command suggestions, anomaly detection.

### 9.3 Deprecation Plan

- **v1.x (embedded PuTTY)**: Already deprecated; no further development.
- **Current stubs**: The macOS and Linux adapter stubs should be replaced with full implementations before declaring them stable.

---

## 10. References

### 10.1 Essential Files

| File | Purpose |
|------|---------|
| `PairAdmin_PRD_v2.0.md` | Product requirements and design rationale. |
| `docs/user-guide.md` | End-user documentation. |
| `docs/installation.md` | Installation instructions. |
| `PROJECT_HISTORY_AND_MAINTENANCE.md` | This document. |
| `wails.json` | Wails build configuration. |
| `internal/config/manager.go` | Configuration management and keychain integration. |
| `internal/llm/gateway.go` | LLM provider orchestration. |
| `internal/terminal/windows/adapter.go` | Reference adapter implementation. |
| `.github/workflows/release.yml` | CI/CD pipeline definition. |
| `scripts/install-deps.sh` | Development environment setup. |

### 10.2 External Dependencies

- [Wails v2 Documentation](https://wails.io/docs/)
- [Windows UI Automation](https://docs.microsoft.com/en-us/windows/win32/winauto/entry-uiauto-win32)
- [macOS Accessibility API](https://developer.apple.com/documentation/appkit/nsaccessibility)
- [AT-SPI2 (Linux)](https://developer.gnome.org/atspi/)
- [Ollama](https://ollama.ai/)

---

## 11. License

PairAdmin is released under the **MIT License**. See `LICENSE` file for details.

---

*Document version: 1.0*  
*Last updated: April 3, 2026*  
*Maintainer: AI-assisted workflow with OpenCode*