# Linux Terminal Permissions Guide (AT-SPI2)

This guide explains how to configure AT-SPI2 (Assistive Technology Service Provider Interface) permissions for PairAdmin on Linux.

---

## Overview

**AT-SPI2** is the accessibility framework used by GNOME and other Linux desktop environments. PairAdmin uses AT-SPI2 to:

- Detect open terminal windows (GNOME Terminal, Konsole, Terminator)
- Extract text content from terminal windows
- Monitor terminal activity in real-time

---

## Requirements

### System Packages

Install AT-SPI2 development files:

```bash
# Debian/Ubuntu
sudo apt install libatspi2.0-dev at-spi2-core

# Fedora/RHEL
sudo dnf install at-spi2-core-devel

# Arch Linux
sudo pacman -S at-spi2-atk libatspi
```

### Desktop Environment

AT-SPI2 works best with:
- ✅ GNOME (full support)
- ✅ KDE Plasma (good support)
- ✅ MATE (good support)
- ⚠️ XFCE (limited support)
- ❌ TTY/Console (no support - use tmux adapter instead)

---

## Enabling Accessibility

### GNOME

1. **Enable via Settings:**
   - Open **Settings** → **Accessibility**
   - Toggle **Accessibility** to ON

2. **Enable via Command Line:**
   ```bash
   gsettings set org.gnome.desktop.interface accessibility true
   ```

3. **Verify:**
   ```bash
   gsettings get org.gnome.desktop.interface accessibility
   # Should output: true
   ```

### KDE Plasma

1. **Enable via System Settings:**
   - Open **System Settings** → **Accessibility**
   - Check **Enable Accessibility**

2. **AT-SPI2 should start automatically** when a terminal is opened

---

## Checking AT-SPI2 Status

### Is AT-SPI2 Running?

```bash
# Check if at-spi2-daemon is running
ps aux | grep at-spi

# Check D-Bus service
dbus-send --session --print-reply \
  --dest=org.a11y.Bus \
  /org/a11y/bus \
  org.a11y.Bus.GetAddress
```

### Start AT-SPI2 Manually

```bash
# Start the AT-SPI2 daemon
/usr/lib/at-spi2-core/at-spi2-daemon &

# Or on some systems:
at-spi2-daemon &
```

---

## Troubleshooting

### Issue: "No terminals found"

**Possible causes:**
1. AT-SPI2 daemon not running
2. Terminal doesn't support AT-SPI2
3. Running in a sandbox (Flatpak/Snap)

**Solutions:**

```bash
# 1. Start AT-SPI2 daemon
/usr/lib/at-spi2-core/at-spi2-daemon &

# 2. Try a different terminal
# GNOME Terminal is recommended:
sudo apt install gnome-terminal

# 3. If using Flatpak, grant accessibility permission
flatpak override --user --socket=at-spi org.gnome.Terminal
```

### Issue: "Permission denied" or "Access denied"

**Solution:**

```bash
# Check D-Bus permissions
ls -la ~/.cache/at-spi/

# Clear AT-SPI2 cache and restart
rm -rf ~/.cache/at-spi/
killall at-spi2-daemon
/usr/lib/at-spi2-core/at-spi2-daemon &
```

### Issue: "Adapter not available"

**Check if AT-SPI2 is accessible:**

```bash
# Test AT-SPI2 connection
accerciser &

# If accerciser doesn't start, AT-SPI2 is not properly configured
```

**Install accerciser for debugging:**

```bash
# Debian/Ubuntu
sudo apt install accerciser

# Fedora
sudo dnf install accerciser
```

### Issue: Flatpak/Snap Permissions

**Flatpak:**

```bash
# Grant AT-SPI2 access
flatpak override --user --socket=at-spi <app-id>

# For PairAdmin specifically (if installed as Flatpak)
flatpak override --user --socket=at-spi com.pairadmin.PairAdmin
```

**Snap:**

```bash
# Grant accessibility interface
sudo snap connect pairadmin:accessibility :accessibility
```

---

## Supported Terminals

| Terminal | AT-SPI2 Support | Notes |
|----------|----------------|-------|
| GNOME Terminal | ✅ Full | Recommended |
| Konsole | ✅ Full | KDE default |
| Terminator | ✅ Good | |
| Tilix | ✅ Good | |
| xterm | ⚠️ Limited | May not expose text |
| st | ❌ None | Use tmux adapter |
| Alacritty | ❌ None | Use tmux adapter |

---

## Verification Steps

### 1. Test AT-SPI2 Availability

```bash
# Simple test - should return address
dbus-send --session --print-reply \
  --dest=org.a11y.Bus \
  /org/a11y/bus \
  org.a11y.Bus.GetAddress
```

Expected output:
```
method return time=... sender=... -> destination=... serial=... reply_serial=...
   string "unix:abstract=/tmp/dbus-..."
```

### 2. Test Terminal Detection

```bash
# Using PairAdmin CLI
pairadmin terminals list

# Should show detected terminals like:
# Found 2 terminals:
#   - linux-12345-user@host:~ (GNOME Terminal)
#   - linux-67890-user@host:~/projects (Konsole)
```

### 3. Test Content Capture

```bash
# Capture from a specific terminal
pairadmin terminals capture --id linux-12345-user@host:~

# Should return terminal content
```

---

## Advanced Configuration

### Environment Variables

```bash
# Force AT-SPI2 to use specific D-Bus session
export AT_SPI_BUS_ADDRESS="unix:abstract=/tmp/dbus-..."

# Enable AT-SPI2 debugging
export NO_AT_BRIDGE=1
```

### Systemd User Service

Create `~/.config/systemd/user/at-spi2.service`:

```ini
[Unit]
Description=AT-SPI2 Accessibility Service

[Service]
Type=simple
ExecStart=/usr/lib/at-spi2-core/at-spi2-daemon
Restart=on-failure

[Install]
WantedBy=default.target
```

Enable and start:

```bash
systemctl --user enable at-spi2
systemctl --user start at-spi2
```

---

## Getting Help

If issues persist:

1. **Check logs:**
   ```bash
   journalctl --user -u at-spi2
   ```

2. **Report bugs:**
   - AT-SPI2: https://gitlab.gnome.org/GNOME/at-spi2-core
   - PairAdmin: https://github.com/pairadmin/pairadmin/issues

3. **Community support:**
   - GNOME Accessibility: https://wiki.gnome.org/Accessibility

---

**Last Updated:** 2026-03-31  
**Applicable Version:** PairAdmin 2.0
