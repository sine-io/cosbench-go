#!/usr/bin/env python3

import json
import sys

from compare_local_manifest import ManifestError, read_manifest, select_fixtures, validate_filter


def main() -> int:
    if len(sys.argv) not in (2, 3, 4):
        raise SystemExit("usage: list_compare_local_fixtures.py <manifest> [--names|--pairs] [filter]")

    try:
        fixtures = read_manifest(sys.argv[1])
    except ManifestError as err:
        raise SystemExit(str(err))
    names_only = False
    pairs_only = False
    filter_args = []
    for arg in sys.argv[2:]:
        if arg == "--names":
            names_only = True
        elif arg == "--pairs":
            pairs_only = True
        elif arg.startswith("--"):
            raise SystemExit(f"unknown option: {arg}")
        else:
            filter_args.append(arg)
    if names_only and pairs_only:
        raise SystemExit("choose only one of --names or --pairs")
    if len(filter_args) > 1:
        joined = " ".join(filter_args)
        raise SystemExit(f"expected at most one filter argument, got: {joined}")
    raw_filter = filter_args[0] if filter_args else ""

    try:
        validate_filter(fixtures, raw_filter)
    except ValueError as err:
        names = "".join(f"  - {fixture['name']}\n" for fixture in fixtures)
        raise SystemExit(f"unknown compare-local fixture: {err}\nknown fixtures:\n{names}")

    fixtures = select_fixtures(fixtures, raw_filter)

    if names_only:
        for fixture in fixtures:
            print(fixture["name"])
        return 0

    if pairs_only:
        for fixture in fixtures:
            print(f"{fixture['name']} {fixture['workload']}")
        return 0

    print(json.dumps(fixtures, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
