#!/usr/bin/env python3

import sys

from worktree_output import build_meta, generated_at, load_json_script, run_script

def main():
    json_mode = "--json" in sys.argv[1:]
    args = [arg for arg in sys.argv[1:] if arg != "--json"]
    base_ref = args[0] if args else "origin/main"
    output_path = sys.argv[2] if len(sys.argv) > 2 else ""
    audit = load_json_script("worktree_audit.py", "--json", base_ref)
    prune_plan = load_json_script("worktree_prune_plan.py", "--json", base_ref)
    merged_text = run_script("worktree_audit.py", "--merged-only", base_ref).rstrip()
    integrated_text = run_script("worktree_audit.py", "--integrated-only", base_ref).rstrip()
    prune_candidates_text = run_script("worktree_audit.py", "--prune-only", base_ref).rstrip()
    stale_text = run_script("worktree_audit.py", "--stale-only", base_ref).rstrip()
    prune_text = run_script("worktree_prune_plan.py", base_ref).rstrip()

    summary = audit["summary"]
    report_generated_at = generated_at()
    current_worktree = prune_plan["summary"].get("current_worktree", "")
    if json_mode:
        merged_view = load_json_script("worktree_audit.py", "--json", "--merged-only", base_ref)
        integrated_view = load_json_script("worktree_audit.py", "--json", "--integrated-only", base_ref)
        stale_view = load_json_script("worktree_audit.py", "--json", "--stale-only", base_ref)
        prune_candidates_view = load_json_script("worktree_audit.py", "--json", "--prune-only", base_ref)
        meta = build_meta(report_generated_at, summary["base_ref"], current_worktree)
        payload = {
            "generated_at": report_generated_at,
            "meta": meta,
            "summary": summary,
            "views": {
                "merged": merged_view,
                "integrated": integrated_view,
                "stale": stale_view,
                "prune_candidates": prune_candidates_view,
                "prune_plan": prune_plan,
            },
            "merged": merged_view,
            "integrated": integrated_view,
            "stale": stale_view,
            "prune_candidates": prune_candidates_view,
            "prune_plan": prune_plan,
        }
        print(json.dumps(payload, indent=2))
        return

    lines = [
        "# Worktree Cleanup Report",
        "",
        "## Summary",
        "",
        f"- Generated at: `{report_generated_at}`",
        f"- Base ref: `{summary['base_ref']}`",
        f"- Current worktree: `{current_worktree}`",
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
        "## Integrated",
        "",
        "```text",
        integrated_text,
        "```",
        "",
        "## Stale",
        "",
        "```text",
        stale_text,
        "```",
        "",
        "## Prune Candidates",
        "",
        "```text",
        prune_candidates_text,
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
