# Windows UI Automation Permissions

## Overview
PairAdmin uses Windows UI Automation API to capture terminal content on Windows. No special permissions are required for standard user accounts.

## Requirements
- **Windows Version:** Windows 7 or later (UI Automation introduced in Windows 7)
- **Terminal Support:**
  - Windows Terminal (full support)
  - PowerShell Console (good support)
  - Command Prompt (cmd.exe) (limited support)
  - PuTTY (basic support)

## UI Automation Initialization
The adapter automatically initializes UI Automation COM API with apartment-threaded model (`COINIT_APARTMENTTHREADED`). No manual configuration needed.

## Troubleshooting

### Common Issues

1. **"UI Automation not initialized" error**
   - Ensure UI Automation COM components are registered (runs as part of Windows)
   - Check COM permissions (standard user should have access)
   - Verify Windows SDK is installed for development

2. **No terminal windows detected**
   - Open a terminal window (Windows Terminal, PowerShell, or cmd.exe)
   - Ensure terminal is not running as administrator if PairAdmin runs as standard user
   - Some terminal emulators may not expose UI Automation interfaces

3. **"Text pattern not available" error**
   - The terminal may not support UI Automation Text pattern
   - Try a different terminal (Windows Terminal has best support)
   - Some older terminals may require accessibility settings

### Development Setup
For building from source:
```cmd
REM Install Windows SDK (if not present)
REM Typically included with Visual Studio Build Tools

REM Verify headers exist
dir "C:\Program Files (x86)\Windows Kits\10\Include\*\um\uiautomation.h"

REM Build with CGO
set CGO_ENABLED=1
go build ./internal/terminal/windows/...
```

### Security Considerations
- UI Automation requires no special permissions for standard users
- Only terminal window content is captured (not other applications)
- No screen capture or keystroke logging
- COM security uses default impersonation level
