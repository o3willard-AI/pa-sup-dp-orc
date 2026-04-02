# Windows UI Automation Integration Testing

## Prerequisites
1. Windows 7 or later
2. Windows SDK installed (for UI Automation headers)
3. At least one terminal window open (Windows Terminal, PowerShell, or cmd.exe)
4. CGO enabled: `set CGO_ENABLED=1`

## Running Integration Tests
```bash
# Build with CGO
set CGO_ENABLED=1
go build ./internal/terminal/windows/...

# Run integration tests
go test -v ./internal/terminal/windows/... -tags=integration

# Run unit tests only (mocked)
go test -v ./internal/terminal/windows/... -tags=!integration
```

## Notes
- Integration tests require actual terminal windows to be open
- Tests will skip if UI Automation is not available
- Use `-short` flag to skip integration tests
- Mocked unit tests work on any platform with `windows` build tag