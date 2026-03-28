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
    summary = payload["summary"]
    assert "local_env_ready" in summary
    assert "local_workflow_ready" in summary
    assert "remote_happy_ready" in summary
    assert "remote_recovery_ready" in summary
    assert "ready" in summary


def test_smoke_ready_text_reports_remote_categories():
    proc = run_helper()
    text = proc.stdout
    assert "Smoke Local" in text
    assert "Remote Smoke Local" in text
    assert "Remote Smoke Matrix" in text
    assert "Remote Smoke Recovery" in text
    assert "Remote Smoke Recovery Matrix" in text
    assert "Local Env Ready" in text
    assert "Local Workflow Ready" in text
    assert "Remote Happy Ready" in text
    assert "Remote Recovery Ready" in text
