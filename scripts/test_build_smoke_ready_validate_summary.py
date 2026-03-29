import json
import subprocess
from pathlib import Path


def test_build_smoke_ready_validate_summary_writes_summary_files(tmp_path):
    src = tmp_path / "smoke-ready-validate"
    dst = tmp_path / "smoke-ready-validate-summary"
    src.mkdir(parents=True)
    payload = {
        "valid": True,
        "schema_path": "docs/smoke-ready.schema.json",
        "schema_version": 1,
        "generated_at": "2026-03-30T00:00:00Z",
        "error": "",
    }
    (src / "validation.json").write_text(json.dumps(payload) + "\n", encoding="utf-8")

    subprocess.run(
        [
            "python3",
            "scripts/build_smoke_ready_validate_summary.py",
            str(src),
            str(dst),
        ],
        check=True,
        cwd=Path.cwd(),
    )

    assert json.loads((dst / "summary.json").read_text(encoding="utf-8")) == payload
    summary_md = (dst / "summary.md").read_text(encoding="utf-8")
    assert "# Smoke Ready Validate" in summary_md
    assert "Valid: `yes`" in summary_md
    assert "Schema Path: `docs/smoke-ready.schema.json`" in summary_md
