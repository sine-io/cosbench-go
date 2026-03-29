#!/usr/bin/env python3

import json
import sys
from pathlib import Path


def yes_no(value):
    return "yes" if value else "no"


def main(argv):
    if len(argv) != 3:
        raise SystemExit("usage: build_smoke_ready_validate_summary.py <src-dir> <dst-dir>")
    src = Path(argv[1])
    dst = Path(argv[2])
    validation_path = src / "validation.json"
    if not validation_path.exists():
        raise SystemExit(f"missing required validation file: {validation_path}")

    payload = json.loads(validation_path.read_text(encoding="utf-8"))
    dst.mkdir(parents=True, exist_ok=True)
    (dst / "summary.json").write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")

    lines = [
        "# Smoke Ready Validate",
        "",
        f"- Valid: `{yes_no(payload.get('valid') is True)}`",
        f"- Schema Path: `{payload.get('schema_path', '')}`",
        f"- Schema Version: `{payload.get('schema_version', '')}`",
        f"- Generated At: `{payload.get('generated_at', '')}`",
        f"- Error: `{payload.get('error', '')}`",
    ]
    (dst / "summary.md").write_text("\n".join(lines) + "\n", encoding="utf-8")


if __name__ == "__main__":
    main(sys.argv)
