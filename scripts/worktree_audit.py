#!/usr/bin/env python3

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
        key, value = raw_line.split(" ", 1)
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
        return "detached", ""
    merged = run("git", "merge-base", "--is-ancestor", branch, base_ref)
    if merged.returncode == 0:
        return "merged", base_ref
    ahead_behind = run("git", "rev-list", "--left-right", "--count", f"{base_ref}...{branch}")
    if ahead_behind.returncode != 0:
        return "unknown", ahead_behind.stderr.strip() or ahead_behind.stdout.strip()
    behind, ahead = ahead_behind.stdout.strip().split()
    return "active", f"ahead={ahead} behind={behind}"


def main():
    base_ref = sys.argv[1] if len(sys.argv) > 1 else "origin/main"
    print("PATH\tBRANCH\tSTATE\tDETAILS")
    for entry in worktree_entries():
        branch = branch_name(entry)
        state, details = classify(branch, base_ref)
        print(f"{entry['worktree']}\t{branch}\t{state}\t{details}")


if __name__ == "__main__":
    main()
