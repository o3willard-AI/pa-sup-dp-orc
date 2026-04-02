# Multi-Tier LLM Development Workflow Design

**Version:** 1.0  
**Date:** March 30, 2026  
**Status:** DRAFT  
**Author:** AI Development Team

---

## 1. Overview

This document defines the multi-tier LLM cascade workflow for developing PairAdmin v2.0. The workflow uses different model tiers for different development activities to optimize cost while maintaining quality.

### 1.1 Core Principles

1. **End-users always receive the best available model** - No cost optimization in the shipped product
2. **Development uses hierarchical cascade** - Tasks flow through tiers with clear handoff criteria
3. **Learning is explicit and documented** - Top-tier reviews produce actionable guidelines
4. **Task specs are the primary interface** - Well-formed specs enable lower-tier success

### 1.2 Model Tier Assignments

| Tier | Role | Models | Provider |
|------|------|--------|----------|
| **L0** | Planner, Reviewer | Qwen3.5 397B A17B, DeepSeek V3.2 | OpenRouter |
| **L0** | Coder | Qwen3-Coder (local), Step 3.5 Flash (free) | LM Studio @ 192.168.101.21, OpenRouter |
| **L1** | Coder | Grok 4.1 Fast | OpenRouter (xAI) |
| **L2** | Coder | MiniMax M2.7 | OpenRouter (MiniMax) |
| **L3** | Coder | Claude Sonnet 4.6 | OpenRouter (Anthropic) |
| **L3** | Architect | Claude Opus 4.6 | OpenRouter (Anthropic) |

### 1.3 Workflow Summary

```
┌─────────────────┐
│  MID-TIER       │
│  (Planner)      │
│  Creates Task   │
│  Specification  │
└────────┬────────┘
         │ Task Spec
         ▼
┌─────────────────┐
│  LOWER-TIER     │
│  (Coder)        │
│  Implements     │
│  Self-Reviews   │
│  (2 passes)     │
└────────┬────────┘
         │ Code + Self-Assessment
         ▼
┌─────────────────┐
│  MID-TIER       │
│  (Reviewer)     │
│  Validates      │
│  Accept/Reject  │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
 ACCEPT   REJECT (max 2 cycles)
    │         │
    │    ┌────▼────────────┐
    │    │  ESCALATION L1  │
    │    │  (Mid-Tier Coder)
    │    │  Re-implement   │
    │    │  Annotate       │
    │    └────┬────────────┘
    │         │
    │         ▼
    │    ┌─────────────────┐
    │    │  MID-TIER       │
    │    │  (Reviewer)     │
    │    │  Review Again   │
    │    └────┬────────────┘
    │         │
    │    ┌────┴────┐
    │    │         │
    │ ACCEPT  REJECT
    │    │         │
    │    │    ┌────▼────────────┐
    │    │    │  ESCALATION L2  │
    │    │    │  (Top-Tier Coder)
    │    │    │  Re-implement   │
    │    │    │  Or Rewrite Spec│
    │    │    └────┬────────────┘
    │    │         │
    │    │         ▼
    │    │    ┌─────────────────┐
    │    │    │  MID-TIER       │
    │    │    │  (Reviewer)     │
    │    │    │  Review Again   │
    │    │    └────┬────────────┘
    │    │         │
    │    │    ┌────┴────┐
    │    │    │         │
    │    │ ACCEPT  REJECT/FLAG
    │    │    │         │
    │    │    │    ┌────▼────────────┐
    │    │    │    │  ESCALATION L3  │
    │    │    │    │  (Architect)    │
    │    │    │    │  Final Review   │
    │    │    │    │  Pattern Analysis
    │    └────┴────│  Checkpoint     │
    │              └─────────────────┘
     ▼
 TASK COMPLETE
```

**Escalation Path:** L0 (Lower-Tier) → L1 (Mid-Tier Coder) → L2 (Top-Tier Coder) → L3 (Architect)

**Key Principle:** Each escalation tier brings greater capability AND documents why previous tiers failed for workflow learning.

### 1.4 Escalation Tier Capabilities

| Level | Tier | Models | Capabilities | Typical Use Cases |
|-------|------|--------|-------------|-------------------|
| **L0** | Planner/Reviewer | Qwen3.5 397B, DeepSeek V3.2 | Task decomposition, spec writing, review | Atomic task specs, compliance checks |
| **L0** | Coder | Qwen3-Coder (local), Step 3.5 Flash | Straightforward implementation, pattern-following | File creation, simple functions, tests, config |
| **L1** | Coder | Grok 4.1 Fast | Debugging failed implementations, multi-file changes | CGO integration, cross-platform adapters, error handling |
| **L2** | Coder | MiniMax M2.7 | Architectural implementation, spec rewriting, decomposition | Security-critical code, performance, novel patterns |
| **L3** | Coder | Claude Sonnet 4.6 | Final escalation, complex implementation | Implementation requiring highest coding capability |
| **L3** | Architect | Claude Opus 4.6 | Workflow analysis, pattern recognition, final authority | Checkpoint reviews, template evolution, technical debt decisions |

### 1.5 Provider Configuration

**Local LM Studio (L0 Coder):**
- URL: `http://192.168.101.21:1234/v1`
- Model: `Qwen3-coder`
- Use for: Cost-free first implementation attempts

**OpenRouter (All Other Tiers):**
- Base URL: `https://openrouter.ai/api/v1`
- L0 Planner/Reviewer: `qwen/qwen3.5-397b-a17b`, `deepseek/deepseek-v3.2`
- L0 Coder (fallback): `stepfun/step-3.5-flash` (free)
- L1 Coder: `xai/grok-4.1-fast`
- L2 Coder: `minimax/minimax-m2.7`
- L3 Coder: `anthropic/claude-sonnet-4.6`
- L3 Architect: `anthropic/claude-opus-4.6` (highest authority)

### 1.6 Tool Capabilities by Tier

| Tier | file_read | file_write | Notes |
|------|-----------|------------|-------|
| **L0 Planner** | ✅ | ❌ | Can read existing files for context |
| **L0 Coder** | ✅ | ✅ | Can read and write implementation files |
| **L0 Reviewer** | ✅ | ❌ | Can read files for comparison |
| **L1 Coder** | ✅ | ✅ | Full file access for re-implementation |
| **L2 Coder** | ✅ | ✅ | Full file access for complex fixes |
| **L3 Coder** | ✅ | ✅ | Full file access for final escalation |
| **L3 Architect** | ✅ | ❌ | Can read for review, writes to docs/ |

### 1.7 Retry Logic

**Each coder tier gets 3 attempts before escalation:**

| Attempt | Context Size | Temperature | Notes |
|---------|-------------|-------------|-------|
| **1st** | 100% | 0.7 | Full context, standard sampling |
| **2nd** | 70% | 0.7 | Truncated context, same temperature |
| **3rd** | 40% | 0.7 | Minimal context, focused prompt |

**Progressive Simplification Rationale:**
- If full context fails, the model may be overwhelmed
- Smaller context forces focus on core requirements
- Same temperature maintains creativity across attempts
- After 3 failures, problem likely needs higher-tier reasoning

### 1.6 Tool Capabilities by Tier

| Tier | file_read | file_write | Notes |
|------|-----------|------------|-------|
| **L0 Planner** | ✅ | ❌ | Can read existing files for context |
| **L0 Coder** | ✅ | ✅ | Can read and write implementation files |
| **L0 Reviewer** | ✅ | ❌ | Can read files for comparison |
| **L1 Coder** | ✅ | ✅ | Full file access for re-implementation |
| **L2 Coder** | ✅ | ✅ | Full file access for complex fixes |
| **L3 Coder** | ✅ | ✅ | Full file access for final escalation |
| **L3 Architect** | ✅ | ❌ | Can read for review, writes to docs/ |

### 1.7 Retry Logic

**Each coder tier gets 3 attempts before escalation:**

| Attempt | Context Size | Temperature | Notes |
|---------|-------------|-------------|-------|
| **1st** | 100% | 0.7 | Full context, standard sampling |
| **2nd** | 70% | 0.7 | Truncated context, same temperature |
| **3rd** | 40% | 0.7 | Minimal context, focused prompt |

**Progressive Simplification Rationale:**
- If full context fails, the model may be overwhelmed
- Smaller context forces focus on core requirements
- Same temperature maintains creativity across attempts
- After 3 failures, problem likely needs higher-tier reasoning

┌─────────────────┐
│  MID-TIER       │
│  (Planner)      │
│  Creates Task   │
│  Specification  │
└────────┬────────┘
         │ Task Spec
         ▼
┌─────────────────┐
│  LOWER-TIER     │
│  (Coder)        │
│  Implements     │
│  Self-Reviews   │
│  (2 passes)     │
└────────┬────────┘
         │ Code + Self-Assessment
         ▼
┌─────────────────┐
│  MID-TIER       │
│  (Reviewer)     │
│  Validates      │
│  Accept/Reject  │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
 ACCEPT   REJECT (max 2 cycles)
    │         │
    │    ┌────▼────────────┐
    │    │  ESCALATION L1  │
    │    │  (Mid-Tier Coder)
    │    │  Re-implement   │
    │    │  Annotate       │
    │    └────┬────────────┘
    │         │
    │         ▼
    │    ┌─────────────────┐
    │    │  MID-TIER       │
    │    │  (Reviewer)     │
    │    │  Review Again   │
    │    └────┬────────────┘
    │         │
    │    ┌────┴────┐
    │    │         │
    │ ACCEPT  REJECT
    │    │         │
    │    │    ┌────▼────────────┐
    │    │    │  ESCALATION L2  │
    │    │    │  (Top-Tier Coder)
    │    │    │  Re-implement   │
    │    │    │  Or Rewrite Spec│
    │    │    └────┬────────────┘
    │    │         │
    │    │         ▼
    │    │    ┌─────────────────┐
    │    │    │  MID-TIER       │
    │    │    │  (Reviewer)     │
    │    │    │  Review Again   │
    │    │    └────┬────────────┘
    │    │         │
    │    │    ┌────┴────┐
    │    │    │         │
    │    │ ACCEPT  REJECT/FLAG
    │    │    │    ┌────▼────────────┐
    │    │    │    │  ESCALATION L3  │
    │    │    │    │  (Architect)    │
    │    │    │    │  Final Review   │
    │    │    │    │  Pattern Analysis
    │    └────┴────│  Checkpoint     │
    │              └─────────────────┘
     ▼
 TASK COMPLETE
```

**Escalation Path:** L0 (Lower-Tier) → L1 (Mid-Tier Coder) → L2 (Top-Tier Coder) → L3 (Architect)

**Key Principle:** Each escalation tier brings greater capability AND documents why previous tiers failed for workflow learning.
┌─────────────────┐
│  MID-TIER       │
│  (Planner)      │
│  Creates Task   │
│  Specification  │
└────────┬────────┘
         │ Task Spec
         ▼
┌─────────────────┐
│  LOWER-TIER     │
│  (Coder)        │
│  Implements     │
│  Self-Reviews   │
└────────┬────────┘
         │ Code + Self-Assessment
         ▼
┌─────────────────┐
│  MID-TIER       │
│  (Reviewer)     │
│  Validates      │
│  Accept/Reject  │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
 ACCEPT   REJECT (max 2x)
    │         │
    │    ┌────┴────────────┐
    │    │  Escalate To:   │
    │    │  Mid-Tier Coder │
    │    └────┬────────────┘
    │         │
    │    ┌────┴────────────┐
    │    │  Still Failing? │
    │    └────┬────────────┘
    │         │
    │     YES │ NO (resolved)
    │         │
    │    ┌────┴────────────┐
    │    │  Top-Tier Coder │
    │    └────┬────────────┘
    │         │
    │    ┌────┴────────────┐
    │    │  Still Failing? │
    │    └────┬────────────┘
    │         │
    │     YES │ NO (resolved)
    │         │
    │    ┌────┴────────────┐
    └───►│  TOP-TIER       │
         │  (Architect)    │
         │  Final Escalation
         │  Checkpoint Review
         └─────────────────┘
```

---

## 2. Prompt Templates

### 2.1 Mid-Tier Planner Prompt Template

**Purpose:** Create atomic task specifications from the Implementation Plan

**When to Use:** At the start of each development task, after reading prior checkpoint learnings

**Template:**

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

---

### 2.2 Lower-Tier Coder Prompt Template

**Purpose:** Implement code from a task specification

**When to Use:** After receiving a task spec from Mid-Tier Planner

**Template:**

```markdown
# ROLE: PairAdmin Coder

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

---

### 2.3 Mid-Tier Reviewer Prompt Template

**Purpose:** Review coder submissions against task specifications

**When to Use:** After Lower-Tier Coder submits implementation

**Template:**

```markdown
# ROLE: PairAdmin Code Reviewer

You are the Mid-Tier Reviewer for the PairAdmin v2.0 project. Your job is to validate that coder submissions meet the task specification requirements.

## REVIEW CONTEXT

**Task Specification:**
{INSERT: Original task spec from Mid-Tier Planner}

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

### ESCALATE TO TOP-TIER
This is the 2nd rejection OR the issue is beyond Coder capability.
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
{Coder revises / Task complete / Escalating to Top-Tier}
```

## ESCALATION CRITERIA

Escalate to Top-Tier when:
- This is the 2nd rejection of the same task
- The Coder consistently misunderstands the spec
- The spec itself appears to have gaps or ambiguities
- The task requires architectural decisions beyond Coder scope
- CGO/platform-specific issues are blocking progress

## BE FAIR BUT RIGOROUS

- If the spec is ambiguous, note it (don't penalize Coder)
- If the Coder followed spec but spec is wrong, flag the spec
- If the Coder added unrequested features, note it (even if helpful)
- Document patterns for Top-Tier checkpoint review

---

**Begin by comparing:** Implementation against specification

**Then decide:** ACCEPT, REJECT, or ESCALATE
```

---

### 2.4 Mid-Tier Coder Prompt Template (Escalation Level 1)

**Purpose:** Re-implement tasks that Lower-Tier Coder failed twice

**When to Use:** When Mid-Tier Reviewer escalates after 2 rejections from Lower-Tier

**Template:**

```markdown
# ROLE: PairAdmin Mid-Tier Coder

You are the Mid-Tier Coder for the PairAdmin v2.0 project. Your job is to re-implement tasks that the Lower-Tier Coder could not complete successfully after two review cycles.

## ESCALATION CONTEXT

**Task Specification:**
{INSERT: Original task spec from Mid-Tier Planner}

**Lower-Tier Attempt 1:**
{INSERT: First implementation + review comments}

**Lower-Tier Attempt 2:**
{INSERT: Second implementation + review comments}

**Reviewer Escalation Notes:**
{INSERT: Why this was escalated to you}

## YOUR RESPONSIBILITIES

1. **Diagnose the failure** - Why did Lower-Tier fail twice?
2. **Re-implement correctly** - Produce working code that meets spec
3. **Document the gap** - What did Lower-Tier miss or misunderstand?
4. **Annotate for learning** - Add comments that help Lower-Tier learn

## IMPLEMENTATION PROCESS

### Step 1: Analyze Previous Failures
- Read the original spec completely
- Review both Lower-Tier attempts
- Read reviewer comments carefully
- Identify the root cause of failure:
  - Spec misunderstanding?
  - Missing knowledge/patterns?
  - Complexity too high?
  - Technical blocker (CGO, platform)?

### Step 2: Plan Your Approach
- Note what Lower-Tier got wrong
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
- Compare to Lower-Tier attempts - what's different?

## LEARNING NOTES

Document what Lower-Tier should learn:

```
## Why Lower-Tier Failed

**Root Cause:** {Spec misunderstanding / Knowledge gap / Technical complexity}

**Key Differences in My Approach:**
1. {What you did differently}
2. {Second difference}
3. {Third difference}

**Comments Added for Learning:**
- {file.go:line} - {what the comment explains}

**Recommendation:**
{Should Lower-Tier attempt similar tasks, or should this task type always come to Mid-Tier?}
```

## SUBMISSION FORMAT

```
MID-TIER IMPLEMENTATION COMPLETE

Task: {TASK_ID} - {TASK_NAME}
Escalation Level: 1 (Mid-Tier Coder)

Files Created/Modified:
- path/to/file1.go
- path/to/file2.go

Verification Results:
- Command: `go build ./...` → PASSED
- Command: `{spec verification}` → PASSED

Why Lower-Tier Failed:
{2-3 sentence summary}

Key Differences:
1. {Difference 1}
2. {Difference 2}

Learning Notes: {See above}

Ready for Review: YES
```

## WHEN TO FURTHER ESCALATE

If you also cannot complete this task, escalate to Top-Tier Coder with:
- Explanation of what you attempted
- Why you're blocked
- What additional knowledge/context is needed

---

**Begin by analyzing:** Why did Lower-Tier fail twice?

**Then implement:** Correctly with learning annotations
```

---

### 2.5 Top-Tier Coder Prompt Template (Escalation Level 2)

**Purpose:** Re-implement tasks that Mid-Tier Coder also struggled with

**When to Use:** When Mid-Tier Coder fails or identifies fundamental complexity issues

**Template:**

```markdown
# ROLE: PairAdmin Top-Tier Coder

You are the Top-Tier Coder for the PairAdmin v2.0 project. Your job is to resolve tasks that Mid-Tier Coder could not complete, or that require architectural-level understanding.

## ESCALATION CONTEXT

**Task Specification:**
{INSERT: Original task spec}

**Lower-Tier Attempts:**
{INSERT: Both attempts + review comments}

**Mid-Tier Attempt:**
{INSERT: Mid-Tier implementation + notes}

**Mid-Tier Escalation Notes:**
{INSERT: Why Mid-Tier escalated}

## YOUR RESPONSIBILITIES

1. **Diagnose the root cause** - Why did both tiers fail?
2. **Determine spec validity** - Is the task spec itself problematic?
3. **Implement or redesign** - Fix the code OR the spec
4. **Document for workflow** - Should this task type change tiers permanently?

## ANALYSIS

Determine which category this falls into:

### Category A: Implementation Complexity
The spec is correct, but the implementation requires expertise beyond Mid-Tier.
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

**Current Assignment:** Lower-Tier → Mid-Tier → Top-Tier

**Recommendation:**
- Option A: Assign directly to Mid-Tier (Lower-Tier not suitable)
- Option B: Assign directly to Top-Tier (requires architectural knowledge)
- Option C: Decompose all similar tasks into smaller specs
- Option D: Current cascade is appropriate, this was an edge case

**Rationale:**
{Why this recommendation}

**Template Updates Needed:**
- Planner: {What to change in task spec template}
- Lower-Tier: {What patterns to emphasize}
- Mid-Tier: {What to watch for}
```

## SUBMISSION FORMAT

```
TOP-TIER IMPLEMENTATION COMPLETE

Task: {TASK_ID} - {TASK_NAME}
Escalation Level: 2 (Top-Tier Coder)

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

Architect Handoff Required: YES/NO
- If YES: This reveals a pattern that Top-Tier Architect should review
```

## ARCHITECT HANDOFF

If this reveals a workflow pattern (not just a one-off complex task), flag for Top-Tier Architect review at next checkpoint.

---

**Begin by categorizing:** What type of problem is this?

**Then:** Implement/fix/decompose + recommend workflow changes
```

---

### 2.6 Top-Tier Architect Prompt Template (Escalation)

**Purpose:** Resolve tasks that Top-Tier Coder flagged for architectural review, and conduct checkpoint reviews

**When to Use:** When Top-Tier Coder identifies workflow patterns OR at milestone checkpoints

**Template:**

```markdown
# ROLE: PairAdmin Chief Architect

You are the Top-Tier Architect for the PairAdmin v2.0 project. Your job is to resolve escalated tasks, fix implementation issues, and document learnings to improve the workflow.

## ESCALATION CONTEXT

**Task Specification:**
{INSERT: Original task spec}

**Coder Attempt 1:**
{INSERT: First implementation + review comments}

**Coder Attempt 2:**
{INSERT: Second implementation + review comments}

**Reviewer Escalation Notes:**
{INSERT: Why this was escalated}

## YOUR RESPONSIBILITIES

1. **Diagnose the root cause** - Why did this task fail twice?
2. **Fix the implementation** - Produce working code
3. **Document the learning** - What should change to prevent recurrence?
4. **Advise the Planner** - How should future specs be written?

## ROOT CAUSE ANALYSIS

Determine which category this falls into:

### Spec Ambiguity
The task specification was unclear, incomplete, or impossible to follow.
**Fix:** Rewrite the spec template to prevent this ambiguity

### Model Limitation
The task requires capabilities the Lower-Tier model doesn't have.
**Fix:** Adjust task assignment (this task type needs Mid-Tier implementation)

### Complexity Mismatch
The task was too large for a single atomic spec.
**Fix:** Decompose into smaller subtasks

### Knowledge Gap
Missing context about existing codebase or patterns.
**Fix:** Improve context inclusion in spec template

### Technical Blocker
CGO, platform, or dependency issue blocking progress.
**Fix:** Document workaround or escalate to human

## IMPLEMENTATION

Produce the correct implementation:

1. Read the original spec
2. Understand what went wrong in previous attempts
3. Write working code that meets acceptance criteria
4. Add comments explaining non-obvious decisions

## LEARNING DOCUMENT

Write a learning entry for the checkpoint document:

```
## Learning: {Short title}

**Task:** {TASK_ID}
**Category:** Spec Ambiguity / Model Limitation / Complexity / Knowledge / Technical

**What Happened:**
{2-3 sentences describing the failure pattern}

**Root Cause:**
{Specific diagnosis}

**Fix Applied:**
{What you changed to resolve it}

**Prevention:**
{How to prevent this in future tasks}

**Guideline Update:**
{Specific change to Planner or Coder prompt templates}
```

## PLANNER ADVICE

Provide specific guidance to the Mid-Tier Planner:

```
To: Mid-Tier Planner
From: Top-Tier Architect
Re: Task {TASK_ID} Learnings

For future tasks of this type:

1. **Spec Structure:** {What to include/clarify}
2. **Context Needed:** {What existing code to reference}
3. **Complexity Limit:** {How to decompose similar tasks}
4. **Model Assignment:** {Should this task type go to Mid-Tier instead?}

Specifically for next task ({NEXT_TASK_ID}):
{Any special considerations}
```

## OUTPUT FORMAT

```
ESCALATION RESOLVED

Task: {TASK_ID} - {TASK_NAME}

Root Cause: {Category + explanation}

Implementation: COMPLETE
- Files modified: {list}
- Verification: {commands + results}

Learning Document: {See above}

Planner Advice: {See above}

Workflow Impact:
- Template changes needed: YES/NO
- Task assignment changes: YES/NO
- No impact: YES/NO
```

---

**Begin by diagnosing:** Why did this task fail twice?

**Then:** Fix implementation + document learnings
```

---

### 2.7 Top-Tier Architect Prompt Template (Checkpoint Review)

**Purpose:** Review completed milestones and produce checkpoint learnings

**When to Use:** After each QA Checkpoint milestone is complete

**Template:**

```markdown
# ROLE: PairAdmin Chief Architect - Checkpoint Review

You are the Top-Tier Architect conducting a checkpoint review for PairAdmin v2.0. Your job is to review all work from the milestone, identify patterns, and produce actionable guidance for the next phase.

## CHECKPOINT CONTEXT

**Milestone:** {MILESTONE_ID} - {MILESTONE_NAME}
**Tasks Completed:** {LIST of task IDs}
**Escalations:** {LIST of tasks that required escalation}
**QA Results:** {Checkpoint validation results}

## ARTIFACTS TO REVIEW

1. All task specifications from Mid-Tier Planner
2. All implementations from Lower-Tier Coder
3. All review decisions from Mid-Tier Reviewer
4. All escalation resolutions from previous Top-Tier reviews
5. QA checkpoint validation results

## REVIEW DIMENSIONS

### 1. Task Specification Quality
- Are specs becoming more precise over time?
- Which specs led to smooth implementation?
- Which specs caused confusion or rework?

### 2. Model Performance Patterns
- Which task types consistently require escalation?
- Where does Lower-Tier Coder excel?
- Where does Lower-Tier Coder struggle?

### 3. Review Effectiveness
- Is Mid-Tier Reviewer catching issues early?
- Are review comments actionable?
- Is the 2-strike escalation threshold appropriate?

### 4. Architecture Integrity
- Does the codebase maintain consistent patterns?
- Are interfaces stable across tasks?
- Is technical debt accumulating?

### 5. Workflow Efficiency
- How many review cycles per task on average?
- What % of tasks require escalation?
- Where are the bottlenecks?

## CHECKPOINT REPORT

Write a comprehensive report:

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

2. **{Issue 2}**
   - Evidence: {examples}
   - Impact: {cost/delay}
   - Recommendation: {specific change}

## Model Performance Analysis

**Lower-Tier Coder:**
- Strong at: {task types}
- Struggles with: {task types}
- Recommendation: {adjust assignments or templates}

**Mid-Tier Planner:**
- Spec quality trend: {improving/stable/declining}
- Common gaps: {patterns}
- Recommendation: {template updates}

**Mid-Tier Reviewer:**
- Review consistency: {assessment}
- Escalation appropriateness: {assessment}
- Recommendation: {guideline updates}

## Architecture Health

- Code consistency: {assessment}
- Interface stability: {assessment}
- Technical debt: {list any accumulating debt}
- Recommendation: {refactoring needs}

## Next Phase Guidance

**For Mid-Tier Planner:**
1. {Specific guideline for next task specs}
2. {Context to include}
3. {Complexity limits to observe}

**For Lower-Tier Coder:**
1. {Patterns to follow}
2. {Common mistakes to avoid}
3. {Focus areas}

**For Mid-Tier Reviewer:**
1. {What to watch for}
2. {When to escalate earlier}
3. {Review checklist updates}

## Updated Workflow Guidelines

{Any changes to prompt templates or handoff protocols}

## Open Questions for Human Team

{Anything requiring human decision}
```

## OUTPUT

Write the checkpoint report to: `docs/checkpoints/checkpoint-{MILESTONE_ID}-review.md`

Commit this document to version control.

---

**Begin by reviewing:** All artifacts from this milestone

**Then produce:** The checkpoint report with actionable guidance
```

---

## 3. Handoff Protocols

### 3.1 Planner → Coder Handoff

**Trigger:** Task specification is complete and saved

**Protocol:**

1. **Planner completes spec** at `docs/tasks/{TASK_ID}-{task-name}.md`
2. **Planner verifies spec quality** using checklist in template
3. **Planner outputs handoff message:**

```
HANDOFF: Task {TASK_ID} Ready for Implementation

Spec Location: docs/tasks/{TASK_ID}-{task-name}.md

Dependencies Verified:
- [ ] {DEPENDENCY_TASK_1} - COMPLETE
- [ ] {DEPENDENCY_TASK_2} - COMPLETE

Context Files:
- {list of existing files Coder should read}

Estimated Effort: {N} hours

Ready for Coder: YES
```

4. **Coder acknowledges receipt:**

```
ACKNOWLEDGED: Task {TASK_ID}

Spec Read: YES
Dependencies Verified: YES
Context Files Reviewed: YES

Beginning Implementation: {TIMESTAMP}
```

**Artifact:** Handoff logged in `docs/handoffs/{TASK_ID}-handoff.log`

---

### 3.2 Coder → Reviewer Handoff

**Trigger:** Implementation complete with self-review passes

**Protocol:**

1. **Coder completes implementation** and self-review
2. **Coder submits** using submission format from template
3. **Coder outputs handoff message:**

```
HANDOFF: Task {TASK_ID} Ready for Review

Implementation Complete: YES

Files Created/Modified:
- {file1.go}
- {file2.go}

Verification Results:
- {command1} → PASSED
- {command2} → PASSED

Self-Assessment:
- Confidence: {HIGH/MEDIUM/LOW}
- Notes: {any concerns}
- Technical Debt: {any shortcuts}

Review Cycles: 0 (first submission)

Ready for Reviewer: YES
```

4. **Reviewer acknowledges:**

```
ACKNOWLEDGED: Task {TASK_ID}

Submission Received: YES
Beginning Review: {TIMESTAMP}
Expected Completion: {TIMESTAMP + 2 hours}
```

**Artifact:** Handoff logged in `docs/handoffs/{TASK_ID}-review-handoff.log`

---

### 3.3 Reviewer → Coder Handoff (Rejection)

**Trigger:** Review identifies issues requiring fixes

**Protocol:**

1. **Reviewer completes review** with REJECT decision
2. **Reviewer outputs rejection:**

```
HANDOFF: Task {TASK_ID} Returned for Revision

Decision: REJECT (Cycle {N} of 2)

Issues Found:
1. {file.go:line} - {description}
   Fix: {specific change required}

2. {file.go:line} - {description}
   Fix: {specific change required}

Review Checklist:
- Specification Compliance: FAIL ({details})
- Code Quality: FAIL ({details})
- Integration: PASS/FAIL ({details})

Revision Deadline: {TIMESTAMP}

After Revision: Resubmit for second review
```

3. **Coder acknowledges:**

```
ACKNOWLEDGED: Task {TASK_ID} Revision

Issues Understood: YES
Revision Beginning: {TIMESTAMP}

Questions: {any clarification needed, or "None"}
```

**Artifact:** Handoff logged in `docs/handoffs/{TASK_ID}-revision-{N}.log`

---

### 3.4 Reviewer → Mid-Tier Coder Handoff (Escalation Level 1)

**Trigger:** Second rejection of Lower-Tier Coder implementation

**Protocol:**

1. **Reviewer decides to escalate to Mid-Tier Coder**
2. **Reviewer outputs escalation:**

```
ESCALATION L1: Task {TASK_ID} to Mid-Tier Coder

Reason: Second rejection of Lower-Tier implementation

Task Spec: docs/tasks/{TASK_ID}-{task-name}.md

Attempt History:
- Attempt 1 (Lower-Tier): {summary of what was wrong}
- Attempt 2 (Lower-Tier): {summary of what was wrong}

Reviewer Analysis:
- Root cause appears to be: {Spec ambiguity / Model limitation / Complexity / Knowledge gap}
- Lower-Tier effort: {Adequate / Struggling / Confused}
- Spec clarity: {Clear / Ambiguous / Incomplete}

Recommendation for Mid-Tier Coder:
{What you think Mid-Tier should focus on}

Urgency: {Normal / Blocking other tasks}
```

3. **Mid-Tier Coder acknowledges:**

```
ACKNOWLEDGED: Task {TASK_ID} Escalation L1

Escalation Received: YES
Beginning Review: {TIMESTAMP}
Expected Resolution: {TIMESTAMP + 4 hours}

Initial Assessment:
- Will re-implement: YES
- Anticipated complexity: {Low / Medium / High}
- May require further escalation: YES/NO/MAYBE
```

**Artifact:** Handoff logged in `docs/escalations/{TASK_ID}-escalation-l1.log`

---

### 3.5 Mid-Tier Coder → Reviewer Handoff (L1 Resolution)

**Trigger:** Mid-Tier Coder completes re-implementation

**Protocol:**

1. **Mid-Tier Coder completes implementation**
2. **Mid-Tier Coder outputs handoff:**

```
HANDOFF: Task {TASK_ID} Ready for Review (L1 Resolution)

Implementation Complete: YES

Files Created/Modified:
- {file1.go}
- {file2.go}

Verification Results:
- {command1} → PASSED
- {command2} → PASSED

Why Lower-Tier Failed:
{2-3 sentence summary}

Key Differences in My Approach:
1. {Difference 1}
2. {Difference 2}

Self-Assessment:
- Confidence: {HIGH/MEDIUM/LOW}
- Ready for Review: YES

Recommendation:
{Should similar tasks go to Mid-Tier directly in future?}
```

3. **Reviewer acknowledges:**

```
ACKNOWLEDGED: Task {TASK_ID} L1 Resolution

Submission Received: YES
Beginning Review: {TIMESTAMP}
Expected Completion: {TIMESTAMP + 2 hours}
```

**Artifact:** Handoff logged in `docs/handoffs/{TASK_ID}-l1-resolution.log`

---

### 3.6 Mid-Tier Coder → Top-Tier Coder Handoff (Escalation Level 2)

**Trigger:** Mid-Tier Coder also cannot complete the task

**Protocol:**

1. **Mid-Tier Coder decides to escalate to Top-Tier**
2. **Mid-Tier Coder outputs escalation:**

```
ESCALATION L2: Task {TASK_ID} to Top-Tier Coder

Reason: {Cannot complete despite Mid-Tier capabilities}

Task Spec: docs/tasks/{TASK_ID}-{task-name}.md

Attempt History:
- Attempt 1 (Lower-Tier): {summary}
- Attempt 2 (Lower-Tier): {summary}
- Attempt 3 (Mid-Tier): {summary of what you tried}

Mid-Tier Analysis:
- Root cause: {Specific technical or architectural blocker}
- What I attempted: {Your approach}
- Why I'm blocked: {Specific gap}

Recommendation for Top-Tier:
{What you think Top-Tier should do}

Urgency: {Normal / Blocking other tasks}
```

3. **Top-Tier Coder acknowledges:**

```
ACKNOWLEDGED: Task {TASK_ID} Escalation L2

Escalation Received: YES
Beginning Review: {TIMESTAMP}
Expected Resolution: {TIMESTAMP + 8 hours}

Initial Assessment:
- Will re-implement: YES/NO
- May need to rewrite spec: YES/NO/MAYBE
- May need to decompose: YES/NO/MAYBE
```

**Artifact:** Handoff logged in `docs/escalations/{TASK_ID}-escalation-l2.log`

---

### 3.7 Top-Tier Coder → Reviewer Handoff (L2 Resolution)

**Trigger:** Top-Tier Coder completes re-implementation

**Protocol:**

1. **Top-Tier Coder completes implementation**
2. **Top-Tier Coder outputs handoff:**

```
HANDOFF: Task {TASK_ID} Ready for Review (L2 Resolution)

Implementation Complete: YES / SPEC REWRITTEN / DECOMPOSED

Files Created/Modified:
- {file1.go}
- {file2.go}

Verification Results:
- {command1} → PASSED
- {command2} → PASSED

Analysis Category: A (Implementation) / B (Spec) / C (Decomposition) / D (Architectural)

Why Previous Tiers Failed:
{Root cause summary}

Workflow Recommendation:
{Should similar tasks be reassigned to different tier?}

Architect Handoff Required: YES/NO
```

3. **Reviewer acknowledges:**

```
ACKNOWLEDGED: Task {TASK_ID} L2 Resolution

Submission Received: YES
Beginning Review: {TIMESTAMP}
Expected Completion: {TIMESTAMP + 2 hours}
```

**Artifact:** Handoff logged in `docs/handoffs/{TASK_ID}-l2-resolution.log`

---

### 3.8 Top-Tier Coder → Architect Handoff (Pattern Flag)

**Trigger:** Top-Tier Coder identifies workflow pattern requiring architectural review

**Protocol:**

1. **Top-Tier Coder flags pattern for Architect**
2. **Top-Tier Coder outputs handoff:**

```
PATTERN FLAG: Task {TASK_ID} for Architect Review

Pattern Type: {Task type that needs reassignment}

Current Assignment: Lower-Tier → Mid-Tier → Top-Tier

Recommended Assignment:
- Option A: Direct to Mid-Tier
- Option B: Direct to Top-Tier
- Option C: Decompose all similar tasks
- Option D: No change (edge case)

Rationale:
{Why this pattern matters}

Affected Future Tasks:
- {TASK_ID_X}
- {TASK_ID_Y}

Architect Action Needed:
{What should Architect decide at checkpoint}
```

3. **Architect acknowledges:**

```
ACKNOWLEDGED: Pattern Flag for Task {TASK_ID}

Flag Received: YES
Will Review At: Checkpoint {MILESTONE_ID}

Interim Guidance:
{Any immediate guidance for Planner}
```

**Artifact:** Handoff logged in `docs/escalations/{TASK_ID}-pattern-flag.log`

---

### 3.9 Reviewer → Architect Handoff (Final Escalation)

**Trigger:** Top-Tier Coder implementation also fails OR task requires architectural decision

---

### 3.10 Architect → Planner Handoff (Learning)

**Trigger:** Escalation resolved or checkpoint complete

**Protocol:**

1. **Architect completes resolution** or checkpoint review
2. **Architect outputs learning:**

```
LEARNING HANDOFF: Task {TASK_ID} / Checkpoint {MILESTONE_ID}

Type: {Escalation Resolution / Checkpoint Review}

Document Location: docs/learnings/{document-name}.md

Key Takeaways:
1. {Most important learning}
2. {Second most important}
3. {Third most important}

Template Changes Required:
- Planner: {YES/NO - what to change}
- Coder: {YES/NO - what to change}
- Reviewer: {YES/NO - what to change}

Task Assignment Changes:
- {Task types that should shift tiers}

Next Task Guidance:
- {Specific advice for next task spec}

Planner Action Items:
1. {Action 1}
2. {Action 2}
```

3. **Planner acknowledges:**

```
ACKNOWLEDGED: Learning Handoff

Document Read: YES
Guidance Understood: YES

Applying to Next Task ({NEXT_TASK_ID}):
- {Specific changes being made}

Template Updates:
- Will update: {which templates}
- After current task: YES/NO
```

**Artifact:** Handoff logged in `docs/handoffs/learning-{TASK_ID}.log`

---

## 4. Artifact Tracking

### 4.1 Directory Structure

```
docs/
├── tasks/
│   ├── 1.1-initialize-go-module.md
│   ├── 1.2-scaffold-wails.md
│   └── ...
├── handoffs/
│   ├── 1.1-handoff.log
│   ├── 1.1-review-handoff.log
│   └── ...
├── escalations/
│   ├── {TASK_ID}-escalation.log
│   └── ...
├── learnings/
│   ├── escalation-{TASK_ID}-resolution.md
│   └── checkpoint-{MILESTONE_ID}-review.md
├── checkpoints/
│   ├── checkpoint-1-review.md
│   └── ...
└── workflow/
    ├── planner-template.md
    ├── coder-template.md
    ├── reviewer-template.md
    └── architect-template.md
```

### 4.2 Version Control

All artifacts committed to version control:

```bash
# After each task completion
git add docs/tasks/{TASK_ID}*.md
git add docs/handoffs/{TASK_ID}*.log
git commit -m "Task {TASK_ID}: {task name} - {status}"

# After each checkpoint
git add docs/checkpoints/
git add docs/learnings/
git commit -m "Checkpoint {MILESTONE_ID}: Review and learnings"
```

### 4.3 Status Tracking

Task status tracked in `PROJECT_CHECKLIST.json`:

```json
{
  "tasks": [
    {
      "id": "1.1",
      "name": "Initialize Go module",
      "status": "completed",
      "review_cycles": 1,
      "escalated": false,
      "completed_date": "2026-03-30"
    }
  ]
}
```

---

## 5. Escalation Criteria

### 5.1 Automatic Escalation Triggers

**Escalate to Mid-Tier Coder (Level 1) when:**
1. **Third Rejection:** L0 Coder fails after 3 retry attempts (with progressive context simplification)
2. **Model Limitation:** L0 demonstrates inability despite clear spec and file tool access

**Escalate to Top-Tier Coder (Level 2) when:**
1. **L1 Exhausted:** L1 Coder fails after 3 retry attempts
2. **Spec Ambiguity Detected:** L1 identifies spec as root cause
3. **Complexity Mismatch:** Task requires more than atomic implementation

**Escalate to Top-Tier Architect when:**
1. **L2 Exhausted:** L2 Coder fails after 3 retry attempts  
2. **Top-Tier Coder Flags:** Implementation reveals workflow pattern issue
3. **Architectural Decision:** Task requires design choices affecting multiple components
4. **Cross-Task Impact:** Issue affects multiple dependent tasks
5. **Platform Blocker:** CGO or OS-specific issue blocking all progress
6. **All Tiers Failed:** L0, L1, and L2 all exhausted retries

### 5.2 Escalation Decision Tree

```
Task Submission (Lower-Tier Coder)
       ↓
[Reviewer Evaluation]
       ↓
┌──────────────────────┐
│  Meets Spec?         │
└──────┬───────────────┘
       │
   YES │ NO
       │  ↓
       │  ┌──────────────────────┐
       │  │  First Rejection?    │
       │  └──────┬───────────────┘
       │         │
       │     YES │ NO (Second)
       │         │  ↓
       │         │  ┌──────────────────────┐
       │         │  │  Return to Coder     │
       │         │  │  (Revision Cycle)    │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         │  │  Second Rejection?   │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │     YES │ NO (resolved)
       │         │         │
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         │  │  ESCALATE L1         │
       │         │  │  (Mid-Tier Coder)    │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │         ↓
       │         │  [Mid-Tier Implements]
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         │  │  Review Again        │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │     PASS│ FAIL
       │         │         │
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         │  │  ESCALATE L2         │
       │         │  │  (Top-Tier Coder)    │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │         ↓
       │         │  [Top-Tier Implements]
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         │  │  Review Again        │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │     PASS│ FAIL
       │         │         │
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         └─►│  ESCALATE L3         │
       │            │  (Architect)         │
       │            └──────────────────────┘
       ↓
   ACCEPT
```
Task Submission (Lower-Tier Coder)
       ↓
[Reviewer Evaluation]
       ↓
┌──────────────────────┐
│  Meets Spec?         │
└──────┬───────────────┘
       │
   YES │ NO
       │  ↓
       │  ┌──────────────────────┐
       │  │  First Rejection?    │
       │  └──────┬───────────────┘
       │         │
       │     YES │ NO (Second)
       │         │  ↓
       │         │  ┌──────────────────────┐
       │         │  │  Return to Coder     │
       │         │  │  (Revision Cycle)    │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         │  │  Second Rejection?   │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │     YES │ NO (resolved)
       │         │         │
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         │  │  ESCALATE L1         │
       │         │  │  (Mid-Tier Coder)    │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │         ↓
       │         │  [Mid-Tier Implements]
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         │  │  Review Again        │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │     PASS│ FAIL
       │         │         │
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         │  │  ESCALATE L2         │
       │         │  │  (Top-Tier Coder)    │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │         ↓
       │         │  [Top-Tier Implements]
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         │  │  Review Again        │
       │         │  └──────┬───────────────┘
       │         │         │
       │         │     PASS│ FAIL
       │         │         │
       │         │         ↓
       │         │  ┌──────────────────────┐
       │         └─►│  ESCALATE L3         │
       │            │  (Architect)         │
       │            └──────────────────────┘
       ↓
   ACCEPT
```

### 5.3 Escalation Response SLA

| Escalation Level | Priority | Response Time | Resolution Time |
|-----------------|----------|--------------|-----------------|
| **L1 (Mid-Tier Coder)** | Normal | 2 hours | 8 hours |
| **L2 (Top-Tier Coder)** | Normal | 4 hours | 24 hours |
| **L3 (Architect)** | Normal | 8 hours | 48 hours |
| **Any Level - Blocking** | High | 1 hour | 50% reduction |
| **Any Level - Critical** | Critical | 30 minutes | 4 hours |

---

## 6. Quality Gates

### 6.1 Spec Quality Gate (Planner)

Before handing off to Coder, verify:

- [ ] Task ID matches Implementation Plan
- [ ] All dependencies are marked complete
- [ ] File paths are absolute and correct
- [ ] Function/struct names match existing conventions
- [ ] Verification commands are executable
- [ ] Acceptance criteria are testable
- [ ] Constraints/gotchas are documented

### 6.2 Code Quality Gate (Coder)

Before submitting for review, verify:

- [ ] All spec requirements implemented
- [ ] Code compiles without warnings
- [ ] Verification commands pass
- [ ] Self-review pass 1 complete (spec compliance)
- [ ] Self-review pass 2 complete (code quality)
- [ ] No unrequested features added
- [ ] Self-assessment completed

### 6.3 Review Quality Gate (Reviewer)

Before finalizing decision, verify:

- [ ] Compared every spec requirement to implementation
- [ ] Ran verification mentally (or actually)
- [ ] Issues have file:line references
- [ ] Fixes requested are specific and actionable
- [ ] Decision rationale is clear
- [ ] Escalation criteria applied correctly

### 6.4 Escalation Quality Gate (Architect)

Before closing escalation, verify:

- [ ] Root cause diagnosed (not just symptoms)
- [ ] Implementation fixed and verified
- [ ] Learning document written
- [ ] Planner advice provided
- [ ] Workflow impact assessed

---

## 7. Metrics & Continuous Improvement

### 7.1 Tracked Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| First-pass acceptance rate | >70% | Tasks accepted on first review / Total tasks |
| Escalation rate | <15% | Escalated tasks / Total tasks |
| Average review cycles | <1.5 | Total reviews / Total tasks |
| Spec rework rate | <10% | Specs revised after handoff / Total specs |
| Checkpoint velocity | Stable | Tasks per week per phase |

### 7.2 Improvement Cadence

- **Per Task:** Coder self-review improves output
- **Per 5 Tasks:** Planner reviews acceptance rate, adjusts spec detail
- **Per Milestone:** Top-Tier checkpoint review produces guidelines
- **Per Phase:** Human team reviews metrics, approves template changes

### 7.3 Template Evolution

Templates are living documents. Update when:

- Same issue causes 3+ escalations
- Checkpoint review identifies pattern
- Human team approves change
- Model capabilities change (new models available)

All template changes logged in `docs/workflow/CHANGELOG.md`

---

## 8. Quick Reference

### 8.1 Role Summary

| Role | Tier | Models | Primary Output | Key Template |
|------|------|--------|---------------|--------------|
| **Planner** | L0 | Qwen3.5 397B, DeepSeek V3.2 | Task specifications | Section 2.1 |
| **Coder** | L0 | Qwen3-Coder (local), Step 3.5 Flash | Implementations (first attempt) | Section 2.2 |
| **Reviewer** | L0 | Qwen3.5 397B, DeepSeek V3.2 | Review decisions | Section 2.3 |
| **Coder** | L1 | Grok 4.1 Fast | Re-implementation (L1 escalation) | Section 2.4 |
| **Coder** | L2 | MiniMax M2.7 | Re-implementation (L2 escalation) | Section 2.5 |
| **Coder** | L3 | Claude Sonnet 4.6 | L3 escalation implementation | Section 2.6 |
| **Architect** | L3 | Claude Opus 4.6 | Checkpoint reviews, final decisions | Section 2.7 |

### 8.2 Handoff Summary

| Handoff | From | To | Models | Trigger |
|---------|------|-----|--------|---------|
| Task Ready | Planner (L0) | Coder (L0) | Qwen3.5 → Qwen3-Coder/Step Flash | Spec complete |
| Review Ready | Coder (L0) | Reviewer (L0) | Qwen3-Coder → Qwen3.5 | Implementation complete |
| Revision | Reviewer (L0) | Coder (L0) | Qwen3.5 → Qwen3-Coder | Rejection (cycle 1) |
| Escalation L1 | Reviewer (L0) | Coder (L1) | Qwen3.5 → Grok 4.1 Fast | Rejection (cycle 2) |
| L1 Resolution | Coder (L1) | Reviewer (L0) | Grok 4.1 → Qwen3.5 | Re-implementation complete |
| Escalation L2 | Coder (L1) | Coder (L2) | Grok 4.1 → MiniMax M2.7 | L1 cannot complete |
| L2 Resolution | Coder (L2) | Reviewer (L0) | MiniMax M2.7 → Qwen3.5 | Re-implementation complete |
| Pattern Flag | Coder (L2) | Architect (L3) | MiniMax M2.7 → Claude Opus 4.6 | Workflow pattern identified |
| Escalation L3 | Reviewer (L0) | Architect (L3) | Qwen3.5 → Claude Opus 4.6 | All coder tiers failed |
| Learning | Architect (L3) | Planner (L0) | Claude Opus 4.6 → Qwen3.5 | Resolution/Checkpoint |

### 8.3 File Locations

| Artifact | Location |
|----------|----------|
| Task specs | `docs/tasks/` |
| Handoff logs | `docs/handoffs/` |
| Escalations | `docs/escalations/` |
| Learnings | `docs/learnings/` |
| Checkpoints | `docs/checkpoints/` |
| Templates | `docs/workflow/` |

---

## Appendix A: Example Task Flow

See `docs/examples/complete-task-flow-example.md` for a full walkthrough of Task 1.1 through the entire cascade.

## Appendix B: Prompt Template Quick-Copy

All templates available as copy-paste files in `docs/workflow/templates/` directory.

---

*This document is a living specification. Update after each checkpoint review based on Top-Tier recommendations.*
