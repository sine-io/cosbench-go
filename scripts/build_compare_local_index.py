#!/usr/bin/env python3

import json
import sys
from pathlib import Path


def main() -> int:
    if len(sys.argv) not in (3, 4):
        raise SystemExit("usage: build_compare_local_index.py <manifest> <output_dir> [filter]")

    manifest_path = Path(sys.argv[1])
    output_dir = Path(sys.argv[2])
    selected = sys.argv[3] if len(sys.argv) == 4 else ""
    fixtures = []

    for raw_line in manifest_path.read_text().splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        name, workload = line.split()
        if selected and name != selected:
            continue
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
            "filter": selected,
            "fixture_count": len(fixtures),
        },
        "fixtures": fixtures,
    }
    (output_dir / "index.json").write_text(json.dumps(payload, indent=2) + "\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
