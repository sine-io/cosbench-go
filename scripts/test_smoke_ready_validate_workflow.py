from pathlib import Path


def test_smoke_ready_validate_workflow_shape():
    workflow = Path(".github/workflows/smoke-ready-validate.yml").read_text(encoding="utf-8")
    assert "workflow_dispatch:" in workflow
    assert "schedule:" in workflow
    assert '- cron: "7 4 * * *"' in workflow
    assert "GH_TOKEN: ${{ github.token }}" in workflow
    assert "make --no-print-directory smoke-ready-json" in workflow
    assert "make --no-print-directory smoke-ready-validate-json" in workflow
    assert ".artifacts/smoke-ready-validate/smoke-ready.json" in workflow
    assert ".artifacts/smoke-ready-validate/validation.json" in workflow
    assert "python3 ./scripts/build_smoke_ready_validate_summary.py .artifacts/smoke-ready-validate .artifacts/smoke-ready-validate-summary" in workflow
    assert "uses: actions/upload-artifact@v7.0.0" in workflow
    assert "name: smoke-ready-validate-output" in workflow
    assert "name: smoke-ready-validate-summary" in workflow
    assert 'cat .artifacts/smoke-ready-validate-summary/summary.md >> "$GITHUB_STEP_SUMMARY"' in workflow
