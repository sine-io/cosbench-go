#!/usr/bin/env python3

import shutil
import sys
from pathlib import Path


def main(argv):
    if len(argv) != 3:
        raise SystemExit("usage: build_remote_smoke_recovery_summary.py <src-dir> <dst-dir>")
    src = Path(argv[1])
    dst = Path(argv[2])
    dst.mkdir(parents=True, exist_ok=True)
    for name in ("summary.json", "summary.md"):
        source = src / name
        if not source.exists():
            raise SystemExit(f"missing required summary file: {source}")
        shutil.copy2(source, dst / name)


if __name__ == "__main__":
    main(sys.argv)
