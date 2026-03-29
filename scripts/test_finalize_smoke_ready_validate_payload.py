import json
import subprocess
from pathlib import Path


def test_finalize_smoke_ready_validate_payload_updates_schema_validation_to_current_run(tmp_path):
    payload_path = tmp_path / "smoke-ready.json"
    validation_path = tmp_path / "validation.json"
    run_path = tmp_path / "current-run.json"

    payload = {
        "generated_at": "2026-03-30T00:10:00Z",
        "current_head_sha": "sha-current",
        "workflows": {
            "latest": {
                "Smoke Ready Validate": {
                    "database_id": 100,
                    "status": "completed",
                    "conclusion": "success",
                    "event": "schedule",
                    "head_sha": "sha-old",
                    "head_branch": "main",
                    "created_at": "2026-03-29T00:00:00Z",
                    "url": "https://example.test/old",
                }
            }
        },
        "summary": {
            "schema_validation_latest_success": False,
            "schema_validation_latest_result": "failed",
            "schema_validation_latest_source": "Smoke Ready Validate",
            "schema_validation_latest_event": "schedule",
            "schema_validation_latest_head_sha": "sha-old",
            "schema_validation_latest_head_branch": "main",
            "schema_validation_latest_matches_head": False,
            "schema_validation_latest_run_id": 100,
            "schema_validation_latest_duration_seconds": 0,
            "schema_validation_latest_age_seconds": 999999,
            "schema_validation_latest_fresh": False,
            "schema_validation_current": False,
            "schema_validation_current_reason": "not_successful",
            "schema_validation_latest_url": "https://example.test/old",
            "schema_validation_latest_artifact": "smoke-ready-validate-summary",
            "schema_validation_latest_created_at": "2026-03-29T00:00:00Z",
            "freshness_thresholds_seconds": {
                "schema_validation": 172800,
                "remote": 172800,
                "real_endpoint": 2592000,
                "legacy_live": 2592000,
            },
        },
    }
    validation = {
        "valid": True,
        "schema_version": 1,
        "schema_path": "docs/smoke-ready.schema.json",
        "generated_at": "2026-03-30T00:10:00Z",
        "error": "",
    }
    current_run = {
        "databaseId": 1234,
        "event": "workflow_dispatch",
        "headSha": "sha-current",
        "headBranch": "main",
        "createdAt": "2026-03-30T00:09:00Z",
        "startedAt": "2026-03-30T00:09:10Z",
        "updatedAt": "2026-03-30T00:09:40Z",
        "url": "https://example.test/current",
    }

    payload_path.write_text(json.dumps(payload) + "\n", encoding="utf-8")
    validation_path.write_text(json.dumps(validation) + "\n", encoding="utf-8")
    run_path.write_text(json.dumps(current_run) + "\n", encoding="utf-8")

    subprocess.run(
        [
            "python3",
            "scripts/finalize_smoke_ready_validate_payload.py",
            str(payload_path),
            str(validation_path),
            str(run_path),
        ],
        check=True,
        cwd=Path.cwd(),
    )

    updated = json.loads(payload_path.read_text(encoding="utf-8"))
    summary = updated["summary"]
    latest = updated["workflows"]["latest"]["Smoke Ready Validate"]

    assert summary["schema_validation_latest_success"] is True
    assert summary["schema_validation_latest_result"] == "validated"
    assert summary["schema_validation_latest_event"] == "workflow_dispatch"
    assert summary["schema_validation_latest_head_sha"] == "sha-current"
    assert summary["schema_validation_latest_matches_head"] is True
    assert summary["schema_validation_latest_run_id"] == 1234
    assert summary["schema_validation_latest_duration_seconds"] == 30
    assert summary["schema_validation_latest_age_seconds"] == 60
    assert summary["schema_validation_latest_fresh"] is True
    assert summary["schema_validation_current"] is True
    assert summary["schema_validation_current_reason"] == "current"
    assert summary["schema_validation_latest_url"] == "https://example.test/current"
    assert summary["schema_validation_latest_artifact"] == "smoke-ready-validate-summary"
    assert summary["schema_validation_latest_created_at"] == "2026-03-30T00:09:00Z"

    assert latest["database_id"] == 1234
    assert latest["status"] == "completed"
    assert latest["conclusion"] == "success"
    assert latest["event"] == "workflow_dispatch"
    assert latest["head_sha"] == "sha-current"
    assert latest["head_branch"] == "main"
    assert latest["created_at"] == "2026-03-30T00:09:00Z"
    assert latest["url"] == "https://example.test/current"
