import json
import os
import subprocess


def run_helper(*args, env_overrides=None):
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
            "Smoke Local": {
                "databaseId": 1001,
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:00:00Z",
                "url": "https://example.test/smoke-local",
            },
            "Smoke S3": {
                "databaseId": 1002,
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:05:00Z",
                "url": "https://example.test/smoke-s3",
            },
            "Smoke S3 Matrix": {
                "databaseId": 1009,
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:06:00Z",
                "url": "https://example.test/smoke-s3-matrix",
            },
            "Legacy Live Compare": {
                "databaseId": 1003,
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:07:00Z",
                "url": "https://example.test/legacy-live-compare",
            },
            "Legacy Live Compare Matrix": {
                "databaseId": 1004,
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:08:00Z",
                "url": "https://example.test/legacy-live-compare-matrix",
            },
            "Remote Smoke Local": {
                "databaseId": 1005,
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:10:00Z",
                "url": "https://example.test/remote-smoke-local",
            },
            "Remote Smoke Matrix": {
                "databaseId": 1006,
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:20:00Z",
                "url": "https://example.test/remote-smoke-matrix",
            },
            "Remote Smoke Recovery": {
                "databaseId": 1007,
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:30:00Z",
                "url": "https://example.test/remote-smoke-recovery",
            },
            "Remote Smoke Recovery Matrix": {
                "databaseId": 1008,
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:40:00Z",
                "url": "https://example.test/remote-smoke-recovery-matrix",
            },
        }
    )
    env["SMOKE_READY_MOCK_WORKFLOW_DETAILS_JSON"] = json.dumps(
        {
            "Smoke S3": {
                "summary": {"result": "skipped"},
                "output": "",
            },
            "Smoke S3 Matrix": {
                "rows": [
                    {
                        "backend": "s3",
                        "status": "skipped",
                        "output": "",
                    },
                    {
                        "backend": "sio",
                        "status": "skipped",
                        "output": "",
                    },
                ]
            },
            "Remote Smoke Local": {
                "summary": {
                    "overall": "pass",
                    "job_status": "succeeded",
                }
            },
            "Remote Smoke Matrix": {
                "rows": [
                    {"backend": "s3", "scenario": "single", "status": "present", "summary": {"overall": "pass"}},
                    {"backend": "s3", "scenario": "multistage", "status": "present", "summary": {"overall": "pass"}},
                    {"backend": "sio", "scenario": "single", "status": "present", "summary": {"overall": "pass"}},
                    {"backend": "sio", "scenario": "multistage", "status": "present", "summary": {"overall": "pass"}},
                ],
                "overall": "pass",
            },
            "Remote Smoke Recovery": {
                "summary": {
                    "overall": "pass",
                    "job_status": "succeeded",
                }
            },
            "Remote Smoke Recovery Matrix": {
                "rows": [
                    {"backend": "s3", "scenario": "recovery", "status": "present", "summary": {"overall": "pass"}},
                    {"backend": "sio", "scenario": "recovery", "status": "present", "summary": {"overall": "pass"}},
                ],
                "overall": "pass",
            },
            "Legacy Live Compare": {
                "result": {
                    "result": "skipped",
                    "fixture": "testdata/legacy/sio-config-sample.xml",
                    "backend": "sio",
                    "reason": "missing secrets",
                }
            },
            "Legacy Live Compare Matrix": {
                "rows": [
                    {"backend": "s3", "status": "skipped"},
                    {"backend": "sio", "status": "skipped"},
                ]
            },
        }
    )
    if env_overrides:
        env.update(env_overrides)
    return subprocess.run(
        ["python3", "scripts/smoke_ready.py", *args],
        cwd=os.getcwd(),
        env=env,
        text=True,
        capture_output=True,
        check=True,
    )


def test_smoke_ready_json_reports_full_workflow_surface():
    proc = run_helper("--json")
    payload = json.loads(proc.stdout)
    present = payload["workflows"]["present"]
    assert present["Smoke Local"] is True
    assert present["Smoke S3"] is True
    assert present["Smoke S3 Matrix"] is True
    assert present["Legacy Live Compare"] is True
    assert present["Legacy Live Compare Matrix"] is True
    assert present["Remote Smoke Local"] is True
    assert present["Remote Smoke Matrix"] is True
    assert present["Remote Smoke Recovery"] is True
    assert present["Remote Smoke Recovery Matrix"] is True
    latest = payload["workflows"]["latest"]
    assert latest["Smoke Local"]["conclusion"] == "success"
    assert latest["Smoke S3"]["created_at"] == "2026-03-29T00:05:00Z"
    assert latest["Smoke S3 Matrix"]["url"] == "https://example.test/smoke-s3-matrix"
    assert latest["Legacy Live Compare"]["url"] == "https://example.test/legacy-live-compare"
    assert latest["Legacy Live Compare Matrix"]["url"] == "https://example.test/legacy-live-compare-matrix"
    assert latest["Remote Smoke Local"]["status"] == "completed"
    assert latest["Remote Smoke Matrix"]["created_at"] == "2026-03-29T00:20:00Z"
    assert latest["Remote Smoke Recovery"]["url"] == "https://example.test/remote-smoke-recovery"
    assert latest["Remote Smoke Recovery Matrix"]["conclusion"] == "success"
    summary = payload["summary"]
    assert "local_env_ready" in summary
    assert "local_workflow_ready" in summary
    assert "remote_happy_ready" in summary
    assert "remote_recovery_ready" in summary
    assert "legacy_live_ready" in summary
    assert "legacy_live_matrix_ready" in summary
    assert "real_endpoint_matrix_ready" in summary
    assert summary["real_endpoint_latest_success"] is False
    assert summary["real_endpoint_matrix_latest_success"] is False
    assert summary["real_endpoint_latest_result"] == "skipped"
    assert summary["real_endpoint_matrix_latest_result"] == "skipped"
    assert summary["legacy_live_latest_success"] is False
    assert summary["legacy_live_matrix_latest_success"] is False
    assert summary["legacy_live_latest_result"] == "skipped"
    assert summary["legacy_live_matrix_latest_result"] == "skipped"
    assert summary["remote_happy_latest_success"] is True
    assert summary["remote_recovery_latest_success"] is True
    assert summary["remote_happy_latest_result"] == "executed"
    assert summary["remote_recovery_latest_result"] == "executed"
    assert summary["remote_happy_latest_source"] == "Remote Smoke Matrix"
    assert summary["remote_recovery_latest_source"] == "Remote Smoke Recovery Matrix"
    assert summary["real_endpoint_latest_url"] == "https://example.test/smoke-s3"
    assert summary["real_endpoint_matrix_latest_url"] == "https://example.test/smoke-s3-matrix"
    assert summary["legacy_live_latest_url"] == "https://example.test/legacy-live-compare"
    assert summary["legacy_live_matrix_latest_url"] == "https://example.test/legacy-live-compare-matrix"
    assert summary["real_endpoint_latest_source"] == "Smoke S3"
    assert summary["real_endpoint_matrix_latest_source"] == "Smoke S3 Matrix"
    assert summary["legacy_live_latest_source"] == "Legacy Live Compare"
    assert summary["legacy_live_matrix_latest_source"] == "Legacy Live Compare Matrix"
    assert summary["remote_happy_latest_url"] == "https://example.test/remote-smoke-matrix"
    assert summary["remote_recovery_latest_url"] == "https://example.test/remote-smoke-recovery-matrix"
    assert summary["real_endpoint_latest_created_at"] == "2026-03-29T00:05:00Z"
    assert summary["real_endpoint_matrix_latest_created_at"] == "2026-03-29T00:06:00Z"
    assert summary["legacy_live_latest_created_at"] == "2026-03-29T00:07:00Z"
    assert summary["legacy_live_matrix_latest_created_at"] == "2026-03-29T00:08:00Z"
    assert summary["remote_happy_latest_created_at"] == "2026-03-29T00:20:00Z"
    assert summary["remote_recovery_latest_created_at"] == "2026-03-29T00:40:00Z"
    assert "ready" in summary


def test_smoke_ready_text_reports_remote_categories():
    proc = run_helper()
    text = proc.stdout
    assert "Smoke Local" in text
    assert "Smoke S3" in text
    assert "Smoke S3 Matrix" in text
    assert "Legacy Live Compare" in text
    assert "Legacy Live Compare Matrix" in text
    assert "Remote Smoke Local" in text
    assert "Remote Smoke Matrix" in text
    assert "Remote Smoke Recovery" in text
    assert "Remote Smoke Recovery Matrix" in text
    assert "Latest Runs" in text
    assert "completed/success" in text
    assert "Local Env Ready" in text
    assert "Local Workflow Ready" in text
    assert "Remote Happy Ready" in text
    assert "Remote Recovery Ready" in text
    assert "Legacy Live Ready" in text
    assert "Legacy Live Matrix Ready" in text
    assert "Real Endpoint Latest Success" in text
    assert "Real Endpoint Matrix Ready" in text
    assert "Real Endpoint Matrix Latest Success" in text
    assert "Real Endpoint Latest Result" in text
    assert "Real Endpoint Matrix Latest Result" in text
    assert "Legacy Live Latest Success" in text
    assert "Legacy Live Matrix Latest Success" in text
    assert "Legacy Live Latest Result" in text
    assert "Legacy Live Matrix Latest Result" in text
    assert "Remote Happy Latest Result" in text
    assert "Remote Recovery Latest Result" in text
    assert "Remote Happy Latest Source" in text
    assert "Remote Recovery Latest Source" in text
    assert "Real Endpoint Latest URL" in text
    assert "Real Endpoint Matrix Latest URL" in text
    assert "Legacy Live Latest URL" in text
    assert "Legacy Live Matrix Latest URL" in text
    assert "Real Endpoint Latest Source" in text
    assert "Real Endpoint Matrix Latest Source" in text
    assert "Legacy Live Latest Source" in text
    assert "Legacy Live Matrix Latest Source" in text
    assert "Remote Happy Latest URL" in text
    assert "Remote Recovery Latest URL" in text
    assert "Real Endpoint Latest Created At" in text
    assert "Real Endpoint Matrix Latest Created At" in text
    assert "Legacy Live Latest Created At" in text
    assert "Legacy Live Matrix Latest Created At" in text
    assert "Remote Happy Latest Created At" in text
    assert "Remote Recovery Latest Created At" in text
    assert "skipped" in text
    assert "Remote Happy Latest Success" in text
    assert "Remote Recovery Latest Success" in text
