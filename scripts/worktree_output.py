from datetime import datetime, timezone
import json
import os
from pathlib import Path
import shlex
import subprocess
import sys

python_env = {"PYTHONDONTWRITEBYTECODE": "1"}


def configure_utf8_stdout():
    for stream in (sys.stdout, sys.stderr):
        if hasattr(stream, "reconfigure"):
            stream.reconfigure(encoding="utf-8")


def display_text(value: str):
    return value.encode("utf-8", "surrogateescape").decode("utf-8", "replace")


def utf8_argv(*args):
    return [arg.encode("utf-8") if isinstance(arg, str) else arg for arg in args]


def run_git(*args):
    return subprocess.run(
        utf8_argv("git", *args),
        check=False,
        text=True,
        encoding="utf-8",
        capture_output=True,
    )


def parse_porcelain_entries(text):
    entry = {}
    for raw_line in text.splitlines():
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


def generated_at():
    return datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")


def current_worktree():
    proc = run_git("rev-parse", "--show-toplevel")
    return proc.stdout.strip() if proc.returncode == 0 else ""


def current_branch():
    proc = run_git("symbolic-ref", "--quiet", "--short", "HEAD")
    return proc.stdout.strip() if proc.returncode == 0 else ""


def validate_base_ref(base_ref):
    proc = run_git("rev-parse", "--verify", "--quiet", f"{base_ref}^{{commit}}")
    if proc.returncode != 0:
        raise SystemExit(f"unknown base ref: {base_ref}")


def resolve_base_ref(base_ref: str, default_ref: str = "origin/main"):
    if base_ref:
        validate_base_ref(base_ref)
        return base_ref
    branch = current_branch()
    candidates = [default_ref, "origin/trunk", "main", "master", "trunk"]
    if branch and branch not in candidates:
        candidates.append(branch)
    candidates.append("HEAD")
    for candidate in candidates:
        proc = run_git("rev-parse", "--verify", "--quiet", f"{candidate}^{{commit}}")
        if proc.returncode == 0:
            return candidate
    raise SystemExit(f"unknown base ref: {default_ref}")


def load_worktree_entries():
    proc = run_git("worktree", "list", "--porcelain")
    if proc.returncode != 0:
        raise SystemExit(proc.stderr or proc.stdout)
    return list(parse_porcelain_entries(proc.stdout))


def branch_name(entry):
    ref = entry.get("branch", "")
    if ref.startswith("refs/heads/"):
        return ref.removeprefix("refs/heads/")
    return "(detached)"


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


def parse_known_flags(args, allowed_flags):
    flags = {flag: False for flag in allowed_flags}
    positionals = []
    for arg in args:
        if arg in flags:
            if flags[arg]:
                raise SystemExit(f"duplicate option: {arg}")
            flags[arg] = True
        elif arg.startswith("--"):
            raise SystemExit(f"unknown option: {arg}")
        else:
            positionals.append(arg)
    return flags, positionals


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


def python_command():
    configured = os.environ.get("PYTHON")
    if configured:
        tokens = shlex.split(configured)
        if tokens:
            return tokens
    return [sys.executable or "python3"]


def run_script(name, *args):
    command = [*python_command(), script_path(name), *args]
    command_display = " ".join(display_text(str(part)) for part in (command[:-len(args)-1] if args else command[:-1]))
    try:
        proc = subprocess.run(
            command,
            text=True,
            encoding="utf-8",
            capture_output=True,
            env={**os.environ, **python_env},
        )
    except (OSError, UnicodeEncodeError) as err:
        raise SystemExit(f"unable to execute {name} via configured python command {command_display}: {display_text(str(err))}")
    if proc.returncode != 0:
        message = proc.stderr.strip() or proc.stdout.strip() or f"{name} failed with exit code {proc.returncode}"
        raise SystemExit(message)
    return proc.stdout


def load_json_script(name, *args):
    return json.loads(run_script(name, *args))


def load_worktree_audit_text(base_ref, *flags):
    return run_script("worktree_audit.py", *flags, base_ref).rstrip()


def load_worktree_audit_view(base_ref, *flags):
    return load_json_script("worktree_audit.py", "--json", *flags, base_ref)


def classify_branch(branch, base_ref):
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


def should_include_audit_row(state, branch, path, behind, current_worktree_path, *, merged_only, integrated_only, prune_only, stale_only):
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
