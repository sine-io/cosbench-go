#!/usr/bin/env python3

import json
import sys

from compare_local_manifest import parse_filter, read_manifest, validate_filter


def main() -> int:
    if len(sys.argv) not in (2, 3, 4):
        raise SystemExit("usage: list_compare_local_fixtures.py <manifest> [--names] [filter]")

    fixtures = read_manifest(sys.argv[1])
    names_only = False
    raw_filter = ""
    for arg in sys.argv[2:]:
        if arg == "--names":
            names_only = True
        else:
            raw_filter = arg

    try:
        validate_filter(fixtures, raw_filter)
    except ValueError as err:
        names = "".join(f"  - {fixture['name']}\n" for fixture in fixtures)
        raise SystemExit(f"unknown compare-local fixture: {err}\nknown fixtures:\n{names}")

    selected = set(parse_filter(raw_filter))
    if selected:
        fixtures = [fixture for fixture in fixtures if fixture["name"] in selected]

    if names_only:
        for fixture in fixtures:
            print(fixture["name"])
        return 0

    print(json.dumps(fixtures, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
