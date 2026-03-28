import json
import subprocess
from pathlib import Path


def test_build_remote_smoke_recovery_summary_copies_summary_files(tmp_path):
    src = tmp_path / "remote-smoke"
    dst = tmp_path / "remote-smoke-recovery-summary"
    src.mkdir(parents=True)
    (src / "summary.json").write_text(json.dumps({"overall": "pass"}) + "\n", encoding="utf-8")
    (src / "summary.md").write_text("# Remote Smoke\n", encoding="utf-8")

    subprocess.run(
        [
            "python3",
            "scripts/build_remote_smoke_recovery_summary.py",
            str(src),
            str(dst),
        ],
        check=True,
        cwd=Path.cwd(),
    )

    assert json.loads((dst / "summary.json").read_text(encoding="utf-8"))["overall"] == "pass"
    assert (dst / "summary.md").read_text(encoding="utf-8") == "# Remote Smoke\n"
