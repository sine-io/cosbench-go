#!/usr/bin/env python3

import json
import sys
from datetime import datetime
from pathlib import Path


def parse_time(value):
    if not value:
        return None
    try:
        return datetime.fromisoformat(value.replace("Z", "+00:00"))
    except ValueError:
        return None


def duration_seconds(started_at, updated_at):
    start = parse_time(started_at)
    end = parse_time(updated_at)
    if start is None or end is None:
        return None
    seconds = int((end - start).total_seconds())
    return seconds if seconds >= 0 else None


def age_seconds(generated_at, created_at):
    generated = parse_time(generated_at)
    created = parse_time(created_at)
    if generated is None or created is None:
        return None
    seconds = int((generated - created).total_seconds())
    return seconds if seconds >= 0 else None


def current_reason(valid, fresh, matches_head):
    if valid and fresh and matches_head:
        return "current"
    if not valid:
        return "not_successful"
    if not fresh:
        return "stale"
    if not matches_head:
        return "head_mismatch"
    return "missing"


def main(argv):
    if len(argv) != 4:
        raise SystemExit(
            "usage: finalize_smoke_ready_validate_payload.py <smoke-ready.json> <validation.json> <current-run.json>"
        )

    payload_path = Path(argv[1])
    validation_path = Path(argv[2])
    run_path = Path(argv[3])

    payload = json.loads(payload_path.read_text(encoding="utf-8"))
    validation = json.loads(validation_path.read_text(encoding="utf-8"))
    current_run = json.loads(run_path.read_text(encoding="utf-8"))

    generated_at = payload.get("generated_at", "")
    head_sha = current_run.get("headSha", "")
    head_branch = current_run.get("headBranch", "")
    created_at = current_run.get("createdAt", "") or generated_at
    duration = duration_seconds(current_run.get("startedAt", ""), current_run.get("updatedAt", ""))
    age = age_seconds(generated_at, created_at)
    valid = validation.get("valid") is True
    fresh = valid and age is not None and age <= payload["summary"]["freshness_thresholds_seconds"]["schema_validation"]
    matches_head = payload.get("current_head_sha", "") == head_sha
    current = valid and fresh and matches_head

    latest = {
        "database_id": current_run.get("databaseId"),
        "status": "completed",
        "conclusion": "success" if valid else "failure",
        "event": current_run.get("event", ""),
        "head_sha": head_sha,
        "head_branch": head_branch,
        "created_at": created_at,
        "url": current_run.get("url", ""),
    }
    payload.setdefault("workflows", {}).setdefault("latest", {})["Smoke Ready Validate"] = latest

    summary = payload.setdefault("summary", {})
    summary["schema_validation_latest_success"] = valid
    summary["schema_validation_latest_result"] = "validated" if valid else "failed"
    summary["schema_validation_latest_source"] = "Smoke Ready Validate"
    summary["schema_validation_latest_event"] = current_run.get("event", "")
    summary["schema_validation_latest_head_sha"] = head_sha
    summary["schema_validation_latest_head_branch"] = head_branch
    summary["schema_validation_latest_matches_head"] = matches_head
    summary["schema_validation_latest_run_id"] = current_run.get("databaseId")
    summary["schema_validation_latest_duration_seconds"] = duration
    summary["schema_validation_latest_age_seconds"] = age
    summary["schema_validation_latest_fresh"] = fresh
    summary["schema_validation_current"] = current
    summary["schema_validation_current_reason"] = current_reason(valid, fresh, matches_head)
    summary["schema_validation_latest_url"] = current_run.get("url", "")
    summary["schema_validation_latest_artifact"] = "smoke-ready-validate-summary"
    summary["schema_validation_latest_created_at"] = created_at

    payload_path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")


if __name__ == "__main__":
    main(sys.argv)
