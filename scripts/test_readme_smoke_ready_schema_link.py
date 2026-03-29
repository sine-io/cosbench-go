from pathlib import Path


def test_readme_uses_repo_local_smoke_ready_schema_link():
    readme = Path("README.md").read_text(encoding="utf-8")
    assert "[docs/smoke-ready.schema.json](docs/smoke-ready.schema.json)" in readme
    assert ".worktrees/smoke-ready-json-schema" not in readme
