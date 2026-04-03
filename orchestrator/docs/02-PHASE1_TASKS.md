# Phase 1: Stabilization - Detailed Tasks

**Duration:** 2 days  
**Goal:** Working baseline with all core features functional

---

## Task 1.1: Write Clean orchestrator.py From Scratch

**Estimated:** 4 hours  
**Priority:** CRITICAL  
**Status:** TODO

### Acceptance Criteria
- [ ] File has no syntax errors
- [ ] All classes properly defined
- [ ] No duplicate code
- [ ] Proper imports
- [ ] Type hints throughout

### Implementation Notes
- Write entire file in one operation (no incremental edits)
- Use simple, clear code structure
- Include all 3 main classes: FileTools, ToolExecutor, LLMOrchestrator
- Include main() function with proper argument parsing

### Files to Create
- `orchestrator/src/core/orchestrator.py` (main implementation)
- `orchestrator/src/__init__.py` (package init)
- `orchestrator/src/main.py` (CLI entry point)

---

## Task 1.2: Implement All 6 Tool Parsing Patterns

**Estimated:** 3 hours  
**Priority:** CRITICAL  
**Status:** TODO

### Patterns to Support
1. Custom `file_write("path", """content""")`
2. Custom `file_read("path")`
3. Claude JSON write: `{"name": "write_file", ...}`
4. Claude JSON read: `{"name": "read_file", ...}`
5. MiniMax XML write: `<invoke name="Write">...`
6. MiniMax XML read: `<invoke name="Read">...`

### Acceptance Criteria
- [ ] Each pattern has dedicated parser class
- [ ] Unit tests for each pattern
- [ ] Handles edge cases (empty content, special chars)
- [ ] Returns structured results

### Files to Create
- `orchestrator/src/parsers/base.py`
- `orchestrator/src/parsers/custom.py`
- `orchestrator/src/parsers/claude.py`
- `orchestrator/src/parsers/minimax.py`

---

## Task 1.3: Add Comprehensive Error Handling

**Estimated:** 2 hours  
**Priority:** HIGH  
**Status:** TODO

### Error Types to Handle
- API connection failures
- Rate limiting (429)
- Invalid API keys (401)
- Model not found (404)
- Context too long (400)
- File I/O errors
- Template not found
- Invalid tool calls

### Acceptance Criteria
- [ ] Custom exception classes
- [ ] Meaningful error messages
- [ ] Error context preserved
- [ ] Graceful degradation

### Files to Create
- `orchestrator/src/core/exceptions.py`

---

## Task 1.4: Add Startup Validation

**Estimated:** 2 hours  
**Priority:** HIGH  
**Status:** TODO

### Validations
- [ ] Model IDs exist (via OpenRouter API)
- [ ] API keys set for required providers
- [ ] Template files exist
- [ ] Output directories writable
- [ ] Python version compatible

### Acceptance Criteria
- [ ] Validation runs before first task
- [ ] Clear error messages
- [ ] Fix suggestions provided
- [ ] Can skip validation for offline mode

### Files to Create
- `orchestrator/src/validators/models.py`
- `orchestrator/src/validators/api_keys.py`
- `orchestrator/src/validators/templates.py`
- `orchestrator/src/validators/__init__.py`

---

## Task 1.5: Verify with Smoke Tests

**Estimated:** 3 hours  
**Priority:** HIGH  
**Status:** TODO

### Smoke Tests
1. L0-Planner creates task spec
2. L0-Coder creates file via tool call
3. L0-Reviewer validates file
4. Escalation to L1-Coder works
5. Error handling works (bad API key)

### Acceptance Criteria
- [ ] All smoke tests pass
- [ ] Test execution < 5 minutes
- [ ] Clear pass/fail output
- [ ] Logs captured for debugging

### Files to Create
- `orchestrator/tests/e2e/test_smoke.py`
- `orchestrator/scripts/run_smoke_tests.sh`

---

## Phase 1 Exit Checklist

- [ ] All 5 tasks complete
- [ ] No syntax errors in any file
- [ ] All smoke tests pass
- [ ] Documentation updated
- [ ] Code reviewed

---

## Dependencies

```
Task 1.1 → Task 1.2 → Task 1.3 → Task 1.4 → Task 1.5
     ↓                                         ↑
     └─────────────────────────────────────────┘
```

---

## Time Estimate

| Task | Hours |
|------|-------|
| 1.1 | 4 |
| 1.2 | 3 |
| 1.3 | 2 |
| 1.4 | 2 |
| 1.5 | 3 |
| **Total** | **14 hours** (~2 days) |
