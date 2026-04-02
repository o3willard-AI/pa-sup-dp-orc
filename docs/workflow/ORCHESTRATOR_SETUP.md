# Multi-Tier Orchestrator Setup Guide

## Overview

The orchestrator routes tasks to appropriate LLM models based on tier, enabling the multi-tier workflow as designed.

## Model Configuration

| Tier | Role | Model | Provider |
|------|------|-------|----------|
| L0-Planner | Task specs | Qwen3.5 397B A17B | OpenRouter |
| L0-Reviewer | Reviews | Qwen3.5 397B A17B | OpenRouter |
| L0-Coder | Implementation | Qwen3-Coder 30B | LM Studio (local) |
| L1-Coder | Re-implementation | Grok 4.1 Fast | OpenRouter |
| L2-Coder | Complex tasks | MiniMax M2.7 | OpenRouter |
| L3-Coder | Final escalation | Claude Sonnet 4.6 | OpenRouter |
| L3-Architect | Checkpoint reviews | Claude Opus 4.6 | OpenRouter |

## Setup Steps

### 1. Get OpenRouter API Key

1. Visit https://openrouter.ai/keys
2. Create account or sign in
3. Create new API key
4. Copy the key (starts with `sk-or-...`)

### 2. Set Environment Variable

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
export OPENROUTER_API_KEY="sk-or-your-key-here"
```

Or set temporarily for current session:

```bash
export OPENROUTER_API_KEY="sk-or-your-key-here"
```

### 3. Verify LM Studio Connection

The orchestrator expects LM Studio running at `http://192.168.101.21:1234`

Test connection:

```bash
curl http://192.168.101.21:1234/v1/models
```

If LM Studio is at different address, update `docs/workflow/lib/orchestrator.py`:

```python
"L0-Coder": {
    "provider": "lmstudio",
    "model": "qwen/qwen3-coder-30b",
    "base_url": "http://YOUR-IP:1234/v1"  # Update this
}
```

### 4. Test the Orchestrator

```bash
cd /home/sblanken/code/paa5
./docs/workflow/orchestrate.sh --help
```

Should show help with all available tiers.

### 5. Test L0-Planner (Free/Open)

```bash
cd /home/sblanken/code/paa5

# Create a simple test
cat > /tmp/test-context.json << 'EOF'
{
  "task_spec": "Create a simple Go function that returns 'Hello, World!'",
  "existing_context": "Fresh project, no existing code"
}
EOF

./docs/workflow/orchestrate.sh \
  --task TEST-001 \
  --tier L0-Planner \
  --context /tmp/test-context.json \
  --output /tmp/planner-output.txt
```

### 6. Test L0-Coder (Local/Free)

```bash
./docs/workflow/orchestrate.sh \
  --task TEST-001 \
  --tier L0-Coder \
  --context /tmp/test-context.json \
  --output /tmp/coder-output.txt
```

## Usage Examples

### Full Task Workflow

```bash
# Step 1: L0 Planner creates task spec
./docs/workflow/orchestrate.sh \
  --task 1.3 \
  --tier L0-Planner \
  --spec docs/tasks/1.3-template.md \
  --output docs/tasks/1.3-final.md

# Step 2: L0 Coder implements
./docs/workflow/orchestrate.sh \
  --task 1.3 \
  --tier L0-Coder \
  --spec docs/tasks/1.3-final.md \
  --output /tmp/impl.go

# Step 3: L0 Reviewer reviews
./docs/workflow/orchestrate.sh \
  --task 1.3 \
  --tier L0-Reviewer \
  --spec docs/tasks/1.3-final.md \
  --context '{"implementation": "..."}' \
  --output /tmp/review.txt
```

### Escalation Workflow

```bash
# L0 failed twice, escalate to L1
./docs/workflow/orchestrate.sh \
  --task 2.5 \
  --tier L1-Coder \
  --context escalation-context.json \
  --output /tmp/l1-impl.go

# L1 failed, escalate to L2
./docs/workflow/orchestrate.sh \
  --task 2.5 \
  --tier L2-Coder \
  --context escalation-context.json \
  --output /tmp/l2-impl.go

# L2 failed, escalate to L3
./docs/workflow/orchestrate.sh \
  --task 2.5 \
  --tier L3-Coder \
  --context escalation-context.json \
  --output /tmp/l3-impl.go
```

### Checkpoint Review

```bash
./docs/workflow/orchestrate.sh \
  --task CHECKPOINT-1 \
  --tier L3-Architect \
  --context checkpoint-1-context.json \
  --output docs/checkpoints/checkpoint-1-review.md
```

## Cost Tracking

The orchestrator logs all calls to `docs/workflow/handoffs/`. To estimate costs:

```bash
# View handoff logs
cat docs/workflow/handoffs/*.json | jq '.tier, .duration_seconds'

# OpenRouter usage dashboard: https://openrouter.ai/activity
```

**Approximate costs per 1K tokens:**
- Qwen3.5 397B: ~$0.003
- Qwen3-Coder (local): $0.00 (your hardware)
- Grok 4.1 Fast: ~$0.001
- MiniMax M2.7: ~$0.002
- Claude Sonnet 4.6: ~$0.003
- Claude Opus 4.6: ~$0.015

## Troubleshooting

### "OPENROUTER_API_KEY not set"
```bash
export OPENROUTER_API_KEY="your-key"
```

### "Connection refused" (LM Studio)
- Ensure LM Studio is running
- Check IP address: `http://192.168.101.21:1234`
- Verify model is loaded in LM Studio

### "Rate limit exceeded"
- OpenRouter has rate limits on free tier
- Wait a few minutes or upgrade plan

### "Model not found"
- Check model name in `orchestrator.py`
- Some models may require specific OpenRouter access

## Next Steps

1. Set `OPENROUTER_API_KEY` environment variable
2. Test with a simple task (e.g., Task 1.3)
3. Run full multi-tier workflow
4. Review handoff logs in `docs/workflow/handoffs/`
5. Update `WORKFLOW_LEARNINGS.md` with observations
