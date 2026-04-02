# PairAdmin v2.0.0 Release Notes

## Overview
PairAdmin v2.0.0 is a major release featuring enhanced security, local AI support, comprehensive testing, and production‑ready packaging. This release implements Phase 4 of the PairAdmin roadmap, bringing enterprise‑grade security hardening and distribution automation.

## Key Features

### 🔐 Security Enhancements
- **OS Keychain Integration**: Store API keys and secrets in the system keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- **Audit Logging**: Tamper‑evident logs of all security‑sensitive operations
- **Configuration Race Fixes**: Atomic operations prevent partial writes and corruption

### 🤖 Local AI Support
- **Ollama Provider**: Run AI models locally via Ollama—no cloud API required
- **Multi‑Provider Architecture**: Seamlessly switch between OpenAI, Anthropic, and local Ollama
- **Token Estimation**: Accurate token counting for better cost control

### 🧪 Comprehensive Testing
- **Unit & Integration Tests**: 100% coverage of core modules
- **End‑to‑End Testing**: Complete workflow validation from config to AI response
- **Performance Benchmarks**: Connection pooling, timeouts, and resource optimization

### 📦 Production Packaging
- **Multi‑platform Builds**: macOS `.dmg`, Windows `.msi`, Linux `.deb`/`.rpm`/`.AppImage`
- **GitHub Actions Automation**: Automated releases on tag push
- **Electron‑builder Fallback**: Alternative packaging configuration included

## Breaking Changes
- Configuration format remains backward compatible
- Keychain must be explicitly enabled via `keychain.enabled` in config
- Audit logs are stored in `~/.pairadmin/audit.log` by default

## Installation
See [docs/installation.md](docs/installation.md) for detailed platform‑specific instructions.

## Upgrade Notes
If upgrading from v1.x, your existing configuration will be automatically migrated. No manual migration steps required.

## Full Changelog

### Features
- Add OS keychain integration interface
- Integrate keychain with configuration manager
- Add Ollama provider for local AI
- Add end‑to‑end testing script
- Add connection pooling, timeouts, and token estimation
- Add Electron‑builder fallback packaging configuration
- Add GitHub Actions release workflow for tagged releases

### Fixes
- Make Delete idempotent and preserve sentinel error in Get
- Address code quality issues in keychain implementation
- Address code quality issues in hotkey manager
- Prevent app panic on unsupported LLM providers; add mutex to hotkey manager
- Improve safety of integration test script
- Token estimation uses rune count, add missing newlines
- End‑to‑end test script improvements

### Documentation
- Add user guide and installation instructions
- Add release checklist

### Build & CI
- Update wails.json with packaging configuration
- Clean up temporary files and add missing source code

## Known Issues
- Linux builds require `webkit2gtk‑4.1` (Ubuntu 24.04) – symlink provided for compatibility
- End‑to‑end test may fail on systems missing webkit2gtk development packages

## Contributors
This release was developed with the assistance of AI‑powered workflows using the Multi‑Tier Orchestrator tool.

---
**Release Date**: April 2, 2026  
**Version**: 2.0.0  
**License**: MIT