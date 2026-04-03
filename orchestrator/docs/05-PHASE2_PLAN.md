# Phase 2: Hardening - Implementation Plan

**Duration:** 3 days  
**Goal:** Robust, production-ready code with comprehensive tests

---

## Task 2.1: Unit Tests (90%+ Coverage)

### Components to Test
- [ ] FileTools (file_read, file_write, error cases)
- [ ] ToolExecutor (all 6 parsing patterns)
- [ ] LLMOrchestrator (retry logic, context simplification)
- [ ] Validators (models, api_keys, templates, startup)
- [ ] Exceptions (all 10 types)

### Test Structure
```
tests/unit/
├── test_file_tools.py
├── test_tool_executor.py
├── test_orchestrator.py
├── test_validators.py
└── test_exceptions.py
```

---

## Task 2.2: Integration Tests

### Test Scenarios
- [ ] L0-Planner end-to-end
- [ ] L0-Coder with file creation
- [ ] L0-Reviewer validation
- [ ] Escalation path (L0 → L1)
- [ ] Error handling (bad API key, network failure)

---

## Task 2.3: Exponential Backoff

### Implementation
- [ ] Replace fixed delay with exponential backoff
- [ ] Add jitter to prevent thundering herd
- [ ] Configurable max delay
- [ ] Progress logging

---

## Task 2.4: Context Window Management

### Features
- [ ] Auto-detect context window limits per model
- [ ] Smart truncation (preserve recent messages)
- [ ] Warning when approaching limits
- [ ] Fallback to more aggressive truncation

---

## Task 2.5: Rate Limiting Protection

### Implementation
- [ ] Track requests per minute/hour
- [ ] Automatic backoff on 429 responses
- [ ] Configurable rate limits per provider
- [ ] Queue system for burst handling

---

## Task 2.6: Structured Logging

### Implementation
- [ ] JSON log format
- [ ] Log levels (DEBUG, INFO, WARN, ERROR)
- [ ] Request/response logging
- [ ] Performance metrics in logs
- [ ] Correlation IDs for tracing

---

## Task 2.7: Metrics Collection

### Metrics to Track
- [ ] Requests per tier
- [ ] Success/failure rates
- [ ] Latency (p50, p95, p99)
- [ ] Cost per task
- [ ] Token usage
- [ ] Tool execution success rate

---

## Timeline

| Day | Tasks |
|-----|-------|
| 1 | 2.1 (Unit Tests), 2.2 (Integration Tests) |
| 2 | 2.3 (Backoff), 2.4 (Context), 2.5 (Rate Limit) |
| 3 | 2.6 (Logging), 2.7 (Metrics), Final Testing |

---

## Exit Criteria

- [ ] 90%+ unit test coverage
- [ ] All integration tests pass
- [ ] Retry handles transient failures gracefully
- [ ] Rate limiting prevents API bans
- [ ] Structured logging implemented
- [ ] Metrics dashboard functional

---

**Started:** 2026-03-30  
**Target Complete:** 2026-04-04
