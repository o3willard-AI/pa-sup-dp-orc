# Final Release Checklist Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create release checklist document, run final verification tests, and commit changes.

**Architecture:** Simple file creation and script execution tasks. No architectural changes.

**Tech Stack:** Go, Wails, shell scripts.

---

### Task 1: Create release checklist

**Files:**
- Create: `RELEASE_CHECKLIST.md`

- [ ] **Step 1: Create RELEASE_CHECKLIST.md file**

```markdown
# PairAdmin v2.0 Release Checklist

## Pre‑Release
- [ ] All unit tests pass (`go test ./internal/...`)
- [ ] Integration tests pass (`./scripts/test-integration.sh`)
- [ ] End‑to‑end tests pass (`./scripts/test-e2e.sh`)
- [ ] No known critical bugs
- [ ] Documentation updated (user guide, installation)
- [ ] Version number updated in `wails.json`

## Packaging
- [ ] macOS `.dmg` builds successfully
- [ ] Windows `.msi` builds successfully
- [ ] Linux packages (`.deb`, `.rpm`, `.AppImage`) build successfully
- [ ] Installers tested on clean VMs

## Release
- [ ] Create GitHub release with changelog
- [ ] Upload all platform binaries
- [ ] Update website/download page
- [ ] Announce on social channels

## Post‑Release
- [ ] Monitor error logs
- [ ] Gather user feedback
- [ ] Plan next iteration
```

- [ ] **Step 2: Verify file created**

Run: `ls -la RELEASE_CHECKLIST.md`
Expected: File exists with correct permissions

### Task 2: Run final verification

**Files:**
- Test: `scripts/test-integration.sh`
- Test: `scripts/test-e2e.sh`

- [ ] **Step 1: Run integration tests**

Run: `./scripts/test-integration.sh`
Expected: All tests pass (or script exits with 0)

- [ ] **Step 2: Run end-to-end tests**

Run: `./scripts/test-e2e.sh`
Expected: All tests pass (or script exits with 0)

### Task 3: Commit changes

**Files:**
- Modify: `RELEASE_CHECKLIST.md`

- [ ] **Step 1: Add file to git**

Run: `git add RELEASE_CHECKLIST.md`

- [ ] **Step 2: Commit with message**

Run: `git commit -m "docs: add release checklist"`

- [ ] **Step 3: Verify commit**

Run: `git log --oneline -1`
Expected: Commit message appears
