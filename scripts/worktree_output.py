from datetime import datetime, timezone
import json
from pathlib import Path
import subprocess


def run_git(*args):
    return subprocess.run(
        ["git", *args],
        check=False,
        text=True,
        capture_output=True,
    )


def generated_at():
    return datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")


def current_worktree():
    proc = run_git("rev-parse", "--show-toplevel")
    return proc.stdout.strip() if proc.returncode == 0 else ""


def build_meta(generated_at_value, base_ref, current_worktree_path):
    return {
        "generated_at": generated_at_value,
        "base_ref": base_ref,
        "current_worktree": current_worktree_path,
    }


def print_text_header(generated_at_value, base_ref, current_worktree_path):
    print(f"# Generated at: {generated_at_value}")
    print(f"# Base ref: {base_ref}")
    print(f"# Current worktree: {current_worktree_path}")


def build_single_view_payload(generated_at_value, base_ref, current_worktree_path, view_name, summary, rows):
    meta = build_meta(generated_at_value, base_ref, current_worktree_path)
    view = {"summary": summary, "rows": rows}
    return {
        "generated_at": meta["generated_at"],
        "meta": meta,
        "views": {view_name: view},
        "summary": summary,
        "rows": rows,
    }


def is_prune_candidate(state, branch, path, current_worktree_path):
    return (
        state in ("merged", "integrated")
        and branch not in ("main", "master")
        and bool(path)
        and path != current_worktree_path
    )


def is_stale_row(state, behind):
    return state == "active" and behind > 0


def build_prune_plan_row(path, branch, state, details, ahead, behind):
    return {
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


def markdown_text_section(title, content):
    return [
        "",
        f"## {title}",
        "",
        "```text",
        content,
        "```",
    ]


def script_path(name):
    return str(Path(__file__).resolve().parent / name)


def run_script(name, *args):
    proc = subprocess.run(
        ["python3", script_path(name), *args],
        check=True,
        text=True,
        capture_output=True,
    )
    return proc.stdout


def load_json_script(name, *args):
    return json.loads(run_script(name, *args))


def load_worktree_audit_text(base_ref, *flags):
    return run_script("worktree_audit.py", *flags, base_ref).rstrip()


def load_worktree_audit_view(base_ref, *flags):
    return load_json_script("worktree_audit.py", "--json", *flags, base_ref)
