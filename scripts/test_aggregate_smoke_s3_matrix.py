import importlib.util
import json
from pathlib import Path


def load_module():
    path = Path("scripts/aggregate_smoke_s3_matrix.py")
    spec = importlib.util.spec_from_file_location("aggregate_smoke_s3_matrix", path)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


def write_output(path: Path, text: str):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(text, encoding="utf-8")


def write_summary(path: Path, payload: dict):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")


def test_aggregate_smoke_s3_rows_reports_present_and_missing(tmp_path):
    module = load_module()
    root = tmp_path / "downloads"
    write_output(root / "smoke-s3-s3" / "smoke-s3-output.txt", "s3 success")
    write_summary(root / "smoke-s3-s3" / "summary.json", {"result": "executed"})
    rows = module.aggregate_rows(root, ["s3", "sio"])
    assert rows[0]["status"] == "executed"
    assert rows[0]["backend"] == "s3"
    assert rows[1]["status"] == "missing"
    assert rows[1]["backend"] == "sio"


def test_render_markdown_includes_row_statuses():
    module = load_module()
    rows = [
        {"backend": "s3", "status": "executed", "output": "s3 success"},
        {"backend": "sio", "status": "missing"},
    ]
    markdown = module.render_markdown(rows)
    assert "s3" in markdown
    assert "executed" in markdown
    assert "sio" in markdown
    assert "missing" in markdown
