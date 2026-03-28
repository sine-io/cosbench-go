#!/usr/bin/env python3

import json
import os
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
ARTIFACT_DIR = ROOT / ".artifacts" / "remote-smoke"


def build_summary(
    *,
    controller_url,
    driver_urls,
    job_id,
    job_status,
    drivers_seen,
    units_claimed,
    drivers_participated,
    operation_count,
    byte_count,
    checks,
):
    overall = "pass" if all(value == "pass" for value in checks.values()) else "fail"
    return {
        "controller_url": controller_url,
        "driver_urls": driver_urls,
        "job_id": job_id,
        "job_status": job_status,
        "drivers_seen": drivers_seen,
        "units_claimed": units_claimed,
        "drivers_participated": drivers_participated,
        "operation_count": operation_count,
        "byte_count": byte_count,
        "checks": checks,
        "overall": overall,
    }


def build_failure_summary(failed_at, error):
    return {
        "failed_at": failed_at,
        "error": error,
        "overall": "fail",
    }


def render_summary_md(summary):
    lines = ["# Remote Smoke", ""]
    for key in [
        "controller_url",
        "driver_urls",
        "job_id",
        "job_status",
        "drivers_seen",
        "units_claimed",
        "drivers_participated",
        "operation_count",
        "byte_count",
        "overall",
    ]:
        if key in summary:
            lines.append(f"- `{key}`: `{summary[key]}`")
    if "checks" in summary:
        lines.append("")
        lines.append("## Checks")
        for name, status in summary["checks"].items():
            lines.append(f"- `{name}`: `{status}`")
    if "error" in summary:
        lines.append("")
        lines.append("## Error")
        lines.append(summary["error"])
    return "\n".join(lines) + "\n"


def write_summary(summary):
    ARTIFACT_DIR.mkdir(parents=True, exist_ok=True)
    json_path = ARTIFACT_DIR / "summary.json"
    md_path = ARTIFACT_DIR / "summary.md"
    json_path.write_text(json.dumps(summary, indent=2) + "\n", encoding="utf-8")
    md_path.write_text(render_summary_md(summary), encoding="utf-8")


def run_mock():
    mode = os.environ.get("SMOKE_REMOTE_LOCAL_MOCK", "").strip()
    if not mode:
        return None
    if mode == "success":
        summary = build_summary(
            controller_url="http://127.0.0.1:19088",
            driver_urls=["http://127.0.0.1:18081", "http://127.0.0.1:18082"],
            job_id="job-1",
            job_status="succeeded",
            drivers_seen=2,
            units_claimed=2,
            drivers_participated=2,
            operation_count=2,
            byte_count=2000,
            checks={
                "process_ready": "pass",
                "drivers_healthy": "pass",
                "units_distributed": "pass",
                "job_succeeded": "pass",
                "visibility": "pass",
            },
        )
        write_summary(summary)
        sys.stdout.write(render_summary_md(summary))
        return 0
    summary = build_failure_summary("controller", "mocked remote smoke failure")
    write_summary(summary)
    sys.stdout.write(render_summary_md(summary))
    return 1


def main():
    mock = run_mock()
    if mock is not None:
        raise SystemExit(mock)
    summary = build_failure_summary("orchestration", "real remote smoke helper not implemented yet")
    write_summary(summary)
    sys.stdout.write(render_summary_md(summary))
    raise SystemExit(1)


if __name__ == "__main__":
    main()
