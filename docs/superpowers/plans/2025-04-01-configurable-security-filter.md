# Configurable Security Filter Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Modify `NewChatHandlers` to use configurable security filter patterns from `config.yaml` instead of hardcoded default patterns.

**Architecture:** The filter will be created from `cfg.Security.FilterPatterns`. If no patterns configured or error compiling regex, fall back to default filter with logging.

**Tech Stack:** Go, wails runtime logging, regexp

---

### Task 1: Modify chat_handlers.go to use configurable patterns

**Files:**
- Modify: `internal/ui/chat_handlers.go:64`

- [ ] **Step 1: Write the failing test**

Add a unit test for the private helper function we'll create. In `internal/ui/chat_handlers_test.go`, add:

```go
func TestCreateFilterFromConfig(t *testing.T) {
	ctx := context.Background()
	
	// Test nil config -> default filter
	filter := createFilterFromConfig(ctx, nil)
	defaultFilter := security.DefaultFilter()
	text := "password = secret"
	if filter.Redact(text) != defaultFilter.Redact(text) {
		t.Errorf("nil config should return default filter")
	}

	// Test empty config -> default filter
	cfg := &config.Config{}
	filter = createFilterFromConfig(ctx, cfg)
	if filter.Redact(text) != defaultFilter.Redact(text) {
		t.Errorf("empty patterns should return default filter")
	}

	// Test custom pattern
	cfg.Security.FilterPatterns = []config.FilterPattern{
		{Pattern: "FOO_\\d+"},
	}
	filter = createFilterFromConfig(ctx, cfg)
	redacted := filter.Redact("secret FOO_123 bar")
	if redacted != "secret [REDACTED] bar" {
		t.Errorf("custom pattern not applied, got %q", redacted)
	}

	// Test invalid regex -> fallback with logging (cannot test logging easily)
	cfg.Security.FilterPatterns = []config.FilterPattern{
		{Pattern: "["}, // invalid regex
	}
	filter = createFilterFromConfig(ctx, cfg)
	if filter.Redact(text) != defaultFilter.Redact(text) {
		t.Errorf("invalid regex should fall back to default filter")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/ui -run TestCreateFilterFromConfig -v`
Expected: FAIL because createFilterFromConfig is undefined.

- [ ] **Step 3: Write minimal implementation**

Add helper function in `internal/ui/chat_handlers.go` after imports (before NewChatHandlers):

```go
func createFilterFromConfig(ctx context.Context, cfg *config.Config) *security.Filter {
	if cfg == nil || len(cfg.Security.FilterPatterns) == 0 {
		return security.DefaultFilter()
	}
	var rawPatterns []string
	for _, fp := range cfg.Security.FilterPatterns {
		rawPatterns = append(rawPatterns, fp.Pattern)
	}
	filter, err := security.NewFilter(rawPatterns)
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("failed to compile custom filter patterns: %v, falling back to default", err))
		return security.DefaultFilter()
	}
	return filter
}
```

Replace line 64 (`filter := security.DefaultFilter()`) with:

```go
filter := createFilterFromConfig(ctx, cfg)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/ui -run TestCreateFilterFromConfig -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/ui/chat_handlers.go internal/ui/chat_handlers_test.go
git commit -m "feat: add configurable security filter patterns"
```

### Task 2: Integration test with config file

**Files:**
- Modify: `internal/ui/chat_handlers_test.go`

Add a test that initializes config with custom patterns and verifies the filter in ChatHandlers uses those patterns (via reflection).

- [ ] **Step 1: Write the failing test**

Add to chat_handlers_test.go:

```go
func TestNewChatHandlers_CustomFilterPatterns(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configContent := `llm:
  provider: openai
  openai:
    api_key: dummy
    model: gpt-4
    base_url: http://localhost
security:
  filter_patterns:
    - name: "test"
      pattern: "SECRET_\\w+"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = config.Init(configPath)
	if err != nil {
		t.Fatal(err)
	}
	defer config.Init(configPath) // reset

	ctx := context.Background()
	store, err := session.NewStore(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory store: %v", err)
	}
	defer store.Close()

	handlers, err := NewChatHandlers(ctx, store)
	if err != nil {
		t.Fatalf("NewChatHandlers failed: %v", err)
	}
	// Use reflection to get filter field
	filterField := reflect.ValueOf(handlers).Elem().FieldByName("filter")
	if !filterField.IsValid() {
		t.Fatal("filter field not found")
	}
	filter := filterField.Interface().(*security.Filter)
	// Verify filter redacts custom pattern
	redacted := filter.Redact("hello SECRET_KEY world")
	if redacted != "hello [REDACTED] world" {
		t.Errorf("custom pattern not applied, got %q", redacted)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/ui -run TestNewChatHandlers_CustomFilterPatterns -v`
Expected: FAIL because filter still uses default patterns.

- [ ] **Step 3: Ensure implementation passes**

The implementation from Task 1 should satisfy this test.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/ui -run TestNewChatHandlers_CustomFilterPatterns -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/ui/chat_handlers_test.go
git commit -m "test: integration test for custom filter patterns"
```

### Task 3: Run lint and typecheck

**Files:**
- Run lint command if available.

- [ ] **Step 1: Find lint command**

Check for `go fmt`, `go vet`, `golangci-lint`. Run:

```bash
go fmt ./...
go vet ./...
golangci-lint run ./...
```

- [ ] **Step 2: Fix any issues**

- [ ] **Step 3: Commit any fixes**

```bash
git add -u
git commit -m "style: fix formatting"
```

### Task 4: Final verification

**Files:**
- Run existing test suite.

- [ ] **Step 1: Run all tests**

```bash
go test ./...
```

- [ ] **Step 2: Ensure no regression**

Check that default filter still works when no config patterns (by running integration test with empty filter_patterns). We'll add a quick test or manually verify.

- [ ] **Step 3: Commit final changes**

```bash
git add -u
git commit -m "fix: use configurable security filter patterns"
```

**Plan complete and saved to `docs/superpowers/plans/2025-04-01-configurable-security-filter.md`. Two execution options:**

**1. Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

**Which approach?**