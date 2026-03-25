#!/usr/bin/env python3

import json
import subprocess
import sys


def run(*args):
    proc = subprocess.run(args, check=True, text=True, capture_output=True)
    return proc.stdout


def main():
    base_ref = sys.argv[1] if len(sys.argv) > 1 else "origin/main"
    output_path = sys.argv[2] if len(sys.argv) > 2 else ""
    audit = json.loads(run("python3", "./scripts/worktree_audit.py", "--json", base_ref))
    merged = run("python3", "./scripts/worktree_audit.py", "--merged-only", base_ref).rstrip()
    stale = run("python3", "./scripts/worktree_audit.py", "--stale-only", base_ref).rstrip()
    prune = run("python3", "./scripts/worktree_prune_plan.py", base_ref).rstrip()

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
    report = "\n".join(lines) + "\n"
    if output_path:
        with open(output_path, "w", encoding="utf-8") as fh:
            fh.write(report)
    print(report, end="")


if __name__ == "__main__":
    main()
