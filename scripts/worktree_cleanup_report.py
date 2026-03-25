#!/usr/bin/env python3

import json
import subprocess
import sys


def run(*args):
    proc = subprocess.run(args, check=True, text=True, capture_output=True)
    return proc.stdout


def main():
    json_mode = "--json" in sys.argv[1:]
    args = [arg for arg in sys.argv[1:] if arg != "--json"]
    base_ref = args[0] if args else "origin/main"
    output_path = sys.argv[2] if len(sys.argv) > 2 else ""
    audit = json.loads(run("python3", "./scripts/worktree_audit.py", "--json", base_ref))
    prune_plan = json.loads(run("python3", "./scripts/worktree_prune_plan.py", "--json", base_ref))
    merged_text = run("python3", "./scripts/worktree_audit.py", "--merged-only", base_ref).rstrip()
    stale_text = run("python3", "./scripts/worktree_audit.py", "--stale-only", base_ref).rstrip()
    prune_text = run("python3", "./scripts/worktree_prune_plan.py", base_ref).rstrip()

    summary = audit["summary"]
    if json_mode:
        payload = {
            "summary": summary,
            "merged": json.loads(run("python3", "./scripts/worktree_audit.py", "--json", "--merged-only", base_ref)),
            "stale": json.loads(run("python3", "./scripts/worktree_audit.py", "--json", "--stale-only", base_ref)),
            "prune_plan": prune_plan,
        }
        print(json.dumps(payload, indent=2))
        return

    lines = [
        "# Worktree Cleanup Report",
        "",
        "## Summary",
        "",
        f"- Base ref: `{summary['base_ref']}`",
        f"- Total worktrees: {summary['total']}",
        f"- Merged: {summary['merged']}",
        f"- Integrated: {summary['integrated']}",
        f"- Stale: {summary['stale']}",
        f"- Prune candidates: {prune_plan['summary']['total']}",
        f"- Active: {summary['active']}",
        f"- Detached: {summary['detached']}",
        f"- Unknown: {summary['unknown']}",
        "",
        "## Merged",
        "",
        "```text",
        merged_text,
        "```",
        "",
        "## Stale",
        "",
        "```text",
        stale_text,
        "```",
        "",
        "## Prune Plan",
        "",
        "```text",
        prune_text,
        "```",
    ]
    report = "\n".join(lines) + "\n"
    if output_path:
        with open(output_path, "w", encoding="utf-8") as fh:
            fh.write(report)
    print(report, end="")


if __name__ == "__main__":
    main()
