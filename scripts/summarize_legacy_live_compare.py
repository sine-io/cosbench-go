#!/usr/bin/env python3

import json
import sys
from pathlib import Path


def classify_summary(payload: dict) -> str:
    if payload.get("status") == "skipped":
        return "skipped"
    if payload.get("workload"):
        return "executed"
    return "failed"


def build_payload(raw_payload: dict, fixture: str, backend: str) -> dict:
    payload = {
        "result": classify_summary(raw_payload),
        "fixture": fixture,
        "backend": backend,
    }
    if payload["result"] == "skipped" and "reason" in raw_payload:
        payload["reason"] = raw_payload["reason"]
    return payload


def main(argv):
    if len(argv) != 5:
        raise SystemExit("usage: summarize_legacy_live_compare.py <input> <output> <fixture> <backend>")
    input_path = Path(argv[1])
    output_path = Path(argv[2])
    fixture = argv[3]
    backend = argv[4]
    raw_payload = json.loads(input_path.read_text(encoding="utf-8"))
    payload = build_payload(raw_payload, fixture=fixture, backend=backend)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
    sys.stdout.write(json.dumps(payload, indent=2) + "\n")


if __name__ == "__main__":
    main(sys.argv)
