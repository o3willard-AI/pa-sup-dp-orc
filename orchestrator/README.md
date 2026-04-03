# Multi-Tier LLM Orchestrator

**Version:** 0.1.0-dev  
**Status:** RECOVERY MODE  
**License:** MIT

A world-class orchestration system for cost-optimized AI-assisted development using multi-tier LLM cascades.

---

## Overview

The Multi-Tier Orchestrator routes development tasks through a hierarchy of LLM models based on complexity and cost:

```
┌─────────────────┐
│  L0-Planner     │  Qwen3.5 397B - Task specifications
│  (OpenRouter)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  L0-Coder       │  Qwen3-Coder (local) - First implementation
│  (LM Studio)    │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
 ACCEPT   REJECT → Escalate
    │         │
    │         ▼
    │  ┌──────────────┐
    │  │  L1-Coder    │  Grok 4.1 Fast - Re-implementation
    │  │  (OpenRouter)│
    │  └──────┬───────┘
    │         │
    │    ┌────┴────┐
    │    │         │
    │ ACCEPT  REJECT → Escalate
    │    │         │
    │    │         ▼
    │    │  ┌──────────────┐
    │    │  │  L2-Coder    │  MiniMax M2.7 - Complex tasks
    │    │  │  (OpenRouter)│
    │    │  └──────┬───────┘
    │    │         │
    │    │    ┌────┴────┐
    │    │    │         │
    │    │ ACCEPT  REJECT → Escalate
    │    │    │         │
    │    │    │         ▼
    │    │    │  ┌──────────────┐
    │    │    │  │  L3-Coder    │  Claude Sonnet 4.6 - Final
    │    │    │  │  (OpenRouter)│
    │    │    │  └──────────────┘
    ▼    ▼    ▼
┌─────────────────┐
│  L3-Architect   │  Claude Opus 4.6 - Checkpoint reviews
│  (OpenRouter)   │
└─────────────────┘
```

---

## Features

- **Multi-Tier Cascade** - Automatic escalation through model tiers
- **Native Tool Support** - Parses Claude JSON, MiniMax XML, and custom formats
- **Retry Logic** - 3 attempts per tier with exponential backoff
- **Model Validation** - Pre-flight checks for model availability
- **Cost Tracking** - Monitor spending per task and tier
- **Structured Logging** - JSON logs for observability
- **Type Safe** - Full mypy type checking

---

## Quick Start

### Installation

```bash
cd orchestrator
pip install -e .[dev]
```

### Configuration

```bash
export OPENROUTER_API_KEY="sk-or-your-key-here"
export LMSTUDIO_BASE_URL="http://192.168.101.21:1234/v1"
```

### Validate Setup

```bash
python -m src.validators.models
```

### Run Task

```bash
orchestrator --task 1.0 \
  --tier L0-Planner \
  --spec docs/tasks/1.0-spec.md \
  --output result.md
```

---

## Project Status

| Phase | Status | Progress |
|-------|--------|----------|
| Phase 1: Stabilization | 🟡 In Progress | 0% |
| Phase 2: Hardening | ⚪ Pending | 0% |
| Phase 3: Enhancement | ⚪ Pending | 0% |
| Phase 4: Polish | ⚪ Pending | 0% |

**Current Focus:** Phase 1 - Writing clean orchestrator from scratch

---

## Documentation

- [Get Well Plan](docs/01-GET_WELL_PLAN.md)
- [Phase 1 Tasks](docs/02-PHASE1_TASKS.md)
- [API Reference](docs/03-API_REFERENCE.md) _(TODO)_
- [User Guide](docs/04-USER_GUIDE.md) _(TODO)_

---

## Development

### Running Tests

```bash
# Unit tests
pytest tests/unit/

# Integration tests
pytest tests/integration/

# End-to-end tests
pytest tests/e2e/

# With coverage
pytest --cov=src --cov-report=html
```

### Code Quality

```bash
# Type checking
mypy src/

# Formatting
black src/ tests/

# Linting
ruff check src/ tests/
```

---

## Architecture

```
src/
├── main.py              # CLI entry point
├── config.py            # Configuration management
├── core/
│   ├── orchestrator.py  # Main orchestration logic
│   ├── executor.py      # Task execution engine
│   ├── retry.py         # Retry with backoff
│   └── exceptions.py    # Custom exceptions
├── tools/
│   ├── base.py          # Tool interface
│   ├── file_tools.py    # File read/write
│   └── registry.py      # Tool registration
├── parsers/
│   ├── base.py          # Parser interface
│   ├── custom.py        # file_write() syntax
│   ├── claude.py        # Claude JSON format
│   └── minimax.py       # MiniMax XML format
└── validators/
    ├── models.py        # Model ID validation
    ├── api_keys.py      # API key validation
    └── templates.py     # Template validation
```

---

## License

MIT License - see [LICENSE](LICENSE) for details.

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests
4. Ensure all tests pass
5. Submit pull request

---

**Last Updated:** 2026-03-30  
**Maintained By:** PairAdmin Team

---

## ⚠️ WORKFLOW DISCIPLINE - MANDATORY

**CRITICAL:** Never intervene manually until the full escalation cascade is exhausted:

```
L0-Coder (3 attempts) → L0-Reviewer (2 cycles) → L1-Coder → L2-Coder → L3-Coder → MANUAL
```

**See:** [docs/09-WORKFLOW_DISCIPLINE.md](docs/09-WORKFLOW_DISCIPLINE.md) for complete requirements.

### Quick Reference

| When | Action |
|------|--------|
| L0-Coder fails | → Auto-retry (3 attempts) |
| L0-Coder succeeds partially | → L0-Reviewer (NOT manual) |
| L0-Reviewer rejects | → Back to L0-Coder (2 cycles) |
| All L0 cycles exhausted | → Escalate to L1-Coder |
| All tiers exhausted | → THEN manual intervention |

**Violation = Broken workflow + Wasted human time + L0-Coder never learns**


---

## ⚠️ WORKFLOW DISCIPLINE - MANDATORY

**CRITICAL:** Never intervene manually until the full escalation cascade is exhausted:

```
L0-Coder (3 attempts) → L0-Reviewer (2 cycles) → L1-Coder → L2-Coder → L3-Coder → MANUAL
```

**See:** [docs/09-WORKFLOW_DISCIPLINE.md](docs/09-WORKFLOW_DISCIPLINE.md) for complete requirements.

### Quick Reference

| When | Action |
|------|--------|
| L0-Coder fails | → Auto-retry (3 attempts) |
| L0-Coder succeeds partially | → L0-Reviewer (NOT manual) |
| L0-Reviewer rejects | → Back to L0-Coder (2 cycles) |
| All L0 cycles exhausted | → Escalate to L1-Coder |
| All tiers exhausted | → THEN manual intervention |

**Violation = Broken workflow + Wasted human time + L0-Coder never learns**

