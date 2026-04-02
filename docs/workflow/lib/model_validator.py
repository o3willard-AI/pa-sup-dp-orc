#!/usr/bin/env python3
"""
Model ID Validator for Multi-Tier Workflow

Validates all configured model IDs against OpenRouter API before starting work.
"""

import requests
import sys
from pathlib import Path

# Add parent directory to path
sys.path.insert(0, str(Path(__file__).parent))

from orchestrator import MODELS

def validate_models(verbose: bool = True) -> tuple[bool, dict]:
    """
    Validate all configured model IDs against OpenRouter.
    
    Returns:
        tuple: (all_valid: bool, results: dict)
    """
    results = {"valid": [], "invalid": [], "unreachable": []}
    
    # Get available models from OpenRouter
    try:
        api_key = input("Enter OpenRouter API key (or set OPENROUTER_API_KEY): ").strip() if not __import__('os').environ.get('OPENROUTER_API_KEY') else __import__('os').environ['OPENROUTER_API_KEY']
        response = requests.get(
            "https://openrouter.ai/api/v1/models",
            headers={"Authorization": f"Bearer {api_key}"},
            timeout=30
        )
        if response.status_code != 200:
            if verbose:
                print(f"Failed to fetch models: {response.status_code}")
            return False, {"error": f"API error: {response.status_code}"}
        
        available_ids = {m['id'] for m in response.json().get('data', [])}
    except Exception as e:
        if verbose:
            print(f"Error fetching models: {e}")
        return False, {"error": str(e)}
    
    # Validate each configured model
    for tier, config in MODELS.items():
        model_id = config.get('model', '')
        provider = config.get('provider', '')
        
        # Skip non-OpenRouter models (LM Studio local)
        if provider != 'openrouter':
            if verbose:
                print(f"✓ {tier}: {model_id} (local/{provider} - not validated)")
            results["valid"].append(tier)
            continue
        
        # Check if model is available
        if model_id in available_ids:
            if verbose:
                print(f"✓ {tier}: {model_id}")
            results["valid"].append(tier)
        else:
            # Check for similar models (typo detection)
            similar = [m for m in available_ids if model_id.split('/')[-1] in m.split('/')[-1]]
            if verbose:
                if similar:
                    print(f"✗ {tier}: {model_id} (INVALID)")
                    print(f"  Did you mean: {', '.join(similar[:3])}")
                else:
                    print(f"✗ {tier}: {model_id} (INVALID - not found)")
            results["invalid"].append({"tier": tier, "model": model_id, "similar": similar[:3]})
    
    all_valid = len(results["invalid"]) == 0 and len(results["unreachable"]) == 0
    
    if verbose:
        print(f"\n{'='*60}")
        print(f"Validation Result: {'PASS' if all_valid else 'FAIL'}")
        print(f"Valid: {len(results['valid'])}, Invalid: {len(results['invalid'])}")
        if results["invalid"]:
            print(f"\nInvalid models must be fixed in orchestrator.py MODELS config")
    
    return all_valid, results


def main():
    import argparse
    parser = argparse.ArgumentParser(description="Validate model IDs")
    parser.add_argument("--quiet", "-q", action="store_true", help="Suppress output")
    parser.add_argument("--json", action="store_true", help="Output as JSON")
    args = parser.parse_args()
    
    import os
    if not os.environ.get('OPENROUTER_API_KEY'):
        print("ERROR: OPENROUTER_API_KEY environment variable not set")
        sys.exit(1)
    
    all_valid, results = validate_models(verbose=not args.quiet)
    
    if args.json:
        import json
        print(json.dumps(results, indent=2))
    
    sys.exit(0 if all_valid else 1)


if __name__ == "__main__":
    main()
