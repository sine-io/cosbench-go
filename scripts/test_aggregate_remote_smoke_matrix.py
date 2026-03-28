import importlib.util
import json
from pathlib import Path


def load_module():
    path = Path("scripts/aggregate_remote_smoke_matrix.py")
    spec = importlib.util.spec_from_file_location("aggregate_remote_smoke_matrix", path)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


def write_summary(path: Path, payload: dict):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload) + "\n", encoding="utf-8")


def test_aggregate_rows_reports_present_and_missing(tmp_path):
    module = load_module()
    root = tmp_path / "downloads"
    write_summary(
        root / "remote-smoke-s3-single" / "summary.json",
        {
            "backend": "s3",
            "scenario": "single",
            "overall": "pass",
            "job_status": "succeeded",
            "drivers_seen": 2,
            "units_claimed": 2,
            "stages_seen": 1,
        },
    )
    rows = module.aggregate_rows(root, [("s3", "single"), ("sio", "multistage")])
    assert rows[0]["status"] == "present"
    assert rows[0]["summary"]["backend"] == "s3"
    assert rows[1]["status"] == "missing"
    assert rows[1]["backend"] == "sio"
    assert rows[1]["scenario"] == "multistage"


def test_render_markdown_includes_row_statuses(tmp_path):
    module = load_module()
    rows = [
        {
            "backend": "s3",
            "scenario": "single",
            "status": "present",
            "summary": {
                "overall": "pass",
                "job_status": "succeeded",
                "drivers_seen": 2,
                "units_claimed": 2,
                "stages_seen": 1,
            },
        },
        {
            "backend": "sio",
            "scenario": "multistage",
            "status": "missing",
        },
    ]
    markdown = module.render_markdown(rows)
    assert "s3" in markdown
    assert "single" in markdown
    assert "pass" in markdown
    assert "sio" in markdown
    assert "multistage" in markdown
    assert "missing" in markdown
