#!/usr/bin/env python3

import json
import subprocess
import sys


def run(*args):
    return subprocess.run(args, check=False, text=True, capture_output=True)


def worktree_entries():
    proc = run("git", "worktree", "list", "--porcelain")
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
    merged = run("git", "merge-base", "--is-ancestor", branch, base_ref)
    if merged.returncode == 0:
        return "merged", base_ref, 0, 0
    ahead_behind = run("git", "rev-list", "--left-right", "--count", f"{base_ref}...{branch}")
    if ahead_behind.returncode != 0:
        return "unknown", ahead_behind.stderr.strip() or ahead_behind.stdout.strip(), 0, 0
    behind, ahead = ahead_behind.stdout.strip().split()
    return "active", f"ahead={ahead} behind={behind}", int(ahead), int(behind)


def sort_key(row):
    state_rank = {
        "merged": 0,
        "active": 1,
        "detached": 2,
        "unknown": 3,
    }
    return (
        state_rank.get(row["state"], 9),
        -row["behind"],
        row["branch"],
    )


def main():
    json_mode = "--json" in sys.argv[1:]
    merged_only = "--merged-only" in sys.argv[1:]
    stale_only = "--stale-only" in sys.argv[1:]
    args = [arg for arg in sys.argv[1:] if arg not in ("--json", "--merged-only", "--stale-only")]
    base_ref = args[0] if args else "origin/main"
    current_proc = run("git", "rev-parse", "--show-toplevel")
    current_worktree = current_proc.stdout.strip() if current_proc.returncode == 0 else ""

    rows = []
    for entry in worktree_entries():
        branch = branch_name(entry)
        state, details, ahead, behind = classify(branch, base_ref)
        if merged_only and state != "merged":
            continue
        if stale_only and not (state == "active" and behind > 0):
            continue
        rows.append(
            {
                "path": entry["worktree"],
                "branch": branch,
                "state": state,
                "details": details,
                "ahead": ahead,
                "behind": behind,
                "current": entry["worktree"] == current_worktree,
            }
        )

    rows.sort(key=sort_key)

    if json_mode:
        summary = {
            "base_ref": base_ref,
            "total": len(rows),
            "merged": sum(1 for row in rows if row["state"] == "merged"),
            "active": sum(1 for row in rows if row["state"] == "active"),
            "detached": sum(1 for row in rows if row["state"] == "detached"),
            "unknown": sum(1 for row in rows if row["state"] == "unknown"),
        }
        print(json.dumps({"summary": summary, "rows": rows}, indent=2))
        return

    print("PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS")
    for row in rows:
        current = "yes" if row["current"] else "no"
        print(f"{row['path']}\t{row['branch']}\t{current}\t{row['state']}\t{row['details']}")


if __name__ == "__main__":
    main()
