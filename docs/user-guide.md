# PairAdmin User Guide

## Getting Started

1. Install PairAdmin from [releases page]
2. Launch the application
3. Configure your LLM provider in Settings → LLM
4. Select a terminal session
5. Ask for AI assistance in the chat area

## Features

### AI‑Assisted Commands
- Type natural language questions in the chat area
- AI suggests terminal commands
- One‑click copy to clipboard

### Command History
- All suggested commands are saved per terminal
- See usage counts and timestamps
- Filter by terminal session

### Security
- API keys stored in OS keychain
- Sensitive data filtered before sending to LLM
- Audit logging of all actions

## Troubleshooting

**Q: AI responses are slow**
A: Check your internet connection. For local LLM (Ollama), ensure Ollama is running.

**Q: Terminal not detected**
A: Make sure you're using a supported terminal (bash, zsh, tmux, Windows Terminal).

**Q: Configuration not saving**
A: Check write permissions for `~/.pairadmin/config.yaml`.
