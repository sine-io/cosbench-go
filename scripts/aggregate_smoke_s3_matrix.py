#!/usr/bin/env python3

import json
import sys
from pathlib import Path


EXPECTED_ROWS = ["s3", "sio"]


def find_output_path(root: Path, backend: str):
    candidates = [
        root / f"smoke-s3-{backend}" / "smoke-s3-output.txt",
        root / f"smoke-s3-{backend}" / backend / "smoke-s3-output.txt",
    ]
    for candidate in candidates:
        if candidate.exists():
            return candidate
    return None


def find_summary_path(root: Path, backend: str):
    candidates = [
        root / f"smoke-s3-{backend}" / "summary.json",
        root / f"smoke-s3-{backend}" / backend / "summary.json",
        root / f"smoke-s3-{backend}" / ".artifacts" / "smoke-s3-summary" / "summary.json",
        root / f"smoke-s3-{backend}" / backend / ".artifacts" / "smoke-s3-summary" / "summary.json",
    ]
    for candidate in candidates:
        if candidate.exists():
            return candidate
    return None


def aggregate_rows(root: Path, expected_rows=EXPECTED_ROWS):
    rows = []
    for backend in expected_rows:
        output_path = find_output_path(root, backend)
        summary_path = find_summary_path(root, backend)
        if output_path is None:
            rows.append({"backend": backend, "status": "missing"})
            continue
        status = "present"
        if summary_path is not None:
            payload = json.loads(summary_path.read_text(encoding="utf-8"))
            status = payload.get("result", "present")
        rows.append(
            {
                "backend": backend,
                "status": status,
                "output": output_path.read_text(encoding="utf-8"),
            }
        )
    return rows


def render_markdown(rows):
    lines = [
        "# Smoke S3 Matrix Summary",
        "",
        "| backend | status |",
        "| --- | --- |",
    ]
    for row in rows:
        lines.append(f"| {row['backend']} | {row['status']} |")
    return "\n".join(lines) + "\n"


def build_payload(rows):
    return {
        "rows": rows,
        "overall": "pass" if all(row["status"] in {"executed", "skipped"} for row in rows) else "partial",
    }


def main(argv):
    if len(argv) != 3:
        raise SystemExit("usage: aggregate_smoke_s3_matrix.py <download-dir> <output-dir>")
    download_dir = Path(argv[1])
    output_dir = Path(argv[2])
    rows = aggregate_rows(download_dir)
    if not any(row["status"] == "present" for row in rows):
        raise SystemExit("no smoke-s3 outputs found")
    output_dir.mkdir(parents=True, exist_ok=True)
    payload = build_payload(rows)
    (output_dir / "summary.json").write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
    markdown = render_markdown(rows)
    (output_dir / "summary.md").write_text(markdown, encoding="utf-8")
    sys.stdout.write(markdown)


if __name__ == "__main__":
    main(sys.argv)
