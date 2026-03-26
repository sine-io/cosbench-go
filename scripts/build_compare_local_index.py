#!/usr/bin/env python3

import json
import sys
from datetime import datetime, timezone
from pathlib import Path

from compare_local_manifest import ManifestFormatError, normalize_filter, read_manifest, select_fixtures


def build_summary(payload):
    meta = payload["meta"]
    lines = [
        "## Compare Local",
        "",
        "Artifact directory: `.artifacts/compare-local/`",
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


def main() -> int:
    if len(sys.argv) not in (3, 4):
        raise SystemExit("usage: build_compare_local_index.py <manifest> <output_dir> [filter]")

    output_dir = Path(sys.argv[2])
    selected = sys.argv[3] if len(sys.argv) == 4 else ""
    filter_label = normalize_filter(selected)
    fixtures = []

    try:
        manifest_fixtures = read_manifest(sys.argv[1])
    except ManifestFormatError as err:
        raise SystemExit(str(err))

    for fixture in select_fixtures(manifest_fixtures, selected):
        name = fixture["name"]
        workload = fixture["workload"]
        summary_name = f"{name}.json"
        summary = json.loads((output_dir / summary_name).read_text())
        fixtures.append(
            {
                "name": name,
                "workload": workload,
                "summary": summary_name,
                "stages": summary["stages"],
                "works": summary["works"],
                "samples": summary["samples"],
                "errors": summary["errors"],
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
    (output_dir / "index.json").write_text(json.dumps(payload, indent=2) + "\n")
    (output_dir / "summary.md").write_text(build_summary(payload))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
