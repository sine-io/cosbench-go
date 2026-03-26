#!/usr/bin/env python3

import json
import sys

from worktree_output import (
    build_prune_plan_row,
    build_single_view_payload,
    current_worktree,
    generated_at,
    is_prune_candidate,
    load_json_script,
    parse_known_flags,
    print_text_header,
    validate_base_ref,
)


def main():
    flags, args = parse_known_flags(sys.argv[1:], ("--json",))
    if len(args) > 1:
        raise SystemExit("usage: worktree_prune_plan.py [--json] [base_ref]")
    json_mode = flags["--json"]
    base_ref = args[0] if args else "origin/main"
    validate_base_ref(base_ref)
    current_worktree_path = current_worktree()
    payload = load_json_script("worktree_audit.py", "--json", base_ref)
    source_rows = payload.get("rows", payload)

    rows = []
    for row in source_rows:
        state = row.get("state", "")
        branch = row.get("branch", "")
        path = row.get("path", "")
        details = row.get("details", "")
        ahead = row.get("ahead", 0)
        behind = row.get("behind", 0)
        if not is_prune_candidate(state, branch, path, current_worktree_path):
            continue
        rows.append(build_prune_plan_row(path, branch, state, details, ahead, behind))

    plan_generated_at = generated_at()

    if json_mode:
        summary = {
            "base_ref": base_ref,
            "current_worktree": current_worktree_path,
            "total": len(rows),
            "merged": sum(1 for row in rows if row["state"] == "merged"),
            "integrated": sum(1 for row in rows if row["state"] == "integrated"),
        }
        print(
            json.dumps(
                build_single_view_payload(plan_generated_at, base_ref, current_worktree_path, "prune_plan", summary, rows),
                indent=2,
            )
        )
        return

    print("# Suggested cleanup commands")
    print_text_header(plan_generated_at, base_ref, current_worktree_path)
    if not rows:
        print("# no prune-candidate worktrees to prune")
        return
    for row in rows:
        for command in row["commands"]:
            print(command)


if __name__ == "__main__":
    main()
