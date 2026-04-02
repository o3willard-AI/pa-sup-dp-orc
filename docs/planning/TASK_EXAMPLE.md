# Example Atomic Task Specification

## Task ID: 1.1
## Title: Initialize Go module
## Phase: 1 (Foundation)
## Estimated effort: 0.5 hour
## Dependencies: None

### Description
Create the `go.mod` file for the PairAdmin project, specifying Go 1.21+ as the minimum version and setting the module path to `github.com/yourorg/pairadmin` (or a placeholder).

### Inputs
- None (fresh directory)

### Outputs
- File: `go.mod` with content:
```
module github.com/yourorg/pairadmin

go 1.21

require (
    // dependencies will be added later
)
```

### Steps
1. Navigate to project root directory.
2. Run `go mod init github.com/yourorg/pairadmin`.
3. Edit `go.mod` to explicitly set `go 1.21` (or higher).
4. Verify file exists and is syntactically correct.

### Verification
- Run `go mod tidy` (should succeed with no errors).
- Run `go version` to confirm Go 1.21+ is installed.

---

## Task ID: 1.2
## Title: Scaffold Wails project
## Phase: 1 (Foundation)
## Estimated effort: 1 hour
## Dependencies: Task 1.1 (Go module)

### Description
Use the Wails CLI to scaffold a new Wails v2 project in the current directory, generating the basic project structure and configuration files.

### Inputs
- `go.mod` from Task 1.1
- Wails CLI installed globally (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Outputs
- `wails.json` configuration file
- `frontend/` directory with placeholder assets
- `build/` directory (optional)
- `app.go` in project root (or `internal/ui/app.go` depending on template)

### Steps
1. Ensure Wails CLI is installed: `wails version`.
2. Run `wails init -n pairadmin -t svelte` (or `-t vue`/`-t react`). Choose Vue template for easier integration with PRD wireframes.
3. Answer prompts (author, description, etc.) with placeholder values.
4. Verify that `wails.json` is created and contains correct project name.
5. Move `app.go` to `internal/ui/app.go` and update package declaration accordingly (optional but aligns with directory structure).

### Verification
- Run `wails build` to ensure project compiles without errors.
- Run `wails dev` to see a blank application window (optional).

---

## Task ID: 2.1
## Title: Terminal adapter interface
## Phase: 2 (Terminal Integration)
## Estimated effort: 1 hour
## Dependencies: Phase 1 tasks (basic project structure)

### Description
Define the Go interface and data structures that all terminal adapters must implement. This provides a contract for tmux, macOS, Linux, and Windows adapters.

### Inputs
- Existing project structure (`internal/terminal/` directory)

### Outputs
- File: `internal/terminal/types.go` with:
  - `TerminalType` enum
  - `TerminalAdapter` interface
  - `DetectedTerminal` struct
  - `TerminalEvent` struct and `EventType` enum

### Steps
1. Create directory `internal/terminal/` if it doesn't exist.
2. Create `types.go` with package `terminal`.
3. Define `TerminalType` as a string enum: `Tmux`, `MacOSTerminal`, `GNOMETerminal`, `Konsole`, `PuTTY`, `Unknown`.
4. Define `DetectedTerminal` struct with fields: `ID`, `Type`, `Name`, `PID`, `Accessible`, `ErrorMessage`.
5. Define `EventType` enum: `Created`, `Closed`, `TitleChanged`.
6. Define `TerminalEvent` struct with `Type` and `Terminal`.
7. Define `TerminalAdapter` interface with methods: `Capture(terminalID string) (string, error)`, `ListSessions() ([]DetectedTerminal, error)`, `Watch(ctx context.Context) (<-chan TerminalEvent, error)`, `Close()`.

### Verification
- Compile project: `go build ./internal/terminal`.
- Run `go test ./internal/terminal` (no tests yet, but should pass).

---

## Notes for LLM Execution
- Each task should be provided as a standalone prompt with the above format.
- The LLM should have access to the current directory state (via tools like `ls`, `cat`).
- After completing a task, the LLM should verify the outputs and report any issues.
- Dependencies must be satisfied before starting a task; if not, the LLM should request the missing artifacts.
- Tasks can be parallelized where dependencies allow (e.g., UI components can be built simultaneously with backend services as long as interfaces are stable).
