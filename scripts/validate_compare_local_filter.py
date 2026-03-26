#!/usr/bin/env python3

import sys

from compare_local_manifest import FilterError, ManifestError, format_filter_error, read_manifest, validate_filter


def main() -> int:
    if len(sys.argv) < 3:
        raise SystemExit("usage: validate_compare_local_filter.py <manifest> <filter>")

    manifest_path = sys.argv[1]
    filter_args = sys.argv[2:]
    for arg in filter_args:
        if arg.startswith("--"):
            raise SystemExit(f"unknown option: {arg}")
    if len(filter_args) > 1:
        raise SystemExit(f"expected exactly one filter argument, got: {' '.join(filter_args)}")
    raw_filter = filter_args[0]
    try:
        fixtures = read_manifest(manifest_path)
    except ManifestError as err:
        raise SystemExit(str(err))
    try:
        validate_filter(fixtures, raw_filter)
    except FilterError as err:
        raise SystemExit(format_filter_error(fixtures, err))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
