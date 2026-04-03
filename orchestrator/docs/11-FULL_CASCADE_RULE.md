# Full Cascade Rule - MANDATORY WORKFLOW DISCIPLINE

**Date:** 2026-03-31  
**Priority:** CRITICAL - RESEARCH INTEGRITY  
**Status:** MANDATORY COMPLIANCE

---

## The Rule

> **Manual intervention is ONLY permitted after ALL coder tiers (L0, L1, L2, L3) have each received 3 attempts with reviewer feedback between each attempt.**

This is NOT optional. This is NOT negotiable. This is the core research methodology.

---

## Purpose

### 1. Iterative Learning
Each tier receives feedback and has opportunity to improve. We learn:
- What feedback is most effective for each model
- How models respond to correction
- Learning patterns across model families

### 2. Pressure Testing
Systematically test capabilities of every tier:
- L0 (Qwen3-Coder local): Baseline capability
- L1 (Grok 4.1 Fast): First escalation
- L2 (MiniMax M2.7): Complex reasoning
- L3 (Claude Sonnet 4.6): Best general coder

### 3. Data Collection
Generate benchmark dataset:
- Failure modes per model
- Attempt counts per success
- Cost vs. quality per tier
- Feedback effectiveness metrics

### 4. Research Value
Publishable findings:
- "Multi-Tier LLM Development: A Systematic Study"
- Model comparison data
- Cost optimization strategies
- Feedback loop effectiveness

---

## Complete Workflow Diagram

```
┌─────────────────┐
│  L0-Planner     │  Create specification
│  (Qwen3.5)      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  L0-Coder       │  Attempt 1
│  (Qwen3-Coder)  │  (Local, $0)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  L0-Reviewer    │  Review + Feedback
│  (Qwen3.5)      │  Accept or Reject
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
 ACCEPT   REJECT → L0-Coder Attempt 2
    │              (with feedback)
    │                   │
    │                   ▼
    │            ┌──────────────┐
    │            │  L0-Reviewer │  Re-review
    │            └──────┬───────┘
    │                   │
    │              ┌────┴────┐
    │              │         │
    │           ACCEPT   REJECT → L0-Coder Attempt 3
    │              │              (with feedback)
    │              │                   │
    │              │                   ▼
    │              │            ┌──────────────┐
    │              │            │  L0-Reviewer │  Final review
    │              │            └──────┬───────┘
    │              │                   │
    │              │              ┌────┴────┐
    │              │              │         │
    │              │           ACCEPT   REJECT/ESCALATE
    │              │              │
    │              │              ▼
    │              │       ┌──────────────┐
    │              │       │  L1-Coder    │  Attempt 1
    │              │       │  (Grok 4.1)  │  ($0.002/1K tokens)
    │              │       └──────┬───────┘
    │              │              │
    │              │              ▼
    │              │       ┌──────────────┐
    │              │       │  L0-Reviewer │  Review + Feedback
    │              │       └──────┬───────┘
    │              │              │
    │              │         [Repeat 3 attempts]
    │              │              │
    │              │              ▼
    │              │       ┌──────────────┐
    │              │       │  L2-Coder    │  Attempt 1
    │              │       │  (MiniMax)   │  ($0.0002/1K tokens)
    │              │       └──────┬───────┘
    │              │              │
    │              │              ▼
    │              │       ┌──────────────┐
    │              │       │  L0-Reviewer │  Review + Feedback
    │              │       └──────┬───────┘
    │              │              │
    │              │         [Repeat 3 attempts]
    │              │              │
    │              │              ▼
    │              │       ┌──────────────┐
    │              │       │  L3-Coder    │  Attempt 1
    │              │       │  (Claude     │  ($0.003/1K tokens)
    │              │       │   Sonnet)    │
    │              │       └──────┬───────┘
    │              │              │
    │              │              ▼
    │              │       ┌──────────────┐
    │              │       │  L0-Reviewer │  Review + Feedback
    │              │       └──────┬───────┘
    │              │              │
    │              │         [Repeat 3 attempts]
    │              │              │
    │              │         ┌────┴────┐
    │              │         │         │
    │              │      ACCEPT   REJECT
    │              │         │         │
    │              │         │    ┌────▼─────────────┐
    │              │         │    │  MANUAL          │
    │              │         │    │  INTERVENTION    │
    │              │         │    │  (Human takes    │
    │              │         │    │   over, documents│
    │              │         │    │   why all tiers  │
    │              │         │    │   failed)        │
    │              │         │    └──────────────────┘
    ▼              ▼
┌─────────────────────────┐
│  TASK COMPLETE          │  Continue to next task
└─────────────────────────┘
```

---

## Compliance Checklist

Before ANY manual intervention, verify ALL boxes are checked:

### L0-Coder Cycle
- [ ] L0-Coder Attempt 1 completed
- [ ] L0-Reviewer reviewed and rejected with feedback
- [ ] L0-Coder Attempt 2 completed (with feedback)
- [ ] L0-Reviewer re-reviewed and rejected with feedback
- [ ] L0-Coder Attempt 3 completed (with feedback)
- [ ] L0-Reviewer final review rejected

### L1-Coder Cycle
- [ ] L1-Coder Attempt 1 completed
- [ ] L0-Reviewer reviewed and rejected with feedback
- [ ] L1-Coder Attempt 2 completed (with feedback)
- [ ] L0-Reviewer re-reviewed and rejected with feedback
- [ ] L1-Coder Attempt 3 completed (with feedback)
- [ ] L0-Reviewer final review rejected

### L2-Coder Cycle
- [ ] L2-Coder Attempt 1 completed
- [ ] L0-Reviewer reviewed and rejected with feedback
- [ ] L2-Coder Attempt 2 completed (with feedback)
- [ ] L0-Reviewer re-reviewed and rejected with feedback
- [ ] L2-Coder Attempt 3 completed (with feedback)
- [ ] L0-Reviewer final review rejected

### L3-Coder Cycle
- [ ] L3-Coder Attempt 1 completed
- [ ] L0-Reviewer reviewed and rejected with feedback
- [ ] L3-Coder Attempt 2 completed (with feedback)
- [ ] L0-Reviewer re-reviewed and rejected with feedback
- [ ] L3-Coder Attempt 3 completed (with feedback)
- [ ] L0-Reviewer final review rejected

### Documentation
- [ ] All handoff logs preserved
- [ ] All reviewer feedback documented
- [ ] Failure analysis per tier completed
- [ ] Cost tracking complete
- [ ] Ready for publication

**If ANY box is unchecked → DO NOT INTERVENE MANUALLY**

---

## Cost Implications

| Tier | Model | Approx Cost/Task | 3 Attempts |
|------|-------|-----------------|------------|
| L0 | Qwen3-Coder (local) | $0.00 | $0.00 |
| L1 | Grok 4.1 Fast | $0.002 | $0.006 |
| L2 | MiniMax M2.7 | $0.0002 | $0.0006 |
| L3 | Claude Sonnet 4.6 | $0.003 | $0.009 |
| **Total** | | | **~$0.016/task** |

**Manual intervention cost:** $50-100/hour × 0.5-2 hours = **$25-200/task**

**Research value:** Priceless benchmark data on multi-tier LLM development

---

## What We're Measuring

### Per Task
- Attempts per tier
- Reviewer accept/reject decisions
- Feedback types and effectiveness
- Time per attempt
- Cost per tier
- Success/failure per tier

### Per Model
- Success rate by task type
- Common failure modes
- Response to feedback
- Cost-effectiveness ranking

### Per Workflow
- Optimal tier assignment by task complexity
- When escalation is most valuable
- Feedback loop optimization
- Human intervention frequency

---

## Publication Goals

### Planned Papers
1. "Multi-Tier LLM Development: Methodology and Initial Results"
2. "Cost-Optimized AI-Assisted Development: A Systematic Approach"
3. "Model Comparison for Code Generation: L0 Through L3"
4. "Feedback Loops in LLM Development Workflows"

### Dataset
- All handoff logs (anonymized)
- Reviewer decisions and feedback
- Attempt counts and outcomes
- Cost tracking data
- Human intervention analysis

---

## Enforcement

### Orchestrator Changes (TODO)
```python
class WorkflowEnforcer:
    def can_manual_intervene(self) -> bool:
        """Returns True ONLY after all tiers exhausted."""
        return all([
            self.l0_exhausted(),
            self.l1_exhausted(),
            self.l2_exhausted(),
            self.l3_exhausted()
        ])
    
    def l0_exhausted(self) -> bool:
        return self.coder_attempts['L0'] >= 3 and \
               self.reviewer_rejections['L0'] >= 3
    
    # ... similar for L1, L2, L3
```

### Audit Trail
Every task must have:
- Complete handoff chain (all tiers)
- Reviewer decisions logged
- Feedback preserved
- Final outcome documented

---

## Acknowledgment

I understand and commit to this workflow discipline:

- [ ] Manual intervention ONLY after ALL tiers exhausted
- [ ] Each tier gets 3 attempts with reviewer feedback
- [ ] Data collection is as important as task completion
- [ ] Research value justifies additional cost/time
- [ ] All handoffs preserved for publication

---

**Signed:** Development Team  
**Date:** 2026-03-31  
**Next Review:** After 10 tasks completed with full cascade

---

## Appendix: Task 2.3 Retrospective

**What Happened:**
- L0-Coder: 2 attempts, partial success
- L0-Reviewer: Invoked, rejected with feedback ✓
- **Manual intervention: PREMATURE** ❌

**Should Have Happened:**
- L0-Coder: 3 attempts with reviewer feedback
- L1-Coder: 3 attempts with reviewer feedback
- L2-Coder: 3 attempts with reviewer feedback
- L3-Coder: 3 attempts with reviewer feedback
- Manual: Only after all above exhausted

**Corrective Action:**
- Task 2.4+ will follow full cascade rule
- Task 2.3 data preserved as "partial compliance" example
- This document created to prevent future violations

---

**Remember:** We're not just building software. We're generating research data.
