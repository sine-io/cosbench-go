#!/usr/bin/env python3

import json
import sys
from datetime import datetime, timezone
from pathlib import Path

from compare_local_manifest import ManifestError, normalize_filter, read_manifest, select_fixtures, validate_filter


def build_summary(payload, output_dir: Path):
    meta = payload["meta"]
    lines = [
        "## Compare Local",
        "",
        f"Artifact directory: `{output_dir}`",
        "",
        f"Filter: `{meta.get('filter', 'all')}`",
        "",
        f"Fixture count: {meta.get('fixture_count', 0)}",
        "",
        f"Generated at: `{meta.get('generated_at', '')}`",
        "",
        "| Fixture | Workload | Stages | Works | Samples | Errors | Summary |",
        "| --- | --- | --- | --- | --- | --- | --- |",
    ]
    for fixture in payload.get("fixtures", []):
        lines.append(
            f"| `{fixture['name']}` | `{fixture['workload']}` | {fixture['stages']} | {fixture['works']} | {fixture['samples']} | {fixture['errors']} | `{fixture['summary']}` |"
        )
    return "\n".join(lines) + "\n"


def write_output_file(path: Path, content: str):
    try:
        path.write_text(content)
    except OSError as err:
        raise SystemExit(f"unable to write compare-local artifact {path}: {err}")


def load_fixture_summary(output_dir: Path, summary_name: str, fixture_name: str):
    summary_path = output_dir / summary_name
    try:
        summary = json.loads(summary_path.read_text())
    except FileNotFoundError:
        raise SystemExit(f"missing compare-local summary for fixture {fixture_name}: {summary_path}")
    except OSError as err:
        raise SystemExit(f"unable to read compare-local summary for fixture {fixture_name}: {summary_path}: {err}")
    except json.JSONDecodeError as err:
        raise SystemExit(f"invalid compare-local summary for fixture {fixture_name}: {summary_path}: {err}")
    if not isinstance(summary, dict):
        raise SystemExit(
            f"invalid compare-local summary for fixture {fixture_name}: {summary_path}: summary payload must be a JSON object"
        )
    return summary


def require_summary_field(summary, field: str, fixture_name: str, summary_path: Path):
    if field not in summary:
        raise SystemExit(
            f"invalid compare-local summary for fixture {fixture_name}: {summary_path}: missing required field {field}"
        )
    return summary[field]


def require_summary_int(summary, field: str, fixture_name: str, summary_path: Path):
    value = require_summary_field(summary, field, fixture_name, summary_path)
    if isinstance(value, bool) or not isinstance(value, int):
        raise SystemExit(
            f"invalid compare-local summary for fixture {fixture_name}: {summary_path}: field {field} must be an integer"
        )
    return value


def main() -> int:
    if len(sys.argv) not in (3, 4):
        raise SystemExit("usage: build_compare_local_index.py <manifest> <output_dir> [filter]")

    output_dir = Path(sys.argv[2])
    selected = sys.argv[3] if len(sys.argv) == 4 else ""
    if selected.startswith("--"):
        raise SystemExit(f"unknown option: {selected}")
    filter_label = normalize_filter(selected)
    fixtures = []

    try:
        manifest_fixtures = read_manifest(sys.argv[1])
    except ManifestError as err:
        raise SystemExit(str(err))
    try:
        validate_filter(manifest_fixtures, selected)
    except ValueError as err:
        names = "".join(f"  - {fixture['name']}\n" for fixture in manifest_fixtures)
        raise SystemExit(f"unknown compare-local fixture: {err}\nknown fixtures:\n{names}")

    for fixture in select_fixtures(manifest_fixtures, selected):
        name = fixture["name"]
        workload = fixture["workload"]
        summary_name = f"{name}.json"
        summary_path = output_dir / summary_name
        summary = load_fixture_summary(output_dir, summary_name, name)
        fixtures.append(
            {
                "name": name,
                "workload": workload,
                "summary": summary_name,
                "stages": require_summary_int(summary, "stages", name, summary_path),
                "works": require_summary_int(summary, "works", name, summary_path),
                "samples": require_summary_int(summary, "samples", name, summary_path),
                "errors": require_summary_int(summary, "errors", name, summary_path),
            }
        )

    payload = {
        "meta": {
            "filter": filter_label,
            "fixture_count": len(fixtures),
            "generated_at": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
        },
        "fixtures": fixtures,
    }
    output_dir.mkdir(parents=True, exist_ok=True)
    write_output_file(output_dir / "index.json", json.dumps(payload, indent=2) + "\n")
    write_output_file(output_dir / "summary.md", build_summary(payload, output_dir))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
