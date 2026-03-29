import json
import os
import subprocess
from pathlib import Path

import jsonschema


def run_helper_json():
    env = os.environ.copy()
    env["SMOKE_READY_MOCK_REPO_SECRETS"] = "COSBENCH_SMOKE_ENDPOINT,COSBENCH_SMOKE_ACCESS_KEY,COSBENCH_SMOKE_SECRET_KEY"
    env["SMOKE_READY_MOCK_WORKFLOWS"] = ",".join(
        [
            "Smoke Local",
            "Smoke S3",
            "Smoke S3 Matrix",
            "Legacy Live Compare",
            "Legacy Live Compare Matrix",
            "Remote Smoke Local",
            "Remote Smoke Matrix",
            "Remote Smoke Recovery",
            "Remote Smoke Recovery Matrix",
        ]
    )
    env["SMOKE_READY_MOCK_WORKFLOW_RUNS_JSON"] = json.dumps(
        {
            "Smoke Local": {"databaseId": 1, "status": "completed", "conclusion": "success", "created_at": "2026-03-29T00:00:00Z", "url": "https://example.test/smoke-local"},
            "Smoke S3": {"databaseId": 2, "status": "completed", "conclusion": "success", "created_at": "2026-03-29T00:05:00Z", "url": "https://example.test/smoke-s3"},
            "Smoke S3 Matrix": {"databaseId": 3, "status": "completed", "conclusion": "success", "created_at": "2026-03-29T00:06:00Z", "url": "https://example.test/smoke-s3-matrix"},
            "Legacy Live Compare": {"databaseId": 4, "status": "completed", "conclusion": "success", "created_at": "2026-03-29T00:07:00Z", "url": "https://example.test/legacy-live-compare"},
            "Legacy Live Compare Matrix": {"databaseId": 5, "status": "completed", "conclusion": "success", "created_at": "2026-03-29T00:08:00Z", "url": "https://example.test/legacy-live-compare-matrix"},
            "Remote Smoke Local": {"databaseId": 6, "status": "completed", "conclusion": "success", "created_at": "2026-03-29T00:10:00Z", "url": "https://example.test/remote-smoke-local"},
            "Remote Smoke Matrix": {"databaseId": 7, "status": "completed", "conclusion": "success", "created_at": "2026-03-29T00:20:00Z", "url": "https://example.test/remote-smoke-matrix"},
            "Remote Smoke Recovery": {"databaseId": 8, "status": "completed", "conclusion": "success", "created_at": "2026-03-29T00:30:00Z", "url": "https://example.test/remote-smoke-recovery"},
            "Remote Smoke Recovery Matrix": {"databaseId": 9, "status": "completed", "conclusion": "success", "created_at": "2026-03-29T00:40:00Z", "url": "https://example.test/remote-smoke-recovery-matrix"},
        }
    )
    env["SMOKE_READY_MOCK_WORKFLOW_DETAILS_JSON"] = json.dumps(
        {
            "Smoke S3": {"summary": {"result": "skipped"}, "output": ""},
            "Smoke S3 Matrix": {"rows": [{"backend": "s3", "status": "skipped"}, {"backend": "sio", "status": "skipped"}]},
            "Remote Smoke Local": {"summary": {"overall": "pass", "job_status": "succeeded"}},
            "Remote Smoke Matrix": {"overall": "pass", "rows": [{"backend": "s3", "scenario": "single", "status": "present", "summary": {"overall": "pass"}}]},
            "Remote Smoke Recovery": {"summary": {"overall": "pass", "job_status": "succeeded"}},
            "Remote Smoke Recovery Matrix": {"overall": "pass", "rows": [{"backend": "s3", "scenario": "recovery", "status": "present", "summary": {"overall": "pass"}}]},
            "Legacy Live Compare": {"result": {"result": "skipped", "fixture": "testdata/legacy/sio-config-sample.xml", "backend": "sio", "reason": "missing secrets"}},
            "Legacy Live Compare Matrix": {"rows": [{"backend": "s3", "status": "skipped"}, {"backend": "sio", "status": "skipped"}]},
        }
    )
    proc = subprocess.run(
        ["python3", "scripts/smoke_ready.py", "--json"],
        cwd=os.getcwd(),
        env=env,
        text=True,
        capture_output=True,
        check=True,
    )
    return json.loads(proc.stdout)


def load_schema():
    with Path("docs/smoke-ready.schema.json").open("r", encoding="utf-8") as f:
        return json.load(f)


def test_smoke_ready_schema_contract():
    payload = run_helper_json()
    schema = load_schema()
    jsonschema.validate(payload, schema)
    assert payload["schema_version"] == 1
    for key in ["repo", "required", "local_env", "repo_secrets", "workflows", "summary", "blockers"]:
        assert key in payload

    summary = payload["summary"]
    required_summary_keys = [
        "real_endpoint_latest_result",
        "real_endpoint_latest_source",
        "real_endpoint_latest_url",
        "real_endpoint_latest_artifact",
        "real_endpoint_latest_created_at",
        "real_endpoint_matrix_latest_result",
        "real_endpoint_matrix_latest_source",
        "real_endpoint_matrix_latest_url",
        "real_endpoint_matrix_latest_artifact",
        "real_endpoint_matrix_latest_created_at",
        "legacy_live_latest_result",
        "legacy_live_latest_source",
        "legacy_live_latest_url",
        "legacy_live_latest_artifact",
        "legacy_live_latest_created_at",
        "legacy_live_matrix_latest_result",
        "legacy_live_matrix_latest_source",
        "legacy_live_matrix_latest_url",
        "legacy_live_matrix_latest_artifact",
        "legacy_live_matrix_latest_created_at",
        "remote_happy_latest_result",
        "remote_happy_latest_source",
        "remote_happy_latest_url",
        "remote_happy_latest_artifact",
        "remote_happy_latest_created_at",
        "remote_recovery_latest_result",
        "remote_recovery_latest_source",
        "remote_recovery_latest_url",
        "remote_recovery_latest_artifact",
        "remote_recovery_latest_created_at",
    ]
    for key in required_summary_keys:
        assert key in summary
