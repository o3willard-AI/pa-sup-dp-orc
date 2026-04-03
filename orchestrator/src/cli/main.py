#!/usr/bin/env python3
"""Enhanced CLI for Multi-Tier Orchestrator."""

import sys, os, json
from pathlib import Path
from datetime import datetime

PROJECT_ROOT = Path(__file__).parent.parent.parent
sys.path.insert(0, str(PROJECT_ROOT))

from src.core.orchestrator import LLMOrchestrator, MODELS
from src.validators.startup import StartupValidator
from src.core.metrics import MetricsCollector
from src.core.cost import CostTracker, Budget, TokenCount

try:
    from rich.console import Console
    from rich.panel import Panel
    from rich.table import Table
    from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn
    from rich import box
    RICH_AVAILABLE = True
except ImportError:
    RICH_AVAILABLE = False
    Console = None

class OrchestratorCLI:
    def __init__(self, use_rich=True):
        self.console = Console() if (use_rich and RICH_AVAILABLE) else None
        self.orchestrator = LLMOrchestrator(str(PROJECT_ROOT.parent))
        self.metrics = MetricsCollector()
        self.use_rich = use_rich and RICH_AVAILABLE
    
    def print_header(self):
        if self.use_rich:
            self.console.print(Panel.fit("[bold blue]Multi-Tier LLM Orchestrator[/bold blue]", border_style="blue"))
        else:
            print("=" * 60)
            print("Multi-Tier LLM Orchestrator")
            print("=" * 60)
    
    def print_status(self, msg, success=True):
        if self.use_rich:
            icon = "✓" if success else "✗"
            color = "green" if success else "red"
            self.console.print(f"[{color}]{icon} {msg}[/{color}]")
        else:
            print(f"[{'PASS' if success else 'FAIL'}] {msg}")
    
    def run_validation(self):
        self.console.print("\n[bold]Running Startup Validation...[/bold]\n")
        validator = StartupValidator(str(PROJECT_ROOT / "templates"))
        ok, results = validator.validate_all()
        
        self.console.print("[bold cyan]API Keys:[/bold cyan]")
        for r in results["api_keys"]["valid"]:
            self.console.print(f"  [green]✓[/green] {r['provider']}")
        for r in results["api_keys"]["invalid"]:
            self.console.print(f"  [red]✗[/red] {r['provider']}: {r['message']}")
        
        self.console.print("\n[bold cyan]Models:[/bold cyan]")
        for r in results["models"]["valid"]:
            self.console.print(f"  [green]✓[/green] {r['tier']}: {r['model']}")
        
        self.console.print("\n[bold cyan]Templates:[/bold cyan]")
        for r in results["templates"]["valid"]:
            self.console.print(f"  [green]✓[/green] {r['file']}")
        
        return ok
    
    def execute_task(self, task_id, tier, context, output_file=None):
        self.console.print(f"\n[bold]Executing Task {task_id} with Tier {tier}[/bold]")
        self.console.print(f"Model: [cyan]{MODELS[tier]['model']}[/cyan]")
        
        with Progress(SpinnerColumn(), TextColumn("[progress.description]{task.description}"), BarColumn(), console=self.console) as progress:
            task = progress.add_task(f"Running {tier}...", total=None)
            try:
                result = self.orchestrator.execute_task(task_id, tier, context)
                progress.update(task, completed=True)
                
                if result["success"]:
                    self.print_status(f"Completed in {result['duration_seconds']:.1f}s (attempt {result['attempts']})", True)
                    if result.get("tool_results"):
                        tr = result["tool_results"]
                        self.console.print(f"\n[bold]Tools:[/bold]")
                        for tool in tr.get("tool_results", []):
                            icon = "[green]✓[/green]" if tool.get("success") else "[red]✗[/red]"
                            self.console.print(f"  {icon} {tool['tool']}('{tool.get('path', '')}')")
                    self.metrics.record_task(task_id, tier, True, result["duration_seconds"], result["attempts"], tr.get("tools_executed", 0))
                else:
                    self.print_status(f"Failed after {result['attempts']} attempts: {result.get('error', 'Unknown')}", False)
                    self.metrics.record_task(task_id, tier, False, result.get("duration_seconds", 0), result["attempts"])
            except Exception as e:
                progress.update(task, completed=True)
                self.print_status(f"Error: {str(e)}", False)
                raise
        
        if result.get("success") and output_file:
            Path(output_file).write_text(result.get("output", ""))
            self.console.print(f"\n[dim]Output saved to: {output_file}[/dim]")
        
        return result
    
    def print_metrics(self):
        self.console.print("\n[bold]Metrics Summary[/bold]")
        summary = self.metrics.get_summary()
        table = Table(box=box.SIMPLE)
        table.add_column("Metric", style="cyan")
        table.add_column("Value", style="white")
        table.add_row("Total Tasks", str(summary["total_tasks"]))
        table.add_row("Success Rate", f"{summary['overall_success_rate']:.1%}")
        table.add_row("Est. Cost", f"${summary['total_cost_usd']:.4f}")
        self.console.print(table)

def main():
    import argparse
    p = argparse.ArgumentParser(description="Multi-Tier LLM Orchestrator CLI")
    p.add_argument("--task", required=True)
    p.add_argument("--tier", required=True, choices=list(MODELS.keys()))
    p.add_argument("--spec")
    p.add_argument("--context")
    p.add_argument("--output")
    p.add_argument("--validate", action="store_true")
    p.add_argument("--no-rich", action="store_true")
    args = p.parse_args()
    
    cli = OrchestratorCLI(use_rich=not args.no_rich)
    cli.print_header()
    
    if args.validate:
        success = cli.run_validation()
        sys.exit(0 if success else 1)
    
    if MODELS[args.tier].get("provider") == "openrouter" and not os.environ.get("OPENROUTER_API_KEY"):
        cli.print_status("OPENROUTER_API_KEY not set", False)
        sys.exit(1)
    
    context = {}
    if args.spec: context["task_spec"] = Path(args.spec).read_text()
    if args.context:
        with open(args.context) as f: context.update(json.load(f))
    
    result = cli.execute_task(args.task, args.tier, context, args.output)
    cli.print_metrics()
    sys.exit(0 if result.get("success") else 1)

if __name__ == "__main__":
    main()
