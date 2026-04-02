#!/bin/bash
#
# Multi-Tier Workflow Orchestrator CLI
# 
# Routes tasks to appropriate LLM models based on tier.
#
# Setup:
#   export OPENROUTER_API_KEY="your-api-key-here"
#
# Usage:
#   ./orchestrate.sh --task 1.3 --tier L0-Planner --spec docs/tasks/1.3-task.md
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ORCHESTRATOR="$SCRIPT_DIR/lib/orchestrator.py"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_error() { echo -e "${RED}✗ $1${NC}"; }

# Check for OpenRouter API key
if [ -z "$OPENROUTER_API_KEY" ]; then
    print_error "OPENROUTER_API_KEY environment variable not set"
    echo ""
    echo "Set it with:"
    echo "  export OPENROUTER_API_KEY=\"your-key-here\""
    echo ""
    echo "Get your key from: https://openrouter.ai/keys"
    exit 1
fi

# Check for Python
if ! command -v python3 &> /dev/null; then
    print_error "Python3 is required but not installed"
    exit 1
fi

# Show help
if [ "$1" == "--help" ] || [ "$1" == "-h" ] || [ -z "$1" ]; then
    print_header "Multi-Tier Workflow Orchestrator"
    echo ""
    echo "Usage: $0 --task <TASK_ID> --tier <TIER> [options]"
    echo ""
    echo "Required:"
    echo "  --task <TASK_ID>     Task identifier (e.g., 1.3)"
    echo "  --tier <TIER>        Tier to use (see below)"
    echo ""
    echo "Options:"
    echo "  --spec <FILE>        Path to task specification file"
    echo "  --context <FILE>     Path to context JSON file"
    echo "  --output <FILE>      Save response to file"
    echo "  --help, -h           Show this help"
    echo ""
    echo "Available Tiers:"
    echo "  L0-Planner    - Qwen3.5 397B (OpenRouter) - Task specs, reviews"
    echo "  L0-Reviewer   - Qwen3.5 397B (OpenRouter) - Compliance checks"
    echo "  L0-Coder      - Qwen3-Coder (LM Studio)   - First implementation"
    echo "  L1-Coder      - Grok 4.1 Fast (OpenRouter) - Re-implementation"
    echo "  L2-Coder      - MiniMax M2.7 (OpenRouter) - Complex implementation"
    echo "  L3-Coder      - Claude Sonnet 4.6 (OpenRouter) - Final escalation"
    echo "  L3-Architect  - Claude Opus 4.6 (OpenRouter) - Checkpoint reviews"
    echo ""
    echo "Examples:"
    echo "  # Create task spec for Task 1.3"
    echo "  $0 --task 1.3 --tier L0-Planner --spec docs/tasks/1.3-task.md"
    echo ""
    echo "  # Implement with local Qwen3-Coder"
    echo "  $0 --task 1.3 --tier L0-Coder --spec docs/tasks/1.3-task.md --output impl.txt"
    echo ""
    echo "  # Review implementation"
    echo "  $0 --task 1.3 --tier L0-Reviewer --spec docs/tasks/1.3-task.md --context review.json"
    echo ""
    exit 0
fi

# Run orchestrator
python3 "$ORCHESTRATOR" "$@"
