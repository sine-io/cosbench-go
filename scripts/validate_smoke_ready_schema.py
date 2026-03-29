#!/usr/bin/env python3

import argparse
import json
import subprocess
import sys
from pathlib import Path

import jsonschema


SCHEMA_PATH = Path("docs/smoke-ready.schema.json")


def run_smoke_ready_json():
    proc = subprocess.run(
        [sys.executable, "scripts/smoke_ready.py", "--json"],
        text=True,
        capture_output=True,
        check=False,
    )
    if proc.returncode != 0:
        error = (proc.stderr or proc.stdout).strip()
        raise RuntimeError(f"smoke-ready-json failed: {error}")
    try:
        return json.loads(proc.stdout)
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"smoke-ready-json emitted invalid JSON: {exc}") from exc


def load_schema():
    with SCHEMA_PATH.open("r", encoding="utf-8") as f:
        return json.load(f)


def build_result():
    try:
        payload = run_smoke_ready_json()
        schema = load_schema()
        jsonschema.validate(payload, schema)
        return {
            "schema_path": str(SCHEMA_PATH),
            "schema_version": payload.get("schema_version"),
            "repo": payload.get("repo", ""),
            "generated_at": payload.get("generated_at", ""),
            "valid": True,
            "error": "",
        }, 0
    except Exception as exc:  # pragma: no cover - exercised by exit status paths
        return {
            "schema_path": str(SCHEMA_PATH),
            "schema_version": None,
            "repo": "",
            "generated_at": "",
            "valid": False,
            "error": str(exc),
        }, 1


def print_text(result):
    print("# Smoke Ready Schema Validation")
    print()
    print(f"Schema: `{result['schema_path']}`")
    print(f"Schema Version: `{result['schema_version']}`")
    print(f"Repository: `{result['repo']}`")
    print(f"Generated At: `{result['generated_at']}`")
    print(f"Valid: `{'yes' if result['valid'] else 'no'}`")
    if result["error"]:
        print(f"Error: `{result['error']}`")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--json", action="store_true", help="emit machine-readable validation result")
    args = parser.parse_args()

    result, exit_code = build_result()
    if args.json:
        print(json.dumps(result, indent=2))
    else:
        print_text(result)
    raise SystemExit(exit_code)


if __name__ == "__main__":
    main()
