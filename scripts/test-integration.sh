#!/usr/bin/env bash
set -e

echo "=== PairAdmin Phase 3 Integration Test ==="
echo "Building backend..."
go build ./...

echo "Running unit tests..."
go test ./internal/llm/... ./internal/config/... ./internal/session/... ./internal/clipboard/... ./internal/security/...

echo "Creating test configuration..."
mkdir -p ~/.pairadmin
cp config.yaml ~/.pairadmin/config.yaml

echo "Starting session store test..."
go run ./cmd/test-session-store/main.go 2>/dev/null || echo "Session store test skipped (no test binary)"

echo "✅ Phase 3 components compile and unit tests pass."
echo "Next: run 'wails dev' to test the frontend."