import importlib.util
import json
from pathlib import Path
from types import SimpleNamespace


def load_smoke_ready_module():
    spec = importlib.util.spec_from_file_location("smoke_ready_module", Path("scripts/smoke_ready.py"))
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


class FakeTempDir:
    def __init__(self, path):
        self.path = path

    def __enter__(self):
        self.path.mkdir(parents=True, exist_ok=True)
        return str(self.path)

    def __exit__(self, exc_type, exc, tb):
        return False


def test_load_schema_validation_details_prefers_summary_artifact(tmp_path, monkeypatch):
    smoke_ready = load_smoke_ready_module()
    payload = {
        "valid": True,
        "schema_path": "docs/smoke-ready.schema.json",
        "schema_version": 1,
        "generated_at": "2026-03-30T00:00:00Z",
        "error": "",
    }
    calls = []

    def fake_run(*args):
        calls.append(args)
        artifact_name = args[args.index("-n") + 1]
        download_dir = Path(args[args.index("-D") + 1])
        if artifact_name == "smoke-ready-validate-summary":
            download_dir.mkdir(parents=True, exist_ok=True)
            (download_dir / "summary.json").write_text(json.dumps(payload), encoding="utf-8")
            return SimpleNamespace(returncode=0, stdout="", stderr="")
        return SimpleNamespace(returncode=1, stdout="", stderr="unexpected artifact")

    monkeypatch.setattr(smoke_ready, "run", fake_run)
    monkeypatch.setattr(smoke_ready.tempfile, "TemporaryDirectory", lambda prefix: FakeTempDir(tmp_path / "download"))

    details, ok, error = smoke_ready.load_schema_validation_details(
        "sine-io/cosbench-go",
        {
            smoke_ready.SMOKE_READY_VALIDATE_WORKFLOW: {
                "database_id": 123,
                "status": "completed",
            }
        },
    )

    assert ok is True
    assert error == ""
    assert details[smoke_ready.SMOKE_READY_VALIDATE_WORKFLOW] == payload
    assert calls[0][calls[0].index("-n") + 1] == "smoke-ready-validate-summary"
