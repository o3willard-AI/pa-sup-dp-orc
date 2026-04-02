# Orchestrator Improvements to Prevent Corrupted JSON

**Date:** 2026-03-31  
**Issue:** L0-Coder corrupted output in Tasks 2.5, 2.6

---

## Root Cause Analysis

### Task 2.6 L0-Coder Failure

**Model Output:**
```json
{"name": "read_file", "parameters": {"path": "internal/terminal/types.go"}}
{"name": "read_file", "parameters": {"path": "internal/terminal/detector.go"}}
{"name": "read_file", "parameters": {"path": "internal/terminal/macos/accessibility_cgo.go"}}
```

**Extracted (Wrong):**
- Path 1: `}}\n{` (corrupted)
- Path 2: `internal/terminal/macos/accessibility_cgo.go` (correct)

**Regex Pattern (orchestrator.py:116):**
```python
claude_read = r'\{[^}]*"name"[^}]*"read_file"[^}]*"parameters"[^}]*\{[^}]*"path"[^}]*"([^"]+)"'
```

**Problem:** `[^}]*` is greedy and spans across JSON object boundaries (`}}\n{`).

---

## Fix Options

### Option 1: Fix Regex Patterns (Quick Fix)

Change `[^}]*` to `[^{}]*` to prevent crossing object boundaries:

```python
# Before (buggy)
claude_read = r'\{[^}]*"name"[^}]*"read_file"[^}]*"parameters"[^}]*\{[^}]*"path"[^}]*"([^"]+)"'

# After (fixed)
claude_read = r'\{[^{}]*"name"[^{}]*"read_file"[^{}]*"parameters"[^{}]*\{[^{}]*"path"[^{}]*"([^"]+)"'
```

**Pros:** Quick, minimal change  
**Cons:** Still regex-based, may have other edge cases

---

### Option 2: Proper JSON Parsing (Recommended)

Parse each line as JSON and validate:

```python
import json

def parse_and_execute_tools(self, response: str) -> Dict[str, Any]:
    results = []
    
    # Try to parse as JSON lines (Claude/OpenRouter format)
    for line in response.strip().split('\n'):
        line = line.strip()
        if not line.startswith('{'):
            continue
        try:
            tool_call = json.loads(line)
            if tool_call.get('name') == 'read_file':
                path = tool_call.get('parameters', {}).get('path')
                if path and isinstance(path, str):
                    result = self.file_tools.file_read(path)
                    results.append({...})
            elif tool_call.get('name') == 'write_file':
                path = tool_call.get('parameters', {}).get('path')
                content = tool_call.get('parameters', {}).get('content')
                if path and content:
                    result = self.file_tools.file_write(path, content)
                    results.append({...})
        except json.JSONDecodeError:
            continue  # Skip invalid JSON lines
    
    # Also check for custom file_write() format
    # ... existing patterns ...
    
    return {"tool_results": results, ...}
```

**Pros:** Robust, handles edge cases, clear error messages  
**Cons:** More code changes

---

### Option 3: Model-Specific Prompt Engineering

Add explicit format instructions to system prompt:

```
TOOL CALL FORMAT (REQUIRED):
You MUST use exactly ONE tool call per response in this format:
file_write("path/to/file", """file content here""")

Do NOT use JSON format. Do NOT use multiple tool calls in one response.
After each file_write, wait for confirmation before proceeding.
```

**Pros:** No code changes, guides model behavior  
**Cons:** May not work consistently across models

---

### Option 4: Response Validation Layer

Add validation before tool execution:

```python
def validate_tool_call(self, tool: str, path: str) -> bool:
    """Validate tool call parameters before execution."""
    if tool == "file_read" or tool == "file_write":
        # Check path looks like a valid file path
        if not path or len(path) > 500:
            return False
        if path.startswith('}') or path.startswith('{'):
            return False
        if '\n' in path or '\r' in path:
            return False
        # Check for valid path characters
        if not re.match(r'^[a-zA-Z0-9_./\-]+$', path):
            return False
    return True
```

**Pros:** Catches corrupted data before execution  
**Cons:** Doesn't fix root cause, adds complexity

---

## Recommended Approach

**Combine Options 1 + 2 + 4:**

1. **Immediate:** Fix regex patterns (Option 1) - prevents current corruption
2. **Short-term:** Add JSON line parsing (Option 2) - proper handling of Claude format
3. **Defense:** Add validation layer (Option 4) - catch any remaining issues

---

## Implementation Priority

| Priority | Fix | Effort | Impact |
|----------|-----|--------|--------|
| P0 | Fix regex `[^}]*` → `[^{}]*` | 10 min | High |
| P1 | Add JSON line parser | 1-2 hours | High |
| P2 | Add validation layer | 30 min | Medium |
| P3 | Improve system prompts | 15 min | Medium |

---

## Expected Improvement

**Current L0-Coder Success Rate:** ~33% (1/3 tasks with usable output)

**Expected After Fixes:** ~80%+ (based on Task 2.4 L1 success pattern)

**Rationale:**
- Most L0 failures are parsing issues, not model capability issues
- Task 2.6 L0 output showed correct file paths in response, just corrupted during parsing
- Proper JSON parsing handles the Claude/OpenRouter native format correctly

---

## Files to Modify

1. `orchestrator/src/core/orchestrator.py` - ToolExecutor.parse_and_execute_tools()
2. `orchestrator/src/core/orchestrator.py` - Add validate_tool_call() method
3. `orchestrator/templates/` - Update system prompts with format examples

---

## Testing

After fixes, re-run failed tasks:
- Task 2.5: L0-Coder should produce valid output
- Task 2.6: L0-Coder should produce valid output

Compare against known-good outputs from manual implementation.
