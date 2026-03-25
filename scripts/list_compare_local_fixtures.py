#!/usr/bin/env python3

import json
import sys
from pathlib import Path


def main() -> int:
    if len(sys.argv) != 2:
        raise SystemExit("usage: list_compare_local_fixtures.py <manifest>")

    manifest_path = Path(sys.argv[1])
    fixtures = []

    for raw_line in manifest_path.read_text().splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        name, workload = line.split()
        fixtures.append({"name": name, "workload": workload})

    print(json.dumps(fixtures, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
