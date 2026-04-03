# Orchestrator Quick Reference

**Print this and keep it at your desk**

---

## Workflow Cascade - NEVER SKIP

```
┌──────────────┐
│ L0-Planner   │  Create spec
└──────┬───────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐
│ L0-Coder     │────▶│ L0-Reviewer  │
│ (3 attempts) │     │ (2 cycles)   │
└──────────────┘     └──────┬───────┘
                            │
                       ┌────┴────┐
                       │         │
                    ACCEPT   REJECT (back to L0-Coder)
                       │
                       ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ L1-Coder     │─▶│ L2-Coder     │─▶│ L3-Coder     │
└──────────────┘  └──────────────┘  └──────┬───────┘
                                           │
                                           ▼
                                    MANUAL (LAST)
```

---

## Command Reference

### Generate Spec (L0-Planner)
```bash
export OPENROUTER_API_KEY="sk-or-..."
python3 orchestrator/src/core/orchestrator.py \
  --task 2.3 \
  --tier L0-Planner \
  --context docs/tasks/2.3-context.json \
  --output docs/tasks/2.3-spec.md
```

### Implement (L0-Coder)
```bash
python3 orchestrator/src/core/orchestrator.py \
  --task 2.3 \
  --tier L0-Coder \
  --spec docs/tasks/2.3-spec.md \
  --context docs/tasks/2.3-coder-context.json \
  --output docs/tasks/2.3-implementation.md
```

### Review (L0-Reviewer) - MANDATORY
```bash
python3 orchestrator/src/core/orchestrator.py \
  --task 2.3 \
  --tier L0-Reviewer \
  --spec docs/tasks/2.3-spec.md \
  --context docs/tasks/2.3-review-context.json
```

### Validate Setup
```bash
python3 -m src.validators.startup
```

### Run Smoke Tests
```bash
PYTHONPATH=/path/to/orchestrator python3 tests/e2e/test_smoke.py
```

---

## Tier Configuration

| Tier | Model | Provider | Use |
|------|-------|----------|-----|
| L0-Planner | qwen/qwen3.5-397b-a17b | OpenRouter | Task specs |
| L0-Reviewer | qwen/qwen3.5-397b-a17b | OpenRouter | Reviews |
| L0-Coder | qwen/qwen3-coder-30b | LM Studio | Implementation |
| L1-Coder | x-ai/grok-4.1-fast | OpenRouter | Escalation 1 |
| L2-Coder | minimax/minimax-m2.7 | OpenRouter | Escalation 2 |
| L3-Coder | anthropic/claude-sonnet-4.6 | OpenRouter | Final |

---

## Troubleshooting

### "OPENROUTER_API_KEY not set"
```bash
export OPENROUTER_API_KEY="sk-or-your-key-here"
```

### LM Studio connection failed
```bash
# Verify LM Studio is running on 192.168.101.21:1234
curl http://192.168.101.21:1234/v1/models
```

### Tool execution failed
1. Check LLM output for actual tool calls
2. Verify regex patterns match LLM format
3. Review handoff JSON for details

### Context limit exceeded
- Orchestrator auto-truncates on retries
- If persistent, escalate to higher tier

---

## Before Manual Intervention - CHECKLIST

- [ ] L0-Coder had 3 attempts
- [ ] L0-Reviewer reviewed output
- [ ] L0-Reviewer sent feedback
- [ ] L0-Coder had 2 correction cycles
- [ ] L1-Coder escalation attempted
- [ ] All tiers exhausted

**If ANY box unchecked → DO NOT INTERVENE**

---

## File Locations

```
orchestrator/
├── src/core/
│   ├── orchestrator.py      # Main orchestrator
│   ├── exceptions.py        # Error types
│   ├── retry.py             # Backoff logic
│   ├── feedback.py          # Tool feedback
│   ├── metrics.py           # Metrics collection
│   └── cost.py              # Cost tracking
├── src/validators/
│   ├── startup.py           # Pre-flight checks
│   ├── models.py            # Model validation
│   ├── api_keys.py          # API key check
│   └── templates.py         # Template check
├── tests/
│   ├── unit/                # Unit tests
│   └── e2e/                 # Smoke tests
├── templates/               # Prompt templates
└── docs/                    # Documentation
```

---

## Emergency Contacts

| Issue | Resolution |
|-------|------------|
| API key invalid | Generate new key at openrouter.ai |
| LM Studio down | Restart on 192.168.101.21 |
| Orchestrator error | Check `docs/workflow/handoffs/` for logs |
| Tool parsing fails | Review regex in `ToolExecutor` |

---

**Remember:** The workflow only works if you follow it.

**Last Updated:** 2026-03-31
