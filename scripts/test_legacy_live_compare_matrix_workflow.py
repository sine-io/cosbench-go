from pathlib import Path


def test_legacy_live_compare_matrix_workflow_shape():
    workflow = Path(".github/workflows/legacy-live-compare-matrix.yml").read_text(encoding="utf-8")
    assert "name: Legacy Live Compare Matrix" in workflow
    assert "workflow_dispatch:" in workflow
    assert "region:" in workflow
    assert "path_style:" in workflow
    assert "strategy:" in workflow
    assert "backend: s3" in workflow
    assert "fixture: testdata/legacy/s3-config-sample.xml" in workflow
    assert "backend: sio" in workflow
    assert "fixture: testdata/legacy/sio-config-sample.xml" in workflow
    assert "Check live compare credentials" in workflow
    assert "render_legacy_live_compare_workload.py" in workflow
    assert "summarize_legacy_live_compare.py" in workflow
    assert "actions/download-artifact@v8.0.1" in workflow
    assert "aggregate:" in workflow
    assert "aggregate_legacy_live_compare_matrix.py" in workflow
    assert "legacy-live-compare-${{ matrix.backend }}" in workflow
    assert "legacy-live-compare-matrix-aggregate" in workflow
