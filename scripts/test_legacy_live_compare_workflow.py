from pathlib import Path


def test_legacy_live_compare_workflow_shape():
    workflow = Path(".github/workflows/legacy-live-compare.yml").read_text(encoding="utf-8")
    assert "workflow_dispatch:" in workflow
    assert "fixture:" in workflow
    assert "backend:" in workflow
    assert "region:" in workflow
    assert "path_style:" in workflow
    assert "render_legacy_live_compare_workload.py" in workflow
    assert "go run ./cmd/cosbench-go" in workflow
    assert "-json -quiet -summary-file" in workflow
    assert "legacy-live-compare-output" in workflow
