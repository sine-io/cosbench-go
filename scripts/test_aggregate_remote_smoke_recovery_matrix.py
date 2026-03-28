import importlib.util
import json
from pathlib import Path


def load_module():
    path = Path("scripts/aggregate_remote_smoke_recovery_matrix.py")
    spec = importlib.util.spec_from_file_location("aggregate_remote_smoke_recovery_matrix", path)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


def write_summary(path: Path, payload: dict):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload) + "\n", encoding="utf-8")


def test_aggregate_recovery_rows_reports_present_and_missing(tmp_path):
    module = load_module()
    root = tmp_path / "downloads"
    write_summary(
        root / "remote-smoke-recovery-s3-summary" / "summary.json",
        {
            "backend": "s3",
            "scenario": "recovery",
            "overall": "pass",
            "job_status": "succeeded",
            "recovery_observed": True,
            "reclaimed_units": 1,
            "drivers_participated": 2,
        },
    )
    rows = module.aggregate_rows(root, [("s3", "recovery"), ("sio", "recovery")])
    assert rows[0]["status"] == "present"
    assert rows[0]["summary"]["backend"] == "s3"
    assert rows[1]["status"] == "missing"
    assert rows[1]["backend"] == "sio"
    assert rows[1]["scenario"] == "recovery"


def test_render_markdown_includes_recovery_statuses():
    module = load_module()
    rows = [
        {
            "backend": "s3",
            "scenario": "recovery",
            "status": "present",
            "summary": {
                "overall": "pass",
                "job_status": "succeeded",
                "recovery_observed": True,
                "reclaimed_units": 1,
                "drivers_participated": 2,
            },
        },
        {
            "backend": "sio",
            "scenario": "recovery",
            "status": "missing",
        },
    ]
    markdown = module.render_markdown(rows)
    assert "s3" in markdown
    assert "recovery" in markdown
    assert "pass" in markdown
    assert "sio" in markdown
    assert "missing" in markdown
