# PairAdmin

**AI‑assisted terminal administration**

PairAdmin enables **“Pair Administration”**—a collaboration model where human system administrators work alongside AI agents to manage Linux, Unix, macOS, and Windows systems via terminal interfaces.

## Quick Links

- [User Guide](docs/user-guide.md)
- [Installation Guide](docs/installation.md)
- [Product Requirements (PRD)](docs/planning/PairAdmin_PRD_v2.0.md)
- [Project History & Maintenance Guide](PROJECT_HISTORY_AND_MAINTENANCE.md) – **single source of truth for maintainers**

## Building

```bash
# Install dependencies
./scripts/install-deps.sh

# Development (hot‑reload)
wails dev

# Production build
wails build
```

## Testing

```bash
go test ./internal/... -short          # Unit tests
./scripts/test-integration.sh          # Integration tests
./scripts/test-e2e.sh                  # End‑to‑end tests
```

## License

MIT – see [LICENSE](LICENSE) file.

---

*For detailed architecture, release history, CI/CD, and maintenance instructions, refer to [`PROJECT_HISTORY_AND_MAINTENANCE.md`](PROJECT_HISTORY_AND_MAINTENANCE.md).*