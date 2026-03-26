#!/usr/bin/env python3

import json
import sys

from worktree_output import (
    branch_name,
    build_single_view_payload,
    classify_branch,
    current_worktree,
    generated_at,
    is_prune_candidate,
    is_stale_row,
    load_worktree_entries,
    print_text_header,
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


def should_include_row(state, branch, path, behind, current_worktree_path, *, merged_only, integrated_only, prune_only, stale_only):
    if merged_only and state != "merged":
        return False
    if integrated_only and state != "integrated":
        return False
    if prune_only and not is_prune_candidate(state, branch, path, current_worktree_path):
        return False
    if stale_only and not is_stale_row(state, behind):
        return False
    return True


def build_audit_row(entry, base_ref, current_worktree_path):
    branch = branch_name(entry)
    state, details, ahead, behind = classify_branch(branch, base_ref)
    path = entry["worktree"]
    return {
        "path": path,
        "branch": branch,
        "state": state,
        "details": details,
        "ahead": ahead,
        "behind": behind,
        "current": path == current_worktree_path,
    }


def build_audit_summary(rows, base_ref, current_worktree_path):
    return {
        "base_ref": base_ref,
        "total": len(rows),
        "merged": sum(1 for row in rows if row["state"] == "merged"),
        "integrated": sum(1 for row in rows if row["state"] == "integrated"),
        "active": sum(1 for row in rows if row["state"] == "active"),
        "detached": sum(1 for row in rows if row["state"] == "detached"),
        "unknown": sum(1 for row in rows if row["state"] == "unknown"),
        "stale": sum(1 for row in rows if is_stale_row(row["state"], row["behind"])),
        "prune_candidates": sum(
            1 for row in rows if is_prune_candidate(row["state"], row["branch"], row["path"], current_worktree_path)
        ),
    }


def main():
    json_mode = "--json" in sys.argv[1:]
    merged_only = "--merged-only" in sys.argv[1:]
    integrated_only = "--integrated-only" in sys.argv[1:]
    prune_only = "--prune-only" in sys.argv[1:]
    stale_only = "--stale-only" in sys.argv[1:]
    args = [
        arg
        for arg in sys.argv[1:]
        if arg not in ("--json", "--merged-only", "--integrated-only", "--prune-only", "--stale-only")
    ]
    base_ref = args[0] if args else "origin/main"
    current_worktree_path = current_worktree()

    rows = []
    for entry in load_worktree_entries():
        row = build_audit_row(entry, base_ref, current_worktree_path)
        if not should_include_row(
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
