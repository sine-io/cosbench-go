#!/usr/bin/env python3

import json
import sys

from worktree_output import (
    branch_name,
    build_single_view_payload,
    build_audit_row,
    build_audit_summary,
    current_worktree,
    generated_at,
    load_worktree_entries,
    parse_known_flags,
    print_text_header,
    should_include_audit_row,
)


def sort_key(row):
    state_rank = {
        "merged": 0,
        "integrated": 1,
        "active": 2,
        "detached": 3,
        "unknown": 4,
    }
    return (
        state_rank.get(row["state"], 9),
        -row["behind"],
        row["branch"],
    )

def main():
    flags, args = parse_known_flags(
        sys.argv[1:],
        ("--json", "--merged-only", "--integrated-only", "--prune-only", "--stale-only"),
    )
    if len(args) > 1:
        raise SystemExit("usage: worktree_audit.py [--json] [--merged-only|--integrated-only|--prune-only|--stale-only] [base_ref]")
    json_mode = flags["--json"]
    merged_only = flags["--merged-only"]
    integrated_only = flags["--integrated-only"]
    prune_only = flags["--prune-only"]
    stale_only = flags["--stale-only"]
    selected_views = sum(1 for enabled in (merged_only, integrated_only, prune_only, stale_only) if enabled)
    if selected_views > 1:
        raise SystemExit("choose at most one of --merged-only, --integrated-only, --prune-only, or --stale-only")
    base_ref = args[0] if args else "origin/main"
    current_worktree_path = current_worktree()

    rows = []
    for entry in load_worktree_entries():
        row = build_audit_row(entry, base_ref, current_worktree_path)
        if not should_include_audit_row(
            row["state"],
            row["branch"],
            row["path"],
            row["behind"],
            current_worktree_path,
            merged_only=merged_only,
            integrated_only=integrated_only,
            prune_only=prune_only,
            stale_only=stale_only,
        ):
            continue
        rows.append(row)

    rows.sort(key=sort_key)
    audit_generated_at = generated_at()

    if json_mode:
        summary = build_audit_summary(rows, base_ref, current_worktree_path)
        print(json.dumps(build_single_view_payload(generated_at(), base_ref, current_worktree_path, "audit", summary, rows), indent=2))
        return

    print_text_header(audit_generated_at, base_ref, current_worktree_path)
    print("PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS")
    for row in rows:
        current = "yes" if row["current"] else "no"
        print(f"{row['path']}\t{row['branch']}\t{current}\t{row['state']}\t{row['details']}")


if __name__ == "__main__":
    main()
