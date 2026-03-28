import json
import os
import subprocess


def run_helper(*args, env_overrides=None):
    env = os.environ.copy()
    env["SMOKE_READY_MOCK_REPO_SECRETS"] = "COSBENCH_SMOKE_ENDPOINT,COSBENCH_SMOKE_ACCESS_KEY,COSBENCH_SMOKE_SECRET_KEY"
    env["SMOKE_READY_MOCK_WORKFLOWS"] = ",".join(
        [
            "Smoke Local",
            "Remote Smoke Local",
            "Remote Smoke Matrix",
            "Remote Smoke Recovery",
            "Remote Smoke Recovery Matrix",
        ]
    )
    env["SMOKE_READY_MOCK_WORKFLOW_RUNS_JSON"] = json.dumps(
        {
            "Smoke Local": {
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:00:00Z",
                "url": "https://example.test/smoke-local",
            },
            "Remote Smoke Local": {
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:10:00Z",
                "url": "https://example.test/remote-smoke-local",
            },
            "Remote Smoke Matrix": {
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:20:00Z",
                "url": "https://example.test/remote-smoke-matrix",
            },
            "Remote Smoke Recovery": {
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:30:00Z",
                "url": "https://example.test/remote-smoke-recovery",
            },
            "Remote Smoke Recovery Matrix": {
                "status": "completed",
                "conclusion": "success",
                "created_at": "2026-03-29T00:40:00Z",
                "url": "https://example.test/remote-smoke-recovery-matrix",
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
    assert present["Remote Smoke Local"] is True
    assert present["Remote Smoke Matrix"] is True
    assert present["Remote Smoke Recovery"] is True
    assert present["Remote Smoke Recovery Matrix"] is True
    latest = payload["workflows"]["latest"]
    assert latest["Smoke Local"]["conclusion"] == "success"
    assert latest["Remote Smoke Local"]["status"] == "completed"
    assert latest["Remote Smoke Matrix"]["created_at"] == "2026-03-29T00:20:00Z"
    assert latest["Remote Smoke Recovery"]["url"] == "https://example.test/remote-smoke-recovery"
    assert latest["Remote Smoke Recovery Matrix"]["conclusion"] == "success"
    summary = payload["summary"]
    assert "local_env_ready" in summary
    assert "local_workflow_ready" in summary
    assert "remote_happy_ready" in summary
    assert "remote_recovery_ready" in summary
    assert "remote_happy_latest_success" in summary
    assert "remote_recovery_latest_success" in summary
    assert "ready" in summary


def test_smoke_ready_text_reports_remote_categories():
    proc = run_helper()
    text = proc.stdout
    assert "Smoke Local" in text
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
    assert "Remote Happy Latest Success" in text
    assert "Remote Recovery Latest Success" in text
