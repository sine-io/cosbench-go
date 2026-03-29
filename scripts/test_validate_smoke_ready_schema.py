import json
import os
import subprocess


def build_env():
    env = os.environ.copy()
    env["SMOKE_READY_MOCK_REPO_SECRETS"] = "COSBENCH_SMOKE_ENDPOINT,COSBENCH_SMOKE_ACCESS_KEY,COSBENCH_SMOKE_SECRET_KEY"
    env["SMOKE_READY_MOCK_WORKFLOWS"] = ",".join(
        [
            "Smoke Local",
            "Smoke S3",
            "Smoke S3 Matrix",
            "Smoke Ready Validate",
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
            "Smoke Ready Validate": {"databaseId": 10, "status": "completed", "conclusion": "success", "created_at": "2026-03-29T00:06:30Z", "url": "https://example.test/smoke-ready-validate"},
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
            "Smoke Ready Validate": {"valid": True, "schema_path": "docs/smoke-ready.schema.json", "schema_version": 1, "repo": "sine-io/cosbench-go", "generated_at": "2026-03-29T00:06:31Z", "error": ""},
            "Remote Smoke Local": {"summary": {"overall": "pass", "job_status": "succeeded"}},
            "Remote Smoke Matrix": {"overall": "pass", "rows": [{"backend": "s3", "scenario": "single", "status": "present", "summary": {"overall": "pass"}}]},
            "Remote Smoke Recovery": {"summary": {"overall": "pass", "job_status": "succeeded"}},
            "Remote Smoke Recovery Matrix": {"overall": "pass", "rows": [{"backend": "s3", "scenario": "recovery", "status": "present", "summary": {"overall": "pass"}}]},
            "Legacy Live Compare": {"result": {"result": "skipped", "fixture": "testdata/legacy/sio-config-sample.xml", "backend": "sio", "reason": "missing secrets"}},
            "Legacy Live Compare Matrix": {"rows": [{"backend": "s3", "status": "skipped"}, {"backend": "sio", "status": "skipped"}]},
        }
    )
    return env


def test_validate_smoke_ready_schema_json_reports_valid():
    proc = subprocess.run(
        ["python3", "scripts/validate_smoke_ready_schema.py", "--json"],
        cwd=os.getcwd(),
        env=build_env(),
        text=True,
        capture_output=True,
    )
    assert proc.returncode == 0
    payload = json.loads(proc.stdout)
    assert payload["valid"] is True
    assert payload["schema_path"] == "docs/smoke-ready.schema.json"
    assert payload["schema_version"] == 1


def test_validate_smoke_ready_schema_text_reports_valid():
    proc = subprocess.run(
        ["python3", "scripts/validate_smoke_ready_schema.py"],
        cwd=os.getcwd(),
        env=build_env(),
        text=True,
        capture_output=True,
    )
    assert proc.returncode == 0
    assert "Smoke Ready Schema Validation" in proc.stdout
    assert "Valid: `yes`" in proc.stdout
