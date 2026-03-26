#!/usr/bin/env python3

from datetime import datetime, timezone
import json
import subprocess
import sys
from pathlib import Path


def generated_at():
    return datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")


def main():
    json_mode = "--json" in sys.argv[1:]
    args = [arg for arg in sys.argv[1:] if arg != "--json"]
    base_ref = args[0] if args else "origin/main"
    script_dir = Path(__file__).resolve().parent
    audit_script = str(script_dir / "worktree_audit.py")
    cwd_proc = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        check=True,
        text=True,
        capture_output=True,
    )
    current_worktree = cwd_proc.stdout.strip()
    proc = subprocess.run(
        ["python3", audit_script, "--json", base_ref],
        check=True,
        text=True,
        capture_output=True,
    )
    payload = json.loads(proc.stdout)
    source_rows = payload.get("rows", payload)

    rows = []
    for row in source_rows:
        state = row.get("state", "")
        branch = row.get("branch", "")
        path = row.get("path", "")
        details = row.get("details", "")
        ahead = row.get("ahead", 0)
        behind = row.get("behind", 0)
        if state not in ("merged", "integrated"):
            continue
        if branch in ("main", "master") or not path or path == current_worktree:
            continue
        rows.append(
            {
                "path": path,
                "branch": branch,
                "state": state,
                "details": details,
                "ahead": ahead,
                "behind": behind,
                "commands": [
                    f"git worktree remove '{path}'",
                    f"git branch -D {branch}",
                ],
            }
        )

    plan_generated_at = generated_at()

    if json_mode:
        summary = {
            "base_ref": base_ref,
            "current_worktree": current_worktree,
            "total": len(rows),
            "merged": sum(1 for row in rows if row["state"] == "merged"),
            "integrated": sum(1 for row in rows if row["state"] == "integrated"),
        }
        view = {"summary": summary, "rows": rows}
        meta = {
            "generated_at": plan_generated_at,
            "base_ref": base_ref,
            "current_worktree": current_worktree,
        }
        print(
            json.dumps(
                {
                    "generated_at": plan_generated_at,
                    "meta": meta,
                    "views": {"prune_plan": view},
                    "summary": summary,
                    "rows": rows,
                },
                indent=2,
            )
        )
        return

    print("# Suggested cleanup commands")
    print(f"# Generated at: {plan_generated_at}")
    print(f"# Base ref: {base_ref}")
    print(f"# Current worktree: {current_worktree}")
    if not rows:
        print("# no merged worktrees to prune")
        return
    for row in rows:
        for command in row["commands"]:
            print(command)


if __name__ == "__main__":
    main()
