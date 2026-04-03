#!/usr/bin/env python3
"""Multi-Tier LLM Orchestrator - Core Implementation"""

import os
import sys
import json
import time
import requests
import re
from pathlib import Path
from typing import Optional, Dict, Any, List
from datetime import datetime, timezone

# Configuration
MAX_RETRIES = 3
RETRY_DELAY = 2
CONTEXT_SIMPLIFICATION = [1.0, 0.7, 0.4]

# Model Configuration
MODELS = {
    "L0-Planner": {"provider": "openrouter", "model": "qwen/qwen3.5-397b-a17b", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.3, "tools": ["file_read"]},
    "L0-Reviewer": {"provider": "openrouter", "model": "qwen/qwen3.5-397b-a17b", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.3, "tools": ["file_read"]},
    "L0-Coder": {"provider": "lmstudio", "model": "qwen/qwen3-coder-30b", "base_url": "http://192.168.101.21:1234/v1", "temperature": 0.7, "tools": ["file_read", "file_write"]},
    "L1-Coder": {"provider": "openrouter", "model": "x-ai/grok-4.1-fast", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.7, "tools": ["file_read", "file_write"]},
    "L2-Coder": {"provider": "openrouter", "model": "minimax/minimax-m2.7", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.7, "tools": ["file_read", "file_write"]},
    "L3-Coder": {"provider": "openrouter", "model": "anthropic/claude-sonnet-4.6", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.7, "tools": ["file_read", "file_write"]},
    "L3-Architect": {"provider": "openrouter", "model": "anthropic/claude-opus-4.6", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.3, "tools": ["file_read"]}
}


class FileTools:
    """File read/write operations."""
    
    def __init__(self, project_root: Path):
        self.project_root = project_root
    
    def _fix_escaped_content(self, content: str) -> str:
        """P4: Fix escaped newlines/characters from L0-Coder output.
        
        L0-Coder (Qwen3-Coder 30B) often outputs \\n instead of actual newlines
        in triple-quoted strings. This method detects and fixes that pattern.
        """
        # Detect escaped newline pattern (common L0-Coder issue)
        # If content has literal backslash-n but no actual newlines, fix it
        if '\\n' in content and content.count('\n') < 5:
            # Likely escaped - convert to actual characters
            content = content.replace('\\n', '\n')
            content = content.replace('\\t', '\t')
            content = content.replace('\\"', '"')
            content = content.replace('\\\\', '\\')
        return content
    
    def file_read(self, path: str) -> Dict[str, Any]:
        """Read a file and return contents."""
        try:
            file_path = Path(path)
            if not file_path.is_absolute():
                file_path = self.project_root / file_path
            if not file_path.exists():
                return {"success": False, "error": f"File not found: {file_path}"}
            content = file_path.read_text()
            return {"success": True, "path": str(file_path), "content": content, "lines": len(content.splitlines())}
        except Exception as e:
            return {"success": False, "error": str(e)}
    
    def file_write(self, path: str, content: str) -> Dict[str, Any]:
        """Write content to a file with post-processing for L0-Coder issues."""
        try:
            # P4: Fix escaped newlines from L0-Coder (common issue)
            content = self._fix_escaped_content(content)
            
            file_path = Path(path)
            if not file_path.is_absolute():
                file_path = self.project_root / file_path
            file_path.parent.mkdir(parents=True, exist_ok=True)
            file_path.write_text(content)
            return {"success": True, "path": str(file_path), "bytes": len(content)}
        except Exception as e:
            return {"success": False, "error": str(e)}


class ToolExecutor:
    """Executes tool calls from LLM responses."""
    
    def __init__(self, file_tools: FileTools):
        self.file_tools = file_tools
    
    def _validate_tool_call(self, tool: str, path: str) -> bool:
        """Validate tool call parameters before execution (P2)."""
        if not path:
            return False
        
        # Check path length
        if len(path) > 500:
            return False
        
        # Check for corrupted JSON artifacts
        if path.startswith('}') or path.startswith('{'):
            return False
        if path.startswith(']]') or path.startswith('[['):
            return False
        
        # Check for newlines in path (indicates parsing error)
        if '\n' in path or '\r' in path:
            return False
        
        # Check for valid path characters (allow common path chars)
        if not re.match(r'^[a-zA-Z0-9_./\-\\]+$', path):
            return False
        
        # Check path doesn't start with special chars
        if path.startswith('/') and len(path) > 1 and path[1] in '{}[]()':
            return False
        
        return True
    
    def parse_and_execute_tools(self, response: str) -> Dict[str, Any]:
        """Parse LLM response for tool calls and execute them.
        
        P8 FIX: Removed broken JSON line parser that failed on multi-line content.
        Now relies on regex patterns which handle multi-line content correctly with re.DOTALL.
        """
        results = []
        
        # Pattern 1: Custom file_write - file_write("path", """content""")
        write_pattern = r'file_write\(["\']([^"\']+)["\'],\s*"""(.+?)"""'
        for match in re.finditer(write_pattern, response, re.DOTALL):
            path = match.group(1)
            content = match.group(2)
            result = self.file_tools.file_write(path, content)
            results.append({
                "tool": "file_write",
                "path": path,
                "success": result.get("success", False),
                "error": result.get("error"),
                "bytes": result.get("bytes", 0)
            })
        
        # Pattern 2: Custom file_read - file_read("path")
        read_pattern = r'file_read\(["\']([^"\']+)["\']\)'
        for match in re.finditer(read_pattern, response):
            path = match.group(1)
            result = self.file_tools.file_read(path)
            results.append({
                "tool": "file_read",
                "path": path,
                "success": result.get("success", False),
                "content": result.get("content"),
                "error": result.get("error")
            })
        
        # Pattern 3: Claude native write - {"name": "write_file", "parameters": {"path": "x", "content": "y"}}
        # P9 FIX: Improved pattern to handle multi-line content and escaped characters
        # Matches entire JSON object to avoid partial matches
        claude_write = r'\{\s*"name"\s*:\s*"write_file"\s*,\s*"parameters"\s*:\s*\{\s*"path"\s*:\s*"([^"]+)"\s*,\s*"content"\s*:\s*"((?:[^"\\]|\\.)+)"\s*\}\s*\}'
        for match in re.finditer(claude_write, response, re.DOTALL):
            path = match.group(1)
            content = match.group(2)
            result = self.file_tools.file_write(path, content)
            results.append({
                "tool": "file_write",
                "path": path,
                "success": result.get("success", False),
                "error": result.get("error"),
                "bytes": result.get("bytes", 0),
                "format": "claude_native"
            })
        
        # Pattern 4: Claude native read - {"name": "read_file", "parameters": {"path": "x"}}
        # Fixed: Use [^{}]* to prevent matching across JSON object boundaries
        claude_read = r'\{\s*"name"\s*:\s*"read_file"\s*,\s*"parameters"\s*:\s*{\s*"path"\s*:\s*"([^"]+)"\s*}\s*}'
        for match in re.finditer(claude_read, response):
            path = match.group(1)
            result = self.file_tools.file_read(path)
            results.append({
                "tool": "file_read",
                "path": path,
                "success": result.get("success", False),
                "content": result.get("content"),
                "error": result.get("error"),
                "format": "claude_native"
            })
        
        # Pattern 5: MiniMax XML write - <invoke name="Write">
        minimax_write = r'<invoke\s+name="Write"[^>]*>\s*<parameter\s+name="file_path"[^>]*>([^<]+)[\s\S]*?<parameter\s+name="content"[^>]*>([^<]+)'
        for match in re.finditer(minimax_write, response, re.DOTALL):
            path = match.group(1).strip()
            content = match.group(2).strip()
            result = self.file_tools.file_write(path, content)
            results.append({
                "tool": "file_write",
                "path": path,
                "success": result.get("success", False),
                "error": result.get("error"),
                "bytes": result.get("bytes", 0),
                "format": "minimax_native"
            })
        
        # Pattern 6: MiniMax XML read - <invoke name="Read">
        minimax_read = r'<invoke\s+name="Read"[^>]*>\s*<parameter\s+name="file_path"[^>]*>([^<]+)'
        for match in re.finditer(minimax_read, response, re.DOTALL):
            path = match.group(1).strip()
            result = self.file_tools.file_read(path)
            results.append({
                "tool": "file_read",
                "path": path,
                "success": result.get("success", False),
                "content": result.get("content"),
                "error": result.get("error"),
                "format": "minimax_native"
            })
        
        return {
            "tool_results": results,
            "tools_executed": len(results),
            "all_succeeded": all(r.get("success", False) for r in results) if results else True
        }


class LLMOrchestrator:
    """Orchestrates multi-tier LLM workflow."""
    
    def __init__(self, project_root: str = str(Path(__file__).parent.parent.parent.parent)):
        self.project_root = Path(project_root)
        self.workflow_dir = self.project_root / "docs" / "workflow"
        self.tasks_dir = self.project_root / "docs" / "tasks"
        self.handoffs_dir = self.workflow_dir / "handoffs"
        self.escalations_dir = self.workflow_dir / "escalations"
        self.file_tools = FileTools(self.project_root)
        self.tool_executor = ToolExecutor(self.file_tools)
        
        # Create directories
        for d in [self.handoffs_dir, self.escalations_dir, self.tasks_dir]:
            d.mkdir(parents=True, exist_ok=True)
    
    def get_api_key(self, tier: str) -> Optional[str]:
        """Get API key for specified tier."""
        config = MODELS.get(tier, {})
        env_var = config.get("env_var")
        return os.environ.get(env_var) if env_var else None
    
    def call_llm_with_retry(self, tier: str, system_prompt: str, user_prompt: str, temperature: float = 0.7) -> Dict[str, Any]:
        """Call LLM with retry logic."""
        last_error = None
        
        for attempt in range(MAX_RETRIES):
            try:
                # Simplify context on retries
                multiplier = CONTEXT_SIMPLIFICATION[min(attempt, len(CONTEXT_SIMPLIFICATION) - 1)]
                prompt = self._simplify_context(user_prompt, multiplier) if attempt > 0 else user_prompt
                
                start_time = datetime.now(timezone.utc)
                response = self.call_llm(tier, system_prompt, prompt, temperature)
                end_time = datetime.now(timezone.utc)
                
                return {
                    "success": True,
                    "output": response,
                    "attempt": attempt + 1,
                    "duration_seconds": (end_time - start_time).total_seconds(),
                    "context_simplified": attempt > 0
                }
            except Exception as e:
                last_error = str(e)
                print(f"  Attempt {attempt + 1}/{MAX_RETRIES} failed: {last_error}")
                if attempt < MAX_RETRIES - 1:
                    print(f"  Retrying in {RETRY_DELAY}s with simplified context...")
                    time.sleep(RETRY_DELAY)
        
        return {
            "success": False,
            "error": last_error,
            "attempts": MAX_RETRIES,
            "ready_for_escalation": True
        }
    
    def _simplify_context(self, prompt: str, multiplier: float) -> str:
        """Simplify prompt by truncating."""
        if multiplier >= 1.0:
            return prompt
        lines = prompt.splitlines()
        max_lines = max(int(len(lines) * multiplier), 10)
        return '\n'.join(lines[:max_lines]) + f"\n\n[Context truncated from {len(lines)} to {max_lines} lines]"
    
    def call_llm(self, tier: str, system_prompt: str, user_prompt: str, temperature: float = 0.7) -> str:
        """Call the appropriate LLM for the specified tier."""
        config = MODELS.get(tier)
        if not config:
            raise ValueError(f"Unknown tier: {tier}")
        
        provider = config.get("provider")
        model = config.get("model")
        base_url = config.get("base_url")
        api_key = self.get_api_key(tier)
        
        if not model or not base_url:
            raise ValueError(f"Invalid configuration for tier {tier}")
        
        messages = [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_prompt}
        ]
        
        if provider == "openrouter":
            return self._call_openrouter(base_url, model, api_key, messages, temperature)
        elif provider == "lmstudio":
            return self._call_lmstudio(base_url, model, messages, temperature)
        else:
            raise ValueError(f"Unknown provider: {provider}")
    
    def _call_openrouter(self, base_url: str, model: str, api_key: Optional[str], messages: List[Dict], temperature: float) -> str:
        """Call OpenRouter API."""
        if not api_key:
            raise ValueError("OPENROUTER_API_KEY not set")
        
        headers = {
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json",
            "HTTP-Referer": "https://github.com/pairadmin/orchestrator",
            "X-Title": "Multi-Tier Orchestrator"
        }
        payload = {
            "model": model,
            "messages": messages,
            "temperature": temperature,
            "max_tokens": 8192
        }
        
        response = requests.post(f"{base_url}/chat/completions", headers=headers, json=payload, timeout=300)
        if response.status_code != 200:
            raise Exception(f"OpenRouter API error: {response.status_code} - {response.text}")
        
        return response.json()["choices"][0]["message"]["content"]
    
    def _call_lmstudio(self, base_url: str, model: str, messages: List[Dict], temperature: float) -> str:
        """Call LM Studio API."""
        headers = {"Content-Type": "application/json"}
        payload = {
            "model": model,
            "messages": messages,
            "temperature": temperature,
            "max_tokens": 8192
        }
        
        response = requests.post(f"{base_url}/chat/completions", headers=headers, json=payload, timeout=300)
        if response.status_code != 200:
            raise Exception(f"LM Studio API error: {response.status_code} - {response.text}")
        
        return response.json()["choices"][0]["message"]["content"]
    
    
    def _load_prompt_template(self, tier: str) -> str:
        """Load prompt template for tier."""
        template_map = {
            "L0-Planner": "01-planner.md",
            "L0-Coder": "02-l0-coder.md",
            "L0-Reviewer": "03-reviewer.md",
            "L1-Coder": "04-l1-coder.md",
            "L2-Coder": "05-l2-coder.md",
            "L3-Coder": "05-l2-coder.md",
            "L3-Architect": "06-l3-architect.md"
        }
        template_file = template_map.get(tier)
        if not template_file:
            raise ValueError(f"No template mapped for tier: {tier}")
        template_path = self.workflow_dir / "templates" / template_file
        if not template_path.exists():
            raise FileNotFoundError(f"Template not found: {template_path}")
        return template_path.read_text()
    
    def _build_system_prompt(self, tier: str, template: str) -> str:
        """Extract system prompt from template."""
        lines = template.split('\n')
        system_lines = []
        for line in lines:
            if line.startswith('# ROLE:'):
                system_lines.append(line)
            elif line.startswith('##'):
                break
            else:
                system_lines.append(line)
        return '\n'.join(system_lines)
    
    def _build_user_prompt(self, task_id: str, tier: str, context: Dict, template: str) -> str:
        """Build user prompt from template and context."""
        user_prompt = template
        if "task_spec" in context:
            user_prompt = user_prompt.replace("{INSERT: Full task specification from Mid-Tier Planner}", context.get("task_spec", ""))
        if "existing_context" in context:
            user_prompt = user_prompt.replace("{INSERT: Relevant existing files, interfaces, patterns}", context.get("existing_context", "First task - no existing code"))
        if "implementation" in context:
            user_prompt = user_prompt.replace("{INSERT: Coder's implementation + self-assessment}", context.get("implementation", ""))
        user_prompt += f"\n\n## Current Task: {task_id}\n\n"
        for key, value in context.items():
            if key not in ["task_spec", "implementation", "existing_context"]:
                user_prompt += f"**{key}:** {value}\n"
        return user_prompt
    
    def _log_handoff(self, task_id: str, tier: str, context: Dict, response: str, timestamp: datetime, duration: float, attempt: int, tool_result: Dict) -> str:
        """Log successful handoff."""
        safe_task_id = task_id.replace(".", "_")
        timestamp_str = timestamp.isoformat().replace(':', '-')
        log_entry = {
            "timestamp": timestamp.isoformat() + "Z",
            "task_id": task_id,
            "tier": tier,
            "attempt": attempt,
            "duration_seconds": duration,
            "success": True,
            "input_context": {k: v for k, v in context.items() if len(str(v)) < 1000},
            "output_preview": response[:500] + "..." if len(response) > 500 else response,
            "tool_results": tool_result
        }
        log_file = self.handoffs_dir / f"{safe_task_id}-{tier.lower()}-{timestamp_str}.json"
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        return str(log_file)
    
    def _log_failure(self, task_id: str, tier: str, context: Dict, error: str, timestamp: datetime, attempts: int) -> str:
        """Log failure for escalation."""
        safe_task_id = task_id.replace(".", "_")
        timestamp_str = timestamp.isoformat().replace(':', '-')
        log_entry = {
            "timestamp": timestamp.isoformat() + "Z",
            "task_id": task_id,
            "tier": tier,
            "attempts": attempts,
            "success": False,
            "error": error,
            "ready_for_escalation": True,
            "context_summary": {k: str(v)[:200] for k, v in context.items()}
        }
        log_file = self.escalations_dir / f"{safe_task_id}-{tier.lower()}-failed-{timestamp_str}.json"
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        return str(log_file)
    
    def execute_task(self, task_id: str, tier: str, context: Dict[str, Any]) -> Dict[str, Any]:
        """Execute a task using the specified tier."""
        template = self._load_prompt_template(tier)
        tier_config = MODELS.get(tier, {})
        temperature = tier_config.get("temperature", 0.7)
        
        system_prompt = self._build_system_prompt(tier, template)
        user_prompt = self._build_user_prompt(task_id, tier, context, template)
        
        result = self.call_llm_with_retry(tier, system_prompt, user_prompt, temperature)
        timestamp = datetime.now(timezone.utc)
        
        if result["success"]:
            # Execute tool calls from LLM output
            tool_result = self.tool_executor.parse_and_execute_tools(result["output"])
            
            # Check if file_write tools succeeded
            file_writes = [r for r in tool_result["tool_results"] if r["tool"] == "file_write"]
            if file_writes:
                failed_writes = [w for w in file_writes if not w["success"]]
                if failed_writes:
                    result["success"] = False
                    result["error"] = f"file_write failed: {failed_writes[0].get('error', 'Unknown error')}"
            
            result["tool_results"] = tool_result
            handoff = self._log_handoff(task_id, tier, context, result["output"], timestamp, result["duration_seconds"], result["attempt"], tool_result)
        else:
            handoff = self._log_failure(task_id, tier, context, result.get("error", "Unknown error"), timestamp, result["attempts"])
        
        return {
            "task_id": task_id,
            "tier": tier,
            "success": result["success"],
            "output": result.get("output", ""),
            "attempts": result.get("attempt", result.get("attempts", 0)),
            "duration_seconds": result.get("duration_seconds", 0),
            "ready_for_escalation": result.get("ready_for_escalation", False),
            "tool_results": result.get("tool_results", {"tool_results": [], "tools_executed": 0, "all_succeeded": True}),
            "handoff_log": handoff
        }
    
    def _load_prompt_template(self, tier: str) -> str:
        """Load prompt template for tier."""
        template_map = {
            "L0-Planner": "01-planner.md",
            "L0-Coder": "02-l0-coder.md",
            "L0-Reviewer": "03-reviewer.md",
            "L1-Coder": "04-l1-coder.md",
            "L2-Coder": "05-l2-coder.md",
            "L3-Coder": "05-l2-coder.md",
            "L3-Architect": "06-l3-architect.md"
        }
        template_file = template_map.get(tier)
        if not template_file:
            raise ValueError(f"No template mapped for tier: {tier}")
        template_path = self.workflow_dir / "templates" / template_file
        if not template_path.exists():
            raise FileNotFoundError(f"Template not found: {template_path}")
        return template_path.read_text()
    
    def _build_system_prompt(self, tier: str, template: str) -> str:
        """Extract system prompt from template."""
        lines = template.split('\n')
        system_lines = []
        for line in lines:
            if line.startswith('# ROLE:'):
                system_lines.append(line)
            elif line.startswith('##'):
                break
            else:
                system_lines.append(line)
        return '\n'.join(system_lines)
    
    def _build_user_prompt(self, task_id: str, tier: str, context: Dict, template: str) -> str:
        """Build user prompt from template and context."""
        user_prompt = template
        if "task_spec" in context:
            user_prompt = user_prompt.replace("{INSERT: Full task specification from Mid-Tier Planner}", context.get("task_spec", ""))
        if "existing_context" in context:
            user_prompt = user_prompt.replace("{INSERT: Relevant existing files, interfaces, patterns}", context.get("existing_context", "First task - no existing code"))
        if "implementation" in context:
            user_prompt = user_prompt.replace("{INSERT: Coder's implementation + self-assessment}", context.get("implementation", ""))
        user_prompt += f"\n\n## Current Task: {task_id}\n\n"
        for key, value in context.items():
            if key not in ["task_spec", "implementation", "existing_context"]:
                user_prompt += f"**{key}:** {value}\n"
        return user_prompt
    
    def _log_handoff(self, task_id: str, tier: str, context: Dict, response: str, timestamp: datetime, duration: float, attempt: int, tool_result: Dict) -> str:
        """Log successful handoff."""
        safe_task_id = task_id.replace(".", "_")
        timestamp_str = timestamp.isoformat().replace(':', '-')
        log_entry = {
            "timestamp": timestamp.isoformat() + "Z",
            "task_id": task_id,
            "tier": tier,
            "attempt": attempt,
            "duration_seconds": duration,
            "success": True,
            "input_context": {k: v for k, v in context.items() if len(str(v)) < 1000},
            "output_preview": response[:500] + "..." if len(response) > 500 else response,
            "tool_results": tool_result
        }
        log_file = self.handoffs_dir / f"{safe_task_id}-{tier.lower()}-{timestamp_str}.json"
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        return str(log_file)
    
    def _log_failure(self, task_id: str, tier: str, context: Dict, error: str, timestamp: datetime, attempts: int) -> str:
        """Log failure for escalation."""
        safe_task_id = task_id.replace(".", "_")
        timestamp_str = timestamp.isoformat().replace(':', '-')
        log_entry = {
            "timestamp": timestamp.isoformat() + "Z",
            "task_id": task_id,
            "tier": tier,
            "attempts": attempts,
            "success": False,
            "error": error,
            "ready_for_escalation": True,
            "context_summary": {k: str(v)[:200] for k, v in context.items()}
        }
        log_file = self.escalations_dir / f"{safe_task_id}-{tier.lower()}-failed-{timestamp_str}.json"
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        return str(log_file)


def main():
    """CLI entry point."""
    import argparse
    parser = argparse.ArgumentParser(description="Multi-Tier LLM Orchestrator")
    parser.add_argument("--task", required=True, help="Task ID (e.g., 1.3)")
    parser.add_argument("--tier", required=True, choices=list(MODELS.keys()), help="Tier to use")
    parser.add_argument("--spec", help="Path to task spec file")
    parser.add_argument("--context", help="Path to context JSON file")
    parser.add_argument("--output", help="Save response to file")
    
    args = parser.parse_args()
    
    # Validate API key for OpenRouter tiers
    if MODELS[args.tier].get("provider") == "openrouter" and not os.environ.get("OPENROUTER_API_KEY"):
        print("ERROR: OPENROUTER_API_KEY environment variable not set")
        print("Set with: export OPENROUTER_API_KEY=\"your-key-here\"")
        sys.exit(1)
    
    orchestrator = LLMOrchestrator()
    
    # Build context
    context = {}
    if args.spec:
        context["task_spec"] = Path(args.spec).read_text()
    if args.context:
        with open(args.context) as f:
            context.update(json.load(f))
    
    print(f"Executing Task {args.task} with Tier {args.tier}...")
    print(f"Model: {MODELS[args.tier]['model']}")
    print(f"Provider: {MODELS[args.tier]['provider']}")
    print(f"Max retries: {MAX_RETRIES} (3 attempts per tier before escalation)")
    print(f"Tools: {MODELS[args.tier].get('tools', [])}")
    print()
    
    try:
        result = orchestrator.execute_task(args.task, args.tier, context)
        
        if result["success"]:
            print(f"\n✓ Completed in {result['duration_seconds']:.1f}s (attempt {result['attempts']})")
            if result.get("tool_results"):
                tr = result["tool_results"]
                print(f"  Tools executed: {tr.get('tools_executed', 0)}")
                print(f"  All succeeded: {tr.get('all_succeeded', True)}")
                for tool in tr.get("tool_results", []):
                    status = "✓" if tool.get("success") else "✗"
                    bytes_info = f" - {tool.get('bytes', 0)} bytes" if tool['tool'] == 'file_write' and tool.get('success') else ""
                    print(f"  {status} {tool['tool']}('{tool.get('path', '')}'){bytes_info}")
            print(f"\nOutput preview:")
            preview = result['output'][:500] + "..." if len(result['output']) > 500 else result['output']
            print(preview)
        else:
            print(f"\n✗ Failed after {result['attempts']} attempts: {result.get('error', 'Unknown error')}")
            print(f"\nReady for escalation to next tier.")
        
        print(f"\nHandoff logged to: {result['handoff_log']}")
        
        if result["success"] and args.output:
            Path(args.output).write_text(result['output'])
            print(f"Full output saved to: {args.output}")
        
    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()


def main():
    """CLI entry point."""
    import argparse
    parser = argparse.ArgumentParser(description="Multi-Tier LLM Orchestrator")
    parser.add_argument("--task", required=True, help="Task ID (e.g., 1.3)")
    parser.add_argument("--tier", required=True, choices=list(MODELS.keys()), help="Tier to use")
    parser.add_argument("--spec", help="Path to task spec file")
    parser.add_argument("--context", help="Path to context JSON file")
    parser.add_argument("--output", help="Save response to file")
    
    args = parser.parse_args()
    
    # Validate API key for OpenRouter tiers
    if MODELS[args.tier].get("provider") == "openrouter" and not os.environ.get("OPENROUTER_API_KEY"):
        print("ERROR: OPENROUTER_API_KEY environment variable not set")
        print("Set with: export OPENROUTER_API_KEY=\"your-key-here\"")
        sys.exit(1)
    
    orchestrator = LLMOrchestrator()
    
    # Build context
    context = {}
    if args.spec:
        context["task_spec"] = Path(args.spec).read_text()
    if args.context:
        with open(args.context) as f:
            context.update(json.load(f))
    
    print(f"Executing Task {args.task} with Tier {args.tier}...")
    print(f"Model: {MODELS[args.tier]['model']}")
    print(f"Provider: {MODELS[args.tier]['provider']}")
    print(f"Max retries: {MAX_RETRIES} (3 attempts per tier before escalation)")
    print(f"Tools: {MODELS[args.tier].get('tools', [])}")
    print()
    
    try:
        result = orchestrator.execute_task(args.task, args.tier, context)
        
        if result["success"]:
            print(f"\n✓ Completed in {result['duration_seconds']:.1f}s (attempt {result['attempts']})")
            if result.get("tool_results"):
                tr = result["tool_results"]
                print(f"  Tools executed: {tr.get('tools_executed', 0)}")
                print(f"  All succeeded: {tr.get('all_succeeded', True)}")
                for tool in tr.get("tool_results", []):
                    status = "✓" if tool.get("success") else "✗"
                    bytes_info = f" - {tool.get('bytes', 0)} bytes" if tool['tool'] == 'file_write' and tool.get('success') else ""
                    print(f"  {status} {tool['tool']}('{tool.get('path', '')}'){bytes_info}")
            print(f"\nOutput preview:")
            preview = result['output'][:500] + "..." if len(result['output']) > 500 else result['output']
            print(preview)
        else:
            print(f"\n✗ Failed after {result['attempts']} attempts: {result.get('error', 'Unknown error')}")
            print(f"\nReady for escalation to next tier.")
        
        print(f"\nHandoff logged to: {result['handoff_log']}")
        
        if result["success"] and args.output:
            Path(args.output).write_text(result['output'])
            print(f"Full output saved to: {args.output}")
        
    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
