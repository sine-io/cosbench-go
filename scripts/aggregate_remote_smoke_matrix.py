#!/usr/bin/env python3

import json
import sys
from pathlib import Path


EXPECTED_ROWS = [
    ("s3", "single"),
    ("s3", "multistage"),
    ("sio", "single"),
    ("sio", "multistage"),
]


def find_summary_path(root: Path, backend: str, scenario: str):
    artifact_dir = root / f"remote-smoke-{backend}-{scenario}"
    direct = artifact_dir / "summary.json"
    nested = artifact_dir / "remote-smoke" / "summary.json"
    if direct.exists():
        return direct
    if nested.exists():
        return nested
    return None


def aggregate_rows(root: Path, expected_rows=EXPECTED_ROWS):
    rows = []
    for backend, scenario in expected_rows:
        summary_path = find_summary_path(root, backend, scenario)
        if summary_path is None:
            rows.append(
                {
                    "backend": backend,
                    "scenario": scenario,
                    "status": "missing",
                }
            )
            continue
        rows.append(
            {
                "backend": backend,
                "scenario": scenario,
                "status": "present",
                "summary": json.loads(summary_path.read_text(encoding="utf-8")),
            }
        )
    return rows


def render_markdown(rows):
    lines = [
        "# Remote Smoke Matrix Summary",
        "",
        "| backend | scenario | status | overall | job_status | drivers_seen | units_claimed | stages_seen |",
        "| --- | --- | --- | --- | --- | --- | --- | --- |",
    ]
    for row in rows:
        if row["status"] != "present":
            lines.append(
                f"| {row['backend']} | {row['scenario']} | missing | missing | missing | missing | missing | missing |"
            )
            continue
        summary = row["summary"]
        lines.append(
            "| {backend} | {scenario} | present | {overall} | {job_status} | {drivers_seen} | {units_claimed} | {stages_seen} |".format(
                backend=row["backend"],
                scenario=row["scenario"],
                overall=summary.get("overall", "missing"),
                job_status=summary.get("job_status", "missing"),
                drivers_seen=summary.get("drivers_seen", "missing"),
                units_claimed=summary.get("units_claimed", "missing"),
                stages_seen=summary.get("stages_seen", "missing"),
            )
        )
    return "\n".join(lines) + "\n"


def build_payload(rows):
    return {
        "rows": rows,
        "overall": "pass" if all(row["status"] == "present" and row["summary"].get("overall") == "pass" for row in rows) else "partial",
    }


def main(argv):
    if len(argv) != 3:
        raise SystemExit("usage: aggregate_remote_smoke_matrix.py <download-dir> <output-dir>")
    download_dir = Path(argv[1])
    output_dir = Path(argv[2])
    rows = aggregate_rows(download_dir)
    present = [row for row in rows if row["status"] == "present"]
    if not present:
        raise SystemExit("no remote smoke summaries found")
    output_dir.mkdir(parents=True, exist_ok=True)
    payload = build_payload(rows)
    (output_dir / "summary.json").write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
    markdown = render_markdown(rows)
    (output_dir / "summary.md").write_text(markdown, encoding="utf-8")
    sys.stdout.write(markdown)


if __name__ == "__main__":
    main(sys.argv)
