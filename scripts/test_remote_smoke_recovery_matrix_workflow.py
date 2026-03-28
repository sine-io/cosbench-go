from pathlib import Path


def test_remote_smoke_recovery_matrix_workflow_shape():
    workflow = Path(".github/workflows/remote-smoke-recovery-matrix.yml").read_text(encoding="utf-8")
    assert "workflow_dispatch:" in workflow
    assert "schedule:" in workflow
    assert "fail-fast: false" in workflow
    assert "- backend: s3" in workflow
    assert "- backend: sio" in workflow
    assert "scenario: recovery" in workflow
    assert "SMOKE_REMOTE_LOCAL_BACKEND='${{ matrix.backend }}'" in workflow
    assert "SMOKE_REMOTE_LOCAL_SCENARIO='${{ matrix.scenario }}'" in workflow
    assert "aggregate:" in workflow
    assert "aggregate_remote_smoke_recovery_matrix.py" in workflow
