#!/usr/bin/env python3
"""Enhanced cost tracking and budget management."""

from typing import Dict, Optional, List
from dataclasses import dataclass, field
from datetime import datetime, timezone
from pathlib import Path
import json


@dataclass
class TokenCount:
    """Token count for a request/response."""
    prompt_tokens: int = 0
    completion_tokens: int = 0
    total_tokens: int = 0
    
    def __post_init__(self):
        if self.total_tokens == 0:
            self.total_tokens = self.prompt_tokens + self.completion_tokens


@dataclass
class CostEntry:
    """Single cost entry."""
    timestamp: str
    task_id: str
    tier: str
    model: str
    tokens: TokenCount
    cost_usd: float
    duration_seconds: float


@dataclass
class Budget:
    """Budget configuration."""
    daily_limit_usd: float = 10.0
    task_limit_usd: float = 1.0
    warning_threshold: float = 0.8  # 80% of limit
    
    def is_exceeded(self, current: float) -> bool:
        return current >= self.daily_limit_usd
    
    def is_warning(self, current: float) -> bool:
        return current >= (self.daily_limit_usd * self.warning_threshold)


class CostTracker:
    """Tracks costs and enforces budgets."""
    
    # Cost per 1K tokens (approximate OpenRouter pricing)
    MODEL_COSTS = {
        "qwen/qwen3.5-397b-a17b": {"prompt": 0.000001, "completion": 0.000001},
        "x-ai/grok-4.1-fast": {"prompt": 0.000002, "completion": 0.000006},
        "minimax/minimax-m2.7": {"prompt": 0.0000002, "completion": 0.0000006},
        "anthropic/claude-sonnet-4.6": {"prompt": 0.000003, "completion": 0.000015},
        "anthropic/claude-opus-4.6": {"prompt": 0.000015, "completion": 0.000075},
    }
    
    def __init__(self, budget: Optional[Budget] = None):
        self.budget = budget or Budget()
        self.entries: List[CostEntry] = []
        self.daily_total = 0.0
        self.task_totals: Dict[str, float] = {}
        self.tier_totals: Dict[str, float] = {}
    
    def calculate_cost(
        self,
        model: str,
        tokens: TokenCount
    ) -> float:
        """Calculate cost for given tokens and model."""
        pricing = self.MODEL_COSTS.get(model, {"prompt": 0.000001, "completion": 0.000001})
        
        prompt_cost = (tokens.prompt_tokens / 1000) * pricing["prompt"]
        completion_cost = (tokens.completion_tokens / 1000) * pricing["completion"]
        
        return prompt_cost + completion_cost
    
    def record(
        self,
        task_id: str,
        tier: str,
        model: str,
        tokens: TokenCount,
        duration: float
    ) -> CostEntry:
        """Record a cost entry."""
        cost = self.calculate_cost(model, tokens)
        
        entry = CostEntry(
            timestamp=datetime.now(timezone.utc).isoformat(),
            task_id=task_id,
            tier=tier,
            model=model,
            tokens=tokens,
            cost_usd=cost,
            duration_seconds=duration
        )
        
        self.entries.append(entry)
        self.daily_total += cost
        self.task_totals[task_id] = self.task_totals.get(task_id, 0) + cost
        self.tier_totals[tier] = self.tier_totals.get(tier, 0) + cost
        
        # Check budget
        if self.budget.is_exceeded(self.daily_total):
            raise BudgetExceededError(
                f"Daily budget exceeded: ${self.daily_total:.4f} / ${self.budget.daily_limit_usd:.2f}"
            )
        
        if self.budget.is_warning(self.daily_total):
            print(f"⚠️  Budget warning: ${self.daily_total:.4f} / ${self.budget.daily_limit_usd:.2f}")
        
        return entry
    
    def get_daily_total(self) -> float:
        """Get total cost for current session."""
        return self.daily_total
    
    def get_task_total(self, task_id: str) -> float:
        """Get total cost for specific task."""
        return self.task_totals.get(task_id, 0.0)
    
    def get_summary(self) -> Dict:
        """Get cost summary."""
        return {
            "daily_total": self.daily_total,
            "budget_limit": self.budget.daily_limit_usd,
            "budget_remaining": self.budget.daily_limit_usd - self.daily_total,
            "budget_used_percent": (self.daily_total / self.budget.daily_limit_usd * 100) if self.budget.daily_limit_usd > 0 else 0,
            "task_totals": self.task_totals,
            "tier_totals": self.tier_totals,
            "total_requests": len(self.entries)
        }
    
    def save_report(self, filepath: str = None) -> Path:
        """Save cost report to JSON file."""
        if not filepath:
            filepath = f"cost_report_{datetime.now().strftime('%Y%m%d')}.json"
        
        report = {
            "generated": datetime.now(timezone.utc).isoformat(),
            "summary": self.get_summary(),
            "entries": [
                {
                    "timestamp": e.timestamp,
                    "task_id": e.task_id,
                    "tier": e.tier,
                    "model": e.model,
                    "tokens": {
                        "prompt": e.tokens.prompt_tokens,
                        "completion": e.tokens.completion_tokens,
                        "total": e.tokens.total_tokens
                    },
                    "cost_usd": e.cost_usd,
                    "duration": e.duration_seconds
                }
                for e in self.entries
            ]
        }
        
        path = Path(filepath)
        with open(path, 'w') as f:
            json.dump(report, f, indent=2)
        
        return path


class BudgetExceededError(Exception):
    """Raised when budget is exceeded."""
    pass
