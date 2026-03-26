#!/usr/bin/env python3

import json
import sys

from worktree_output import (
    build_single_view_payload,
    current_worktree,
    generated_at,
    is_prune_candidate,
    is_stale_row,
    print_text_header,
    run_git,
)


def worktree_entries():
    proc = run_git("worktree", "list", "--porcelain")
    if proc.returncode != 0:
        raise SystemExit(proc.stderr or proc.stdout)
    entry = {}
    for raw_line in proc.stdout.splitlines():
        if not raw_line:
            if entry:
                yield entry
                entry = {}
            continue
        if " " in raw_line:
            key, value = raw_line.split(" ", 1)
        else:
            key, value = raw_line, ""
        entry[key] = value
    if entry:
        yield entry


def branch_name(entry):
    ref = entry.get("branch", "")
    if ref.startswith("refs/heads/"):
        return ref.removeprefix("refs/heads/")
    return "(detached)"


def classify(branch, base_ref):
    if branch == "(detached)":
        return "detached", "", 0, 0
    merged = run_git("merge-base", "--is-ancestor", branch, base_ref)
    if merged.returncode == 0:
        return "merged", base_ref, 0, 0
    cherry = run_git("cherry", base_ref, branch)
    cherry_lines = [line for line in cherry.stdout.splitlines() if line.strip()]
    if cherry.returncode == 0 and cherry_lines and all(line.startswith("- ") for line in cherry_lines):
        ahead_behind = run_git("rev-list", "--left-right", "--count", f"{base_ref}...{branch}")
        if ahead_behind.returncode == 0:
            behind, ahead = ahead_behind.stdout.strip().split()
            return "integrated", f"patch-equivalent to {base_ref}", int(ahead), int(behind)
        return "integrated", f"patch-equivalent to {base_ref}", 0, 0
    ahead_behind = run_git("rev-list", "--left-right", "--count", f"{base_ref}...{branch}")
    if ahead_behind.returncode != 0:
        return "unknown", ahead_behind.stderr.strip() or ahead_behind.stdout.strip(), 0, 0
    behind, ahead = ahead_behind.stdout.strip().split()
    return "active", f"ahead={ahead} behind={behind}", int(ahead), int(behind)


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
    for entry in worktree_entries():
        branch = branch_name(entry)
        state, details, ahead, behind = classify(branch, base_ref)
        if merged_only and state != "merged":
            continue
        if integrated_only and state != "integrated":
            continue
        if prune_only and not is_prune_candidate(state, branch, entry["worktree"], current_worktree_path):
            continue
        if stale_only and not is_stale_row(state, behind):
            continue
        rows.append(
            {
                "path": entry["worktree"],
                "branch": branch,
                "state": state,
                "details": details,
                "ahead": ahead,
                "behind": behind,
                "current": entry["worktree"] == current_worktree_path,
            }
        )

    rows.sort(key=sort_key)
    audit_generated_at = generated_at()

    if json_mode:
        summary = {
            "base_ref": base_ref,
            "total": len(rows),
            "merged": sum(1 for row in rows if row["state"] == "merged"),
            "integrated": sum(1 for row in rows if row["state"] == "integrated"),
            "active": sum(1 for row in rows if row["state"] == "active"),
            "detached": sum(1 for row in rows if row["state"] == "detached"),
            "unknown": sum(1 for row in rows if row["state"] == "unknown"),
            "stale": sum(1 for row in rows if is_stale_row(row["state"], row["behind"])),
            "prune_candidates": sum(1 for row in rows if is_prune_candidate(row["state"], row["branch"], row["path"], current_worktree_path)),
        }
        print(json.dumps(build_single_view_payload(generated_at(), base_ref, current_worktree_path, "audit", summary, rows), indent=2))
        return

    print_text_header(audit_generated_at, base_ref, current_worktree_path)
    print("PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS")
    for row in rows:
        current = "yes" if row["current"] else "no"
        print(f"{row['path']}\t{row['branch']}\t{current}\t{row['state']}\t{row['details']}")


if __name__ == "__main__":
    main()
