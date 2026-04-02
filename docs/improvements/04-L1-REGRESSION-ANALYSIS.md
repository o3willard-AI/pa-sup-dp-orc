# L1-Coder Corruption Regression Analysis

**Date:** 2026-03-31  
**Issue:** L1-Coder producing corrupted files in Task 2.8 (worked in Task 2.7)

---

## Root Cause Identified

### The JSON Line Parser Bug (P1)

**Location:** `orchestrator/src/core/orchestrator.py:121-166`

**Problem:** The JSON line parser splits response by newlines:
```python
for line in response.strip().split('\n'):
    line = line.strip()
    if not line.startswith('{'):
        continue
    try:
        tool_call = json.loads(line)  # FAILS for multi-line JSON
```

**When Model Outputs Multi-line JSON:**
```json
{"name": "write_file", "parameters": {"path": "file.h", "content": "#ifndef FOO
#define FOO
#endif
"}}
```

**After split('\n'):**
- Line 1: `{"name": "write_file", "parameters": {"path": "file.h", "content": "#ifndef FOO` ❌ Invalid JSON
- Line 2: `#define FOO` ❌ Doesn't start with `{`
- Line 3: `#endif` ❌ Doesn't start with `{`
- Line 4: `"}}` ❌ Doesn't start with `{`

**Result:** JSON line parser skips all lines, falls through to regex patterns.

---

### The Claude Native Regex Bug (Line 197)

**Pattern:**
```python
claude_write = r'\{[^{}]*"name"[^{}]*"write_file"[^{}]*"parameters"[^{}]*\{[^{}]*"path"[^{}]*"([^"]+)"[^{}]*"content"[^{}]*"((?:[^"\\]|\\.)*)"'
```

**Content Capture Group:** `((?:[^"\\]|\\.)*)`

**Problem:** This pattern:
1. Matches `[^"\\]` - any char except quote or backslash
2. OR `\\.` - backslash followed by any char (escape sequence)
3. **STOPS** at first unescaped quote

**When Content Has Escaped Quotes:**
```json
{"name": "write_file", "parameters": {"path": "file.h", "content": "// Header \"guard\"\n#ifndef FOO"}}
```

The pattern may truncate at unexpected points due to complex interaction with `[^{}]*` wildcards.

---

## Why Task 2.7 Worked

**Task 2.7 Output:** 4 files created successfully

**Likely Reason:** The model output shorter file contents that fit within regex matching limits, or the content didn't have characters that triggered the edge cases.

---

## Why Task 2.8 Failed

**Task 2.8 Output:** 3 files (1 corrupted - 5 bytes, 2 missing)

**atspi2.h Issue:** Only 5 bytes written (`"}}\n{`)

**Root Cause:** 
1. Model output multi-line JSON content
2. JSON line parser failed (splits on newlines)
3. Claude regex matched incorrectly, capturing only `"}}\n{` as content
4. file_write succeeded but with garbage content

---

## Evidence

### Task 2.8 Handoff Analysis
```
Handoff file size: 2463 bytes
output_preview length: 503 chars (truncated for logging)
Tools reported: 3 file_writes, all "succeeded"
Actual file sizes: atspi2.h = 5 bytes (CORRUPTED)
```

### Model Response Truncation
The handoff file is only 2463 bytes total, suggesting the model response itself may have been truncated by:
1. OpenRouter API response limit
2. Model hitting max_tokens (set to 8192)
3. Model stopping early due to complexity

---

## Additional Bug: Duplicate Method Definition

**Location:** Lines 393-432 and 520-559

**Issue:** `execute_task()` method is defined TWICE in the class.

**Impact:** Python uses the second definition (line 520), but this indicates code quality issues from previous edits. LSP correctly warns: "Function declaration is obscured by a declaration of the same name"

---

## Fixes Required

### P8: Remove Broken JSON Line Parser

**Rationale:** The regex patterns (with P0 fix) handle Claude native format correctly. The JSON line parser adds complexity but fails on multi-line content.

**Change:** Remove lines 121-166, keep regex patterns.

### P9: Fix Claude Native Content Regex

**Current:**
```python
claude_write = r'\{[^{}]*"name"[^{}]*"write_file"[^{}]*"parameters"[^{}]*\{[^{}]*"path"[^{}]*"([^"]+)"[^{}]*"content"[^{}]*"((?:[^"\\]|\\.)*)"'
```

**Problem:** Content pattern doesn't handle all escape sequences properly.

**Fix:** Use non-greedy match until closing `"}` of parameters:
```python
claude_write = r'\{"name":"write_file","parameters":\{"path":"([^"]+)","content":"((?:[^"\\]|\\.)+)"\}\}'
```

Or better - use JSON parsing for Claude format:
```python
# Find JSON objects in response
for match in re.finditer(r'\{[^{}]*"name"[^{}]*\}', response):
    try:
        obj = json.loads(match.group())
        if obj.get('name') == 'write_file':
            # Extract from obj['parameters']
    except json.JSONDecodeError:
        pass
```

### P10: Remove Duplicate Method

Delete lines 393-432 (first `execute_task`), keep lines 520-559.

### P11: Add Response Length Logging

Log actual response length to detect truncation:
```python
print(f"  LLM response length: {len(response)} chars")
```

---

## Testing Plan

1. **Remove JSON line parser** (P8)
2. **Test with Task 2.8** - re-run L1-Coder
3. **Verify atspi2.h** - should have full content (~800 bytes)
4. **Test with new task** - verify no regression

---

## Expected Outcome

**After Fixes:**
- L1-Coder success rate: 90%+ (back to Task 2.7 levels)
- No corrupted files from regex parsing
- Cleaner codebase (no duplicate methods)

---

**Priority:** CRITICAL - Blocker for continuing tasks
