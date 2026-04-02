# Config Race Fixes Design

**Date:** 2026-04-01  
**Author:** opencode  
**Status:** Approved

## Overview

Fix critical and important issues identified by code quality review for the config package (`internal/config/manager.go`). The issues involve thread safety, validation improvements, and missing concurrency tests.

## Changes

### 1. Critical: Viper writes under read lock – data race risk
- **Location:** `manager.go:138-153`
- **Issue:** `Save` holds an `RLock` but calls `viper.Set()` and `viper.WriteConfig()`. Viper is not thread-safe; concurrent writes can cause data corruption.
- **Fix:** Change `Save` to use `Lock` (write lock) for the entire operation.

### 2. Important: Missing concurrency tests
- **Issue:** No tests exercise concurrent calls to `Init`, `Get`, or `Save`. The race detector will not catch issues without parallel goroutines.
- **Fix:** Add a test that spawns multiple goroutines calling these functions (with `-race` enabled).

### 3. Important: Error messages could be more informative
- **Issue:** Errors like “openai provider requires api_key” don’t indicate which config key or environment variable to set.
- **Fix:** Include the field name and config key (e.g., `llm.openai.api_key`) and environment variable (e.g., `PAIRADMIN_LLM_OPENAI_API_KEY`).

### 4. Minor: Validate URL format for Ollama
- **Issue:** Only checks non‑empty `BaseURL`, not a valid URL format.
- **Fix:** Add `url.Parse` validation requiring scheme `http` or `https` and non‑empty host.

### 5. Minor: Thread‑safety documentation
- **Issue:** Exported config fields allow unprotected mutation.
- **Fix:** Add comment in `Get` function that the returned Config pointer is shared and must not be modified concurrently by callers.

## Implementation Details

### manager.go changes

#### `Save()` method
- Replace `configMu.RLock()` with `configMu.Lock()`
- Keep same validation and write logic

#### `validate()` method
- Update error messages to include:
  - OpenAI: `"openai provider requires api_key (config: llm.openai.api_key, env: PAIRADMIN_LLM_OPENAI_API_KEY)"`
  - Anthropic: `"anthropic provider requires api_key (config: llm.anthropic.api_key, env: PAIRADMIN_LLM_ANTHROPIC_API_KEY)"`
  - Ollama: `"ollama provider requires base_url (config: llm.ollama.base_url, env: PAIRADMIN_LLM_OLLAMA_BASE_URL)"`
  - Unknown provider: keep existing message

#### URL validation for Ollama
- After checking BaseURL not empty, parse with `url.Parse`
- Require `Scheme` be `"http"` or `"https"`
- Require `Host` not empty
- Return descriptive error if invalid

#### `Get()` documentation
- Add comment: `// Get returns the global configuration. The returned Config pointer is shared and must not be modified concurrently by callers.`

### manager_test.go changes

#### New test `TestConcurrentAccess`
- Initialize config once with valid provider and API keys
- Spawn 10 reader goroutines calling `Get` in loop (e.g., 100 iterations)
- Spawn 2 writer goroutines calling `Save` with alternating provider values
- Use `sync.WaitGroup` for coordination
- Run with `-race` flag (test will be part of normal test suite)
- No concurrent `Init` calls (should be called once in production)

## Testing

All existing tests must continue to pass:
- `TestConfigInitAndSave`
- `TestConfigValidation`
- `TestConfigInvalidYAML`
- `TestConfigPermissionError`
- `TestConfigEmptyConfig`

Run verification:
```bash
go test -race ./internal/config/... -v
go test ./internal/config/...
```

## Constraints

- Do not change the public API signatures (`Init`, `Get`, `Save`).
- Keep existing test behavior (the test `TestConfigInitAndSave` modifies `cfg.LLM.Provider` directly). This must still pass.
- Ensure `go test -race ./internal/config/...` passes.
- Ensure `go test ./internal/config/...` passes.

## Risks

- Changing `Save` to use `Lock` may reduce throughput in high‑concurrency scenarios, but config saving is infrequent operation.
- URL validation may reject valid but unusual URLs (e.g., `http://localhost:11434/v1`). `url.Parse` handles these correctly.
- Error message changes may break external tools parsing error strings. Unlikely as this is internal package.

## Dependencies

- Go standard library: `net/url` for URL parsing
- No external dependencies added

## Commit Strategy

Single commit with message:
```
fix(config): address race conditions and improve validation

- Change Save() to use write lock (fix data race with Viper)
- Add concurrent access test with -race
- Improve validation error messages with config/env keys
- Validate Ollama BaseURL format
- Add thread-safety comment to Get()
```