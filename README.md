# README

## About

This is the official Wails Svelte template.

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

## Phase 3: AI Collaboration Features

- **LLM Provider Integration**: Support for OpenAI, Anthropic, and Ollama (local)
- **Command History**: Stores AI‑suggested commands per terminal session in SQLite
- **Clipboard Manager**: One‑click copy to clipboard with optional terminal focus
- **Sensitive Data Filtering**: Redacts passwords, API keys, etc. before sending to LLM
- **Settings Dialog**: Configure providers, terminals, hotkeys, and appearance
- **Hotkey Support**: Global shortcuts for copying last command and focusing the app

## Testing

### Unit Tests
```bash
go test ./internal/... -short
```

### Integration Tests
```bash
./scripts/test-integration.sh
```

### End‑to‑End Tests
```bash
./scripts/test-e2e.sh
```
