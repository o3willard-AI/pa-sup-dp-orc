# Tool Usage Guide for Multi-Tier Workflow

**Purpose:** Ensure all LLM subagents properly use file_read and file_write tools.

---

## The Problem

LLMs can describe what files should contain but cannot directly create files. Our workflow provides tool functions, but the LLM must be explicitly instructed to call them.

**Common Failure Mode:**
```
LLM Response: "Create a file called config.go with this content:
package config
..."
```
❌ This does NOT create a file - it's just text description.

**Correct Usage:**
```
LLM Response: "Creating the file:
file_write("internal/config/config.go", "package config\n...")
```
✅ This calls the file_write tool which creates the actual file.

---

## Tool Syntax

### file_write(path, content)

**Purpose:** Create or modify a file

**Syntax:**
```
file_write("relative/or/absolute/path.txt", "file content here")
```

**Examples:**
```python
# Create a Go file
file_write("internal/config/config.go", "package config\n\nfunc New() *Config {}")

# Create a shell script
file_write("scripts/install-deps.sh", "#!/bin/bash\necho 'Installing...'")

# Create .gitignore
file_write(".gitignore", "*.exe\n*.dll\nnode_modules/\n")
```

**Return Value:**
```json
{
  "success": true,
  "path": "internal/config/config.go",
  "bytes": 1234
}
```

### file_read(path)

**Purpose:** Read an existing file

**Syntax:**
```
file_read("relative/or/absolute/path.txt")
```

**Examples:**
```python
# Read go.mod
file_read("go.mod")

# Read existing config
file_read("internal/config/config.go")
```

**Return Value:**
```json
{
  "success": true,
  "path": "go.mod",
  "content": "module github.com/pairadmin/pairadmin\n\ngo 1.21",
  "lines": 5
}
```

---

## Coder Tier Requirements

**ALL coder tiers (L0, L1, L2, L3 Coder) MUST:**

1. **Use file_write() for every file** - No descriptions, only tool calls
2. **Verify each file_write succeeded** - Check return value
3. **Report tool execution results** in submission format

**Submission Format:**
```
IMPLEMENTATION COMPLETE

Task: 1.4 - Create install-deps script

Files Created (via file_write):
- scripts/install-deps.sh (SUCCESS - 1371 bytes)

Tool Execution Results:
- file_write("scripts/install-deps.sh"): SUCCESS

All file_write calls succeeded: YES
Ready for Review: YES
```

---

## Reviewer Tier Requirements

**ALL reviewer tiers (L0, L3 Architect) MUST:**

1. **Verify files exist** - Don't trust claims, check actual files
2. **Use file_read() to compare** - Read actual file content
3. **Reject if files missing** - No file_write = not implemented

**Review Checklist:**
- [ ] Each file in spec has corresponding file_write() call
- [ ] Each file_write() returned success
- [ ] Files actually exist at specified paths
- [ ] File contents match spec requirements

---

## Planner Tier Requirements

**Planner tiers (L0 Planner, L3 Architect) SHOULD:**

1. **Use file_read() for context** - Read existing files instead of asking
2. **Reference existing patterns** - Read similar files for consistency
3. **Specify exact file paths** - Make it clear where files should be created

**Example:**
```
Instead of: "Create a config file somewhere"
Use: "Create: internal/config/config.go"
```

---

## Common Mistakes

### Mistake 1: Describing Instead of Creating

❌ **Wrong:**
```
"I'll create a .gitignore file with the following patterns:
*.exe
*.dll
node_modules/
"
```

✅ **Correct:**
```
file_write(".gitignore", "*.exe\n*.dll\nnode_modules/\n")
```

### Mistake 2: Claiming Success Without Tool Call

❌ **Wrong:**
```
IMPLEMENTATION COMPLETE
Files Created:
- scripts/install-deps.sh
```
(File doesn't actually exist)

✅ **Correct:**
```
IMPLEMENTATION COMPLETE
Tool Execution:
- file_write("scripts/install-deps.sh"): SUCCESS (1371 bytes)
Files Created:
- scripts/install-deps.sh (verified)
```

### Mistake 3: Wrong Syntax

❌ **Wrong:**
```
file_write(scripts/install-deps.sh, #!/bin/bash...)
```
(Missing quotes)

✅ **Correct:**
```
file_write("scripts/install-deps.sh", "#!/bin/bash\n...")
```

---

## Troubleshooting

### "file_write not recognized"
- Ensure you're using exact syntax: `file_write("path", "content")`
- Include quotes around both path and content
- Content with newlines should use `\n`

### "file_write failed"
- Check if parent directory exists (create with mkdir if needed)
- Check file path is valid (no special characters)
- Check disk space

### "File not found after file_write"
- Verify file_write returned success
- Check the exact path (relative vs absolute)
- Reviewer should use file_read() to verify

---

## Process for Task Execution

### For Coders:

1. **Read spec** → Identify files to create
2. **For each file:** Call `file_write("path", "content")`
3. **Verify** → Each file_write returned success
4. **Submit** → Include tool execution results

### For Reviewers:

1. **Read submission** → Check tool execution results
2. **Verify files** → Use file_read() on each claimed file
3. **Compare to spec** → Ensure all requirements met
4. **Decision** → ACCEPT if all files exist and match spec

---

## Enforcement

**Reviewer Authorization:**

Reviewers are authorized to **REJECT immediately** if:
- No file_write() calls in coder submission
- file_write() calls but files don't exist
- Files exist but content doesn't match spec

**Escalation Path:**

If L0-Coder fails to use file_write after 3 retries:
1. L0-Reviewer rejects with specific feedback
2. Escalate to L1-Coder with note: "Must use file_write tool"
3. L1-Coder implements with proper tool usage

---

## Examples

### Good Coder Submission

```
Task: 1.3 - Set up .gitignore

Files Created (via file_write):
- .gitignore (SUCCESS - 847 bytes)

Tool Execution:
file_write(".gitignore", "# Wails defaults\n/build/bin/\nfrontend/node_modules/\n...")
  → SUCCESS: 847 bytes written

All tools executed successfully: YES
Ready for Review: YES
```

### Good Reviewer Response

```
REVIEW: ACCEPT

Verification:
- file_read(".gitignore"): EXISTS (847 bytes)
- Contains /build/bin/: YES
- Contains frontend/node_modules/: YES
- Contains *.exe: YES
- Contains .DS_Store: YES

All acceptance criteria met: YES
Decision: ACCEPT
```

---

**This guide applies to:** L0-Coder, L1-Coder, L2-Coder, L3-Coder, L0-Planner, L0-Reviewer, L3-Architect

**Questions?** Refer to this guide before implementing.
