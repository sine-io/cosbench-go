#!/usr/bin/env python3

import json
import subprocess


def main():
    proc = subprocess.run(
        ["python3", "./scripts/worktree_audit.py", "--json", "--merged-only", "origin/main"],
        check=True,
        text=True,
        capture_output=True,
    )
    payload = json.loads(proc.stdout)
    rows = payload.get("rows", payload)

    print("# Suggested cleanup commands")
    printed = False
    for row in rows:
        branch = row.get("branch", "")
        path = row.get("path", "")
        if branch in ("main", "master") or not path:
            continue
        print(f"git worktree remove '{path}'")
        print(f"git branch -D {branch}")
        printed = True
    if not printed:
        print("# no merged worktrees to prune")


if __name__ == "__main__":
    main()
