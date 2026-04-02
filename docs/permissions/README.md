# PairAdmin Terminal Permissions Guide

This guide explains the accessibility permissions required for PairAdmin to capture terminal content on different platforms.

---

## Platform Support

- **Linux:** AT-SPI2 accessibility API ([Linux Guide](./linux.md))
- **macOS:** Accessibility API ([macOS Guide](./macos.md))
- **Windows:** UI Automation API ([Windows Guide](./windows.md))
- **tmux:** Cross-platform via `tmux` command

---

## Quick Start

### Linux (GNOME/KDE)
```bash
# Enable accessibility
gsettings set org.gnome.desktop.interface accessibility true

# Verify AT-SPI2 is running
ps aux | grep at-spi

# Test PairAdmin terminal detection
pairadmin detect
```

### macOS
```bash
# Open System Preferences → Security & Privacy → Privacy → Accessibility
# Add PairAdmin to the list

# Verify permission
defaults read com.apple.universalaccess assistiveAccessAllowed
```

### Windows
```powershell
# Run as Administrator if experiencing issues
# Test UI Automation availability
pairadmin detect
```

---

## Detailed Guides

- **[Linux Permissions Guide](linux.md)** - AT-SPI2 setup and troubleshooting
- **[Windows Permissions Guide](windows.md)** - UI Automation setup and troubleshooting
- **[macOS Permissions Guide](macos.md)** - Accessibility API setup and troubleshooting

---

## Common Issues

| Issue | Linux | macOS | Windows |
|-------|-------|-------|---------|
| "Permission denied" | Enable AT-SPI2 | Add to Accessibility | Run as Admin |
| "No terminals found" | Check daemon | Restart Terminal.app | Check UAC |
| "Adapter not available" | Install libatspi2 | Grant permission | Enable UIA |
| Flatpak/Sandbox issues | Use native package | N/A | N/A |

---

## Testing Your Setup

After configuring permissions, verify with:

```bash
# List detected terminals
pairadmin terminals list

# Capture from first terminal
pairadmin terminals capture --id <terminal-id>
```

---

## Getting Help

If you're still experiencing issues:

1. Check the platform-specific guide for detailed troubleshooting
2. Review the [FAQ](../faq.md)
3. Open an issue on [GitHub](https://github.com/pairadmin/pairadmin/issues)

---

**Last Updated:** 2026-03-31  
**Version:** 2.0
