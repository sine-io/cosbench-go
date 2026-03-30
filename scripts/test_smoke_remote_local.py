import json
import pathlib

import smoke_remote_local as smoke


def test_remote_smoke_fixture_has_two_workers():
    fixture = pathlib.Path("testdata/workloads/remote-smoke-s3-two-driver.xml")
    text = fixture.read_text(encoding="utf-8")
    assert 'workers="2"' in text
    assert 'storage type="s3"' in text
    assert 'operation type="write"' in text


def test_remote_sio_smoke_fixture_has_two_workers():
    fixture = pathlib.Path("testdata/workloads/remote-smoke-sio-two-driver.xml")
    text = fixture.read_text(encoding="utf-8")
    assert 'workers="2"' in text
    assert 'storage type="sio"' in text
    assert 'operation type="write"' in text


def test_remote_multistage_smoke_fixture_has_two_stages_and_two_workers():
    fixture = pathlib.Path("testdata/workloads/remote-smoke-s3-multistage-two-driver.xml")
    text = fixture.read_text(encoding="utf-8")
    assert text.count("<workstage ") == 2
    assert text.count('workers="2"') == 2
    assert 'storage type="s3"' in text
    assert 'name="stage-a"' in text
    assert 'name="stage-b"' in text


def test_remote_sio_multistage_smoke_fixture_has_two_stages_and_two_workers():
    fixture = pathlib.Path("testdata/workloads/remote-smoke-sio-multistage-two-driver.xml")
    text = fixture.read_text(encoding="utf-8")
    assert text.count("<workstage ") == 2
    assert text.count('workers="2"') == 2
    assert 'storage type="sio"' in text
    assert 'name="stage-a"' in text
    assert 'name="stage-b"' in text


def test_remote_recovery_smoke_fixture_has_delay_and_two_workers():
    fixture = pathlib.Path("testdata/workloads/remote-smoke-s3-recovery-two-driver.xml")
    text = fixture.read_text(encoding="utf-8")
    assert text.count("<workstage ") == 1
    assert 'workers="2"' in text
    assert 'storage type="s3"' in text
    assert 'operation type="delay"' in text
    assert 'duration=45s' in text


def test_remote_sio_recovery_smoke_fixture_has_delay_and_two_workers():
    fixture = pathlib.Path("testdata/workloads/remote-smoke-sio-recovery-two-driver.xml")
    text = fixture.read_text(encoding="utf-8")
    assert text.count("<workstage ") == 1
    assert 'workers="2"' in text
    assert 'storage type="sio"' in text
    assert 'operation type="delay"' in text
    assert 'duration=45s' in text


def test_build_summary_json_shape():
    summary = smoke.build_summary(
        backend="s3",
        scenario="single",
        controller_url="http://127.0.0.1:19088",
        driver_urls=["http://127.0.0.1:18081", "http://127.0.0.1:18082"],
        job_id="job-1",
        job_status="succeeded",
        drivers_seen=2,
        units_claimed=2,
        drivers_participated=2,
        operation_count=2,
        byte_count=2000,
        stage_names=["main"],
        stages_seen=1,
        recovery_observed=None,
        reclaimed_units=None,
        checks={"drivers_healthy": "pass", "units_distributed": "pass"},
    )
    assert summary["controller_url"] == "http://127.0.0.1:19088"
    assert summary["driver_urls"] == ["http://127.0.0.1:18081", "http://127.0.0.1:18082"]
    assert summary["job_id"] == "job-1"
    assert summary["backend"] == "s3"
    assert summary["scenario"] == "single"
    assert summary["stage_names"] == ["main"]
    assert summary["stages_seen"] == 1
    assert summary["overall"] == "pass"
    json.dumps(summary)


def test_build_summary_can_include_recovery_fields():
    summary = smoke.build_summary(
        backend="s3",
        scenario="recovery",
        controller_url="http://127.0.0.1:19088",
        driver_urls=["http://127.0.0.1:18081", "http://127.0.0.1:18082"],
        job_id="job-1",
        job_status="succeeded",
        drivers_seen=2,
        units_claimed=2,
        drivers_participated=2,
        operation_count=2,
        byte_count=0,
        stage_names=["main"],
        stages_seen=1,
        recovery_observed=True,
        reclaimed_units=1,
        lease_expiry_event_observed=True,
        driver_unhealthy_event_observed=True,
        checks={
            "recovery_observed": "pass",
            "lease_expiry_event": "pass",
            "driver_unhealthy_event": "pass",
        },
    )
    assert summary["recovery_observed"] is True
    assert summary["reclaimed_units"] == 1
    assert summary["lease_expiry_event_observed"] is True
    assert summary["driver_unhealthy_event_observed"] is True
    json.dumps(summary)


def test_recovery_event_flags_detects_expected_events():
    events = [
        {"message": "mission lease expired"},
        {"message": "driver driver-1 marked unhealthy by heartbeat timeout"},
    ]
    flags = smoke.recovery_event_flags(events, "driver-1")
    assert flags == {
        "lease_expiry_event_observed": True,
        "driver_unhealthy_event_observed": True,
    }


def test_controller_server_cmd_omits_driver_heartbeat_timeout_for_single_scenario():
    cmd = smoke.controller_server_cmd(
        controller_port=19088,
        controller_data=pathlib.Path("/tmp/controller-data"),
        shared_token="remote-smoke-token",
        scenario="single",
    )
    assert "-driver-heartbeat-timeout" not in cmd


def test_controller_server_cmd_includes_short_driver_heartbeat_timeout_for_recovery():
    cmd = smoke.controller_server_cmd(
        controller_port=19088,
        controller_data=pathlib.Path("/tmp/controller-data"),
        shared_token="remote-smoke-token",
        scenario="recovery",
    )
    assert "-driver-heartbeat-timeout" in cmd
    assert "2s" in cmd


def test_build_failure_summary_for_missing_process():
    summary = smoke.build_failure_summary("controller", "controller failed to start")
    assert summary["overall"] == "fail"
    assert summary["failed_at"] == "controller"
    assert "controller failed to start" in summary["error"]


def test_fixture_path_selection_by_backend():
    assert smoke.fixture_for_backend("s3").name == "remote-smoke-s3-two-driver.xml"
    assert smoke.fixture_for_backend("sio").name == "remote-smoke-sio-two-driver.xml"


def test_fixture_path_selection_by_backend_and_scenario():
    assert smoke.fixture_for_selection("s3", "single").name == "remote-smoke-s3-two-driver.xml"
    assert smoke.fixture_for_selection("s3", "multistage").name == "remote-smoke-s3-multistage-two-driver.xml"
    assert smoke.fixture_for_selection("sio", "multistage").name == "remote-smoke-sio-multistage-two-driver.xml"
    assert smoke.fixture_for_selection("s3", "recovery").name == "remote-smoke-s3-recovery-two-driver.xml"
    assert smoke.fixture_for_selection("sio", "recovery").name == "remote-smoke-sio-recovery-two-driver.xml"


def test_unknown_backend_is_rejected():
    try:
        smoke.fixture_for_backend("swift")
    except ValueError as err:
        assert "unsupported backend" in str(err)
    else:
        raise AssertionError("expected unsupported backend error")
