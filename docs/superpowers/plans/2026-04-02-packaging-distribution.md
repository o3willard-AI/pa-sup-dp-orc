# Packaging and Distribution Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Configure Wails native packaging, create GitHub Actions release workflow, and provide Electron‑builder fallback configuration.

**Architecture:** Update wails.json with packaging metadata, create cross‑platform release.yml workflow, add root package.json with Electron‑builder config matching wails.json settings.

**Tech Stack:** Wails 2, GitHub Actions, electron‑builder (fallback), Go 1.25+, Node.js 20.

---

### Task 1: Update wails.json with packaging configuration

**Files:**
- Modify: `/home/sblanken/code/paa5/wails.json`

- [ ] **Step 1: Read current wails.json**

```bash
cat wails.json
```

Expected: See current config with `$schema`, `name: "wails-pairadmin"`, `outputfilename: "wails-pairadmin"`, minimal `author`.

- [ ] **Step 2: Update wails.json with packaging fields**

Edit `wails.json` to merge new packaging configuration while preserving `$schema`:

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "PairAdmin",
  "outputfilename": "pairadmin",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "PairAdmin Team",
    "email": "team@pairadmin.dev"
  },
  "info": {
    "companyName": "PairAdmin",
    "productName": "PairAdmin",
    "productVersion": "2.0.0",
    "copyright": "Copyright © 2026 PairAdmin",
    "comments": "AI‑assisted terminal administration"
  },
  "deb": {
    "depends": ["libgtk-3-0", "libwebkit2gtk-4.0-37"]
  },
  "nsis": {
    "installDirectory": "$PROGRAMFILES\\PairAdmin"
  },
  "mac": {
    "bundle": "dev.pairadmin.app",
    "category": "public.app-category.developer-tools"
  }
}
```

- [ ] **Step 3: Verify JSON syntax**

```bash
python3 -m json.tool wails.json > /dev/null && echo "JSON valid"
```

Expected: "JSON valid"

- [ ] **Step 4: Commit wails.json changes**

```bash
git add wails.json
git commit -m "build: update wails.json with packaging configuration"
```

---

### Task 2: Create GitHub Actions release workflow

**Files:**
- Create: `/home/sblanken/code/paa5/.github/workflows/release.yml`

- [ ] **Step 1: Create directory if missing**

```bash
mkdir -p .github/workflows
```

- [ ] **Step 2: Write release.yml**

Create `.github/workflows/release.yml`:

```yaml
name: Release
on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
      - name: Install dependencies
        run: ./scripts/install-deps.sh
      - name: Build
        run: wails build -platform ${{ matrix.os == 'ubuntu-latest' && 'linux' || matrix.os == 'windows-latest' && 'windows' || 'darwin' }}
      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: pairadmin-${{ matrix.os }}
          path: build/bin/*
```

- [ ] **Step 3: Validate YAML syntax**

```bash
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/release.yml'))" && echo "YAML valid"
```

Expected: "YAML valid"

- [ ] **Step 4: Commit release workflow**

```bash
git add .github/workflows/release.yml
git commit -m "ci: add GitHub Actions release workflow for tagged releases"
```

---

### Task 3: Create root package.json with Electron‑builder fallback

**Files:**
- Create: `/home/sblanken/code/paa5/package.json`

- [ ] **Step 1: Write package.json with Electron‑builder configuration**

Create `package.json` at project root:

```json
{
  "name": "pairadmin",
  "version": "2.0.0",
  "description": "AI‑assisted terminal administration",
  "scripts": {
    "build": "wails build"
  },
  "devDependencies": {
    "electron": "^latest",
    "electron-builder": "^latest"
  },
  "build": {
    "appId": "dev.pairadmin.app",
    "productName": "PairAdmin",
    "directories": {
      "output": "dist"
    },
    "files": [
      "build/bin/*",
      "frontend/dist/**/*"
    ],
    "mac": {
      "category": "public.app-category.developer-tools"
    },
    "linux": {
      "target": ["deb", "rpm", "AppImage"],
      "category": "Development"
    },
    "win": {
      "target": "nsis"
    }
  }
}
```

- [ ] **Step 2: Validate JSON syntax**

```bash
python3 -m json.tool package.json > /dev/null && echo "JSON valid"
```

Expected: "JSON valid"

- [ ] **Step 3: Commit package.json**

```bash
git add package.json
git commit -m "build: add Electron‑builder fallback packaging configuration"
```

---

### Task 4: Test packaging locally

**Files:**
- Test: Build output in `build/bin/`

- [ ] **Step 1: Run Wails build for Linux platform**

```bash
wails build -platform linux
```

Expected: Build succeeds (may warn about missing system dependencies but still produce binary).

- [ ] **Step 2: Verify binary exists**

```bash
ls -la build/bin/
```

Expected: See `pairadmin` binary (or similar).

- [ ] **Step 3: Commit any generated build artifacts (optional)**

```bash
git add build/bin/ 2>/dev/null || echo "No new build artifacts to commit"
```

---

### Task 5: Final commit of all packaging changes

**Files:**
- All modified/created files

- [ ] **Step 1: Verify all changes are staged**

```bash
git status --porcelain
```

Expected: No unstaged changes (or only expected untracked files).

- [ ] **Step 2: Create final commit summarizing packaging work**

```bash
git commit -m "build: complete packaging and distribution configuration

- Update wails.json with platform‑specific packaging metadata
- Add GitHub Actions release workflow for tagged releases
- Provide Electron‑builder fallback configuration in root package.json
- Test local Linux build"
```

---

**Plan complete and saved to `docs/superpowers/plans/2026-04-02-packaging-distribution.md`. Two execution options:**

**1. Subagent‑Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration.

**2. Inline Execution** - Execute tasks in this session using executing‑plans, batch execution with checkpoints.

**Which approach?**