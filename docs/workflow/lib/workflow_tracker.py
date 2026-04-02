#!/usr/bin/env python3
"""
Workflow Tracker for Multi-Tier LLM Development Cascade

Logs handoffs, escalations, and learnings to JSONL files for audit and analysis.
"""

import json
import os
from datetime import datetime
from pathlib import Path
from typing import Optional, Dict, Any, List


class WorkflowTracker:
    """Tracks workflow state across L0-L3 tier cascade."""
    
    def __init__(self, base_dir: str = "docs/workflow"):
        self.base_dir = Path(base_dir)
        self.handoffs_dir = self.base_dir / "handoffs"
        self.escalations_dir = self.base_dir / "escalations"
        self.learnings_dir = self.base_dir / "learnings"
        
        # Ensure directories exist
        self.handoffs_dir.mkdir(parents=True, exist_ok=True)
        self.escalations_dir.mkdir(parents=True, exist_ok=True)
        self.learnings_dir.mkdir(parents=True, exist_ok=True)
    
    def log_handoff(
        self,
        task_id: str,
        from_tier: str,
        to_tier: str,
        handoff_type: str,
        details: Dict[str, Any],
        timestamp: Optional[str] = None
    ) -> str:
        """
        Log a handoff between tiers.
        
        Args:
            task_id: Task identifier (e.g., "1.1")
            from_tier: Source tier (e.g., "Planner", "L0")
            to_tier: Destination tier (e.g., "L0", "Reviewer")
            handoff_type: Type of handoff (e.g., "Task Ready", "Review Ready")
            details: Handoff-specific data
            timestamp: ISO format timestamp (auto-generated if not provided)
        
        Returns:
            Path to the log file
        """
        if timestamp is None:
            timestamp = datetime.utcnow().isoformat() + "Z"
        
        log_entry = {
            "timestamp": timestamp,
            "task_id": task_id,
            "from_tier": from_tier,
            "to_tier": to_tier,
            "handoff_type": handoff_type,
            "details": details
        }
        
        # Create log filename
        safe_task_id = task_id.replace(".", "_")
        log_file = self.handoffs_dir / f"{safe_task_id}-{handoff_type.replace(' ', '-').lower()}.json"
        
        # Write log (overwrite if exists - each handoff type is unique per task)
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        
        return str(log_file)
    
    def log_escalation(
        self,
        task_id: str,
        from_tier: str,
        to_tier: str,
        escalation_level: str,
        reason: str,
        attempt_history: List[Dict[str, str]],
        timestamp: Optional[str] = None
    ) -> str:
        """
        Log an escalation event.
        
        Args:
            task_id: Task identifier
            from_tier: Escalating tier
            to_tier: Receiving tier
            escalation_level: "L1", "L2", or "L3"
            reason: Why this was escalated
            attempt_history: List of previous attempts with outcomes
            timestamp: ISO format timestamp
        
        Returns:
            Path to the log file
        """
        if timestamp is None:
            timestamp = datetime.utcnow().isoformat() + "Z"
        
        log_entry = {
            "timestamp": timestamp,
            "task_id": task_id,
            "escalation_level": escalation_level,
            "from_tier": from_tier,
            "to_tier": to_tier,
            "reason": reason,
            "attempt_history": attempt_history
        }
        
        # Create log filename
        safe_task_id = task_id.replace(".", "_")
        log_file = self.escalations_dir / f"{safe_task_id}-{escalation_level.lower()}.json"
        
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        
        return str(log_file)
    
    def log_learning(
        self,
        task_id: str,
        learning_type: str,
        category: str,
        content: Dict[str, Any],
        timestamp: Optional[str] = None
    ) -> str:
        """
        Log a learning document.
        
        Args:
            task_id: Task identifier
            learning_type: "Escalation Resolution" or "Checkpoint Review"
            category: Root cause category
            content: Learning content
            timestamp: ISO format timestamp
        
        Returns:
            Path to the log file
        """
        if timestamp is None:
            timestamp = datetime.utcnow().isoformat() + "Z"
        
        log_entry = {
            "timestamp": timestamp,
            "task_id": task_id,
            "learning_type": learning_type,
            "category": category,
            "content": content
        }
        
        # Create log filename
        safe_task_id = task_id.replace(".", "_")
        log_file = self.learnings_dir / f"{safe_task_id}-{learning_type.replace(' ', '-').lower()}.json"
        
        with open(log_file, 'w') as f:
            json.dump(log_entry, f, indent=2)
        
        return str(log_file)
    
    def get_task_history(self, task_id: str) -> Dict[str, Any]:
        """
        Get complete history for a task.
        
        Args:
            task_id: Task identifier
        
        Returns:
            Dictionary with handoffs, escalations, and learnings for the task
        """
        safe_task_id = task_id.replace(".", "_")
        
        history = {
            "task_id": task_id,
            "handoffs": [],
            "escalations": [],
            "learnings": []
        }
        
        # Find handoffs
        for log_file in self.handoffs_dir.glob(f"{safe_task_id}-*.json"):
            with open(log_file) as f:
                history["handoffs"].append(json.load(f))
        
        # Find escalations
        for log_file in self.escalations_dir.glob(f"{safe_task_id}-*.json"):
            with open(log_file) as f:
                history["escalations"].append(json.load(f))
        
        # Find learnings
        for log_file in self.learnings_dir.glob(f"{safe_task_id}-*.json"):
            with open(log_file) as f:
                history["learnings"].append(json.load(f))
        
        # Sort by timestamp
        history["handoffs"].sort(key=lambda x: x["timestamp"])
        history["escalations"].sort(key=lambda x: x["timestamp"])
        history["learnings"].sort(key=lambda x: x["timestamp"])
        
        return history
    
    def get_workflow_metrics(self) -> Dict[str, Any]:
        """
        Calculate workflow-wide metrics.
        
        Returns:
            Dictionary with aggregate metrics
        """
        # Count all handoffs
        handoff_files = list(self.handoffs_dir.glob("*.json"))
        handoff_count = len(handoff_files)
        
        # Count escalations by level
        escalation_counts = {"L1": 0, "L2": 0, "L3": 0}
        for log_file in self.escalations_dir.glob("*.json"):
            with open(log_file) as f:
                entry = json.load(f)
                level = entry.get("escalation_level", "unknown")
                if level in escalation_counts:
                    escalation_counts[level] += 1
        
        # Count learnings
        learning_count = len(list(self.learnings_dir.glob("*.json")))
        
        # Calculate escalation rate (approximate - would need task count)
        total_escalations = sum(escalation_counts.values())
        
        return {
            "total_handoffs": handoff_count,
            "total_escalations": total_escalations,
            "escalation_by_level": escalation_counts,
            "total_learnings": learning_count,
            "escalation_rate": f"{total_escalations}/{handoff_count} (approximate)" if handoff_count > 0 else "0/0"
        }


def main():
    """CLI entry point for workflow tracker."""
    import argparse
    
    parser = argparse.ArgumentParser(description="Workflow Tracker CLI")
    parser.add_argument("--command", choices=["history", "metrics", "log-handoff", "log-escalation"], required=True)
    parser.add_argument("--task-id", help="Task ID for history or logging")
    parser.add_argument("--from-tier", help="Source tier for logging")
    parser.add_argument("--to-tier", help="Destination tier for logging")
    parser.add_argument("--type", help="Handoff or escalation type")
    parser.add_argument("--level", help="Escalation level (L1/L2/L3)")
    parser.add_argument("--reason", help="Escalation reason")
    
    args = parser.parse_args()
    
    tracker = WorkflowTracker()
    
    if args.command == "history":
        if not args.task_id:
            print("Error: --task-id required for history command")
            return
        history = tracker.get_task_history(args.task_id)
        print(json.dumps(history, indent=2))
    
    elif args.command == "metrics":
        metrics = tracker.get_workflow_metrics()
        print(json.dumps(metrics, indent=2))
    
    elif args.command == "log-handoff":
        # Interactive or JSON input for handoff details
        print("Handoff logging - provide details via stdin JSON")
        details = json.loads(input())
        log_file = tracker.log_handoff(
            task_id=args.task_id,
            from_tier=args.from_tier,
            to_tier=args.to_tier,
            handoff_type=args.type,
            details=details
        )
        print(f"Logged to: {log_file}")
    
    elif args.command == "log-escalation":
        print("Escalation logging - provide attempt history via stdin JSON")
        attempt_history = json.loads(input())
        log_file = tracker.log_escalation(
            task_id=args.task_id,
            from_tier=args.from_tier,
            to_tier=args.to_tier,
            escalation_level=args.level,
            reason=args.reason,
            attempt_history=attempt_history
        )
        print(f"Logged to: {log_file}")


if __name__ == "__main__":
    main()
