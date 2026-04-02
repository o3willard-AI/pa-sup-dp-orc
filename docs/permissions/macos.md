# macOS Terminal Permissions Guide (Accessibility API)

This guide explains how to configure Accessibility permissions for PairAdmin on macOS.

---

## Overview

**macOS Accessibility API** allows applications to interact with UI elements of other applications. PairAdmin uses the Accessibility API to:

- Detect open Terminal.app windows
- Extract text content from terminal windows
- Monitor terminal activity in real-time

---

## Requirements

### macOS Version

Accessibility API is built into macOS and requires:

- ✅ macOS 10.15 (Catalina) or later
- ✅ macOS 11 (Big Sur)
- ✅ macOS 12 (Monterey)
- ✅ macOS 13 (Ventura)
- ✅ macOS 14 (Sonoma)

### Administrator Privileges

You need administrator access to grant Accessibility permissions.

---

## Enabling Accessibility Permissions

### First Launch

When you first launch PairAdmin, macOS will display a permission dialog:

1. Click **Open System Preferences** when prompted
2. Go to **Security & Privacy** → **Privacy** → **Accessibility**
3. Click the **lock icon** 🔒 and enter your password
4. Find **PairAdmin** in the list and check the box
5. Restart PairAdmin

### Manual Configuration

If you dismissed the dialog or need to re-enable:

#### macOS 11+ (Big Sur and later)

1. Click **Apple menu** → **System Preferences**
2. Click **Security & Privacy**
3. Select **Privacy** tab
4. Select **Accessibility** in the left sidebar
5. Click the **lock icon** 🔒 and enter your password
6. Click **+** to add PairAdmin (if not listed)
7. Navigate to `/Applications/PairAdmin.app` and add it
8. Check the box next to PairAdmin

#### macOS 10.15 (Catalina)

1. Click **Apple menu** → **System Preferences**
2. Click **Security & Privacy**
3. Select **Privacy** tab
4. Select **Accessibility** in the left sidebar
5. Click the **lock icon** 🔒 and enter your password
6. Check the box next to PairAdmin

### Terminal Command

You can also enable via Terminal:

```bash
# Grant accessibility permission
sudo tccutil reset Accessibility com.pairadmin.PairAdmin

# Or manually add to accessibility list
defaults write com.apple.universalaccess AXManualAccessibility -int 1
```

---

## Verifying Permissions

### Check Permission Status

```bash
# Check if PairAdmin has accessibility access
osascript -e 'tell application "System Events" to get assistive enabled'
# Should return: true
```

### Test Terminal Detection

```bash
# Using PairAdmin CLI
pairadmin terminals list

# Expected output:
# Found 2 terminals:
#   - macos-12345-Administrator:~ (Terminal)
#   - macos-67890-user@host:~ (iTerm2)
```

---

## Supported Terminals

| Terminal | Accessibility Support | Notes |
|----------|----------------------|-------|
| **Terminal.app** | ✅ Full | Built-in, recommended |
| **iTerm2** | ✅ Full | Excellent support |
| **Kitty** | ⚠️ Limited | May require configuration |
| **Alacritty** | ❌ None | Use tmux adapter |
| **WezTerm** | ⚠️ Limited | Experimental support |

---

## Troubleshooting

### Issue: "Permission denied" or "Not authorized"

**Solution 1: Re-grant Permission**

1. Open **System Preferences** → **Security & Privacy** → **Privacy** → **Accessibility**
2. Uncheck PairAdmin
3. Wait 5 seconds
4. Check PairAdmin again
5. Restart PairAdmin

**Solution 2: Reset TCC Database**

```bash
# Reset accessibility permissions for PairAdmin
tccutil reset Accessibility com.pairadmin.PairAdmin

# Then re-launch PairAdmin and grant permission when prompted
```

### Issue: "No terminals found"

**Possible causes:**
1. No supported terminals are open
2. Accessibility permission not granted
3. Terminal doesn't support Accessibility API

**Solutions:**

```bash
# 1. Open Terminal.app (built-in)
open -a Terminal

# 2. Verify permission
osascript -e 'tell application "System Events" to get assistive enabled'

# 3. Try iTerm2 if Terminal.app doesn't work
open -a iTerm
```

### Issue: "Adapter not available"

**Check System Events:**

```bash
# Ensure System Events is running
ps aux | grep "System Events"

# If not running, it should auto-start when needed
```

### Issue: Permission Dialog Doesn't Appear

**Force Reset:**

```bash
# Kill PairAdmin
killall PairAdmin

# Reset TCC
tccutil reset Accessibility com.pairadmin.PairAdmin

# Re-launch
open -a PairAdmin
```

**Manual Add:**

1. Open **System Preferences** → **Security & Privacy** → **Privacy** → **Accessibility**
2. Click **+** button
3. Navigate to `/Applications/PairAdmin.app`
4. Add and check the box

### Issue: iTerm2 Not Detected

**Configure iTerm2:**

1. Open iTerm2
2. Go to **iTerm2** → **Preferences** → **General**
3. Check **Enable Apple Scripting**
4. Restart iTerm2

---

## Verification Steps

### 1. Test Accessibility Availability

**AppleScript Test:**

```bash
# Test if accessibility is enabled
osascript -e 'tell application "System Events" to return assistive enabled'
# Should return: true
```

### 2. Test Terminal Detection

**List Windows:**

```bash
# Using AppleScript
osascript -e 'tell application "Terminal" to count of windows'
# Should return number of open Terminal windows
```

### 3. Test Content Capture

```bash
# Using PairAdmin CLI
pairadmin terminals capture --id macos-12345-Administrator:~

# Should return terminal content
```

---

## Advanced Configuration

### Privacy Preferences Control (TCC)

View current accessibility permissions:

```bash
# List all apps with accessibility access
sqlite3 ~/Library/Application\ Support/com.apple.TCC/TCC.db \
  "SELECT * FROM access WHERE service='kTCCServiceAccessibility';"
```

### Reset All Accessibility Permissions

```bash
# Warning: This resets ALL apps, not just PairAdmin
tccutil reset Accessibility
```

### Enterprise Deployment

For MDM deployment (Jamf, Kandji, etc.):

```xml
<!-- Accessibility payload configuration -->
<key>Accessibility</key>
<array>
    <dict>
        <key>Identifier</key>
        <string>com.pairadmin.PairAdmin</string>
        <key>CodeRequirement</key>
        <string>anchor apple generic and identifier "com.pairadmin.PairAdmin"</string>
    </dict>
</array>
```

---

## Security Considerations

### What PairAdmin Can Access

- ✅ Terminal window titles
- ✅ Terminal text content
- ✅ Terminal process information
- ✅ Window position and size

### What PairAdmin Cannot Access

- ❌ Other application windows (unless explicitly selected)
- ❌ System clipboard (without permission)
- ❌ Network activity
- ❌ File system (without permission)
- ❌ Keystrokes (PairAdmin does not log input)

### Data Privacy

- Terminal content is processed locally
- No data is sent to external servers without explicit configuration
- Audit logs are stored in `~/.pairadmin/logs/`

### Revoking Access

To revoke accessibility access:

1. Open **System Preferences** → **Security & Privacy** → **Privacy** → **Accessibility**
2. Uncheck PairAdmin
3. Or remove PairAdmin from the list with **-** button

---

## Getting Help

If issues persist:

1. **Check Console logs:**
   ```bash
   log show --predicate 'process == "PairAdmin"' --last 1h
   ```

2. **Enable PairAdmin Debug Logging:**
   ```bash
   pairadmin --debug
   ```

3. **Report bugs:**
   - PairAdmin: https://github.com/pairadmin/pairadmin/issues

4. **Apple Documentation:**
   - Accessibility: https://developer.apple.com/documentation/applescript
   - TCC: https://developer.apple.com/documentation/security/accessing_protected_data

---

**Last Updated:** 2026-03-31  
**Applicable Version:** PairAdmin 2.0
