# Multi-Tier Workflow Learnings & Improvements

**Purpose:** This document captures learnings, improvements, and adjustments made to the multi-tier LLM development workflow during the PairAdmin v2.0 project. The goal is to refine the workflow for reuse in future projects.

**Document Status:** Living Document  
**Created:** March 30, 2026  
**Last Updated:** March 30, 2026  
**Project:** PairAdmin v2.0

---

## Workflow Configuration (Current)

| Level | Tier | Role | Models | Provider |
|-------|------|------|--------|----------|
| **L0** | Planner/Reviewer | Task specs, reviews | Qwen3.5 397B A17B, DeepSeek V3.2 | OpenRouter |
| **L0** | Coder | First implementation | Qwen3-Coder (local), Step 3.5 Flash | LM Studio @ 192.168.101.21, OpenRouter |
| **L1** | Coder | Re-implementation (escalation 1) | Grok 4.1 Fast | OpenRouter (xAI) |
| **L2** | Coder | Re-implementation (escalation 2) | MiniMax M2.7 | OpenRouter (MiniMax) |
| **L3** | Coder | Final escalation implementation | Claude Sonnet 4.6 | OpenRouter (Anthropic) |
| **L3** | Architect | Checkpoint reviews, final authority | Claude Opus 4.6 | OpenRouter (Anthropic) |

---

## Change Log

### 2026-03-30 - Task 1.3 First Execution (CRITICAL LEARNINGS)

**Changes:**
- First real multi-tier execution (L0-Planner → L0-Coder → L0-Reviewer)
- Orchestrator created and tested with OpenRouter + LM Studio
- 3-retry rule established per coder tier before escalation
- File read/write tools added to coder subagents

**Task 1.3 Results:**
- L0-Planner: ✅ Spec created (15.9s, ~$0.001)
- L0-Coder: ❌ LM Studio crashed mid-request
- L0-Coder (fallback): Human implemented
- L0-Reviewer: ✅ REJECT - caught critical syntax error (commas vs newlines in .gitignore)
- Final: ✅ Fixed and accepted

**Critical Issues Found:**
1. **LM Studio reliability** - Model crashed without detailed error info
2. **Syntax validation** - L0-Coder produced invalid .gitignore (comma-separated patterns)
3. **Reviewer value proven** - L0-Reviewer caught the syntax error before commit
4. **File tool gap** - Coder couldn't read existing files, leading to context gaps

**Adjustments Made:**
- Changed escalation rule: **3 retries per tier** (was 2) before escalation
- Added file read/write tools to orchestrator for coder subagents

### 2026-03-30 - Tool Integration Complete (RESOLVED)

**Changes:**
- Fixed missing `import time` in orchestrator.py (line 20)
- Fixed `execute_task()` to store full tool_result dict instead of just the list
- Updated L0-Coder prompt template with explicit `file_write()` syntax requirements

**Root Cause Analysis:**
1. LLMs were describing tool usage in markdown rather than outputting actual function calls
2. Prompt template showed examples but didn't emphasize exact syntax strongly enough
3. Common pitfalls section updated with explicit WRONG vs CORRECT examples

**Test Results:**
```
✓ Completed in 1.5s (attempt 1)
  Tools executed: 1
  All succeeded: True
  ✓ file_write('test_workflow.txt') - 36 bytes
```

**Key Learnings:**
1. **LLMs need explicit syntax instructions** - Saying "use file_write()" isn't enough; must show exact format
2. **Triple-quote format works best** - `file_write("path", """content""")` handles multi-line content
3. **Prompt iteration is critical** - First prompt version still allowed descriptive output
4. **Tool execution verification** - System now confirms files are actually created

**Files Modified:**
- `docs/workflow/lib/orchestrator.py` - Added import, fixed tool_results handling
- `docs/workflow/templates/02-l0-coder.md` - Enhanced tool usage instructions
- `docs/workflow/TODO_TOOL_INTEGRATION.md` - Marked as RESOLVED

**Status:** Tool integration now fully functional. LLMs can autonomously create files.

### 2026-03-30 - Tool Integration Complete (RESOLVED)

**Changes:**
- Fixed missing `import time` in orchestrator.py (line 20)
- Fixed `execute_task()` to store full tool_result dict instead of just the list
- Updated L0-Coder prompt template with explicit `file_write()` syntax requirements

**Root Cause Analysis:**
1. LLMs were describing tool usage in markdown rather than outputting actual function calls
2. Prompt template showed examples but didn't emphasize exact syntax strongly enough
3. Common pitfalls section updated with explicit WRONG vs CORRECT examples

**Test Results:**
```
✓ Completed in 1.5s (attempt 1)
  Tools executed: 1
  All succeeded: True
  ✓ file_write('test_workflow.txt') - 36 bytes
```

**Key Learnings:**
1. **LLMs need explicit syntax instructions** - Saying "use file_write()" isn't enough; must show exact format
2. **Triple-quote format works best** - `file_write("path", """content""")` handles multi-line content
3. **Prompt iteration is critical** - First prompt version still allowed descriptive output
4. **Tool execution verification** - System now confirms files are actually created

**Files Modified:**
- `docs/workflow/lib/orchestrator.py` - Added import, fixed tool_results handling
- `docs/workflow/templates/02-l0-coder.md` - Enhanced tool usage instructions
- `docs/workflow/TODO_TOOL_INTEGRATION.md` - Marked as RESOLVED

**Status:** Tool integration now fully functional. LLMs can autonomously create files.
- Added file read tool to all tiers (Planner, Reviewer, Architect)
- Simplified context format for L0-Coder (smaller prompts)
- Added retry logic with progressive simplification

**Cost Data:**
- L0-Planner (Qwen3.5 397B): ~$0.001 per task spec
- L0-Reviewer (Qwen3.5 397B): ~$0.001 per review
- L0-Coder (local): $0.00 but reliability concerns

---

### 2026-03-30 - Initial Workflow Created

**Changes:**
- Created initial multi-tier workflow design
- Established L0-L3 cascade with 6 prompt templates
- Added Python workflow tracker and CLI tool
- Integrated with PairAdmin v2.0 project

**Rationale:**
- Need cost-effective development process for complex Go + Wails application
- Want to learn model capability boundaries across tiers
- Workflow should be reusable for future projects

**Initial Configuration Decisions:**
- L0 uses local Qwen3-Coder for cost-free first attempts
- L1-L3 use progressively more capable OpenRouter models
- Claude Opus 4.6 as final authority (Architect role)
- Qwen3.5 397B for planning/review (strong reasoning, lower cost than Claude)

---

## Outstanding Issues & Resolutions

### Issue 1: LM Studio Reliability
**Status:** ✅ MITIGATED (3-retry logic handles transient failures)  
**Date Resolved:** 2026-03-30  
**Resolution:** Added automatic retry with context simplification. If LM Studio returns 500 error, retry up to 3 times with progressively smaller context.

### Issue 2: file_write Tool Not Being Used
**Status:** ⚠️ PARTIALLY RESOLVED (Template fixed, integration pending)  
**Date Identified:** 2026-03-30  
**Resolution:**
1. ✅ Updated `02-l0-coder.md` template with explicit file_write instructions and examples
2. ✅ Created `TOOL_USAGE_GUIDE.md` documenting proper tool syntax
3. ✅ Added ToolExecutor class to parse tool calls from LLM output
4. ✅ Updated reviewer template to verify file_write was actually called
5. ❌ **PENDING:** ToolExecutor not integrated into execute_task() flow
6. ❌ **PENDING:** LLM output not parsed for tool calls, files not actually created

**Workaround:** Human creates files manually until tool integration is complete.

**Tracking:** `docs/workflow/TODO_TOOL_INTEGRATION.md`

### Issue 3: File Tool Integration
**Status:** ✅ DOCUMENTED  
**Date Resolved:** 2026-03-30  
**Resolution:** Created comprehensive tool usage guide with:
- Correct syntax for file_read and file_write
- Common mistakes and how to avoid them
- Submission format requirements
- Reviewer verification checklist

---

## Learnings by Category

### Model Performance

#### L0 Coder (Qwen3-Coder 30B local via LM Studio)
**Date:** 2026-03-30  
**Observation:** Model crashed during API call without detailed error message. When working, produced invalid .gitignore syntax (comma-separated patterns instead of newlines). Context window limitations - struggled with large context payloads. On Task 1.4, first attempt returned 500 error, retry succeeded but claimed file was created when it wasn't (file_write tool not used).  
**Impact:** Task 1.3 implementation required human intervention. Task 1.4 required human to create file. Reviewer caught syntax error before commit.  
**Adjustment:** 
- Added 3-retry logic with progressive context simplification
- Added file read/write tools so coder doesn't need context in prompt
- Reduced default context size for L0-Coder calls
- Added fallback to L1-Coder after 3 failed attempts
- **FIXED:** Updated 02-l0-coder.md template with explicit file_write instructions and examples
- **FIXED:** Created TOOL_USAGE_GUIDE.md documenting proper tool syntax
- **FIXED:** Added ToolExecutor class to parse and execute tool calls from LLM output
- **PROCESS:** Reviewers now authorized to REJECT immediately if file_write not used
- **CRITICAL:** file_write tool not being used by L0-Coder - needs prompt engineering fix

#### L0 Planner (Qwen3.5 397B via OpenRouter)
**Date:** 2026-03-30  
**Observation:** Produces detailed, well-structured task specs. Cannot read files directly - needs content in prompt. Response time ~15-30s.  
**Impact:** Need to provide file content in context or add file read tool.  
**Adjustment:** Added file_read() tool to planner template.

#### L0 Reviewer (Qwen3.5 397B via OpenRouter)
**Date:** 2026-03-30  
**Observation:** Excellent at catching syntax errors and spec compliance issues. Caught critical .gitignore syntax error that would have broken git functionality. Response time ~30-40s.  
**Impact:** Validated the reviewer tier - caught issue human would have missed.  
**Adjustment:** Added file_read() tool to review template for direct file comparison.

#### L1 Coder (Grok 4.1 Fast via OpenRouter)
**Date:** Not yet tested  
**Observation:** N/A  
**Impact:** N/A  
**Adjustment:** N/A

#### L2 Coder (MiniMax M2.7 via OpenRouter)
**Date:** Not yet tested  
**Observation:** N/A  
**Impact:** N/A  
**Adjustment:** N/A

#### L3 Coder (Claude Sonnet 4.6 via OpenRouter)
**Date:** Not yet tested  
**Observation:** N/A  
**Impact:** N/A  
**Adjustment:** N/A

#### L3 Architect (Claude Opus 4.6 via OpenRouter)
**Date:** Not yet tested  
**Observation:** N/A  
**Impact:** N/A  
**Adjustment:** N/A

#### L2 Coder (MiniMax M2.7)
**Date:** _YYYY-MM-DD_  
**Observation:** _What did you notice?_  
**Impact:** _How did this affect the workflow?_  
**Adjustment:** _What changed as a result?_

#### L3 Coder (Claude Sonnet 4.6)
**Date:** _YYYY-MM-DD_  
**Observation:** _What did you notice?_  
**Impact:** _How did this affect the workflow?_  
**Adjustment:** _What changed as a result?_

#### L3 Architect (Claude Opus 4.6)
**Date:** _YYYY-MM-DD_  
**Observation:** _What did you notice?_  
**Impact:** _How did this affect the workflow?_  
**Adjustment:** _What changed as a result?_

#### Planner/Reviewer (Qwen3.5 397B, DeepSeek V3.2)
**Date:** _YYYY-MM-DD_  
**Observation:** _What did you notice?_  
**Impact:** _How did this affect the workflow?_  
**Adjustment:** _What changed as a result?_

---

### Escalation Patterns

#### Common Escalation Reasons

| Pattern | Frequency | Typical Resolution | Prevention |
|---------|-----------|-------------------|------------|
| _e.g., CGO binding misunderstandings_ | _High/Med/Low_ | _L1 or L2?_ | _Better spec templates?_ |
| _e.g., Cross-platform file paths_ | High/Med/Low | _L1 or L2?_ | _Better spec templates?_ |
| _e.g., Interface contract violations_ | High/Med/Low | _L1 or L2?_ | _Better spec templates?_ |

#### Escalation Statistics (Per Checkpoint)

**Checkpoint 1 (Tasks 1-15):**
- Total tasks: _N_
- L0 first-pass success: _N%_
- Escalated to L1: _N (_%)_
- Escalated to L2: _N (_%)_
- Escalated to L3: _N (_%)_
- Most common escalation reason: _reason_

**Checkpoint 2 (Tasks 16-28):**
- _Same format as above_

**Checkpoint 3 (Tasks 29-50):**
- _Same format as above_

---

### Template Effectiveness

#### Planner Template (01-planner.md)
**What's Working:**
- _List effective elements_

**What Needs Improvement:**
- _List gaps or ambiguities_

**Changes Made:**
- _Date:_ _Change description_

#### Coder Templates (02-l0-coder.md, 04-l1-coder.md, 05-l2-coder.md)
**What's Working:**
- _List effective elements_

**What Needs Improvement:**
- _List gaps or ambiguities_

**Changes Made:**
- _Date:_ _Change description_

#### Reviewer Template (03-reviewer.md)
**What's Working:**
- _List effective elements_

**What Needs Improvement:**
- _List gaps or ambiguities_

**Changes Made:**
- _Date:_ _Change description_

#### Architect Template (06-l3-architect.md)
**What's Working:**
- _List effective elements_

**What Needs Improvement:**
- _List gaps or ambiguities_

**Changes Made:**
- _Date:_ _Change description_

---

### Workflow Process

#### Handoff Protocols
**What's Working:**
- _e.g., JSON logging is clear and auditable_

**What Needs Improvement:**
- _e.g., Manual logging is tedious_

**Changes Made:**
- _Date:_ _Change description_

#### Review Cycles
**What's Working:**
- _e.g., Two-stage review catches issues early_

**What Needs Improvement:**
- _e.g., Review loops can stall_

**Changes Made:**
- _Date:_ _Change description_

#### Task Granularity
**What's Working:**
- _e.g., Atomic tasks fit context windows_

**What Needs Improvement:**
- _e.g., Some tasks still too large_

**Changes Made:**
- _Date:_ _Change description_

---

## Cost Analysis

### Token Usage by Tier (Per Checkpoint)

**Checkpoint 1:**
| Tier | Model | Tokens In | Tokens Out | Est. Cost |
|------|-------|-----------|------------|-----------|
| L0 Planner | Qwen3.5 397B | | | |
| L0 Coder | Qwen3-Coder (local) | | | $0.00 |
| L0 Reviewer | Qwen3.5 397B | | | |
| L1 Coder | Grok 4.1 Fast | | | |
| L2 Coder | MiniMax M2.7 | | | |
| L3 Architect | Claude Opus 4.6 | | | |
| **Total** | | | | **$X.XX** |

**Checkpoint 2:**
- _Same format_

**Checkpoint 3:**
- _Same format_

### Cost Per Task Type

| Task Type | Avg Cost | Notes |
|-----------|----------|-------|
| File creation (simple) | $X.XX | _Mostly L0_ |
| CGO integration | $X.XX | _Often escalates to L1/L2_ |
| Interface design | $X.XX | _Varies_ |
| Cross-platform adapters | $X.XX | _Often escalates_ |

---

## Reusability Assessment

### What Translates Well to Other Projects

**High Reusability:**
- _e.g., Prompt templates require minimal adjustment_
- _e.g., Escalation cascade works for any complex development_
- _e.g., Workflow tracker is project-agnostic_

**Medium Reusability:**
- _e.g., Model assignments may need adjustment per project_
- _e.g., Task granularity depends on domain_

**Low Reusability:**
- _e.g., Go/Wails-specific patterns in specs_
- _e.g., CGO-specific escalation patterns_

### Recommended Adjustments for Future Projects

**For Smaller Projects:**
- _e.g., Skip L2, go L1 → L3 directly_
- _e.g., Use only L0 + L3_

**For Larger Projects:**
- _e.g., Add L4 for cross-project coordination_
- _e.g., Multiple parallel cascades per subsystem_

**For Different Domains:**
- _Frontend-heavy:_ _Adjustments needed_
- _Data/ML:_ _Adjustments needed_
- _Infrastructure:_ _Adjustments needed_

---

## Open Questions

1. _What is the optimal number of tiers? Is 4 (L0-L3) right, or should we have 3 or 5?_
2. _Should Planner and Reviewer always be the same model, or split them?_
3. _Is the 2-strike escalation rule optimal, or should it be 1 or 3?_
4. _Should L0 Coder use local-only, or is the free Step 3.5 Flash worth the minimal cost?_
5. _How do we better predict which tasks will escalate before assignment?_

---

## Future Workflow Versions

### v1.1 (Planned)
**Changes:**
- _List planned improvements_

**Rationale:**
- _Why these changes_

**Target Date:**
- _When to implement_

### v2.0 (Future)
**Changes:**
- _Major redesigns_

**Rationale:**
- _Why these changes_

**Target Date:**
- _When to implement_

---

## Appendix: Model Comparison Notes

### Qwen3.5 397B A17B (L0 Planner/Reviewer)
**Strengths:**
- _e.g., Strong reasoning for spec writing_
- _e.g., Good at catching compliance issues_

**Weaknesses:**
- _e.g., Sometimes overly verbose_
- _e.g., Misses subtle Go patterns_

**Best For:**
- Task specification
- Compliance review

### Qwen3-Coder Local (L0 Coder)
**Strengths:**
- _e.g., Free_
- _e.g., Fast for simple tasks_

**Weaknesses:**
- _e.g., Misses edge cases_
- _e.g., Struggles with CGO_

**Best For:**
- Boilerplate
- Simple functions
- Test files

### Step 3.5 Flash (L0 Coder Fallback)
**Strengths:**
- _e.g., Free via OpenRouter_
- _e.g., Faster than local sometimes_

**Weaknesses:**
- _e.g., Smaller context window_

**Best For:**
- Tasks that exceed local model context

### Grok 4.1 Fast (L1 Coder)
**Strengths:**
- _e.g., Good debugging capability_
- _e.g., Fast response_

**Weaknesses:**
- _e.g., _

**Best For:**
- L0 failure recovery
- Multi-file changes

### MiniMax M2.7 (L2 Coder)
**Strengths:**
- _e.g., _
- _e.g., _

**Weaknesses:**
- _e.g., _

**Best For:**
- Complex implementation
- Spec rewriting

### Claude Sonnet 4.6 (L3 Coder)
**Strengths:**
- _e.g., _
- _e.g., _

**Weaknesses:**
- _e.g., _

**Best For:**
- Final escalation
- Security-critical code

### Claude Opus 4.6 (L3 Architect)
**Strengths:**
- _e.g., Highest capability_
- _e.g., Excellent pattern recognition_
- _e.g., Clear decision documentation_

**Weaknesses:**
- _e.g., Highest cost_
- _e.g., Slower response_

**Best For:**
- Checkpoint reviews
- Architectural decisions
- Template evolution
- Final authority on disputes

---

_This document is a living record. Update after each checkpoint review and whenever a significant learning or adjustment occurs._
