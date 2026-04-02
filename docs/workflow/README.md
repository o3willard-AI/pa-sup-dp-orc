# Multi-Tier LLM Development Workflow

This directory contains the prompt templates, tools, and documentation for the L0-L3 multi-tier development cascade used to build PairAdmin v2.0.

## Overview

The workflow uses a tiered cascade of LLM capabilities:

| Level | Tier | Role | Models | Provider |
|-------|------|------|--------|----------|
| **L0** | Planner/Reviewer | Task specs, reviews | Qwen3.5 397B A17B, DeepSeek V3.2 | OpenRouter |
| **L0** | Coder | First implementation | Qwen3-Coder (local), Step 3.5 Flash | LM Studio @ 192.168.101.21, OpenRouter |
| **L1** | Coder | Re-implementation (escalation 1) | Grok 4.1 Fast | OpenRouter (xAI) |
| **L2** | Coder | Re-implementation (escalation 2) | MiniMax M2.7 | OpenRouter (MiniMax) |
| **L3** | Coder | Final escalation implementation | Claude Sonnet 4.6 | OpenRouter (Anthropic) |
| **L3** | Architect | Checkpoint reviews, final authority | Claude Opus 4.6 | OpenRouter (Anthropic) |

## Quick Start

### For Planners (L0 - Qwen3.5 397B / DeepSeek V3.2)

1. Read the planner template: `cat templates/01-planner.md`
2. Read IMPLEMENTATION_PLAN.md section for your task
3. Create task spec in `docs/tasks/`
4. Log handoff: `./workflow.sh log-handoff`

### For Coders

**L0 (Qwen3-Coder local / Step 3.5 Flash):**
1. Read template: `cat templates/02-l0-coder.md`
2. Read task spec from `docs/tasks/`
3. Implement according to spec
4. Submit using submission format

**L1 (Grok 4.1 Fast):**
1. Read template: `cat templates/04-l1-coder.md`
2. Analyze why L0 failed
3. Re-implement with learning annotations

**L2 (MiniMax M2.7):**
1. Read template: `cat templates/05-l2-coder.md`
2. Categorize problem (A/B/C/D)
3. Implement/fix/decompose + recommend workflow changes

**L3 Coder (Claude Sonnet 4.6):**
1. Read template: `cat templates/05-l2-coder.md` (use L2 template for L3 coding)
2. Handle tasks that L2 could not complete
3. Flag patterns for Architect review

**L3 Architect (Claude Opus 4.6):**
1. Read template: `cat templates/06-l3-architect.md`
2. Make architectural decisions or checkpoint reviews
3. Document precedent and update guidelines (highest authority)

### For Reviewers (L0 - Qwen3.5 397B / DeepSeek V3.2)

1. Read reviewer template: `cat templates/03-reviewer.md`
2. Compare implementation against spec
3. Decide: ACCEPT, REJECT, or ESCALATE
4. Log the decision

## Directory Structure

```
workflow/
├── templates/          # Prompt templates for each tier
│   ├── 01-planner.md
│   ├── 02-l0-coder.md
│   ├── 03-reviewer.md
│   ├── 04-l1-coder.md
│   ├── 05-l2-coder.md
│   └── 06-l3-architect.md
├── handoffs/           # Handoff logs (JSON)
├── escalations/        # Escalation logs (JSON)
├── learnings/          # Learning documents (JSON)
├── lib/
│   └── workflow_tracker.py   # Python library for logging
├── workflow.sh         # CLI tool
└── README.md           # This file
```

## CLI Commands

```bash
# Show workflow status
./workflow.sh status

# View task history
./workflow.sh history 1.1

# List templates
./workflow.sh templates

# Find next task
./workflow.sh next-task
```

## Escalation Path

```
L0 Coder → Review → (Reject x2) → L1 Coder → Review → (Reject) → L2 Coder → Review → (Pattern Flag) → L3 Architect
```

## Logging Handoffs

Every tier transition should be logged:

```bash
# Example: L0 → Reviewer handoff
./workflow.sh log-handoff \
  --task-id 1.1 \
  --from-tier L0 \
  --to-tier Reviewer \
  --type "Review Ready"
```

Or provide JSON via stdin:

```bash
echo '{"files": ["main.go"], "verification": "go build ./... PASSED"}' | \
  ./workflow.sh log-handoff --task-id 1.1 --from-tier L0 --to-tier Reviewer --type "Review Ready"
```

## Learning Documents

After each escalation resolution or checkpoint, a learning document is created in `docs/learnings/`. These documents capture:

- Why previous tiers failed
- What approach succeeded
- Workflow recommendations
- Template update suggestions

## Checkpoint Reviews

At each QA Checkpoint milestone (see `QA_CHECKPOINTS.md`), the L3 Architect produces a comprehensive review:

1. All task specs and implementations reviewed
2. Escalation patterns analyzed
3. Model performance assessed
4. Workflow guidelines updated

Checkpoint reports are saved to `docs/checkpoints/`.

## Template Updates

Templates are living documents. Update when:

- Same issue causes 3+ escalations
- Checkpoint review identifies pattern
- Human team approves change
- Model capabilities change

Log all template changes in a CHANGELOG (to be created).

## Related Documents

- **Design Spec:** `docs/superpowers/specs/2026-03-30-multi-tier-development-workflow-design.md`
- **Implementation Plan:** `IMPLEMENTATION_PLAN.md`
- **QA Checkpoints:** `QA_CHECKPOINTS.md`
- **Task Examples:** `TASK_EXAMPLE.md`
