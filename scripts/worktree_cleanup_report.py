#!/usr/bin/env python3

import json
import subprocess


def run(*args):
    proc = subprocess.run(args, check=True, text=True, capture_output=True)
    return proc.stdout


def main():
    audit = json.loads(run("python3", "./scripts/worktree_audit.py", "--json", "origin/main"))
    merged = run("python3", "./scripts/worktree_audit.py", "--merged-only", "origin/main").rstrip()
    stale = run("python3", "./scripts/worktree_audit.py", "--stale-only", "origin/main").rstrip()
    prune = run("python3", "./scripts/worktree_prune_plan.py", "origin/main").rstrip()

    summary = audit["summary"]
    lines = [
        "# Worktree Cleanup Report",
        "",
        "## Summary",
        "",
        f"- Base ref: `{summary['base_ref']}`",
        f"- Total worktrees: {summary['total']}",
        f"- Merged: {summary['merged']}",
        f"- Active: {summary['active']}",
        f"- Detached: {summary['detached']}",
        f"- Unknown: {summary['unknown']}",
        "",
        "## Merged",
        "",
        "```text",
        merged,
        "```",
        "",
        "## Stale",
        "",
        "```text",
        stale,
        "```",
        "",
        "## Prune Plan",
        "",
        "```text",
        prune,
        "```",
    ]
    print("\n".join(lines))


if __name__ == "__main__":
    main()
