import json
import pathlib

import smoke_remote_local as smoke


def test_remote_smoke_fixture_has_two_workers():
    fixture = pathlib.Path("testdata/workloads/remote-smoke-s3-two-driver.xml")
    text = fixture.read_text(encoding="utf-8")
    assert 'workers="2"' in text
    assert 'storage type="s3"' in text
    assert 'operation type="write"' in text


def test_build_summary_json_shape():
    summary = smoke.build_summary(
        controller_url="http://127.0.0.1:19088",
        driver_urls=["http://127.0.0.1:18081", "http://127.0.0.1:18082"],
        job_id="job-1",
        job_status="succeeded",
        drivers_seen=2,
        units_claimed=2,
        drivers_participated=2,
        operation_count=2,
        byte_count=2000,
        checks={"drivers_healthy": "pass", "units_distributed": "pass"},
    )
    assert summary["controller_url"] == "http://127.0.0.1:19088"
    assert summary["driver_urls"] == ["http://127.0.0.1:18081", "http://127.0.0.1:18082"]
    assert summary["job_id"] == "job-1"
    assert summary["overall"] == "pass"
    json.dumps(summary)


def test_build_failure_summary_for_missing_process():
    summary = smoke.build_failure_summary("controller", "controller failed to start")
    assert summary["overall"] == "fail"
    assert summary["failed_at"] == "controller"
    assert "controller failed to start" in summary["error"]
