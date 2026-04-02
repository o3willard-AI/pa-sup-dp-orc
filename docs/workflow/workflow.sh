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
