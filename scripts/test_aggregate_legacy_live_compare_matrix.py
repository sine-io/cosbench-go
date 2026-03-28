import json
import subprocess
from pathlib import Path


def write_summary(path: Path, payload: dict):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")


def test_aggregate_legacy_live_compare_matrix(tmp_path: Path):
    root = tmp_path / "downloads"
    write_summary(root / "legacy-live-compare-s3" / "summary.json", {"operations": 12})
    write_summary(root / "legacy-live-compare-s3" / "result.json", {"result": "executed"})
    write_summary(root / "legacy-live-compare-sio" / "summary.json", {"status": "skipped", "reason": "missing secrets"})
    write_summary(root / "legacy-live-compare-sio" / "result.json", {"result": "skipped", "reason": "missing secrets"})

    output = tmp_path / "aggregate"
    subprocess.run(
        ["python3", "scripts/aggregate_legacy_live_compare_matrix.py", str(root), str(output)],
        check=True,
        text=True,
        capture_output=True,
    )

    payload = json.loads((output / "summary.json").read_text(encoding="utf-8"))
    assert payload["overall"] == "pass"
    assert payload["rows"][0]["backend"] == "s3"
    assert payload["rows"][0]["status"] == "executed"
    assert payload["rows"][1]["backend"] == "sio"
    assert payload["rows"][1]["status"] == "skipped"
    assert payload["rows"][1]["reason"] == "missing secrets"

    markdown = (output / "summary.md").read_text(encoding="utf-8")
    assert "# Legacy Live Compare Matrix Summary" in markdown
    assert "| s3 | executed |" in markdown
    assert "| sio | skipped |" in markdown
