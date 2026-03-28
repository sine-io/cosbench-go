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
    assert "aggregate:" in workflow
    assert "needs: remote_smoke_matrix" in workflow
    assert "pattern: remote-smoke-*" in workflow
    assert "aggregate_remote_smoke_matrix.py" in workflow
    aggregate_section = workflow.split("aggregate:", 1)[1]
    assert "FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true" in aggregate_section
    assert "uses: actions/upload-artifact@v6.0.0" in aggregate_section
    assert "remote-smoke-matrix-aggregate" in aggregate_section
    assert ".artifacts/remote-smoke-matrix-aggregate" in aggregate_section
