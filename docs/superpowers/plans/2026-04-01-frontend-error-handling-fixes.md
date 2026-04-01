# Frontend Error Handling Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix critical and important issues identified in code quality review for PairAdmin frontend components.

**Architecture:** Fixes involve modifying `fetchCommands` utility to propagate errors, adding request tracking to prevent race conditions, adding reactive cleanup, and adding user-facing error feedback.

**Tech Stack:** Svelte, JavaScript, Wails Go bindings.

---

### Task 1: Remove error‑swallowing try‑catch from fetchCommands

**Files:**
- Modify: `frontend/src/lib/commands.js:6-12`

- [ ] **Step 1: Remove try‑catch block**

```javascript
// Current content lines 6‑12:
//    try {
//        const commands = await GetCommandsByTerminal(terminalId);
//        commandHistory.set(commands);
//    } catch (error) {
//        console.error('Failed to fetch commands:', error);
//    }

// Replace with:
    const commands = await GetCommandsByTerminal(terminalId);
    commandHistory.set(commands);
```

- [ ] **Step 2: Verify file after edit**

Run: `cat frontend/src/lib/commands.js`
Expected: No try‑catch block, function returns promise that propagates errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/lib/commands.js
git commit -m "fix: remove error‑swallowing try‑catch from fetchCommands (critical)"
```

---

### Task 2: Implement request‑ID counter to prevent race condition

**Files:**
- Modify: `frontend/src/components/CommandSidebar.svelte:6-18`

- [ ] **Step 1: Add request‑ID variable and update loadCommands**

```javascript
// After line 6 (let loadingCommands = false;), add:
  let currentRequestId = 0;

// Replace loadCommands function (lines 8‑18) with:
  async function loadCommands(terminalId) {
    if (!terminalId) return;
    const requestId = ++currentRequestId;
    loadingCommands = true;
    try {
      await fetchCommands(terminalId);
    } catch (error) {
      console.error('Failed to fetch commands:', error);
    } finally {
      // Only update loading state if this request is still the latest
      if (requestId === currentRequestId) {
        loadingCommands = false;
      }
    }
  }
```

- [ ] **Step 2: Verify the changes**

Run: `grep -n "currentRequestId\|loadCommands" frontend/src/components/CommandSidebar.svelte`
Expected: See the variable and updated function.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/CommandSidebar.svelte
git commit -m "fix: add request‑ID tracking to prevent loading‑state race (important)"
```

---

### Task 3: Clear command history when no terminal selected

**Files:**
- Modify: `frontend/src/components/CommandSidebar.svelte:31-34`

- [ ] **Step 1: Add reactive statement to clear commandHistory**

```javascript
// After line 33 (the reactive loadCommands call), add:
  $: if (!$activeTerminalId) commandHistory.set([]);
```

- [ ] **Step 2: Verify the reactive statement**

Run: `tail -n 10 frontend/src/components/CommandSidebar.svelte`
Expected: See the new line after the existing reactive statement.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/CommandSidebar.svelte
git commit -m "fix: clear command history when terminal deselected (minor)"
```

---

### Task 4: Add user‑facing error feedback for command‑fetch failures

**Files:**
- Modify: `frontend/src/components/CommandSidebar.svelte`

- [ ] **Step 1: Add error state variable**

```javascript
// After line 6 (let loadingCommands = false;), add:
  let error = '';

// After currentRequestId line (from Task 2), add helper to clear error:
  function clearError() {
    error = '';
  }
```

- [ ] **Step 2: Update loadCommands to set and auto‑clear error**

```javascript
// Inside loadCommands catch block, replace console.error line with:
      error = 'Failed to load commands. Please try again.';
      console.error('Failed to fetch commands:', error);
      setTimeout(clearError, 5000);
```

- [ ] **Step 3: Add error display in template**

```html
<!-- In template, after the <h3> line (around line 37), add: -->
  {#if error}
    <div class="error">{error}</div>
  {/if}
```

- [ ] **Step 4: Add CSS for error styling**

```css
<!-- In <style> block (anywhere), add: -->
  .error {
    background: #ffebee;
    color: #c62828;
    padding: 0.5rem;
    border-radius: 0.25rem;
    margin-bottom: 1rem;
    font-size: 0.9rem;
  }
```

- [ ] **Step 5: Verify all changes**

Run: `grep -n "error\|clearError" frontend/src/components/CommandSidebar.svelte`
Expected: See variable, function, usage in catch, and template display.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/CommandSidebar.svelte
git commit -m "feat: add transient error feedback for command‑fetch failures (minor)"
```

---

### Task 5: Build verification

**Files:**
- Build frontend to ensure no syntax errors.

- [ ] **Step 1: Run frontend build**

```bash
cd frontend && npm run build
```

Expected: Build completes without errors (exit code 0).

- [ ] **Step 2: Return to root directory**

```bash
cd ..
```

- [ ] **Step 3: Commit any build artifacts if needed**

```bash
git status
# If there are any auto‑generated files (e.g., dist/), add and commit them.
```

- [ ] **Step 4: Final commit summary**

```bash
git log --oneline -5
```
Expected: See the four fix commits.

---

### Task 6: Self‑review

**Files:**
- Review all modified files.

- [ ] **Step 1: Check diff for each file**

```bash
git diff HEAD~4 HEAD -- frontend/src/lib/commands.js
git diff HEAD~4 HEAD -- frontend/src/components/CommandSidebar.svelte
```

- [ ] **Step 2: Verify each issue is addressed**

1. Critical: fetchCommands no longer catches errors.
2. Important: loadCommands uses request‑ID counter.
3. Minor: commandHistory cleared when terminal deselected.
4. Minor: error state displayed and auto‑cleared.

- [ ] **Step 3: Run a quick lint/type check if available**

```bash
cd frontend && npm run lint 2>/dev/null || echo "No lint script"
```

No lint script is fine; just ensure build passed.

- [ ] **Step 4: Report status**

Output: DONE, DONE_WITH_CONCERNS, BLOCKED, or NEEDS_CONTEXT based on results.