#!/usr/bin/env python3
"""
Multi-Tier LLM Orchestrator for PairAdmin Workflow

Routes tasks to appropriate models based on tier:
- L0 Planner/Reviewer: Qwen3.5 397B (OpenRouter)
- L0 Coder: Qwen3-Coder (LM Studio local)
- L1 Coder: Grok 4.1 Fast (OpenRouter)
- L2 Coder: MiniMax M2.7 (OpenRouter)
- L3 Coder: Claude Sonnet 4.6 (OpenRouter)
- L3 Architect: Claude Opus 4.6 (OpenRouter)

Usage:
    export OPENROUTER_API_KEY="your-key-here"
    python3 orchestrator.py --task 1.3 --tier L0-Planner
"""

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
SCRIPT_DIR = Path(__file__).parent
WORKFLOW_DIR = SCRIPT_DIR.parent
PROJECT_ROOT = WORKFLOW_DIR.parent.parent

# Retry Configuration
MAX_RETRIES = 3
RETRY_DELAY = 2
CONTEXT_SIMPLIFICATION = [1.0, 0.7, 0.4]

# Model Configuration
MODELS = {
    "L0-Planner": {"provider": "openrouter", "model": "qwen/qwen3.5-397b-a17b", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.3, "tools": ["file_read"]},
    "L0-Reviewer": {"provider": "openrouter", "model": "qwen/qwen3.5-397b-a17b", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.3, "tools": ["file_read"]},
    "L0-Coder": {"provider": "lmstudio", "model": "qwen/qwen3-coder-30b", "base_url": "http://192.168.101.21:1234/v1", "temperature": 0.7, "tools": ["file_read", "file_write"]},
    "L1-Coder": {"provider": "openrouter", "model": "xai/grok-4.1-fast", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.7, "tools": ["file_read", "file_write"]},
    "L2-Coder": {"provider": "openrouter", "model": "minimax/minimax-m2.7", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.7, "tools": ["file_read", "file_write"]},
    "L3-Coder": {"provider": "openrouter", "model": "anthropic/claude-sonnet-4.6", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.7, "tools": ["file_read", "file_write"]},
    "L3-Architect": {"provider": "openrouter", "model": "anthropic/claude-opus-4.6", "env_var": "OPENROUTER_API_KEY", "base_url": "https://openrouter.ai/api/v1", "temperature": 0.3, "tools": ["file_read"]}
}


class FileTools:
    """File read/write tools for subagents."""
    
    def __init__(self, project_root: Path):
        self.project_root = project_root
    
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
        """Write content to a file."""
        try:
            file_path = Path(path)
            if not file_path.is_absolute():
                file_path = self.project_root / file_path
            file_path.parent.mkdir(parents=True, exist_ok=True)
            file_path.write_text(content)
            return {"success": True, "path": str(file_path), "bytes": len(content)}
        except Exception as e:
            return {"success": False, "error": str(e)}



class ToolExecutor:
    """Executes tool calls from LLM responses using multiple format parsers."""
    
    def __init__(self, file_tools: FileTools):
        self.file_tools = file_tools
    
    def parse_and_execute_tools(self, response: str) -> Dict[str, Any]:
        """Parse LLM response for tool calls and execute them."""
        results = []
        
        # Pattern 1: Custom file_write syntax - file_write("path", """content""")
        write_pattern = r'file_write\(["\']([^"\']+)["\'],\s*"""(.+?)"""|file_write\(["\']([^"\']+)["\'],\s*["\'](.+?)["\']\)'
        for match in re.finditer(write_pattern, response, re.DOTALL):
            path = match.group(1) or match.group(3)
            content = match.group(2) or match.group(4)
            if path and content:
                result = self.file_tools.file_write(path, content)
                results.append({"tool": "file_write", "path": path, "success": result.get("success", False), "error": result.get("error"), "bytes": result.get("bytes", 0)})
        
        # Pattern 2: Custom file_read syntax - file_read("path")
        read_pattern = r'file_read\(["\']([^"\']+)["\']\)'
        for match in re.finditer(read_pattern, response):
            path = match.group(1)
            result = self.file_tools.file_read(path)
            results.append({"tool": "file_read", "path": path, "success": result.get("success", False), "content": result.get("content") if result.get("success") else None, "error": result.get("error")})
        
        # Pattern 3: Claude native read_file - {"name": "read_file", "parameters": {"path": "..."}}
        claude_read = r'\{[^}]*"name"[^}]*"read_file"[^}]*"parameters"[^}]*\{[^}]*"path"[^}]*"([^"]+)"[^}]*\}[^}]*\}'
        for match in re.finditer(claude_read, response):
            path_match = re.search(r'"path"[^:]*:\s*"([^"]+)"', match.group(0))
            if path_match:
                result = self.file_tools.file_read(path_match.group(1))
                results.append({"tool": "file_read", "path": path_match.group(1), "success": result.get("success", False), "content": result.get("content"), "error": result.get("error"), "format": "claude_native"})
        
        # Pattern 4: Claude native write_file
        claude_write = r'\{[^}]*"name"[^}]*"write_file"[^}]*"parameters"[^}]*\{[^}]*"path"[^}]*"([^"]+)"[^}]*"content"[^}]*"([^"]*)"[^}]*\}[^}]*\}'
        for match in re.finditer(claude_write, response, re.DOTALL):
            path_match = re.search(r'"path"[^:]*:\s*"([^"]+)"', match.group(0))
            content_match = re.search(r'"content"[^:]*:\s*"([^"]*)"', match.group(0))
            if path_match and content_match:
                result = self.file_tools.file_write(path_match.group(1), content_match.group(1))
                results.append({"tool": "file_write", "path": path_match.group(1), "success": result.get("success", False), "error": result.get("error"), "bytes": result.get("bytes", 0), "format": "claude_native"})
        
        # Pattern 5: MiniMax XML read
        minimax_read = r'<invoke\s+name="Read"[^>]*>\s*<parameter\s+name="file_path"[^>]*>([^<]+)'
        for match in re.finditer(minimax_read, response, re.DOTALL):
            if match.group(1):
                path = match.group(1).strip()
                result = self.file_tools.file_read(path)
                results.append({"tool": "file_read", "path": path, "success": result.get("success", False), "content": result.get("content"), "error": result.get("error"), "format": "minimax_native"})
        
        # Pattern 6: MiniMax XML write - <invoke name="Write">
        minimax_write = r'<invoke\s+name="Write"[^>]*>\s*<parameter\s+name="file_path"[^>]*>([^<]+)[\s\S]*?<parameter\s+name="content"[^>]*>([^<]+)'
        for match in re.finditer(minimax_write, response, re.DOTALL):
            if match.group(1) and match.group(2):
                path = match.group(1).strip()
                content = match.group(2).strip()
                result = self.file_tools.file_write(path, content)
                results.append({"tool": "file_write", "path": path, "success": result.get("success", False), "error": result.get("error"), "bytes": result.get("bytes", 0), "format": "minimax_native"})
        
    
    def _simplify_context(self, prompt: str, multiplier: float) -> str:
        if multiplier >= 1.0:
            return prompt
        lines = prompt.splitlines()
        max_lines = max(int(len(lines) * multiplier), 10)
        return '\n'.join(lines[:max_lines]) + f"\n\n[Context truncated from {len(lines)} to {max_lines} lines]"
    
    def call_llm(self, tier: str, system_prompt: str, user_prompt: str, temperature: float = 0.7) -> str:
        config = MODELS.get(tier)
        if not config:
            raise ValueError(f"Unknown tier: {tier}")
        provider = config.get("provider")
        model = config.get("model")
        base_url = config.get("base_url")
        api_key = self.get_api_key(tier)
        if not model or not base_url:
            raise ValueError(f"Invalid configuration for tier {tier}")
        messages = [{"role": "system", "content": system_prompt}, {"role": "user", "content": user_prompt}]
        if provider == "openrouter":
            return self._call_openrouter(base_url, model, api_key, messages, temperature)
        elif provider == "lmstudio":
            return self._call_lmstudio(base_url, model, messages, temperature)
        else:
            raise ValueError(f"Unknown provider: {provider}")
    
    def _call_openrouter(self, base_url: str, model: str, api_key: Optional[str], messages: List[Dict], temperature: float) -> str:
        if not api_key:
            raise ValueError("OPENROUTER_API_KEY not set")
        headers = {"Authorization": f"Bearer {api_key}", "Content-Type": "application/json", "HTTP-Referer": "https://github.com/pairadmin/pairadmin", "X-Title": "PairAdmin Multi-Tier Workflow"}
        payload = {"model": model, "messages": messages, "temperature": temperature, "max_tokens": 8192}
        response = requests.post(f"{base_url}/chat/completions", headers=headers, json=payload, timeout=300)
        if response.status_code != 200:
            raise Exception(f"OpenRouter API error: {response.status_code} - {response.text}")
        return response.json()["choices"][0]["message"]["content"]
    
    def _call_lmstudio(self, base_url: str, model: str, messages: List[Dict], temperature: float) -> str:
        headers = {"Content-Type": "application/json"}
        payload = {"model": model, "messages": messages, "temperature": temperature, "max_tokens": 8192}
        response = requests.post(f"{base_url}/chat/completions", headers=headers, json=payload, timeout=300)
        if response.status_code != 200:
            raise Exception(f"LM Studio API error: {response.status_code} - {response.text}")
        return response.json()["choices"][0]["message"]["content"]
    
    def execute_task(self, task_id: str, tier: str, context: Dict[str, Any]) -> Dict[str, Any]:
        template = self._load_prompt_template(tier)
        tier_config = MODELS.get(tier, {})
        temperature = tier_config.get("temperature", 0.7)
        system_prompt = self._build_system_prompt(tier, template)
        user_prompt = self._build_user_prompt(task_id, tier, context, template)
        result = self.call_llm_with_retry(tier, system_prompt, user_prompt, temperature)
        timestamp = datetime.now(timezone.utc)
        if result["success"]:
            tool_result = self.tool_executor.parse_and_execute_tools(result["output"])
            file_writes = [r for r in tool_result["tool_results"] if r["tool"] == "file_write"]
            if file_writes:
                failed_writes = [w for w in file_writes if not w["success"]]
                if failed_writes:
                    result["success"] = False
                    result["error"] = f"file_write failed: {failed_writes[0]['error']}"
            result["tool_results"] = tool_result
            handoff = self._log_handoff(task_id, tier, context, result["output"], timestamp, result["duration_seconds"], result["attempt"], tool_result)
        else:
            handoff = self._log_failure(task_id, tier, context, result["error"], timestamp, result["attempts"])
        return {"task_id": task_id, "tier": tier, "success": result["success"], "output": result.get("output", ""), "attempts": result.get("attempt", result.get("attempts", 0)), "duration_seconds": result.get("duration_seconds", 0), "ready_for_escalation": result.get("ready_for_escalation", False), "tool_results": result.get("tool_results", {"tool_results": [], "tools_executed": 0, "all_succeeded": True}), "handoff_log": handoff}
    
    def _load_prompt_template(self, tier: str) -> str:
        template_map = {"L0-Planner": "01-planner.md", "L0-Coder": "02-l0-coder.md", "L0-Reviewer": "03-reviewer.md", "L1-Coder": "04-l1-coder.md", "L2-Coder": "05-l2-coder.md", "L3-Coder": "05-l2-coder.md", "L3-Architect": "06-l3-architect.md"}
        template_file = template_map.get(tier)
        if not template_file:
            raise ValueError(f"No template mapped for tier: {tier}")
        template_path = self.workflow_dir / "templates" / template_file
        if not template_path.exists():
            raise FileNotFoundError(f"Template not found: {template_path}")
        return template_path.read_text()
    
    def _build_system_prompt(self, tier: str, template: str) -> str:
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
        safe_task_id = task_id.replace(".", "_")
        timestamp_str = timestamp.isoformat().replace(':', '-')
        log_entry = {"timestamp": timestamp.isoformat() + "Z", "task_id": task_id, "tier": tier, "attempt": attempt, "duration_seconds": duration, "success": True, "input_context": {k: v for k, v in context.items() if len(str(v)) < 1000}, "output_preview": response[:500] + "..." if len(response) > 500 else response, "tool_results": tool_result}
        log_file = self.handoffs_dir / f"{safe_task_id}-{tier.lower()}-{timestamp_str}.json"
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        return str(log_file)
    
    def _log_failure(self, task_id: str, tier: str, context: Dict, error: str, timestamp: datetime, attempts: int) -> str:
        safe_task_id = task_id.replace(".", "_")
        timestamp_str = timestamp.isoformat().replace(':', '-')
        log_entry = {"timestamp": timestamp.isoformat() + "Z", "task_id": task_id, "tier": tier, "attempts": attempts, "success": False, "error": error, "ready_for_escalation": True, "context_summary": {k: str(v)[:200] for k, v in context.items()}}
        log_file = self.escalations_dir / f"{safe_task_id}-{tier.lower()}-failed-{timestamp_str}.json"
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        return str(log_file)


def main():
    import argparse
    parser = argparse.ArgumentParser(description="Multi-Tier LLM Orchestrator")
    parser.add_argument("--task", required=True, help="Task ID (e.g., 1.3)")
    parser.add_argument("--tier", required=True, choices=list(MODELS.keys()), help="Tier to use")
    parser.add_argument("--spec", help="Path to task spec file")
    parser.add_argument("--context", help="Path to context JSON file")
    parser.add_argument("--output", help="Output file for response")
    args = parser.parse_args()
    orchestrator = LLMOrchestrator()
    context = {}
    if args.spec:
        context["task_spec"] = Path(args.spec).read_text()
    if args.context:
        with open(args.context) as f:
            context.update(json.load(f))
    print(f"Executing Task {args.task} with Tier {args.tier}...")
    print(f"Model: {MODELS[args.tier]['model']}")
    print(f"Provider: {MODELS[args.tier]['provider']}")
    print(f"Max retries: {MAX_RETRIES}")
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
                    print(f"  {status} {tool['tool']}('{tool.get('path', '')}')" + (f" - {tool.get('bytes', 0)} bytes" if tool['tool'] == 'file_write' and tool.get('success') else ""))
            print(f"\nOutput preview:")
            print(result['output'][:500] + "..." if len(result['output']) > 500 else result['output'])
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
        # Pattern 5: MiniMax XML format
