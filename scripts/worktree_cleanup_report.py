#!/usr/bin/env python3

import json
from pathlib import Path
import sys

from worktree_output import (
    build_meta,
    configure_utf8_stdout,
    display_text,
    generated_at,
    load_json_script,
    load_worktree_audit_text,
    load_worktree_audit_view,
    markdown_text_section,
    parse_known_flags,
    run_script,
    resolve_base_ref,
)

def format_os_error(err: OSError) -> str:
    parts = []
    if getattr(err, "errno", None) is not None:
        parts.append(f"[Errno {err.errno}]")
    if getattr(err, "strerror", None):
        parts.append(display_text(str(err.strerror)))
    elif str(err):
        parts.append(display_text(str(err)))
    return " ".join(parts) or err.__class__.__name__


def markdown_code(value: str) -> str:
    longest_run = 0
    current_run = 0
    for ch in value:
        if ch == "`":
            current_run += 1
            if current_run > longest_run:
                longest_run = current_run
        else:
            current_run = 0
    fence = "`" * (longest_run + 1)
    return f"{fence}{value}{fence}"


def main():
    configure_utf8_stdout()
    flags, args = parse_known_flags(sys.argv[1:], ("--json",))
    json_mode = flags["--json"]
    if json_mode and len(args) > 1:
        raise SystemExit(f"expected at most one base_ref argument in --json mode, got: {' '.join(args)}")
    if not json_mode and len(args) > 2:
        raise SystemExit(f"expected at most base_ref and output_path arguments, got: {' '.join(args)}")
    base_ref = resolve_base_ref(args[0] if args else "")
    output_path = args[1] if len(args) > 1 else ""
    audit = load_json_script("worktree_audit.py", "--json", base_ref)
    prune_plan = load_json_script("worktree_prune_plan.py", "--json", base_ref)
    merged_text = load_worktree_audit_text(base_ref, "--merged-only")
    integrated_text = load_worktree_audit_text(base_ref, "--integrated-only")
    prune_candidates_text = load_worktree_audit_text(base_ref, "--prune-only")
    stale_text = load_worktree_audit_text(base_ref, "--stale-only")
    prune_text = run_script("worktree_prune_plan.py", base_ref).rstrip()

    summary = audit["summary"]
    report_generated_at = generated_at()
    current_worktree = prune_plan["summary"].get("current_worktree", "")
    if json_mode:
        merged_view = load_worktree_audit_view(base_ref, "--merged-only")
        integrated_view = load_worktree_audit_view(base_ref, "--integrated-only")
        stale_view = load_worktree_audit_view(base_ref, "--stale-only")
        prune_candidates_view = load_worktree_audit_view(base_ref, "--prune-only")
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
        print(json.dumps(payload, indent=2, ensure_ascii=False))
        return

    lines = [
        "# Worktree Cleanup Report",
        "",
        "## Summary",
        "",
        f"- Generated at: {markdown_code(report_generated_at)}",
        f"- Base ref: {markdown_code(summary['base_ref'])}",
        f"- Current worktree: {markdown_code(current_worktree)}",
        f"- Total worktrees: {summary['total']}",
        f"- Merged: {summary['merged']}",
        f"- Integrated: {summary['integrated']}",
        f"- Stale: {summary['stale']}",
        f"- Prune candidates: {prune_plan['summary']['total']}",
        f"- Active: {summary['active']}",
        f"- Detached: {summary['detached']}",
        f"- Unknown: {summary['unknown']}",
    ]
    lines.extend(markdown_text_section("Merged", merged_text))
    lines.extend(markdown_text_section("Integrated", integrated_text))
    lines.extend(markdown_text_section("Stale", stale_text))
    lines.extend(markdown_text_section("Prune Candidates", prune_candidates_text))
    lines.extend(markdown_text_section("Prune Plan", prune_text))
    report = "\n".join(lines) + "\n"
    if output_path:
        parent_dir = Path(output_path).parent
        try:
            parent_dir.mkdir(parents=True, exist_ok=True)
        except OSError as err:
            raise SystemExit(
                f"unable to prepare worktree cleanup report parent dir {display_text(str(parent_dir))}: {format_os_error(err)}"
            )
        try:
            with open(output_path, "w", encoding="utf-8") as fh:
                fh.write(report)
        except OSError as err:
            raise SystemExit(
                f"unable to write worktree cleanup report {display_text(str(output_path))}: {format_os_error(err)}"
            )
    print(report, end="")


if __name__ == "__main__":
    main()
