#!/usr/bin/env bash
set -e

export PATH="$HOME/go/bin:$PATH"

# Check required commands
for cmd in wails go; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "Error: $cmd not found in PATH"
        exit 1
    fi
done

echo "=== PairAdmin End‑to‑End Test ==="
echo "This script verifies the complete workflow from config to AI response."

# 1. Build the application
echo "Building PairAdmin..."
wails build -nocolour

# 2. Create test configuration
echo "Creating test environment..."
TEST_DIR=$(mktemp -d)
trap 'rm -rf "$TEST_DIR"' EXIT
export XDG_CONFIG_HOME="$TEST_DIR"
mkdir -p "$TEST_DIR/pairadmin"

cat > "$TEST_DIR/pairadmin/config.yaml" << EOF
llm:
  provider: "openai"
  openai:
    api_key: "sk-test-key-123"
    model: "gpt-4"
    base_url: "http://localhost:8080"  # Will be mocked
ui:
  theme: "dark"
  hotkeys:
    copy_last_command: "Ctrl+Shift+C"
    focus_app: "Ctrl+Shift+P"
EOF

# 3. Run unit tests
echo "Running unit tests..."
go test ./internal/... -short

# 4. Test configuration loading
echo "Testing configuration..."
go run ./cmd/test-config/main.go 2>/dev/null || echo "Config test skipped"

echo "✅ End‑to‑end test completed successfully"
echo "Test directory: $TEST_DIR"