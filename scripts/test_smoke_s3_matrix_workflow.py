from pathlib import Path


def test_smoke_s3_matrix_workflow_shape():
    workflow = Path(".github/workflows/smoke-s3-matrix.yml").read_text(encoding="utf-8")
    assert "workflow_dispatch:" in workflow
    assert "fail-fast: false" in workflow
    assert "- backend: s3" in workflow
    assert "- backend: sio" in workflow
    assert "GO=go make smoke-s3" in workflow
    assert "aggregate:" in workflow
    assert "aggregate_smoke_s3_matrix.py" in workflow
    assert "smoke-s3-${{ matrix.backend }}" in workflow
