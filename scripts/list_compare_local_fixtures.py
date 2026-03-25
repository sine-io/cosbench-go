#!/usr/bin/env python3

import json
import sys

from compare_local_manifest import read_manifest


def main() -> int:
    if len(sys.argv) not in (2, 3):
        raise SystemExit("usage: list_compare_local_fixtures.py <manifest> [--names]")

    fixtures = read_manifest(sys.argv[1])
    if len(sys.argv) == 3 and sys.argv[2] == "--names":
        for fixture in fixtures:
            print(fixture["name"])
        return 0

    print(json.dumps(fixtures, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
