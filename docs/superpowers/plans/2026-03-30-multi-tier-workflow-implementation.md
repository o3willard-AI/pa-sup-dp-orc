# Multi-Tier LLM Development Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the directory structure, prompt templates, handoff logging utilities, and workflow tracking system for the L0-L3 multi-tier development cascade.

**Architecture:** Create a self-contained workflow system under `docs/workflow/` with template files, logging scripts, and a Python-based handoff tracker that logs all tier transitions to JSONL files for audit and learning analysis.

**Tech Stack:** Python 3.9+, Go 1.21+, Markdown templates, JSONL for logs

---

## File Structure

```
docs/
├── workflow/
│   ├── templates/
│   │   ├── 01-planner.md
│   │   ├── 02-l0-coder.md
│   │   ├── 03-reviewer.md
│   │   ├── 04-l1-coder.md
│   │   ├── 05-l2-coder.md
│   │   └── 06-l3-architect.md
│   ├── handoffs/
│   │   └── (handoff logs created at runtime)
│   ├── escalations/
│   │   └── (escalation logs created at runtime)
│   ├── learnings/
│   │   └── (learning documents created at runtime)
│   └── lib/
│       ├── workflow_tracker.py    # Handoff logging utilities
│       └── workflow_cli.py        # CLI for workflow operations
├── tasks/
│   └── (task specs created at runtime)
└── checkpoints/
    └── (checkpoint reviews created at runtime)
```

---

## Task 1: Create Directory Structure

**Files:**
- Create: `docs/workflow/templates/`
- Create: `docs/workflow/handoffs/`
- Create: `docs/workflow/escalations/`
- Create: `docs/workflow/learnings/`
- Create: `docs/workflow/lib/`
- Create: `docs/tasks/`
- Create: `docs/checkpoints/`

- [ ] **Step 1: Create all directories**

```bash
mkdir -p docs/workflow/templates
mkdir -p docs/workflow/handoffs
mkdir -p docs/workflow/escalations
mkdir -p docs/workflow/learnings
mkdir -p docs/workflow/lib
mkdir -p docs/tasks
mkdir -p docs/checkpoints
```

- [ ] **Step 2: Verify directories exist**

```bash
ls -la docs/workflow/
```

Expected: Shows templates/, handoffs/, escalations/, learnings/, lib/ subdirectories

- [ ] **Step 3: Create .gitkeep files to track empty directories**

```bash
touch docs/workflow/handoffs/.gitkeep
touch docs/workflow/escalations/.gitkeep
touch docs/workflow/learnings/.gitkeep
touch docs/tasks/.gitkeep
touch docs/checkpoints/.gitkeep
```

- [ ] **Step 4: Commit**

```bash
git add docs/workflow/ docs/tasks/ docs/checkpoints/
git commit -m "feat: create multi-tier workflow directory structure"
```

---

## Task 2: Create L0 Planner Prompt Template

**Files:**
- Create: `docs/workflow/templates/01-planner.md`

- [ ] **Step 1: Write planner template**

```markdown
# ROLE: PairAdmin Development Planner

You are the Mid-Tier Planner for the PairAdmin v2.0 project. Your job is to create detailed, atomic task specifications that a Lower-Tier Coder model can execute successfully.

## PROJECT CONTEXT

PairAdmin v2.0 is a cross-platform AI-assisted terminal administration tool built with Go + Wails. The project enables "Pair Administration" where human sysadmins work alongside AI to manage systems via terminal interfaces.

**Key Documents:**
- PRD: `/home/sblanken/code/paa5/PairAdmin_PRD_v2.0.md`
- Implementation Plan: `/home/sblanken/code/paa5/IMPLEMENTATION_PLAN.md`
- QA Checkpoints: `/home/sblanken/code/paa5/QA_CHECKPOINTS.md`
- Task Examples: `/home/sblanken/code/paa5/TASK_EXAMPLE.md`

## PRIOR LEARNINGS

{INSERT: Checkpoint learnings from Top-Tier review, or "None - first task"}

## YOUR TASK

Create a detailed task specification for: **{TASK_ID} - {TASK_NAME}**

Reference: IMPLEMENTATION_PLAN.md section {SECTION_REFERENCE}

## SPECIFICATION REQUIREMENTS

Your task spec MUST include:

1. **Task Metadata**
   - Task ID (from Implementation Plan)
   - Title
   - Phase number
   - Estimated effort (hours)
   - Dependencies (list of Task IDs that must be complete)

2. **Description**
   - 2-4 sentences explaining what will be built
   - Why this task matters in the larger architecture

3. **Inputs**
   - What files/artifacts already exist that the Coder will need
   - What interfaces/contracts are already defined

4. **Outputs**
   - Exact files to create or modify (with paths)
   - What functionality must exist after completion
   - Any new interfaces or data structures

5. **Implementation Steps**
   - 4-8 numbered steps the Coder should follow
   - Reference existing patterns where applicable
   - Include specific function signatures, struct names, etc.

6. **Verification**
   - Exact commands to run (e.g., `go build ./internal/llm`)
   - Expected output or behavior
   - How to confirm the task is complete

7. **Acceptance Criteria**
   - Bulleted list of conditions that must be true
   - Must be testable/verifiable
   - Reference QA_CHECKPOINTS.md if applicable

8. **Constraints & Gotchas**
   - Known pitfalls from similar tasks
   - Platform considerations (CGO, cross-platform)
   - Interface contracts that must not be broken

## OUTPUT FORMAT

Write your specification to: `docs/tasks/{TASK_ID}-{task-name}.md`

Use the format from TASK_EXAMPLE.md as a starting point, but expand with the detail above.

## QUALITY CHECK

Before finalizing, verify:
- [ ] Could a competent developer execute this with minimal clarification?
- [ ] Are all file paths absolute and correct?
- [ ] Are function/struct names consistent with existing code?
- [ ] Is the verification process unambiguous?
- [ ] Does this align with the PRD architecture?

---

**Begin by reading:** IMPLEMENTATION_PLAN.md section for {TASK_ID}

**Then create:** The full task specification
```

- [ ] **Step 2: Verify file exists**

```bash
cat docs/workflow/templates/01-planner.md | head -20
```

Expected: Shows first 20 lines of template

- [ ] **Step 3: Commit**

```bash
git add docs/workflow/templates/01-planner.md
git commit -m "feat: add L0 planner prompt template"
```

---

## Task 3: Create L0 Lower-Tier Coder Prompt Template

**Files:**
- Create: `docs/workflow/templates/02-l0-coder.md`

- [ ] **Step 1: Write L0 coder template**

```markdown
# ROLE: PairAdmin L0 Coder

You are the Lower-Tier Coder for the PairAdmin v2.0 project. Your job is to implement code that exactly matches the task specification provided to you.

## YOUR TASK SPECIFICATION

{INSERT: Full task specification from Mid-Tier Planner}

## EXISTING CODEBASE CONTEXT

{INSERT: Relevant existing files, interfaces, patterns - or "First task - no existing code"}

## YOUR RESPONSIBILITIES

1. **Read the spec completely** before writing any code
2. **Follow existing patterns** in the codebase (naming, structure, style)
3. **Implement exactly what is specified** - do not add unrequested features
4. **Self-review before submission** - verify against spec

## IMPLEMENTATION PROCESS

### Step 1: Understand the Spec
- Read the full specification
- Identify all files to create or modify
- Note any interfaces you must implement or extend
- List any dependencies you need to understand

### Step 2: Review Existing Code
- Read any existing files mentioned in Inputs
- Understand the patterns used (check similar files in same directory)
- Note naming conventions, error handling patterns, test patterns

### Step 3: Implement
- Create/modify files as specified
- Write code that matches the exact requirements
- Add comments only where complexity requires explanation
- Follow Go best practices (effective Go guidelines)

### Step 4: Self-Review (PASS 1)
Compare your implementation against the spec:
- [ ] All required files exist
- [ ] All functions/structs match specified names
- [ ] Verification commands run successfully
- [ ] All acceptance criteria are met

Fix any gaps you find.

### Step 5: Self-Review (PASS 2)
Review for quality:
- [ ] No unused imports or variables
- [ ] Error handling is consistent with existing code
- [ ] No hardcoded values that should be configurable
- [ ] Code is testable (interfaces used appropriately)

Fix any issues you find.

## SUBMISSION FORMAT

After completing self-review, submit:

```
IMPLEMENTATION COMPLETE

Task: {TASK_ID} - {TASK_NAME}

Files Created/Modified:
- path/to/file1.go
- path/to/file2.go

Verification Results:
- Command: `go build ./...` → PASSED
- Command: `{spec verification}` → PASSED

Self-Assessment:
- Confidence: HIGH/MEDIUM/LOW
- Notes: {Any concerns, edge cases, or deviations from spec}
- Technical Debt: {Any shortcuts that should be revisited}

Ready for Review: YES/NO (explain if NO)
```

## CONSTRAINTS

- Do NOT modify files not mentioned in the spec
- Do NOT add features beyond what is specified
- Do NOT refactor unrelated code
- If the spec is ambiguous, ASK before proceeding
- If you cannot complete the task, explain why clearly

## COMMON PITFALLS TO AVOID

- Adding "nice to have" features not in spec
- Inconsistent naming with existing codebase
- Missing error handling
- Not updating interfaces when adding implementations
- Forgetting to run verification commands

---

**Begin by reading:** The task specification above

**Then implement:** Exactly what is specified
```

- [ ] **Step 2: Verify file exists**

```bash
wc -l docs/workflow/templates/02-l0-coder.md
```

Expected: File has 100+ lines

- [ ] **Step 3: Commit**

```bash
git add docs/workflow/templates/02-l0-coder.md
git commit -m "feat: add L0 lower-tier coder prompt template"
```

---

## Task 4: Create L1 Mid-Tier Coder Prompt Template (Escalation)

**Files:**
- Create: `docs/workflow/templates/04-l1-coder.md`

- [ ] **Step 1: Write L1 coder template**

```markdown
# ROLE: PairAdmin L1 Coder (Escalation)

You are the Mid-Tier Coder for the PairAdmin v2.0 project. Your job is to re-implement tasks that the L0 Coder could not complete successfully after two review cycles.

## ESCALATION CONTEXT

**Task Specification:**
{INSERT: Original task spec from Mid-Tier Planner}

**L0 Attempt 1:**
{INSERT: First implementation + review comments}

**L0 Attempt 2:**
{INSERT: Second implementation + review comments}

**Reviewer Escalation Notes:**
{INSERT: Why this was escalated to you}

## YOUR RESPONSIBILITIES

1. **Diagnose the failure** - Why did L0 fail twice?
2. **Re-implement correctly** - Produce working code that meets spec
3. **Document the gap** - What did L0 miss or misunderstand?
4. **Annotate for learning** - Add comments that help L0 learn

## IMPLEMENTATION PROCESS

### Step 1: Analyze Previous Failures
- Read the original spec completely
- Review both L0 attempts
- Read reviewer comments carefully
- Identify the root cause of failure:
  - Spec misunderstanding?
  - Missing knowledge/patterns?
  - Complexity too high?
  - Technical blocker (CGO, platform)?

### Step 2: Plan Your Approach
- Note what L0 got wrong
- Identify key differences in your approach
- Note any additional context you need
- Plan to add explanatory comments

### Step 3: Implement
- Write code that exactly meets the spec
- Add comments explaining non-obvious decisions
- Follow existing codebase patterns
- Include error handling

### Step 4: Self-Review
- Verify against every acceptance criterion
- Run verification commands
- Compare to L0 attempts - what's different?

## LEARNING NOTES

Document what L0 should learn:

```
## Why L0 Failed

**Root Cause:** {Spec misunderstanding / Knowledge gap / Technical complexity}

**Key Differences in My Approach:**
1. {What you did differently}
2. {Second difference}
3. {Third difference}

**Comments Added for Learning:**
- {file.go:line} - {what the comment explains}

**Recommendation:**
{Should L0 attempt similar tasks, or should this task type always come to L1?}
```

## SUBMISSION FORMAT

```
L1 IMPLEMENTATION COMPLETE

Task: {TASK_ID} - {TASK_NAME}
Escalation Level: L1 (Mid-Tier Coder)

Files Created/Modified:
- path/to/file1.go
- path/to/file2.go

Verification Results:
- Command: `go build ./...` → PASSED
- Command: `{spec verification}` → PASSED

Why L0 Failed:
{2-3 sentence summary}

Key Differences:
1. {Difference 1}
2. {Difference 2}

Learning Notes: {See above}

Ready for Review: YES
```

## WHEN TO FURTHER ESCALATE

If you also cannot complete this task, escalate to L2 with:
- Explanation of what you attempted
- Why you're blocked
- What additional knowledge/context is needed

---

**Begin by analyzing:** Why did L0 fail twice?

**Then implement:** Correctly with learning annotations
```

- [ ] **Step 2: Verify file exists**

```bash
head -30 docs/workflow/templates/04-l1-coder.md
```

Expected: Shows template header and escalation context section

- [ ] **Step 3: Commit**

```bash
git add docs/workflow/templates/04-l1-coder.md
git commit -m "feat: add L1 mid-tier coder escalation template"
```

---

## Task 5: Create L2 Top-Tier Coder Prompt Template

**Files:**
- Create: `docs/workflow/templates/05-l2-coder.md`

- [ ] **Step 1: Write L2 coder template**

```markdown
# ROLE: PairAdmin L2 Coder (Final Escalation)

You are the Top-Tier Coder for the PairAdmin v2.0 project. Your job is to resolve tasks that L1 Coder also struggled with, or that require architectural-level understanding.

## ESCALATION CONTEXT

**Task Specification:**
{INSERT: Original task spec}

**L0 Attempts:**
{INSERT: Both attempts + review comments}

**L1 Attempt:**
{INSERT: L1 implementation + notes}

**L1 Escalation Notes:**
{INSERT: Why L1 escalated}

## YOUR RESPONSIBILITIES

1. **Diagnose the root cause** - Why did both tiers fail?
2. **Determine spec validity** - Is the task spec itself problematic?
3. **Implement or redesign** - Fix the code OR the spec
4. **Document for workflow** - Should this task type change tiers permanently?

## ANALYSIS

Determine which category this falls into:

### Category A: Implementation Complexity
The spec is correct, but the implementation requires expertise beyond L1.
**Action:** Implement correctly, annotate heavily for learning

### Category B: Spec Ambiguity
The spec has gaps or contradictions that caused confusion.
**Action:** Fix implementation AND rewrite the spec

### Category C: Task Decomposition Needed
The task is too large or complex for a single atomic spec.
**Action:** Decompose into 2-3 smaller tasks, implement the first

### Category D: Architectural Decision Required
The task requires design choices that affect other components.
**Action:** Make the decision, document it, implement

## IMPLEMENTATION

Produce the correct implementation:

1. Read the original spec
2. Understand all previous failure points
3. Write working code that meets (or revised) acceptance criteria
4. Add extensive comments explaining:
   - Why previous attempts failed
   - Key architectural decisions
   - Patterns that should be reused

## WORKFLOW RECOMMENDATION

Document how similar tasks should be handled:

```
## Workflow Recommendation

**Task Type:** {Description of task category}

**Current Assignment:** L0 → L1 → L2

**Recommendation:**
- Option A: Assign directly to L1 (L0 not suitable)
- Option B: Assign directly to L2 (requires architectural knowledge)
- Option C: Decompose all similar tasks into smaller specs
- Option D: Current cascade is appropriate, this was an edge case

**Rationale:**
{Why this recommendation}

**Template Updates Needed:**
- Planner: {What to change in task spec template}
- L0: {What patterns to emphasize}
- L1: {What to watch for}
```

## SUBMISSION FORMAT

```
L2 IMPLEMENTATION COMPLETE

Task: {TASK_ID} - {TASK_NAME}
Escalation Level: L2 (Top-Tier Coder)

Analysis Category: A / B / C / D

Implementation: COMPLETE / SPEC REWRITTEN / DECOMPOSED

Files Created/Modified:
- path/to/file1.go
- path/to/file2.go

Verification Results:
- Command: `go build ./...` → PASSED
- Command: `{spec verification}` → PASSED

Why Previous Tiers Failed:
{Summary of root cause}

Workflow Recommendation:
{See above}

L3 Architect Handoff Required: YES/NO
- If YES: This reveals a pattern that L3 Architect should review
```

## L3 ARCHITECT HANDOFF

If this reveals a workflow pattern (not just a one-off complex task), flag for L3 Architect review at next checkpoint.

---

**Begin by categorizing:** What type of problem is this?

**Then:** Implement/fix/decompose + recommend workflow changes
```

- [ ] **Step 2: Verify file exists**

```bash
grep "Category A\|Category B\|Category C\|Category D" docs/workflow/templates/05-l2-coder.md
```

Expected: Shows all 4 analysis categories

- [ ] **Step 3: Commit**

```bash
git add docs/workflow/templates/05-l2-coder.md
git commit -m "feat: add L2 top-tier coder final escalation template"
```

---

## Task 6: Create L3 Architect Prompt Templates

**Files:**
- Create: `docs/workflow/templates/06-l3-architect.md`

- [ ] **Step 1: Write L3 architect template (escalation + checkpoint)**

```markdown
# ROLE: PairAdmin L3 Architect

You are the Top-Tier Architect for the PairAdmin v2.0 project. Your job is to resolve escalated tasks that L2 Coder flagged for architectural review, conduct checkpoint reviews, and document workflow learnings.

## USE CASE 1: L3 Escalation Review

Use this section when L2 Coder flags a pattern requiring architectural decision.

### Input Context

**Task Specification:**
{INSERT: Original task spec}

**L2 Implementation:**
{INSERT: L2 implementation + workflow recommendation}

**L2 Flag:**
{INSERT: Why L2 flagged this for architect review}

### Your Responsibilities

1. **Make the architectural decision** - Choose the design approach
2. **Document the precedent** - This becomes a guideline for future tasks
3. **Update workflow if needed** - Change task assignment rules
4. **Advise the Planner** - How should future specs be written?

### Decision Record

```
## Architect Decision: {TASK_ID}

**Decision:** {Chosen approach}

**Alternatives Considered:**
1. {Alternative 1} - Rejected because: {reason}
2. {Alternative 2} - Rejected because: {reason}

**Precedent:**
{How this decision should guide future similar tasks}

**Workflow Impact:**
- Task assignment change: YES/NO
- Template update needed: YES/NO
- New guideline: {description}
```

---

## USE CASE 2: Checkpoint Review

Use this section at milestone completion (per QA_CHECKPOINTS.md).

### Checkpoint Context

**Milestone:** {MILESTONE_ID} - {MILESTONE_NAME}
**Tasks Completed:** {LIST of task IDs}
**Escalations:** {LIST of tasks that required escalation}
**QA Results:** {Checkpoint validation results}

### Review Dimensions

Analyze across these dimensions:

1. **Task Specification Quality** - Are specs becoming more precise?
2. **Model Performance Patterns** - Which task types consistently require escalation?
3. **Review Effectiveness** - Is L1/L2 catching issues early?
4. **Architecture Integrity** - Is codebase maintaining consistent patterns?
5. **Workflow Efficiency** - Where are the bottlenecks?

### Checkpoint Report

```
# Checkpoint Review: {MILESTONE_ID}

## Summary
- Tasks completed: {N}
- First-pass acceptance rate: {N}%
- Escalations: {N} ({%})
- Average review cycles: {N}

## What's Working Well

1. **{Strength 1}**
   - Evidence: {examples}
   - Continue: {specific practice}

2. **{Strength 2}**
   - Evidence: {examples}
   - Continue: {specific practice}

## Areas for Improvement

1. **{Issue 1}**
   - Evidence: {examples}
   - Impact: {cost/delay}
   - Recommendation: {specific change}

## Model Performance Analysis

**L0 Coder:**
- Strong at: {task types}
- Struggles with: {task types}
- Recommendation: {adjust assignments or templates}

**L1 Coder:**
- Escalation patterns: {observations}
- Recommendation: {guideline updates}

## Architecture Health

- Code consistency: {assessment}
- Interface stability: {assessment}
- Technical debt: {list any accumulating debt}
- Recommendation: {refactoring needs}

## Next Phase Guidance

**For Planner:**
1. {Specific guideline for next task specs}
2. {Context to include}
3. {Complexity limits to observe}

## Updated Workflow Guidelines

{Any changes to prompt templates or handoff protocols}

## Open Questions for Human Team

{Anything requiring human decision}
```

## OUTPUT

Write the checkpoint report to: `docs/checkpoints/checkpoint-{MILESTONE_ID}-review.md`

Commit this document to version control.

---

**Begin by reviewing:** All artifacts from this milestone/escalation

**Then produce:** The decision record or checkpoint report with actionable guidance
```

- [ ] **Step 2: Verify file exists**

```bash
grep "USE CASE\|Checkpoint Review\|L3 Escalation" docs/workflow/templates/06-l3-architect.md
```

Expected: Shows both use cases

- [ ] **Step 3: Commit**

```bash
git add docs/workflow/templates/06-l3-architect.md
git commit -m "feat: add L3 architect escalation and checkpoint templates"
```

---

## Task 7: Create Reviewer Prompt Template

**Files:**
- Create: `docs/workflow/templates/03-reviewer.md`

- [ ] **Step 1: Write reviewer template**

```markdown
# ROLE: PairAdmin Code Reviewer

You are the Mid-Tier Reviewer for the PairAdmin v2.0 project. Your job is to validate that coder submissions meet the task specification requirements.

## REVIEW CONTEXT

**Task Specification:**
{INSERT: Original task spec from Planner}

**Coder Submission:**
{INSERT: Coder's implementation + self-assessment}

**Review History:**
{INSERT: Previous review cycles if applicable, or "First review"}

## YOUR RESPONSIBILITIES

1. **Compare implementation against spec** - line by line
2. **Run verification mentally** - would the verification commands pass?
3. **Check acceptance criteria** - every single one
4. **Be strict but fair** - spec is the contract

## REVIEW CHECKLIST

### Specification Compliance
- [ ] All files specified in Outputs exist
- [ ] All function/struct names match spec
- [ ] All acceptance criteria are satisfied
- [ ] Verification commands would pass

### Code Quality
- [ ] Follows existing codebase patterns
- [ ] Consistent error handling
- [ ] No obvious bugs or race conditions
- [ ] Interfaces properly implemented

### Integration
- [ ] Does not break existing functionality
- [ ] Compatible with dependent tasks
- [ ] Cross-platform considerations addressed (if applicable)

## REVIEW DECISION

Choose one:

### ACCEPT
Implementation meets all spec requirements. Ready for next task.

### REJECT WITH COMMENTS
Implementation has gaps. Return to Coder with specific fixes needed.
- List each issue with file:line reference
- Explain what change is needed
- Coder will revise and resubmit

### ESCALATE TO L1
This is the 2nd rejection OR the issue is beyond L0 capability.
- Summarize what was attempted
- Explain why Coder struggled
- Note any spec ambiguities that contributed

## OUTPUT FORMAT

```
REVIEW DECISION: ACCEPT / REJECT / ESCALATE

Task: {TASK_ID} - {TASK_NAME}
Review Cycle: 1 / 2 / Escalation

Checklist Results:
- Specification Compliance: PASS / FAIL (details)
- Code Quality: PASS / FAIL (details)
- Integration: PASS / FAIL (details)

Issues Found:
1. {file.go:line} - {description} - {fix required}
2. ...

Decision Rationale:
{2-3 sentences explaining the decision}

Next Steps:
{Coder revises / Task complete / Escalating to L1}
```

## ESCALATION CRITERIA

Escalate to L1 when:
- This is the 2nd rejection of the same task
- The Coder consistently misunderstands the spec
- The spec itself appears to have gaps or ambiguities
- The task requires architectural decisions beyond Coder scope
- CGO/platform-specific issues are blocking progress

## BE FAIR BUT RIGOROUS

- If the spec is ambiguous, note it (don't penalize Coder)
- If the Coder followed spec but spec is wrong, flag the spec
- If the Coder added unrequested features, note it (even if helpful)
- Document patterns for L3 checkpoint review

---

**Begin by comparing:** Implementation against specification

**Then decide:** ACCEPT, REJECT, or ESCALATE
```

- [ ] **Step 2: Verify file exists**

```bash
wc -l docs/workflow/templates/03-reviewer.md
```

Expected: File has 100+ lines

- [ ] **Step 3: Commit**

```bash
git add docs/workflow/templates/03-reviewer.md
git commit -m "feat: add mid-tier reviewer prompt template"
```

---

## Task 8: Create Workflow Tracker Python Library

**Files:**
- Create: `docs/workflow/lib/workflow_tracker.py`

- [ ] **Step 1: Write workflow tracker module**

```python
#!/usr/bin/env python3
"""
Workflow Tracker for Multi-Tier LLM Development Cascade

Logs handoffs, escalations, and learnings to JSONL files for audit and analysis.
"""

import json
import os
from datetime import datetime
from pathlib import Path
from typing import Optional, Dict, Any, List


class WorkflowTracker:
    """Tracks workflow state across L0-L3 tier cascade."""
    
    def __init__(self, base_dir: str = "docs/workflow"):
        self.base_dir = Path(base_dir)
        self.handoffs_dir = self.base_dir / "handoffs"
        self.escalations_dir = self.base_dir / "escalations"
        self.learnings_dir = self.base_dir / "learnings"
        
        # Ensure directories exist
        self.handoffs_dir.mkdir(parents=True, exist_ok=True)
        self.escalations_dir.mkdir(parents=True, exist_ok=True)
        self.learnings_dir.mkdir(parents=True, exist_ok=True)
    
    def log_handoff(
        self,
        task_id: str,
        from_tier: str,
        to_tier: str,
        handoff_type: str,
        details: Dict[str, Any],
        timestamp: Optional[str] = None
    ) -> str:
        """
        Log a handoff between tiers.
        
        Args:
            task_id: Task identifier (e.g., "1.1")
            from_tier: Source tier (e.g., "Planner", "L0")
            to_tier: Destination tier (e.g., "L0", "Reviewer")
            handoff_type: Type of handoff (e.g., "Task Ready", "Review Ready")
            details: Handoff-specific data
            timestamp: ISO format timestamp (auto-generated if not provided)
        
        Returns:
            Path to the log file
        """
        if timestamp is None:
            timestamp = datetime.utcnow().isoformat() + "Z"
        
        log_entry = {
            "timestamp": timestamp,
            "task_id": task_id,
            "from_tier": from_tier,
            "to_tier": to_tier,
            "handoff_type": handoff_type,
            "details": details
        }
        
        # Create log filename
        safe_task_id = task_id.replace(".", "_")
        log_file = self.handoffs_dir / f"{safe_task_id}-{handoff_type.replace(' ', '-').lower()}.json"
        
        # Write log (overwrite if exists - each handoff type is unique per task)
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        
        return str(log_file)
    
    def log_escalation(
        self,
        task_id: str,
        from_tier: str,
        to_tier: str,
        escalation_level: str,
        reason: str,
        attempt_history: List[Dict[str, str]],
        timestamp: Optional[str] = None
    ) -> str:
        """
        Log an escalation event.
        
        Args:
            task_id: Task identifier
            from_tier: Escalating tier
            to_tier: Receiving tier
            escalation_level: "L1", "L2", or "L3"
            reason: Why this was escalated
            attempt_history: List of previous attempts with outcomes
            timestamp: ISO format timestamp
        
        Returns:
            Path to the log file
        """
        if timestamp is None:
            timestamp = datetime.utcnow().isoformat() + "Z"
        
        log_entry = {
            "timestamp": timestamp,
            "task_id": task_id,
            "escalation_level": escalation_level,
            "from_tier": from_tier,
            "to_tier": to_tier,
            "reason": reason,
            "attempt_history": attempt_history
        }
        
        # Create log filename
        safe_task_id = task_id.replace(".", "_")
        log_file = self.escalations_dir / f"{safe_task_id}-{escalation_level.lower()}.json"
        
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        
        return str(log_file)
    
    def log_learning(
        self,
        task_id: str,
        learning_type: str,
        category: str,
        content: Dict[str, Any],
        timestamp: Optional[str] = None
    ) -> str:
        """
        Log a learning document.
        
        Args:
            task_id: Task identifier
            learning_type: "Escalation Resolution" or "Checkpoint Review"
            category: Root cause category
            content: Learning content
            timestamp: ISO format timestamp
        
        Returns:
            Path to the log file
        """
        if timestamp is None:
            timestamp = datetime.utcnow().isoformat() + "Z"
        
        log_entry = {
            "timestamp": timestamp,
            "task_id": task_id,
            "learning_type": learning_type,
            "category": category,
            "content": content
        }
        
        # Create log filename
        safe_task_id = task_id.replace(".", "_")
        log_file = self.learnings_dir / f"{safe_task_id}-{learning_type.replace(' ', '-').lower()}.json"
        
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        
        return str(log_file)
    
    def get_task_history(self, task_id: str) -> Dict[str, Any]:
        """
        Get complete history for a task.
        
        Args:
            task_id: Task identifier
        
        Returns:
            Dictionary with handoffs, escalations, and learnings for the task
        """
        safe_task_id = task_id.replace(".", "_")
        
        history = {
            "task_id": task_id,
            "handoffs": [],
            "escalations": [],
            "learnings": []
        }
        
        # Find handoffs
        for log_file in self.handoffs_dir.glob(f"{safe_task_id}-*.json"):
            with open(log_file) as f:
                history["handoffs"].append(json.load(f))
        
        # Find escalations
        for log_file in self.escalations_dir.glob(f"{safe_task_id}-*.json"):
            with open(log_file) as f:
                history["escalations"].append(json.load(f))
        
        # Find learnings
        for log_file in self.learnings_dir.glob(f"{safe_task_id}-*.json"):
            with open(log_file) as f:
                history["learnings"].append(json.load(f))
        
        # Sort by timestamp
        history["handoffs"].sort(key=lambda x: x["timestamp"])
        history["escalations"].sort(key=lambda x: x["timestamp"])
        history["learnings"].sort(key=lambda x: x["timestamp"])
        
        return history
    
    def get_workflow_metrics(self) -> Dict[str, Any]:
        """
        Calculate workflow-wide metrics.
        
        Returns:
            Dictionary with aggregate metrics
        """
        # Count all handoffs
        handoff_count = len(list(self.handoffs_dir.glob("*.json")))
        
        # Count escalations by level
        escalation_counts = {"L1": 0, "L2": 0, "L3": 0}
        for log_file in self.escalations_dir.glob("*.json"):
            with open(log_file) as f:
                entry = json.load(f)
                level = entry.get("escalation_level", "unknown")
                if level in escalation_counts:
                    escalation_counts[level] += 1
        
        # Count learnings
        learning_count = len(list(self.learnings_dir.glob("*.json")))
        
        # Calculate escalation rate (approximate - would need task count)
        total_escalations = sum(escalation_counts.values())
        
        return {
            "total_handoffs": handoff_count,
            "total_escalations": total_escalations,
            "escalation_by_level": escalation_counts,
            "total_learnings": learning_count,
            "escalation_rate": f"{total_escalations}/{handoff_count} (approximate)"
        }


def main():
    """CLI entry point for workflow tracker."""
    import argparse
    
    parser = argparse.ArgumentParser(description="Workflow Tracker CLI")
    parser.add_argument("--command", choices=["history", "metrics", "log-handoff", "log-escalation"], required=True)
    parser.add_argument("--task-id", help="Task ID for history or logging")
    parser.add_argument("--from-tier", help="Source tier for logging")
    parser.add_argument("--to-tier", help="Destination tier for logging")
    parser.add_argument("--type", help="Handoff or escalation type")
    parser.add_argument("--level", help="Escalation level (L1/L2/L3)")
    parser.add_argument("--reason", help="Escalation reason")
    
    args = parser.parse_args()
    
    tracker = WorkflowTracker()
    
    if args.command == "history":
        if not args.task_id:
            print("Error: --task-id required for history command")
            return
        history = tracker.get_task_history(args.task_id)
        print(json.dumps(history, indent=2))
    
    elif args.command == "metrics":
        metrics = tracker.get_workflow_metrics()
        print(json.dumps(metrics, indent=2))
    
    elif args.command == "log-handoff":
        # Interactive or JSON input for handoff details
        print("Handoff logging - provide details via stdin JSON")
        details = json.loads(input())
        log_file = tracker.log_handoff(
            task_id=args.task_id,
            from_tier=args.from_tier,
            to_tier=args.to_tier,
            handoff_type=args.type,
            details=details
        )
        print(f"Logged to: {log_file}")
    
    elif args.command == "log-escalation":
        print("Escalation logging - provide attempt history via stdin JSON")
        attempt_history = json.loads(input())
        log_file = tracker.log_escalation(
            task_id=args.task_id,
            from_tier=args.from_tier,
            to_tier=args.to_tier,
            escalation_level=args.level,
            reason=args.reason,
            attempt_history=attempt_history
        )
        print(f"Logged to: {log_file}")


if __name__ == "__main__":
    main()
```

- [ ] **Step 2: Create requirements.txt for Python dependencies**

```
# docs/workflow/lib/requirements.txt
# No external dependencies - uses only Python stdlib
```

- [ ] **Step 3: Test the workflow tracker module**

```bash
cd docs/workflow/lib
python3 workflow_tracker.py --command metrics
```

Expected: Shows JSON metrics output with zeros (no data yet)

- [ ] **Step 4: Commit**

```bash
git add docs/workflow/lib/workflow_tracker.py docs/workflow/lib/requirements.txt
git commit -m "feat: add workflow tracker Python library"
```

---

## Task 9: Create Workflow CLI Wrapper Script

**Files:**
- Create: `docs/workflow/workflow.sh`

- [ ] **Step 1: Write CLI wrapper script**

```bash
#!/bin/bash
#
# Workflow CLI - Command-line interface for multi-tier development workflow
#
# Usage:
#   ./workflow.sh <command> [options]
#
# Commands:
#   status              Show workflow status and metrics
#   history <TASK_ID>   Show complete history for a task
#   templates           List available prompt templates
#   next-task           Show next pending task from IMPLEMENTATION_PLAN.md
#   log-handoff         Interactive handoff logging
#   log-escalation      Interactive escalation logging
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LIB_DIR="$SCRIPT_DIR/lib"
TRACKER="$LIB_DIR/workflow_tracker.py"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

cmd_status() {
    print_header "Workflow Status"
    echo ""
    python3 "$TRACKER" --command metrics
    echo ""
    echo "Template files:"
    ls -1 "$SCRIPT_DIR/templates/" 2>/dev/null || echo "  (no templates found)"
    echo ""
    echo "Recent handoffs:"
    ls -lt "$SCRIPT_DIR/handoffs/" 2>/dev/null | head -6 || echo "  (no handoffs logged)"
}

cmd_history() {
    local task_id="$1"
    if [ -z "$task_id" ]; then
        print_error "Task ID required"
        echo "Usage: $0 history <TASK_ID>"
        exit 1
    fi
    
    print_header "Task History: $task_id"
    echo ""
    python3 "$TRACKER" --command history --task-id "$task_id"
}

cmd_templates() {
    print_header "Available Prompt Templates"
    echo ""
    for template in "$SCRIPT_DIR/templates/"*.md; do
        if [ -f "$template" ]; then
            basename "$template"
            echo "  $(head -3 "$template" | tail -1)"
            echo ""
        fi
    done
}

cmd_next_task() {
    print_header "Next Pending Task"
    echo ""
    
    # Check if PROJECT_CHECKLIST.json exists
    if [ -f "PROJECT_CHECKLIST.json" ]; then
        python3 -c "
import json
with open('PROJECT_CHECKLIST.json') as f:
    data = json.load(f)
    pending = [t for t in data.get('tasks', []) if t.get('status') == 'pending']
    if pending:
        task = pending[0]
        print(f\"Task ID: {task.get('id')}\")
        print(f\"Title: {task.get('name')}\")
        print(f\"Dependencies: {task.get('dependencies', [])}\")
    else:
        print('No pending tasks found')
"
    else
        print_warning "PROJECT_CHECKLIST.json not found"
        echo "Showing first task from IMPLEMENTATION_PLAN.md instead:"
        echo ""
        grep -A 5 "#### Task 1.1" IMPLEMENTATION_PLAN.md | head -6
    fi
}

cmd_help() {
    print_header "Workflow CLI Help"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  status              Show workflow status and metrics"
    echo "  history <TASK_ID>   Show complete history for a task"
    echo "  templates           List available prompt templates"
    echo "  next-task           Show next pending task"
    echo "  log-handoff         Interactive handoff logging"
    echo "  log-escalation      Interactive escalation logging"
    echo "  help                Show this help message"
    echo ""
}

# Main command dispatcher
case "${1:-help}" in
    status)
        cmd_status
        ;;
    history)
        cmd_history "$2"
        ;;
    templates)
        cmd_templates
        ;;
    next-task)
        cmd_next_task
        ;;
    help|--help|-h)
        cmd_help
        ;;
    *)
        print_error "Unknown command: $1"
        cmd_help
        exit 1
        ;;
esac
```

- [ ] **Step 2: Make script executable**

```bash
chmod +x docs/workflow/workflow.sh
```

- [ ] **Step 3: Test the CLI**

```bash
cd /home/sblanken/code/paa5
./docs/workflow/workflow.sh help
./docs/workflow/workflow.sh templates
```

Expected: Shows help and template list

- [ ] **Step 4: Commit**

```bash
git add docs/workflow/workflow.sh
git commit -m "feat: add workflow CLI wrapper script"
```

---

## Task 10: Create Workflow README Documentation

**Files:**
- Create: `docs/workflow/README.md`

- [ ] **Step 1: Write workflow documentation**

```markdown
# Multi-Tier LLM Development Workflow

This directory contains the prompt templates, tools, and documentation for the L0-L3 multi-tier development cascade used to build PairAdmin v2.0.

## Overview

The workflow uses a tiered cascade of LLM capabilities:

| Level | Tier | Role | Models |
|-------|------|------|--------|
| **L0** | Lower-Tier | First implementation | DeepSeek Coder, StarCoder2 |
| **L1** | Mid-Tier Coder | Re-implementation after L0 fails 2x | Claude Sonnet, GPT-4-turbo |
| **L2** | Top-Tier Coder | Re-implementation after L1 fails | Claude Opus, GPT-4 |
| **L3** | Architect | Final escalation + checkpoint reviews | Claude Opus, GPT-4 |

## Quick Start

### For Planners (L0 Task Setup)

1. Read the planner template: `cat templates/01-planner.md`
2. Read IMPLEMENTATION_PLAN.md section for your task
3. Create task spec in `docs/tasks/`
4. Log handoff: `./workflow.sh log-handoff`

### For Coders (L0/L1/L2 Implementation)

1. Read your tier's template from `templates/`
2. Read the task spec from `docs/tasks/`
3. Implement according to spec
4. Submit using the template's submission format

### For Reviewers

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
```

- [ ] **Step 2: Verify file exists**

```bash
head -40 docs/workflow/README.md
```

Expected: Shows overview and quick start sections

- [ ] **Step 3: Commit**

```bash
git add docs/workflow/README.md
git commit -m "docs: add multi-tier workflow README"
```

---

## Task 11: Create Example Handoff Log Files

**Files:**
- Create: `docs/workflow/handoffs/example-handoff.json`
- Create: `docs/workflow/escalations/example-escalation.json`

- [ ] **Step 1: Create example handoff log**

```json
{
  "timestamp": "2026-03-30T10:00:00Z",
  "task_id": "1.1",
  "from_tier": "Planner",
  "to_tier": "L0",
  "handoff_type": "Task Ready",
  "details": {
    "spec_location": "docs/tasks/1.1-initialize-go-module.md",
    "dependencies_verified": [],
    "estimated_effort": "0.5 hours",
    "ready_for_coder": true
  }
}
```

- [ ] **Step 2: Create example escalation log**

```json
{
  "timestamp": "2026-03-30T14:00:00Z",
  "task_id": "2.5",
  "escalation_level": "L1",
  "from_tier": "Reviewer",
  "to_tier": "L1",
  "reason": "Second rejection - L0 consistently misunderstood CGO binding requirements",
  "attempt_history": [
    {
      "attempt": 1,
      "coder": "L0",
      "outcome": "REJECT",
      "issues": ["Missing CGO preamble", "Incorrect C function signatures"]
    },
    {
      "attempt": 2,
      "coder": "L0",
      "outcome": "REJECT",
      "issues": ["Still missing CGO preamble", "Memory management concerns"]
    }
  ]
}
```

- [ ] **Step 3: Verify files exist**

```bash
cat docs/workflow/handoffs/example-handoff.json
cat docs/workflow/escalations/example-escalation.json
```

- [ ] **Step 4: Commit**

```bash
git add docs/workflow/handoffs/example-handoff.json docs/workflow/escalations/example-escalation.json
git commit -m "docs: add example handoff and escalation logs"
```

---

## Task 12: Update Root README with Workflow Reference

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Read current README**

```bash
cat README.md
```

- [ ] **Step 2: Add workflow section to README**

Add after the "Key Documents" table:

```markdown
## Multi-Tier Workflow

This project uses a multi-tier LLM cascade (L0-L3) for cost-optimized development:

- **L0 (Lower-Tier):** First implementation attempt
- **L1 (Mid-Tier Coder):** Re-implementation after L0 fails 2x
- **L2 (Top-Tier Coder):** Re-implementation after L1 fails
- **L3 (Architect):** Final escalation + checkpoint reviews

**Documentation:** `docs/workflow/README.md`

**CLI:** `./docs/workflow/workflow.sh`

**Templates:** `docs/workflow/templates/`
```

- [ ] **Step 3: Verify update**

```bash
grep -A 10 "Multi-Tier Workflow" README.md
```

Expected: Shows new workflow section

- [ ] **Step 4: Commit**

```bash
git add README.md
git commit -m "docs: add multi-tier workflow reference to root README"
```

---

## Self-Review

After completing all tasks above, verify:

1. **Spec coverage:** All 6 prompt templates from design spec exist?
2. **No placeholders:** All template files have complete content?
3. **Type consistency:** Python module works without errors?
4. **CLI functional:** All commands return expected output?
5. **Documentation complete:** README explains workflow clearly?

Run verification:

```bash
# Verify all templates exist
ls -la docs/workflow/templates/

# Verify Python module imports
cd docs/workflow/lib && python3 -c "import workflow_tracker; print('OK')"

# Verify CLI works
./docs/workflow/workflow.sh status

# Verify directory structure
tree docs/workflow/
```

---

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-03-30-multi-tier-workflow-implementation.md`.

**Two execution options:**

**1. Subagent-Driven (recommended)** - Dispatch fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

**Which approach?**
