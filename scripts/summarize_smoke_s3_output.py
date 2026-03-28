#!/usr/bin/env python3

import json
import sys
from pathlib import Path


def classify_output(text: str) -> str:
    if "--- FAIL:" in text or "\nFAIL\n" in text:
        return "failed"
    if "--- SKIP:" in text and "--- PASS:" not in text:
        return "skipped"
    if "--- PASS:" in text and "\nPASS\n" in text:
        return "executed"
    return "failed"


def build_payload(text: str) -> dict:
    return {
        "result": classify_output(text),
    }


def main(argv):
    if len(argv) != 3:
        raise SystemExit("usage: summarize_smoke_s3_output.py <input> <output>")
    input_path = Path(argv[1])
    output_path = Path(argv[2])
    text = input_path.read_text(encoding="utf-8")
    payload = build_payload(text)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
    sys.stdout.write(json.dumps(payload, indent=2) + "\n")


if __name__ == "__main__":
    main(sys.argv)
