from pathlib import Path


def test_remote_smoke_recovery_workflow_shape():
    workflow = Path(".github/workflows/remote-smoke-recovery.yml").read_text(encoding="utf-8")
    assert "workflow_dispatch:" in workflow
    assert "backend:" in workflow
    assert 'default: "s3"' in workflow
    assert "SMOKE_REMOTE_LOCAL_BACKEND='${{ inputs.backend }}'" in workflow
    assert "SMOKE_REMOTE_LOCAL_SCENARIO=recovery" in workflow
    assert "make --no-print-directory smoke-remote-local" in workflow
    assert "build_remote_smoke_recovery_summary.py" in workflow
    assert "uses: actions/upload-artifact@v7.0.0" in workflow
    assert "path: .artifacts/remote-smoke" in workflow
    assert "remote-smoke-recovery-summary" in workflow
    assert "path: .artifacts/remote-smoke-recovery-summary" in workflow
    assert 'cat .artifacts/remote-smoke/summary.md >> "$GITHUB_STEP_SUMMARY"' in workflow
