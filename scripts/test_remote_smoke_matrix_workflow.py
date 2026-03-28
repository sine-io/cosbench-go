import pathlib


def test_remote_smoke_matrix_workflow_shape():
    workflow = pathlib.Path(".github/workflows/remote-smoke-matrix.yml").read_text(encoding="utf-8")
    assert "schedule:" in workflow
    assert "workflow_dispatch:" in workflow
    assert "fail-fast: false" in workflow
    assert "- backend: s3" in workflow
    assert "scenario: single" in workflow
    assert "scenario: multistage" in workflow
    assert "- backend: sio" in workflow
    assert "SMOKE_REMOTE_LOCAL_BACKEND='${{ matrix.backend }}'" in workflow
    assert "SMOKE_REMOTE_LOCAL_SCENARIO='${{ matrix.scenario }}'" in workflow
    assert "remote-smoke-${{ matrix.backend }}-${{ matrix.scenario }}" in workflow
