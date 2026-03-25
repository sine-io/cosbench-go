#!/usr/bin/env python3

import json
import subprocess
import sys


def main():
    json_mode = "--json" in sys.argv[1:]
    args = [arg for arg in sys.argv[1:] if arg != "--json"]
    base_ref = args[0] if args else "origin/main"
    cwd_proc = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        check=True,
        text=True,
        capture_output=True,
    )
    current_worktree = cwd_proc.stdout.strip()
    proc = subprocess.run(
        ["python3", "./scripts/worktree_audit.py", "--json", "--merged-only", base_ref],
        check=True,
        text=True,
        capture_output=True,
    )
    payload = json.loads(proc.stdout)
    source_rows = payload.get("rows", payload)

    rows = []
    for row in source_rows:
        branch = row.get("branch", "")
        path = row.get("path", "")
        if branch in ("main", "master") or not path or path == current_worktree:
            continue
        rows.append(
            {
                "path": path,
                "branch": branch,
                "commands": [
                    f"git worktree remove '{path}'",
                    f"git branch -D {branch}",
                ],
            }
        )

    if json_mode:
        print(json.dumps(rows, indent=2))
        return

    print("# Suggested cleanup commands")
    if not rows:
        print("# no merged worktrees to prune")
        return
    for row in rows:
        for command in row["commands"]:
            print(command)


if __name__ == "__main__":
    main()
