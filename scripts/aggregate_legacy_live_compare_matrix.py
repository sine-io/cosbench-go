#!/usr/bin/env python3

import json
import sys
from pathlib import Path


EXPECTED_ROWS = ["s3", "sio"]


def find_summary_path(root: Path, backend: str):
    candidates = [
        root / f"legacy-live-compare-{backend}" / "summary.json",
        root / f"legacy-live-compare-{backend}" / backend / "summary.json",
    ]
    for candidate in candidates:
        if candidate.exists():
            return candidate
    return None


def classify_row(backend: str, summary_path: Path | None):
    if summary_path is None:
        return {"backend": backend, "status": "missing"}
    payload = json.loads(summary_path.read_text(encoding="utf-8"))
    if payload.get("status") == "skipped":
        return {
            "backend": backend,
            "status": "skipped",
            "reason": payload.get("reason", ""),
        }
    return {
        "backend": backend,
        "status": "executed",
    }


def aggregate_rows(root: Path, expected_rows=EXPECTED_ROWS):
    return [classify_row(backend, find_summary_path(root, backend)) for backend in expected_rows]


def build_payload(rows):
    return {
        "rows": rows,
        "overall": "pass" if all(row["status"] in {"executed", "skipped"} for row in rows) else "partial",
    }


def render_markdown(rows):
    lines = [
        "# Legacy Live Compare Matrix Summary",
        "",
        "| backend | status |",
        "| --- | --- |",
    ]
    for row in rows:
        lines.append(f"| {row['backend']} | {row['status']} |")
    return "\n".join(lines) + "\n"


def main(argv):
    if len(argv) != 3:
        raise SystemExit("usage: aggregate_legacy_live_compare_matrix.py <download-dir> <output-dir>")
    download_dir = Path(argv[1])
    output_dir = Path(argv[2])
    rows = aggregate_rows(download_dir)
    if not any(row["status"] in {"executed", "skipped"} for row in rows):
        raise SystemExit("no legacy live compare outputs found")
    output_dir.mkdir(parents=True, exist_ok=True)
    payload = build_payload(rows)
    (output_dir / "summary.json").write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
    markdown = render_markdown(rows)
    (output_dir / "summary.md").write_text(markdown, encoding="utf-8")
    sys.stdout.write(markdown)


if __name__ == "__main__":
    main(sys.argv)
