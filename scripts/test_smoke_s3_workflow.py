from pathlib import Path


def test_smoke_s3_workflow_shape():
    workflow = Path(".github/workflows/smoke-s3.yml").read_text(encoding="utf-8")
    assert "workflow_dispatch:" in workflow
    assert "backend:" in workflow
    assert "region:" in workflow
    assert "path_style:" in workflow
    assert "bucket_prefix:" in workflow
    assert "COSBENCH_SMOKE_ENDPOINT" in workflow
    assert "COSBENCH_SMOKE_ACCESS_KEY" in workflow
    assert "COSBENCH_SMOKE_SECRET_KEY" in workflow
    assert "GO=go make smoke-s3" in workflow
    assert "tee smoke-s3-output.txt" in workflow
    assert "summarize_smoke_s3_output.py" in workflow
    assert ".artifacts/smoke-s3-summary/summary.json" in workflow
    assert "uses: actions/upload-artifact@v7.0.0" in workflow
    assert "name: smoke-s3-output" in workflow
    assert 'cat smoke-s3-output.txt >> "$GITHUB_STEP_SUMMARY"' in workflow
