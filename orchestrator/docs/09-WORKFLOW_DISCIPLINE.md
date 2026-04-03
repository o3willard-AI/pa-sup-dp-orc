# Workflow Discipline - Critical Requirements

**Date:** 2026-03-31  
**Priority:** CRITICAL  
**Status:** MANDATORY COMPLIANCE

---

## Incident Report: Task 2.2 Workflow Violation

### What Happened

| Step | Expected | Actual | Status |
|------|----------|--------|--------|
| 1. L0-Coder executes | Creates files | Created files (partial) | ✓ |
| 2. Tool results analyzed | 2 success, 2 failures | 2 success, 2 failures | ✓ |
| 3. **L0-Reviewer invoked** | **Review output** | **SKIPPED** | ❌ **VIOLATION** |
| 4. Reviewer rejects | Send corrective feedback | Never happened | ❌ **VIOLATION** |
| 5. L0-Coder retry 1 | Fix issues | Never happened | ❌ **VIOLATION** |
| 6. Reviewer re-reviews | Check fixes | Never happened | ❌ **VIOLATION** |
| 7. L0-Coder retry 2 | Final correction | Never happened | ❌ **VIOLATION** |
| 8. Reviewer accepts | Pass to next task | Never happened | ❌ **VIOLATION** |
| 9. Manual intervention | Only after escalation | Happened at step 3 | ❌ **VIOLATION** |

### Root Cause

**Human impatience bypassed the workflow.**

The engineer (human) saw tool failures and immediately took over instead of:
1. Invoking L0-Reviewer
2. Allowing 2 correction cycles
3. Only escalating to manual after exhaustion

---

## Correct Workflow Behavior - MANDATORY

### The Golden Rule

> **NEVER intervene manually until the full cascade has been exhausted.**
> 
> This means: 3 L0-Coder attempts → 2 L0-Reviewer cycles → THEN manual

### Complete Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    TASK EXECUTION FLOW                          │
└─────────────────────────────────────────────────────────────────┘

    ┌──────────────────┐
    │  L0-Planner      │  Create task specification
    │  (Qwen3.5 397B)  │
    └────────┬─────────┘
             │ Spec
             ▼
    ┌──────────────────┐
    │  L0-Coder        │  Implement (Attempt 1)
    │  (Qwen3-Coder)   │  - 3 retry attempts with
    │                  │    progressive context
    └────────┬─────────┘    - Tool execution required
             │              - Files must be created
             │ Output + Tool Results
             ▼
    ┌──────────────────┐
    │  L0-Reviewer     │  Review implementation
    │  (Qwen3.5 397B)  │  - Check against spec
    │                  │  - Verify tool results
    └────────┬─────────┘  - Validate files exist
             │
      ┌──────┴──────┐
      │             │
   ACCEPT        REJECT
      │             │
      │      ┌──────┴──────────┐
      │      │  L0-Coder       │  Corrective Attempt 2
      │      │  (With feedback)│  - Must address all
      │      │                 │    reviewer comments
      │      └──────┬──────────┘
      │             │
      │      ┌──────┴──────┐
      │      │  L0-Reviewer│  Re-review
      │      └──────┬──────┘
      │             │
      │      ┌──────┴──────┐
      │      │             │
      │   ACCEPT        REJECT
      │      │             │
      │      │      ┌──────┴──────────┐
      │      │      │  L0-Coder       │  Final Attempt 3
      │      │      │  (With feedback)│
      │      │      └──────┬──────────┘
      │      │             │
      │      │      ┌──────┴──────┐
      │      │      │  L0-Reviewer│  Final review
      │      │      └──────┬──────┘
      │      │             │
      │      │      ┌──────┴──────┐
      │      │      │             │
      │      │   ACCEPT        REJECT/ESCALATE
      │      │      │             │
      │      │      │      ┌──────┴──────────────┐
      │      │      │      │  L1-Coder           │
      │      │      │      │  (Grok 4.1 Fast)    │
      │      │      │      │  Re-implement       │
      │      │      │      └─────────────────────┘
      │      │      │
      ▼      ▼      ▼
┌─────────────────────────┐
│  MANUAL INTERVENTION    │  ← ONLY AFTER ALL
│  (Human takes over)     │     ESCALATION EXHAUSTED
└─────────────────────────┘
```

---

## Why This Discipline Matters

### 1. Learning Loop

```
L0-Coder makes mistake → Reviewer provides feedback → L0-Coder learns
```

**Without reviewer:** L0-Coder never improves, same mistakes repeat

### 2. Cost Optimization

| Tier | Cost (approx) | Use Case |
|------|---------------|----------|
| L0-Coder | $0 (local) | First attempts, corrections |
| L0-Reviewer | $0.001 | Review cycles |
| L1-Coder | $0.002 | Escalation |
| Manual | $50-100/hr | Last resort only |

**Skipping reviewer → costs 10-100x more in human time**

### 3. Quality Assurance

- Reviewer catches issues before they propagate
- Reviewer validates tool execution actually succeeded
- Reviewer ensures spec compliance

### 4. Workflow Integrity

- Each tier has a purpose
- Skipping breaks the system
- Trust the cascade

---

## Mandatory Compliance Checklist

Before ANY manual intervention, verify:

- [ ] L0-Coder had 3 retry attempts
- [ ] L0-Reviewer reviewed output
- [ ] L0-Reviewer sent corrective feedback
- [ ] L0-Coder had 2 correction cycles
- [ ] L0-Reviewer rejected after final attempt
- [ ] L1-Coder escalation attempted
- [ ] L1-Coder failed or rejected
- [ ] L2-Coder escalation attempted (if needed)
- [ ] L3-Coder escalation attempted (if needed)
- [ ] ALL tiers exhausted

**If ANY box is unchecked → DO NOT INTERVENE MANUALLY**

---

## Corrective Actions for Task 2.2

### What Should Have Happened

1. **L0-Coder output received** (2 success, 2 failures)
2. **Invoke L0-Reviewer:**
   ```bash
   python3 orchestrator.py --task 2.2 --tier L0-Reviewer \
     --spec docs/tasks/2.2-spec.md \
     --context review_context.json
   ```
3. **Reviewer identifies issues:**
   - Escape characters in generated files
   - Incomplete implementation
   - Missing interface methods
4. **Reviewer rejects with feedback** → Back to L0-Coder
5. **L0-Coder retry 1** → Fix issues
6. **Reviewer re-reviews** → Still issues?
7. **L0-Coder retry 2** → Final correction
8. **Reviewer accepts** → Continue to Task 2.3
9. **Only if still failing** → Escalate to L1-Coder

### What We'll Do Differently

1. **Add workflow enforcement to orchestrator:**
   - Track correction cycles
   - Prevent manual intervention until exhausted
   - Log all reviewer decisions

2. **Add reviewer auto-invocation:**
   - After L0-Coder completes → Auto-invoke reviewer
   - No human decision point

3. **Add workflow validation:**
   - Pre-flight check before manual tasks
   - Verify all escalation paths exhausted

---

## Consequences of Non-Compliance

| Violation | Impact | Mitigation |
|-----------|--------|------------|
| Skipping reviewer | L0-Coder doesn't learn | Enforce auto-invocation |
| Early manual intervention | Wastes human time | Workflow gate in orchestrator |
| Not documenting deviations | Can't improve workflow | Mandatory incident reports |
| Ignoring tool failures | Broken files committed | Reviewer must verify files |

---

## Enforcement Mechanisms

### Orchestrator Changes (TODO)

```python
class WorkflowEnforcer:
    def __init__(self):
        self.coder_attempts = 0
        self.reviewer_cycles = 0
        self.max_coder_attempts = 3
        self.max_reviewer_cycles = 2
    
    def can_manual_intervene(self) -> bool:
        """Returns True only if all escalation paths exhausted."""
        return (
            self.coder_attempts >= self.max_coder_attempts and
            self.reviewer_cycles >= self.max_reviewer_cycles
        )
    
    def next_step(self) -> str:
        """Returns next required step in workflow."""
        if self.coder_attempts == 0:
            return "L0-Coder (Attempt 1)"
        elif self.reviewer_cycles < self.max_reviewer_cycles:
            return "L0-Reviewer (Cycle {self.reviewer_cycles + 1})"
        elif self.coder_attempts < self.max_coder_attempts:
            return "L0-Coder (Correction {self.coder_attempts})"
        else:
            return "ESCALATE to L1-Coder"
```

### Documentation Requirements

Every workflow deviation MUST be documented:
- Task ID
- What was skipped
- Why (emergency? time pressure?)
- Impact assessment
- Prevention plan

---

## Commitment

**I commit to following the workflow discipline:**

1. ✓ Never skip the reviewer
2. ✓ Allow full correction cycles
3. ✓ Only intervene manually after exhaustion
4. ✓ Document any deviations
5. ✓ Trust the cascade

---

**Acknowledged By:** Development Team  
**Date:** 2026-03-31  
**Next Review:** After Task 2.3 completion

---

## Appendix: Task 2.2 Incident Timeline

```
23:50:15 - L0-Planner generates spec (SUCCESS)
23:52:40 - L0-Coder implements (PARTIAL - 2/4 tools succeeded)
23:52:41 - ❌ HUMAN INTERVENES (Workflow violation)
23:53:00 - Human creates complete implementation
23:55:00 - Tests pass
23:56:00 - Task marked complete

MISSING:
- L0-Reviewer invocation
- Corrective feedback to L0-Coder
- L0-Coder correction attempts
- Final reviewer acceptance
```

**Lesson:** The workflow only works if we follow it.
