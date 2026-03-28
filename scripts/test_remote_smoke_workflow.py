import pathlib


def test_remote_smoke_workflow_accepts_scenario_and_threads_it_to_helper():
    workflow = pathlib.Path(".github/workflows/remote-smoke-local.yml").read_text(encoding="utf-8")
    assert "scenario:" in workflow
    assert "default: \"single\"" in workflow
    assert "SMOKE_REMOTE_LOCAL_SCENARIO='${{ inputs.scenario }}'" in workflow
    assert "SMOKE_REMOTE_LOCAL_BACKEND='${{ inputs.backend }}'" in workflow
