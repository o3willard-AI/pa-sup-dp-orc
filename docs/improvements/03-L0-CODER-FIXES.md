# L0-Coder Parsing Issues Analysis & Fixes

**Date:** 2026-03-31  
**Status:** Additional improvements needed

---

## Observed Failure Patterns

### Pattern 1: Escaped Newlines in File Content (Task 2.7)

**L0 Output:**
```python
file_write("detector.go", """package terminal\n\nimport (\n\t"context"\n)""")
```

**Expected:**
```python
file_write("detector.go", """package terminal

import (
	"context"
)""")
```

**Root Cause:** L0-Coder (Qwen3-Coder 30B) outputs escaped newlines (`\n`) instead of actual newlines in triple-quoted strings.

**Impact:** Files created have all content on single line, won't compile.

---

### Pattern 2: JSON Object Boundary Corruption (Task 2.6)

**L0 Output:**
```json
{"name": "read_file", "parameters": {"path": "file1.go"}}
{"name": "read_file", "parameters": {"path": "file2.go"}}
{"name": "read_file", "parameters": {"path": "file3.go"}}
```

**Parsed (Before P0 Fix):**
- Path 1: `}}\n{` (corrupted - regex crossed boundary)
- Path 2: `file3.go` (correct)

**Root Cause:** Regex `[^}]*` was greedy across JSON objects.

**Current Status:** ✅ Fixed by P0 (changed to `[^{}]*`)

---

### Pattern 3: Context Size Exceeded (Task 2.5)

**L0 Error:** `400 - {"error":"Context size has been exceeded."}`

**Root Cause:** Spec file + conversation history exceeded LM Studio context window.

**Current Status:** ⚠️ Partially addressed by context simplification, but still fails on large specs.

---

## Additional Improvements (P4-P7)

### P4: Post-Processing for Escaped Newlines

**Location:** `orchestrator/src/core/orchestrator.py` - `FileTools.file_write()`

**Fix:** Detect and convert escaped newlines before writing:

```python
def file_write(self, path: str, content: str) -> Dict[str, Any]:
    """Write content to a file with post-processing."""
    try:
        # P4: Fix escaped newlines from L0-Coder
        content = self._fix_escaped_newlines(content)
        
        file_path = Path(path)
        if not file_path.is_absolute():
            file_path = self.project_root / file_path
        file_path.parent.mkdir(parents=True, exist_ok=True)
        file_path.write_text(content)
        return {"success": True, "path": str(file_path), "bytes": len(content)}
    except Exception as e:
        return {"success": False, "error": str(e)}

def _fix_escaped_newlines(self, content: str) -> str:
    """Convert escaped newlines to actual newlines."""
    # Detect if content has escaped newlines (common L0-Coder issue)
    if '\\n' in content and '\n' not in content:
        # Likely escaped - convert to actual newlines
        content = content.replace('\\n', '\n')
        content = content.replace('\\t', '\t')
        content = content.replace('\\"', '"')
        content = content.replace('\\\\', '\\')
    return content
```

**Expected Impact:** L0-Coder file writes become usable (80%+ success rate)

---

### P5: Smarter Context Simplification

**Location:** `orchestrator/src/core/orchestrator.py` - `_simplify_context()`

**Current:** Simple multiplier-based truncation

**Improved:** Priority-based content preservation:

```python
def _simplify_context(self, user_prompt: str, multiplier: float) -> str:
    """Simplify context while preserving critical information."""
    # Keep instructions and file paths (critical)
    # Truncate file content samples (non-critical)
    
    lines = user_prompt.split('\n')
    preserved = []
    truncatable = []
    
    for line in lines:
        if any(kw in line.lower() for kw in ['instruction', 'must', 'create', 'file:', 'path:']):
            preserved.append(line)
        elif line.startswith('```') or line.strip().startswith('//') or line.strip().startswith('#'):
            truncatable.append(line)  # Code comments can be truncated
        else:
            truncatable.append(line)
    
    # Truncate truncatable content
    target_lines = int(len(truncatable) * multiplier)
    truncatable = truncatable[:target_lines]
    
    return '\n'.join(preserved + truncatable)
```

**Expected Impact:** L0-Coder stays within context limits while keeping instructions

---

### P6: Explicit Tool Format in System Prompt

**Location:** `orchestrator/templates/02-l0-coder.md`

**Add:**

```markdown
## CRITICAL: FILE CONTENT FORMAT

When using file_write(), content MUST have ACTUAL newlines:

✅ CORRECT:
file_write("main.go", """package main

func main() {
    fmt.Println("hello")
}
""")

❌ WRONG (will create broken files):
file_write("main.go", """package main\n\nfunc main() {\n    fmt.Println("hello")\n}
""")

**RULE:** Press Enter for newlines in your response. Do NOT type \n characters.
```

**Expected Impact:** L0-Coder understands format requirement

---

### P7: Early Escalation Detection

**Location:** `orchestrator/src/core/orchestrator.py` - `execute_task()`

**Add:**

```python
def _should_escalate_early(self, response: str, tool_results: Dict) -> bool:
    """Detect patterns indicating L0 will fail, escalate early."""
    
    # Pattern 1: Multiple context size errors
    if 'context' in response.lower() and 'exceed' in response.lower():
        return True
    
    # Pattern 2: Tool calls but all failed
    if tool_results.get('tools_executed', 0) > 0:
        success_rate = sum(1 for r in tool_results.get('tool_results', []) if r.get('success'))
        if success_rate == 0:
            return True
    
    # Pattern 3: Response mentions inability
    if any(kw in response.lower() for kw in ['cannot', 'unable to', 'sorry', 'apologize']):
        if 'implement' in response.lower() or 'create' in response.lower():
            return True
    
    return False
```

**Expected Impact:** Faster escalation, less wasted attempts

---

## Implementation Priority

| Priority | Improvement | Effort | Expected Impact |
|----------|-------------|--------|-----------------|
| P4 | Escaped newline fix | 15 min | HIGH (fixes Task 2.7 issue) |
| P6 | System prompt update | 10 min | MEDIUM (prevents future issues) |
| P7 | Early escalation | 30 min | MEDIUM (saves time) |
| P5 | Context simplification | 45 min | LOW (complex, marginal gain) |

---

## Testing Plan

After implementing P4-P7:

1. **Re-run Task 2.7** with L0-Coder - verify escaped newlines are fixed
2. **Run new Task 2.8** with L0-Coder - verify full cascade works
3. **Measure success rate** over 5 tasks - target 60%+ L0 success

---

## Expected L0-Coder Success Rate

| Issue | Current | After P4-P7 |
|-------|---------|-------------|
| Escaped newlines | 100% fail | 90%+ pass |
| Context exceeded | 100% fail | 50% pass (early escalation) |
| JSON corruption | 40% fail | 5% fail (P0 fix) |
| **Overall** | ~30% | ~70% |

---

**Recommendation:** Implement P4 and P6 immediately (25 min total), then test.
