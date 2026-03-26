from datetime import datetime, timezone
import json
from pathlib import Path
import subprocess


def generated_at():
    return datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")


def current_worktree():
    proc = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        check=False,
        text=True,
        capture_output=True,
    )
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
